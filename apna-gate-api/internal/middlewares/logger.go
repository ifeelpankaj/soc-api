package middleware

import (
	"github.com/google/uuid"
	"go-server/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerConfig defines configuration for the logger middleware
type LoggerConfig struct {
	// SkipPaths is a list of paths to skip logging (e.g., health checks)
	SkipPaths []string
	// SkipHealthCheck skips /health and /healthz endpoints
	SkipHealthCheck bool
	// EnableBody logs request/response bodies (use cautiously in prod)
	EnableBody bool
}

// DefaultLoggerConfig returns sensible defaults
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipPaths:       []string{},
		SkipHealthCheck: true,
		EnableBody:      false,
	}
}

// Logger returns a gin.HandlerFunc middleware for request logging
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig())
}

// LoggerWithConfig returns a gin.HandlerFunc with custom config
func LoggerWithConfig(cfg LoggerConfig) gin.HandlerFunc {
	// Build skip paths map for O(1) lookup
	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}
	if cfg.SkipHealthCheck {
		skipPaths["/health"] = true
		skipPaths["/healthz"] = true
		skipPaths["/ping"] = true
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		if isImagePath(path) {
			c.Header("Cache-Control", "no-store")
		}

		// Skip logging for configured paths
		if skipPaths[path] {
			c.Next()
			return
		}

		// Process request
		c.Next()

		// Calculate metrics
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		message := c.GetString("response_message")
		source := c.GetString("response_message_source")
		if source == "" {
			message = http.StatusText(statusCode)
			source = "status_fallback"
		}
		// Prepare log fields
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("event", "request_completed"),
			zap.String("route", route),
			zap.String("response_message", logger.Sanitize(message)),
			zap.String("response_message_source", source),
			zap.String("error_code", c.GetString("error_code")),
			zap.String("internal_error", c.GetString("internal_error")),
			zap.Float64("duration_ms", float64(latency)/float64(time.Millisecond)),
			zap.Bool("panic_recovered", c.GetBool("panic_recovered")),
			zap.Bool("connection_failure", c.GetBool("connection_failure")),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add request ID if exists (common pattern)
		if requestID := c.GetString("request_id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Add user ID if exists (for authenticated requests)
		if userID, ok := GetUserIDFromContext(c); ok {
			fields = append(fields, zap.Int64("user_id", userID))
		}

		// Add content length
		if c.Request.ContentLength > 0 {
			fields = append(fields, zap.Int64("request_size", c.Request.ContentLength))
		}
		fields = append(fields, zap.Int("response_size", c.Writer.Size()))

		// Log any errors that occurred during request processing
		if len(c.Errors) > 0 {
			// Log each error
			for _, ginErr := range c.Errors {
				errFields := make([]zap.Field, 0, len(fields))
				for _, field := range fields {
					if field.Key != "event" && field.Key != "internal_error" {
						errFields = append(errFields, field)
					}
				}
				errFields = append(errFields, zap.String("event", "request_error"), zap.String("internal_error", logger.Sanitize(ginErr.Err.Error())))
				logger.Error("Request error occurred", errFields...)
			}
		}

		// Determine log level and message based on status code
		msg := "Request completed"
		switch {
		case statusCode >= 500:
			logger.Error(msg, fields...)
		case statusCode >= 400:
			logger.Warn(msg, fields...)
		case statusCode >= 300:
			logger.Info(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}

func isImagePath(path string) bool {
	return strings.HasSuffix(path, "/auth/profile/avatar") || (strings.Contains(path, "/visitor-entries/") && strings.HasSuffix(path, "/photo"))
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID exists in header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" || len(requestID) > 128 || logger.Sanitize(requestID) != requestID {
			// Generate a collision-resistant correlation ID for missing/unsafe input.
			requestID = generateRequestID()
		}

		// Set request ID in context and response header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID creates a UUID without exposing time or user information.
func generateRequestID() string { return uuid.NewString() }
