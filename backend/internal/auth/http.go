package auth

import (
	"errors"
	"net/http"

	"confiq/internal/httpx"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		httpx.Unauthorized(w)

	default:
		httpx.InternalServerError(w, err)
	}
}
