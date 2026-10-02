package jobs

import "context"

// RunResult reports whether this process acquired the distributed job lock.
type RunResult struct {
	LockAcquired bool
}

// Runner is a background job that can be triggered once.
type Runner interface {
	RunOnce(ctx context.Context) (RunResult, error)
}
