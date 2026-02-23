package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
)

// CategoryRepository defines the port for category persistence
type CategoryRepository interface {
	// Create persists a new category
	Create(ctx context.Context, category *entities.Category) error

	// FindByID returns a category by ID
	FindByID(ctx context.Context, id string) (*entities.Category, error)

	// FindAllByUserID returns user-owned + system categories
	FindAllByUserID(ctx context.Context, userID string) ([]*entities.Category, error)

	// FindSystemCategories returns all system/global categories
	FindSystemCategories(ctx context.Context) ([]*entities.Category, error)

	// Update persists changes to an existing category
	Update(ctx context.Context, category *entities.Category) error

	// Delete removes a user-defined category
	Delete(ctx context.Context, id string, userID string) error
}
