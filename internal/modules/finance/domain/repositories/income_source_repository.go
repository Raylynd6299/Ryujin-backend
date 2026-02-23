package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// IncomeSourceRepository defines the port for income source persistence
type IncomeSourceRepository interface {
	// Create persists a new income source
	Create(ctx context.Context, income *entities.IncomeSource) error

	// FindByID returns an income source by ID
	FindByID(ctx context.Context, id string, userID string) (*entities.IncomeSource, error)

	// FindAllByUserID returns all income sources for a user with pagination
	FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.IncomeSource, int64, error)

	// FindActiveByUserID returns only active income sources
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.IncomeSource, error)

	// Update persists changes to an existing income source
	Update(ctx context.Context, income *entities.IncomeSource) error

	// Delete removes an income source
	Delete(ctx context.Context, id string, userID string) error
}
