package endpoints

import (
	"confiq/internal/configtypes"
	"confiq/internal/httpx"
	"confiq/internal/identity"
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
	endpoints, err := h.service.List()
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponseList(endpoints))
}

func (h *Handler) ListForUser(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	endpoints, err := h.service.ListForUser(claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponseList(endpoints))
}

func (h *Handler) ListConfigTypes(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	endpointID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.BadRequest(w, "invalid endpoint id")
		return
	}

	configTypes, err := h.service.ListConfigTypes(
		claims.UserID,
		uint(endpointID),
	)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, configtypes.ToShortResponseList(configTypes))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidEndpointID.Error())
		return
	}

	endpoint, err := h.service.GetByID(uint(id))
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*endpoint))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEndpointRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	endpoint, err := h.service.Create(req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.Created(w, ToResponse(*endpoint))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidEndpointID.Error())
		return
	}

	var req UpdateEndpointRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	endpoint, err := h.service.Update(uint(id), req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*endpoint))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidEndpointID.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		handleError(w, err)
		return
	}

	httpx.NoContent(w)
}

func (h *Handler) ListGroupsEndpoints(w http.ResponseWriter, r *http.Request) {
	groupEndpoints, err := h.service.ListGroupsEndpoints()
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToGroupEndpointResponseList(groupEndpoints))
}

func (h *Handler) AddGroupEndpoint(w http.ResponseWriter, r *http.Request) {
	var req AddGroupEndpointRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	if err := h.service.AddGroupEndpoint(req.GroupID, req.EndpointID); err != nil {
		handleError(w, err)
		return
	}

	httpx.Created(w, nil)
}

func (h *Handler) RemoveGroupEndpoint(w http.ResponseWriter, r *http.Request) {
	groupIDParam := chi.URLParam(r, "group_id")
	endpointIDParam := chi.URLParam(r, "endpoint_id")

	groupID, err := strconv.ParseUint(groupIDParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidGroupEndpointID.Error())
		return
	}

	endpointID, err := strconv.ParseUint(endpointIDParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidGroupEndpointID.Error())
		return
	}

	if err := h.service.RemoveGroupEndpoint(uint(groupID), uint(endpointID)); err != nil {
		handleError(w, err)
		return
	}

	httpx.NoContent(w)
}
