package groups

import (
	"errors"
	"net/http"

	"confiq/internal/httpx"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrGroupAlreadyExists):
		httpx.BadRequest(w, err.Error())

	case errors.Is(err, ErrGroupNotFound):
		httpx.NotFound(w, err.Error())

	case errors.Is(err, ErrInvalidGroupID):
		httpx.BadRequest(w, err.Error())

	default:
		httpx.InternalServerError(w, err)
	}
}
