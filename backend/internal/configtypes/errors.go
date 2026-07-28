package configtypes

import "errors"

var (
	ErrConfigTypeNotFound      = errors.New("config type not found")
	ErrConfigTypeAlreadyExists = errors.New("config type already exists")
	ErrInvalidConfigTypeID     = errors.New("invalid config type id")
)
