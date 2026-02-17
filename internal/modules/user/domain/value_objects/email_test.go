package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple email", "user@example.com"},
		{"with numbers", "user123@example.com"},
		{"with dots in local", "first.last@example.com"},
		{"with hyphen in domain", "user@my-domain.com"},
		{"multiple subdomains", "user@mail.example.co.uk"},
		{"gmail", "john@gmail.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, email.String())
		})
	}
}

func TestNewEmailInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{"empty", "", "empty email"},
		{"whitespace", "   ", "whitespace email"},
		{"no at sign", "userexample.com", "missing @"},
		{"multiple at signs", "user@exam@ple.com", "multiple @"},
		{"no local part", "@example.com", "no local part"},
		{"no domain", "user@", "no domain"},
		{"no dot in domain", "user@example", "domain without dot"},
		{"at start", "@example.com", "@ at start"},
		{"space in email", "user @example.com", "space in email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			assert.Error(t, err, tt.desc)
			assert.Empty(t, email.String())
		})
	}
}

func TestEmailEquals(t *testing.T) {
	email1, _ := NewEmail("user@example.com")
	email2, _ := NewEmail("user@example.com")
	email3, _ := NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

func TestEmailIsZero(t *testing.T) {
	email := Email("")
	assert.True(t, email.IsZero())

	email2, _ := NewEmail("user@example.com")
	assert.False(t, email2.IsZero())
}
