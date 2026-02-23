package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/application/dto"
	appServices "github.com/Raylynd6299/ryujin/internal/modules/finance/application/services"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"
)

// AccountController handles financial account HTTP endpoints
type AccountController struct {
	accountService *appServices.AccountService
}

func NewAccountController(accountService *appServices.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

// ListAccounts godoc
// GET /api/v1/accounts
func (ctrl *AccountController) ListAccounts(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	accounts, err := ctrl.accountService.ListAccounts(c.Request.Context(), userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, accounts, "")
}

// GetAccount godoc
// GET /api/v1/accounts/:id
func (ctrl *AccountController) GetAccount(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	account, err := ctrl.accountService.GetAccount(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, account, "")
}

// CreateAccount godoc
// POST /api/v1/accounts
func (ctrl *AccountController) CreateAccount(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	account, err := ctrl.accountService.CreateAccount(c.Request.Context(), userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, account, "account created successfully")
}

// UpdateAccount godoc
// PUT /api/v1/accounts/:id
func (ctrl *AccountController) UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	account, err := ctrl.accountService.UpdateAccount(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, account, "account updated successfully")
}

// UpdateBalance godoc
// PATCH /api/v1/accounts/:id/balance
func (ctrl *AccountController) UpdateBalance(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	account, err := ctrl.accountService.UpdateBalance(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, account, "balance updated successfully")
}

// DeactivateAccount godoc
// PATCH /api/v1/accounts/:id/deactivate
func (ctrl *AccountController) DeactivateAccount(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	account, err := ctrl.accountService.DeactivateAccount(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, account, "account deactivated successfully")
}

// DeleteAccount godoc
// DELETE /api/v1/accounts/:id
func (ctrl *AccountController) DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.accountService.DeleteAccount(c.Request.Context(), id, userID); err != nil {
		handleFinanceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
