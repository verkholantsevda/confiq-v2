package totp

import (
	"github.com/pquerna/otp/totp"
)

func GenerateSecret() (string, error) {
	key, err := totp.Generate(
		totp.GenerateOpts{
			Issuer:      "Confiq",
			AccountName: "user",
		},
	)

	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

func Validate(secret, code string) bool {
	return totp.Validate(code, secret)
}
