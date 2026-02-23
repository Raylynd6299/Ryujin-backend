package dto

import "time"

// --- Request DTOs ---

// CreateAccountRequest is used to create a new financial account
type CreateAccountRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=500"`
	AccountType string  `json:"accountType" binding:"required,oneof=checking savings cash wallet"`
	Balance     float64 `json:"balance" binding:"gte=0"`
	Currency    string  `json:"currency" binding:"required,len=3"`
}

// UpdateAccountRequest is used to update account metadata
type UpdateAccountRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=150"`
	Description string `json:"description" binding:"max=500"`
	AccountType string `json:"accountType" binding:"required,oneof=checking savings cash wallet"`
}

// UpdateBalanceRequest is used for manual balance reconciliation
type UpdateBalanceRequest struct {
	Balance  float64 `json:"balance" binding:"gte=0"`
	Currency string  `json:"currency" binding:"required,len=3"`
}

// --- Response DTOs ---

// AccountResponse is returned for account operations
type AccountResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AccountType string    `json:"accountType"`
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
