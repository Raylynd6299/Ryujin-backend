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

// DebtService handles debt use cases
type DebtService struct {
	debtRepo repositories.DebtRepository
}

// NewDebtService creates a new DebtService
func NewDebtService(debtRepo repositories.DebtRepository) *DebtService {
	return &DebtService{debtRepo: debtRepo}
}

// CreateDebt creates a new debt record for a user
func (s *DebtService) CreateDebt(ctx context.Context, userID string, req dto.CreateDebtRequest) (*dto.DebtResponse, error) {
	totalCents := decimalToCents(req.TotalAmount)
	remainingCents := decimalToCents(req.RemainingAmount)
	monthlyCents := decimalToCents(req.MonthlyPayment)

	var startDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, financeErrors.NewDebtInvalidError("invalid startDate format, expected YYYY-MM-DD")
		}
		startDate = &t
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, financeErrors.NewDebtInvalidError("invalid dueDate format, expected YYYY-MM-DD")
		}
		dueDate = &t
	}

	debt, err := entities.NewDebt(
		userID,
		req.Name,
		req.Description,
		req.DebtType,
		totalCents,
		remainingCents,
		monthlyCents,
		req.Currency,
		req.InterestRate,
		startDate,
		dueDate,
	)
	if err != nil {
		return nil, err
	}

	if err := s.debtRepo.Create(ctx, debt); err != nil {
		return nil, fmt.Errorf("failed to create debt: %w", err)
	}

	return mapDebtToResponse(debt), nil
}

// GetDebt returns a debt by ID
func (s *DebtService) GetDebt(ctx context.Context, id string, userID string) (*dto.DebtResponse, error) {
	debt, err := s.debtRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return mapDebtToResponse(debt), nil
}

// ListDebts returns paginated debts for a user
func (s *DebtService) ListDebts(ctx context.Context, userID string, page, perPage int) (*dto.DebtListResponse, error) {
	pagination := utils.NormalizePagination(utils.Pagination{Page: page, PerPage: perPage})

	debts, total, err := s.debtRepo.FindAllByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list debts: %w", err)
	}

	responses := make([]*dto.DebtResponse, len(debts))
	for i, debt := range debts {
		responses[i] = mapDebtToResponse(debt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PerPage)))

	return &dto.DebtListResponse{
		Data:       responses,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateDebt updates debt metadata
func (s *DebtService) UpdateDebt(ctx context.Context, id string, userID string, req dto.UpdateDebtRequest) (*dto.DebtResponse, error) {
	debt, err := s.debtRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	monthlyCents := decimalToCents(req.MonthlyPayment)

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, financeErrors.NewDebtInvalidError("invalid dueDate format, expected YYYY-MM-DD")
		}
		dueDate = &t
	}

	if err := debt.Update(
		req.Name,
		req.Description,
		monthlyCents,
		req.Currency,
		req.InterestRate,
		dueDate,
	); err != nil {
		return nil, err
	}

	if err := s.debtRepo.Update(ctx, debt); err != nil {
		return nil, fmt.Errorf("failed to update debt: %w", err)
	}

	return mapDebtToResponse(debt), nil
}

// RecordPayment registers a payment and reduces remaining balance
func (s *DebtService) RecordPayment(ctx context.Context, id string, userID string, req dto.RecordPaymentRequest) (*dto.DebtResponse, error) {
	debt, err := s.debtRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	paymentCents := decimalToCents(req.PaymentAmount)

	if err := debt.RecordPayment(paymentCents); err != nil {
		return nil, err
	}

	if err := s.debtRepo.Update(ctx, debt); err != nil {
		return nil, fmt.Errorf("failed to record payment: %w", err)
	}

	return mapDebtToResponse(debt), nil
}

// DeleteDebt removes a debt record
func (s *DebtService) DeleteDebt(ctx context.Context, id string, userID string) error {
	_, err := s.debtRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.debtRepo.Delete(ctx, id, userID)
}

// --- Mapper ---

func mapDebtToResponse(d *entities.Debt) *dto.DebtResponse {
	return &dto.DebtResponse{
		ID:              d.ID,
		UserID:          d.UserID,
		Name:            d.Name,
		Description:     d.Description,
		DebtType:        d.DebtType.String(),
		TotalAmount:     centsToDecimal(d.TotalAmount.Amount()),
		RemainingAmount: centsToDecimal(d.RemainingAmount.Amount()),
		MonthlyPayment:  centsToDecimal(d.MonthlyPayment.Amount()),
		Currency:        d.TotalAmount.Currency(),
		InterestRate:    d.InterestRate,
		StartDate:       d.StartDate,
		DueDate:         d.DueDate,
		IsActive:        d.IsActive,
		ProgressPercent: d.ProgressPercent(),
		MonthsToPayoff:  d.MonthsToPayoff(),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
