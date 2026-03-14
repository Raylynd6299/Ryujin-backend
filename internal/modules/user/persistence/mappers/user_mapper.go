package mappers

import (
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/entities"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/value_objects"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/persistence/models"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

// ToDomain converts a GORM UserModel into a domain User entity.
// If any value object is invalid (shouldn't happen for persisted data), it panics —
// data integrity is the DB's responsibility at this layer.
func ToDomain(m *models.UserModel) *entities.User {
	email, _ := value_objects.NewEmail(m.Email)
	locale, _ := value_objects.NewLocale(m.Locale)

	savingsCurrency, _ := sharedVO.NewCurrency(m.DefaultSavingsCurrency)
	investmentCurrency, _ := sharedVO.NewCurrency(m.DefaultInvestmentCurrency)

	return &entities.User{
		ID:                        m.ID,
		Email:                     email,
		HashedPassword:            value_objects.NewHashedPassword(m.HashedPassword),
		FirstName:                 m.FirstName,
		LastName:                  m.LastName,
		DefaultSavingsCurrency:    *savingsCurrency,
		DefaultInvestmentCurrency: *investmentCurrency,
		Locale:                    locale,
		CreatedAt:                 m.CreatedAt,
		UpdatedAt:                 m.UpdatedAt,
		DeletedAt:                 m.DeletedAt,
	}
}

// ToModel converts a domain User entity into a GORM UserModel.
func ToModel(u *entities.User) *models.UserModel {
	return &models.UserModel{
		ID:                        u.ID,
		Email:                     u.Email.String(),
		HashedPassword:            u.HashedPassword.String(),
		FirstName:                 u.FirstName,
		LastName:                  u.LastName,
		DefaultSavingsCurrency:    u.DefaultSavingsCurrency.String(),
		DefaultInvestmentCurrency: u.DefaultInvestmentCurrency.String(),
		Locale:                    u.Locale.String(),
		CreatedAt:                 u.CreatedAt,
		UpdatedAt:                 u.UpdatedAt,
		DeletedAt:                 u.DeletedAt,
	}
}
