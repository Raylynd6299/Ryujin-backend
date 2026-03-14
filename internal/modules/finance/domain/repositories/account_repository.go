package repositories

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
)

// AccountRepository defines the port for account persistence
type AccountRepository interface {
	// Create persists a new account
	Create(ctx context.Context, account *entities.Account) error

	// FindByID returns an account by ID
	FindByID(ctx context.Context, id string, userID string) (*entities.Account, error)

	// FindAllByUserID returns all accounts for a user
	FindAllByUserID(ctx context.Context, userID string) ([]*entities.Account, error)

	// FindActiveByUserID returns only active accounts
	FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Account, error)

	// Update persists changes to an account
	Update(ctx context.Context, account *entities.Account) error

	// Delete removes an account
	Delete(ctx context.Context, id string, userID string) error
}
