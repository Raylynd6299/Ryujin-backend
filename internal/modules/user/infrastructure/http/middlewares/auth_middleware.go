package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

const (
	// ContextUserID is the key used to store the authenticated user ID in the Gin context.
	ContextUserID = "userID"
	// ContextUserEmail is the key used to store the authenticated user email in the Gin context.
	ContextUserEmail = "userEmail"
)

// AuthMiddleware validates the JWT Bearer token and injects user claims into the context.
func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sharedHTTP.UnauthorizedResponse(c, "authorization header is required")
			c.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			sharedHTTP.UnauthorizedResponse(c, "authorization header must be in format: Bearer <token>")
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			sharedHTTP.UnauthorizedResponse(c, "token cannot be empty")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			sharedHTTP.UnauthorizedResponse(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Inject user data into context for downstream handlers
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserEmail, claims.Email)

		c.Next()
	}
}

// GetUserIDFromContext extracts the authenticated user ID from the Gin context.
// Returns empty string if not set (should not happen on protected routes).
func GetUserIDFromContext(c *gin.Context) string {
	userID, _ := c.Get(ContextUserID)
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}
