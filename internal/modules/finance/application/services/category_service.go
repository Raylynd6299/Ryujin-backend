package services

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/repositories"
)

// CategoryService handles category use cases
type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(categoryRepo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory creates a new user-defined category
func (s *CategoryService) CreateCategory(ctx context.Context, userID string, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := entities.NewCategory(
		userID,
		req.Name,
		entities.CategoryType(req.Type),
		req.Icon,
		req.Color,
	)
	if err != nil {
		return nil, err
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return mapCategoryToResponse(category), nil
}

// GetCategory returns a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id string, userID string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !category.BelongsTo(userID) {
		return nil, financeErrors.NewUnauthorizedError("you do not have access to this category")
	}

	return mapCategoryToResponse(category), nil
}

// ListCategories returns all categories available to the user (system + user-owned)
func (s *CategoryService) ListCategories(ctx context.Context, userID string) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = mapCategoryToResponse(c)
	}

	return responses, nil
}

// UpdateCategory updates a user-defined category
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, userID string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !category.BelongsTo(userID) {
		return nil, financeErrors.NewUnauthorizedError("you do not have access to this category")
	}

	if err := category.Update(req.Name, req.Icon, req.Color); err != nil {
		return nil, err
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return mapCategoryToResponse(category), nil
}

// DeleteCategory removes a user-defined category
func (s *CategoryService) DeleteCategory(ctx context.Context, id string, userID string) error {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !category.BelongsTo(userID) {
		return financeErrors.NewUnauthorizedError("you do not have access to this category")
	}

	if category.IsDefault {
		return financeErrors.NewCategoryInvalidError("cannot delete a system default category")
	}

	return s.categoryRepo.Delete(ctx, id, userID)
}

// --- Mapper ---

func mapCategoryToResponse(c *entities.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Type:      string(c.Type),
		Icon:      c.Icon,
		Color:     c.Color,
		IsDefault: c.IsDefault,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
