package users

import (
	"errors"
	"net/http"

	"confiq/internal/httpx"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		httpx.BadRequest(w, err.Error())

	case errors.Is(err, ErrUserNotFound):
		httpx.NotFound(w, err.Error())

	case errors.Is(err, ErrInvalidUserID):
		httpx.BadRequest(w, err.Error())

	default:
		httpx.InternalServerError(w, err)
	}
}
