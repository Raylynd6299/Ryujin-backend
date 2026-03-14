package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/entities"
	userErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/repositories"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/value_objects"
	sharedErrors "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// AuthService handles authentication use cases: register, login, refresh token.
type AuthService struct {
	userRepo             repositories.UserRepository
	jwtService           *utils.JWTService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewAuthService creates a new AuthService with its required dependencies.
func NewAuthService(
	userRepo repositories.UserRepository,
	jwtService *utils.JWTService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		jwtService:           jwtService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// Register creates a new user account and returns the user with auth tokens.
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Build and validate value objects
	email, err := value_objects.NewEmail(req.Email)
	if err != nil {
		return nil, sharedErrors.NewValidationError("INVALID_EMAIL", err.Error())
	}

	password, err := value_objects.NewPassword(req.Password)
	if err != nil {
		return nil, sharedErrors.NewValidationError("INVALID_PASSWORD", err.Error())
	}

	locale, err := value_objects.NewLocale(req.Locale)
	if err != nil {
		// Default to English if locale is invalid or empty
		locale, _ = value_objects.NewLocale("en")
	}

	// Check email uniqueness
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("auth service: checking email availability: %w", err)
	}
	if exists {
		return nil, userErrors.NewDuplicateEmailError(email.String())
	}

	// Create domain entity (hashes password internally)
	user, err := entities.NewUser(email, password, req.FirstName, req.LastName, locale)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("auth service: persisting user: %w", err)
	}

	// Issue tokens
	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, fmt.Errorf("auth service: issuing tokens: %w", err)
	}

	return &dto.AuthResponse{
		User:   toUserResponse(user),
		Tokens: tokens,
	}, nil
}

// Login authenticates a user and returns the user with fresh tokens.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email, err := value_objects.NewEmail(req.Email)
	if err != nil {
		return nil, sharedErrors.NewUnauthorizedError("invalid credentials")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Never reveal whether email exists or not
		return nil, sharedErrors.NewUnauthorizedError("invalid credentials")
	}

	password, err := value_objects.NewPassword(req.Password)
	if err != nil {
		return nil, sharedErrors.NewUnauthorizedError("invalid credentials")
	}

	if !user.VerifyPassword(password) {
		return nil, sharedErrors.NewUnauthorizedError("invalid credentials")
	}

	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, fmt.Errorf("auth service: issuing tokens: %w", err)
	}

	return &dto.AuthResponse{
		User:   toUserResponse(user),
		Tokens: tokens,
	}, nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	userID, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, sharedErrors.NewUnauthorizedError("invalid or expired refresh token")
	}

	// Verify user still exists and is active
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, sharedErrors.NewUnauthorizedError("invalid or expired refresh token")
	}

	tokens, err := s.issueTokens(user)
	if err != nil {
		return nil, fmt.Errorf("auth service: issuing tokens: %w", err)
	}

	return &tokens, nil
}

// issueTokens generates a new access + refresh token pair for a user.
func (s *AuthService) issueTokens(user *entities.User) (dto.TokenPair, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email.String(), s.accessTokenDuration)
	if err != nil {
		return dto.TokenPair{}, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, s.refreshTokenDuration)
	if err != nil {
		return dto.TokenPair{}, fmt.Errorf("generating refresh token: %w", err)
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
