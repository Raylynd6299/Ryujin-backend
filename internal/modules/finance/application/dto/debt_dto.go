package dto

import "time"

// --- Request DTOs ---

// CreateDebtRequest is used to create a new debt record
type CreateDebtRequest struct {
	Name            string  `json:"name" binding:"required,min=1,max=150"`
	Description     string  `json:"description" binding:"max=500"`
	DebtType        string  `json:"debtType" binding:"required,oneof=credit_card personal_loan mortgage car_loan student_loan other"`
	TotalAmount     float64 `json:"totalAmount" binding:"required,gt=0"`
	RemainingAmount float64 `json:"remainingAmount" binding:"required,gte=0"`
	MonthlyPayment  float64 `json:"monthlyPayment" binding:"required,gt=0"`
	Currency        string  `json:"currency" binding:"required,len=3"`
	InterestRate    float64 `json:"interestRate" binding:"gte=0"`
	StartDate       *string `json:"startDate"` // ISO 8601 date string, optional
	DueDate         *string `json:"dueDate"`   // ISO 8601 date string, optional
}

// UpdateDebtRequest is used to update debt metadata
type UpdateDebtRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=150"`
	Description    string  `json:"description" binding:"max=500"`
	MonthlyPayment float64 `json:"monthlyPayment" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required,len=3"`
	InterestRate   float64 `json:"interestRate" binding:"gte=0"`
	DueDate        *string `json:"dueDate"` // ISO 8601 date string, optional
}

// RecordPaymentRequest is used to register a payment against a debt
type RecordPaymentRequest struct {
	PaymentAmount float64 `json:"paymentAmount" binding:"required,gt=0"`
}

// --- Response DTOs ---

// DebtResponse is returned for debt operations
type DebtResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	DebtType        string     `json:"debtType"`
	TotalAmount     float64    `json:"totalAmount"`
	RemainingAmount float64    `json:"remainingAmount"`
	MonthlyPayment  float64    `json:"monthlyPayment"`
	Currency        string     `json:"currency"`
	InterestRate    float64    `json:"interestRate"`
	StartDate       *time.Time `json:"startDate,omitempty"`
	DueDate         *time.Time `json:"dueDate,omitempty"`
	IsActive        bool       `json:"isActive"`
	// Computed fields
	ProgressPercent float64   `json:"progressPercent"`
	MonthsToPayoff  int       `json:"monthsToPayoff"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// DebtListResponse wraps a paginated list of debts
type DebtListResponse struct {
	Data       []*DebtResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	TotalPages int             `json:"totalPages"`
}
