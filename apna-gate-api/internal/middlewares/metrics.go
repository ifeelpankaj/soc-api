package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "society",
			Subsystem: "api",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests handled by the API.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "society",
			Subsystem: "api",
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "society",
			Subsystem: "api",
			Name:      "http_requests_in_flight",
			Help:      "Current number of HTTP requests being served.",
		},
	)

	httpRequestSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "society",
			Subsystem: "api",
			Name:      "http_request_size_bytes",
			Help:      "HTTP request body size in bytes.",
			Buckets:   []float64{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000},
		},
		[]string{"method", "route", "status"},
	)

	httpResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "society",
			Subsystem: "api",
			Name:      "http_response_size_bytes",
			Help:      "HTTP response body size in bytes.",
			Buckets:   []float64{100, 1_000, 10_000, 100_000, 1_000_000, 10_000_000},
		},
		[]string{"method", "route", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
		httpRequestSizeBytes,
		httpResponseSizeBytes,
	)
}

// Metrics records HTTP request metrics using Gin's resolved route pattern.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipMetrics(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		status := strconv.Itoa(c.Writer.Status())
		labels := prometheus.Labels{
			"method": c.Request.Method,
			"route":  route,
			"status": status,
		}

		httpRequestsTotal.With(labels).Inc()
		httpRequestDuration.With(labels).Observe(time.Since(start).Seconds())

		if c.Request.ContentLength > 0 {
			httpRequestSizeBytes.With(labels).Observe(float64(c.Request.ContentLength))
		}
		if responseSize := c.Writer.Size(); responseSize >= 0 {
			httpResponseSizeBytes.With(labels).Observe(float64(responseSize))
		}
	}
}

func shouldSkipMetrics(path string) bool {
	switch path {
	case "/metrics", "/health", "/healthz", "/ping", "/api/health/live", "/api/health/ready":
		return true
	default:
		return false
	}
}
