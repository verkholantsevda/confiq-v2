package configs

import (
	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
	"confiq/internal/logger"
	"confiq/internal/users"
	"confiq/internal/warp"
)

type Service struct {
	repo *Repository

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
) *Service {
	return &Service{
		repo: repo,

		users:       users,
		endpoints:   endpoints,
		configTypes: configTypes,
		generator:   generator,
	}
}

func (s *Service) Create(userID uint, req CreateConfigRequest) (*Config, error) {
	// Проверяем пользователя
	if _, err := s.users.GetByID(userID); err != nil {
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
		UserID:       userID,
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

func (s *Service) Update(id uint, userID uint, isAdmin bool, req UpdateConfigRequest) (*Config, error) {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrConfigNotFound
	}

	if !isAdmin && config.UserID != userID {
		return nil, ErrConfigNotFound
	}

	config.Name = req.Name

	if err := s.repo.Update(config); err != nil {
		return nil, err
	}

	logger.ConfigUpdated(config.UserID, config.ID, config.Name)
	return config, nil
}

func (s *Service) Delete(id uint, userID uint, isAdmin bool) error {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return ErrConfigNotFound
	}

	if !isAdmin && config.UserID != userID {
		return ErrConfigNotFound
	}

	if err := s.generator.DeleteDevice(config.CloudflareID, config.CloudflareToken); err != nil {
		return err
	}

	if err := s.repo.Delete(config.ID); err != nil {
		return err
	}

	logger.ConfigDeleted(config.UserID, config.ID)
	return nil
}
