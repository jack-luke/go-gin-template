// A simple Gin app that implements all the required boilerplate for a production app.
package main

import (
	"log/slog"
	"main/controllers"
	"main/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// By default, trust no proxies
	// See: https://gin-gonic.com/en/docs/deployment/#dont-trust-all-proxies
	if err := r.SetTrustedProxies(nil); err != nil {
		logger.Error("Error setting trusted proxies", "error", err)
		return
	}

	// Apply middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Slogger(logger))
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.PrometheusMetrics())
	r.Use(middleware.ErrorHandler)

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Kubernetes liveliness & readiness probes
	r.GET("/healthz", controllers.Liveliness)
	r.GET("/readyz", controllers.Readiness())

	logger.Info("Starting HTTP server")
	if err := r.Run(); err != nil {
		logger.Error("Error starting HTTP server", "error", err)
	}
}
