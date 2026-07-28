package configs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"confiq/internal/identity"

	"github.com/go-chi/chi/v5"
)

func List(w http.ResponseWriter, r *http.Request, service *Service) {
	claims := identity.GetClaims(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	configs, err := service.List(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ToResponseList(configs))
}

func ListAll(w http.ResponseWriter, r *http.Request, service *Service) {
	configs, err := service.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ToResponseList(configs))
}

func Get(w http.ResponseWriter, r *http.Request, service *Service) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	claims := identity.GetClaims(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	config, err := service.GetByID(uint(id), claims.UserID, claims.IsAdmin)
	if err != nil {
		if err == ErrConfigNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ToResponse(*config))
}

func Create(w http.ResponseWriter, r *http.Request, service *Service) {
	var req CreateConfigRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	claims := identity.GetClaims(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	config, err := service.Create(claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToResponse(*config))
}

func Update(w http.ResponseWriter, r *http.Request, service *Service) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateConfigRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	claims := identity.GetClaims(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	config, err := service.Update(uint(id), claims.UserID, claims.IsAdmin, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(ToResponse(*config))
}

func Delete(w http.ResponseWriter, r *http.Request, service *Service) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	claims := identity.GetClaims(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := service.Delete(uint(id), claims.UserID, claims.IsAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
