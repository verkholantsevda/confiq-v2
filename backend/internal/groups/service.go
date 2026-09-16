package groups

import (
	"confiq/internal/audit"
	"confiq/internal/endpoints"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type Service struct {
	repo      *Repository
	endpoints *endpoints.Repository
	audit     *audit.Service
}

func NewService(
	repo *Repository,
	endpoints *endpoints.Repository,
	auditService *audit.Service,
) *Service {
	return &Service{
		repo:      repo,
		endpoints: endpoints,
		audit:     auditService,
	}
}

func (s *Service) List() ([]Group, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id uint) (*Group, error) {
	group, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *Service) Create(actorID uint, req CreateGroupRequest) (*Group, error) {

	existing, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrGroupAlreadyExists
	}

	group := &Group{
		Name:        req.Name,
		Description: req.Description,
	}

	var eps []endpoints.Endpoint

	for _, id := range req.EndpointIDs {
		ep, err := s.endpoints.GetByID(id)
		if err != nil || ep == nil {
			return nil, endpoints.ErrEndpointNotFound
		}

		eps = append(eps, *ep)
	}

	group.Endpoints = eps

	if err := s.repo.Create(group); err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"group.created",
		"group",
		&group.ID,
		"Создана группа "+group.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return group, nil
}

func (s *Service) Update(actorID uint, id uint, req UpdateGroupRequest) (*Group, error) {

	group, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGroupNotFound
		}

		return nil, err
	}

	// Если меняется имя — проверяем уникальность
	if group.Name != req.Name {
		existing, err := s.repo.GetByName(req.Name)
		if err != nil {
			return nil, err
		}

		if existing != nil {
			return nil, ErrGroupAlreadyExists
		}
	}

	group.Name = req.Name
	group.Description = req.Description

	var eps []endpoints.Endpoint

	for _, id := range req.EndpointIDs {
		ep, err := s.endpoints.GetByID(id)
		if err != nil || ep == nil {
			return nil, endpoints.ErrEndpointNotFound
		}

		eps = append(eps, *ep)
	}

	group.Endpoints = eps

	if err := s.repo.Update(group); err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"group.updated",
		"group",
		&group.ID,
		"Изменена группа "+group.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return group, nil
}

func (s *Service) Delete(actorID uint, id uint) error {

	group, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrGroupNotFound
		}

		return err
	}

	if err := s.repo.Delete(group.ID); err != nil {
		return err
	}

	if err := s.audit.Log(
		&actorID,
		"group.deleted",
		"group",
		&group.ID,
		"Удалена группа "+group.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return nil
}
