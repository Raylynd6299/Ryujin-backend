package entities

import (
	"fmt"
	"strings"
)

// ============================================================
// AssetType
// ============================================================

// AssetType represents the type of a financial asset
type AssetType string

const (
	AssetTypeStock       AssetType = "stock"
	AssetTypeETF         AssetType = "etf"
	AssetTypeFixedIncome AssetType = "fixed_income"
	AssetTypeCrypto      AssetType = "crypto"
	AssetTypeREIT        AssetType = "reit"
)

var validAssetTypes = map[AssetType]bool{
	AssetTypeStock:       true,
	AssetTypeETF:         true,
	AssetTypeFixedIncome: true,
	AssetTypeCrypto:      true,
	AssetTypeREIT:        true,
}

// NewAssetType validates and returns an AssetType
func NewAssetType(s string) (AssetType, error) {
	at := AssetType(strings.ToLower(strings.TrimSpace(s)))
	if !validAssetTypes[at] {
		return "", fmt.Errorf("invalid asset type: %s (valid: stock, etf, fixed_income, crypto, reit)", s)
	}
	return at, nil
}

// String returns the string representation of the AssetType
func (a AssetType) String() string {
	return string(a)
}

// ============================================================
// Quantity
// ============================================================

// Quantity represents a number of shares/units in micro-units.
// 1 share = 1_000_000 micro-units, allowing up to 6 decimal places.
type Quantity struct {
	microUnits int64
}

// NewQuantity creates a Quantity from micro-units (must be > 0)
func NewQuantity(microUnits int64) (Quantity, error) {
	if microUnits <= 0 {
		return Quantity{}, fmt.Errorf("quantity must be greater than zero")
	}
	return Quantity{microUnits: microUnits}, nil
}

// MicroUnits returns the raw micro-unit value
func (q Quantity) MicroUnits() int64 {
	return q.microUnits
}

// ToFloat returns the quantity as a float64 (microUnits / 1_000_000)
func (q Quantity) ToFloat() float64 {
	return float64(q.microUnits) / 1_000_000.0
}

// String returns a human-readable representation
func (q Quantity) String() string {
	return fmt.Sprintf("%.6f", q.ToFloat())
}

// ============================================================
// Symbol
// ============================================================

// Symbol represents a stock/asset ticker symbol (1-10 uppercase characters)
type Symbol struct {
	value string
}

// NewSymbol creates a Symbol after validating length and format
func NewSymbol(s string) (Symbol, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return Symbol{}, fmt.Errorf("symbol cannot be empty")
	}
	if len(s) > 10 {
		return Symbol{}, fmt.Errorf("symbol cannot exceed 10 characters")
	}

	upper := strings.ToUpper(s)
	// Symbols may include letters, numbers, dots, and hyphens (e.g. BRK.B, BTC-USD)
	for _, ch := range upper {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '.' || ch == '-') {
			return Symbol{}, fmt.Errorf("symbol contains invalid character: %c", ch)
		}
	}

	return Symbol{value: upper}, nil
}

// Value returns the symbol string
func (s Symbol) Value() string {
	return s.value
}

// String returns the symbol as a string
func (s Symbol) String() string {
	return s.value
}
