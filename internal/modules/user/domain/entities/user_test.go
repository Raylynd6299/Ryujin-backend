package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	entityVO "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/domain/value_objects"
	sharedVO "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/value_objects"
)

func TestNewUserValid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")

	user, err := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, entityVO.LocaleSpanish, user.Locale)
	assert.NotEmpty(t, user.HashedPassword.String())
	assert.False(t, user.IsDeleted())
	assert.True(t, user.IsActive())
}

func TestNewUserInvalidFirstName(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")

	_, err := NewUser(email, password, "", "Doe", entityVO.LocaleSpanish)

	assert.Error(t, err)
}

func TestNewUserInvalidLastName(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")

	_, err := NewUser(email, password, "John", "", entityVO.LocaleSpanish)
	assert.Error(t, err)
}

func TestUpdateProfileValid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	err := user.UpdateProfile("Jane", "Smith", entityVO.LocaleEnglish)

	assert.NoError(t, err)
	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, entityVO.LocaleEnglish, user.Locale)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestUpdateProfileInvalidFirstName(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	err := user.UpdateProfile("", "Smith", entityVO.LocaleEnglish)

	assert.Error(t, err)
}

func TestChangePasswordValid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	oldPassword, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, oldPassword, "John", "Doe", entityVO.LocaleSpanish)

	newPassword, _ := entityVO.NewPassword("NewSecure456")
	err := user.ChangePassword(oldPassword, newPassword)

	assert.NoError(t, err)
	assert.True(t, user.VerifyPassword(newPassword))
	assert.False(t, user.VerifyPassword(oldPassword))
}

func TestChangePasswordWrongOldPassword(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	wrongOldPassword, _ := entityVO.NewPassword("WrongPass789")
	newPassword, _ := entityVO.NewPassword("NewSecure456")

	err := user.ChangePassword(wrongOldPassword, newPassword)

	assert.Error(t, err)
	assert.True(t, user.VerifyPassword(password))
}

func TestVerifyPasswordValid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	assert.True(t, user.VerifyPassword(password))
}

func TestVerifyPasswordInvalid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)
	wrongPassword, _ := entityVO.NewPassword("WrongPass789")
	assert.False(t, user.VerifyPassword(wrongPassword))
}

func TestGetFullName(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	assert.Equal(t, "John Doe", user.GetFullName())
}

func TestSoftDelete(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	assert.False(t, user.IsDeleted())
	assert.True(t, user.IsActive())
	assert.Nil(t, user.DeletedAt)

	user.SoftDelete()

	assert.True(t, user.IsDeleted())
	assert.False(t, user.IsActive())
	assert.NotNil(t, user.DeletedAt)
}

func TestUpdateCurrenciesValid(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	savingsCurrency, _ := sharedVO.NewCurrency("MXN")
	investmentCurrency, _ := sharedVO.NewCurrency("USD")

	err := user.UpdateCurrencies(*savingsCurrency, *investmentCurrency)

	assert.NoError(t, err)
	assert.Equal(t, *savingsCurrency, user.DefaultSavingsCurrency)
	assert.Equal(t, *investmentCurrency, user.DefaultInvestmentCurrency)
}

func TestDefaultCurrencies(t *testing.T) {
	email, _ := entityVO.NewEmail("user@example.com")
	password, _ := entityVO.NewPassword("SecurePass123")
	user, _ := NewUser(email, password, "John", "Doe", entityVO.LocaleSpanish)

	// Defaults should be USD for both
	defaultCurrency, _ := sharedVO.NewCurrency("USD")
	assert.Equal(t, *defaultCurrency, user.DefaultSavingsCurrency)
	assert.Equal(t, *defaultCurrency, user.DefaultInvestmentCurrency)
}
