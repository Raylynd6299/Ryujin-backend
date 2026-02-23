package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	domainRepos "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/repositories"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/external"
)

// StockAnalysisService handles find-or-fetch logic for stock quotes.
// It reads from the local cache first; only calls the external provider when
// the quote is missing or stale (NeedsRefresh).
type StockAnalysisService struct {
	stockQuoteRepo   domainRepos.StockQuoteRepository
	stockHistoryRepo domainRepos.StockPriceHistoryRepository
	quoteProvider    external.QuoteProvider
}

// NewStockAnalysisService creates a new StockAnalysisService.
func NewStockAnalysisService(
	stockQuoteRepo domainRepos.StockQuoteRepository,
	stockHistoryRepo domainRepos.StockPriceHistoryRepository,
	quoteProvider external.QuoteProvider,
) *StockAnalysisService {
	return &StockAnalysisService{
		stockQuoteRepo:   stockQuoteRepo,
		stockHistoryRepo: stockHistoryRepo,
		quoteProvider:    quoteProvider,
	}
}

// GetStockQuote returns enriched quote data for the given symbol.
// Fast path: if the cached quote is still fresh it is returned immediately.
// Slow path: the external provider is called, the cache is updated, and a
// price-history entry is appended.
func (s *StockAnalysisService) GetStockQuote(ctx context.Context, symbol string) (*dto.StockQuoteResponse, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	existing, err := s.stockQuoteRepo.FindBySymbol(ctx, symbol)

	switch {
	case errors.Is(err, investErrors.ErrStockQuoteNotFound):
		// Cache miss — fetch from provider, persist, record history
		quote, fetchErr := s.fetchAndPersist(ctx, symbol)
		if fetchErr != nil {
			return nil, fetchErr
		}
		resp := dto.ToStockQuoteResponse(quote)
		return &resp, nil

	case err != nil:
		// Unexpected repository error
		return nil, fmt.Errorf("failed to query stock quote for %s: %w", symbol, err)

	case existing.NeedsRefresh():
		// Cache hit but stale — refresh and update
		quote, fetchErr := s.refreshAndPersist(ctx, existing)
		if fetchErr != nil {
			// Return stale data rather than failing completely
			resp := dto.ToStockQuoteResponse(existing)
			return &resp, nil
		}
		resp := dto.ToStockQuoteResponse(quote)
		return &resp, nil

	default:
		// Cache hit and fresh — fast path
		resp := dto.ToStockQuoteResponse(existing)
		return &resp, nil
	}
}

// ListStockQuotes returns all cached stock quotes.
func (s *StockAnalysisService) ListStockQuotes(ctx context.Context) ([]*dto.StockQuoteResponse, error) {
	quotes, err := s.stockQuoteRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock quotes: %w", err)
	}

	responses := make([]*dto.StockQuoteResponse, 0, len(quotes))
	for _, q := range quotes {
		resp := dto.ToStockQuoteResponse(q)
		responses = append(responses, &resp)
	}
	return responses, nil
}

// GetPriceHistory returns the most recent `limit` price history entries for a symbol.
func (s *StockAnalysisService) GetPriceHistory(ctx context.Context, symbol string, limit int) ([]*dto.StockPriceHistoryResponse, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	entries, err := s.stockHistoryRepo.FindBySymbol(ctx, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history for %s: %w", symbol, err)
	}

	responses := make([]*dto.StockPriceHistoryResponse, 0, len(entries))
	for _, e := range entries {
		resp := dto.ToStockPriceHistoryResponse(e)
		responses = append(responses, &resp)
	}
	return responses, nil
}

// ── private helpers ─────────────────────────────────────────────────────────

// fetchAndPersist calls the external provider, creates a new domain entity,
// upserts it into the cache, and appends a history entry.
func (s *StockAnalysisService) fetchAndPersist(ctx context.Context, symbol string) (*entities.StockQuote, error) {
	raw, err := s.quoteProvider.FetchQuote(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("symbol %s not found: %w", symbol, err)
	}

	quote, err := entities.NewStockQuote(raw.Symbol, raw.Name, raw.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock quote entity for %s: %w", symbol, err)
	}

	now := time.Now()
	quote.Update(
		raw.Name,
		toCents(raw.RegularMarketPrice),
		toCents(raw.PreviousClose),
		toCents(raw.Open),
		toCents(raw.DayHigh),
		toCents(raw.DayLow),
		raw.Volume,
		raw.MarketCap,
		toCents(raw.FiftyTwoWeekHigh),
		toCents(raw.FiftyTwoWeekLow),
		raw.TrailingPE,
		raw.ForwardPE,
		raw.DividendYield,
		raw.EPS,
		now,
	)

	if err := s.stockQuoteRepo.Upsert(ctx, quote); err != nil {
		return nil, fmt.Errorf("failed to upsert stock quote for %s: %w", symbol, err)
	}

	s.appendHistory(ctx, quote)

	return quote, nil
}

// refreshAndPersist updates an existing stale entity with fresh provider data.
func (s *StockAnalysisService) refreshAndPersist(ctx context.Context, existing *entities.StockQuote) (*entities.StockQuote, error) {
	raw, err := s.quoteProvider.FetchQuote(ctx, existing.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh quote for %s: %w", existing.Symbol, err)
	}

	now := time.Now()
	existing.Update(
		raw.Name,
		toCents(raw.RegularMarketPrice),
		toCents(raw.PreviousClose),
		toCents(raw.Open),
		toCents(raw.DayHigh),
		toCents(raw.DayLow),
		raw.Volume,
		raw.MarketCap,
		toCents(raw.FiftyTwoWeekHigh),
		toCents(raw.FiftyTwoWeekLow),
		raw.TrailingPE,
		raw.ForwardPE,
		raw.DividendYield,
		raw.EPS,
		now,
	)

	if err := s.stockQuoteRepo.Upsert(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to upsert refreshed quote for %s: %w", existing.Symbol, err)
	}

	s.appendHistory(ctx, existing)

	return existing, nil
}

// appendHistory creates a price-history entry for the given quote.
// Errors are logged but do not fail the parent operation — history is
// best-effort; the quote itself is what callers depend on.
func (s *StockAnalysisService) appendHistory(ctx context.Context, quote *entities.StockQuote) {
	if quote.PriceCents <= 0 {
		return
	}
	entry, err := entities.NewStockPriceHistory(quote.Symbol, quote.PriceCents, quote.Currency)
	if err != nil {
		return
	}
	// Ignore history persistence errors — best-effort
	_ = s.stockHistoryRepo.Create(ctx, entry)
}

// toCents converts a float64 price (e.g. 153.42) to int64 cents (15342).
// This is the standard conversion at the infrastructure boundary.
func toCents(price float64) int64 {
	return int64(price * 100)
}
