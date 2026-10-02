package jobs

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type runnerFunc func(context.Context) (RunResult, error)

func (f runnerFunc) RunOnce(ctx context.Context) (RunResult, error) {
	return f(ctx)
}

func TestManagerConstructionDoesNotRunWebhookJobs(t *testing.T) {
	var calls atomic.Int32
	runner := runnerFunc(func(context.Context) (RunResult, error) {
		calls.Add(1)
		return RunResult{LockAcquired: true}, nil
	})
	manager := newManager(context.Background(), nil, map[string]Runner{
		JobCleanup: runner, JobMonthlyVisitorReport: runner,
	})
	defer manager.Shutdown(time.Second)
	if calls.Load() != 0 {
		t.Fatalf("webhook job calls during construction = %d, want 0", calls.Load())
	}
}

func TestManagerTriggerIsAsyncAndReturnsActiveTriggerID(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan struct{})
	runner := runnerFunc(func(context.Context) (RunResult, error) {
		close(started)
		<-release
		close(finished)
		return RunResult{LockAcquired: true}, nil
	})
	manager := newManager(context.Background(), nil, map[string]Runner{JobCleanup: runner})

	first, err := manager.Trigger(JobCleanup)
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != "accepted" || first.TriggerID == "" {
		t.Fatalf("first trigger = %+v", first)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("asynchronous job did not start")
	}

	duplicate, err := manager.Trigger(JobCleanup)
	if err != nil {
		t.Fatal(err)
	}
	if duplicate.Status != "already_running" || duplicate.TriggerID != first.TriggerID {
		t.Fatalf("duplicate trigger = %+v, want active ID %q", duplicate, first.TriggerID)
	}
	close(release)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("job did not continue after Trigger returned")
	}
	if !manager.Shutdown(time.Second) {
		t.Fatal("manager did not shut down")
	}
}

func TestManagerRegistryExcludesExpiry(t *testing.T) {
	manager := newManager(context.Background(), nil, map[string]Runner{JobCleanup: runnerFunc(func(context.Context) (RunResult, error) {
		return RunResult{}, nil
	})})
	defer manager.Shutdown(time.Second)

	if _, err := manager.Trigger("expiry"); !errors.Is(err, ErrUnknownJob) {
		t.Fatalf("Trigger(expiry) error = %v, want ErrUnknownJob", err)
	}
}

func TestManagerShutdownCancelsAndWaitsForWebhookJob(t *testing.T) {
	started := make(chan struct{})
	finished := make(chan struct{})
	runner := runnerFunc(func(ctx context.Context) (RunResult, error) {
		close(started)
		<-ctx.Done()
		close(finished)
		return RunResult{LockAcquired: true}, ctx.Err()
	})
	manager := newManager(context.Background(), nil, map[string]Runner{JobCleanup: runner})
	if _, err := manager.Trigger(JobCleanup); err != nil {
		t.Fatal(err)
	}
	<-started
	if !manager.Shutdown(time.Second) {
		t.Fatal("manager timed out during shutdown")
	}
	select {
	case <-finished:
	default:
		t.Fatal("Shutdown returned before webhook job exited")
	}
	if _, err := manager.Trigger(JobCleanup); !errors.Is(err, ErrJobManagerStopped) {
		t.Fatalf("Trigger after shutdown error = %v, want ErrJobManagerStopped", err)
	}
}

type trackedExpiryStore struct {
	once    sync.Once
	started chan struct{}
}

func (s *trackedExpiryStore) WithAdvisoryLock(ctx context.Context, _ string, fn func(context.Context) error) (bool, error) {
	s.once.Do(func() { close(s.started) })
	return true, fn(ctx)
}
func (*trackedExpiryStore) ExpireWaitingVisitorEntries(context.Context) (int64, error) {
	return 0, nil
}
func (*trackedExpiryStore) ExpireApprovedVisitorEntries(context.Context) (int64, error) {
	return 0, nil
}
func (*trackedExpiryStore) ExpireVisitorInvites(context.Context) (int64, error) { return 0, nil }
func (*trackedExpiryStore) ExpireFlatMemberInvites(context.Context) (int64, error) {
	return 0, nil
}
func (*trackedExpiryStore) ExpireSubscriptions(context.Context) (int64, error) { return 0, nil }

func TestManagerTracksInternalExpiryLoop(t *testing.T) {
	store := &trackedExpiryStore{started: make(chan struct{})}
	manager := newManager(context.Background(), NewExpiryJob(store), map[string]Runner{})
	if !manager.StartExpiry(time.Hour) {
		t.Fatal("expiry loop was not started")
	}
	select {
	case <-store.started:
	case <-time.After(time.Second):
		t.Fatal("expiry did not run immediately")
	}
	if manager.StartExpiry(time.Hour) {
		t.Fatal("expiry loop started twice")
	}
	if !manager.Shutdown(time.Second) {
		t.Fatal("manager did not wait for expiry loop")
	}
}
