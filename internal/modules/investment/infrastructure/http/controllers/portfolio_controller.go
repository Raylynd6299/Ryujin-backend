package controllers

import (
	appServices "github.com/Raylynd6299/ryujin/internal/modules/investment/application/services"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"

	"github.com/gin-gonic/gin"
)

// PortfolioController handles portfolio-level HTTP endpoints
type PortfolioController struct {
	portfolioService *appServices.PortfolioService
}

// NewPortfolioController creates a new PortfolioController
func NewPortfolioController(portfolioService *appServices.PortfolioService) *PortfolioController {
	return &PortfolioController{portfolioService: portfolioService}
}

// GetSummary godoc
// GET /api/v1/portfolio/summary
func (ctrl *PortfolioController) GetSummary(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	summary, err := ctrl.portfolioService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, summary, "")
}

// GetPerformance godoc
// GET /api/v1/portfolio/performance
func (ctrl *PortfolioController) GetPerformance(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	performance, err := ctrl.portfolioService.GetPerformance(c.Request.Context(), userID)
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, performance, "")
}
