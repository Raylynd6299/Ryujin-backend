package dto

import "time"

// --- Request DTOs ---

// CreateExpenseRequest is used to create a new expense
type CreateExpenseRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3"`
	Priority    string  `json:"priority" binding:"required,oneof=essential important optional low"`
	Recurrence  string  `json:"recurrence" binding:"required,oneof=none daily weekly biweekly monthly quarterly annually"`
	ExpenseDate string  `json:"expenseDate" binding:"required"` // ISO 8601 date string
	CategoryID  *string `json:"categoryId"`
}

// UpdateExpenseRequest is used to update an existing expense
type UpdateExpenseRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required,len=3"`
	Priority    string  `json:"priority" binding:"required,oneof=essential important optional low"`
	Recurrence  string  `json:"recurrence" binding:"required,oneof=none daily weekly biweekly monthly quarterly annually"`
	CategoryID  *string `json:"categoryId"`
}

// --- Response DTOs ---

// ExpenseResponse is returned for expense operations
type ExpenseResponse struct {
	ID            string     `json:"id"`
	UserID        string     `json:"userId"`
	CategoryID    *string    `json:"categoryId,omitempty"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	Priority      string     `json:"priority"`
	Recurrence    string     `json:"recurrence"`
	ExpenseDate   time.Time  `json:"expenseDate"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	IsActive      bool       `json:"isActive"`
	IsUnnecessary bool       `json:"isUnnecessary"`
	// Computed field for dashboard aggregation
	MonthlyEquivalent float64   `json:"monthlyEquivalent"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// ExpenseListResponse wraps a paginated list of expenses
type ExpenseListResponse struct {
	Data       []*ExpenseResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"perPage"`
	TotalPages int                `json:"totalPages"`
}
