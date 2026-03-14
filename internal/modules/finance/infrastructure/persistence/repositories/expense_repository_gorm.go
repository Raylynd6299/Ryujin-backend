package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/persistence/models"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// ExpenseRepositoryGorm implements domain/repositories.ExpenseRepository using GORM
type ExpenseRepositoryGorm struct {
	db *gorm.DB
}

func NewExpenseRepositoryGorm(db *gorm.DB) *ExpenseRepositoryGorm {
	return &ExpenseRepositoryGorm{db: db}
}

func (r *ExpenseRepositoryGorm) Create(ctx context.Context, expense *entities.Expense) error {
	model := mappers.ExpenseToModel(expense)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ExpenseRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.Expense, error) {
	var model models.ExpenseModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, financeErrors.NewExpenseNotFoundError(id)
		}
		return nil, err
	}
	return mappers.ExpenseToDomain(&model)
}

func (r *ExpenseRepositoryGorm) FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.Expense, int64, error) {
	var modelList []models.ExpenseModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ExpenseModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("expense_date DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entities.Expense, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.ExpenseToDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domain)
	}
	return result, total, nil
}

func (r *ExpenseRepositoryGorm) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Expense, error) {
	var modelList []models.ExpenseModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Expense, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.ExpenseToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *ExpenseRepositoryGorm) Update(ctx context.Context, expense *entities.Expense) error {
	model := mappers.ExpenseToModel(expense)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *ExpenseRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.ExpenseModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return financeErrors.NewExpenseNotFoundError(id)
	}
	return nil
}
