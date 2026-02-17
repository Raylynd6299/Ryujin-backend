package value_objects

import (
	"errors"
	"strings"
)

// Locale represents a supported application locale/language.
type Locale string

const (
	LocaleSpanish Locale = "es"
	LocaleEnglish Locale = "en"
)

// NewLocale creates a new Locale value object with validation.
// Valid values: "es", "en"
func NewLocale(value string) (Locale, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return "", errors.New("locale cannot be empty")
	}

	locale := Locale(trimmed)

	// Validate against allowed values
	switch locale {
	case LocaleSpanish, LocaleEnglish:
		return locale, nil
	default:
		return "", errors.New("unsupported locale: " + trimmed)
	}
}

// String returns the locale code.
func (l Locale) String() string {
	return string(l)
}

// IsSpanish checks if locale is Spanish.
func (l Locale) IsSpanish() bool {
	return l == LocaleSpanish
}

// IsEnglish checks if locale is English.
func (l Locale) IsEnglish() bool {
	return l == LocaleEnglish
}

// DefaultLocale returns the default locale (Spanish).
func DefaultLocale() Locale {
	return LocaleSpanish
}
