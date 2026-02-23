package dto

import "time"

// --- Request DTOs ---

// CreateIncomeSourceRequest is used to create a new income source
type CreateIncomeSourceRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"` // decimal input (e.g. 1500.00)
	Currency    string  `json:"currency" binding:"required,len=3"`
	IncomeType  string  `json:"incomeType" binding:"required,oneof=salary freelance rental dividend business other"`
	Recurrence  string  `json:"recurrence" binding:"required,oneof=none daily weekly biweekly monthly quarterly annually"`
	StartDate   string  `json:"startDate" binding:"required"` // ISO 8601 date string
	CategoryID  *string `json:"categoryId"`
}

// UpdateIncomeSourceRequest is used to update an existing income source
type UpdateIncomeSourceRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3"`
	IncomeType  string  `json:"incomeType" binding:"required,oneof=salary freelance rental dividend business other"`
	Recurrence  string  `json:"recurrence" binding:"required,oneof=none daily weekly biweekly monthly quarterly annually"`
	CategoryID  *string `json:"categoryId"`
}

// DeactivateIncomeSourceRequest is used to stop a recurring income source
type DeactivateIncomeSourceRequest struct {
	EndDate string `json:"endDate" binding:"required"` // ISO 8601 date string
}

// --- Response DTOs ---

// IncomeSourceResponse is returned for income source operations
type IncomeSourceResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	CategoryID  *string    `json:"categoryId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"` // converted from cents to decimal
	Currency    string     `json:"currency"`
	IncomeType  string     `json:"incomeType"`
	Recurrence  string     `json:"recurrence"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	IsActive    bool       `json:"isActive"`
	// Computed field for dashboard aggregation
	MonthlyEquivalent float64   `json:"monthlyEquivalent"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// IncomeSourceListResponse wraps a paginated list of income sources
type IncomeSourceListResponse struct {
	Data       []*IncomeSourceResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"perPage"`
	TotalPages int                     `json:"totalPages"`
}
