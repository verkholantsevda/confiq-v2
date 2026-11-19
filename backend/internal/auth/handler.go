package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"

	"confiq/internal/audit"
	"confiq/internal/httpx"
	"confiq/internal/ratelimit"
)

type Handler struct {
	service *Service
	guard   *ratelimit.LoginGuard
	audit   *audit.Service
}

func NewHandler(service *Service, guard *ratelimit.LoginGuard, auditService *audit.Service) *Handler {
	return &Handler{
		service: service,
		guard:   guard,
		audit:   auditService,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, httpx.ErrInvalidJSON)
		return
	}

	ip := clientIP(r)

	if allowed, retryAfter := h.guard.Check(req.Username, ip); !allowed {
		h.logFailure("auth.rate_limited", req.Username, ip)
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
		httpx.Error(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	token, err := h.service.Login(req.Username, req.Password, req.TOTPCode)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			h.guard.Fail(req.Username, ip, false)
			h.logFailure("auth.login_failed", req.Username, ip)
		case errors.Is(err, ErrInvalidTOTP):
			h.guard.Fail(req.Username, ip, true)
			h.logFailure("auth.totp_failed", req.Username, ip)
		}

		handleError(w, err)
		return
	}

	h.guard.Success(req.Username, ip)

	httpx.OK(w, LoginResponse{
		Token: token,
	})
}

func (h *Handler) logFailure(action, username, ip string) {
	if h.audit == nil {
		return
	}

	if err := h.audit.Log(
		nil,
		action,
		"user",
		nil,
		fmt.Sprintf("%s (IP: %s)", username, ip),
	); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
}

// clientIP берёт X-Real-IP, который nginx перезаписывает на $remote_addr.
// X-Forwarded-For не используется: nginx дописывает туда значение клиента.
func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
