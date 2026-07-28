package configtypes

import (
	"confiq/internal/httpx"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	configTypes, err := h.service.List()
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponseList(configTypes))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidConfigTypeID.Error())
		return
	}

	configType, err := h.service.GetByID(uint(id))
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*configType))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateConfigTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	configType, err := h.service.Create(req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.Created(w, ToResponse(*configType))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidConfigTypeID.Error())
		return
	}

	var req UpdateConfigTypeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	configType, err := h.service.Update(uint(id), req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*configType))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidConfigTypeID.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		handleError(w, err)
		return
	}

	httpx.NoContent(w)
}
