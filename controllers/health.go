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
		// 	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// 		c.String(http.StatusServiceUnavailable, "Database not connected")
		// 		return
		// 	}

		// MQTT Example:
		// 	if !client.IsConnectionOpen() {
		//		c.String(http.StatusServiceUnavailable, "MQTT broker not connected")
		// 		return
		// 	}

		c.String(http.StatusOK, "OK")
	}
}
