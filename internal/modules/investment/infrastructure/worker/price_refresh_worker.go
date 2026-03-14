package worker

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/entities"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/repositories"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/external"
)

// PriceRefreshWorker periodically refreshes prices for all symbols with active holdings.
// It fetches enriched quote data from a QuoteProvider and persists the results via
// the domain repository ports.
type PriceRefreshWorker struct {
	stockQuoteRepo   repositories.StockQuoteRepository
	stockHistoryRepo repositories.StockPriceHistoryRepository
	quoteProvider    external.QuoteProvider
	interval         time.Duration
}

// NewPriceRefreshWorker creates a new PriceRefreshWorker with the given dependencies.
func NewPriceRefreshWorker(
	stockQuoteRepo repositories.StockQuoteRepository,
	stockHistoryRepo repositories.StockPriceHistoryRepository,
	quoteProvider external.QuoteProvider,
	interval time.Duration,
) *PriceRefreshWorker {
	return &PriceRefreshWorker{
		stockQuoteRepo:   stockQuoteRepo,
		stockHistoryRepo: stockHistoryRepo,
		quoteProvider:    quoteProvider,
		interval:         interval,
	}
}

// Start launches the worker in a goroutine and returns immediately.
// The worker runs until ctx is cancelled.
func (w *PriceRefreshWorker) Start(ctx context.Context) {
	go w.run(ctx)
}

// run is the main loop. It performs an immediate first refresh, then repeats on
// every tick. Panics inside refresh are recovered so the loop stays alive.
func (w *PriceRefreshWorker) run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[price_refresh_worker] recovered from panic: %v", r)
		}
	}()

	// Immediate first refresh so callers don't wait a full interval on startup.
	w.refresh(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[price_refresh_worker] context cancelled, shutting down")
			return
		case <-ticker.C:
			w.refresh(ctx)
		}
	}
}

// refresh performs one full refresh cycle:
//  1. Find all symbols that have at least one active holding.
//  2. Fetch enriched quote data for each symbol from the external provider.
//  3. Upsert the quote into stock_quotes.
//  4. Append an entry to stock_price_history.
//
// Errors are logged per-symbol and the loop continues so a single bad symbol
// does not block the rest. A 300 ms sleep between symbols avoids API rate limits.
func (w *PriceRefreshWorker) refresh(ctx context.Context) {
	symbols, err := w.stockQuoteRepo.FindSymbolsWithActiveHoldings(ctx)
	if err != nil {
		log.Printf("[price_refresh_worker] failed to fetch active symbols: %v", err)
		return
	}

	if len(symbols) == 0 {
		return
	}

	log.Printf("[price_refresh_worker] refreshing %d symbol(s)", len(symbols))

	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := w.refreshSymbol(ctx, symbol); err != nil {
			log.Printf("[price_refresh_worker] symbol %s: %v", symbol, err)
		}

		// Avoid hitting rate limits on free-tier APIs.
		time.Sleep(300 * time.Millisecond)
	}
}

// refreshSymbol fetches and persists data for a single symbol.
func (w *PriceRefreshWorker) refreshSymbol(ctx context.Context, symbol string) error {
	quote, err := w.quoteProvider.FetchQuote(ctx, symbol)
	if err != nil {
		return err
	}

	// Convert float64 prices from the external DTO to int64 cents for the domain entity.
	priceCents := toCents(quote.RegularMarketPrice)
	previousCloseCents := toCents(quote.PreviousClose)
	openCents := toCents(quote.Open)
	dayHighCents := toCents(quote.DayHigh)
	dayLowCents := toCents(quote.DayLow)
	week52HighCents := toCents(quote.FiftyTwoWeekHigh)
	week52LowCents := toCents(quote.FiftyTwoWeekLow)

	sq, err := entities.NewStockQuote(quote.Symbol, quote.Name, quote.Currency)
	if err != nil {
		return err
	}

	sq.Update(
		quote.Name,
		priceCents,
		previousCloseCents,
		openCents,
		dayHighCents,
		dayLowCents,
		quote.Volume,
		quote.MarketCap,
		week52HighCents,
		week52LowCents,
		quote.TrailingPE,
		quote.ForwardPE,
		quote.DividendYield,
		quote.EPS,
		quote.FetchedAt,
	)

	if err := w.stockQuoteRepo.Upsert(ctx, sq); err != nil {
		return err
	}

	history, err := entities.NewStockPriceHistory(sq.Symbol, priceCents, sq.Currency)
	if err != nil {
		return err
	}

	return w.stockHistoryRepo.Create(ctx, history)
}

// toCents converts a float64 price to int64 cents, rounding to the nearest integer.
func toCents(price float64) int64 {
	return int64(math.Round(price * 100))
}
