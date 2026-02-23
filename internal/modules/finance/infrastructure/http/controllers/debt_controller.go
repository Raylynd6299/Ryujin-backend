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

// DebtController handles debt HTTP endpoints
type DebtController struct {
	debtService *appServices.DebtService
}

func NewDebtController(debtService *appServices.DebtService) *DebtController {
	return &DebtController{debtService: debtService}
}

// ListDebts godoc
// GET /api/v1/debts?page=1&per_page=20
func (ctrl *DebtController) ListDebts(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := ctrl.debtService.ListDebts(c.Request.Context(), userID, page, perPage)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetDebt godoc
// GET /api/v1/debts/:id
func (ctrl *DebtController) GetDebt(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	debt, err := ctrl.debtService.GetDebt(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, debt, "")
}

// CreateDebt godoc
// POST /api/v1/debts
func (ctrl *DebtController) CreateDebt(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateDebtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	debt, err := ctrl.debtService.CreateDebt(c.Request.Context(), userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, debt, "debt created successfully")
}

// UpdateDebt godoc
// PUT /api/v1/debts/:id
func (ctrl *DebtController) UpdateDebt(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateDebtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	debt, err := ctrl.debtService.UpdateDebt(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, debt, "debt updated successfully")
}

// RecordPayment godoc
// POST /api/v1/debts/:id/payments
func (ctrl *DebtController) RecordPayment(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	debt, err := ctrl.debtService.RecordPayment(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, debt, "payment recorded successfully")
}

// DeleteDebt godoc
// DELETE /api/v1/debts/:id
func (ctrl *DebtController) DeleteDebt(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.debtService.DeleteDebt(c.Request.Context(), id, userID); err != nil {
		handleFinanceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
