package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raylynd6299/ryujin/internal/modules/finance/application/dto"
	appServices "github.com/Raylynd6299/ryujin/internal/modules/finance/application/services"
	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
	userMiddlewares "github.com/Raylynd6299/ryujin/internal/modules/user/infrastructure/http/middlewares"
	sharedHTTP "github.com/Raylynd6299/ryujin/internal/shared/infrastructure/http"
)

// CategoryController handles category HTTP endpoints
type CategoryController struct {
	categoryService *appServices.CategoryService
}

func NewCategoryController(categoryService *appServices.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// ListCategories godoc
// GET /api/v1/categories
func (ctrl *CategoryController) ListCategories(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	categories, err := ctrl.categoryService.ListCategories(c.Request.Context(), userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, categories, "")
}

// GetCategory godoc
// GET /api/v1/categories/:id
func (ctrl *CategoryController) GetCategory(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	category, err := ctrl.categoryService.GetCategory(c.Request.Context(), id, userID)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, category, "")
}

// CreateCategory godoc
// POST /api/v1/categories
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	category, err := ctrl.categoryService.CreateCategory(c.Request.Context(), userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.CreatedResponse(c, category, "category created successfully")
}

// UpdateCategory godoc
// PUT /api/v1/categories/:id
func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHTTP.BadRequestResponse(c, "invalid request body", []string{err.Error()})
		return
	}

	category, err := ctrl.categoryService.UpdateCategory(c.Request.Context(), id, userID, req)
	if err != nil {
		handleFinanceError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, category, "category updated successfully")
}

// DeleteCategory godoc
// DELETE /api/v1/categories/:id
func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	userID := userMiddlewares.GetUserIDFromContext(c)

	if err := ctrl.categoryService.DeleteCategory(c.Request.Context(), id, userID); err != nil {
		handleFinanceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleFinanceError maps finance domain errors to HTTP responses
func handleFinanceError(c *gin.Context, err error) {
	fe, ok := err.(*financeErrors.FinanceError)
	if !ok {
		sharedHTTP.InternalServerErrorResponse(c, "an unexpected error occurred")
		return
	}

	switch fe.Code {
	case financeErrors.ErrCodeCategoryNotFound,
		financeErrors.ErrCodeIncomeSourceNotFound,
		financeErrors.ErrCodeExpenseNotFound,
		financeErrors.ErrCodeDebtNotFound,
		financeErrors.ErrCodeAccountNotFound:
		sharedHTTP.NotFoundResponse(c, fe.Message)
	case financeErrors.ErrCodeUnauthorized:
		sharedHTTP.ForbiddenResponse(c, fe.Message)
	default:
		sharedHTTP.BadRequestResponse(c, fe.Message, nil)
	}
}
