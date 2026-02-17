package utils

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password cannot be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// CheckPassword compares a bcrypt password hash with a plaintext password.
func CheckPassword(hashedPassword, password string) error {
	if strings.TrimSpace(hashedPassword) == "" {
		return errors.New("stored password hash cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
