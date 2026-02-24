package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	goalErrors "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/errors"
	vo "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/value_objects"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// PurchaseGoal represents a savings target for a specific purchase.
// Examples: "New laptop", "Vacation to Japan", "Emergency fund".
type PurchaseGoal struct {
	ID     string
	UserID string

	Name        string
	Description string
	Icon        string // emoji or icon identifier (e.g. "💻", "✈️")

	// Financial data
	TargetAmount sharedVO.Money

	// Optional deadline
	Deadline *time.Time

	// Priority classification
	Priority vo.GoalPriority

	// Status
	IsCompleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPurchaseGoal creates a new purchase goal with validation
func NewPurchaseGoal(
	userID string,
	name string,
	description string,
	icon string,
	targetAmountCents int64,
	currency string,
	priority string,
	deadline *time.Time,
) (*PurchaseGoal, error) {
	if strings.TrimSpace(name) == "" {
		return nil, goalErrors.NewGoalInvalidError("goal name cannot be empty")
	}

	if targetAmountCents <= 0 {
		return nil, goalErrors.NewGoalInvalidError("target amount must be greater than zero")
	}

	targetMoney, err := sharedVO.NewMoney(targetAmountCents, currency)
	if err != nil {
		return nil, goalErrors.NewGoalInvalidError("invalid currency: " + err.Error())
	}

	p, err := vo.NewGoalPriority(priority)
	if err != nil {
		return nil, goalErrors.NewGoalInvalidError(err.Error())
	}

	now := time.Now()
	return &PurchaseGoal{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         strings.TrimSpace(name),
		Description:  strings.TrimSpace(description),
		Icon:         strings.TrimSpace(icon),
		TargetAmount: *targetMoney,
		Deadline:     deadline,
		Priority:     *p,
		IsCompleted:  false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates the goal metadata
func (g *PurchaseGoal) Update(
	name string,
	description string,
	icon string,
	targetAmountCents int64,
	currency string,
	priority string,
	deadline *time.Time,
) error {
	if strings.TrimSpace(name) == "" {
		return goalErrors.NewGoalInvalidError("goal name cannot be empty")
	}

	if targetAmountCents <= 0 {
		return goalErrors.NewGoalInvalidError("target amount must be greater than zero")
	}

	targetMoney, err := sharedVO.NewMoney(targetAmountCents, currency)
	if err != nil {
		return goalErrors.NewGoalInvalidError("invalid currency: " + err.Error())
	}

	p, err := vo.NewGoalPriority(priority)
	if err != nil {
		return goalErrors.NewGoalInvalidError(err.Error())
	}

	g.Name = strings.TrimSpace(name)
	g.Description = strings.TrimSpace(description)
	g.Icon = strings.TrimSpace(icon)
	g.TargetAmount = *targetMoney
	g.Priority = *p
	g.Deadline = deadline
	g.UpdatedAt = time.Now()

	return nil
}

// MarkCompleted marks the goal as completed
func (g *PurchaseGoal) MarkCompleted() {
	g.IsCompleted = true
	g.UpdatedAt = time.Now()
}

// BelongsTo checks if this goal belongs to the given user
func (g *PurchaseGoal) BelongsTo(userID string) bool {
	return g.UserID == userID
}

// ProgressPercent returns the percentage of the target amount contributed (0-100+)
// Can exceed 100 if total contributions surpass the target
func (g *PurchaseGoal) ProgressPercent(contributions []*GoalContribution) float64 {
	target := g.TargetAmount.Amount()
	if target == 0 {
		return 100.0
	}
	total := totalContributed(contributions)
	return float64(total) / float64(target) * 100.0
}

// MissingAmount returns how much is still needed to reach the target.
// Returns 0 if already over-funded.
func (g *PurchaseGoal) MissingAmount(contributions []*GoalContribution) sharedVO.Money {
	target := g.TargetAmount.Amount()
	total := totalContributed(contributions)
	missing := target - total
	if missing < 0 {
		missing = 0
	}
	m, _ := sharedVO.NewMoney(missing, g.TargetAmount.Currency())
	return *m
}

// TotalContributed returns the sum of all contributions in cents
func (g *PurchaseGoal) TotalContributed(contributions []*GoalContribution) sharedVO.Money {
	total := totalContributed(contributions)
	m, _ := sharedVO.NewMoney(total, g.TargetAmount.Currency())
	return *m
}

// EstimatedCompletionDate calculates when the goal will be reached
// based on the rolling average of the last 3 months of contributions.
// Returns nil if there are no contributions or average is zero.
func (g *PurchaseGoal) EstimatedCompletionDate(contributions []*GoalContribution) *time.Time {
	if len(contributions) == 0 {
		return nil
	}

	avgMonthly := rollingAverageMonthly(contributions, 3)
	if avgMonthly <= 0 {
		return nil
	}

	target := g.TargetAmount.Amount()
	totalSoFar := totalContributed(contributions)
	remaining := target - totalSoFar
	if remaining <= 0 {
		now := time.Now()
		return &now
	}

	// months needed (ceiling)
	monthsNeeded := int((remaining + avgMonthly - 1) / avgMonthly)
	eta := time.Now().AddDate(0, monthsNeeded, 0)
	return &eta
}

// IsOverFunded returns true if total contributions exceed the target
func (g *PurchaseGoal) IsOverFunded(contributions []*GoalContribution) bool {
	return totalContributed(contributions) >= g.TargetAmount.Amount()
}

// --- internal helpers ---

// totalContributed sums up all contribution amounts in cents
func totalContributed(contributions []*GoalContribution) int64 {
	var total int64
	for _, c := range contributions {
		total += c.Amount.Amount()
	}
	return total
}

// rollingAverageMonthly computes the average monthly contribution
// over the last `months` calendar months.
// Falls back to overall average if there is less history.
func rollingAverageMonthly(contributions []*GoalContribution, months int) int64 {
	if len(contributions) == 0 {
		return 0
	}

	cutoff := time.Now().AddDate(0, -months, 0)
	var recentTotal int64
	var recentCount int

	for _, c := range contributions {
		if c.Date.After(cutoff) {
			recentTotal += c.Amount.Amount()
			recentCount++
		}
	}

	// If we have recent data, divide by the number of months window
	if recentCount > 0 {
		return recentTotal / int64(months)
	}

	// Fallback: overall average across all time
	var allTotal int64
	for _, c := range contributions {
		allTotal += c.Amount.Amount()
	}
	return allTotal / int64(len(contributions))
}
