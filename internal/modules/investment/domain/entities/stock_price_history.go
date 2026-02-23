package entities

import (
	"strings"
	"time"

	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	"github.com/google/uuid"
)

// StockPriceHistory is an append-only record of a point-in-time price for a symbol.
// Records are never updated — only created.
type StockPriceHistory struct {
	ID         string
	Symbol     string // FK to stock_quotes.symbol
	PriceCents int64
	Currency   string
	RecordedAt time.Time
}

// NewStockPriceHistory creates and validates a new StockPriceHistory entry.
// symbol must be non-empty and at most 10 characters.
// priceCents must be > 0.
// currency must be non-empty.
func NewStockPriceHistory(symbol string, priceCents int64, currency string) (*StockPriceHistory, error) {
	sym, err := NewSymbol(symbol)
	if err != nil {
		return nil, investErrors.NewStockQuoteValidationError("symbol", err.Error())
	}

	if priceCents <= 0 {
		return nil, investErrors.NewStockQuoteValidationError("price_cents", "price must be greater than zero")
	}

	if strings.TrimSpace(currency) == "" {
		return nil, investErrors.NewStockQuoteValidationError("currency", "currency cannot be empty")
	}

	return &StockPriceHistory{
		ID:         uuid.New().String(),
		Symbol:     sym.Value(),
		PriceCents: priceCents,
		Currency:   strings.ToUpper(strings.TrimSpace(currency)),
		RecordedAt: time.Now(),
	}, nil
}
