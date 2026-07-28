package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MinLength = 3

var ErrPasswordTooShort = errors.New("password too short")

func Validate(password string) error {
	if len(password) < MinLength {
		return ErrPasswordTooShort
	}

	return nil
}

func HashPassword(password string) (string, error) {
	if err := Validate(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}
