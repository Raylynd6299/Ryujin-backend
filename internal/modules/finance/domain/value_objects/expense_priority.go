package value_objects

import "errors"

// ExpensePriority defines how important an expense is
type ExpensePriority string

const (
	PriorityEssential ExpensePriority = "essential" // Rent, food, utilities — can't skip
	PriorityImportant ExpensePriority = "important" // Insurance, subscriptions — hard to cut
	PriorityOptional  ExpensePriority = "optional"  // Dining out, entertainment — can reduce
	PriorityLow       ExpensePriority = "low"       // Impulse buys, luxuries — easy to cut
)

var validPriorities = map[ExpensePriority]bool{
	PriorityEssential: true,
	PriorityImportant: true,
	PriorityOptional:  true,
	PriorityLow:       true,
}

// Priority is a value object representing expense priority
type Priority struct {
	value ExpensePriority
}

// NewPriority creates a validated Priority value object
func NewPriority(p string) (*Priority, error) {
	ep := ExpensePriority(p)
	if !validPriorities[ep] {
		return nil, errors.New("invalid expense priority: " + p)
	}
	return &Priority{value: ep}, nil
}

func (p *Priority) Value() ExpensePriority {
	return p.value
}

func (p *Priority) String() string {
	return string(p.value)
}

// IsUnnecessary returns true for expenses that can easily be reduced
func (p *Priority) IsUnnecessary() bool {
	return p.value == PriorityOptional || p.value == PriorityLow
}
