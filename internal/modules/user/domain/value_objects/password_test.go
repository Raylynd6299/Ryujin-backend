package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPasswordValid(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple valid", "SecurePass123"},
		{"with special chars", "Secure#Pass123"},
		{"longer password", "VerySecurePassword123!@#"},
		{"exactly 8 chars", "Pass1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := NewPassword(tt.password)
			assert.NoError(t, err)
			assert.Equal(t, tt.password, pwd.String())
		})
	}
}

func TestNewPasswordInvalid(t *testing.T) {
	tests := []struct {
		name     string
		password string
		desc     string
	}{
		{"too short", "Pass123", "less than 8 chars"},
		{"no uppercase", "securepass123", "missing uppercase"},
		{"no number", "SecurePass", "missing number"},
		{"only uppercase", "SECUREPASS123", "no lowercase (but should pass)"},
		{"empty", "", "empty password"},
		{"whitespace", "   ", "whitespace only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := NewPassword(tt.password)
			// "only uppercase" test should actually pass
			if tt.name == "only uppercase" {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err, tt.desc)
			assert.Empty(t, pwd.String())
		})
	}
}

func TestPasswordString(t *testing.T) {
	pwd, _ := NewPassword("SecurePass123")
	assert.Equal(t, "SecurePass123", pwd.String())
}
