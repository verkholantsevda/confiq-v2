package configtypes

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

func (r *Repository) Create(configType *ConfigType) error {
	return r.db.Create(configType).Error
}

func (r *Repository) GetByID(id uint) (*ConfigType, error) {
	var configType ConfigType

	if err := r.db.First(&configType, id).Error; err != nil {
		return nil, err
	}

	return &configType, nil
}

func (r *Repository) GetByName(name string) (*ConfigType, error) {
	var configType ConfigType

	if err := r.db.Where("name = ?", name).First(&configType).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &configType, nil
}

func (r *Repository) List() ([]ConfigType, error) {
	var configTypes []ConfigType

	if err := r.db.Order("id").Find(&configTypes).Error; err != nil {
		return nil, err
	}

	return configTypes, nil
}

func (r *Repository) Update(configType *ConfigType) error {
	return r.db.Save(configType).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&ConfigType{}, id).Error
}
