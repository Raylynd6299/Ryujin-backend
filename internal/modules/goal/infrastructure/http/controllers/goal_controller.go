package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/modules/goal/application/dto"
	appServices "github.com/Raylynd6299/ryujin/internal/modules/goal/application/services"
	goalErrors "github.com/Raylynd6299/ryujin/internal/modules/goal/domain/errors"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"
)

// GoalController handles purchase goal HTTP endpoints
type GoalController struct {
	goalService *appServices.GoalService
}

func NewGoalController(goalService *appServices.GoalService) *GoalController {
	return &GoalController{goalService: goalService}
}

// ---- Goal endpoints ----

// ListGoals godoc
// GET /api/v1/goals?page=1&per_page=20
func (ctrl *GoalController) ListGoals(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := ctrl.goalService.ListGoals(c.Request.Context(), userID, page, perPage)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// GetGoal godoc
// GET /api/v1/goals/:id
func (ctrl *GoalController) GetGoal(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	goal, err := ctrl.goalService.GetGoal(c.Request.Context(), id, userID)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, goal, "")
}

// CreateGoal godoc
// POST /api/v1/goals
func (ctrl *GoalController) CreateGoal(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	goal, err := ctrl.goalService.CreateGoal(c.Request.Context(), userID, req)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, goal, "goal created successfully")
}

// UpdateGoal godoc
// PUT /api/v1/goals/:id
func (ctrl *GoalController) UpdateGoal(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	goal, err := ctrl.goalService.UpdateGoal(c.Request.Context(), id, userID, req)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, goal, "goal updated successfully")
}

// DeleteGoal godoc
// DELETE /api/v1/goals/:id
func (ctrl *GoalController) DeleteGoal(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.goalService.DeleteGoal(c.Request.Context(), id, userID); err != nil {
		handleGoalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ---- Contribution endpoints ----

// ListContributions godoc
// GET /api/v1/goals/:id/contributions
func (ctrl *GoalController) ListContributions(c *gin.Context) {
	goalID := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	result, err := ctrl.goalService.ListContributions(c.Request.Context(), goalID, userID)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, result, "")
}

// AddContribution godoc
// POST /api/v1/goals/:id/contributions
func (ctrl *GoalController) AddContribution(c *gin.Context) {
	goalID := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateContributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	contribution, err := ctrl.goalService.AddContribution(c.Request.Context(), goalID, userID, req)
	if err != nil {
		handleGoalError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, contribution, "contribution added successfully")
}

// DeleteContribution godoc
// DELETE /api/v1/goals/:id/contributions/:cid
func (ctrl *GoalController) DeleteContribution(c *gin.Context) {
	goalID := c.Param("id")
	contributionID := c.Param("cid")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.goalService.DeleteContribution(c.Request.Context(), goalID, contributionID, userID); err != nil {
		handleGoalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ---- error handler ----

func handleGoalError(c *gin.Context, err error) {
	ge, ok := err.(*goalErrors.GoalError)
	if !ok {
		sharedHTTP.InternalServerErrorResponse(c, "an unexpected error occurred")
		return
	}

	switch ge.Code {
	case goalErrors.ErrCodeGoalNotFound, goalErrors.ErrCodeContributionNotFound:
		sharedHTTP.NotFoundResponse(c, ge.Message)
	case goalErrors.ErrCodeGoalUnauthorized:
		sharedHTTP.ForbiddenResponse(c, ge.Message)
	default:
		sharedHTTP.BadRequestResponse(c, ge.Message, nil)
	}
}
