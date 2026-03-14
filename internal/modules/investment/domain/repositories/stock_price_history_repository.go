package repositories

import (
	"context"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/entities"
)

// StockPriceHistoryRepository defines the port for stock price history persistence.
// Records are append-only — there is no update or delete method.
type StockPriceHistoryRepository interface {
	// Create persists a new price history entry.
	Create(ctx context.Context, entry *entities.StockPriceHistory) error

	// FindBySymbol returns the most recent `limit` entries for the given symbol,
	// ordered by RecordedAt descending.
	FindBySymbol(ctx context.Context, symbol string, limit int) ([]*entities.StockPriceHistory, error)
}
