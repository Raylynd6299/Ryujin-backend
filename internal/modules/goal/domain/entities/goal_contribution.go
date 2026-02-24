package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	goalErrors "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/errors"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// GoalContribution represents a single monetary contribution toward a purchase goal.
// Contributions are immutable once created — delete and re-create to correct mistakes.
type GoalContribution struct {
	ID     string
	GoalID string
	UserID string

	Amount sharedVO.Money
	Date   time.Time
	Notes  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewGoalContribution creates a new contribution with validation
func NewGoalContribution(
	goalID string,
	userID string,
	amountCents int64,
	currency string,
	date time.Time,
	notes string,
) (*GoalContribution, error) {
	if strings.TrimSpace(goalID) == "" {
		return nil, goalErrors.NewContributionInvalidError("goalID cannot be empty")
	}

	if amountCents <= 0 {
		return nil, goalErrors.NewContributionInvalidError("contribution amount must be greater than zero")
	}

	money, err := sharedVO.NewMoney(amountCents, currency)
	if err != nil {
		return nil, goalErrors.NewContributionInvalidError("invalid currency: " + err.Error())
	}

	now := time.Now()
	return &GoalContribution{
		ID:        uuid.New().String(),
		GoalID:    goalID,
		UserID:    userID,
		Amount:    *money,
		Date:      date,
		Notes:     strings.TrimSpace(notes),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// BelongsTo checks if this contribution belongs to the given user
func (c *GoalContribution) BelongsTo(userID string) bool {
	return c.UserID == userID
}

// BelongsToGoal checks if this contribution belongs to the given goal
func (c *GoalContribution) BelongsToGoal(goalID string) bool {
	return c.GoalID == goalID
}
