package errors

import "fmt"

// DomainError represents a generic domain-level error
type DomainError struct {
	message string
	code    string
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		code:    code,
		message: message,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

// Code returns the error code
func (e *DomainError) Code() string {
	return e.code
}

// Message returns the error message
func (e *DomainError) Message() string {
	return e.message
}
