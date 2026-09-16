package users

import (
	"confiq/internal/config"
	"confiq/internal/logger"
	"confiq/internal/password"
	"confiq/internal/totp"
	"errors"
)

type Service struct {
	repo     *Repository
	otpUser  bool
	otpAdmin bool
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo:     repo,
		otpUser:  cfg.OTPUserEnabled,
		otpAdmin: cfg.OTPAdminEnabled,
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

func (s *Service) GetMe(id uint) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	count, err := s.repo.CountConfigurationsByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	response := ToResponse(*user)
	response.Configurations = int(count)
	response.TOTPAvailable = s.CanUseTOTP(user)

	return &response, nil
}

func (s *Service) GetByUsername(username string) (*User, error) {
	return s.repo.GetByUsername(username)
}

func (s *Service) List() ([]User, error) {
	return s.repo.List()
}

func (s *Service) ListWithConfigurations() ([]UserResponse, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	result := make([]UserResponse, 0, len(users))

	for _, user := range users {
		count, err := s.repo.CountConfigurationsByUserID(user.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, UserResponse{
			ID:             user.ID,
			Username:       user.Username,
			ConfigLimit:    user.ConfigLimit,
			GroupID:        user.GroupID,
			IsAdmin:        user.IsAdmin,
			CreatedAt:      user.CreatedAt.Format("2006-01-02 15:04:05"),
			Configurations: int(count),
		})
	}

	return result, nil
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

func (s *Service) CanUseTOTP(user *User) bool {
	if user.IsAdmin {
		return s.otpAdmin
	}

	return s.otpUser
}

func (s *Service) CreateTOTP(id uint) (string, error) {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if !s.CanUseTOTP(user) {
		if user.IsAdmin {
			return "", ErrTOTPDisabledForAdmin
		}
		return "", ErrTOTPDisabledForUser
	}

	if user.TotpEnabled {
		return "", errors.New("totp already enabled")
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		return "", err
	}

	user.TotpSecret = secret

	if err := s.repo.Update(user); err != nil {
		return "", err
	}

	return secret, nil
}

func (s *Service) ConfirmTOTP(id uint, code string) error {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.TotpSecret == "" {
		return errors.New("totp not initialized")
	}

	if !totp.Validate(user.TotpSecret, code) {
		return errors.New("invalid totp code")
	}

	user.TotpEnabled = true

	return s.repo.Update(user)
}

func (s *Service) EnableTOTP(id uint) (string, error) {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if !s.CanUseTOTP(user) {
		if user.IsAdmin {
			return "", ErrTOTPDisabledForAdmin
		}
		return "", ErrTOTPDisabledForUser
	}

	if user.TotpEnabled {
		return "", errors.New("totp already enabled")
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		return "", err
	}

	user.TotpSecret = secret
	user.TotpEnabled = true

	if err := s.repo.Update(user); err != nil {
		return "", err
	}

	return secret, nil
}

func (s *Service) DisableTOTP(id uint) error {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	user.TotpSecret = ""
	user.TotpEnabled = false

	return s.repo.Update(user)
}

type TOTPStatus struct {
	Enabled   bool `json:"enabled"`
	Available bool `json:"available"`
}

func (s *Service) GetTOTPStatus(id uint) (*TOTPStatus, error) {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &TOTPStatus{
		Enabled:   user.TotpEnabled,
		Available: s.CanUseTOTP(user),
	}, nil
}
