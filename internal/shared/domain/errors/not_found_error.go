package errors

import "fmt"

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	resource string
	id       string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		resource: resource,
		id:       id,
	}
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.resource, e.id)
}

// Resource returns the resource type
func (e *NotFoundError) Resource() string {
	return e.resource
}

// ID returns the resource ID
func (e *NotFoundError) ID() string {
	return e.id
}
