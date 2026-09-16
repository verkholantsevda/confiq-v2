package audit

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(log *AuditLog) error {
	return r.db.Create(log).Error
}

func (r *Repository) ListRecent(limit int) ([]AuditLog, error) {
	var logs []AuditLog

	err := r.db.
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).
		Error

	return logs, err
}
