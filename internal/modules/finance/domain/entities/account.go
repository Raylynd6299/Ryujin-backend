package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// AccountType classifies the kind of financial account
type AccountType string

const (
	AccountTypeChecking AccountType = "checking" // Bank checking account
	AccountTypeSavings  AccountType = "savings"  // Bank savings account
	AccountTypeCash     AccountType = "cash"     // Physical cash
	AccountTypeWallet   AccountType = "wallet"   // Digital wallet (PayPal, etc.)
)

// Account represents a financial account owned by a user.
// Used to track liquid assets and current balances.
type Account struct {
	ID     string
	UserID string

	Name        string
	Description string
	AccountType AccountType

	// Financial data
	Balance sharedVO.Money

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAccount creates a new account with validation
func NewAccount(
	userID string,
	name string,
	description string,
	accountType AccountType,
	balanceCents int64,
	currency string,
) (*Account, error) {
	if strings.TrimSpace(name) == "" {
		return nil, financeErrors.NewAccountInvalidError("account name cannot be empty")
	}

	validTypes := map[AccountType]bool{
		AccountTypeChecking: true,
		AccountTypeSavings:  true,
		AccountTypeCash:     true,
		AccountTypeWallet:   true,
	}
	if !validTypes[accountType] {
		return nil, financeErrors.NewAccountInvalidError("invalid account type")
	}

	money, err := sharedVO.NewMoney(balanceCents, currency)
	if err != nil {
		return nil, financeErrors.NewAccountInvalidError("invalid currency: " + err.Error())
	}

	now := time.Now()
	return &Account{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		AccountType: accountType,
		Balance:     *money,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates account metadata
func (a *Account) Update(name string, description string, accountType AccountType) error {
	if strings.TrimSpace(name) == "" {
		return financeErrors.NewAccountInvalidError("account name cannot be empty")
	}

	validTypes := map[AccountType]bool{
		AccountTypeChecking: true,
		AccountTypeSavings:  true,
		AccountTypeCash:     true,
		AccountTypeWallet:   true,
	}
	if !validTypes[accountType] {
		return financeErrors.NewAccountInvalidError("invalid account type")
	}

	a.Name = strings.TrimSpace(name)
	a.Description = strings.TrimSpace(description)
	a.AccountType = accountType
	a.UpdatedAt = time.Now()

	return nil
}

// UpdateBalance directly sets a new balance (for manual reconciliation)
func (a *Account) UpdateBalance(balanceCents int64, currency string) error {
	money, err := sharedVO.NewMoney(balanceCents, currency)
	if err != nil {
		return financeErrors.NewAccountInvalidError("invalid currency: " + err.Error())
	}
	a.Balance = *money
	a.UpdatedAt = time.Now()
	return nil
}

// Deactivate marks the account as closed/inactive
func (a *Account) Deactivate() {
	a.IsActive = false
	a.UpdatedAt = time.Now()
}

// BelongsTo checks if this account belongs to the given user
func (a *Account) BelongsTo(userID string) bool {
	return a.UserID == userID
}
