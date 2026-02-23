package external

import (
	"context"
	"fmt"
)

// CompositePriceProvider tries a primary provider first and falls back to a secondary.
// If both fail, the errors from both are combined and returned.
type CompositePriceProvider struct {
	primary  PriceProvider
	fallback PriceProvider
}

// NewCompositePriceProvider creates a composite provider with primary and fallback
func NewCompositePriceProvider(primary, fallback PriceProvider) *CompositePriceProvider {
	return &CompositePriceProvider{
		primary:  primary,
		fallback: fallback,
	}
}

// FetchPrice attempts to fetch a price from the primary provider.
// On failure it tries the fallback. If both fail, it returns a combined error.
func (c *CompositePriceProvider) FetchPrice(ctx context.Context, symbol string) (*PriceQuote, error) {
	quote, primaryErr := c.primary.FetchPrice(ctx, symbol)
	if primaryErr == nil {
		return quote, nil
	}

	quote, fallbackErr := c.fallback.FetchPrice(ctx, symbol)
	if fallbackErr == nil {
		return quote, nil
	}

	return nil, fmt.Errorf("all price providers failed for symbol %s: [%s: %w] [%s: %v]",
		symbol,
		c.primary.Name(), primaryErr,
		c.fallback.Name(), fallbackErr,
	)
}

// Name returns a composite name describing both providers
func (c *CompositePriceProvider) Name() string {
	return fmt.Sprintf("composite(%s/%s)", c.primary.Name(), c.fallback.Name())
}
