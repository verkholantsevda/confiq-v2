package middleware

import (
	"net/http"
	"strings"
	"time"

	"confiq/internal/auth"
	"confiq/internal/identity"
	"confiq/internal/users"
)

// TokenUserProvider нужен middleware, чтобы проверять, что пользователь из
// токена ещё существует и его версия токенов совпадает с версией в БД.
type TokenUserProvider interface {
	GetByID(id uint) (*users.User, error)
}

type Auth struct {
	jwt   *auth.JWT
	users TokenUserProvider
}

func NewAuth(jwt *auth.JWT, userProvider TokenUserProvider) *Auth {
	return &Auth{
		jwt:   jwt,
		users: userProvider,
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

		user, err := m.users.GetByID(claims.UserID)
		if err != nil || user == nil || user.TokenVersion != claims.TokenVersion {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Права берём из БД, а не из токена: разжалование действует сразу.
		claims.Username = user.Username
		claims.IsAdmin = user.IsAdmin

		// Скользящий токен: активные пользователи не разлогиниваются,
		// хотя срок жизни одного токена короткий.
		if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) < m.jwt.ExpiresIn()/2 {
			refreshed, err := m.jwt.GenerateToken(
				user.ID,
				user.Username,
				user.IsAdmin,
				user.TokenVersion,
			)
			if err == nil {
				w.Header().Set("X-Refresh-Token", refreshed)
			}
		}

		next.ServeHTTP(w, identity.WithClaims(r, claims))
	})
}
