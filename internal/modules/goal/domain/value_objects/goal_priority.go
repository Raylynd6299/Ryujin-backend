package value_objects

import "errors"

// GoalPriorityLevel classifies the importance of a purchase goal
type GoalPriorityLevel string

const (
	GoalPriorityLow    GoalPriorityLevel = "low"
	GoalPriorityMedium GoalPriorityLevel = "medium"
	GoalPriorityHigh   GoalPriorityLevel = "high"
)

var validGoalPriorities = map[GoalPriorityLevel]bool{
	GoalPriorityLow:    true,
	GoalPriorityMedium: true,
	GoalPriorityHigh:   true,
}

// GoalPriority is a value object for goal importance classification
type GoalPriority struct {
	value GoalPriorityLevel
}

// NewGoalPriority creates a validated GoalPriority
func NewGoalPriority(p string) (*GoalPriority, error) {
	gp := GoalPriorityLevel(p)
	if !validGoalPriorities[gp] {
		return nil, errors.New("invalid goal priority: " + p + " (must be low, medium, or high)")
	}
	return &GoalPriority{value: gp}, nil
}

func (g *GoalPriority) Value() GoalPriorityLevel {
	return g.value
}

func (g *GoalPriority) String() string {
	return string(g.value)
}

// IsHigh returns true for high-priority goals
func (g *GoalPriority) IsHigh() bool {
	return g.value == GoalPriorityHigh
}
