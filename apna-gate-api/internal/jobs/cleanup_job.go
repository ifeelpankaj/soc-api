package jobs

import (
	"context"
	"errors"
	"fmt"
	"go-server/pkg/logger"
	"time"

	"go-server/internal/models"

	"go.uber.org/zap"
)

const cleanupJobLock = "apna-gate/jobs/cleanup"

type cleanupStore interface {
	DeleteOrphanVisitorsBatch(context.Context, time.Time, int32) (int64, error)
	ListVisitorImageDeletions(context.Context, int64, int32) ([]models.PendingImageDeletion, error)
	CompleteVisitorImageDeletion(context.Context, int64) error
	WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error)
	DeleteNotificationsBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
	DeleteVerificationsBatch(ctx context.Context, batchSize int32) (int64, error)
	ListVisitorEntryCleanupWindows(ctx context.Context, timezone string, cutoff time.Time) ([]models.VisitorReportCleanupWindow, error)
	DeleteVisitorEntriesBatch(ctx context.Context, window models.VisitorReportCleanupWindow, from, to, cutoff time.Time, batchSize int32) (int64, error)
	DeleteVisitorInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
	DeleteFlatMemberInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
}

type CleanupJobConfig struct {
	HubCleanup   func(context.Context) error
	ImageStorage interface {
		Delete(context.Context, string) error
	}
	Location                  *time.Location
	Timezone                  string
	VisitorEntryRetention     time.Duration
	VisitorInviteRetention    time.Duration
	FlatMemberInviteRetention time.Duration
	NotificationRetention     time.Duration
	BatchSize                 int32
}

type CleanupJob struct {
	repo cleanupStore
	cfg  CleanupJobConfig
	now  func() time.Time
}

func NewCleanupJob(repo cleanupStore, cfg CleanupJobConfig) *CleanupJob {
	if cfg.Location == nil {
		cfg.Location = time.UTC
	}
	if cfg.Timezone == "" {
		cfg.Timezone = cfg.Location.String()
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 500
	}
	return &CleanupJob{repo: repo, cfg: cfg, now: time.Now}
}

