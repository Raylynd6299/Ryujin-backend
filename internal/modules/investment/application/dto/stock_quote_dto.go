package dto

import (
	"time"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
)

// StockQuoteResponse is the API response for stock quote data
type StockQuoteResponse struct {
	Symbol             string    `json:"symbol"`
	Name               string    `json:"name"`
	Currency           string    `json:"currency"`
	RegularMarketPrice float64   `json:"regularMarketPrice"` // priceCents / 100
	PreviousClose      float64   `json:"previousClose"`
	Open               float64   `json:"open"`
	DayHigh            float64   `json:"dayHigh"`
	DayLow             float64   `json:"dayLow"`
	ChangeAmount       float64   `json:"changeAmount"`
	ChangePct          float64   `json:"changePct"`
	Volume             int64     `json:"volume"`
	MarketCap          int64     `json:"marketCap"` // in cents
	FiftyTwoWeekHigh   float64   `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow    float64   `json:"fiftyTwoWeekLow"`
	TrailingPE         float64   `json:"trailingPE"`
	ForwardPE          float64   `json:"forwardPE"`
	DividendYield      float64   `json:"dividendYield"`
	EPS                float64   `json:"eps"`
	IsFresh            bool      `json:"isFresh"`
	FetchedAt          time.Time `json:"fetchedAt"`
}

// StockPriceHistoryResponse is the API response for a single history entry
type StockPriceHistoryResponse struct {
	Symbol     string    `json:"symbol"`
	Price      float64   `json:"price"` // priceCents / 100
	Currency   string    `json:"currency"`
	RecordedAt time.Time `json:"recordedAt"`
}

// ToStockQuoteResponse converts a domain StockQuote entity to a StockQuoteResponse DTO.
// All cent-stored values are divided by 100 for the API response.
func ToStockQuoteResponse(q *entities.StockQuote) StockQuoteResponse {
	changeAmount := float64(q.PriceCents-q.PreviousCloseCents) / 100.0
	var changePct float64
	if q.PreviousCloseCents > 0 {
		changePct = changeAmount / (float64(q.PreviousCloseCents) / 100.0) * 100.0
	}

	return StockQuoteResponse{
		Symbol:             q.Symbol,
		Name:               q.Name,
		Currency:           q.Currency,
		RegularMarketPrice: float64(q.PriceCents) / 100.0,
		PreviousClose:      float64(q.PreviousCloseCents) / 100.0,
		Open:               float64(q.OpenCents) / 100.0,
		DayHigh:            float64(q.DayHighCents) / 100.0,
		DayLow:             float64(q.DayLowCents) / 100.0,
		ChangeAmount:       changeAmount,
		ChangePct:          changePct,
		Volume:             q.Volume,
		MarketCap:          q.MarketCapCents,
		FiftyTwoWeekHigh:   float64(q.Week52HighCents) / 100.0,
		FiftyTwoWeekLow:    float64(q.Week52LowCents) / 100.0,
		TrailingPE:         q.TrailingPE,
		ForwardPE:          q.ForwardPE,
		DividendYield:      q.DividendYield,
		EPS:                q.EPS,
		IsFresh:            !q.NeedsRefresh(),
		FetchedAt:          q.FetchedAt,
	}
}

// ToStockPriceHistoryResponse converts a domain StockPriceHistory entity to a DTO.
func ToStockPriceHistoryResponse(h *entities.StockPriceHistory) StockPriceHistoryResponse {
	return StockPriceHistoryResponse{
		Symbol:     h.Symbol,
		Price:      float64(h.PriceCents) / 100.0,
		Currency:   h.Currency,
		RecordedAt: h.RecordedAt,
	}
}
