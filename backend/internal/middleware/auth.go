package middleware

import (
	"net/http"
	"strings"

	"confiq/internal/auth"
	"confiq/internal/identity"
)

type Auth struct {
	jwt *auth.JWT
}

func NewAuth(jwt *auth.JWT) *Auth {
	return &Auth{
		jwt: jwt,
	}
}

func (m *Auth) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		claims, err := m.jwt.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, identity.WithClaims(r, claims))
	})
}
