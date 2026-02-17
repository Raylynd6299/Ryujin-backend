package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvalidUserError(t *testing.T) {
	err := NewInvalidUserError("name cannot be empty")
	assert.NotNil(t, err)
	assert.Equal(t, "name cannot be empty", err.Message())
}

func TestNewUserNotFoundError(t *testing.T) {
	err := NewUserNotFoundError("user-id-123")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "user-id-123")
}

func TestNewInvalidPasswordError(t *testing.T) {
	err := NewInvalidPasswordError("password is incorrect")
	assert.NotNil(t, err)
	assert.Equal(t, "password is incorrect", err.Message())
}

func TestNewDuplicateEmailError(t *testing.T) {
	err := NewDuplicateEmailError("user@example.com")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "user@example.com")
}
