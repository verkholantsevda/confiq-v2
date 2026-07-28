package groups

import (
	"errors"

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

func (r *Repository) Create(group *Group) error {
	return r.db.Create(group).Error
}

func (r *Repository) GetByID(id uint) (*Group, error) {
	var group Group

	if err := r.db.Preload("Endpoints").First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &group, nil
}

func (r *Repository) GetByName(name string) (*Group, error) {
	var group Group

	err := r.db.
		Where("name = ?", name).
		First(&group).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &group, nil
}

func (r *Repository) List() ([]Group, error) {
	var groups []Group

	if err := r.db.Preload("Endpoints").Order("id").Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *Repository) Update(group *Group) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(group).Error; err != nil {
			return err
		}

		if err := tx.Model(group).
			Association("Endpoints").
			Replace(group.Endpoints); err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Group{}, id).Error
}
