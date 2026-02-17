package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharedErrors "github.com/Raylynd6299/ryujin/internal/shared/domain/errors"
)

// HandleError handles domain errors and converts them to appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *sharedErrors.NotFoundError:
		NotFoundResponse(c, e.Error())
	case *sharedErrors.ValidationError:
		BadRequestResponse(c, "Validation failed", []string{e.Error()})
	case *sharedErrors.ValidationErrors:
		errors := make([]string, 0, len(e.Errors()))
		for _, validationErr := range e.Errors() {
			errors = append(errors, validationErr.Error())
		}
		BadRequestResponse(c, "Validation failed", errors)
	case *sharedErrors.UnauthorizedError:
		UnauthorizedResponse(c, e.Error())
	case *sharedErrors.DomainError:
		BadRequestResponse(c, e.Error(), nil)
	default:
		// Log the error for debugging (in production, use proper logging)
		InternalServerErrorResponse(c, "An unexpected error occurred")
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (in production, use proper logging)
				c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
					Success: false,
					Message: "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
