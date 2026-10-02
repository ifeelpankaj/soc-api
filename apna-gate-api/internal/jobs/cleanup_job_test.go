package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-server/internal/models"
)

type cleanupStoreFake struct {
	pending           []models.PendingImageDeletion
	completed         []int64
	orphanCalls       int
	notificationCalls int
	entryCalls        int
	entryFrom         time.Time
	entryTo           time.Time
	entryWindow       models.VisitorReportCleanupWindow
	timezone          string
	notificationErr   error
	verificationCalls int
}

func (f *cleanupStoreFake) DeleteOrphanVisitorsBatch(context.Context, time.Time, int32) (int64, error) {
	f.orphanCalls++
	return 0, nil
}
func (f *cleanupStoreFake) ListVisitorImageDeletions(_ context.Context, afterID int64, size int32) ([]models.PendingImageDeletion, error) {
	var result []models.PendingImageDeletion
	for _, item := range f.pending {
		if item.ID > afterID && len(result) < int(size) {
			result = append(result, item)
		}
	}
	return result, nil
}
func (f *cleanupStoreFake) CompleteVisitorImageDeletion(_ context.Context, id int64) error {
	f.completed = append(f.completed, id)
	for i, item := range f.pending {
		if item.ID == id {
			f.pending = append(f.pending[:i], f.pending[i+1:]...)
			break
		}
	}
	return nil
}

func (f *cleanupStoreFake) WithAdvisoryLock(ctx context.Context, _ string, fn func(context.Context) error) (bool, error) {
	return true, fn(ctx)
}
func (f *cleanupStoreFake) DeleteNotificationsBatch(context.Context, time.Time, int32) (int64, error) {
	f.notificationCalls++
	if f.notificationErr != nil {
		return 0, f.notificationErr
	}
	if f.notificationCalls == 1 {
		return 2, nil
	}
	return 0, nil
}
func (f *cleanupStoreFake) DeleteVerificationsBatch(context.Context, int32) (int64, error) {
	f.verificationCalls++
	return 0, nil
}
func (f *cleanupStoreFake) ListVisitorEntryCleanupWindows(_ context.Context, timezone string, _ time.Time) ([]models.VisitorReportCleanupWindow, error) {
	f.timezone = timezone
	return []models.VisitorReportCleanupWindow{{
		SocietyID: 4, ReportMonth: time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
	}}, nil
}
func (f *cleanupStoreFake) DeleteVisitorEntriesBatch(_ context.Context, window models.VisitorReportCleanupWindow, from, to, _ time.Time, _ int32) (int64, error) {
	f.entryCalls++
	f.entryWindow, f.entryFrom, f.entryTo = window, from, to
	if f.entryCalls == 1 {
		return 2, nil
	}
	return 0, nil
}
func (f *cleanupStoreFake) DeleteVisitorInvitesBatch(context.Context, time.Time, int32) (int64, error) {
	return 0, nil
}
func (f *cleanupStoreFake) DeleteFlatMemberInvitesBatch(context.Context, time.Time, int32) (int64, error) {
	return 0, nil
}

func TestCleanupJobUsesBoundedBatchesAndSharedMonthBounds(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	store := &cleanupStoreFake{}
	job := NewCleanupJob(store, CleanupJobConfig{
		Location: loc, Timezone: "Asia/Kolkata", BatchSize: 2,
		VisitorEntryRetention: 90 * 24 * time.Hour, VisitorInviteRetention: 90 * 24 * time.Hour,
		FlatMemberInviteRetention: 90 * 24 * time.Hour, NotificationRetention: 30 * 24 * time.Hour,
	})
	job.now = func() time.Time { return time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC) }
	if _, err := job.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	if store.notificationCalls != 2 || store.entryCalls != 2 {
		t.Fatalf("batch calls = notifications:%d entries:%d, want 2 each", store.notificationCalls, store.entryCalls)
	}
	if store.timezone != "Asia/Kolkata" || store.entryWindow.SocietyID != 4 {
		t.Fatalf("cleanup scope = timezone:%q window:%+v", store.timezone, store.entryWindow)
	}
	wantFrom := time.Date(2026, time.August, 31, 18, 30, 0, 0, time.UTC)
	wantTo := time.Date(2026, time.September, 30, 18, 30, 0, 0, time.UTC)
	if !store.entryFrom.Equal(wantFrom) || !store.entryTo.Equal(wantTo) {
		t.Fatalf("entry bounds = (%v, %v), want (%v, %v)", store.entryFrom, store.entryTo, wantFrom, wantTo)
	}
}

func TestCleanupJobStopsAfterCancellation(t *testing.T) {
	store := &cleanupStoreFake{notificationErr: context.Canceled}
	job := NewCleanupJob(store, CleanupJobConfig{Location: time.UTC, Timezone: "UTC", BatchSize: 2})
	_, err := job.RunOnce(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RunOnce() error = %v, want context canceled", err)
	}
	if store.verificationCalls != 0 {
		t.Fatalf("verification cleanup calls = %d, want 0", store.verificationCalls)
	}
}

type imageDeleterFunc func(context.Context, string) error

func (f imageDeleterFunc) Delete(ctx context.Context, id string) error { return f(ctx, id) }

func TestCleanupImagesRetriesFailuresAcrossRuns(t *testing.T) {
	store := &cleanupStoreFake{pending: []models.PendingImageDeletion{{ID: 1, FileID: "fail"}, {ID: 2, FileID: "ok"}, {ID: 3, FileID: "also-ok"}}}
	var calls []string
	job := NewCleanupJob(store, CleanupJobConfig{BatchSize: 2, ImageStorage: imageDeleterFunc(func(ctx context.Context, id string) error {
		calls = append(calls, id)
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("missing provider timeout")
		}
		if id == "fail" {
			return errors.New("provider unavailable")
		}
		return nil
	})})
	if err := job.cleanupImages(context.Background()); err == nil {
		t.Fatal("expected failure")
	}
	if len(calls) != 3 || len(store.pending) != 1 || store.pending[0].ID != 1 {
		t.Fatalf("calls=%v pending=%v", calls, store.pending)
	}
	// A new job instance models a restart; persisted failed work remains available.
	restarted := NewCleanupJob(store, CleanupJobConfig{BatchSize: 2, ImageStorage: imageDeleterFunc(func(context.Context, string) error { return nil })})
	if err := restarted.cleanupImages(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(store.pending) != 0 {
		t.Fatalf("pending=%v", store.pending)
	}
}

func TestCleanupImagesDisabledAndCancelled(t *testing.T) {
	store := &cleanupStoreFake{pending: []models.PendingImageDeletion{{ID: 1, FileID: "photo"}}}
	job := NewCleanupJob(store, CleanupJobConfig{})
	if err := job.cleanupImages(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(store.pending) != 1 {
		t.Fatal("disabled storage discarded pending work")
	}
	ctx, cancel := context.WithCancel(context.Background())
	job.cfg.ImageStorage = imageDeleterFunc(func(context.Context, string) error { cancel(); return context.Canceled })
	if err := job.cleanupImages(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
	if len(store.pending) != 1 {
		t.Fatal("cancelled deletion discarded pending work")
	}
}
