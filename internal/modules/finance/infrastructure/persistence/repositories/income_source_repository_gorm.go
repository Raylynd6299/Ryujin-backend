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

// IncomeSourceRepositoryGorm implements domain/repositories.IncomeSourceRepository using GORM
type IncomeSourceRepositoryGorm struct {
	db *gorm.DB
}

func NewIncomeSourceRepositoryGorm(db *gorm.DB) *IncomeSourceRepositoryGorm {
	return &IncomeSourceRepositoryGorm{db: db}
}

func (r *IncomeSourceRepositoryGorm) Create(ctx context.Context, income *entities.IncomeSource) error {
	model := mappers.IncomeSourceToModel(income)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *IncomeSourceRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.IncomeSource, error) {
	var model models.IncomeSourceModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, financeErrors.NewIncomeSourceNotFoundError(id)
		}
		return nil, err
	}
	return mappers.IncomeSourceToDomain(&model)
}

func (r *IncomeSourceRepositoryGorm) FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.IncomeSource, int64, error) {
	var modelList []models.IncomeSourceModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.IncomeSourceModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entities.IncomeSource, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.IncomeSourceToDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domain)
	}
	return result, total, nil
}

func (r *IncomeSourceRepositoryGorm) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.IncomeSource, error) {
	var modelList []models.IncomeSourceModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.IncomeSource, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.IncomeSourceToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *IncomeSourceRepositoryGorm) Update(ctx context.Context, income *entities.IncomeSource) error {
	model := mappers.IncomeSourceToModel(income)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *IncomeSourceRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.IncomeSourceModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return financeErrors.NewIncomeSourceNotFoundError(id)
	}
	return nil
}
