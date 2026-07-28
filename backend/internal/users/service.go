package users

import (
	"confiq/internal/logger"
	"confiq/internal/password"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(req CreateUserRequest) (*User, error) {

	existing, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     req.Username,
		PasswordHash: hash,
		ConfigLimit:  req.ConfigLimit,
		GroupID:      req.GroupID,
		IsAdmin:      req.IsAdmin,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	logger.UserCreated(user.ID, user.Username)
	return user, nil

}

func (s *Service) GetByID(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByUsername(username string) (*User, error) {
	return s.repo.GetByUsername(username)
}

func (s *Service) List() ([]User, error) {
	return s.repo.List()
}

func (s *Service) Delete(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	logger.UserDeleted(user.ID, user.Username)
	return nil
}

func (s *Service) UpdateUser(id uint, req UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" && req.Username != user.Username {
		existing, err := s.repo.GetByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != user.ID {
			return nil, errors.New("user already exists")
		}
		user.Username = req.Username
	}

	if req.Password != "" {
		hash, err := password.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	if req.ConfigLimit != nil {
		user.ConfigLimit = *req.ConfigLimit
	}

	if req.GroupID != nil {
		user.GroupID = req.GroupID
	}

	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	logger.UserUpdated(user.ID, user.Username)
	return user, nil
}

func (s *Service) EnsureAdmin(username, passwordHash string) error {
	if username == "" || passwordHash == "" {
		return nil
	}

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return err
	}

	if user != nil {
		return nil
	}

	return s.repo.Create(&User{
		Username:     username,
		PasswordHash: passwordHash,
		IsAdmin:      true,
		ConfigLimit:  1000,
	})
}

func (s *Service) ChangePassword(id uint, currentPassword, newPassword string) error {

	user, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	ok := password.CheckPassword(currentPassword, user.PasswordHash)
	if !ok {

		return errors.New("invalid current password")

	}

	hash, err := password.HashPassword(newPassword)

	if err != nil {
		return err
	}
	user.PasswordHash = hash
	err = s.repo.Update(user)
	if err != nil {
		return err
	}

	logger.PasswordChanged(user.ID, user.Username)
	return nil

}
