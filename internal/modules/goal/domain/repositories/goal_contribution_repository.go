package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/goal/domain/entities"
)

// GoalContributionRepository defines the port for contribution persistence
type GoalContributionRepository interface {
	// Create persists a new contribution
	Create(ctx context.Context, contribution *entities.GoalContribution) error

	// FindByID returns a contribution by ID scoped to the user
	FindByID(ctx context.Context, id string, userID string) (*entities.GoalContribution, error)

	// FindAllByGoalID returns all contributions for a specific goal scoped to the user
	FindAllByGoalID(ctx context.Context, goalID string, userID string) ([]*entities.GoalContribution, error)

	// Delete removes a contribution record
	Delete(ctx context.Context, id string, userID string) error

	// SumByGoalID returns the total contributed cents for a goal
	SumByGoalID(ctx context.Context, goalID string, userID string) (int64, error)
}
