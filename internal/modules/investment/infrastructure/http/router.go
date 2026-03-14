package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/http/controllers"
)

// RegisterRoutes registers all investment module routes under /api/v1.
// All investment routes are protected and require a valid JWT (passed via authMiddleware).
func RegisterRoutes(
	router *gin.RouterGroup,
	holdingCtrl *controllers.HoldingController,
	portfolioCtrl *controllers.PortfolioController,
	stockAnalysisCtrl *controllers.StockAnalysisController,
	authMiddleware gin.HandlerFunc,
) {
	investments := router.Group("/").Use(authMiddleware)
	{
		investments.GET("/holdings", holdingCtrl.ListHoldings)
		investments.POST("/holdings", holdingCtrl.CreateHolding)
		investments.GET("/holdings/:id", holdingCtrl.GetHolding)
		investments.PUT("/holdings/:id", holdingCtrl.UpdateHolding)
		investments.DELETE("/holdings/:id", holdingCtrl.DeleteHolding)
		investments.POST("/holdings/:id/refresh-price", holdingCtrl.RefreshPrice)

		investments.GET("/portfolio/summary", portfolioCtrl.GetSummary)
		investments.GET("/portfolio/performance", portfolioCtrl.GetPerformance)

		// Stock analysis
		investments.GET("/stocks", stockAnalysisCtrl.ListStockQuotes)
		investments.GET("/stocks/:symbol/quote", stockAnalysisCtrl.GetStockQuote)
		investments.GET("/stocks/:symbol/history", stockAnalysisCtrl.GetStockPriceHistory)
	}
}
