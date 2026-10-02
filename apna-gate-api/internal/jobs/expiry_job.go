package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-server/pkg/logger"

	"go.uber.org/zap"
)

const expiryJobLock = "apna-gate/jobs/expiry"

type expiryStore interface {
	WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error)
	ExpireWaitingVisitorEntries(ctx context.Context) (int64, error)
	ExpireApprovedVisitorEntries(ctx context.Context) (int64, error)
	ExpireVisitorInvites(ctx context.Context) (int64, error)
	ExpireFlatMemberInvites(ctx context.Context) (int64, error)
	ExpireSubscriptions(ctx context.Context) (int64, error)
}

type ExpiryJob struct {
	repo expiryStore
}

func NewExpiryJob(repo expiryStore) *ExpiryJob {
	return &ExpiryJob{repo: repo}
}

func (j *ExpiryJob) Start(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 15 * time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	j.runAndLog(ctx)
	for {
		select {
		case <-ticker.C:
			j.runAndLog(ctx)
		case <-ctx.Done():
			logger.Info("expiry job stopped")
			return
		}
	}
}

func (j *ExpiryJob) RunOnce(ctx context.Context) (RunResult, error) {
	locked, err := j.repo.WithAdvisoryLock(ctx, expiryJobLock, func(lockCtx context.Context) error {
		var operationErrors []error
		operations := []struct {
			name string
			run  func(context.Context) (int64, error)
		}{
			{"waiting visitor entries", j.repo.ExpireWaitingVisitorEntries},
			{"approved visitor entries", j.repo.ExpireApprovedVisitorEntries},
			{"visitor invites", j.repo.ExpireVisitorInvites},
			{"flat member invites", j.repo.ExpireFlatMemberInvites},
			{"subscriptions", j.repo.ExpireSubscriptions},
		}
		for _, operation := range operations {
			if err := lockCtx.Err(); err != nil {
				return err
			}
			count, operationErr := operation.run(lockCtx)
			if operationErr != nil {
				if isContextCancellation(operationErr) {
					return operationErr
				}
				jobLogger(lockCtx).Error("expiry operation failed", zap.String("operation", operation.name), zap.String("internal_error", logger.Sanitize(operationErr.Error())))
				operationErrors = append(operationErrors, fmt.Errorf("%s: %w", operation.name, operationErr))
				continue
			}
			if count > 0 {
				jobLogger(lockCtx).Info("expiry operation completed", zap.String("operation", operation.name), zap.Int64("count", count))
			}
		}
		return errors.Join(operationErrors...)
	})
	result := RunResult{LockAcquired: locked}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (j *ExpiryJob) runAndLog(ctx context.Context) {
	executeRun(ctx, "expiry", "", j.RunOnce)
}

func isContextCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
