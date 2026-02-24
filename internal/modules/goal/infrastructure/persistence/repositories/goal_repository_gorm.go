package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/modules/goal/domain/entities"
	goalErrors "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/goal/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/ryujin/internal/modules/goal/infrastructure/persistence/models"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// GoalRepositoryGorm implements domain/repositories.GoalRepository using GORM
type GoalRepositoryGorm struct {
	db *gorm.DB
}

func NewGoalRepositoryGorm(db *gorm.DB) *GoalRepositoryGorm {
	return &GoalRepositoryGorm{db: db}
}

func (r *GoalRepositoryGorm) Create(ctx context.Context, goal *entities.PurchaseGoal) error {
	model := mappers.PurchaseGoalToModel(goal)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GoalRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.PurchaseGoal, error) {
	var model models.PurchaseGoalModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, goalErrors.NewGoalNotFoundError(id)
		}
		return nil, err
	}
	return mappers.PurchaseGoalToDomain(&model)
}

func (r *GoalRepositoryGorm) FindAllByUserID(ctx context.Context, userID string, params utils.Pagination) ([]*entities.PurchaseGoal, int64, error) {
	var modelList []models.PurchaseGoalModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PurchaseGoalModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("is_completed ASC, priority DESC, created_at DESC").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entities.PurchaseGoal, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.PurchaseGoalToDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domain)
	}
	return result, total, nil
}

func (r *GoalRepositoryGorm) Update(ctx context.Context, goal *entities.PurchaseGoal) error {
	model := mappers.PurchaseGoalToModel(goal)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GoalRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	// Cascade delete contributions first
	r.db.WithContext(ctx).
		Where("goal_id = ? AND user_id = ?", id, userID).
		Delete(&models.GoalContributionModel{})

	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.PurchaseGoalModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return goalErrors.NewGoalNotFoundError(id)
	}
	return nil
}
