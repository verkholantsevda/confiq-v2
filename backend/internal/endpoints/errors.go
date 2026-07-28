package endpoints

import "errors"

var (
	ErrEndpointNotFound      = errors.New("endpoint not found")
	ErrEndpointAlreadyExists = errors.New("endpoint already exists")
	ErrInvalidEndpointID     = errors.New("invalid endpoint id")
)

var (
	ErrGroupEndpointNotFound      = errors.New("group endpoint not found")
	ErrGroupEndpointAlreadyExists = errors.New("group endpoint already exists")
	ErrInvalidGroupEndpointID     = errors.New("invalid group endpoint id")
)
