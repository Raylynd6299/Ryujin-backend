package services

import (
	"context"
	"fmt"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/dto"
	financeRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/domain/repositories"
	investRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/repositories"
	pkgFinance "github.com/Raylynd6299/Ryujin-backend/pkg/finance"
)

// IndicesCalculatorService computes the financial health indices for a user.
// It aggregates data from income, expense, debt, account, and (optionally) investment repos.
type IndicesCalculatorService struct {
	incomeRepo  financeRepos.IncomeSourceRepository
	expenseRepo financeRepos.ExpenseRepository
	debtRepo    financeRepos.DebtRepository
	accountRepo financeRepos.AccountRepository
	holdingRepo investRepos.HoldingRepository // optional — may be nil
}

// NewIndicesCalculatorService creates a new IndicesCalculatorService.
// holdingRepo is optional; pass nil if investment data is not available.
func NewIndicesCalculatorService(
	incomeRepo financeRepos.IncomeSourceRepository,
	expenseRepo financeRepos.ExpenseRepository,
	debtRepo financeRepos.DebtRepository,
	accountRepo financeRepos.AccountRepository,
	holdingRepo investRepos.HoldingRepository,
) *IndicesCalculatorService {
	return &IndicesCalculatorService{
		incomeRepo:  incomeRepo,
		expenseRepo: expenseRepo,
		debtRepo:    debtRepo,
		accountRepo: accountRepo,
		holdingRepo: holdingRepo,
	}
}

// CalculateIndices computes all 9 financial health indices for the given user.
func (s *IndicesCalculatorService) CalculateIndices(ctx context.Context, userID string) (*dto.IndicesResponseDTO, error) {
	totals, err := s.gatherTotals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to gather financial totals: %w", err)
	}

	indices := []dto.IndexDTO{
		buildIndex("savings_ratio", "Savings Ratio",
			pkgFinance.SavingsRatio(totals.monthlyIncomeCents, totals.monthlyExpensesCents)),
		buildIndex("debt_ratio", "Debt Ratio",
			pkgFinance.DebtRatio(totals.totalDebtCents, totals.monthlyIncomeCents)),
		buildIndex("unnecessary_expense_ratio", "Unnecessary Expense Ratio",
			pkgFinance.UnnecessaryExpenseRatio(totals.unnecessaryExpensesCents, totals.monthlyIncomeCents)),
		buildIndex("net_cash_flow", "Net Cash Flow",
			pkgFinance.NetCashFlow(totals.monthlyIncomeCents, totals.monthlyExpensesCents, totals.monthlyDebtPaymentsCents)),
		buildIndex("net_worth", "Net Worth",
			pkgFinance.NetWorth(totals.investmentValueCents, totals.liquidAssetsCents, totals.totalDebtCents)),
		buildIndex("emergency_coverage", "Emergency Coverage",
			pkgFinance.EmergencyCoverage(totals.liquidAssetsCents, totals.monthlyExpensesCents)),
		buildIndex("investment_ratio", "Investment Ratio",
			pkgFinance.InvestmentRatio(totals.investmentValueCents, totals.monthlyIncomeCents)),
		buildIndex("liquidity_ratio", "Liquidity Ratio",
			pkgFinance.LiquidityRatio(totals.liquidAssetsCents, totals.monthlyExpensesCents)),
		buildIndex("payment_capacity", "Payment Capacity",
			pkgFinance.PaymentCapacity(totals.monthlyIncomeCents, totals.monthlyDebtPaymentsCents)),
	}

	return &dto.IndicesResponseDTO{
		Indices:         indices,
		CurrencyWarning: totals.mixedCurrencies,
	}, nil
}

