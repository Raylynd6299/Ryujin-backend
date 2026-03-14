package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/dto"
	appServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/services"
	userMiddlewares "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
)

// IncomeSourceController handles income source HTTP endpoints
type IncomeSourceController struct {
	incomeService *appServices.IncomeSourceService
}

func NewIncomeSourceController(incomeService *appServices.IncomeSourceService) *IncomeSourceController {
	return &IncomeSourceController{incomeService: incomeService}
}

// ListIncomeSources godoc
// GET /api/v1/income-sources?page=1&per_page=20
func (ctrl *IncomeSourceController) ListIncomeSources(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := ctrl.incomeService.ListIncomeSources(c.Request.Context(), userID, page, perPage)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetIncomeSource godoc
// GET /api/v1/income-sources/:id
func (ctrl *IncomeSourceController) GetIncomeSource(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	income, err := ctrl.incomeService.GetIncomeSource(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, income, "")
}

// CreateIncomeSource godoc
// POST /api/v1/income-sources
func (ctrl *IncomeSourceController) CreateIncomeSource(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateIncomeSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	income, err := ctrl.incomeService.CreateIncomeSource(c.Request.Context(), userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, income, "income source created successfully")
}

// UpdateIncomeSource godoc
// PUT /api/v1/income-sources/:id
func (ctrl *IncomeSourceController) UpdateIncomeSource(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateIncomeSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	income, err := ctrl.incomeService.UpdateIncomeSource(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, income, "income source updated successfully")
}

// DeactivateIncomeSource godoc
// PATCH /api/v1/income-sources/:id/deactivate
func (ctrl *IncomeSourceController) DeactivateIncomeSource(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.DeactivateIncomeSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	income, err := ctrl.incomeService.DeactivateIncomeSource(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, income, "income source deactivated successfully")
}

// DeleteIncomeSource godoc
// DELETE /api/v1/income-sources/:id
func (ctrl *IncomeSourceController) DeleteIncomeSource(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.incomeService.DeleteIncomeSource(c.Request.Context(), id, userID); err != nil {
		handleFinanceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
