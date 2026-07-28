package groups

import "errors"

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrInvalidGroupID     = errors.New("invalid group id")
)
