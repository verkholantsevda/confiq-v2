package relations

import (
	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
)

type EndpointConfigType struct {
	EndpointID uint               `gorm:"primaryKey"`
	Endpoint   endpoints.Endpoint `gorm:"foreignKey:EndpointID"`

	ConfigTypeID uint                   `gorm:"primaryKey"`
	ConfigType   configtypes.ConfigType `gorm:"foreignKey:ConfigTypeID"`
}
