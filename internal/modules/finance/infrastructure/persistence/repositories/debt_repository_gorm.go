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

// DebtRepositoryGorm implements domain/repositories.DebtRepository using GORM
type DebtRepositoryGorm struct {
	db *gorm.DB
}

func NewDebtRepositoryGorm(db *gorm.DB) *DebtRepositoryGorm {
	return &DebtRepositoryGorm{db: db}
}

func (r *DebtRepositoryGorm) Create(ctx context.Context, debt *entities.Debt) error {
	model := mappers.DebtToModel(debt)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *DebtRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.Debt, error) {
	var model models.DebtModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, financeErrors.NewDebtNotFoundError(id)
		}
		return nil, err
	}
	return mappers.DebtToDomain(&model)
}

func (r *DebtRepositoryGorm) FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.Debt, int64, error) {
	var modelList []models.DebtModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DebtModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("is_active DESC, created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entities.Debt, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.DebtToDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domain)
	}
	return result, total, nil
}

func (r *DebtRepositoryGorm) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Debt, error) {
	var modelList []models.DebtModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("remaining_amount_cents DESC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Debt, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.DebtToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *DebtRepositoryGorm) Update(ctx context.Context, debt *entities.Debt) error {
	model := mappers.DebtToModel(debt)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *DebtRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.DebtModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return financeErrors.NewDebtNotFoundError(id)
	}
	return nil
}
