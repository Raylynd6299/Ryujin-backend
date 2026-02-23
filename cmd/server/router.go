package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	financeHTTP "github.com/Raylynd6299/ryujin/internal/modules/finance/infrastructure/http"
	investmentHTTP "github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/http"
	userHTTP "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	"github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http/middlewares"
)

// SetupRouter configures all routes and middlewares for the application
func SetupRouter(deps *AppDependencies) *gin.Engine {
	engine := deps.Engine

	// Apply global middlewares
	engine.Use(gin.Recovery())
	engine.Use(middlewares.LoggerMiddleware())
	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.RateLimitMiddleware())

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ryujin-backend",
		})
	})

	// API v1 routes group
	v1 := engine.Group("/api/v1")
	{
		// User module: /auth/* and /users/*
		userHTTP.RegisterRoutes(v1, deps.AuthController, deps.ProfileController, deps.JWTService)

		// Finance module: /categories/* /income-sources/* /expenses/* /debts/* /accounts/*
		financeHTTP.RegisterRoutes(
			v1,
			deps.JWTService,
			deps.CategoryController,
			deps.IncomeSourceController,
			deps.ExpenseController,
			deps.DebtController,
			deps.AccountController,
		)

		// Investment module: /holdings/* /portfolio/* /stocks/*
		investmentHTTP.RegisterRoutes(
			v1,
			deps.HoldingController,
			deps.PortfolioController,
			deps.StockAnalysisController,
			userMiddlewares.AuthMiddleware(deps.JWTService),
		)

		// TODO: Register goal module routes
		// TODO: Register dashboard module routes
	}

	return engine
}
