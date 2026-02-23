package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
)

// HoldingRepository defines the port for holding persistence
type HoldingRepository interface {
	// Create persists a new holding
	Create(ctx context.Context, h *entities.Holding) error

	// FindByID returns a holding by ID scoped to the given user.
	// Returns HoldingNotFoundError if not found or if the holding belongs to a different user.
	FindByID(ctx context.Context, id, userID string) (*entities.Holding, error)

	// FindAllByUserID returns paginated holdings for a user with total count
	FindAllByUserID(ctx context.Context, userID string, page, limit int, sort, order string) ([]*entities.Holding, int64, error)

	// FindActiveByUserID returns all holdings for a user (no pagination)
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Holding, error)

	// Update persists changes to an existing holding
	Update(ctx context.Context, h *entities.Holding) error

	// Delete removes a holding scoped to the given user
	Delete(ctx context.Context, id, userID string) error
}
