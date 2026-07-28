package configs

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(config *Config) error {
	return r.db.Create(config).Error
}

func (r *Repository) GetByID(id uint) (*Config, error) {
	var config Config

	if err := r.db.
		Preload("User").
		Preload("Endpoint").
		Preload("ConfigType").
		First(&config, id).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *Repository) List() ([]Config, error) {
	var configs []Config

	if err := r.db.
		Preload("User").
		Preload("Endpoint").
		Preload("ConfigType").
		Order("id").
		Find(&configs).Error; err != nil {
		return nil, err
	}

	return configs, nil
}

func (r *Repository) ListByUser(userID uint) ([]Config, error) {
	var configs []Config

	err := r.db.
		Preload("User").
		Preload("Endpoint").
		Preload("ConfigType").
		Where("user_id = ?", userID).
		Order("id").
		Find(&configs).Error

	return configs, err
}

func (r *Repository) Update(config *Config) error {
	return r.db.Save(config).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Config{}, id).Error
}
