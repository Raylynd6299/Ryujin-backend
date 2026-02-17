package utils

// Pagination defines paging parameters.
type Pagination struct {
	Page    int
	PerPage int
}

// Default pagination values.
const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// NormalizePagination applies defaults and bounds.
func NormalizePagination(p Pagination) Pagination {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PerPage <= 0 {
		p.PerPage = DefaultPerPage
	}
	if p.PerPage > MaxPerPage {
		p.PerPage = MaxPerPage
	}
	return p
}

// Offset returns the offset for database queries.
func (p Pagination) Offset() int {
	if p.Page < 1 {
		return 0
	}
	if p.PerPage <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

// Limit returns the per-page limit.
func (p Pagination) Limit() int {
	if p.PerPage <= 0 {
		return DefaultPerPage
	}
	return p.PerPage
}
