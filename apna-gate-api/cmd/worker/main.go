package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"go-server/internal/app"
	"go-server/internal/config"
	repository "go-server/internal/repositories"
	"go-server/internal/repositories/contracts"
	"go-server/internal/server"
	notificationsvc "go-server/internal/services/notificationSvc"
	"go-server/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	replayID := flag.Int64("replay-push-id", 0, "requeue one dead push delivery")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	if err = server.InitializeLogger(cfg); err != nil {
		return err
	}
	db, err := database.New(ctx, cfg, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	if *replayID > 0 {
		store, ok := repository.NewNotificationRepository(db).(contracts.PipelineRepository)
		if !ok {
			return errors.New("notification replay repository unavailable")
		}
		return store.ReplayPushDelivery(ctx, *replayID)
	}
	deps, err := app.InitializeDependencies(db, cfg)
	if err != nil {
		return err
	}
	defer deps.Shutdown()
	worker, ok := deps.Notification.(interface{ RunPipeline(context.Context) error })
	if !ok {
		return errors.New("notification worker unavailable")
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		if !notificationsvc.PipelineReady() {
			http.Error(w, "worker loop unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := &http.Server{Addr: ":9091", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.ListenAndServe() }()
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(cleanup)
	}()
	err = worker.RunPipeline(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}
