package dto

import "time"

// --- Request DTOs ---

// CreateGoalRequest is used to create a new purchase goal
type CreateGoalRequest struct {
	Name         string  `json:"name" binding:"required,min=1,max=150"`
	Description  string  `json:"description" binding:"max=500"`
	Icon         string  `json:"icon" binding:"max=10"`
	TargetAmount float64 `json:"targetAmount" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required,len=3"`
	Priority     string  `json:"priority" binding:"required,oneof=low medium high"`
	Deadline     *string `json:"deadline"` // ISO 8601 date string YYYY-MM-DD, optional
}

// UpdateGoalRequest is used to update goal metadata
type UpdateGoalRequest struct {
	Name         string  `json:"name" binding:"required,min=1,max=150"`
	Description  string  `json:"description" binding:"max=500"`
	Icon         string  `json:"icon" binding:"max=10"`
	TargetAmount float64 `json:"targetAmount" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required,len=3"`
	Priority     string  `json:"priority" binding:"required,oneof=low medium high"`
	Deadline     *string `json:"deadline"` // ISO 8601 date string YYYY-MM-DD, optional
}

// CreateContributionRequest is used to add a contribution to a goal
type CreateContributionRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,len=3"`
	Date     string  `json:"date" binding:"required"` // ISO 8601 date string YYYY-MM-DD
	Notes    string  `json:"notes" binding:"max=300"`
}

// --- Response DTOs ---

// GoalResponse is returned for goal operations
type GoalResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Icon         string     `json:"icon"`
	TargetAmount float64    `json:"targetAmount"`
	Currency     string     `json:"currency"`
	Priority     string     `json:"priority"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	IsCompleted  bool       `json:"isCompleted"`
	// Computed analytics
	TotalContributed    float64    `json:"totalContributed"`
	ProgressPercent     float64    `json:"progressPercent"`
	MissingAmount       float64    `json:"missingAmount"`
	EstimatedCompletion *time.Time `json:"estimatedCompletion,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// GoalListResponse wraps a paginated list of goals
type GoalListResponse struct {
	Data       []*GoalResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	TotalPages int             `json:"totalPages"`
}

// ContributionResponse is returned for contribution operations
type ContributionResponse struct {
	ID        string    `json:"id"`
	GoalID    string    `json:"goalId"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ContributionListResponse wraps a list of contributions for a goal
type ContributionListResponse struct {
	Data  []*ContributionResponse `json:"data"`
	Total int                     `json:"total"`
}
