package groups

import (
	"confiq/internal/endpoints"
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	repo      *Repository
	endpoints *endpoints.Repository
}

func NewService(
	repo *Repository,
	endpoints *endpoints.Repository,
) *Service {
	return &Service{
		repo:      repo,
		endpoints: endpoints,
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

func (s *Service) Create(req CreateGroupRequest) (*Group, error) {

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

	return group, nil
}

func (s *Service) Update(id uint, req UpdateGroupRequest) (*Group, error) {

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

	return group, nil
}

func (s *Service) Delete(id uint) error {

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

	return nil
}
