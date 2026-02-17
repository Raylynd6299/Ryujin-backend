package errors

import "fmt"

// ValidationError represents a validation error with field-level details
type ValidationError struct {
	field   string
	message string
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		field:   field,
		message: message,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.field, e.message)
}

// Field returns the field that failed validation
func (e *ValidationError) Field() string {
	return e.field
}

// Message returns the validation error message
func (e *ValidationError) Message() string {
	return e.message
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	errors []*ValidationError
}

// NewValidationErrors creates a new collection of validation errors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		errors: make([]*ValidationError, 0),
	}
}

// Add adds a validation error to the collection
func (ve *ValidationErrors) Add(field, message string) {
	ve.errors = append(ve.errors, NewValidationError(field, message))
}

// HasErrors returns true if there are any validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.errors) > 0
}

// Errors returns all validation errors
func (ve *ValidationErrors) Errors() []*ValidationError {
	return ve.errors
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	if len(ve.errors) == 0 {
		return "no validation errors"
	}
	if len(ve.errors) == 1 {
		return ve.errors[0].Error()
	}
	return fmt.Sprintf("multiple validation errors: %d fields failed validation", len(ve.errors))
}
