package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/services"
	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
)

// AuthController handles authentication HTTP endpoints.
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// POST /api/v1/auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	response, err := ctrl.authService.Register(c.Request.Context(), req)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, response, "user registered successfully")
}

// Login godoc
// POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	response, err := ctrl.authService.Login(c.Request.Context(), req)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.SuccessResponse(c, http.StatusOK, response, "login successful")
}

// RefreshToken godoc
// POST /api/v1/auth/refresh
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	tokens, err := ctrl.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, tokens, "token refreshed successfully")
}
