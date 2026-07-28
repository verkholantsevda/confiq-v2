package configs

import "errors"

var (
	ErrConfigNotFound  = errors.New("config not found")
	ErrInvalidConfigID = errors.New("invalid config id")

	ErrUserNotFound       = errors.New("user not found")
	ErrEndpointNotFound   = errors.New("endpoint not found")
	ErrConfigTypeNotFound = errors.New("config type not found")
)
