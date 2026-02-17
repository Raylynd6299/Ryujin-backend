package value_objects

import (
	"errors"
	"strings"
)

// Currency represents an ISO 4217 currency code
type Currency struct {
	code string
}

// Supported currencies
var supportedCurrencies = map[string]bool{
	"USD": true,
	"MXN": true,
	"EUR": true,
}

// NewCurrency creates a new Currency value object
func NewCurrency(code string) (*Currency, error) {
	if code == "" {
		return nil, errors.New("currency code cannot be empty")
	}

	// Normalize to uppercase
	code = strings.ToUpper(code)

	// Validate format (3 letters)
	if len(code) != 3 {
		return nil, errors.New("currency code must be 3 characters")
	}

	// Validate against supported currencies
	if !supportedCurrencies[code] {
		return nil, errors.New("unsupported currency code: " + code)
	}

	return &Currency{code: code}, nil
}

// Code returns the ISO 4217 currency code
func (c *Currency) Code() string {
	return c.code
}

// Equals checks if two currencies are equal
func (c *Currency) Equals(other *Currency) bool {
	return c.code == other.code
}

// String returns the currency code as a string
func (c *Currency) String() string {
	return c.code
}

// IsSupportedCurrency checks if a currency code is supported
func IsSupportedCurrency(code string) bool {
	return supportedCurrencies[strings.ToUpper(code)]
}

// GetSupportedCurrencies returns a list of all supported currency codes
func GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(supportedCurrencies))
	for code := range supportedCurrencies {
		currencies = append(currencies, code)
	}
	return currencies
}
