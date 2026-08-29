package users

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrTOTPDisabledForUser  = errors.New("totp is disabled for regular users")
	ErrTOTPDisabledForAdmin = errors.New("totp is disabled for administrators")
)
