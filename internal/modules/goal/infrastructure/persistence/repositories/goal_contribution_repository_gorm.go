package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/domain/entities"
	goalErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/persistence/models"
)

// GoalContributionRepositoryGorm implements domain/repositories.GoalContributionRepository using GORM
type GoalContributionRepositoryGorm struct {
	db *gorm.DB
}

func NewGoalContributionRepositoryGorm(db *gorm.DB) *GoalContributionRepositoryGorm {
	return &GoalContributionRepositoryGorm{db: db}
}

func (r *GoalContributionRepositoryGorm) Create(ctx context.Context, contribution *entities.GoalContribution) error {
	model := mappers.GoalContributionToModel(contribution)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GoalContributionRepositoryGorm) FindByID(ctx context.Context, id string, userID string) (*entities.GoalContribution, error) {
	var model models.GoalContributionModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, goalErrors.NewContributionNotFoundError(id)
		}
		return nil, err
	}
	return mappers.GoalContributionToDomain(&model)
}

func (r *GoalContributionRepositoryGorm) FindAllByGoalID(ctx context.Context, goalID string, userID string) ([]*entities.GoalContribution, error) {
	var modelList []models.GoalContributionModel
	err := r.db.WithContext(ctx).
		Where("goal_id = ? AND user_id = ?", goalID, userID).
		Order("date DESC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.GoalContribution, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.GoalContributionToDomain(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

func (r *GoalContributionRepositoryGorm) Delete(ctx context.Context, id string, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.GoalContributionModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return goalErrors.NewContributionNotFoundError(id)
	}
	return nil
}

func (r *GoalContributionRepositoryGorm) SumByGoalID(ctx context.Context, goalID string, userID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.GoalContributionModel{}).
		Where("goal_id = ? AND user_id = ?", goalID, userID).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&total).Error
	return total, err
}
