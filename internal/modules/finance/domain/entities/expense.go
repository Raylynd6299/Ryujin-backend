package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	vo "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/value_objects"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// Expense represents a financial outflow for a user.
// Can be one-time or recurring (rent, subscriptions, etc.)
type Expense struct {
	ID         string
	UserID     string
	CategoryID *string // optional category

	Name        string
	Description string

	// Financial data
	Amount     sharedVO.Money
	Priority   vo.Priority
	Recurrence vo.Recurrence

	// Dates
	ExpenseDate time.Time  // when the expense occurred / starts
	EndDate     *time.Time // nil = one-time or ongoing recurring
	IsActive    bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewExpense creates a new expense with validation
func NewExpense(
	userID string,
	name string,
	description string,
	amountCents int64,
	currency string,
	priority string,
	recurrence string,
	expenseDate time.Time,
	categoryID *string,
) (*Expense, error) {
	if strings.TrimSpace(name) == "" {
		return nil, financeErrors.NewExpenseInvalidError("expense name cannot be empty")
	}

	if amountCents <= 0 {
		return nil, financeErrors.NewExpenseInvalidError("amount must be greater than zero")
	}

	money, err := sharedVO.NewMoney(amountCents, currency)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError("invalid currency: " + err.Error())
	}

	prio, err := vo.NewPriority(priority)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError(err.Error())
	}

	rec, err := vo.NewRecurrence(recurrence)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError(err.Error())
	}

	now := time.Now()
	return &Expense{
		ID:          uuid.New().String(),
		UserID:      userID,
		CategoryID:  categoryID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Amount:      *money,
		Priority:    *prio,
		Recurrence:  *rec,
		ExpenseDate: expenseDate,
		EndDate:     nil,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates expense fields
func (e *Expense) Update(
	name string,
	description string,
	amountCents int64,
	currency string,
	priority string,
	recurrence string,
	categoryID *string,
) error {
	if strings.TrimSpace(name) == "" {
		return financeErrors.NewExpenseInvalidError("expense name cannot be empty")
	}

	if amountCents <= 0 {
		return financeErrors.NewExpenseInvalidError("amount must be greater than zero")
	}

	money, err := sharedVO.NewMoney(amountCents, currency)
	if err != nil {
		return financeErrors.NewExpenseInvalidError("invalid currency: " + err.Error())
	}

	prio, err := vo.NewPriority(priority)
	if err != nil {
		return financeErrors.NewExpenseInvalidError(err.Error())
	}

	rec, err := vo.NewRecurrence(recurrence)
	if err != nil {
		return financeErrors.NewExpenseInvalidError(err.Error())
	}

	e.Name = strings.TrimSpace(name)
	e.Description = strings.TrimSpace(description)
	e.Amount = *money
	e.Priority = *prio
	e.Recurrence = *rec
	e.CategoryID = categoryID
	e.UpdatedAt = time.Now()

	return nil
}

// Deactivate stops a recurring expense
func (e *Expense) Deactivate(endDate time.Time) {
	e.IsActive = false
	e.EndDate = &endDate
	e.UpdatedAt = time.Now()
}

// BelongsTo checks if this expense belongs to the given user
func (e *Expense) BelongsTo(userID string) bool {
	return e.UserID == userID
}

// IsUnnecessary returns true if this expense is optional or low priority
func (e *Expense) IsUnnecessary() bool {
	return e.Priority.IsUnnecessary()
}

// IsRecurring returns true if this is a recurring expense
func (e *Expense) IsRecurring() bool {
	return e.Recurrence.IsRecurring()
}

// MonthlyEquivalent returns the monthly cost of this expense
func (e *Expense) MonthlyEquivalent() sharedVO.Money {
	amountCents := e.Amount.Amount()
	var monthlyAmount int64

	switch e.Recurrence.Type() {
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
	default: // one-time
		monthlyAmount = amountCents
	}

	money, _ := sharedVO.NewMoney(monthlyAmount, e.Amount.Currency())
	return *money
}
