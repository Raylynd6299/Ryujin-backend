package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	domainRepos "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/repositories"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/external"
)

// HoldingService handles investment holding use cases
type HoldingService struct {
	holdingRepo    domainRepos.HoldingRepository
	stockQuoteRepo domainRepos.StockQuoteRepository
	priceProvider  external.PriceProvider
	quoteProvider  external.QuoteProvider
}

// NewHoldingService creates a new HoldingService
func NewHoldingService(
	repo domainRepos.HoldingRepository,
	provider external.PriceProvider,
	stockQuoteRepo domainRepos.StockQuoteRepository,
	quoteProvider external.QuoteProvider,
) *HoldingService {
	return &HoldingService{
		holdingRepo:    repo,
		stockQuoteRepo: stockQuoteRepo,
		priceProvider:  provider,
		quoteProvider:  quoteProvider,
	}
}

// CreateHolding creates a new investment holding for a user.
// Before persisting the holding, it ensures the symbol exists in the stock_quotes
// cache (inserting a fresh quote if missing) so that the FK constraint is satisfied.
func (s *HoldingService) CreateHolding(ctx context.Context, userID string, req dto.CreateHoldingRequest) (*dto.HoldingResponse, error) {
	// Normalize symbol before any lookup
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))

	// Ensure the symbol is seeded in stock_quotes (satisfies FK)
	_, err := s.stockQuoteRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		if !errors.Is(err, investErrors.ErrStockQuoteNotFound) {
			return nil, fmt.Errorf("failed to look up symbol %s: %w", symbol, err)
		}

		// Symbol not cached — fetch from provider to validate it exists
		raw, fetchErr := s.quoteProvider.FetchQuote(ctx, symbol)
		if fetchErr != nil {
			return nil, fmt.Errorf("symbol %s not found: %w", symbol, fetchErr)
		}

		quote, createErr := entities.NewStockQuote(raw.Symbol, raw.Name, raw.Currency)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create stock quote entity for %s: %w", symbol, createErr)
		}

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
			raw.FetchedAt,
		)

		if upsertErr := s.stockQuoteRepo.Upsert(ctx, quote); upsertErr != nil {
			return nil, fmt.Errorf("failed to seed stock quote for %s: %w", symbol, upsertErr)
		}
	}
	// If found (err == nil) → symbol already cached, FK satisfied; continue.

	id := uuid.New().String()

	holding, err := entities.NewHolding(
		id,
		userID,
		symbol,
		req.Name,
		req.AssetType,
		req.QuantityMicro,
		req.BuyPriceCents,
		req.Currency,
		req.Notes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.holdingRepo.Create(ctx, holding); err != nil {
		return nil, fmt.Errorf("failed to create holding: %w", err)
	}

	resp := dto.ToHoldingResponse(holding)
	return &resp, nil
}

// GetHolding returns a holding by ID scoped to the given user
func (s *HoldingService) GetHolding(ctx context.Context, id, userID string) (*dto.HoldingResponse, error) {
	holding, err := s.holdingRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	resp := dto.ToHoldingResponse(holding)
	return &resp, nil
}

// ListHoldings returns paginated holdings for a user
func (s *HoldingService) ListHoldings(ctx context.Context, userID string, page, limit int, sort, order string) (*dto.HoldingListResponse, error) {
	holdings, total, err := s.holdingRepo.FindAllByUserID(ctx, userID, page, limit, sort, order)
	if err != nil {
		return nil, fmt.Errorf("failed to list holdings: %w", err)
	}

	responses := make([]dto.HoldingResponse, 0, len(holdings))
	for _, h := range holdings {
		responses = append(responses, dto.ToHoldingResponse(h))
	}

	return &dto.HoldingListResponse{
		Holdings: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// UpdateHolding updates the mutable fields of an existing holding
func (s *HoldingService) UpdateHolding(ctx context.Context, id, userID string, req dto.UpdateHoldingRequest) (*dto.HoldingResponse, error) {
	holding, err := s.holdingRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Use current values as fallback when optional fields are zero-valued
	name := req.Name
	if name == "" {
		name = holding.Name
	}

	quantityMicro := req.QuantityMicro
	if quantityMicro == 0 {
		quantityMicro = holding.Quantity.MicroUnits()
	}

	buyPriceCents := req.BuyPriceCents
	if buyPriceCents == 0 {
		buyPriceCents = holding.BuyPrice.Amount()
	}

	currency := req.Currency
	if currency == "" {
		currency = holding.Currency.Code()
	}

	if err := holding.Update(name, quantityMicro, buyPriceCents, currency, req.Notes); err != nil {
		return nil, err
	}

	if err := s.holdingRepo.Update(ctx, holding); err != nil {
		return nil, fmt.Errorf("failed to update holding: %w", err)
	}

	resp := dto.ToHoldingResponse(holding)
	return &resp, nil
}

// DeleteHolding removes a holding scoped to the given user
func (s *HoldingService) DeleteHolding(ctx context.Context, id, userID string) error {
	// Verify ownership before deletion
	if _, err := s.holdingRepo.FindByID(ctx, id, userID); err != nil {
		return err
	}
	return s.holdingRepo.Delete(ctx, id, userID)
}

// RefreshHoldingPrice fetches the latest market price and updates the holding
func (s *HoldingService) RefreshHoldingPrice(ctx context.Context, id, userID string) (*dto.HoldingResponse, error) {
	// Verify ownership first
	holding, err := s.holdingRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Fetch current market price from external provider
	quote, err := s.priceProvider.FetchPrice(ctx, holding.Symbol.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price for %s: %w", holding.Symbol.Value(), err)
	}

	// Update the domain entity with the new price
	holding.RefreshPrice(quote.PriceCents, quote.FetchedAt)

	// Persist the updated holding
	if err := s.holdingRepo.Update(ctx, holding); err != nil {
		return nil, fmt.Errorf("failed to persist refreshed price: %w", err)
	}

	resp := dto.ToHoldingResponse(holding)
	return &resp, nil
}
