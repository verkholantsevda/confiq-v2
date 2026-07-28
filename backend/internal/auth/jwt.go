package auth

import (
	"time"

	"confiq/internal/identity"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
	expire time.Duration
}

func NewJWT(secret string, expire time.Duration) *JWT {
	return &JWT{
		secret: []byte(secret),
		expire: expire,
	}
}

func (j *JWT) GenerateToken(userID uint, username string, isAdmin bool) (string, error) {
	now := time.Now()

	claims := identity.Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *JWT) ParseToken(tokenString string) (*identity.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&identity.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return j.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*identity.Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (j *JWT) ValidateToken(tokenString string) (*identity.Claims, error) {
	return j.ParseToken(tokenString)
}
