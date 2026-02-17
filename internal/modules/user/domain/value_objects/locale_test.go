package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocaleValid(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"spanish", "es", "es"},
		{"english", "en", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locale, err := NewLocale(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, locale.String())
		})
	}
}

func TestNewLocaleInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{"empty", "", "empty locale"},
		{"fr", "fr", "unsupported locale"},
		{"uppercase ES", "ES", "wrong case"},
		{"whitespace", "   ", "whitespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locale, err := NewLocale(tt.input)
			assert.Error(t, err, tt.desc)
			assert.Empty(t, locale.String())
		})
	}
}

func TestLocaleIsSpanish(t *testing.T) {
	localeES, _ := NewLocale("es")
	assert.True(t, localeES.IsSpanish())

	localeEN, _ := NewLocale("en")
	assert.False(t, localeEN.IsSpanish())
}

func TestLocaleIsEnglish(t *testing.T) {
	localeEN, _ := NewLocale("en")
	assert.True(t, localeEN.IsEnglish())

	localeES, _ := NewLocale("es")
	assert.False(t, localeES.IsEnglish())
}

func TestDefaultLocale(t *testing.T) {
	locale := DefaultLocale()
	assert.Equal(t, "es", locale.String())
	assert.True(t, locale.IsSpanish())
}
