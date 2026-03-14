package repositories

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// ExpenseRepository defines the port for expense persistence
type ExpenseRepository interface {
	// Create persists a new expense
	Create(ctx context.Context, expense *entities.Expense) error

	// FindByID returns an expense by ID
	FindByID(ctx context.Context, id string, userID string) (*entities.Expense, error)

	// FindAllByUserID returns all expenses for a user with pagination
	FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.Expense, int64, error)

	// FindActiveByUserID returns only active (recurring) expenses
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Expense, error)

	// Update persists changes to an existing expense
	Update(ctx context.Context, expense *entities.Expense) error

	// Delete removes an expense
	Delete(ctx context.Context, id string, userID string) error
}
