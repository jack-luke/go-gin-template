// A simple Gin app that implements all the required boilerplate for a production app.
package main

import (
	"fmt"
	"log/slog"
	"main/controllers"
	"main/middleware"
	"runtime/debug"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// envDefault takes an environment variable key and a default value, returning
// the environment variable if set, else the default.
func envDefault(key string, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback	
}


// setupRouter creates a router with all middleware and routes attached
func setupRouter() (*gin.Engine, error) {
	r := gin.New()

	// By default, trust no proxies
	// See: https://gin-gonic.com/en/docs/deployment/#dont-trust-all-proxies
	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("Error setting trusted proxies: %v", err)
	}

	// Apply middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Slogger())
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.PrometheusMetrics(nil))
	r.Use(middleware.ErrorHandler)

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Kubernetes liveness & readiness probes
	r.GET("/healthz", controllers.Liveness)
	r.GET("/readyz", controllers.Readiness())

	return r, nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Set the server mode, with 'release' as the default
	gin.SetMode(envDefault("GIN_MODE", "release"))

	// Startup log
	buildInfo, _ := debug.ReadBuildInfo()
	slog.Info("Go Gin Server", 
		"go_version", buildInfo.GoVersion,
		"gin_version", gin.Version,
		"mode", gin.Mode(),
	)

	// Create router with all routes & middleware
	router, err := setupRouter()
	if err != nil {
		slog.Error("Router setup error", "error", err)
		return
	}	

	tlsKeyFile, keyExists := os.LookupEnv("GIN_TLS_KEY_FILE")
	tlsCertFile, certExists := os.LookupEnv("GIN_TLS_CERT_FILE")

	port := envDefault("PORT", "8080")

	// Start server with TLS if cert and key paths are set, else run on HTTP
	if certExists && keyExists {

		// Start a non-blocking QUIC listener unless HTTP/3 is disabled
		if os.Getenv("GIN_HTTP3_ENABLED") != "false" {
			slog.Info("Starting listener", 
				"port", port,
				"tls", "enabled", 
				"transport", "quic", 
				"protocols", "http/3", 
			)
			go func() {
				if err := router.RunQUIC("", tlsCertFile, tlsKeyFile); err != nil {
					slog.Error("QUIC listener error", "error", err)
				}
			}()
		}

		// Start a TCP listener with TLS
		slog.Info("Starting listener", 
			"port", port,
			"tls", "enabled", 
			"transport", "tcp", 
			"protocols", "http/1.1,http/2", 
		)
		if err := router.RunTLS("", tlsCertFile, tlsKeyFile); err != nil {
			slog.Error("TCP TLS listener error", "error", err)
		}

	} else {
		// Start TCP listener without TLS
		slog.Info("Starting listener", 
			"port", port,
			"tls", "disabled", 
			"transport", "tcp", 
			"protocols", "http/1.1,http/2", 
		)
		if err := router.Run(); err != nil {
			slog.Error("TCP listener error", "error", err)
		}
	}
}
