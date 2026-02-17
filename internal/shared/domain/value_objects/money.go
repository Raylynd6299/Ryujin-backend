package value_objects

import (
	"errors"
	"fmt"
)

// Money represents a monetary value in the smallest currency unit (cents)
// This ensures precision and avoids floating-point arithmetic issues
type Money struct {
	amount   int64  // Amount in cents (smallest unit)
	currency string // ISO 4217 currency code (USD, MXN, EUR, etc.)
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency string) (*Money, error) {
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	// Validate currency code (basic validation - 3 uppercase letters)
	if len(currency) != 3 {
		return nil, fmt.Errorf("invalid currency code: %s (must be 3 characters)", currency)
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount returns the amount in cents
func (m *Money) Amount() int64 {
	return m.amount
}

// Currency returns the ISO 4217 currency code
func (m *Money) Currency() string {
	return m.currency
}

// ToDecimal converts cents to decimal representation (e.g., 1050 cents -> 10.50)
func (m *Money) ToDecimal() float64 {
	return float64(m.amount) / 100.0
}

// Add adds two Money values (must be same currency)
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
	}

	return &Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts two Money values (must be same currency)
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot subtract different currencies: %s and %s", m.currency, other.currency)
	}

	return &Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply multiplies Money by a factor
func (m *Money) Multiply(factor float64) *Money {
	return &Money{
		amount:   int64(float64(m.amount) * factor),
		currency: m.currency,
	}
}

// IsZero checks if the amount is zero
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive checks if the amount is positive
func (m *Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative checks if the amount is negative
func (m *Money) IsNegative() bool {
	return m.amount < 0
}

// Equals checks if two Money values are equal
func (m *Money) Equals(other *Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns a string representation of the Money value
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.ToDecimal(), m.currency)
}
