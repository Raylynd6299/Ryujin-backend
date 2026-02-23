package services

import (
	"context"
	"fmt"

	"github.com/Raylynd6299/ryujin/internal/modules/user/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/repositories"
	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/value_objects"
	sharedErrors "github.com/Raylynd6299/ryujin/internal/shared/domain/errors"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// ProfileService handles user profile management use cases.
type ProfileService struct {
	userRepo repositories.UserRepository
}

// NewProfileService creates a new ProfileService.
func NewProfileService(userRepo repositories.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

// GetProfile returns the public profile of a user by ID.
func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}

// UpdateProfile updates name and locale fields.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	locale, err := value_objects.NewLocale(req.Locale)
	if err != nil {
		locale = value_objects.DefaultLocale()
	}

	if err := user.UpdateProfile(req.FirstName, req.LastName, locale); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("profile service: persisting profile update: %w", err)
	}

	response := toUserResponse(user)
	return &response, nil
}

// UpdateCurrencies updates the user's default currencies.
func (s *ProfileService) UpdateCurrencies(ctx context.Context, userID string, req dto.UpdateCurrenciesRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	savings, err := sharedVO.NewCurrency(req.DefaultSavingsCurrency)
	if err != nil {
		return nil, sharedErrors.NewValidationError("INVALID_CURRENCY", "invalid savings currency: "+err.Error())
	}

	investment, err := sharedVO.NewCurrency(req.DefaultInvestmentCurrency)
	if err != nil {
		return nil, sharedErrors.NewValidationError("INVALID_CURRENCY", "invalid investment currency: "+err.Error())
	}

	if err := user.UpdateCurrencies(*savings, *investment); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("profile service: persisting currencies update: %w", err)
	}

	response := toUserResponse(user)
	return &response, nil
}

// ChangePassword verifies the old password and sets a new one.
func (s *ProfileService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	currentPassword, err := value_objects.NewPassword(req.CurrentPassword)
	if err != nil {
		return sharedErrors.NewUnauthorizedError("current password is incorrect")
	}

	newPassword, err := value_objects.NewPassword(req.NewPassword)
	if err != nil {
		return sharedErrors.NewValidationError("INVALID_PASSWORD", err.Error())
	}

	if err := user.ChangePassword(currentPassword, newPassword); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("profile service: persisting password change: %w", err)
	}

	return nil
}
