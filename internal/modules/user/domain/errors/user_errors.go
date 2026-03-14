package errors

import (
	sharedErrors "github.com/Raylynd6299/Ryujin-backend/internal/shared/domain/errors"
)

// InvalidUserError is a validation error for user validation failures.
type InvalidUserError = sharedErrors.ValidationError

// NewInvalidUserError creates a new InvalidUserError.
func NewInvalidUserError(message string) *InvalidUserError {
	return sharedErrors.NewValidationError("INVALID_USER", message)
}

// UserNotFoundError is a not found error for missing users.
type UserNotFoundError = sharedErrors.NotFoundError

// NewUserNotFoundError creates a new UserNotFoundError.
func NewUserNotFoundError(message string) *UserNotFoundError {
	return sharedErrors.NewNotFoundError("USER_NOT_FOUND", message)
}

// InvalidPasswordError is a validation error for password validation failures.
type InvalidPasswordError = sharedErrors.ValidationError

// NewInvalidPasswordError creates a new InvalidPasswordError.
func NewInvalidPasswordError(message string) *InvalidPasswordError {
	return sharedErrors.NewValidationError("INVALID_PASSWORD", message)
}

// DuplicateEmailError is a domain error for email uniqueness violation.
type DuplicateEmailError = sharedErrors.DomainError

// NewDuplicateEmailError creates a new DuplicateEmailError.
func NewDuplicateEmailError(email string) *DuplicateEmailError {
	return sharedErrors.NewDomainError("DUPLICATE_EMAIL", "email already in use: "+email)
}
