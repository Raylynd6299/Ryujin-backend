package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/modules/user/domain/value_objects"
)

// TestUserDomainIntegration tests the complete user domain layer working together.
func TestUserDomainIntegration(t *testing.T) {
	// Create value objects
	email, err := value_objects.NewEmail("john.doe@example.com")
	assert.NoError(t, err)

	password, err := value_objects.NewPassword("SecurePass123")
	assert.NoError(t, err)

	locale := value_objects.LocaleSpanish

	// Create user entity
	user, err := entities.NewUser(email, password, "John", "Doe", locale)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Verify initial state
	assert.Equal(t, "John Doe", user.GetFullName())
	assert.True(t, user.IsActive())
	assert.False(t, user.IsDeleted())
	assert.True(t, user.VerifyPassword(password))

	// Update profile
	newLocale := value_objects.LocaleEnglish
	err = user.UpdateProfile("Jane", "Smith", newLocale)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Smith", user.GetFullName())
	assert.Equal(t, newLocale, user.Locale)

	// Change password
	newPassword, err := value_objects.NewPassword("NewSecure456")
	assert.NoError(t, err)

	err = user.ChangePassword(password, newPassword)
	assert.NoError(t, err)
	assert.True(t, user.VerifyPassword(newPassword))
	assert.False(t, user.VerifyPassword(password))

	// Soft delete
	user.SoftDelete()
	assert.True(t, user.IsDeleted())
	assert.False(t, user.IsActive())
	assert.NotNil(t, user.DeletedAt)
}

// TestUserDomainErrorHandling tests that domain layer properly handles errors.
func TestUserDomainErrorHandling(t *testing.T) {
	// Invalid email
	_, err := value_objects.NewEmail("not-an-email")
	assert.Error(t, err)

	// Weak password
	_, err = value_objects.NewPassword("weak")
	assert.Error(t, err)

	// Invalid locale
	_, err = value_objects.NewLocale("fr")
	assert.Error(t, err)

	// User with empty name
	email, _ := value_objects.NewEmail("user@example.com")
	password, _ := value_objects.NewPassword("SecurePass123")

	_, err = entities.NewUser(email, password, "", "Doe", value_objects.LocaleSpanish)
	assert.Error(t, err)
}
