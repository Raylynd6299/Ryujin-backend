package errors

import "fmt"

// FinanceErrorCode represents error codes for the finance domain
type FinanceErrorCode string

const (
	ErrCodeCategoryNotFound     FinanceErrorCode = "CATEGORY_NOT_FOUND"
	ErrCodeCategoryInvalid      FinanceErrorCode = "CATEGORY_INVALID"
	ErrCodeIncomeSourceNotFound FinanceErrorCode = "INCOME_SOURCE_NOT_FOUND"
	ErrCodeIncomeSourceInvalid  FinanceErrorCode = "INCOME_SOURCE_INVALID"
	ErrCodeExpenseNotFound      FinanceErrorCode = "EXPENSE_NOT_FOUND"
	ErrCodeExpenseInvalid       FinanceErrorCode = "EXPENSE_INVALID"
	ErrCodeDebtNotFound         FinanceErrorCode = "DEBT_NOT_FOUND"
	ErrCodeDebtInvalid          FinanceErrorCode = "DEBT_INVALID"
	ErrCodeAccountNotFound      FinanceErrorCode = "ACCOUNT_NOT_FOUND"
	ErrCodeAccountInvalid       FinanceErrorCode = "ACCOUNT_INVALID"
	ErrCodeUnauthorized         FinanceErrorCode = "FINANCE_UNAUTHORIZED"
)

// FinanceError represents a domain error in the finance module
type FinanceError struct {
	Code    FinanceErrorCode
	Message string
}

func (e *FinanceError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewCategoryNotFoundError(id string) *FinanceError {
	return &FinanceError{Code: ErrCodeCategoryNotFound, Message: fmt.Sprintf("category %s not found", id)}
}

func NewCategoryInvalidError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeCategoryInvalid, Message: msg}
}

func NewIncomeSourceNotFoundError(id string) *FinanceError {
	return &FinanceError{Code: ErrCodeIncomeSourceNotFound, Message: fmt.Sprintf("income source %s not found", id)}
}

func NewIncomeSourceInvalidError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeIncomeSourceInvalid, Message: msg}
}

func NewExpenseNotFoundError(id string) *FinanceError {
	return &FinanceError{Code: ErrCodeExpenseNotFound, Message: fmt.Sprintf("expense %s not found", id)}
}

func NewExpenseInvalidError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeExpenseInvalid, Message: msg}
}

func NewDebtNotFoundError(id string) *FinanceError {
	return &FinanceError{Code: ErrCodeDebtNotFound, Message: fmt.Sprintf("debt %s not found", id)}
}

func NewDebtInvalidError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeDebtInvalid, Message: msg}
}

func NewAccountNotFoundError(id string) *FinanceError {
	return &FinanceError{Code: ErrCodeAccountNotFound, Message: fmt.Sprintf("account %s not found", id)}
}

func NewAccountInvalidError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeAccountInvalid, Message: msg}
}

func NewUnauthorizedError(msg string) *FinanceError {
	return &FinanceError{Code: ErrCodeUnauthorized, Message: msg}
}
