package services

import (
	"context"
	"fmt"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/repositories"
)

// AccountService handles account use cases
type AccountService struct {
	accountRepo repositories.AccountRepository
}

// NewAccountService creates a new AccountService
func NewAccountService(accountRepo repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

// CreateAccount creates a new financial account for a user
func (s *AccountService) CreateAccount(ctx context.Context, userID string, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	balanceCents := decimalToCents(req.Balance)

	account, err := entities.NewAccount(
		userID,
		req.Name,
		req.Description,
		entities.AccountType(req.AccountType),
		balanceCents,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return mapAccountToResponse(account), nil
}

// GetAccount returns an account by ID
func (s *AccountService) GetAccount(ctx context.Context, id string, userID string) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return mapAccountToResponse(account), nil
}

// ListAccounts returns all accounts for a user
func (s *AccountService) ListAccounts(ctx context.Context, userID string) ([]*dto.AccountResponse, error) {
	accounts, err := s.accountRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	responses := make([]*dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = mapAccountToResponse(account)
	}

	return responses, nil
}

// UpdateAccount updates account metadata
func (s *AccountService) UpdateAccount(ctx context.Context, id string, userID string, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if err := account.Update(req.Name, req.Description, entities.AccountType(req.AccountType)); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return mapAccountToResponse(account), nil
}

// UpdateBalance reconciles the account balance manually
func (s *AccountService) UpdateBalance(ctx context.Context, id string, userID string, req dto.UpdateBalanceRequest) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	balanceCents := decimalToCents(req.Balance)

	if err := account.UpdateBalance(balanceCents, req.Currency); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return mapAccountToResponse(account), nil
}

// DeactivateAccount marks an account as closed/inactive
func (s *AccountService) DeactivateAccount(ctx context.Context, id string, userID string) (*dto.AccountResponse, error) {
	account, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	account.Deactivate()

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to deactivate account: %w", err)
	}

	return mapAccountToResponse(account), nil
}

// DeleteAccount removes an account
func (s *AccountService) DeleteAccount(ctx context.Context, id string, userID string) error {
	_, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.accountRepo.Delete(ctx, id, userID)
}

// --- Mapper ---

func mapAccountToResponse(a *entities.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:          a.ID,
		UserID:      a.UserID,
		Name:        a.Name,
		Description: a.Description,
		AccountType: string(a.AccountType),
		Balance:     centsToDecimal(a.Balance.Amount()),
		Currency:    a.Balance.Currency(),
		IsActive:    a.IsActive,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
