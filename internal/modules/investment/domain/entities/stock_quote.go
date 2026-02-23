package entities

import (
	"strings"
	"time"

	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
)

// StockQuote represents real-time market data for a traded symbol.
// Symbol is the natural primary key (uppercase ticker).
// All monetary values are stored as int64 cents to avoid floating-point errors.
type StockQuote struct {
	Symbol   string // PK, uppercase
	Name     string
	Currency string

	// Price data — all in cents (smallest currency unit)
	PriceCents         int64
	PreviousCloseCents int64
	OpenCents          int64
	DayHighCents       int64
	DayLowCents        int64
	Volume             int64
	MarketCapCents     int64
	Week52HighCents    int64
	Week52LowCents     int64

	// Ratios — float64 is acceptable for dimensionless metrics
	TrailingPE    float64
	ForwardPE     float64
	DividendYield float64
	EPS           float64

	FetchedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewStockQuote creates and validates a new StockQuote entity.
// symbol must be non-empty and at most 10 characters.
// currency must be non-empty.
func NewStockQuote(symbol, name, currency string) (*StockQuote, error) {
	sym, err := NewSymbol(symbol)
	if err != nil {
		return nil, investErrors.NewStockQuoteValidationError("symbol", err.Error())
	}

	if strings.TrimSpace(currency) == "" {
		return nil, investErrors.NewStockQuoteValidationError("currency", "currency cannot be empty")
	}

	now := time.Now()
	return &StockQuote{
		Symbol:    sym.Value(),
		Name:      strings.TrimSpace(name),
		Currency:  strings.ToUpper(strings.TrimSpace(currency)),
		FetchedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update replaces all market data fields with fresh values from an external provider.
func (q *StockQuote) Update(
	name string,
	priceCents int64,
	previousCloseCents int64,
	openCents int64,
	dayHighCents int64,
	dayLowCents int64,
	volume int64,
	marketCapCents int64,
	week52HighCents int64,
	week52LowCents int64,
	trailingPE float64,
	forwardPE float64,
	dividendYield float64,
	eps float64,
	fetchedAt time.Time,
) {
	if strings.TrimSpace(name) != "" {
		q.Name = strings.TrimSpace(name)
	}
	q.PriceCents = priceCents
	q.PreviousCloseCents = previousCloseCents
	q.OpenCents = openCents
	q.DayHighCents = dayHighCents
	q.DayLowCents = dayLowCents
	q.Volume = volume
	q.MarketCapCents = marketCapCents
	q.Week52HighCents = week52HighCents
	q.Week52LowCents = week52LowCents
	q.TrailingPE = trailingPE
	q.ForwardPE = forwardPE
	q.DividendYield = dividendYield
	q.EPS = eps
	q.FetchedAt = fetchedAt
	q.UpdatedAt = time.Now()
}

// IsFresh returns true if the quote was fetched within the given TTL.
func (q *StockQuote) IsFresh(ttl time.Duration) bool {
	return time.Since(q.FetchedAt) < ttl
}

// NeedsRefresh returns true if the quote is older than 15 minutes.
func (q *StockQuote) NeedsRefresh() bool {
	return !q.IsFresh(15 * time.Minute)
}
