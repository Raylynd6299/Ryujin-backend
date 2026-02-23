package errors

import (
	"errors"
	"fmt"
)

// ============================================================
// Sentinel errors
// ============================================================

var ErrStockQuoteNotFound = errors.New("stock quote not found")
var ErrInvalidSymbol = errors.New("invalid symbol")

// InvestmentErrorCode represents error codes for the investment domain
type InvestmentErrorCode string

const (
	ErrCodeHoldingNotFound   InvestmentErrorCode = "HOLDING_NOT_FOUND"
	ErrCodeHoldingForbidden  InvestmentErrorCode = "HOLDING_FORBIDDEN"
	ErrCodeHoldingValidation InvestmentErrorCode = "HOLDING_VALIDATION"
	ErrCodePriceRefresh      InvestmentErrorCode = "PRICE_REFRESH_ERROR"
)

// HoldingNotFoundError is returned when a holding cannot be found
type HoldingNotFoundError struct {
	ID string
}

func (e *HoldingNotFoundError) Error() string {
	return fmt.Sprintf("[%s] holding %s not found", ErrCodeHoldingNotFound, e.ID)
}

// NewHoldingNotFoundError creates a new HoldingNotFoundError
func NewHoldingNotFoundError(id string) *HoldingNotFoundError {
	return &HoldingNotFoundError{ID: id}
}

// HoldingForbiddenError is returned when a user attempts to access another user's holding
type HoldingForbiddenError struct {
	ID string
}

func (e *HoldingForbiddenError) Error() string {
	return fmt.Sprintf("[%s] access to holding %s is forbidden", ErrCodeHoldingForbidden, e.ID)
}

// NewHoldingForbiddenError creates a new HoldingForbiddenError
func NewHoldingForbiddenError(id string) *HoldingForbiddenError {
	return &HoldingForbiddenError{ID: id}
}

// HoldingValidationError is returned when a holding fails validation
type HoldingValidationError struct {
	Field   string
	Message string
}

func (e *HoldingValidationError) Error() string {
	return fmt.Sprintf("[%s] validation failed on field '%s': %s", ErrCodeHoldingValidation, e.Field, e.Message)
}

// NewHoldingValidationError creates a new HoldingValidationError
func NewHoldingValidationError(field, message string) *HoldingValidationError {
	return &HoldingValidationError{Field: field, Message: message}
}

// PriceRefreshError is returned when a price cannot be fetched from an external provider
type PriceRefreshError struct {
	Symbol string
	Reason string
}

func (e *PriceRefreshError) Error() string {
	return fmt.Sprintf("[%s] could not refresh price for symbol '%s': %s", ErrCodePriceRefresh, e.Symbol, e.Reason)
}

// NewPriceRefreshError creates a new PriceRefreshError
func NewPriceRefreshError(symbol, reason string) *PriceRefreshError {
	return &PriceRefreshError{Symbol: symbol, Reason: reason}
}

// ============================================================
// StockQuote errors
// ============================================================

const ErrCodeStockQuoteValidation InvestmentErrorCode = "STOCK_QUOTE_VALIDATION"

// StockQuoteValidationError is returned when a stock quote fails validation
type StockQuoteValidationError struct {
	Field   string
	Message string
}

func (e *StockQuoteValidationError) Error() string {
	return fmt.Sprintf("[%s] validation failed on field '%s': %s", ErrCodeStockQuoteValidation, e.Field, e.Message)
}

// NewStockQuoteValidationError creates a new StockQuoteValidationError
func NewStockQuoteValidationError(field, message string) *StockQuoteValidationError {
	return &StockQuoteValidationError{Field: field, Message: message}
}
