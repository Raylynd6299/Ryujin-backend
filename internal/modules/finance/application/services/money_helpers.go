package services

import "math"

// decimalToCents converts a decimal amount (e.g. 10.50) to integer cents (1050).
// Rounds to nearest cent to avoid floating-point issues.
func decimalToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

// centsToDecimal converts integer cents (1050) to decimal (10.50).
func centsToDecimal(cents int64) float64 {
	return float64(cents) / 100.0
}
