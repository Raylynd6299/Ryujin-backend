package value_objects

import "errors"

// RecurrenceType defines how often an income/expense repeats
type RecurrenceType string

const (
	RecurrenceNone      RecurrenceType = "none"
	RecurrenceDaily     RecurrenceType = "daily"
	RecurrenceWeekly    RecurrenceType = "weekly"
	RecurrenceBiweekly  RecurrenceType = "biweekly"
	RecurrenceMonthly   RecurrenceType = "monthly"
	RecurrenceQuarterly RecurrenceType = "quarterly"
	RecurrenceAnnually  RecurrenceType = "annually"
)

var validRecurrences = map[RecurrenceType]bool{
	RecurrenceNone:      true,
	RecurrenceDaily:     true,
	RecurrenceWeekly:    true,
	RecurrenceBiweekly:  true,
	RecurrenceMonthly:   true,
	RecurrenceQuarterly: true,
	RecurrenceAnnually:  true,
}

// Recurrence is a value object that represents how often something repeats
type Recurrence struct {
	recurrenceType RecurrenceType
}

// NewRecurrence creates a validated Recurrence value object
func NewRecurrence(t string) (*Recurrence, error) {
	rt := RecurrenceType(t)
	if !validRecurrences[rt] {
		return nil, errors.New("invalid recurrence type: " + t)
	}
	return &Recurrence{recurrenceType: rt}, nil
}

func (r *Recurrence) Type() RecurrenceType {
	return r.recurrenceType
}

func (r *Recurrence) String() string {
	return string(r.recurrenceType)
}

func (r *Recurrence) IsRecurring() bool {
	return r.recurrenceType != RecurrenceNone
}
