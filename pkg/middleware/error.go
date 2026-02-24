package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-event-registration/internal/event"
)

// ErrorHandlingMiddleware catches errors added to the context and returns them as JSON
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// If status code hasn't been set, default to 500
			if c.Writer.Status() < 400 {
				c.Status(http.StatusInternalServerError)
			}
			c.JSON(-1, event.ErrorResponse{Error: c.Errors.Last().Error()})
		}
	}
}
