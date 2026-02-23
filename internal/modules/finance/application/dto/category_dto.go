package dto

import "time"

// --- Request DTOs ---

// CreateCategoryRequest is used when a user creates a custom category
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Type  string `json:"type" binding:"required,oneof=income expense both"`
	Icon  string `json:"icon" binding:"max=50"`
	Color string `json:"color" binding:"max=20"`
}

// UpdateCategoryRequest is used when a user updates a category
type UpdateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Icon  string `json:"icon" binding:"max=50"`
	Color string `json:"color" binding:"max=20"`
}

// --- Response DTOs ---

// CategoryResponse is returned for category operations
type CategoryResponse struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"userId,omitempty"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
