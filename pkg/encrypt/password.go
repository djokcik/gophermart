package encrypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string, pepper string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd+pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcryptPassword: %w", err)
	}

	return string(hash), nil
}

func CompareHashAndPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
