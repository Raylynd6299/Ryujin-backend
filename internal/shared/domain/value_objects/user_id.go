package value_objects

import (
	"errors"
	"strings"
)

// UserID represents a user identifier value object.
type UserID string

// NewUserID creates a new UserID value object.
func NewUserID(value string) (UserID, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("user id cannot be empty")
	}

	return UserID(value), nil
}

// String returns the raw identifier.
func (id UserID) String() string {
	return string(id)
}

// IsZero checks if the identifier is empty.
func (id UserID) IsZero() bool {
	return strings.TrimSpace(string(id)) == ""
}

// Equals checks if two identifiers are the same.
func (id UserID) Equals(other UserID) bool {
	return id == other
}
