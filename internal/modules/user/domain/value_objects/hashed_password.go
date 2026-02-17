package value_objects

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// HashedPassword represents a bcrypt-hashed password that was previously validated as a Password.
// This value object does NOT validate the hash format - it assumes the input is already a valid bcrypt hash.
type HashedPassword string

// NewHashedPassword creates a new HashedPassword value object.
// Note: This does not validate the hash format - assumes input is already a valid bcrypt hash.
func NewHashedPassword(hash string) HashedPassword {
	return HashedPassword(hash)
}

// String returns the hash string.
// Should only be used for storage/transmission, never for password comparison.
func (hp HashedPassword) String() string {
	return string(hp)
}

// CompareWith verifies if a raw Password matches this hash.
// Returns true if the password matches, false otherwise.
// Never panics - safe for authentication attempts.
func (hp HashedPassword) CompareWith(pwd Password) bool {
	// Get SECRET_PHRASE from environment
	secretPhrase := os.Getenv("SECRET_PHRASE")
	if secretPhrase == "" {
		secretPhrase = "default-secret-phrase" // Fallback for development
	}

	// Combine password with secret phrase for comparison
	combinedPassword := fmt.Sprintf("%s:%s", pwd.String(), secretPhrase)

	// bcrypt.CompareHashAndPassword returns nil if match, error otherwise
	err := bcrypt.CompareHashAndPassword([]byte(hp.String()), []byte(combinedPassword))
	return err == nil
}

// HashPassword is a utility function to hash a Password into a HashedPassword.
// Used at entity creation and password change.
// Combines the raw password with a SECRET_PHRASE environment variable before hashing.
func HashPassword(pwd Password) (HashedPassword, error) {
	// Get SECRET_PHRASE from environment
	secretPhrase := os.Getenv("SECRET_PHRASE")
	if secretPhrase == "" {
		secretPhrase = "default-secret-phrase" // Fallback for development
	}

	// Combine password with secret phrase
	combinedPassword := fmt.Sprintf("%s:%s", pwd.String(), secretPhrase)

	// Generate bcrypt hash with salt (bcrypt includes salt internally)
	hash, err := bcrypt.GenerateFromPassword([]byte(combinedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return NewHashedPassword(string(hash)), nil
}
