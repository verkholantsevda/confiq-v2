package configs

import (
	"time"

	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
	"confiq/internal/users"
)

type Config struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"size:100"`

	UserID uint       `gorm:"not null"`
	User   users.User `gorm:"foreignKey:UserID"`

	EndpointID uint               `gorm:"not null"`
	Endpoint   endpoints.Endpoint `gorm:"foreignKey:EndpointID"`

	ConfigTypeID *uint
	ConfigType   *configtypes.ConfigType `gorm:"foreignKey:ConfigTypeID"`

	CloudflareID    string `gorm:"size:100;not null"`
	CloudflareToken string `gorm:"size:255;not null"`

	PrivateKey    string `gorm:"type:text;not null"`
	PublicKey     string `gorm:"type:text;not null"`
	PeerPublicKey string `gorm:"type:text;not null"`

	ClientIPv4 string `gorm:"size:50;not null"`
	ClientIPv6 string `gorm:"size:50;not null"`

	ConfigContent string `gorm:"type:text;not null"`

	CreatedAt time.Time
}
