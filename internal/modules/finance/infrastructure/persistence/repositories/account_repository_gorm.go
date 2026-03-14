package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/persistence/models"
)

// AccountRepositoryGorm implements domain/repositories.AccountRepository using GORM
type AccountRepositoryGorm struct {
	db *gorm.DB
}

func NewAccountRepositoryGorm(db *gorm.DB) *AccountRepositoryGorm {
	return &AccountRepositoryGorm{db: db}
}

func (r *AccountRepositoryGorm) Create(ctx context.Context, account *entities.Account) error {
	model := mappers.AccountToModel(account)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *AccountRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.Account, error) {
	var model models.AccountModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, financeErrors.NewAccountNotFoundError(id)
		}
		return nil, err
	}
	return mappers.AccountToDomain(&model)
}

func (r *AccountRepositoryGorm) FindAllByUserID(ctx context.Context, userID string) ([]*entities.Account, error) {
	var modelList []models.AccountModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Account, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.AccountToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *AccountRepositoryGorm) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Account, error) {
	var modelList []models.AccountModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Order("name ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Account, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.AccountToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *AccountRepositoryGorm) Update(ctx context.Context, account *entities.Account) error {
	model := mappers.AccountToModel(account)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *AccountRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.AccountModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return financeErrors.NewAccountNotFoundError(id)
	}
	return nil
}
