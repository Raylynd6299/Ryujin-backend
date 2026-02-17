package value_objects

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Password represents a raw, unencrypted password that validates strength requirements.
type Password string

// NewPassword creates a new Password value object with strength validation.
// Requirements:
// - Minimum 8 characters
// - At least one uppercase letter
// - At least one digit
func NewPassword(value string) (Password, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("password cannot be empty")
	}

	// Check minimum length
	if len(value) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	hasUppercase := false
	for _, char := range value {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}
	if !hasUppercase {
		return "", errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(value)
	if !hasDigit {
		return "", errors.New("password must contain at least one digit")
	}

	return Password(value), nil
}

// String returns the raw password string.
// Note: This should only be used internally for hashing, never exposed in responses.
func (p Password) String() string {
	return string(p)
}
