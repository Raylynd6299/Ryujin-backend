package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/infrastructure/http/controllers"
	"github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	"github.com/Raylynd6299/ryujin/internal/shared/utils"
)

// RegisterRoutes registers all finance module routes under /api/v1.
// All finance routes are protected and require a valid JWT.
func RegisterRoutes(
	v1 *gin.RouterGroup,
	jwtService *utils.JWTService,
	categoryCtrl *controllers.CategoryController,
	incomeCtrl *controllers.IncomeSourceController,
	expenseCtrl *controllers.ExpenseController,
	debtCtrl *controllers.DebtController,
	accountCtrl *controllers.AccountController,
) {
	// All finance routes require authentication
	authMiddleware := middlewares.AuthMiddleware(jwtService)

	// Categories
	categories := v1.Group("/categories")
	categories.Use(authMiddleware)
	{
		categories.GET("", categoryCtrl.ListCategories)
		categories.GET("/:id", categoryCtrl.GetCategory)
		categories.POST("", categoryCtrl.CreateCategory)
		categories.PUT("/:id", categoryCtrl.UpdateCategory)
		categories.DELETE("/:id", categoryCtrl.DeleteCategory)
	}

	// Income Sources
	incomeSources := v1.Group("/income-sources")
	incomeSources.Use(authMiddleware)
	{
		incomeSources.GET("", incomeCtrl.ListIncomeSources)
		incomeSources.GET("/:id", incomeCtrl.GetIncomeSource)
		incomeSources.POST("", incomeCtrl.CreateIncomeSource)
		incomeSources.PUT("/:id", incomeCtrl.UpdateIncomeSource)
		incomeSources.PATCH("/:id/deactivate", incomeCtrl.DeactivateIncomeSource)
		incomeSources.DELETE("/:id", incomeCtrl.DeleteIncomeSource)
	}

	// Expenses
	expenses := v1.Group("/expenses")
	expenses.Use(authMiddleware)
	{
		expenses.GET("", expenseCtrl.ListExpenses)
		expenses.GET("/:id", expenseCtrl.GetExpense)
		expenses.POST("", expenseCtrl.CreateExpense)
		expenses.PUT("/:id", expenseCtrl.UpdateExpense)
		expenses.DELETE("/:id", expenseCtrl.DeleteExpense)
	}

	// Debts
	debts := v1.Group("/debts")
	debts.Use(authMiddleware)
	{
		debts.GET("", debtCtrl.ListDebts)
		debts.GET("/:id", debtCtrl.GetDebt)
		debts.POST("", debtCtrl.CreateDebt)
		debts.PUT("/:id", debtCtrl.UpdateDebt)
		debts.POST("/:id/payments", debtCtrl.RecordPayment)
		debts.DELETE("/:id", debtCtrl.DeleteDebt)
	}

	// Accounts
	accounts := v1.Group("/accounts")
	accounts.Use(authMiddleware)
	{
		accounts.GET("", accountCtrl.ListAccounts)
		accounts.GET("/:id", accountCtrl.GetAccount)
		accounts.POST("", accountCtrl.CreateAccount)
		accounts.PUT("/:id", accountCtrl.UpdateAccount)
		accounts.PATCH("/:id/balance", accountCtrl.UpdateBalance)
		accounts.PATCH("/:id/deactivate", accountCtrl.DeactivateAccount)
		accounts.DELETE("/:id", accountCtrl.DeleteAccount)
	}
}
