package auth

import (
	"confiq/internal/logger"
	"confiq/internal/password"
	"confiq/internal/totp"
	"confiq/internal/users"
)

type Service struct {
	users *users.Service
	jwt   *JWT
}

func NewService(
	usersService *users.Service,
	jwt *JWT,
) *Service {

	return &Service{
		users: usersService,
		jwt:   jwt,
	}
}

func (s *Service) Login(username, plainPassword string, totpCode string) (string, error) {

	user, err := s.users.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	if !password.CheckPassword(plainPassword, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	if s.users.CanUseTOTP(user) && user.TotpEnabled {
		if totpCode == "" {
			return "", ErrTOTPRequired
		}

		if !totp.Validate(user.TotpSecret, totpCode) {
			return "", ErrInvalidTOTP
		}
	}

	token, err := s.jwt.GenerateToken(
		user.ID,
		user.Username,
		user.IsAdmin,
	)
	if err != nil {
		return "", err
	}
	logger.UserLogin(user.ID, user.Username)
	return token, nil
}
