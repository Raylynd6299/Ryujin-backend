package utils

import (
	"errors"
	"fmt"
	"strings"
)

// ValidateRequiredString ensures a string is not empty or whitespace.
func ValidateRequiredString(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

// ValidatePositiveInt ensures an integer is greater than zero.
func ValidatePositiveInt(value int, field string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than zero", field)
	}
	return nil
}

// ValidateNonNil ensures a value is not nil.
func ValidateNonNil(value interface{}, field string) error {
	if value == nil {
		return errors.New(field + " cannot be nil")
	}
	return nil
}
