package backup

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"go-server/internal/config"
	"go-server/internal/jobs"
)

type memoryStore struct {
	mu        sync.Mutex
	runs      map[string]Run
	history   map[string][]string
	createErr error
}

func newMemoryStore() *memoryStore {
	return &memoryStore{runs: map[string]Run{}, history: map[string][]string{}}
}
func (m *memoryStore) Create(_ context.Context, r *Run) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.createErr != nil {
		return m.createErr
	}
	r.Created = time.Now().UTC()
	m.runs[r.ID] = *r
	m.history[r.ID] = []string{"queued"}
	return nil
}
func (m *memoryStore) Save(_ context.Context, r *Run) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runs[r.ID] = *r
	m.history[r.ID] = append(m.history[r.ID], r.Status)
	return nil
}
func (m *memoryStore) Get(_ context.Context, id string) (*Run, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.runs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &r, nil
}

type fakeLocker struct {
	mu       sync.Mutex
	held     bool
	pulseErr bool
}
type fakeLease struct{ l *fakeLocker }

func (l *fakeLocker) Try(context.Context) (Lease, bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.held {
		return nil, false, nil
	}
	l.held = true
	return &fakeLease{l}, true, nil
}
func (l *fakeLease) Pulse(context.Context, *Run) error {
	if l.l.pulseErr {
		return errors.New("lost lock")
	}
	return nil
}
func (l *fakeLease) Recover(context.Context, *Run) ([]string, error) { return nil, nil }
func (l *fakeLease) Close()                                          { l.l.mu.Lock(); defer l.l.mu.Unlock(); l.l.held = false }

type fakeLocal struct {
	preflightErr error
	fail         string
	block        chan struct{}
	entered      chan struct{}
	checks       int
}

func (l *fakeLocal) Preflight(_ context.Context, r *Run, _ string) error {
	if l.preflightErr != nil {
		return l.preflightErr
	}
	r.Checks.Database = "passed"
	r.Checks.PageChecksums = "disabled"
	if l.fail == "preflight" {
		return errors.New("secret-password")
	}
	return nil
}

func TestPreflightDiagnosticIsPersisted(t *testing.T) {
	local := &fakeLocal{preflightErr: preflightFailure("tool_unavailable", "pg_dump --version failed; install PostgreSQL 15 client tools")}
	s, _, manager := setupService(t, local, &fakeCloud{})
	defer manager.Shutdown(time.Second)
	r, err := s.Trigger(context.Background(), "daily")
	if err != nil {
		t.Fatal(err)
	}
	r = awaitRun(t, s, r.ID)
	if r.Status != "failed" || r.ErrorCode != "preflight_tool_unavailable" || r.ErrorMessage != local.preflightErr.Error() {
		t.Fatalf("missing safe preflight diagnostic: %+v", r)
	}
}
func (l *fakeLocal) Check(context.Context, string) error {
	l.checks++
	if l.fail == "checking" {
		return errors.New("corruption")
	}
	return nil
}
func (l *fakeLocal) Dump(ctx context.Context, dir string) error {
	if l.entered != nil {
		close(l.entered)
	}
	if l.block != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-l.block:
		}
	}
	if l.fail == "dumping" {
		return errors.New("password=secret-password")
	}
	return os.WriteFile(filepath.Join(dir, "database.dump"), []byte("test backup bytes"), 0600)
}
func (l *fakeLocal) Validate(context.Context, string) error {
	if l.fail == "validating" {
		return errors.New("malformed archive")
	}
	return nil
}

type fakeCloud struct {
	fail                string
	uploads, retentions int
}

