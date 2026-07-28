package relations

import (
	"confiq/internal/endpoints"
	"confiq/internal/groups"
)

type GroupEndpoint struct {
	GroupID uint         `gorm:"primaryKey"`
	Group   groups.Group `gorm:"foreignKey:GroupID"`

	EndpointID uint               `gorm:"primaryKey"`
	Endpoint   endpoints.Endpoint `gorm:"foreignKey:EndpointID"`
}
