package dto

import "time"

// UserResponse is the public representation of a user — never exposes sensitive fields.
type UserResponse struct {
	ID                        string    `json:"id"`
	Email                     string    `json:"email"`
	FirstName                 string    `json:"firstName"`
	LastName                  string    `json:"lastName"`
	DefaultSavingsCurrency    string    `json:"defaultSavingsCurrency"`
	DefaultInvestmentCurrency string    `json:"defaultInvestmentCurrency"`
	Locale                    string    `json:"locale"`
	CreatedAt                 time.Time `json:"createdAt"`
}

// UpdateProfileRequest contains the fields the user can update on their profile.
type UpdateProfileRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"  binding:"required"`
	Locale    string `json:"locale"`
}

// UpdateCurrenciesRequest allows updating default currencies.
type UpdateCurrenciesRequest struct {
	DefaultSavingsCurrency    string `json:"defaultSavingsCurrency"    binding:"required,len=3"`
	DefaultInvestmentCurrency string `json:"defaultInvestmentCurrency" binding:"required,len=3"`
}

// ChangePasswordRequest contains the old and new passwords.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword"     binding:"required,min=8"`
}
