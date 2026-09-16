package configtypes

import (
	"confiq/internal/audit"
	"log/slog"
)

type Service struct {
	repo  *Repository
	audit *audit.Service
}

func NewService(repo *Repository, auditService *audit.Service) *Service {
	return &Service{
		repo:  repo,
		audit: auditService,
	}
}

func (s *Service) List() ([]ConfigType, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id uint) (*ConfigType, error) {
	configType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return configType, nil
}

func (s *Service) Create(actorID uint, req CreateConfigTypeRequest) (*ConfigType, error) {

	existing, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrConfigTypeAlreadyExists
	}

	configType := &ConfigType{
		Name:              req.Name,
		Description:       req.Description,
		ConfigTemplate:    req.ConfigTemplate,
		UsageInstructions: req.UsageInstructions,
		ClientLinks:       req.ClientLinks,
		IsActive:          req.IsActive,
	}

	if err := s.repo.Create(configType); err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"config_type.created",
		"config_type",
		&configType.ID,
		"Создан тип конфигурации "+configType.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return configType, nil
}

func (s *Service) Update(actorID uint, id uint, req UpdateConfigTypeRequest) (*ConfigType, error) {

	configType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrConfigTypeNotFound
	}

	if configType.Name != req.Name {
		existing, err := s.repo.GetByName(req.Name)
		if err != nil {
			return nil, err
		}

		if existing != nil {
			return nil, ErrConfigTypeAlreadyExists
		}
	}

	configType.Name = req.Name
	configType.Description = req.Description
	configType.ConfigTemplate = req.ConfigTemplate
	configType.UsageInstructions = req.UsageInstructions
	configType.ClientLinks = req.ClientLinks
	configType.IsActive = req.IsActive

	if err := s.repo.Update(configType); err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"config_type.updated",
		"config_type",
		&configType.ID,
		"Изменен тип конфигурации "+configType.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return configType, nil
}

func (s *Service) Delete(actorID uint, id uint) error {

	configType, err := s.repo.GetByID(id)
	if err != nil {
		return ErrConfigTypeNotFound
	}

	if err := s.repo.Delete(configType.ID); err != nil {
		return err
	}

	if err := s.audit.Log(
		&actorID,
		"config_type.deleted",
		"config_type",
		&configType.ID,
		"Удален тип конфигурации "+configType.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return nil
}
