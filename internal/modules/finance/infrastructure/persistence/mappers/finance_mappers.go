package mappers

import (
	"github.com/Raylynd6299/ryujin/internal/modules/finance/domain/entities"
	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
	vo "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/value_objects"
	"github.com/Raylynd6299/ryujin/internal/modules/finance/infrastructure/persistence/models"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// ============================================================
// Category
// ============================================================

func CategoryToModel(c *entities.Category) *models.CategoryModel {
	return &models.CategoryModel{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Type:      string(c.Type),
		Icon:      c.Icon,
		Color:     c.Color,
		IsDefault: c.IsDefault,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func CategoryToDomain(m *models.CategoryModel) *entities.Category {
	return &entities.Category{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Type:      entities.CategoryType(m.Type),
		Icon:      m.Icon,
		Color:     m.Color,
		IsDefault: m.IsDefault,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ============================================================
// IncomeSource
// ============================================================

func IncomeSourceToModel(i *entities.IncomeSource) *models.IncomeSourceModel {
	return &models.IncomeSourceModel{
		ID:          i.ID,
		UserID:      i.UserID,
		CategoryID:  i.CategoryID,
		Name:        i.Name,
		Description: i.Description,
		AmountCents: i.Amount.Amount(),
		Currency:    i.Amount.Currency(),
		IncomeType:  i.IncomeType.String(),
		Recurrence:  i.Recurrence.String(),
		StartDate:   i.StartDate,
		EndDate:     i.EndDate,
		IsActive:    i.IsActive,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func IncomeSourceToDomain(m *models.IncomeSourceModel) (*entities.IncomeSource, error) {
	money, err := sharedVO.NewMoney(m.AmountCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid money in DB: " + err.Error())
	}

	it, err := vo.NewIncomeSourceType(m.IncomeType)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid income type in DB: " + err.Error())
	}

	rec, err := vo.NewRecurrence(m.Recurrence)
	if err != nil {
		return nil, financeErrors.NewIncomeSourceInvalidError("invalid recurrence in DB: " + err.Error())
	}

	return &entities.IncomeSource{
		ID:          m.ID,
		UserID:      m.UserID,
		CategoryID:  m.CategoryID,
		Name:        m.Name,
		Description: m.Description,
		Amount:      *money,
		IncomeType:  *it,
		Recurrence:  *rec,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

// ============================================================
// Expense
// ============================================================

func ExpenseToModel(e *entities.Expense) *models.ExpenseModel {
	return &models.ExpenseModel{
		ID:          e.ID,
		UserID:      e.UserID,
		CategoryID:  e.CategoryID,
		Name:        e.Name,
		Description: e.Description,
		AmountCents: e.Amount.Amount(),
		Currency:    e.Amount.Currency(),
		Priority:    e.Priority.String(),
		Recurrence:  e.Recurrence.String(),
		ExpenseDate: e.ExpenseDate,
		EndDate:     e.EndDate,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func ExpenseToDomain(m *models.ExpenseModel) (*entities.Expense, error) {
	money, err := sharedVO.NewMoney(m.AmountCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError("invalid money in DB: " + err.Error())
	}

	prio, err := vo.NewPriority(m.Priority)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError("invalid priority in DB: " + err.Error())
	}

	rec, err := vo.NewRecurrence(m.Recurrence)
	if err != nil {
		return nil, financeErrors.NewExpenseInvalidError("invalid recurrence in DB: " + err.Error())
	}

	return &entities.Expense{
		ID:          m.ID,
		UserID:      m.UserID,
		CategoryID:  m.CategoryID,
		Name:        m.Name,
		Description: m.Description,
		Amount:      *money,
		Priority:    *prio,
		Recurrence:  *rec,
		ExpenseDate: m.ExpenseDate,
		EndDate:     m.EndDate,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

// ============================================================
// Debt
// ============================================================

func DebtToModel(d *entities.Debt) *models.DebtModel {
	return &models.DebtModel{
		ID:                   d.ID,
		UserID:               d.UserID,
		Name:                 d.Name,
		Description:          d.Description,
		DebtType:             d.DebtType.String(),
		TotalAmountCents:     d.TotalAmount.Amount(),
		RemainingAmountCents: d.RemainingAmount.Amount(),
		MonthlyPaymentCents:  d.MonthlyPayment.Amount(),
		Currency:             d.TotalAmount.Currency(),
		InterestRate:         d.InterestRate,
		StartDate:            d.StartDate,
		DueDate:              d.DueDate,
		IsActive:             d.IsActive,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
}

func DebtToDomain(m *models.DebtModel) (*entities.Debt, error) {
	totalMoney, err := sharedVO.NewMoney(m.TotalAmountCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid total money in DB: " + err.Error())
	}

	remainingMoney, err := sharedVO.NewMoney(m.RemainingAmountCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid remaining money in DB: " + err.Error())
	}

	monthlyMoney, err := sharedVO.NewMoney(m.MonthlyPaymentCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid monthly payment in DB: " + err.Error())
	}

	dt, err := vo.NewDebtCategory(m.DebtType)
	if err != nil {
		return nil, financeErrors.NewDebtInvalidError("invalid debt type in DB: " + err.Error())
	}

	return &entities.Debt{
		ID:              m.ID,
		UserID:          m.UserID,
		Name:            m.Name,
		Description:     m.Description,
		DebtType:        *dt,
		TotalAmount:     *totalMoney,
		RemainingAmount: *remainingMoney,
		MonthlyPayment:  *monthlyMoney,
		InterestRate:    m.InterestRate,
		StartDate:       m.StartDate,
		DueDate:         m.DueDate,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}, nil
}

// ============================================================
// Account
// ============================================================

func AccountToModel(a *entities.Account) *models.AccountModel {
	return &models.AccountModel{
		ID:           a.ID,
		UserID:       a.UserID,
		Name:         a.Name,
		Description:  a.Description,
		AccountType:  string(a.AccountType),
		BalanceCents: a.Balance.Amount(),
		Currency:     a.Balance.Currency(),
		IsActive:     a.IsActive,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func AccountToDomain(m *models.AccountModel) (*entities.Account, error) {
	money, err := sharedVO.NewMoney(m.BalanceCents, m.Currency)
	if err != nil {
		return nil, financeErrors.NewAccountInvalidError("invalid money in DB: " + err.Error())
	}

	return &entities.Account{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Description: m.Description,
		AccountType: entities.AccountType(m.AccountType),
		Balance:     *money,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}
