package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
	vo "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/value_objects"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// Debt represents a financial liability for a user.
// Tracks total amount, remaining balance, monthly payments, and interest rate.
type Debt struct {
	ID     string
	UserID string

	Name        string
	Description string
	DebtType    vo.DebtCategory

	// Financial data (all in smallest currency unit)
	TotalAmount     sharedVO.Money // original loan/debt amount
	RemainingAmount sharedVO.Money // current remaining balance
	MonthlyPayment  sharedVO.Money // minimum monthly payment
	InterestRate    float64        // annual interest rate as percentage (e.g. 18.5 for 18.5%)

	// Dates
	StartDate *time.Time // when the debt was acquired
	DueDate   *time.Time // final payoff date (nil = no fixed date, e.g. credit card)
	IsActive  bool       // false when fully paid

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDebt creates a new debt with validation
func NewDebt(
	userID string,
	name string,
	description string,
	debtType string,
	totalAmountCents int64,
	remainingAmountCents int64,
	monthlyPaymentCents int64,
	currency string,
	interestRate float64,
	startDate *time.Time,
	dueDate *time.Time,
) (*Debt, error) {
	if strings.TrimSpace(name) == "" {
		return nil, financeErrors.NewDebtInvalidError("debt name cannot be empty")
	}

	if totalAmountCents <= 0 {
		return nil, financeErrors.NewDebtInvalidError("total amount must be greater than zero")
	}

	if remainingAmountCents < 0 {
		return nil, financeErrors.NewDebtInvalidError("remaining amount cannot be negative")
	}

	if remainingAmountCents > totalAmountCents {
		return nil, financeErrors.NewDebtInvalidError("remaining amount cannot exceed total amount")
	}

	if monthlyPaymentCents <= 0 {
		return nil, financeErrors.NewDebtInvalidError("monthly payment must be greater than zero")
	}

	if interestRate < 0 {
		return nil, financeErrors.NewDebtInvalidError("interest rate cannot be negative")
	}

	totalMoney, err := sharedVO.NewMoney(totalAmountCents, currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid currency: " + err.Error())
	}

	remainingMoney, err := sharedVO.NewMoney(remainingAmountCents, currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid currency: " + err.Error())
	}

	monthlyMoney, err := sharedVO.NewMoney(monthlyPaymentCents, currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid currency: " + err.Error())
	}

	dt, err := vo.NewDebtCategory(debtType)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError(err.Error())
	}

	now := time.Now()
	return &Debt{
		ID:              uuid.New().String(),
		UserID:          userID,
		Name:            strings.TrimSpace(name),
		Description:     strings.TrimSpace(description),
		DebtType:        *dt,
		TotalAmount:     *totalMoney,
		RemainingAmount: *remainingMoney,
		MonthlyPayment:  *monthlyMoney,
		InterestRate:    interestRate,
		StartDate:       startDate,
		DueDate:         dueDate,
		IsActive:        remainingAmountCents > 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Update updates debt metadata (not payments — use RecordPayment for that)
func (d *Debt) Update(
	name string,
	description string,
	monthlyPaymentCents int64,
	currency string,
	interestRate float64,
	dueDate *time.Time,
) error {
	if strings.TrimSpace(name) == "" {
		return financeErrors.NewDebtInvalidError("debt name cannot be empty")
	}

	if monthlyPaymentCents <= 0 {
		return financeErrors.NewDebtInvalidError("monthly payment must be greater than zero")
	}

	monthlyMoney, err := sharedVO.NewMoney(monthlyPaymentCents, currency)
	if err != nil {
		return financeErrors.NewDebtInvalidError("invalid currency: " + err.Error())
	}

	d.Name = strings.TrimSpace(name)
	d.Description = strings.TrimSpace(description)
	d.MonthlyPayment = *monthlyMoney
	d.InterestRate = interestRate
	d.DueDate = dueDate
	d.UpdatedAt = time.Now()

	return nil
}

// RecordPayment reduces the remaining balance by the payment amount
func (d *Debt) RecordPayment(paymentCents int64) error {
	if paymentCents <= 0 {
		return financeErrors.NewDebtInvalidError("payment amount must be greater than zero")
	}

	remaining := d.RemainingAmount.Amount()
	if paymentCents > remaining {
		paymentCents = remaining // cap at remaining balance
	}

	newRemaining, _ := sharedVO.NewMoney(remaining-paymentCents, d.RemainingAmount.Currency())
	d.RemainingAmount = *newRemaining

	if d.RemainingAmount.IsZero() {
		d.IsActive = false
	}

	d.UpdatedAt = time.Now()
	return nil
}

// BelongsTo checks if this debt belongs to the given user
func (d *Debt) BelongsTo(userID string) bool {
	return d.UserID == userID
}

// IsPaidOff returns true if the debt is fully paid
func (d *Debt) IsPaidOff() bool {
	return d.RemainingAmount.IsZero()
}

// ProgressPercent returns how much of the debt has been paid (0-100)
func (d *Debt) ProgressPercent() float64 {
	total := d.TotalAmount.Amount()
	if total == 0 {
		return 100.0
	}
	remaining := d.RemainingAmount.Amount()
	paid := total - remaining
	return float64(paid) / float64(total) * 100.0
}

// MonthsToPayoff estimates months remaining to pay off the debt
// Simple calculation: remaining / monthly payment (ignores compounding)
func (d *Debt) MonthsToPayoff() int {
	monthly := d.MonthlyPayment.Amount()
	remaining := d.RemainingAmount.Amount()
	if monthly == 0 || remaining == 0 {
		return 0
	}
	months := int((remaining + monthly - 1) / monthly) // ceiling division
	return months
}
