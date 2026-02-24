package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Raylynd6299/ryujin/internal/modules/goal/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/goal/domain/entities"
	goalErrors "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/goal/domain/repositories"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// GoalService handles purchase goal and contribution use cases
type GoalService struct {
	goalRepo         repositories.GoalRepository
	contributionRepo repositories.GoalContributionRepository
}

// NewGoalService creates a new GoalService
func NewGoalService(
	goalRepo repositories.GoalRepository,
	contributionRepo repositories.GoalContributionRepository,
) *GoalService {
	return &GoalService{
		goalRepo:         goalRepo,
		contributionRepo: contributionRepo,
	}
}

// ---- Goal CRUD ----

// CreateGoal creates a new purchase goal for a user
func (s *GoalService) CreateGoal(ctx context.Context, userID string, req dto.CreateGoalRequest) (*dto.GoalResponse, error) {
	targetCents := decimalToCents(req.TargetAmount)

	var deadline *time.Time
	if req.Deadline != nil {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, goalErrors.NewGoalInvalidError("invalid deadline format, expected YYYY-MM-DD")
		}
		deadline = &t
	}

	goal, err := entities.NewPurchaseGoal(
		userID,
		req.Name,
		req.Description,
		req.Icon,
		targetCents,
		req.Currency,
		req.Priority,
		deadline,
	)
	if err != nil {
		return nil, err
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return s.buildGoalResponse(ctx, goal)
}

// GetGoal returns a goal by ID with its analytics
func (s *GoalService) GetGoal(ctx context.Context, id string, userID string) (*dto.GoalResponse, error) {
	goal, err := s.goalRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return s.buildGoalResponse(ctx, goal)
}

// ListGoals returns paginated goals for a user with analytics
func (s *GoalService) ListGoals(ctx context.Context, userID string, page, perPage int) (*dto.GoalListResponse, error) {
	pagination := utils.NormalizePagination(utils.Pagination{Page: page, PerPage: perPage})

	goals, total, err := s.goalRepo.FindAllByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list goals: %w", err)
	}

	responses := make([]*dto.GoalResponse, len(goals))
	for i, goal := range goals {
		resp, err := s.buildGoalResponse(ctx, goal)
		if err != nil {
			return nil, err
		}
		responses[i] = resp
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PerPage)))

	return &dto.GoalListResponse{
		Data:       responses,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateGoal updates a goal's metadata
func (s *GoalService) UpdateGoal(ctx context.Context, id string, userID string, req dto.UpdateGoalRequest) (*dto.GoalResponse, error) {
	goal, err := s.goalRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	targetCents := decimalToCents(req.TargetAmount)

	var deadline *time.Time
	if req.Deadline != nil {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, goalErrors.NewGoalInvalidError("invalid deadline format, expected YYYY-MM-DD")
		}
		deadline = &t
	}

	if err := goal.Update(
		req.Name,
		req.Description,
		req.Icon,
		targetCents,
		req.Currency,
		req.Priority,
		deadline,
	); err != nil {
		return nil, err
	}

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return s.buildGoalResponse(ctx, goal)
}

// DeleteGoal removes a goal and all its contributions
func (s *GoalService) DeleteGoal(ctx context.Context, id string, userID string) error {
	_, err := s.goalRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.goalRepo.Delete(ctx, id, userID)
}

// ---- Contribution CRUD ----

