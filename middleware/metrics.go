package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics simply groups all exported metrics
type metrics struct {
	RequestDuration   *prometheus.HistogramVec
	HTTPRequestsTotal *prometheus.CounterVec
	RequestsInFlight  prometheus.Gauge
}

// NewMetrics creates new instances of Prometheus metrics objects.
// This is called when the metrics middleware is instantiated.
func NewMetrics() *metrics {
	histogramBuckets := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

	httpRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "http_requests_total",
		Help:        "Total number of HTTP requests.",
		ConstLabels: prometheus.Labels{},
	}, []string{"method", "route", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_seconds",
		Help:        "HTTP request duration in seconds.",
		Buckets:     histogramBuckets,
		ConstLabels: prometheus.Labels{},
	}, []string{"method", "route", "status"})

	requestsInFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "http_in_flight_requests",
		Help:        "Number of requests currently being handled by the service.",
		ConstLabels: prometheus.Labels{},
	})

	return &metrics{
		RequestDuration:   requestDuration,
		HTTPRequestsTotal: httpRequestsTotal,
		RequestsInFlight:  requestsInFlight,
	}
}

// RegisterMetrics registers all provided metrics to the specified registerer
func RegisterMetrics(reg prometheus.Registerer, m *metrics) {
	reg.MustRegister(
		m.RequestDuration,
		m.HTTPRequestsTotal,
		m.RequestsInFlight,
	)
}

// PrometheusMetrics is a middleware that records HTTP metrics about requests.
func PrometheusMetrics(m *metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		m.RequestsInFlight.Inc() // mark request as in flight

		c.Next()

		// gather request attributes
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		route := c.FullPath()

		m.HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
		m.RequestDuration.WithLabelValues(method, route, status).Observe(duration)

		m.RequestsInFlight.Dec() // mark request as completed
	}
}
