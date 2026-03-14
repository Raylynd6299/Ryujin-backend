package middlewares

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/config"
)

// CORSMiddleware handles Cross-Origin Resource Sharing.
//
// Rules enforced:
//   - AllowCredentials=true is incompatible with AllowOrigins=["*"] per HTTP spec.
//     When credentials are enabled we fall back to AllowAllOrigins=false and rely
//     on the explicit origin list.
//   - When credentials are disabled and the list is empty/["*"] we allow all origins
//     via AllowAllOrigins=true (no wildcard string needed).
//   - Preflight responses are cached for 12 hours (MaxAge).
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := splitCSV(config.App.CORS.AllowedOrigins)
	allowedMethods := splitCSV(config.App.CORS.AllowedMethods)
	allowedHeaders := splitCSV(config.App.CORS.AllowedHeaders)
	allowCredentials := config.App.CORS.AllowCredentials

	// Default methods/headers when not configured
	if len(allowedMethods) == 0 {
		allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
	if len(allowedHeaders) == 0 {
		allowedHeaders = []string{"Authorization", "Content-Type"}
	}

	// Detect wildcard-only origin list
	isWildcard := len(allowedOrigins) == 0 ||
		(len(allowedOrigins) == 1 && allowedOrigins[0] == "*")

	corsConfig := cors.Config{
		AllowMethods:     allowedMethods,
		AllowHeaders:     allowedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}

	if isWildcard && !allowCredentials {
		// Safe to open everything — no credentials involved
		corsConfig.AllowAllOrigins = true
	} else {
		// Explicit list required (credentials mode or specific origins)
		corsConfig.AllowOrigins = allowedOrigins
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
