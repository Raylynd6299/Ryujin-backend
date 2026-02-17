package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewHashedPasswordValid(t *testing.T) {
	// Create a real bcrypt hash
	password := "SecurePass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	hashedPwd := NewHashedPassword(string(hash))
	assert.NotEmpty(t, hashedPwd.String())
}

func TestHashedPasswordCompareWithValid(t *testing.T) {
	// Use HashPassword to create a proper hashed password (includes SECRET_PHRASE)
	rawPassword := "SecurePass123"
	pwd, _ := NewPassword(rawPassword)
	hashedPwd, _ := HashPassword(pwd)

	// Compare should match
	assert.True(t, hashedPwd.CompareWith(pwd))
}

func TestHashedPasswordCompareWithInvalid(t *testing.T) {
	// Create a hash for one password using HashPassword
	password1 := "SecurePass123"
	pwd1, _ := NewPassword(password1)
	hashedPwd, _ := HashPassword(pwd1)

	// Try comparing with different password
	pwd2, _ := NewPassword("DifferentPass456")

	// Should not match
	assert.False(t, hashedPwd.CompareWith(pwd2))
}

func TestHashedPasswordString(t *testing.T) {
	// Create a simple hash for testing
	password := "SecurePass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPwd := NewHashedPassword(string(hash))

	assert.Equal(t, string(hash), hashedPwd.String())
}

func TestHashedPasswordImmutable(t *testing.T) {
	password := "SecurePass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPwd := NewHashedPassword(string(hash))

	// HashedPassword is immutable, no setter methods should exist
	originalString := hashedPwd.String()
	assert.Equal(t, string(hash), originalString)
}
