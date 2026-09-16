package audit

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     *uint     `gorm:"column:user_id" json:"user_id"`
	Action     string    `gorm:"size:100;not null" json:"action"`
	EntityType string    `gorm:"size:50;not null" json:"entity_type"`
	EntityID   *uint     `gorm:"column:entity_id" json:"entity_id"`
	Message    string    `gorm:"type:text" json:"message"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
}
