package audit

import "time"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Log(
	userID *uint,
	action string,
	entityType string,
	entityID *uint,
	message string,
) error {
	log := &AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Message:    message,
		CreatedAt:  time.Now(),
	}

	return s.repo.Create(log)
}

func (s *Service) ListRecent(limit int) ([]AuditLog, error) {
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return s.repo.ListRecent(limit)
}