func (c *fakeCloud) Preflight(context.Context) error {
	if c.fail == "preflight" {
		return errors.New("oauth-secret")
	}
	return nil
}
func (c *fakeCloud) Upload(context.Context, *Run, string) (string, error) {
	c.uploads++
	if c.fail == "uploading" {
		return "", errors.New("oauth-secret")
	}
	return "drive-file", nil
}
func (c *fakeCloud) Verify(_ context.Context, r *Run) (string, error) {
	if c.fail == "verifying" {
		return "", errors.New("checksum mismatch")
	}
	return r.LocalMD5, nil
}
func (c *fakeCloud) Retain(context.Context, *Run) (int, error) {
	c.retentions++
	if c.fail == "retention" {
		return 1, errors.New("permission denied")
	}
	return 2, nil
}
func setupService(t *testing.T, local *fakeLocal, cloud *fakeCloud) (*Service, *memoryStore, *jobs.Manager) {
	t.Helper()
	store := newMemoryStore()
	manager := jobs.NewManager(context.Background(), nil, nil, nil)
	s := NewService(config.BackupConfig{Enabled: true, Deployment: "test", TempDir: t.TempDir(), Timeout: time.Minute}, "main_db", store, &fakeLocker{}, local, cloud, manager)
	root, err := filepath.EvalSymlinks(s.cfg.TempDir)
	if err != nil {
		t.Fatal("resolve test backup directory:", err)
	}
	if err := os.Chmod(root, 0700); err != nil { //nolint:gosec // Private directory, not a regular file.
		t.Fatal("protect test backup directory:", err)
	}
	t.Cleanup(func() { manager.Shutdown(5 * time.Second) })
	return s, store, manager
}
func awaitRun(t *testing.T, s *Service, id string) *Run {
	t.Helper()
	deadline := time.After(5 * time.Second)
	for {
		r, err := s.Get(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}
		if r.Completed != nil {
			return r
		}
		select {
		case <-deadline:
			t.Fatal("run did not finish")
		case <-time.After(time.Millisecond):
		}
	}
}
func TestPipelineDailyAndWeekly(t *testing.T) {
	for _, kind := range []string{"daily", "weekly"} {
		t.Run(kind, func(t *testing.T) {
			local, cloud := &fakeLocal{}, &fakeCloud{}
			s, store, manager := setupService(t, local, cloud)
			accepted, err := s.Trigger(context.Background(), kind)
			if err != nil {
				t.Fatal(err)
			}
			if accepted.Status != "queued" {
				t.Fatal("response must describe persisted admission")
			}
			r := awaitRun(t, s, accepted.ID)
			manager.Shutdown(time.Second)
			if r.Status != "completed" || r.Checks.PageChecksums != "disabled" || len(r.SHA256) != 64 || len(r.LocalMD5) != 32 || r.DriveMD5 != r.LocalMD5 || r.ArchiveSize != 17 {
				t.Fatalf("unexpected run: %+v", r)
			}
			want := []string{"queued", "preflight"}
			if kind == "weekly" {
				want = append(want, "checking")
				if local.checks != 1 {
					t.Fatal("missing weekly scan")
				}
			} else if local.checks != 0 {
				t.Fatal("daily ran amcheck")
			}
			want = append(want, "dumping", "validating", "uploading", "verifying", "retention", "completed")
			store.mu.Lock()
			got := append([]string(nil), store.history[r.ID]...)
			store.mu.Unlock()
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("states %v, want %v", got, want)
			}
			entries, _ := os.ReadDir(s.cfg.TempDir)
			if len(entries) != 0 {
				t.Fatal("temporary archive not cleaned")
			}
		})
	}
}
func TestPipelineFailuresNeverPruneAndRedact(t *testing.T) {
	for _, stage := range []string{"preflight", "checking", "dumping", "validating", "uploading", "verifying", "retention"} {
		t.Run(stage, func(t *testing.T) {
			local, cloud := &fakeLocal{fail: stage}, &fakeCloud{fail: stage}
			s, _, manager := setupService(t, local, cloud)
			r, err := s.Trigger(context.Background(), "weekly")
			if err != nil {
				t.Fatal(err)
			}
			r = awaitRun(t, s, r.ID)
			manager.Shutdown(time.Second)
			if stage == "retention" {
				if r.Status != "completed" || r.Retention.Status != "warning" || r.Retention.Deleted != 1 {
					t.Fatalf("retention lost success: %+v", r)
				}
				return
			}
			if r.Status != "failed" || cloud.retentions != 0 {
				t.Fatalf("unsafe failure: %+v", r)
			}
			if strings.Contains(r.ErrorMessage, "secret") {
				t.Fatal("leaked credential")
			}
			if stage == "checking" && r.Checks.Amcheck != "failed" {
				t.Fatal("missing amcheck failure")
			}
		})
	}
}
func TestPersistBeforeLaunchAndDisabled(t *testing.T) {
	s, store, _ := setupService(t, &fakeLocal{}, &fakeCloud{})
	store.createErr = errors.New("database down")
	if _, err := s.Trigger(context.Background(), "daily"); err == nil {
		t.Fatal("accepted unpersisted backup")
	}
	if len(store.runs) != 0 {
		t.Fatal("unexpected run")
	}
	s.cfg.Enabled = false
	if _, err := s.Trigger(context.Background(), "daily"); !errors.Is(err, ErrDisabled) {
		t.Fatal(err)
	}
}
func TestConcurrentRequestsHaveSeparateAuditRecords(t *testing.T) {
	local := &fakeLocal{block: make(chan struct{}), entered: make(chan struct{})}
	s, store, manager := setupService(t, local, &fakeCloud{})
	first, err := s.Trigger(context.Background(), "daily")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-local.entered:
	case <-time.After(5 * time.Second):
		r, _ := s.Get(context.Background(), first.ID)
		t.Fatalf("dump never started: %+v", r)
	}
	m2 := jobs.NewManager(context.Background(), nil, nil, nil)
	defer m2.Shutdown(time.Second)
	s2 := NewService(s.cfg, s.database, store, s.locker, &fakeLocal{}, &fakeCloud{}, m2)
	second, err := s2.Trigger(context.Background(), "weekly")
	if err != nil {
		t.Fatal(err)
	}
	r := awaitRun(t, s2, second.ID)
	if r.ID == first.ID || r.Status != "skipped_concurrent" {
		t.Fatalf("overlap: %+v", r)
	}
	manager.Shutdown(time.Second)
	r = awaitRun(t, s, first.ID)
	if r.Status != "interrupted" {
		t.Fatalf("shutdown: %+v", r)
	}
}
func TestTimeoutAndLostLease(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		s, _, _ := setupService(t, &fakeLocal{block: make(chan struct{})}, &fakeCloud{})
		s.cfg.Timeout = 20 * time.Millisecond
		r, err := s.Trigger(context.Background(), "daily")
		if err != nil {
			t.Fatal(err)
		}
		r = awaitRun(t, s, r.ID)
		if r.ErrorCode != "backup_timeout" {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("lease", func(t *testing.T) {
		s, _, _ := setupService(t, &fakeLocal{}, &fakeCloud{})
		s.locker = &fakeLocker{pulseErr: true}
		r, err := s.Trigger(context.Background(), "daily")
		if err != nil {
			t.Fatal(err)
		}
		r = awaitRun(t, s, r.ID)
		if r.Status != "interrupted" {
			t.Fatalf("%+v", r)
		}
	})
}
func TestTempCleanupRejectsEscapes(t *testing.T) {
	root := t.TempDir()
	for _, id := range []string{"..", "../outside", "/", "not-a-uuid"} {
		if removeRunDir(root, id) == nil {
			t.Fatalf("accepted %q", id)
		}
	}
}
