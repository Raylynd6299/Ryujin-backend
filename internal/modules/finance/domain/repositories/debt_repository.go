package repositories

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// DebtRepository defines the port for debt persistence
type DebtRepository interface {
	// Create persists a new debt
	Create(ctx context.Context, debt *entities.Debt) error

	// FindByID returns a debt by ID
	FindByID(ctx context.Context, id string, userID string) (*entities.Debt, error)

	// FindAllByUserID returns all debts for a user with pagination
	FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.Debt, int64, error)

	// FindActiveByUserID returns only active (not fully paid) debts
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Debt, error)

	// Update persists changes to a debt (balance, monthly payment, etc.)
	Update(ctx context.Context, debt *entities.Debt) error

	// Delete removes a debt record
	Delete(ctx context.Context, id string, userID string) error
}
