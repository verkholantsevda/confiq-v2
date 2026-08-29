package users

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:80;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`

	ConfigLimit int  `gorm:"default:5"`
	IsAdmin     bool `gorm:"default:false"`

	GroupID     *uint
	TotpSecret  string
	TotpEnabled bool
	CreatedAt   time.Time
}
