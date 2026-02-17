package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http/middlewares"
)

// SetupRouter configures all routes and middlewares for the application
func SetupRouter(deps *AppDependencies) *gin.Engine {
	engine := deps.Engine

	// Apply global middlewares
	engine.Use(gin.Recovery()) // Recover from panics
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
	// TODO: Uncomment when modules are ready
	// v1 := engine.Group("/api/v1")
	// {
	// 	// TODO: Register user module routes
	// 	// userRoutes := v1.Group("/users")
	// 	// user.RegisterRoutes(userRoutes, deps.DB)

	// 	// TODO: Register finance module routes
	// 	// financeRoutes := v1.Group("/finance")
	// 	// finance.RegisterRoutes(financeRoutes, deps.DB)

	// 	// TODO: Register investment module routes
	// 	// investmentRoutes := v1.Group("/investments")
	// 	// investment.RegisterRoutes(investmentRoutes, deps.DB)

	// 	// TODO: Register goal module routes
	// 	// goalRoutes := v1.Group("/goals")
	// 	// goal.RegisterRoutes(goalRoutes, deps.DB)

	// 	// TODO: Register dashboard module routes
	// 	// dashboardRoutes := v1.Group("/dashboard")
	// 	// dashboard.RegisterRoutes(dashboardRoutes, deps.DB)
	// }

	return engine
}
