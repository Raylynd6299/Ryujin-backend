package value_objects

import (
	"errors"
	"regexp"
	"strings"
)

// Email represents a validated email address value object.
type Email string

// NewEmail creates a new Email value object with validation.
// Returns error if email format is invalid.
func NewEmail(value string) (Email, error) {
	trimmed := strings.TrimSpace(value)

	// Check not empty
	if trimmed == "" {
		return "", errors.New("email cannot be empty")
	}

	// Basic format: must have exactly one @ and contain a dot
	if strings.Count(trimmed, "@") != 1 {
		return "", errors.New("email must contain exactly one @")
	}

	parts := strings.Split(trimmed, "@")
	localPart := parts[0]
	domain := parts[1]

	// Validate local part (before @)
	if localPart == "" {
		return "", errors.New("email local part cannot be empty")
	}

	// Validate domain part (after @)
	if domain == "" {
		return "", errors.New("email domain cannot be empty")
	}

	if !strings.Contains(domain, ".") {
		return "", errors.New("email domain must contain at least one dot")
	}

	// More strict format validation with regex
	// Allow alphanumeric, dots, hyphens, underscores in local part
	// Allow alphanumeric, dots, hyphens in domain
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(trimmed) {
		return "", errors.New("email format is invalid")
	}

	return Email(trimmed), nil
}

// String returns the email as a string.
func (e Email) String() string {
	return string(e)
}

// IsZero checks if the email is empty.
func (e Email) IsZero() bool {
	return strings.TrimSpace(string(e)) == ""
}

// Equals checks if two emails are equal.
func (e Email) Equals(other Email) bool {
	return e == other
}
