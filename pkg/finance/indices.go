package finance

// IndexStatus represents the health classification of a financial index.
type IndexStatus string

const (
	StatusGreen  IndexStatus = "green"
	StatusYellow IndexStatus = "yellow"
	StatusRed    IndexStatus = "red"
)

// IndexResult holds the computed value and its health status for a single index.
type IndexResult struct {
	Value  float64     `json:"value"`
	Status IndexStatus `json:"status"`
}

// SavingsRatio calculates what percentage of income is saved (income - expenses) / income × 100.
// Green > 20%, Yellow 10-20%, Red < 10%.
// All amounts are in cents (int64).
func SavingsRatio(incomeCents, expensesCents int64) IndexResult {
	if incomeCents == 0 {
		return IndexResult{Value: 0, Status: StatusRed}
	}

	ratio := float64(incomeCents-expensesCents) / float64(incomeCents) * 100

	var status IndexStatus
	switch {
	case ratio > 20:
		status = StatusGreen
	case ratio >= 10:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}

// DebtRatio calculates total debts / total income × 100.
// Green < 30%, Yellow 30-50%, Red > 50%.
// All amounts are in cents (int64).
func DebtRatio(totalDebtsCents, totalIncomeCents int64) IndexResult {
	if totalIncomeCents == 0 {
		return IndexResult{Value: 0, Status: StatusRed}
	}

	ratio := float64(totalDebtsCents) / float64(totalIncomeCents) * 100

	var status IndexStatus
	switch {
	case ratio < 30:
		status = StatusGreen
	case ratio <= 50:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}

// UnnecessaryExpenseRatio calculates unnecessary expenses / income × 100.
// Green < 10%, Yellow 10-25%, Red > 25%.
// All amounts are in cents (int64).
func UnnecessaryExpenseRatio(unnecessaryExpensesCents, incomeCents int64) IndexResult {
	if incomeCents == 0 {
		return IndexResult{Value: 0, Status: StatusRed}
	}

	ratio := float64(unnecessaryExpensesCents) / float64(incomeCents) * 100

	var status IndexStatus
	switch {
	case ratio < 10:
		status = StatusGreen
	case ratio <= 25:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}

// NetCashFlow calculates income - expenses - debtPayments.
// Value returned is the actual cash flow amount in cents (as float64).
// Green: result > 20% of income, Yellow: 0 to 20% of income, Red: negative.
// All amounts are in cents (int64).
func NetCashFlow(incomeCents, expensesCents, debtPaymentsCents int64) IndexResult {
	cashFlow := incomeCents - expensesCents - debtPaymentsCents
	value := float64(cashFlow)

	var status IndexStatus

	if incomeCents == 0 {
		switch {
		case cashFlow == 0:
			status = StatusYellow
		case cashFlow < 0:
			status = StatusRed
		default:
			status = StatusGreen
		}
		return IndexResult{Value: value, Status: status}
	}

	threshold := float64(incomeCents) * 0.20

	switch {
	case value > threshold:
		status = StatusGreen
	case value >= 0:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: value, Status: status}
}

// NetWorth calculates investments + accounts - debts.
// Value returned is net worth in cents (as float64).
// Green: positive, Yellow: within ±1000 cents of zero, Red: negative.
// All amounts are in cents (int64).
func NetWorth(investmentValueCents, accountBalancesCents, totalDebtsCents int64) IndexResult {
	netWorth := investmentValueCents + accountBalancesCents - totalDebtsCents
	value := float64(netWorth)

	const zeroBand = 1000.0

	var status IndexStatus
	switch {
	case value > zeroBand:
		status = StatusGreen
	case value >= -zeroBand:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: value, Status: status}
}

// EmergencyCoverage calculates how many months of expenses are covered by account balances.
// Green > 6 months, Yellow 3-6 months, Red < 3 months.
// All amounts are in cents (int64).
func EmergencyCoverage(accountBalancesCents, monthlyExpensesCents int64) IndexResult {
	if monthlyExpensesCents == 0 {
		return IndexResult{Value: 0, Status: StatusGreen}
	}

	months := float64(accountBalancesCents) / float64(monthlyExpensesCents)

	var status IndexStatus
	switch {
	case months > 6:
		status = StatusGreen
	case months >= 3:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: months, Status: status}
}

// InvestmentRatio calculates monthly investment / income × 100.
// Green > 15%, Yellow 5-15%, Red < 5%.
// All amounts are in cents (int64).
func InvestmentRatio(monthlyInvestmentCents, incomeCents int64) IndexResult {
	if incomeCents == 0 {
		return IndexResult{Value: 0, Status: StatusRed}
	}

	ratio := float64(monthlyInvestmentCents) / float64(incomeCents) * 100

	var status IndexStatus
	switch {
	case ratio > 15:
		status = StatusGreen
	case ratio >= 5:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}

// LiquidityRatio calculates account balances / monthly expenses (ratio).
// Green > 3, Yellow 1-3, Red < 1.
// All amounts are in cents (int64).
func LiquidityRatio(accountBalancesCents, monthlyExpensesCents int64) IndexResult {
	if monthlyExpensesCents == 0 {
		return IndexResult{Value: 0, Status: StatusGreen}
	}

	ratio := float64(accountBalancesCents) / float64(monthlyExpensesCents)

	var status IndexStatus
	switch {
	case ratio > 3:
		status = StatusGreen
	case ratio >= 1:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}

// PaymentCapacity calculates income / monthly debt payments (ratio).
// Green > 3, Yellow 1.5-3, Red < 1.5.
// All amounts are in cents (int64).
func PaymentCapacity(incomeCents, monthlyDebtPaymentsCents int64) IndexResult {
	if monthlyDebtPaymentsCents == 0 {
		return IndexResult{Value: 0, Status: StatusGreen}
	}

	ratio := float64(incomeCents) / float64(monthlyDebtPaymentsCents)

	var status IndexStatus
	switch {
	case ratio > 3:
		status = StatusGreen
	case ratio >= 1.5:
		status = StatusYellow
	default:
		status = StatusRed
	}

	return IndexResult{Value: ratio, Status: status}
}
