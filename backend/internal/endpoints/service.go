package endpoints

import (
	"confiq/internal/audit"
	"confiq/internal/configtypes"
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

func (s *Service) List() ([]Endpoint, error) {
	return s.repo.List()
}

func (s *Service) ListForUser(userID uint) ([]Endpoint, error) {
	return s.repo.ListForUser(userID)
}

func (s *Service) ListConfigTypes(userID, endpointID uint) ([]configtypes.ConfigType, error) {

	if err := s.repo.ValidateUserEndpoint(userID, endpointID); err != nil {
		return nil, ErrEndpointNotFound
	}

	return s.repo.ListConfigTypes(endpointID)
}

func (s *Service) GetByID(id uint) (*Endpoint, error) {
	endpoint, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (s *Service) Create(actorID uint, req CreateEndpointRequest) (*Endpoint, error) {
	endpoint := &Endpoint{
		Name:    req.Name,
		Address: req.Address,
		Port:    req.Port,
	}

	if err := s.repo.Create(endpoint); err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceGroups(endpoint.ID, req.GroupIDs); err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceConfigTypes(endpoint.ID, req.ConfigTypeIDs); err != nil {
		return nil, err
	}

	endpoint, err := s.repo.GetByID(endpoint.ID)
	if err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"endpoint.created",
		"endpoint",
		&endpoint.ID,
		"Создан endpoint "+endpoint.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return endpoint, nil
}

func (s *Service) Update(actorID uint, id uint, req UpdateEndpointRequest) (*Endpoint, error) {
	endpoint, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if endpoint == nil {
		return nil, ErrEndpointNotFound
	}

	endpoint.Name = req.Name
	endpoint.Address = req.Address
	endpoint.Port = req.Port

	if err := s.repo.Update(endpoint); err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceGroups(endpoint.ID, req.GroupIDs); err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceConfigTypes(endpoint.ID, req.ConfigTypeIDs); err != nil {
		return nil, err
	}

	endpoint, err = s.repo.GetByID(endpoint.ID)
	if err != nil {
		return nil, err
	}

	if err := s.audit.Log(
		&actorID,
		"endpoint.updated",
		"endpoint",
		&endpoint.ID,
		"Изменен endpoint "+endpoint.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return endpoint, nil
}

func (s *Service) Delete(actorID uint, id uint) error {
	endpoint, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if endpoint == nil {
		return ErrEndpointNotFound
	}

	if err := s.repo.Delete(endpoint.ID); err != nil {
		return err
	}

	if err := s.audit.Log(
		&actorID,
		"endpoint.deleted",
		"endpoint",
		&endpoint.ID,
		"Удален endpoint "+endpoint.Name,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}

	return nil
}

func (s *Service) AddGroupEndpoint(groupID, endpointID uint) error {
	if err := s.repo.AddGroupEndpoint(groupID, endpointID); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveGroupEndpoint(groupID, endpointID uint) error {

	if err := s.repo.RemoveGroupEndpoint(groupID, endpointID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListGroupsEndpoints() ([]GroupEndpoint, error) {
	groupEndpoints, err := s.repo.ListGroupsEndpoints()
	if err != nil {
		return nil, err
	}
	return groupEndpoints, nil
}
