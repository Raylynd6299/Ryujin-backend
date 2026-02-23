package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/application/dto"
	appServices "github.com/Raylynd6299/ryujin/internal/modules/finance/application/services"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"
)

// ExpenseController handles expense HTTP endpoints
type ExpenseController struct {
	expenseService *appServices.ExpenseService
}

func NewExpenseController(expenseService *appServices.ExpenseService) *ExpenseController {
	return &ExpenseController{expenseService: expenseService}
}

// ListExpenses godoc
// GET /api/v1/expenses?page=1&per_page=20
func (ctrl *ExpenseController) ListExpenses(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := ctrl.expenseService.ListExpenses(c.Request.Context(), userID, page, perPage)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetExpense godoc
// GET /api/v1/expenses/:id
func (ctrl *ExpenseController) GetExpense(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	expense, err := ctrl.expenseService.GetExpense(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, expense, "")
}

// CreateExpense godoc
// POST /api/v1/expenses
func (ctrl *ExpenseController) CreateExpense(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	expense, err := ctrl.expenseService.CreateExpense(c.Request.Context(), userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, expense, "expense created successfully")
}

// UpdateExpense godoc
// PUT /api/v1/expenses/:id
func (ctrl *ExpenseController) UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	expense, err := ctrl.expenseService.UpdateExpense(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, expense, "expense updated successfully")
}

// DeleteExpense godoc
// DELETE /api/v1/expenses/:id
func (ctrl *ExpenseController) DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.expenseService.DeleteExpense(c.Request.Context(), id, userID); err != nil {
		handleFinanceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
