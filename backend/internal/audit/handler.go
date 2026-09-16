package audit

import (
	"net/http"
	"strconv"

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := 10

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			httpx.BadRequest(w, "invalid limit")
			return
		}

		limit = parsed
	}

	logs, err := h.service.ListRecent(limit)
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	httpx.OK(w, logs)
}
