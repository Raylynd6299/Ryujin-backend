package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, errors []string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusCreated, data, message)
}

// OkResponse sends a 200 OK response
func OkResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusOK, data, message)
}

// NoContentResponse sends a 204 No Content response
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(c *gin.Context, message string, errors []string) {
	ErrorResponse(c, http.StatusBadRequest, message, errors)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// ConflictResponse sends a 409 Conflict response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message, nil)
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, nil)
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Meta    Pagination  `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedSuccessResponse sends a paginated success response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, page, limit int, totalItems int64, message string) {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Message: message,
		Meta: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	})
}
