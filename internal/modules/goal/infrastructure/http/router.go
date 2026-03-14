package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/http/controllers"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/middlewares"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
)

// RegisterRoutes registers all goal module routes under /api/v1.
// All goal routes are protected and require a valid JWT.
func RegisterRoutes(
	v1 *gin.RouterGroup,
	jwtService *utils.JWTService,
	goalCtrl *controllers.GoalController,
) {
	authMiddleware := middlewares.AuthMiddleware(jwtService)

	goals := v1.Group("/goals")
	goals.Use(authMiddleware)
	{
		// Goal CRUD
		goals.GET("", goalCtrl.ListGoals)
		goals.GET("/:id", goalCtrl.GetGoal)
		goals.POST("", goalCtrl.CreateGoal)
		goals.PUT("/:id", goalCtrl.UpdateGoal)
		goals.DELETE("/:id", goalCtrl.DeleteGoal)

		// Contributions (nested under goals)
		goals.GET("/:id/contributions", goalCtrl.ListContributions)
		goals.POST("/:id/contributions", goalCtrl.AddContribution)
		goals.DELETE("/:id/contributions/:cid", goalCtrl.DeleteContribution)
	}
}