// AddContribution adds a monetary contribution to a goal
func (s *GoalService) AddContribution(ctx context.Context, goalID string, userID string, req dto.CreateContributionRequest) (*dto.ContributionResponse, error) {
	// Ensure goal exists and belongs to the user
	goal, err := s.goalRepo.FindByID(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	amountCents := decimalToCents(req.Amount)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, goalErrors.NewContributionInvalidError("invalid date format, expected YYYY-MM-DD")
	}

	contribution, err := entities.NewGoalContribution(
		goal.ID,
		userID,
		amountCents,
		req.Currency,
		date,
		req.Notes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.contributionRepo.Create(ctx, contribution); err != nil {
		return nil, fmt.Errorf("failed to create contribution: %w", err)
	}

	// Check if goal is now over-funded and auto-complete it
	contributions, err := s.contributionRepo.FindAllByGoalID(ctx, goal.ID, userID)
	if err == nil && goal.IsOverFunded(contributions) && !goal.IsCompleted {
		goal.MarkCompleted()
		_ = s.goalRepo.Update(ctx, goal) // best-effort — don't fail the contribution
	}

	return mapContributionToResponse(contribution), nil
}

// ListContributions returns all contributions for a goal
func (s *GoalService) ListContributions(ctx context.Context, goalID string, userID string) (*dto.ContributionListResponse, error) {
	// Ensure goal exists and belongs to the user
	_, err := s.goalRepo.FindByID(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	contributions, err := s.contributionRepo.FindAllByGoalID(ctx, goalID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list contributions: %w", err)
	}

	responses := make([]*dto.ContributionResponse, len(contributions))
	for i, c := range contributions {
		responses[i] = mapContributionToResponse(c)
	}

	return &dto.ContributionListResponse{
		Data:  responses,
		Total: len(responses),
	}, nil
}

// DeleteContribution removes a contribution from a goal
func (s *GoalService) DeleteContribution(ctx context.Context, goalID string, contributionID string, userID string) error {
	// Ensure goal exists and belongs to the user
	_, err := s.goalRepo.FindByID(ctx, goalID, userID)
	if err != nil {
		return err
	}

	contribution, err := s.contributionRepo.FindByID(ctx, contributionID, userID)
	if err != nil {
		return err
	}

	if !contribution.BelongsToGoal(goalID) {
		return goalErrors.NewContributionInvalidError("contribution does not belong to this goal")
	}

	return s.contributionRepo.Delete(ctx, contributionID, userID)
}

// ---- private helpers ----

// buildGoalResponse fetches contributions and computes analytics for a goal
func (s *GoalService) buildGoalResponse(ctx context.Context, goal *entities.PurchaseGoal) (*dto.GoalResponse, error) {
	contributions, err := s.contributionRepo.FindAllByGoalID(ctx, goal.ID, goal.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load contributions for goal %s: %w", goal.ID, err)
	}

	totalContributed := goal.TotalContributed(contributions)
	missingAmount := goal.MissingAmount(contributions)
	progressPercent := goal.ProgressPercent(contributions)
	estimatedCompletion := goal.EstimatedCompletionDate(contributions)

	return &dto.GoalResponse{
		ID:                  goal.ID,
		UserID:              goal.UserID,
		Name:                goal.Name,
		Description:         goal.Description,
		Icon:                goal.Icon,
		TargetAmount:        centsToDecimal(goal.TargetAmount.Amount()),
		Currency:            goal.TargetAmount.Currency(),
		Priority:            goal.Priority.String(),
		Deadline:            goal.Deadline,
		IsCompleted:         goal.IsCompleted,
		TotalContributed:    centsToDecimal(totalContributed.Amount()),
		ProgressPercent:     progressPercent,
		MissingAmount:       centsToDecimal(missingAmount.Amount()),
		EstimatedCompletion: estimatedCompletion,
		CreatedAt:           goal.CreatedAt,
		UpdatedAt:           goal.UpdatedAt,
	}, nil
}

// mapContributionToResponse maps a domain contribution to a DTO response
func mapContributionToResponse(c *entities.GoalContribution) *dto.ContributionResponse {
	return &dto.ContributionResponse{
		ID:        c.ID,
		GoalID:    c.GoalID,
		UserID:    c.UserID,
		Amount:    centsToDecimal(c.Amount.Amount()),
		Currency:  c.Amount.Currency(),
		Date:      c.Date,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// decimalToCents converts a float64 dollar amount to integer cents
func decimalToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

// centsToDecimal converts integer cents to float64
func centsToDecimal(cents int64) float64 {
	return float64(cents) / 100.0
}
