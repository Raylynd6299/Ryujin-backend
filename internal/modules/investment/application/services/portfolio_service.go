package services

import (
	"context"
	"fmt"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/application/dto"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/repositories"
)

// PortfolioService handles portfolio-level aggregation use cases
type PortfolioService struct {
	holdingRepo repositories.HoldingRepository
}

// NewPortfolioService creates a new PortfolioService
func NewPortfolioService(repo repositories.HoldingRepository) *PortfolioService {
	return &PortfolioService{holdingRepo: repo}
}

// GetSummary returns an aggregated portfolio summary grouped by currency
func (s *PortfolioService) GetSummary(ctx context.Context, userID string) (*dto.PortfolioSummaryResponse, error) {
	holdings, err := s.holdingRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve holdings for summary: %w", err)
	}

	// Accumulate per-currency totals using maps
	type currencyAccumulator struct {
		totalInvested     int64
		totalCurrentValue int64
		count             int
		withoutPrice      int
	}

	currencyMap := make(map[string]*currencyAccumulator)

	for _, h := range holdings {
		currency := h.Currency.Code()

		acc, ok := currencyMap[currency]
		if !ok {
			acc = &currencyAccumulator{}
			currencyMap[currency] = acc
		}

		acc.count++

		// Total invested = buyPrice (cents per unit) * quantity (micro-units) / 1_000_000
		invested := (h.BuyPrice.Amount() * h.Quantity.MicroUnits()) / 1_000_000
		acc.totalInvested += invested

		if h.CurrentPrice != nil {
			// Total current value = currentPrice (cents per unit) * quantity (micro-units) / 1_000_000
			currentValue := (h.CurrentPrice.Amount() * h.Quantity.MicroUnits()) / 1_000_000
			acc.totalCurrentValue += currentValue
		} else {
			acc.withoutPrice++
		}
	}

	subtotals := make([]dto.CurrencySubtotal, 0, len(currencyMap))
	for currency, acc := range currencyMap {
		gainLoss := acc.totalCurrentValue - acc.totalInvested

		var gainLossPct float64
		if acc.totalInvested > 0 {
			gainLossPct = float64(gainLoss) / float64(acc.totalInvested) * 100.0
		}

		subtotals = append(subtotals, dto.CurrencySubtotal{
			Currency:                currency,
			TotalInvestedCents:      acc.totalInvested,
			TotalCurrentValueCents:  acc.totalCurrentValue,
			UnrealizedGainLossCents: gainLoss,
			UnrealizedGainLossPct:   gainLossPct,
			HoldingsCount:           acc.count,
			HoldingsWithoutPrice:    acc.withoutPrice,
		})
	}

	return &dto.PortfolioSummaryResponse{
		Subtotals:     subtotals,
		TotalHoldings: len(holdings),
	}, nil
}

// GetPerformance returns a per-holding performance breakdown
func (s *PortfolioService) GetPerformance(ctx context.Context, userID string) (*dto.PortfolioPerformanceResponse, error) {
	holdings, err := s.holdingRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve holdings for performance: %w", err)
	}

	performances := make([]dto.HoldingPerformance, 0, len(holdings))
	for _, h := range holdings {
		// Total invested = buyPrice * quantity (int64 arithmetic)
		totalInvested := (h.BuyPrice.Amount() * h.Quantity.MicroUnits()) / 1_000_000

		perf := dto.HoldingPerformance{
			HoldingID:          h.ID,
			Symbol:             h.Symbol.Value(),
			Name:               h.Name,
			Currency:           h.Currency.Code(),
			TotalInvestedCents: totalInvested,
		}

		if h.CurrentPrice != nil {
			currentValue := (h.CurrentPrice.Amount() * h.Quantity.MicroUnits()) / 1_000_000
			perf.CurrentValueCents = &currentValue

			gainLoss := currentValue - totalInvested
			perf.UnrealizedGainLossCents = &gainLoss

			if totalInvested > 0 {
				gainLossPct := float64(gainLoss) / float64(totalInvested) * 100.0
				perf.UnrealizedGainLossPct = &gainLossPct
			}
		}

		performances = append(performances, perf)
	}

	return &dto.PortfolioPerformanceResponse{
		Holdings: performances,
	}, nil
}
