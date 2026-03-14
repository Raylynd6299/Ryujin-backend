package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/value_objects"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// User represents a user account in the system.
// This is the aggregate root for the User bounded context.
type User struct {
	// Identity
	ID string // UUID as string

	// Authentication
	Email          value_objects.Email
	HashedPassword value_objects.HashedPassword

	// Profile
	FirstName string
	LastName  string

	// Preferences
	DefaultSavingsCurrency    sharedVO.Currency
	DefaultInvestmentCurrency sharedVO.Currency
	Locale                    value_objects.Locale

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewUser creates a new User with validation of invariants.
func NewUser(
	email value_objects.Email,
	password value_objects.Password,
	firstName string,
	lastName string,
	locale value_objects.Locale,
) (*User, error) {
	// Validate first name
	if strings.TrimSpace(firstName) == "" {
		return nil, errors.NewInvalidUserError("first name cannot be empty")
	}

	// Validate last name
	if strings.TrimSpace(lastName) == "" {
		return nil, errors.NewInvalidUserError("last name cannot be empty")
	}

	// Hash password
	hashedPassword, err := value_objects.HashPassword(password)
	if err != nil {
		return nil, errors.NewInvalidUserError("failed to hash password: " + err.Error())
	}

	// Default currencies (USD for both)
	defaultCurrency, _ := sharedVO.NewCurrency("USD")

	now := time.Now()

	return &User{
		ID:                        uuid.New().String(),
		Email:                     email,
		HashedPassword:            hashedPassword,
		FirstName:                 firstName,
		LastName:                  lastName,
		DefaultSavingsCurrency:    *defaultCurrency,
		DefaultInvestmentCurrency: *defaultCurrency,
		Locale:                    locale,
		CreatedAt:                 now,
		UpdatedAt:                 now,
		DeletedAt:                 nil,
	}, nil
}

// UpdateProfile updates user profile information.
func (u *User) UpdateProfile(firstName string, lastName string, locale value_objects.Locale) error {
	// Validate first name
	if strings.TrimSpace(firstName) == "" {
		return errors.NewInvalidUserError("first name cannot be empty")
	}

	// Validate last name
	if strings.TrimSpace(lastName) == "" {
		return errors.NewInvalidUserError("last name cannot be empty")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.Locale = locale
	u.UpdatedAt = time.Now()

	return nil
}

// UpdateCurrencies updates default currencies.
func (u *User) UpdateCurrencies(savings sharedVO.Currency, investment sharedVO.Currency) error {
	u.DefaultSavingsCurrency = savings
	u.DefaultInvestmentCurrency = investment
	u.UpdatedAt = time.Now()
	return nil
}

// ChangePassword changes the user's password after verifying the old one.
func (u *User) ChangePassword(oldPassword value_objects.Password, newPassword value_objects.Password) error {
	// Verify old password
	if !u.HashedPassword.CompareWith(oldPassword) {
		return errors.NewInvalidPasswordError("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := value_objects.HashPassword(newPassword)
	if err != nil {
		return errors.NewInvalidPasswordError("failed to hash new password: " + err.Error())
	}

	u.HashedPassword = hashedPassword
	u.UpdatedAt = time.Now()

	return nil
}

// VerifyPassword checks if the provided password matches the user's password.
func (u *User) VerifyPassword(password value_objects.Password) bool {
	return u.HashedPassword.CompareWith(password)
}

// GetFullName returns the user's full name.
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// SoftDelete marks the user as deleted.
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// IsDeleted checks if the user is soft-deleted.
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsActive checks if the user is not deleted.
func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}
