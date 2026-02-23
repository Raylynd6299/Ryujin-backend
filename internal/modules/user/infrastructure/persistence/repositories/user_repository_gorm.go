package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/entities"
	userErrors "github.com/Raylynd6299/ryujin/internal/modules/user/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/value_objects"
	"github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/persistence/models"
)

// UserRepositoryGorm implements the domain UserRepository port using GORM.
type UserRepositoryGorm struct {
	db *gorm.DB
}

// NewUserRepositoryGorm creates a new GORM-backed user repository.
func NewUserRepositoryGorm(db *gorm.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{db: db}
}

// Create persists a new user to the database.
func (r *UserRepositoryGorm) Create(ctx context.Context, user *entities.User) error {
	model := mappers.ToModel(user)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return fmt.Errorf("user repository: creating user: %w", result.Error)
	}
	return nil
}

// FindByID retrieves an active user by their ID.
func (r *UserRepositoryGorm) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var model models.UserModel
	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, userErrors.NewUserNotFoundError("user not found with id: " + id)
		}
		return nil, fmt.Errorf("user repository: finding user by id: %w", result.Error)
	}

	return mappers.ToDomain(&model), nil
}

// FindByEmail retrieves an active user by their email.
func (r *UserRepositoryGorm) FindByEmail(ctx context.Context, email value_objects.Email) (*entities.User, error) {
	var model models.UserModel
	result := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email.String()).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, userErrors.NewUserNotFoundError("user not found with email: " + email.String())
		}
		return nil, fmt.Errorf("user repository: finding user by email: %w", result.Error)
	}

	return mappers.ToDomain(&model), nil
}

// Update persists changes to an existing user.
func (r *UserRepositoryGorm) Update(ctx context.Context, user *entities.User) error {
	model := mappers.ToModel(user)
	result := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ? AND deleted_at IS NULL", model.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("user repository: updating user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return userErrors.NewUserNotFoundError("user not found with id: " + user.ID)
	}

	return nil
}

// Delete soft-deletes a user by setting deleted_at.
func (r *UserRepositoryGorm) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		return fmt.Errorf("user repository: deleting user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return userErrors.NewUserNotFoundError("user not found with id: " + id)
	}

	return nil
}

// ExistsByEmail checks if an active user with the given email exists.
func (r *UserRepositoryGorm) ExistsByEmail(ctx context.Context, email value_objects.Email) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("email = ? AND deleted_at IS NULL", email.String()).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("user repository: checking email existence: %w", result.Error)
	}

	return count > 0, nil
}

// FindAll retrieves a paginated list of active users.
func (r *UserRepositoryGorm) FindAll(ctx context.Context, page, pageSize int) ([]*entities.User, int, error) {
	var dbModels []models.UserModel
	var total int64

	offset := (page - 1) * pageSize

	// Count total active users
	if err := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("user repository: counting users: %w", err)
	}

	// Fetch page
	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&dbModels).Error; err != nil {
		return nil, 0, fmt.Errorf("user repository: listing users: %w", err)
	}

	users := make([]*entities.User, 0, len(dbModels))
	for i := range dbModels {
		users = append(users, mappers.ToDomain(&dbModels[i]))
	}

	return users, int(total), nil
}
