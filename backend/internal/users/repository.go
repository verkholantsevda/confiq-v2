package users

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

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetByID(id uint) (*User, error) {
	var user User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var user User

	err := r.db.
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) List() ([]User, error) {
	var users []User

	if err := r.db.Order("id").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
