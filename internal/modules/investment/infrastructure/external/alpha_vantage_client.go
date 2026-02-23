package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	alphaVantageBaseURL = "https://www.alphavantage.co/query"
	alphaVantageTimeout = 5 * time.Second
)

// alphaVantageGlobalQuoteResponse represents the API response for GLOBAL_QUOTE endpoint
type alphaVantageGlobalQuoteResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Price            string `json:"05. price"`
		LatestTradingDay string `json:"07. latest trading day"`
	} `json:"Global Quote"`
}

// AlphaVantageClient fetches prices from the Alpha Vantage API
type AlphaVantageClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewAlphaVantageClient creates a new Alpha Vantage price provider
func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: alphaVantageTimeout,
		},
	}
}

// FetchPrice fetches the current price for a symbol from Alpha Vantage.
// Currency always defaults to "USD" as the GLOBAL_QUOTE endpoint does not return currency.
func (c *AlphaVantageClient) FetchPrice(ctx context.Context, symbol string) (*PriceQuote, error) {
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		alphaVantageBaseURL, symbol, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("alpha_vantage: failed to build request for %s: %w", symbol, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alpha_vantage: request failed for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alpha_vantage: unexpected status %d for symbol %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("alpha_vantage: failed to read response body for %s: %w", symbol, err)
	}

	var data alphaVantageGlobalQuoteResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("alpha_vantage: failed to parse JSON for %s: %w", symbol, err)
	}

	priceStr := strings.TrimSpace(data.GlobalQuote.Price)
	if priceStr == "" {
		return nil, fmt.Errorf("alpha_vantage: empty price returned for symbol %s (check API key or rate limit)", symbol)
	}

	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("alpha_vantage: could not parse price %q for symbol %s: %w", priceStr, symbol, err)
	}

	if priceFloat <= 0 {
		return nil, fmt.Errorf("alpha_vantage: invalid price %.4f for symbol %s", priceFloat, symbol)
	}

	priceCents := int64(math.Round(priceFloat * 100))

	return &PriceQuote{
		Symbol:     symbol,
		PriceCents: priceCents,
		Currency:   "USD", // Alpha Vantage GLOBAL_QUOTE does not return currency
		FetchedAt:  time.Now(),
	}, nil
}

// Name returns the identifier for this price provider
func (c *AlphaVantageClient) Name() string {
	return "alpha_vantage"
}
