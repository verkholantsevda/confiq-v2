package configs

import "net/http"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ListAll(w, r, h.service)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	List(w, r, h.service)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, h.service)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	Create(w, r, h.service)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	Update(w, r, h.service)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, h.service)
}
