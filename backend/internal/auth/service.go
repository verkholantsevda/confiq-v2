package auth

import (
	"confiq/internal/audit"
	"confiq/internal/logger"
	"confiq/internal/password"
	"confiq/internal/totp"
	"confiq/internal/users"
	"log/slog"
)

type Service struct {
	users *users.Service
	jwt   *JWT
	audit *audit.Service
}

func NewService(
	usersService *users.Service,
	jwt *JWT,
	auditService *audit.Service,
) *Service {

	return &Service{
		users: usersService,
		jwt:   jwt,
		audit: auditService,
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
	if err := s.audit.Log(
		&user.ID,
		"auth.login",
		"user",
		&user.ID,
		"Выполнен вход пользователя "+user.Username,
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
	return token, nil
}
