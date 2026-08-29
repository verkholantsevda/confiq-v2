package auth

import (
	"encoding/json"
	"net/http"

	"confiq/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	token, err := h.service.Login(req.Username, req.Password, req.TOTPCode)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, LoginResponse{
		Token: token,
	})
}
