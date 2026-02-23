package dto

// CurrencySubtotal groups portfolio totals for a single currency
type CurrencySubtotal struct {
	Currency                string  `json:"currency"`
	TotalInvestedCents      int64   `json:"totalInvestedCents"`
	TotalCurrentValueCents  int64   `json:"totalCurrentValueCents"`
	UnrealizedGainLossCents int64   `json:"unrealizedGainLossCents"`
	UnrealizedGainLossPct   float64 `json:"unrealizedGainLossPct"`
	HoldingsCount           int     `json:"holdingsCount"`
	HoldingsWithoutPrice    int     `json:"holdingsWithoutPrice"`
}

// PortfolioSummaryResponse is returned for portfolio summary requests
type PortfolioSummaryResponse struct {
	Subtotals     []CurrencySubtotal `json:"subtotals"`
	TotalHoldings int                `json:"totalHoldings"`
}

// HoldingPerformance represents the performance breakdown for a single holding
type HoldingPerformance struct {
	HoldingID               string   `json:"holdingId"`
	Symbol                  string   `json:"symbol"`
	Name                    string   `json:"name"`
	Currency                string   `json:"currency"`
	TotalInvestedCents      int64    `json:"totalInvestedCents"`
	CurrentValueCents       *int64   `json:"currentValueCents"`
	UnrealizedGainLossCents *int64   `json:"unrealizedGainLossCents"`
	UnrealizedGainLossPct   *float64 `json:"unrealizedGainLossPct"`
}

// PortfolioPerformanceResponse is returned for portfolio performance requests
type PortfolioPerformanceResponse struct {
	Holdings []HoldingPerformance `json:"holdings"`
}
