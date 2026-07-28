package endpoints

import (
	"errors"
	"net/http"

	"confiq/internal/httpx"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrEndpointNotFound):
		httpx.NotFound(w, err.Error())

	case errors.Is(err, ErrInvalidEndpointID):
		httpx.BadRequest(w, err.Error())

	default:
		httpx.InternalServerError(w, err)
	}

}
