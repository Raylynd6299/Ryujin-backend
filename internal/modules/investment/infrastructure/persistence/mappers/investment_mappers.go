package mappers

import (
	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/persistence/models"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// ============================================================
// StockQuote mappers
// ============================================================

// StockQuoteToModel converts a domain StockQuote entity to its GORM model.
func StockQuoteToModel(q *entities.StockQuote) *models.StockQuoteModel {
	return &models.StockQuoteModel{
		Symbol:             q.Symbol,
		Name:               q.Name,
		Currency:           q.Currency,
		PriceCents:         q.PriceCents,
		PreviousCloseCents: q.PreviousCloseCents,
		OpenCents:          q.OpenCents,
		DayHighCents:       q.DayHighCents,
		DayLowCents:        q.DayLowCents,
		Volume:             q.Volume,
		MarketCapCents:     q.MarketCapCents,
		Week52HighCents:    q.Week52HighCents,
		Week52LowCents:     q.Week52LowCents,
		TrailingPE:         q.TrailingPE,
		ForwardPE:          q.ForwardPE,
		DividendYield:      q.DividendYield,
		EPS:                q.EPS,
		FetchedAt:          q.FetchedAt,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
	}
}

// StockQuoteFromModel converts a GORM StockQuoteModel to its domain entity.
func StockQuoteFromModel(m *models.StockQuoteModel) *entities.StockQuote {
	return &entities.StockQuote{
		Symbol:             m.Symbol,
		Name:               m.Name,
		Currency:           m.Currency,
		PriceCents:         m.PriceCents,
		PreviousCloseCents: m.PreviousCloseCents,
		OpenCents:          m.OpenCents,
		DayHighCents:       m.DayHighCents,
		DayLowCents:        m.DayLowCents,
		Volume:             m.Volume,
		MarketCapCents:     m.MarketCapCents,
		Week52HighCents:    m.Week52HighCents,
		Week52LowCents:     m.Week52LowCents,
		TrailingPE:         m.TrailingPE,
		ForwardPE:          m.ForwardPE,
		DividendYield:      m.DividendYield,
		EPS:                m.EPS,
		FetchedAt:          m.FetchedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// ============================================================
// StockPriceHistory mappers
// ============================================================

// StockPriceHistoryToModel converts a domain StockPriceHistory to its GORM model.
func StockPriceHistoryToModel(h *entities.StockPriceHistory) *models.StockPriceHistoryModel {
	return &models.StockPriceHistoryModel{
		ID:         h.ID,
		Symbol:     h.Symbol,
		PriceCents: h.PriceCents,
		Currency:   h.Currency,
		RecordedAt: h.RecordedAt,
	}
}

// StockPriceHistoryFromModel converts a GORM StockPriceHistoryModel to its domain entity.
func StockPriceHistoryFromModel(m *models.StockPriceHistoryModel) *entities.StockPriceHistory {
	return &entities.StockPriceHistory{
		ID:         m.ID,
		Symbol:     m.Symbol,
		PriceCents: m.PriceCents,
		Currency:   m.Currency,
		RecordedAt: m.RecordedAt,
	}
}

// HoldingToModel converts a domain Holding entity to its GORM model
func HoldingToModel(h *entities.Holding) *models.HoldingModel {
	m := &models.HoldingModel{
		ID:            h.ID,
		UserID:        h.UserID,
		Symbol:        h.Symbol.Value(),
		Name:          h.Name,
		AssetType:     h.AssetType.String(),
		QuantityMicro: h.Quantity.MicroUnits(),
		BuyPriceCents: h.BuyPrice.Amount(),
		BuyCurrency:   h.Currency.Code(),
		Notes:         h.Notes,
		CreatedAt:     h.CreatedAt,
		UpdatedAt:     h.UpdatedAt,
	}

	if h.CurrentPrice != nil {
		cents := h.CurrentPrice.Amount()
		m.CurrentPriceCents = &cents
	}

	if h.PricedAt != nil {
		t := *h.PricedAt
		m.PricedAt = &t
	}

	return m
}

// ModelToHolding converts a GORM HoldingModel to its domain entity.
// Returns an error if any stored value is corrupt or cannot be mapped.
func ModelToHolding(m *models.HoldingModel) (*entities.Holding, error) {
	sym, err := entities.NewSymbol(m.Symbol)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("symbol", "invalid symbol in DB: "+err.Error())
	}

	at, err := entities.NewAssetType(m.AssetType)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("asset_type", "invalid asset type in DB: "+err.Error())
	}

	qty, err := entities.NewQuantity(m.QuantityMicro)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("quantity", "invalid quantity in DB: "+err.Error())
	}

	cur, err := sharedVO.NewCurrency(m.BuyCurrency)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("currency", "invalid currency in DB: "+err.Error())
	}

	buyPrice, err := sharedVO.NewMoney(m.BuyPriceCents, cur.Code())
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("buy_price", "invalid buy price in DB: "+err.Error())
	}

	h := &entities.Holding{
		ID:        m.ID,
		UserID:    m.UserID,
		Symbol:    sym,
		Name:      m.Name,
		AssetType: at,
		Quantity:  qty,
		BuyPrice:  *buyPrice,
		Currency:  *cur,
		Notes:     m.Notes,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	// Handle nullable current price
	if m.CurrentPriceCents != nil {
		currentPrice, err := sharedVO.NewMoney(*m.CurrentPriceCents, cur.Code())
		if err != nil {
			return nil, investErrors.NewHoldingValidationError("current_price", "invalid current price in DB: "+err.Error())
		}
		h.CurrentPrice = currentPrice
	}

	// Handle nullable priced_at
	if m.PricedAt != nil {
		t := *m.PricedAt
		h.PricedAt = &t
	}

	return h, nil
}
