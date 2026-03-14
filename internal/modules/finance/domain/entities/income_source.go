package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	vo "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/value_objects"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// IncomeSource represents a source of income for a user.
// Examples: salary from company X, freelance project Y, rental property Z.
type IncomeSource struct {
	ID         string
	UserID     string
	CategoryID *string // optional category

	Name        string
	Description string

	// Financial data
	Amount     sharedVO.Money
	IncomeType vo.IncomeSourceType
	Recurrence vo.Recurrence

	// Dates
	StartDate time.Time
	EndDate   *time.Time // nil = ongoing
	IsActive  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewIncomeSource creates a new income source with validation
func NewIncomeSource(
	userID string,
	name string,
	description string,
	amountCents int64,
	currency string,
	incomeType string,
	recurrence string,
	startDate time.Time,
	categoryID *string,
) (*IncomeSource, error) {
	if strings.TrimSpace(name) == "" {
		return nil, financeErrors.NewIncomeSourceInvalidError("income source name cannot be empty")
	}

	if amountCents <= 0 {
		return nil, financeErrors.NewIncomeSourceInvalidError("amount must be greater than zero")
	}

	money, err := sharedVO.NewMoney(amountCents, currency)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid currency: " + err.Error())
	}

	it, err := vo.NewIncomeSourceType(incomeType)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError(err.Error())
	}

	rec, err := vo.NewRecurrence(recurrence)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError(err.Error())
	}

	now := time.Now()
	return &IncomeSource{
		ID:          uuid.New().String(),
		UserID:      userID,
		CategoryID:  categoryID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Amount:      *money,
		IncomeType:  *it,
		Recurrence:  *rec,
		StartDate:   startDate,
		EndDate:     nil,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates the income source fields
func (i *IncomeSource) Update(
	name string,
	description string,
	amountCents int64,
	currency string,
	incomeType string,
	recurrence string,
	categoryID *string,
) error {
	if strings.TrimSpace(name) == "" {
		return financeErrors.NewIncomeSourceInvalidError("income source name cannot be empty")
	}

	if amountCents <= 0 {
		return financeErrors.NewIncomeSourceInvalidError("amount must be greater than zero")
	}

	money, err := sharedVO.NewMoney(amountCents, currency)
	if err != nil {
		return financeErrors.NewIncomeSourceInvalidError("invalid currency: " + err.Error())
	}

	it, err := vo.NewIncomeSourceType(incomeType)
	if err != nil {
		return financeErrors.NewIncomeSourceInvalidError(err.Error())
	}

	rec, err := vo.NewRecurrence(recurrence)
	if err != nil {
		return financeErrors.NewIncomeSourceInvalidError(err.Error())
	}

	i.Name = strings.TrimSpace(name)
	i.Description = strings.TrimSpace(description)
	i.Amount = *money
	i.IncomeType = *it
	i.Recurrence = *rec
	i.CategoryID = categoryID
	i.UpdatedAt = time.Now()

	return nil
}

// Deactivate marks the income source as inactive (ended)
func (i *IncomeSource) Deactivate(endDate time.Time) {
	i.IsActive = false
	i.EndDate = &endDate
	i.UpdatedAt = time.Now()
}

// Reactivate reactivates a previously inactive income source
func (i *IncomeSource) Reactivate() {
	i.IsActive = true
	i.EndDate = nil
	i.UpdatedAt = time.Now()
}

// BelongsTo checks if this income source belongs to the given user
func (i *IncomeSource) BelongsTo(userID string) bool {
	return i.UserID == userID
}

// MonthlyEquivalent returns the monthly equivalent of this income
// Used for dashboard calculations and comparisons
func (i *IncomeSource) MonthlyEquivalent() sharedVO.Money {
	amountCents := i.Amount.Amount()
	var monthlyAmount int64

	switch i.Recurrence.Type() {
	case vo.RecurrenceDaily:
		monthlyAmount = amountCents * 30
	case vo.RecurrenceWeekly:
		monthlyAmount = amountCents * 4
	case vo.RecurrenceBiweekly:
		monthlyAmount = amountCents * 2
	case vo.RecurrenceMonthly:
		monthlyAmount = amountCents
	case vo.RecurrenceQuarterly:
		monthlyAmount = amountCents / 3
	case vo.RecurrenceAnnually:
		monthlyAmount = amountCents / 12
	default:
		monthlyAmount = amountCents
	}

	money, _ := sharedVO.NewMoney(monthlyAmount, i.Amount.Currency())
	return *money
}
