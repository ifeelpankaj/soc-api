package jobs

import (
	"context"
	"time"
)

const JobMaintenanceBilling = "maintenance-billing"

type maintenanceBillingService interface{ RunScheduled(context.Context) error }
type MaintenanceBillingJob struct {
	service maintenanceBillingService
	locks   interface {
		WithAdvisoryLock(context.Context, string, func(context.Context) error) (bool, error)
	}
}

func NewMaintenanceBillingJob(s maintenanceBillingService, locks interface {
	WithAdvisoryLock(context.Context, string, func(context.Context) error) (bool, error)
}) *MaintenanceBillingJob {
	return &MaintenanceBillingJob{s, locks}
}
func (j *MaintenanceBillingJob) RunOnce(ctx context.Context) (RunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	acquired, err := j.locks.WithAdvisoryLock(ctx, "apna-gate/jobs/maintenance-billing", j.service.RunScheduled)
	return RunResult{LockAcquired: acquired}, err
}
