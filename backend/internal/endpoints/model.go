package endpoints

import (
	"confiq/internal/configtypes"
	"time"
)

type Endpoint struct {
	ID uint `gorm:"primaryKey"`

	Name        string                   `gorm:"size:100;not null"`
	Address     string                   `gorm:"size:255;not null"`
	Port        int                      `gorm:"not null"`
	ConfigTypes []configtypes.ConfigType `gorm:"many2many:endpoint_config_types;"`
	CreatedAt   time.Time
}

type GroupEndpoint struct {
	GroupID    uint `gorm:"primaryKey"`
	EndpointID uint `gorm:"primaryKey"`
}
