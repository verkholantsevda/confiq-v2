package httpx

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ErrInvalidJSON  = "invalid json"
	ErrUnauthorized = "unauthorized"
	ErrForbidden    = "forbidden"
	ErrNotFound     = "not found"
)

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Error: message,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, ErrUnauthorized)
}

func Forbidden(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, ErrForbidden)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func InternalServerError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, err.Error())
}
