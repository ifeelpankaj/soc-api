package jobs

import (
	"context"
	"errors"
	"testing"
)

type expiryStoreFake struct {
	calls      []string
	failOn     string
	cancelOn   string
	lockCalled bool
}

func (f *expiryStoreFake) WithAdvisoryLock(ctx context.Context, _ string, fn func(context.Context) error) (bool, error) {
	f.lockCalled = true
	return true, fn(ctx)
}
func (f *expiryStoreFake) result(name string) (int64, error) {
	f.calls = append(f.calls, name)
	if f.cancelOn == name {
		return 0, context.Canceled
	}
	if f.failOn == name {
		return 0, errors.New("test failure")
	}
	return 1, nil
}
func (f *expiryStoreFake) ExpireWaitingVisitorEntries(context.Context) (int64, error) {
	return f.result("waiting")
}
func (f *expiryStoreFake) ExpireApprovedVisitorEntries(context.Context) (int64, error) {
	return f.result("approved")
}
func (f *expiryStoreFake) ExpireVisitorInvites(context.Context) (int64, error) {
	return f.result("visitor_invites")
}
func (f *expiryStoreFake) ExpireFlatMemberInvites(context.Context) (int64, error) {
	return f.result("flat_invites")
}
func (f *expiryStoreFake) ExpireSubscriptions(context.Context) (int64, error) {
	return f.result("subscriptions")
}

func TestExpiryJobRunsEveryOperationWhenOneFails(t *testing.T) {
	store := &expiryStoreFake{failOn: "approved"}
	_, err := NewExpiryJob(store).RunOnce(context.Background())
	if err == nil {
		t.Fatal("expiry job did not return its operation failure")
	}
	want := []string{"waiting", "approved", "visitor_invites", "flat_invites", "subscriptions"}
	if !store.lockCalled {
		t.Fatal("expiry job did not acquire its advisory lock")
	}
	if len(store.calls) != len(want) {
		t.Fatalf("calls = %#v, want %#v", store.calls, want)
	}
	for i := range want {
		if store.calls[i] != want[i] {
			t.Fatalf("calls = %#v, want %#v", store.calls, want)
		}
	}
}

func TestExpiryJobStopsRemainingOperationsWhenCancelled(t *testing.T) {
	store := &expiryStoreFake{cancelOn: "approved"}
	_, err := NewExpiryJob(store).RunOnce(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RunOnce() error = %v, want context canceled", err)
	}
	want := []string{"waiting", "approved"}
	if len(store.calls) != len(want) {
		t.Fatalf("calls after cancellation = %#v, want %#v", store.calls, want)
	}
}
