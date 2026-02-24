package errors

import "fmt"

// GoalErrorCode represents error codes for the goal domain
type GoalErrorCode string

const (
	ErrCodeGoalNotFound         GoalErrorCode = "GOAL_NOT_FOUND"
	ErrCodeGoalInvalid          GoalErrorCode = "GOAL_INVALID"
	ErrCodeContributionNotFound GoalErrorCode = "CONTRIBUTION_NOT_FOUND"
	ErrCodeContributionInvalid  GoalErrorCode = "CONTRIBUTION_INVALID"
	ErrCodeGoalUnauthorized     GoalErrorCode = "GOAL_UNAUTHORIZED"
)

// GoalError represents a domain error in the goal module
type GoalError struct {
	Code    GoalErrorCode
	Message string
}

func (e *GoalError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewGoalNotFoundError(id string) *GoalError {
	return &GoalError{Code: ErrCodeGoalNotFound, Message: fmt.Sprintf("goal %s not found", id)}
}

func NewGoalInvalidError(msg string) *GoalError {
	return &GoalError{Code: ErrCodeGoalInvalid, Message: msg}
}

func NewContributionNotFoundError(id string) *GoalError {
	return &GoalError{Code: ErrCodeContributionNotFound, Message: fmt.Sprintf("contribution %s not found", id)}
}

func NewContributionInvalidError(msg string) *GoalError {
	return &GoalError{Code: ErrCodeContributionInvalid, Message: msg}
}

func NewGoalUnauthorizedError(msg string) *GoalError {
	return &GoalError{Code: ErrCodeGoalUnauthorized, Message: msg}
}
