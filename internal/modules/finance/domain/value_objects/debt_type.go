package value_objects

import "errors"

// DebtType classifies the kind of debt
type DebtType string

const (
	DebtTypeCreditCard   DebtType = "credit_card"   // Credit card balance
	DebtTypePersonalLoan DebtType = "personal_loan" // Personal bank loan
	DebtTypeMortgage     DebtType = "mortgage"      // Home mortgage
	DebtTypeCarLoan      DebtType = "car_loan"      // Car financing
	DebtTypeStudentLoan  DebtType = "student_loan"  // Student loan
	DebtTypeOther        DebtType = "other"         // Any other debt
)

var validDebtTypes = map[DebtType]bool{
	DebtTypeCreditCard:   true,
	DebtTypePersonalLoan: true,
	DebtTypeMortgage:     true,
	DebtTypeCarLoan:      true,
	DebtTypeStudentLoan:  true,
	DebtTypeOther:        true,
}

// DebtCategory is a value object for debt classification
type DebtCategory struct {
	value DebtType
}

// NewDebtCategory creates a validated DebtCategory
func NewDebtCategory(t string) (*DebtCategory, error) {
	dt := DebtType(t)
	if !validDebtTypes[dt] {
		return nil, errors.New("invalid debt type: " + t)
	}
	return &DebtCategory{value: dt}, nil
}

func (d *DebtCategory) Value() DebtType {
	return d.value
}

func (d *DebtCategory) String() string {
	return string(d.value)
}

// IsHighInterest returns true for typically high-interest debt
func (d *DebtCategory) IsHighInterest() bool {
	return d.value == DebtTypeCreditCard
}
