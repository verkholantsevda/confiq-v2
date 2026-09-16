package configs

import (
	"confiq/internal/audit"
	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
	"confiq/internal/logger"
	"confiq/internal/users"
	"confiq/internal/warp"
	"log/slog"
)

type Service struct {
	repo *Repository

	audit *audit.Service

	users       *users.Repository
	endpoints   *endpoints.Repository
	configTypes *configtypes.Repository
	generator   *warp.Generator
}

func NewService(
	repo *Repository,
	users *users.Repository,
	endpoints *endpoints.Repository,
	configTypes *configtypes.Repository,
	generator *warp.Generator,
	auditService *audit.Service,
) *Service {
	return &Service{
		repo: repo,

		audit: auditService,

		users:       users,
		endpoints:   endpoints,
		configTypes: configTypes,
		generator:   generator,
	}
}

func (s *Service) Create(actorID uint, req CreateConfigRequest) (*Config, error) {
	// Проверяем пользователя
	if _, err := s.users.GetByID(actorID); err != nil {
		return nil, users.ErrUserNotFound
	}

	// Проверяем endpoint
	endpoint, err := s.endpoints.GetByID(req.EndpointID)
	if err != nil {
		return nil, endpoints.ErrEndpointNotFound
	}

	// Проверяем тип конфигурации
	var configType *configtypes.ConfigType
	if req.ConfigTypeID != nil {
		configType, err = s.configTypes.GetByID(*req.ConfigTypeID)
		if err != nil {
			return nil, configtypes.ErrConfigTypeNotFound
		}
	}

	// Генерируем WARP-конфигурацию
	generated, err := s.generator.Generate(warp.GenerateRequest{
		EndpointAddress: endpoint.Address,
		EndpointPort:    endpoint.Port,
	})
	if err != nil {
		return nil, err
	}

	configContent := ""
	if configType != nil {
		configContent, err = warp.Render(
			configType.ConfigTemplate,
			warp.TemplateData{
				PrivateKey:    generated.PrivateKey,
				PublicKey:     generated.PublicKey,
				PeerPublicKey: generated.PeerPublicKey,
				ClientIPv4:    generated.ClientIPv4,
				ClientIPv6:    generated.ClientIPv6,
				Endpoint:      endpoint.Address,
				Port:          endpoint.Port,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	config := &Config{
		Name:         req.Name,
		UserID:       actorID,
		EndpointID:   req.EndpointID,
		ConfigTypeID: req.ConfigTypeID,

		CloudflareID:    generated.ID,
		CloudflareToken: generated.Token,

		PrivateKey:    generated.PrivateKey,
		PublicKey:     generated.PublicKey,
		PeerPublicKey: generated.PeerPublicKey,

		ClientIPv4: generated.ClientIPv4,
		ClientIPv6: generated.ClientIPv6,

		ConfigContent: configContent,
	}

	if err := s.repo.Create(config); err != nil {
		_ = s.generator.DeleteDevice(generated.ID, generated.Token) // если добавим такую обертку
		return nil, err
	}

	logger.ConfigCreated(config.UserID, config.ID, config.Name)
	if err := s.audit.Log(
		&actorID,
		"config.created",
		"config",
		&config.ID,
		"Создана конфигурация "+config.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
	return config, nil
}

func (s *Service) List(userID uint) ([]Config, error) {
	return s.repo.ListByUser(userID)
}

func (s *Service) ListAll() ([]Config, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id uint, userID uint, isAdmin bool) (*Config, error) {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrConfigNotFound
	}

	if !isAdmin && config.UserID != userID {
		return nil, ErrConfigNotFound
	}

	return config, nil
}

func (s *Service) Update(id uint, actorID uint, isAdmin bool, req UpdateConfigRequest) (*Config, error) {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrConfigNotFound
	}

	if !isAdmin && config.UserID != actorID {
		return nil, ErrConfigNotFound
	}

	config.Name = req.Name

	if err := s.repo.Update(config); err != nil {
		return nil, err
	}

	logger.ConfigUpdated(config.UserID, config.ID, config.Name)
	if err := s.audit.Log(
		&actorID,
		"config.updated",
		"config",
		&config.ID,
		"Изменена конфигурация "+config.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
	return config, nil
}

func (s *Service) Delete(id uint, actorID uint, isAdmin bool) error {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return ErrConfigNotFound
	}

	if !isAdmin && config.UserID != actorID {
		return ErrConfigNotFound
	}

	if err := s.generator.DeleteDevice(config.CloudflareID, config.CloudflareToken); err != nil {
		return err
	}

	if err := s.repo.Delete(config.ID); err != nil {
		return err
	}

	logger.ConfigDeleted(config.UserID, config.ID)
	if err := s.audit.Log(
		&actorID,
		"config.deleted",
		"config",
		&config.ID,
		"Удалена конфигурация "+config.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
	return nil
}
