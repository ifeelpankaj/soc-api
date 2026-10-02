package backup

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go-server/internal/config"
	"go-server/internal/jobs"
	"go-server/pkg/logger"
	"go.uber.org/zap"
)

var ErrDisabled = errors.New("database backups are disabled")

type Scheduler interface {
	StartTask(func() error, func(context.Context)) error
}
type Service struct {
	cfg              config.BackupConfig
	database, worker string
	store            Store
	locker           Locker
	local            Local
	cloud            Cloud
	scheduler        Scheduler
}

func NewService(cfg config.BackupConfig, database string, store Store, locker Locker, local Local, cloud Cloud, scheduler Scheduler) *Service {
	return &Service{cfg: cfg, database: database, worker: uuid.NewString(), store: store, locker: locker, local: local, cloud: cloud, scheduler: scheduler}
}
func (s *Service) Trigger(ctx context.Context, kind string) (*Run, error) {
	if !s.cfg.Enabled {
		return nil, ErrDisabled
	}
	if kind != "daily" && kind != "weekly" {
		return nil, errors.New("invalid backup type")
	}
	r := &Run{ID: uuid.NewString(), Type: kind, Status: "queued", Database: s.database, Deployment: s.cfg.Deployment, Worker: s.worker,
		Details: Details{Checks: Checks{Database: "pending", PageChecksums: "unknown", Amcheck: "not_performed", Archive: "pending"}, Retention: Retention{Status: "pending"}}}
	var response Run
	err := s.scheduler.StartTask(func() error {
		createCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := s.store.Create(createCtx, r); err != nil {
			return err
		}
		response = *r
		return nil
	}, func(ctx context.Context) { s.execute(ctx, r) })
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func (s *Service) Get(ctx context.Context, id string) (*Run, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.store.Get(ctx, id)
}

// StartRecovery reconciles abandoned runs once, after old heartbeats expire.
// It never schedules a backup and never waits for an active worker's lock.
func (s *Service) StartRecovery() error {
	if !s.cfg.Enabled {
		return nil
	}
	return s.scheduler.StartTask(func() error { return nil }, func(parent context.Context) {
		if pause(parent, 125*time.Second) != nil {
			return
		}
		ctx, cancel := context.WithTimeout(parent, 30*time.Second)
		defer cancel()
		lease, locked, err := s.locker.Try(ctx)
		if err != nil {
			logger.Warn("backup startup recovery could not connect to database")
			return
		}
		if !locked {
			return
		}
		defer lease.Close()
		old, err := lease.Recover(ctx, &Run{ID: uuid.NewString(), Database: s.database, Deployment: s.cfg.Deployment})
		if err != nil {
			logger.Warn("backup startup reconciliation failed")
			return
		}
		root, err := filepath.EvalSymlinks(s.cfg.TempDir)
		if err != nil {
			return
		}
		root, err = filepath.Abs(root)
		if err != nil {
			return
		}
		for _, id := range old {
			if removeRunDir(root, id) != nil {
				logger.Warn("abandoned backup temporary cleanup failed", zap.String("run_id", id))
			}
		}
	})
}
func (s *Service) stage(ctx context.Context, r *Run, state string) error {
	r.Status = state
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.store.Save(ctx, r); err != nil {
		return err
	}
	logger.Info("database backup progress", zap.String("job_name", "database-backup-"+r.Type), zap.String("run_id", r.ID), zap.String("status", state))
	return nil
}
func (s *Service) finish(r *Run, status, code, message string) {
	r.Status = status
	r.ErrorCode = code
	r.ErrorMessage = message
	now := time.Now().UTC()
	r.Completed = &now
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.store.Save(ctx, r); err != nil {
		logger.Error("database backup terminal status could not be persisted", zap.String("job_name", "database-backup-"+r.Type), zap.String("run_id", r.ID), zap.String("status", status))
		return
	}
	fields := []zap.Field{zap.String("job_name", "database-backup-"+r.Type), zap.String("run_id", r.ID), zap.String("status", status), zap.String("error_code", code), zap.String("error_message", message), zap.String("retention", r.Retention.Status)}
	if status == "failed" || status == "interrupted" {
		logger.Error("database backup finished", fields...)
	} else {
		logger.Info("database backup finished", fields...)
	}
}
func (s *Service) execute(parent context.Context, r *Run) {
	started := time.Now()
	log := logger.With(zap.String("job_name", "database-backup-"+r.Type), zap.String("run_id", r.ID))
	log.Info("Job started", zap.String("event", "job_started"))
	completeMetric := jobs.BeginRun("database-backup-" + r.Type)
	defer func() {
		outcome := jobs.OutcomeFailure
		switch r.Status {
		case "completed":
			outcome = jobs.OutcomeSuccess
		case "skipped_concurrent":
			outcome = jobs.OutcomeSkipped
		case "interrupted":
			outcome = jobs.OutcomeInterrupted
		}
		completeMetric(outcome)
		fields := []zap.Field{zap.String("event", "job_completed"), zap.String("outcome", outcome),
			zap.Float64("duration_ms", float64(time.Since(started))/float64(time.Millisecond)),
			zap.String("error_code", r.ErrorCode), zap.String("response_message", logger.Sanitize(r.ErrorMessage))}
		if outcome == jobs.OutcomeFailure {
			log.Error("Job completed", fields...)
		} else {
			log.Info("Job completed", fields...)
		}
	}()
	ctx, cancel := context.WithTimeout(parent, s.cfg.Timeout)
	defer cancel()
	var leaseLost atomic.Bool
	var failure error
	defer func() {
		if recover() != nil {
			failure = errors.New("backup worker panic")
		}
		if r.Completed != nil {
			return
		}
		if failure == nil {
			failure = errors.New("backup worker did not finish")
		}
		if parent.Err() != nil || leaseLost.Load() {
			s.finish(r, "interrupted", "worker_interrupted", "Backup interrupted by shutdown or lost database lease")
			return
		}
		code := r.Status + "_failed"
		message := "Backup failed during " + r.Status + "; check server configuration and database diagnostics"
		var diagnostic *diagnosticError
		if errors.As(failure, &diagnostic) {
			code, message = diagnostic.code, diagnostic.message
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			code = "backup_timeout"
			message = "Backup exceeded its configured time limit"
		}
		// Never persist subprocess stderr, provider responses, URLs, or credentials.
		s.finish(r, "failed", code, message)
	}()
	lease, locked, err := s.locker.Try(ctx)
	if err != nil {
		failure = err
		return
	}
	if !locked {
		s.finish(r, "skipped_concurrent", "", "Another database backup owns the maintenance lock")
		return
	}
	defer lease.Close()
	// Every DB lease operation is serialized; only the heartbeat goroutine uses it below.
	if err := lease.Pulse(ctx, r); err != nil {
		leaseLost.Store(true)
		failure = err
		return
	}
	old, err := lease.Recover(ctx, r)
	if err != nil {
		failure = err
		return
	}
	heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
	heartbeatDone := make(chan struct{})
	go func() {
		defer close(heartbeatDone)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeatCtx.Done():
				return
			case <-ticker.C:
				pulseCtx, done := context.WithTimeout(heartbeatCtx, 5*time.Second)
				err := lease.Pulse(pulseCtx, r)
				done()
				if err != nil && heartbeatCtx.Err() == nil {
					leaseLost.Store(true)
					cancel()
					return
				}
			}
		}
	}()
	defer func() { stopHeartbeat(); <-heartbeatDone }()
	now := time.Now().UTC()
	r.Started = &now
	if err := s.stage(ctx, r, "preflight"); err != nil {
		failure = err
		return
	}
	if err := os.MkdirAll(s.cfg.TempDir, 0700); err != nil {
		failure = preflightFailure("temp_directory", "Cannot create BACKUP_TEMP_DIR; check the directory and service account permissions")
		return
	}
	root, err := filepath.EvalSymlinks(s.cfg.TempDir)
	if err != nil {
		failure = preflightFailure("temp_directory", "Cannot resolve BACKUP_TEMP_DIR; check the directory and symlink target")
		return
	}
	root, err = filepath.Abs(root)
	if err != nil {
		failure = preflightFailure("temp_directory", "Cannot resolve the absolute BACKUP_TEMP_DIR path")
		return
	}
	if err := os.Chmod(root, 0700); err != nil { //nolint:gosec // Directories require owner execute permission; no group/other access.
		failure = preflightFailure("temp_permissions", "Cannot secure BACKUP_TEMP_DIR; check ownership and filesystem permissions")
		return
	}
	for _, id := range old {
		if err := removeRunDir(root, id); err != nil {
			failure = preflightFailure("temp_cleanup", "Cannot remove an abandoned backup directory; check ownership and permissions")
			return
		}
	}
	dir := filepath.Join(root, r.ID)
	if err := os.Mkdir(dir, 0700); err != nil {
		failure = preflightFailure("temp_directory", "Cannot create the backup run directory; check BACKUP_TEMP_DIR permissions and available space")
		return
	}
	defer func() {
		if err := removeRunDir(root, r.ID); err != nil {
			logger.Warn("backup temporary cleanup failed", zap.String("job_name", "database-backup-"+r.Type), zap.String("run_id", r.ID))
		}
	}()
	if err := s.local.Preflight(ctx, r, dir); err != nil {
		failure = err
		return
	}
	if err := s.cloud.Preflight(ctx); err != nil {
		failure = err
		return
	}
	if r.Type == "weekly" {
		if err := s.stage(ctx, r, "checking"); err != nil {
			failure = err
			return
		}
		r.Checks.Amcheck = "failed"
		if err := s.local.Check(ctx, dir); err != nil {
			failure = err
			return
		}
		r.Checks.Amcheck = "passed"
	}
	if err := s.stage(ctx, r, "dumping"); err != nil {
		failure = err
		return
	}
	if err := s.local.Dump(ctx, dir); err != nil {
		failure = err
		return
	}
	if err := s.stage(ctx, r, "validating"); err != nil {
		failure = err
		return
	}
	r.Checks.Archive = "failed"
	if err := s.local.Validate(ctx, dir); err != nil {
		failure = err
		return
	}
	r.Checks.Archive = "passed"
	path := filepath.Join(dir, "database.dump")
	if err := hashArchive(ctx, path, r); err != nil {
		failure = err
		return
	}
	if err := s.stage(ctx, r, "uploading"); err != nil {
		failure = err
		return
	}
	r.DriveFileID, err = s.cloud.Upload(ctx, r, path)
	if err != nil {
		failure = err
		return
	}
	if err := s.stage(ctx, r, "verifying"); err != nil {
		failure = err
		return
	}
	r.DriveMD5, err = s.cloud.Verify(ctx, r)
	if err != nil {
		failure = err
		return
	}
	if err := s.stage(ctx, r, "retention"); err != nil {
		failure = err
		return
	}
	r.Retention.Deleted, err = s.cloud.Retain(ctx, r)
	r.Retention.Status = "completed"
	if err != nil {
		r.Retention.Status = "warning"
	}
	s.finish(r, "completed", "", "")
}

// Only a direct UUID child of the resolved dedicated root can ever be removed.
func removeRunDir(root, id string) error {
	u, err := uuid.Parse(id)
	if err != nil || u.String() != id {
		return errors.New("invalid backup temporary directory identity")
	}
	path := filepath.Join(root, id)
	if filepath.Dir(path) != root {
		return errors.New("backup temporary directory escaped root")
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("backup temporary directory is not a regular directory")
	}
	return os.RemoveAll(path)
}
