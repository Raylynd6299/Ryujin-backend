package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/repositories"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// IncomeSourceService handles income source use cases
type IncomeSourceService struct {
	incomeRepo repositories.IncomeSourceRepository
}

// NewIncomeSourceService creates a new IncomeSourceService
func NewIncomeSourceService(incomeRepo repositories.IncomeSourceRepository) *IncomeSourceService {
	return &IncomeSourceService{incomeRepo: incomeRepo}
}

// CreateIncomeSource creates a new income source for a user
func (s *IncomeSourceService) CreateIncomeSource(ctx context.Context, userID string, req dto.CreateIncomeSourceRequest) (*dto.IncomeSourceResponse, error) {
	// Convert decimal amount to cents
	amountCents := decimalToCents(req.Amount)

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid startDate format, expected YYYY-MM-DD")
	}

	income, err := entities.NewIncomeSource(
		userID,
		req.Name,
		req.Description,
		amountCents,
		req.Currency,
		req.IncomeType,
		req.Recurrence,
		startDate,
		req.CategoryID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.incomeRepo.Create(ctx, income); err != nil {
		return nil, fmt.Errorf("failed to create income source: %w", err)
	}

	return mapIncomeSourceToResponse(income), nil
}

// GetIncomeSource returns an income source by ID
func (s *IncomeSourceService) GetIncomeSource(ctx context.Context, id string, userID string) (*dto.IncomeSourceResponse, error) {
	income, err := s.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return mapIncomeSourceToResponse(income), nil
}

// ListIncomeSources returns paginated income sources for a user
func (s *IncomeSourceService) ListIncomeSources(ctx context.Context, userID string, page, perPage int) (*dto.IncomeSourceListResponse, error) {
	pagination := utils.NormalizePagination(utils.Pagination{Page: page, PerPage: perPage})

	incomes, total, err := s.incomeRepo.FindAllByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list income sources: %w", err)
	}

	responses := make([]*dto.IncomeSourceResponse, len(incomes))
	for i, income := range incomes {
		responses[i] = mapIncomeSourceToResponse(income)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PerPage)))

	return &dto.IncomeSourceListResponse{
		Data:       responses,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateIncomeSource updates an existing income source
func (s *IncomeSourceService) UpdateIncomeSource(ctx context.Context, id string, userID string, req dto.UpdateIncomeSourceRequest) (*dto.IncomeSourceResponse, error) {
	income, err := s.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	amountCents := decimalToCents(req.Amount)

	if err := income.Update(
		req.Name,
		req.Description,
		amountCents,
		req.Currency,
		req.IncomeType,
		req.Recurrence,
		req.CategoryID,
	); err != nil {
		return nil, err
	}

	if err := s.incomeRepo.Update(ctx, income); err != nil {
		return nil, fmt.Errorf("failed to update income source: %w", err)
	}

	return mapIncomeSourceToResponse(income), nil
}

// DeactivateIncomeSource stops a recurring income source
func (s *IncomeSourceService) DeactivateIncomeSource(ctx context.Context, id string, userID string, req dto.DeactivateIncomeSourceRequest) (*dto.IncomeSourceResponse, error) {
	income, err := s.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid endDate format, expected YYYY-MM-DD")
	}

	income.Deactivate(endDate)

	if err := s.incomeRepo.Update(ctx, income); err != nil {
		return nil, fmt.Errorf("failed to deactivate income source: %w", err)
	}

	return mapIncomeSourceToResponse(income), nil
}

// DeleteIncomeSource removes an income source
func (s *IncomeSourceService) DeleteIncomeSource(ctx context.Context, id string, userID string) error {
	_, err := s.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.incomeRepo.Delete(ctx, id, userID)
}

// --- Mapper ---

func mapIncomeSourceToResponse(i *entities.IncomeSource) *dto.IncomeSourceResponse {
	monthly := i.MonthlyEquivalent()
	return &dto.IncomeSourceResponse{
		ID:                i.ID,
		UserID:            i.UserID,
		CategoryID:        i.CategoryID,
		Name:              i.Name,
		Description:       i.Description,
		Amount:            centsToDecimal(i.Amount.Amount()),
		Currency:          i.Amount.Currency(),
		IncomeType:        i.IncomeType.String(),
		Recurrence:        i.Recurrence.String(),
		StartDate:         i.StartDate,
		EndDate:           i.EndDate,
		IsActive:          i.IsActive,
		MonthlyEquivalent: centsToDecimal(monthly.Amount()),
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}
}
