package errors

// UnauthorizedError represents an authentication/authorization error
type UnauthorizedError struct {
	message string
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		message: message,
	}
}

// Error implements the error interface
func (e *UnauthorizedError) Error() string {
	if e.message == "" {
		return "unauthorized access"
	}
	return e.message
}

// Message returns the error message
func (e *UnauthorizedError) Message() string {
	return e.message
}
