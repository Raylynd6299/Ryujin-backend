package dto

// IndexDTO represents a single computed financial health index.
type IndexDTO struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Status string  `json:"status"` // "green" | "yellow" | "red"
	Label  string  `json:"label"`
}

// IndicesResponseDTO wraps all financial indices for the API response.
type IndicesResponseDTO struct {
	Indices         []IndexDTO `json:"indices"`
	CurrencyWarning bool       `json:"currencyWarning"` // true if user has mixed currencies
}

// FinanceSummaryDTO provides a high-level financial summary with both
// raw cents values and human-readable decimal equivalents.
type FinanceSummaryDTO struct {
	TotalIncomeCents     int64   `json:"totalIncomeCents"`
	TotalIncomeDecimal   float64 `json:"totalIncomeDecimal"`
	TotalExpensesCents   int64   `json:"totalExpensesCents"`
	TotalExpensesDecimal float64 `json:"totalExpensesDecimal"`
	NetCashFlowCents     int64   `json:"netCashFlowCents"`
	NetCashFlowDecimal   float64 `json:"netCashFlowDecimal"`
	SavingsAmountCents   int64   `json:"savingsAmountCents"`
	SavingsAmountDecimal float64 `json:"savingsAmountDecimal"`
	Currency             string  `json:"currency"`
}
