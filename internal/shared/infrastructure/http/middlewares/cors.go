package middlewares

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/config"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := splitCSV(config.App.CORS.AllowedOrigins)
	allowedMethods := splitCSV(config.App.CORS.AllowedMethods)
	allowedHeaders := splitCSV(config.App.CORS.AllowedHeaders)

	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     allowedMethods,
		AllowHeaders:     allowedHeaders,
		AllowCredentials: config.App.CORS.AllowCredentials,
	}

	return cors.New(corsConfig)
}

func splitCSV(value string) []string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	parts := strings.Split(trimmed, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		result = append(result, item)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
