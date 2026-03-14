package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/application/dto"
	appServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/application/services"
	investErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/errors"
	userMiddlewares "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
)

// HoldingController handles investment holding HTTP endpoints
type HoldingController struct {
	holdingService *appServices.HoldingService
}

// NewHoldingController creates a new HoldingController
func NewHoldingController(holdingService *appServices.HoldingService) *HoldingController {
	return &HoldingController{holdingService: holdingService}
}

// ListHoldings godoc
// GET /api/v1/holdings?page=1&limit=20&sort=created_at&order=desc
func (ctrl *HoldingController) ListHoldings(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	result, err := ctrl.holdingService.ListHoldings(c.Request.Context(), userID, page, limit, sort, order)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetHolding godoc
// GET /api/v1/holdings/:id
func (ctrl *HoldingController) GetHolding(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	holding, err := ctrl.holdingService.GetHolding(c.Request.Context(), id, userID)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, holding, "")
}

// CreateHolding godoc
// POST /api/v1/holdings
func (ctrl *HoldingController) CreateHolding(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	holding, err := ctrl.holdingService.CreateHolding(c.Request.Context(), userID, req)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, holding, "holding created successfully")
}

// UpdateHolding godoc
// PUT /api/v1/holdings/:id
func (ctrl *HoldingController) UpdateHolding(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	holding, err := ctrl.holdingService.UpdateHolding(c.Request.Context(), id, userID, req)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, holding, "holding updated successfully")
}

// DeleteHolding godoc
// DELETE /api/v1/holdings/:id
func (ctrl *HoldingController) DeleteHolding(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.holdingService.DeleteHolding(c.Request.Context(), id, userID); err != nil {
		handleInvestmentError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RefreshPrice godoc
// POST /api/v1/holdings/:id/refresh-price
func (ctrl *HoldingController) RefreshPrice(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	holding, err := ctrl.holdingService.RefreshHoldingPrice(c.Request.Context(), id, userID)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, holding, "price refreshed successfully")
}

// handleInvestmentError maps investment domain errors to HTTP responses
func handleInvestmentError(c *gin.Context, err error) {
	switch err.(type) {
	case *investErrors.HoldingNotFoundError:
		sharedHTTP.NotFoundResponse(c, err.Error())
	case *investErrors.HoldingValidationError:
		sharedHTTP.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
	case *investErrors.HoldingForbiddenError:
		sharedHTTP.ForbiddenResponse(c, err.Error())
	default:
		sharedHTTP.InternalServerErrorResponse(c, "an unexpected error occurred")
	}
}
