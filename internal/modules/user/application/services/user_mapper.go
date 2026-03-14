package services

import (
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/entities"
)

// toUserResponse converts a domain User entity into a public UserResponse DTO.
// Never exposes sensitive fields like hashed_password.
func toUserResponse(user *entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:                        user.ID,
		Email:                     user.Email.String(),
		FirstName:                 user.FirstName,
		LastName:                  user.LastName,
		DefaultSavingsCurrency:    user.DefaultSavingsCurrency.String(),
		DefaultInvestmentCurrency: user.DefaultInvestmentCurrency.String(),
		Locale:                    user.Locale.String(),
		CreatedAt:                 user.CreatedAt,
	}
}
