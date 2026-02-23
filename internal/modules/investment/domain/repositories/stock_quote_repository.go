package repositories

import (
	"context"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
)

// StockQuoteRepository defines the port for stock quote persistence.
// Implementations live in the infrastructure layer.
type StockQuoteRepository interface {
	// FindBySymbol returns the quote for the given symbol.
	// Returns ErrStockQuoteNotFound if the symbol is not cached.
	FindBySymbol(ctx context.Context, symbol string) (*entities.StockQuote, error)

	// Upsert inserts or updates a quote record keyed on symbol.
	Upsert(ctx context.Context, quote *entities.StockQuote) error

	// FindAll returns every cached quote.
	FindAll(ctx context.Context) ([]*entities.StockQuote, error)

	// FindSymbolsWithActiveHoldings returns the distinct symbols that
	// appear in at least one active user holding.
	FindSymbolsWithActiveHoldings(ctx context.Context) ([]string, error)
}
