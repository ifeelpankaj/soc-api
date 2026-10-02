package main

import (
	"context"
	"fmt"
	"go-server/internal/app"
	"go-server/internal/config"
	"go-server/internal/server"
	"go-server/pkg/logger"
	"go-server/pkg/utils"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Society Management API
// @version 1.0
// @description Backend API for Society Management System. Development base URL: http://localhost:8000/api. Production base URL: http://localhost:8080/api.
// @termsOfService http://localhost:8000/terms
// @contact.name API Support
// @contact.email support@example.com
// @license.name Proprietary
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey AccessToken
// @description Access token stored as an HTTP-only cookie.
// @in cookie
// @name access_token
// @securityDefinitions.apikey RefreshToken
// @description Refresh token stored as an HTTP-only cookie scoped to /api/v1/auth/refresh.
// @in cookie
// @name refresh_token
// @securityDefinitions.apikey OnboardingToken
// @in cookie
// @name onboarding_token

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}

}

func run() error {
	// ========================================
	// 1. Configuration Loading
	// ========================================

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	// ========================================
	// 2. Logger Initialization
	// ========================================

	if err := server.InitializeLogger(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
		}
	}()

	// Log application startup information

	server.LogStartupInfo(cfg)

	// ========================================
	// 3. Configure Response Utilities
	// ========================================
	utils.SetResponseConfig(utils.ResponseConfig{
		EnableDetailedErrors: cfg.IsDevelopment(),
		EnableErrorLogging:   true,
		IncludeRequestID:     true,
	})
	// ========================================
	// 4. Gin Mode Configuration
	// ========================================
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)

	} else {
		gin.SetMode(gin.DebugMode)
	}
	// ========================================
	// 5. Database Connection & Migration
	// ========================================
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := server.InitializeDatabase(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		db.Close()
	}()
	// Run migrations
	if err := server.RunMigrations(ctx, db); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	// ========================================
	// 6. Dependency Injection & Initialization
	// ========================================
	deps, err := app.InitializeDependencies(db, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}
	defer deps.Shutdown()
	// ========================================
	// 7. HTTP Router Setup
	// ========================================
	router := gin.New()

	// Setup middleware
	server.SetupMiddleware(router, cfg)

	// Setup routes
	server.SetupRoutes(router, deps, cfg)

	if cfg.EnableHealthz {
		server.SetupHealthCheck(router, cfg, db)
	}
	// ========================================
	// 8. HTTP Server Configuration
	// ========================================
	srv := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        router,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	listener, err := bindAndStartInternalJobs(srv.Addr, deps)
	if err != nil {
		return err
	}
	// ========================================
	// 9. Graceful Shutdown Setup
	// ========================================
	return server.StartServerWithGracefulShutdown(srv, listener, cfg)
}

type internalJobStarter interface {
	StartInternalJobs() bool
}

func bindAndStartInternalJobs(address string, starter internalJobStarter) (net.Listener, error) {
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(context.Background(), "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("server failed to bind: %w", err)
	}
	if !starter.StartInternalJobs() {
		_ = listener.Close()
		return nil, fmt.Errorf("failed to start internal expiry job")
	}
	return listener, nil
}
