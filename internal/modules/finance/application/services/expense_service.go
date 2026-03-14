package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/repositories"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// ExpenseService handles expense use cases
type ExpenseService struct {
	expenseRepo repositories.ExpenseRepository
}

// NewExpenseService creates a new ExpenseService
func NewExpenseService(expenseRepo repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

// CreateExpense creates a new expense for a user
func (s *ExpenseService) CreateExpense(ctx context.Context, userID string, req dto.CreateExpenseRequest) (*dto.ExpenseResponse, error) {
	amountCents := decimalToCents(req.Amount)

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError("invalid expenseDate format, expected YYYY-MM-DD")
	}

	expense, err := entities.NewExpense(
		userID,
		req.Name,
		req.Description,
		amountCents,
		req.Currency,
		req.Priority,
		req.Recurrence,
		expenseDate,
		req.CategoryID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return mapExpenseToResponse(expense), nil
}

// GetExpense returns an expense by ID
func (s *ExpenseService) GetExpense(ctx context.Context, id string, userID string) (*dto.ExpenseResponse, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return mapExpenseToResponse(expense), nil
}

// ListExpenses returns paginated expenses for a user
func (s *ExpenseService) ListExpenses(ctx context.Context, userID string, page, perPage int) (*dto.ExpenseListResponse, error) {
	pagination := utils.NormalizePagination(utils.Pagination{Page: page, PerPage: perPage})

	expenses, total, err := s.expenseRepo.FindAllByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %w", err)
	}

	responses := make([]*dto.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		responses[i] = mapExpenseToResponse(expense)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PerPage)))

	return &dto.ExpenseListResponse{
		Data:       responses,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateExpense updates an existing expense
func (s *ExpenseService) UpdateExpense(ctx context.Context, id string, userID string, req dto.UpdateExpenseRequest) (*dto.ExpenseResponse, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	amountCents := decimalToCents(req.Amount)

	if err := expense.Update(
		req.Name,
		req.Description,
		amountCents,
		req.Currency,
		req.Priority,
		req.Recurrence,
		req.CategoryID,
	); err != nil {
		return nil, err
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	return mapExpenseToResponse(expense), nil
}

// DeleteExpense removes an expense
func (s *ExpenseService) DeleteExpense(ctx context.Context, id string, userID string) error {
	_, err := s.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.expenseRepo.Delete(ctx, id, userID)
}

// --- Mapper ---

func mapExpenseToResponse(e *entities.Expense) *dto.ExpenseResponse {
	monthly := e.MonthlyEquivalent()
	return &dto.ExpenseResponse{
		ID:                e.ID,
		UserID:            e.UserID,
		CategoryID:        e.CategoryID,
		Name:              e.Name,
		Description:       e.Description,
		Amount:            centsToDecimal(e.Amount.Amount()),
		Currency:          e.Amount.Currency(),
		Priority:          e.Priority.String(),
		Recurrence:        e.Recurrence.String(),
		ExpenseDate:       e.ExpenseDate,
		EndDate:           e.EndDate,
		IsActive:          e.IsActive,
		IsUnnecessary:     e.IsUnnecessary(),
		MonthlyEquivalent: centsToDecimal(monthly.Amount()),
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}
