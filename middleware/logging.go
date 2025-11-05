// Gin middleware handlers.
package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Slogger is a middleware that implements structured request logging via slog.
func Slogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		// 2XX 3XX codes emit INFO
		level := slog.LevelInfo

		switch {
		// 5XX codes emit ERROR
		case status >= 500:
			level = slog.LevelError

		// 4XX codes emit WARN
		case status >= 400:
			level = slog.LevelWarn
		}

		logger.LogAttrs(
			c.Request.Context(),
			level,
			fmt.Sprintf("HTTP %d (%s)", status, http.StatusText(status)),
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
			slog.Duration("duration", duration),
			slog.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
