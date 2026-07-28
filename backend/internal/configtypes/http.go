package configtypes

import (
	"errors"
	"net/http"

	"confiq/internal/httpx"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrConfigTypeAlreadyExists):
		httpx.BadRequest(w, err.Error())

	case errors.Is(err, ErrConfigTypeNotFound):
		httpx.NotFound(w, err.Error())

	case errors.Is(err, ErrInvalidConfigTypeID):
		httpx.BadRequest(w, err.Error())

	default:
		httpx.InternalServerError(w, err)
	}
}
