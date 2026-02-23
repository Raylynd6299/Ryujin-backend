package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/infrastructure/persistence/models"
)

// CategoryRepositoryGorm implements domain/repositories.CategoryRepository using GORM
type CategoryRepositoryGorm struct {
	db *gorm.DB
}

func NewCategoryRepositoryGorm(db *gorm.DB) *CategoryRepositoryGorm {
	return &CategoryRepositoryGorm{db: db}
}

func (r *CategoryRepositoryGorm) Create(ctx context.Context, category *entities.Category) error {
	model := mappers.CategoryToModel(category)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepositoryGorm) FindByID(ctx context.Context, id string) (*entities.Category, error) {
	var model models.CategoryModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, financeErrors.NewCategoryNotFoundError(id)
		}
		return nil, err
	}
	return mappers.CategoryToDomain(&model), nil
}

func (r *CategoryRepositoryGorm) FindAllByUserID(ctx context.Context, userID string) ([]*entities.Category, error) {
	var modelList []models.CategoryModel
	// Returns system categories (user_id IS NULL) + user's own categories
	err := r.db.WithContext(ctx).
		Where("user_id IS NULL OR user_id = ?", userID).
		Order("is_default DESC, name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Category, len(modelList))
	for i, m := range modelList {
		m := m // avoid loop variable capture
		result[i] = mappers.CategoryToDomain(&m)
	}
	return result, nil
}

func (r *CategoryRepositoryGorm) FindSystemCategories(ctx context.Context) ([]*entities.Category, error) {
	var modelList []models.CategoryModel
	err := r.db.WithContext(ctx).
		Where("user_id IS NULL AND is_default = true").
		Order("name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Category, len(modelList))
	for i, m := range modelList {
		m := m
		result[i] = mappers.CategoryToDomain(&m)
	}
	return result, nil
}

func (r *CategoryRepositoryGorm) Update(ctx context.Context, category *entities.Category) error {
	model := mappers.CategoryToModel(category)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CategoryRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.CategoryModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return financeErrors.NewCategoryNotFoundError(id)
	}
	return nil
}
