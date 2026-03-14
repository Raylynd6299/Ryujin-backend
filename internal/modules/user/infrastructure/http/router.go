package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/controllers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/middlewares"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// RegisterRoutes registers all user module routes into the provided router group.
func RegisterRoutes(v1 *gin.RouterGroup, authCtrl *controllers.AuthController, profileCtrl *controllers.ProfileController, jwtService *utils.JWTService) {
	fmt.Println("Registering user module routes...")
	// Public auth routes — no JWT required
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
		auth.POST("/refresh", authCtrl.RefreshToken)
		auth.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"service": "ryujin-backend",
			})
		})
	}

	// Protected user routes — JWT required
	users := v1.Group("/users")
	users.Use(middlewares.AuthMiddleware(jwtService))
	{
		users.GET("/me", profileCtrl.GetMe)
		users.PUT("/me", profileCtrl.UpdateMe)
		users.PATCH("/me/currencies", profileCtrl.UpdateCurrencies)
		users.PATCH("/me/password", profileCtrl.ChangePassword)
	}
}
