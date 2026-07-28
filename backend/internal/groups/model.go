package groups

import (
	"confiq/internal/endpoints"
	"time"
)

type Group struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;uniqueIndex;not null"`
	Description string `gorm:"type:text"`

	Endpoints []endpoints.Endpoint `gorm:"many2many:group_endpoints;"`

	CreatedAt time.Time
}
