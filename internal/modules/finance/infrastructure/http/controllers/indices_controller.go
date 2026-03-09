package controllers

import (
	"github.com/gin-gonic/gin"

	appServices "github.com/Raylynd6299/ryujin/internal/modules/finance/application/services"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"
)

// IndicesController handles financial health index HTTP endpoints
type IndicesController struct {
	indicesService *appServices.IndicesCalculatorService
}

// NewIndicesController creates a new IndicesController
func NewIndicesController(indicesService *appServices.IndicesCalculatorService) *IndicesController {
	return &IndicesController{indicesService: indicesService}
}

// GetIndices godoc
// GET /api/v1/finance/indices
func (ctrl *IndicesController) GetIndices(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	result, err := ctrl.indicesService.CalculateIndices(c.Request.Context(), userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetSummary godoc
// GET /api/v1/finance/summary
func (ctrl *IndicesController) GetSummary(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	result, err := ctrl.indicesService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}
