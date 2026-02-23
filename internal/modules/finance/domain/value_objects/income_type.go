package value_objects

import "errors"

// IncomeType classifies the source of income
type IncomeType string

const (
	IncomeTypeSalary    IncomeType = "salary"    // Regular employment salary
	IncomeTypeFreelance IncomeType = "freelance" // Freelance / consulting work
	IncomeTypeRental    IncomeType = "rental"    // Rental property income
	IncomeTypeDividend  IncomeType = "dividend"  // Dividends from investments
	IncomeTypeBusiness  IncomeType = "business"  // Business income
	IncomeTypeOther     IncomeType = "other"     // Any other income type
)

var validIncomeTypes = map[IncomeType]bool{
	IncomeTypeSalary:    true,
	IncomeTypeFreelance: true,
	IncomeTypeRental:    true,
	IncomeTypeDividend:  true,
	IncomeTypeBusiness:  true,
	IncomeTypeOther:     true,
}

// IncomeSourceType is a value object for income classification
type IncomeSourceType struct {
	value IncomeType
}

// NewIncomeSourceType creates a validated IncomeSourceType
func NewIncomeSourceType(t string) (*IncomeSourceType, error) {
	it := IncomeType(t)
	if !validIncomeTypes[it] {
		return nil, errors.New("invalid income type: " + t)
	}
	return &IncomeSourceType{value: it}, nil
}

func (i *IncomeSourceType) Value() IncomeType {
	return i.value
}

func (i *IncomeSourceType) String() string {
	return string(i.value)
}

func (i *IncomeSourceType) IsPassive() bool {
	return i.value == IncomeTypeRental || i.value == IncomeTypeDividend
}