// GetSummary returns a high-level monthly financial summary for the given user.
func (s *IndicesCalculatorService) GetSummary(ctx context.Context, userID string) (*dto.FinanceSummaryDTO, error) {
	totals, err := s.gatherTotals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to gather financial totals: %w", err)
	}

	// Net cash flow = income - expenses - debt payments
	netCashFlowCents := totals.monthlyIncomeCents - totals.monthlyExpensesCents - totals.monthlyDebtPaymentsCents

	// Savings = income - all expenses (including debt payments)
	savingsAmountCents := totals.monthlyIncomeCents - totals.monthlyExpensesCents

	return &dto.FinanceSummaryDTO{
		TotalIncomeCents:     totals.monthlyIncomeCents,
		TotalIncomeDecimal:   centsToDecimal(totals.monthlyIncomeCents),
		TotalExpensesCents:   totals.monthlyExpensesCents,
		TotalExpensesDecimal: centsToDecimal(totals.monthlyExpensesCents),
		NetCashFlowCents:     netCashFlowCents,
		NetCashFlowDecimal:   centsToDecimal(netCashFlowCents),
		SavingsAmountCents:   savingsAmountCents,
		SavingsAmountDecimal: centsToDecimal(savingsAmountCents),
		Currency:             totals.primaryCurrency,
	}, nil
}

// financialTotals holds all aggregated monetary values needed for index calculation.
type financialTotals struct {
	monthlyIncomeCents       int64
	monthlyExpensesCents     int64
	unnecessaryExpensesCents int64
	monthlyDebtPaymentsCents int64
	totalDebtCents           int64
	liquidAssetsCents        int64
	investmentValueCents     int64
	mixedCurrencies          bool
	primaryCurrency          string
}

// gatherTotals fetches and aggregates all financial data for a user.
func (s *IndicesCalculatorService) gatherTotals(ctx context.Context, userID string) (*financialTotals, error) {
	totals := &financialTotals{}

	// ── Income ───────────────────────────────────────────────────────────────
	incomes, err := s.incomeRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch income sources: %w", err)
	}

	currencies := make(map[string]struct{})
	var primaryCurrency string

	for _, income := range incomes {
		monthly := income.MonthlyEquivalent()
		totals.monthlyIncomeCents += monthly.Amount()
		cur := monthly.Currency()
		currencies[cur] = struct{}{}
		if primaryCurrency == "" {
			primaryCurrency = cur
		}
	}

	// ── Expenses ─────────────────────────────────────────────────────────────
	expenses, err := s.expenseRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %w", err)
	}

	for _, expense := range expenses {
		monthly := expense.MonthlyEquivalent()
		totals.monthlyExpensesCents += monthly.Amount()
		currencies[monthly.Currency()] = struct{}{}

		if expense.IsUnnecessary() {
			totals.unnecessaryExpensesCents += monthly.Amount()
		}
	}

	// ── Debts ────────────────────────────────────────────────────────────────
	debts, err := s.debtRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch debts: %w", err)
	}

	for _, debt := range debts {
		totals.totalDebtCents += debt.RemainingAmount.Amount()
		totals.monthlyDebtPaymentsCents += debt.MonthlyPayment.Amount()
		currencies[debt.RemainingAmount.Currency()] = struct{}{}
	}

	// ── Accounts ─────────────────────────────────────────────────────────────
	accounts, err := s.accountRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	for _, account := range accounts {
		totals.liquidAssetsCents += account.Balance.Amount()
		currencies[account.Balance.Currency()] = struct{}{}
	}

	// ── Holdings (optional) ──────────────────────────────────────────────────
	if s.holdingRepo != nil {
		holdings, err := s.holdingRepo.FindActiveByUserID(ctx, userID)
		if err != nil {
			// Non-fatal — degrade gracefully with zero investment value
			holdings = nil
		}
		for _, holding := range holdings {
			mv := holding.MarketValue()
			if mv != nil {
				totals.investmentValueCents += mv.Amount()
				currencies[mv.Currency()] = struct{}{}
			} else {
				// Fall back to cost basis if no market price is available
				costBasis := holding.BuyPrice.Multiply(holding.Quantity.ToFloat())
				totals.investmentValueCents += costBasis.Amount()
				currencies[holding.BuyPrice.Currency()] = struct{}{}
			}
		}
	}

	// ── Currency warning ─────────────────────────────────────────────────────
	totals.mixedCurrencies = len(currencies) > 1
	totals.primaryCurrency = primaryCurrency
	if totals.primaryCurrency == "" {
		totals.primaryCurrency = "USD"
	}

	return totals, nil
}

// buildIndex creates an IndexDTO from an IndexResult.
func buildIndex(name, label string, result pkgFinance.IndexResult) dto.IndexDTO {
	return dto.IndexDTO{
		Name:   name,
		Value:  result.Value,
		Status: string(result.Status),
		Label:  label,
	}
}
