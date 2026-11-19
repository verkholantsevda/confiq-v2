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

	case errors.Is(err, ErrTOTPRequired):
		httpx.BadRequest(w, err.Error())

	case errors.Is(err, ErrInvalidTOTP):
		httpx.Error(w, http.StatusUnauthorized, err.Error())

	default:
		httpx.InternalServerError(w, err)
	}
}
