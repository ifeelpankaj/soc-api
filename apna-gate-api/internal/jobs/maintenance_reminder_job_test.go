package jobs

import (
	"context"
	"errors"
	"testing"
	"time"
)

type reminderJobService struct {
	calls int
	err   error
}

func (s *reminderJobService) RunReminders(ctx context.Context) error {
	s.calls++
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > 5*time.Minute {
		return errors.New("missing bounded deadline")
	}
	return s.err
}

type reminderJobLock struct {
	acquired bool
	name     string
}

func (l *reminderJobLock) WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error) {
	l.name = name
	if !l.acquired {
		return false, nil
	}
	return true, fn(ctx)
}
func TestMaintenanceReminderJobLockAndErrors(t *testing.T) {
	s := &reminderJobService{}
	lock := &reminderJobLock{}
	job := NewMaintenanceReminderJob(s, lock)
	result, err := job.RunOnce(context.Background())
	if err != nil || result.LockAcquired || s.calls != 0 {
		t.Fatal("ran without lock", result, err)
	}
	lock.acquired = true
	s.err = errors.New("queue failed")
	result, err = job.RunOnce(context.Background())
	if !errors.Is(err, s.err) || !result.LockAcquired || s.calls != 1 || lock.name != "apna-gate/jobs/maintenance-reminders" {
		t.Fatal(result, err, s.calls, lock.name)
	}
}
