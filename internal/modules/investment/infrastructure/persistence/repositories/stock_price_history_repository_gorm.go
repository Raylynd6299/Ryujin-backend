package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/persistence/models"
)

// StockPriceHistoryRepositoryGorm implements domain/repositories.StockPriceHistoryRepository using GORM.
// Records are append-only — there is no update or delete operation.
type StockPriceHistoryRepositoryGorm struct {
	db *gorm.DB
}

// NewStockPriceHistoryRepositoryGorm creates a new GORM-backed price history repository.
func NewStockPriceHistoryRepositoryGorm(db *gorm.DB) *StockPriceHistoryRepositoryGorm {
	return &StockPriceHistoryRepositoryGorm{db: db}
}

// Create persists a new price history entry.
func (r *StockPriceHistoryRepositoryGorm) Create(ctx context.Context, entry *entities.StockPriceHistory) error {
	m := mappers.StockPriceHistoryToModel(entry)
	return r.db.WithContext(ctx).Create(m).Error
}

// FindBySymbol returns the most recent `limit` entries for the given symbol,
// ordered by RecordedAt descending.
func (r *StockPriceHistoryRepositoryGorm) FindBySymbol(ctx context.Context, symbol string, limit int) ([]*entities.StockPriceHistory, error) {
	var modelList []models.StockPriceHistoryModel
	err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("recorded_at DESC").
		Limit(limit).
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.StockPriceHistory, 0, len(modelList))
	for _, m := range modelList {
		m := m
		result = append(result, mappers.StockPriceHistoryFromModel(&m))
	}
	return result, nil
}
