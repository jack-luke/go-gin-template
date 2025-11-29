package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middlewre that handles any errors on the context. If there
// is one, it returns the final error and status code as JSON to the client.
// See: https://gin-gonic.com/en/docs/examples/error-handling-middleware/
func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": c.Writer.Status(),
			"error":  c.Errors.Last().Err.Error(),
		})
	}
}
