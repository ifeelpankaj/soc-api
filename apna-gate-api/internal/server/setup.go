package server

import (
	"context"
	"go-server/internal/app"
	"go-server/internal/config"
	middleware "go-server/internal/middlewares"
	"go-server/internal/middlewares/guards"
	"go-server/internal/routes"

	"go-server/pkg/database"
	"go-server/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "go-server/docs"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupMiddleware initializes all middleware for the Gin engine
func SetupMiddleware(r *gin.Engine, cfg *config.Config) {
	// Request ID (must be first for logging)
	r.Use(middleware.RequestID())

	if cfg.EnableMetrics {
		r.Use(middleware.Metrics())
	}

	// Logger middleware
	loggerCfg := middleware.DefaultLoggerConfig()
	loggerCfg.SkipHealthCheck = false
	loggerCfg.SkipPaths = nil
	r.Use(middleware.LoggerWithConfig(loggerCfg))

	// Recovery middleware
	recoveryCfg := middleware.DefaultRecoveryConfig()
	recoveryCfg.EnableStackTrace = true
	recoveryCfg.EnableRequestDump = cfg.IsDevelopment()
	r.Use(middleware.RecoveryWithConfig(recoveryCfg))

	// CORS + Rate limiting
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.RateLimit(cfg))

	// Trusted proxies
	if len(cfg.TrustedProxies) > 0 {
		if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
			logger.Error("failed to set trusted proxies", zap.Error(err))
		}
	}
}

// setupRoutes configures all application routes
func SetupRoutes(r *gin.Engine, deps *app.Dependencies, cfg *config.Config) {
	logger.Info("🛣️  Setting up routes...")

	if cfg.EnableMetrics {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
	routes.SetupInternalJobRoutes(r.Group("/api/internal"), deps.V1.JobWebhook, cfg.JobWebhookSecret)

	// API version 1
	apiV1 := r.Group("/api/v1")
	{

		authGuards := guards.New(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer, deps.Society, guards.WithSubscriptionService(deps.Subscription), guards.WithSessionUsers(deps.Users))
		routes.SetupAuthRoutesV1(apiV1, deps.V1.Auth, authGuards)
		routes.SetupImageRoutesV1(apiV1, deps.V1.Image, authGuards)
		routes.SetupHubRoutesV1(apiV1, deps.V1.Hub, authGuards)
		routes.SetupMeRoutesV1(apiV1, deps.V1.Notification, authGuards)
		routes.SetupPublicRoutesV1(apiV1, deps.V1, authGuards)
		routes.SetupBootstrapRoutesV1(apiV1, deps.V1.Bootstrap, authGuards)
		routes.SetupResidentRoutesV1(apiV1, deps.V1, authGuards)
		routes.SetupAdminRoutesV1(apiV1, deps.V1, authGuards)
		routes.SetupDeveloperRoutesV1(apiV1, deps.V1, authGuards)
		if cfg.IsDevelopment() {
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
	apiV2 := r.Group("/api/v2")
	_ = apiV2
	{
		// Future v2 routes go here

	}

	logger.Info("✓ Routes configured successfully")
}

func SetupHealthCheck(r *gin.Engine, cfg *config.Config, db *database.Database) {
	r.GET("/health", healthCheckHandler(cfg))
	r.GET("/api/health/ready", readinessCheckHandler(db))
	r.GET("/api/health/live", livenessCheckHandler())

	logger.Info("✓ Health check endpoints configured")
}

// HealthCheck godoc
// @Summary Health check
// @Description Checks if the API server is running.
// @Tags Health
// @Produce json
// @Success 200 {object} models.HealthCheckResponseDoc "API server is healthy"
// @Router /health [get]
func healthCheckHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"app":         cfg.AppName,
			"version":     cfg.Version,
			"environment": cfg.Environment,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Checks if the API server and database are ready.
// @Tags Health
// @Produce json
// @Success 200 {object} models.ReadinessResponseDoc "API and database are ready"
// @Failure 503 {object} models.HealthErrorResponseDoc "Database connection failed"
// @Router /health/ready [get]
func readinessCheckHandler(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			logger.Error("Readiness check failed: database ping error", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  "database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"database":  "connected",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// LivenessCheck godoc
// @Summary Liveness check
// @Description Checks if the API process is alive.
// @Tags Health
// @Produce json
// @Success 200 {object} models.LivenessResponseDoc "API process is alive"
// @Router /health/live [get]
func livenessCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "alive",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
