package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/dto"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/services"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
)

// ProfileController handles user profile HTTP endpoints.
type ProfileController struct {
	profileService *services.ProfileService
}

// NewProfileController creates a new ProfileController.
func NewProfileController(profileService *services.ProfileService) *ProfileController {
	return &ProfileController{profileService: profileService}
}

// GetMe godoc
// GET /api/v1/users/me
func (ctrl *ProfileController) GetMe(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)

	profile, err := ctrl.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, profile, "")
}

// UpdateMe godoc
// PUT /api/v1/users/me
func (ctrl *ProfileController) UpdateMe(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	profile, err := ctrl.profileService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, profile, "profile updated successfully")
}

// UpdateCurrencies godoc
// PATCH /api/v1/users/me/currencies
func (ctrl *ProfileController) UpdateCurrencies(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)

	var req dto.UpdateCurrenciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	profile, err := ctrl.profileService.UpdateCurrencies(c.Request.Context(), userID, req)
	if err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, profile, "currencies updated successfully")
}

// ChangePassword godoc
// PATCH /api/v1/users/me/password
func (ctrl *ProfileController) ChangePassword(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	if err := ctrl.profileService.ChangePassword(c.Request.Context(), userID, req); err != nil {
		sharedHTTP.HandleError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, nil, "password changed successfully")
}
