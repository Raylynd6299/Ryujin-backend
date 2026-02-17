package services

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/value_objects"
)

// UserService defines domain operations that require repository access.
// This service encapsulates domain logic that doesn't belong to a single entity.
type UserService interface {
	// IsEmailAvailable checks if an email can be registered.
	// Returns true if email is not taken by any active user.
	// Returns false if email is already taken (not an error condition).
	// Returns error only if repository operation fails.
	IsEmailAvailable(ctx context.Context, email value_objects.Email) (bool, error)

	// GetUserByID retrieves a user by ID.
	// Centralizes user retrieval at domain level.
	GetUserByID(ctx context.Context, id string) (*entities.User, error)
}
