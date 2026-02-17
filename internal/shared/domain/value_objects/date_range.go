package value_objects

import (
	"errors"
	"time"
)

// DateRange represents a time period with start and end dates
type DateRange struct {
	start time.Time
	end   time.Time
}

// NewDateRange creates a new DateRange value object
func NewDateRange(start, end time.Time) (*DateRange, error) {
	if start.After(end) {
		return nil, errors.New("start date cannot be after end date")
	}

	return &DateRange{
		start: start,
		end:   end,
	}, nil
}

// Start returns the start date
func (d *DateRange) Start() time.Time {
	return d.start
}

// End returns the end date
func (d *DateRange) End() time.Time {
	return d.end
}

// Duration returns the duration between start and end
func (d *DateRange) Duration() time.Duration {
	return d.end.Sub(d.start)
}

// Contains checks if a given date is within the range
func (d *DateRange) Contains(date time.Time) bool {
	return (date.Equal(d.start) || date.After(d.start)) &&
		(date.Equal(d.end) || date.Before(d.end))
}

// Overlaps checks if this date range overlaps with another
func (d *DateRange) Overlaps(other *DateRange) bool {
	return d.start.Before(other.end) && d.end.After(other.start)
}

// String returns a string representation of the date range
func (d *DateRange) String() string {
	return d.start.Format("2006-01-02") + " to " + d.end.Format("2006-01-02")
}
