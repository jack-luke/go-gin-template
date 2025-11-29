// controllers contains HTTP route handlers.
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Liveliness acts as a Kubernetes liveliness probe on /healthz
func Liveliness(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// Readiness acts as a Kubernetes readiness probe on /readyz
func Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check app readiness here
		// e.g. ensure dependencies are connected & working

		// MongoDB Example:
		// 		err := client.Ping(context.Background(), readpref.Primary())
		// 		if err != nil {
		// 			c.AbortWithError(
		// 				http.StatusServiceUnavailable,
		// 				fmt.Errorf("MongoDB ping failed: %v", err),
		// 			)
		// 			return
		// 		}

		// MQTT Example:
		// 		if !client.IsConnectionOpen() {
		// 			c.AbortWithError(
		// 				http.StatusServiceUnavailable,
		// 				errors.New("MQTT broker not connected"),
		// 			)
		// 			return
		// 		}

		c.String(http.StatusOK, "OK")
	}
}
