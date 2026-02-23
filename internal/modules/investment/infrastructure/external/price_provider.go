package external

import (
	"context"
	"time"
)

// PriceQuote holds the result of a price fetch from an external provider
type PriceQuote struct {
	Symbol     string
	PriceCents int64
	Currency   string
	FetchedAt  time.Time
}

// StockQuote holds enriched quote data for stock analysis
type StockQuote struct {
	Symbol             string
	Name               string
	Currency           string
	RegularMarketPrice float64
	PreviousClose      float64
	Open               float64
	DayHigh            float64
	DayLow             float64
	Volume             int64
	MarketCap          int64
	FiftyTwoWeekHigh   float64
	FiftyTwoWeekLow    float64
	TrailingPE         float64 // P/E ratio (trailing twelve months)
	ForwardPE          float64
	DividendYield      float64
	EPS                float64
	FetchedAt          time.Time
}

// PriceProvider is the port for fetching real-time asset prices
type PriceProvider interface {
	// FetchPrice fetches the current price for the given symbol
	FetchPrice(ctx context.Context, symbol string) (*PriceQuote, error)

	// Name returns the name of the price provider
	Name() string
}

// QuoteProvider is the port for fetching enriched stock quote data
type QuoteProvider interface {
	// FetchQuote fetches enriched market data for the given symbol
	FetchQuote(ctx context.Context, symbol string) (*StockQuote, error)
}
