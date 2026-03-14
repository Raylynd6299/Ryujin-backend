package repositories

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/entities"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/value_objects"
)

// UserRepository defines the contract for user persistence.
// This is a port in the hexagonal architecture - the domain defines what it needs,
// and the infrastructure provides the implementation.
type UserRepository interface {
	// Create persists a new user to the database.
	// Returns error if email already exists or database operation fails.
	Create(ctx context.Context, user *entities.User) error

	// FindByID retrieves a user by their ID.
	// Excludes soft-deleted users from results.
	// Returns NotFoundError if user doesn't exist or is deleted.
	FindByID(ctx context.Context, id string) (*entities.User, error)

	// FindByEmail retrieves a user by their email.
	// Excludes soft-deleted users from results.
	// Returns NotFoundError if user doesn't exist or is deleted.
	FindByEmail(ctx context.Context, email value_objects.Email) (*entities.User, error)

	// Update persists changes to an existing user.
	// Returns NotFoundError if user doesn't exist.
	Update(ctx context.Context, user *entities.User) error

	// Delete soft-deletes a user by setting DeletedAt timestamp.
	// Returns NotFoundError if user doesn't exist.
	Delete(ctx context.Context, id string) error

	// ExistsByEmail checks if an email is already taken by an active user.
	// Excludes soft-deleted users.
	// Returns true if email exists, false otherwise.
	ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error)

	// FindAll retrieves a paginated list of active (non-deleted) users.
	// page: 1-indexed page number
	// pageSize: items per page
	// Returns (users, totalCount, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entities.User, int, error)
}
