package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/goal/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// GoalRepository defines the port for purchase goal persistence
type GoalRepository interface {
	// Create persists a new purchase goal
	Create(ctx context.Context, goal *entities.PurchaseGoal) error

	// FindByID returns a goal by ID scoped to the user
	FindByID(ctx context.Context, id string, userID string) (*entities.PurchaseGoal, error)

	// FindAllByUserID returns paginated goals for a user
	FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.PurchaseGoal, int64, error)

	// Update persists changes to a goal
	Update(ctx context.Context, goal *entities.PurchaseGoal) error

	// Delete removes a goal record
	Delete(ctx context.Context, id string, userID string) error
}
