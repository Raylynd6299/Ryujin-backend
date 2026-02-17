package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log request details
		log.Printf("[%s] %s %s | Status: %d | Latency: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}
