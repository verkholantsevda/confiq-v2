package configtypes

import "time"

type ConfigType struct {
	ID uint `gorm:"primaryKey"`

	Name              string `gorm:"size:100;uniqueIndex;not null"`
	Description       string `gorm:"type:text"`
	ConfigTemplate    string `gorm:"type:text;not null"`
	UsageInstructions string `gorm:"type:text"`
	ClientLinks       string `gorm:"type:text"` // JSON

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
}
