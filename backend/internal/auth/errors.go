package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTOTPRequired       = errors.New("totp required")
	ErrInvalidTOTP        = errors.New("invalid totp code")
)
