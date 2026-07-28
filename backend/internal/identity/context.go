package identity

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

func WithClaims(r *http.Request, claims *Claims) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, claims)
	return r.WithContext(ctx)
}

func GetClaims(r *http.Request) *Claims {
	claims, ok := r.Context().Value(UserContextKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
