package users

import (
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
	users, err := h.service.ListWithConfigurations()
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, users)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.Created(w, ToResponse(*user))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidUserID.Error())
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	// Non-admin users may only change their own password.
	if !claims.IsAdmin {
		if claims.UserID != uint(id) {
			httpx.Forbidden(w)
			return
		}

		req.Username = ""
		req.GroupID = nil
		req.ConfigLimit = nil
		req.IsAdmin = nil
	}

	user, err := h.service.UpdateUser(uint(id), req)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*user))
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {

	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, ToResponse(*user))
}

type enableTOTPRequest struct {
	Code string `json:"code"`
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {

	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	var req ChangePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		httpx.BadRequest(w, "password is required")
		return
	}

	err := h.service.ChangePassword(
		claims.UserID,
		req.CurrentPassword,
		req.NewPassword,
	)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.NoContent(w)
}

func (h *Handler) TOTPStatus(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}
	status, err := h.service.GetTOTPStatus(claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.OK(w, status)
}

func (h *Handler) SetupTOTP(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	result, err := h.service.CreateTOTP(claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}

	httpx.OK(w, map[string]string{
		"secret": result,
	})
}

func (h *Handler) EnableTOTP(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}
	var req enableTOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}
	err := h.service.ConfirmTOTP(claims.UserID, req.Code)
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (h *Handler) DisableTOTP(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}
	err := h.service.DisableTOTP(claims.UserID)
	if err != nil {
		handleError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		httpx.BadRequest(w, ErrInvalidUserID.Error())
		return
	}

	claims := identity.GetClaims(r)
	if claims == nil {
		httpx.Unauthorized(w)
		return
	}

	if claims.UserID == uint(id) {
		httpx.BadRequest(w, "cannot delete yourself")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		handleError(w, err)
		return
	}

	httpx.NoContent(w)
}
