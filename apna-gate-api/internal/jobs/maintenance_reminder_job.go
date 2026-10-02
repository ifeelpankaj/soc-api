package jobs

import (
	"context"
	"time"
)

const JobMaintenanceReminders = "maintenance-reminders"

type MaintenanceReminderJob struct {
	service interface{ RunReminders(context.Context) error }
	locks   interface {
		WithAdvisoryLock(context.Context, string, func(context.Context) error) (bool, error)
	}
}

func NewMaintenanceReminderJob(s interface{ RunReminders(context.Context) error }, locks interface {
	WithAdvisoryLock(context.Context, string, func(context.Context) error) (bool, error)
}) *MaintenanceReminderJob {
	return &MaintenanceReminderJob{service: s, locks: locks}
}

func (j *MaintenanceReminderJob) RunOnce(ctx context.Context) (RunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	locked, err := j.locks.WithAdvisoryLock(ctx, "apna-gate/jobs/maintenance-reminders", j.service.RunReminders)
	return RunResult{LockAcquired: locked}, err
}

func (m *Manager) RegisterMaintenanceReminders(job *MaintenanceReminderJob) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runners[JobMaintenanceReminders] = job
}