func (j *CleanupJob) RunOnce(ctx context.Context) (RunResult, error) {
	locked, err := j.repo.WithAdvisoryLock(ctx, cleanupJobLock, func(lockCtx context.Context) error {
		var operationErrors []error
		if j.cfg.HubCleanup != nil {
			if err := j.cfg.HubCleanup(lockCtx); err != nil {
				operationErrors = append(operationErrors, err)
			}
		}
		now := j.now()
		if err := j.cleanupSimpleBatches(lockCtx, "notifications", func(batchCtx context.Context) (int64, error) {
			return j.repo.DeleteNotificationsBatch(batchCtx, now.Add(-j.cfg.NotificationRetention), j.cfg.BatchSize)
		}); err != nil {
			if isContextCancellation(err) {
				return err
			}
			operationErrors = append(operationErrors, err)
		}
		if err := j.cleanupSimpleBatches(lockCtx, "verifications", func(batchCtx context.Context) (int64, error) {
			return j.repo.DeleteVerificationsBatch(batchCtx, j.cfg.BatchSize)
		}); err != nil {
			if isContextCancellation(err) {
				return err
			}
			operationErrors = append(operationErrors, err)
		}

		entryCutoff := now.Add(-j.cfg.VisitorEntryRetention)
		windows, windowsErr := j.repo.ListVisitorEntryCleanupWindows(lockCtx, j.cfg.Timezone, entryCutoff)
		if windowsErr != nil {
			if isContextCancellation(windowsErr) {
				return windowsErr
			}
			operationErrors = append(operationErrors, fmt.Errorf("list visitor entry cleanup windows: %w", windowsErr))
		} else {
			for _, window := range windows {
				if err := lockCtx.Err(); err != nil {
					return err
				}
				from, to := MonthBounds(window.ReportMonth, j.cfg.Location)
				name := fmt.Sprintf("visitor entries society=%d month=%s", window.SocietyID, window.ReportMonth.Format("2006-01"))
				if err := j.cleanupSimpleBatches(lockCtx, name, func(batchCtx context.Context) (int64, error) {
					return j.repo.DeleteVisitorEntriesBatch(batchCtx, window, from, to, entryCutoff, j.cfg.BatchSize)
				}); err != nil {
					if isContextCancellation(err) {
						return err
					}
					operationErrors = append(operationErrors, err)
				}
			}
		}

		if err := j.cleanupSimpleBatches(lockCtx, "visitor invites", func(batchCtx context.Context) (int64, error) {
			return j.repo.DeleteVisitorInvitesBatch(batchCtx, now.Add(-j.cfg.VisitorInviteRetention), j.cfg.BatchSize)
		}); err != nil {
			if isContextCancellation(err) {
				return err
			}
			operationErrors = append(operationErrors, err)
		}
		if err := j.cleanupSimpleBatches(lockCtx, "flat member invites", func(batchCtx context.Context) (int64, error) {
			return j.repo.DeleteFlatMemberInvitesBatch(batchCtx, now.Add(-j.cfg.FlatMemberInviteRetention), j.cfg.BatchSize)
		}); err != nil {
			if isContextCancellation(err) {
				return err
			}
			operationErrors = append(operationErrors, err)
		}
		if err := j.cleanupSimpleBatches(lockCtx, "orphan visitors", func(batchCtx context.Context) (int64, error) {
			return j.repo.DeleteOrphanVisitorsBatch(batchCtx, entryCutoff, j.cfg.BatchSize)
		}); err != nil {
			if isContextCancellation(err) {
				return err
			}
			operationErrors = append(operationErrors, err)
		}
		if err := j.cleanupImages(lockCtx); err != nil {
			operationErrors = append(operationErrors, err)
		}
		return errors.Join(operationErrors...)
	})
	result := RunResult{LockAcquired: locked}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (j *CleanupJob) cleanupImages(ctx context.Context) error {
	if j.cfg.ImageStorage == nil {
		return nil
	}
	var afterID int64
	var failures int
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		items, err := j.repo.ListVisitorImageDeletions(ctx, afterID, j.cfg.BatchSize)
		if err != nil {
			return fmt.Errorf("list image deletions: %w", err)
		}
		for _, item := range items {
			if err := ctx.Err(); err != nil {
				return err
			}
			afterID = item.ID
			deleteCtx, cancel := context.WithTimeout(ctx, models.ImageCleanupTimeout)
			err := j.cfg.ImageStorage.Delete(deleteCtx, item.FileID)
			cancel()
			if err == nil {
				err = j.repo.CompleteVisitorImageDeletion(ctx, item.ID)
			}
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				failures++
				jobLogger(ctx).Warn("visitor image deletion remains queued", zap.Int64("deletion_id", item.ID))
			}
		}
		if len(items) < int(j.cfg.BatchSize) {
			break
		}
	}
	if failures > 0 {
		return fmt.Errorf("%d visitor image deletions remain queued for retry", failures)
	}
	return nil
}

func (j *CleanupJob) cleanupSimpleBatches(ctx context.Context, name string, deleteBatch func(context.Context) (int64, error)) error {
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		count, err := deleteBatch(ctx)
		if err != nil {
			if !isContextCancellation(err) {
				jobLogger(ctx).Error("cleanup operation failed", zap.String("operation", name), zap.String("internal_error", logger.Sanitize(err.Error())))
			}
			return fmt.Errorf("%s: %w", name, err)
		}
		total += count
		if count < int64(j.cfg.BatchSize) {
			break
		}
	}
	if total > 0 {
		jobLogger(ctx).Info("cleanup operation completed", zap.String("operation", name), zap.Int64("count", total))
	}
	return nil
}
