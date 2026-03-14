package mappers

import (
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/domain/entities"
	goalErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/domain/errors"
	vo "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/domain/value_objects"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/persistence/models"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// ============================================================
// PurchaseGoal
// ============================================================

func PurchaseGoalToModel(g *entities.PurchaseGoal) *models.PurchaseGoalModel {
	return &models.PurchaseGoalModel{
		ID:                g.ID,
		UserID:            g.UserID,
		Name:              g.Name,
		Description:       g.Description,
		Icon:              g.Icon,
		TargetAmountCents: g.TargetAmount.Amount(),
		Currency:          g.TargetAmount.Currency(),
		Priority:          g.Priority.String(),
		Deadline:          g.Deadline,
		IsCompleted:       g.IsCompleted,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}
}

func PurchaseGoalToDomain(m *models.PurchaseGoalModel) (*entities.PurchaseGoal, error) {
	targetMoney, err := sharedVO.NewMoney(m.TargetAmountCents, m.Currency)
	if err != nil {
		return nil, goalErrors.NewGoalInvalidError("invalid money in DB: " + err.Error())
	}

	p, err := vo.NewGoalPriority(m.Priority)
	if err != nil {
		return nil, goalErrors.NewGoalInvalidError("invalid priority in DB: " + err.Error())
	}

	return &entities.PurchaseGoal{
		ID:           m.ID,
		UserID:       m.UserID,
		Name:         m.Name,
		Description:  m.Description,
		Icon:         m.Icon,
		TargetAmount: *targetMoney,
		Priority:     *p,
		Deadline:     m.Deadline,
		IsCompleted:  m.IsCompleted,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

// ============================================================
// GoalContribution
// ============================================================

func GoalContributionToModel(c *entities.GoalContribution) *models.GoalContributionModel {
	return &models.GoalContributionModel{
		ID:          c.ID,
		GoalID:      c.GoalID,
		UserID:      c.UserID,
		AmountCents: c.Amount.Amount(),
		Currency:    c.Amount.Currency(),
		Date:        c.Date,
		Notes:       c.Notes,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func GoalContributionToDomain(m *models.GoalContributionModel) (*entities.GoalContribution, error) {
	money, err := sharedVO.NewMoney(m.AmountCents, m.Currency)
	if err != nil {
		return nil, goalErrors.NewContributionInvalidError("invalid money in DB: " + err.Error())
	}

	return &entities.GoalContribution{
		ID:        m.ID,
		GoalID:    m.GoalID,
		UserID:    m.UserID,
		Amount:    *money,
		Date:      m.Date,
		Notes:     m.Notes,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
