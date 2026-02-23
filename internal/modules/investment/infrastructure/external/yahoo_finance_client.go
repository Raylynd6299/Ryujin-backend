package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	yahooFinanceBaseURL  = "https://query1.finance.yahoo.com/v8/finance/chart"
	yahooFinanceQuoteURL = "https://query1.finance.yahoo.com/v8/finance/quoteSummary"
	yahooFinanceTimeout  = 5 * time.Second
)

// yahooChartResponse is the top-level response envelope from the Yahoo Finance chart API
type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				Currency           string  `json:"currency"`
			} `json:"meta"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// yahooQuoteSummaryResponse is the response envelope from the Yahoo Finance quoteSummary API
type yahooQuoteSummaryResponse struct {
	QuoteSummary struct {
		Result []struct {
			Price struct {
				Symbol             string `json:"symbol"`
				ShortName          string `json:"shortName"`
				Currency           string `json:"currency"`
				RegularMarketPrice struct {
					Raw float64 `json:"raw"`
				} `json:"regularMarketPrice"`
				RegularMarketPreviousClose struct {
					Raw float64 `json:"raw"`
				} `json:"regularMarketPreviousClose"`
				RegularMarketOpen struct {
					Raw float64 `json:"raw"`
				} `json:"regularMarketOpen"`
				RegularMarketDayHigh struct {
					Raw float64 `json:"raw"`
				} `json:"regularMarketDayHigh"`
				RegularMarketDayLow struct {
					Raw float64 `json:"raw"`
				} `json:"regularMarketDayLow"`
				RegularMarketVolume struct {
					Raw int64 `json:"raw"`
				} `json:"regularMarketVolume"`
				MarketCap struct {
					Raw int64 `json:"raw"`
				} `json:"marketCap"`
			} `json:"price"`
			SummaryDetail struct {
				FiftyTwoWeekHigh struct {
					Raw float64 `json:"raw"`
				} `json:"fiftyTwoWeekHigh"`
				FiftyTwoWeekLow struct {
					Raw float64 `json:"raw"`
				} `json:"fiftyTwoWeekLow"`
				TrailingPE struct {
					Raw float64 `json:"raw"`
				} `json:"trailingPE"`
				ForwardPE struct {
					Raw float64 `json:"raw"`
				} `json:"forwardPE"`
				DividendYield struct {
					Raw float64 `json:"raw"`
				} `json:"dividendYield"`
			} `json:"summaryDetail"`
			DefaultKeyStatistics struct {
				TrailingEps struct {
					Raw float64 `json:"raw"`
				} `json:"trailingEps"`
			} `json:"defaultKeyStatistics"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"quoteSummary"`
}

// YahooFinanceClient fetches real-time prices from Yahoo Finance
type YahooFinanceClient struct {
	httpClient *http.Client
}

// NewYahooFinanceClient creates a new Yahoo Finance price provider
func NewYahooFinanceClient() *YahooFinanceClient {
	return &YahooFinanceClient{
		httpClient: &http.Client{
			Timeout: yahooFinanceTimeout,
		},
	}
}

// FetchPrice fetches the current market price for a symbol from Yahoo Finance.
// Price is returned in cents (price × 100, rounded to nearest integer).
func (c *YahooFinanceClient) FetchPrice(ctx context.Context, symbol string) (*PriceQuote, error) {
	url := fmt.Sprintf("%s/%s?interval=1d&range=1d", yahooFinanceBaseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to build request for %s: %w", symbol, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: request failed for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo_finance: unexpected status %d for symbol %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to read response body for %s: %w", symbol, err)
	}

	var data yahooChartResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to parse JSON for %s: %w", symbol, err)
	}

	if data.Chart.Error != nil {
		return nil, fmt.Errorf("yahoo_finance: API error for %s: %s - %s",
			symbol, data.Chart.Error.Code, data.Chart.Error.Description)
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("yahoo_finance: no result returned for symbol %s", symbol)
	}

	meta := data.Chart.Result[0].Meta
	if meta.RegularMarketPrice <= 0 {
		return nil, fmt.Errorf("yahoo_finance: invalid price %.4f for symbol %s", meta.RegularMarketPrice, symbol)
	}

	currency := meta.Currency
	if currency == "" {
		currency = "USD"
	}

	priceCents := int64(math.Round(meta.RegularMarketPrice * 100))

	return &PriceQuote{
		Symbol:     symbol,
		PriceCents: priceCents,
		Currency:   currency,
		FetchedAt:  time.Now(),
	}, nil
}

// Name returns the identifier for this price provider
func (c *YahooFinanceClient) Name() string {
	return "yahoo_finance"
}

// FetchQuote fetches enriched market data for a symbol from Yahoo Finance.
// Implements the QuoteProvider port.
func (c *YahooFinanceClient) FetchQuote(ctx context.Context, symbol string) (*StockQuote, error) {
	url := fmt.Sprintf("%s/%s?modules=price,summaryDetail,defaultKeyStatistics", yahooFinanceQuoteURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to build quote request for %s: %w", symbol, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: quote request failed for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo_finance: unexpected status %d for symbol %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to read quote response for %s: %w", symbol, err)
	}

	var data yahooQuoteSummaryResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("yahoo_finance: failed to parse quote JSON for %s: %w", symbol, err)
	}

	if data.QuoteSummary.Error != nil {
		return nil, fmt.Errorf("yahoo_finance: API error for %s: %s - %s",
			symbol, data.QuoteSummary.Error.Code, data.QuoteSummary.Error.Description)
	}

	if len(data.QuoteSummary.Result) == 0 {
		return nil, fmt.Errorf("yahoo_finance: no quote result for symbol %s", symbol)
	}

	r := data.QuoteSummary.Result[0]
	price := r.Price
	summary := r.SummaryDetail
	stats := r.DefaultKeyStatistics

	if price.RegularMarketPrice.Raw <= 0 {
		return nil, fmt.Errorf("yahoo_finance: invalid market price for symbol %s", symbol)
	}

	currency := price.Currency
	if currency == "" {
		currency = "USD"
	}

	_ = math.Round // keep import used
	return &StockQuote{
		Symbol:             price.Symbol,
		Name:               price.ShortName,
		Currency:           currency,
		RegularMarketPrice: price.RegularMarketPrice.Raw,
		PreviousClose:      price.RegularMarketPreviousClose.Raw,
		Open:               price.RegularMarketOpen.Raw,
		DayHigh:            price.RegularMarketDayHigh.Raw,
		DayLow:             price.RegularMarketDayLow.Raw,
		Volume:             price.RegularMarketVolume.Raw,
		MarketCap:          price.MarketCap.Raw,
		FiftyTwoWeekHigh:   summary.FiftyTwoWeekHigh.Raw,
		FiftyTwoWeekLow:    summary.FiftyTwoWeekLow.Raw,
		TrailingPE:         summary.TrailingPE.Raw,
		ForwardPE:          summary.ForwardPE.Raw,
		DividendYield:      summary.DividendYield.Raw,
		EPS:                stats.TrailingEps.Raw,
		FetchedAt:          time.Now(),
	}, nil
}
