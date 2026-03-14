package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/entities"
	investErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/errors"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/persistence/models"
)

// StockQuoteRepositoryGorm implements domain/repositories.StockQuoteRepository using GORM.
type StockQuoteRepositoryGorm struct {
	db *gorm.DB
}

// NewStockQuoteRepositoryGorm creates a new GORM-backed stock quote repository.
func NewStockQuoteRepositoryGorm(db *gorm.DB) *StockQuoteRepositoryGorm {
	return &StockQuoteRepositoryGorm{db: db}
}

// FindBySymbol returns the cached quote for the given symbol.
// Returns ErrStockQuoteNotFound if the symbol is not in the database.
func (r *StockQuoteRepositoryGorm) FindBySymbol(ctx context.Context, symbol string) (*entities.StockQuote, error) {
	var m models.StockQuoteModel
	err := r.db.WithContext(ctx).First(&m, "symbol = ?", symbol).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, investErrors.ErrStockQuoteNotFound
		}
		return nil, err
	}
	return mappers.StockQuoteFromModel(&m), nil
}

// Upsert inserts or updates a quote record keyed on symbol.
// GORM Save performs an upsert when the primary key (Symbol) is present.
func (r *StockQuoteRepositoryGorm) Upsert(ctx context.Context, quote *entities.StockQuote) error {
	m := mappers.StockQuoteToModel(quote)
	return r.db.WithContext(ctx).Save(m).Error
}

// FindAll returns every cached quote in the database.
func (r *StockQuoteRepositoryGorm) FindAll(ctx context.Context) ([]*entities.StockQuote, error) {
	var modelList []models.StockQuoteModel
	if err := r.db.WithContext(ctx).Find(&modelList).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.StockQuote, 0, len(modelList))
	for _, m := range modelList {
		m := m
		result = append(result, mappers.StockQuoteFromModel(&m))
	}
	return result, nil
}

// FindSymbolsWithActiveHoldings returns the distinct symbols that appear in at
// least one holding row, joined against the stock_quotes table.
func (r *StockQuoteRepositoryGorm) FindSymbolsWithActiveHoldings(ctx context.Context) ([]string, error) {
	var symbols []string
	err := r.db.WithContext(ctx).
		Model(&models.StockQuoteModel{}).
		Select("DISTINCT stock_quotes.symbol").
		Joins("INNER JOIN holdings h ON h.symbol = stock_quotes.symbol").
		Pluck("stock_quotes.symbol", &symbols).Error
	if err != nil {
		return nil, err
	}
	return symbols, nil
}
