package jobs

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	JobCleanup              = "cleanup"
	JobMonthlyVisitorReport = "monthly-visitor-report"
)

var (
	ErrUnknownJob        = errors.New("unknown job")
	ErrJobManagerStopped = errors.New("job manager is shutting down")
)

type Trigger = models.JobTriggerResponse

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	expiry  *ExpiryJob
	runners map[string]Runner

	mu            sync.Mutex
	running       map[string]string
	closing       bool
	expiryStarted bool
	wg            sync.WaitGroup
}

func NewManager(parent context.Context, expiry *ExpiryJob, cleanup *CleanupJob, monthly *MonthlyVisitorReportJob, maintenance ...*MaintenanceBillingJob) *Manager {
	runners := map[string]Runner{JobCleanup: cleanup, JobMonthlyVisitorReport: monthly}
	if len(maintenance) > 0 && maintenance[0] != nil {
		runners[JobMaintenanceBilling] = maintenance[0]
	}
	return newManager(parent, expiry, runners)
}

func newManager(parent context.Context, expiry *ExpiryJob, runners map[string]Runner) *Manager {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Manager{
		ctx:     ctx,
		cancel:  cancel,
		expiry:  expiry,
		runners: runners,
		running: make(map[string]string),
	}
}

// StartExpiry starts the tracked internal expiry loop. It is intentionally
// separate from construction so callers can bind the HTTP listener first.
func (m *Manager) StartExpiry(interval time.Duration) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closing || m.expiryStarted || m.expiry == nil {
		return false
	}
	m.expiryStarted = true
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.expiry.Start(m.ctx, interval)
	}()
	return true
}

// Trigger starts a known webhook job asynchronously. The active trigger ID is
// returned when the same job is already running in this process.
func (m *Manager) Trigger(name string) (Trigger, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closing {
		return Trigger{}, ErrJobManagerStopped
	}
	runner, ok := m.runners[name]
	if !ok || runner == nil {
		return Trigger{}, fmt.Errorf("%w: %s", ErrUnknownJob, name)
	}
	if activeTriggerID, running := m.running[name]; running {
		return Trigger{Job: name, TriggerID: activeTriggerID, Status: "already_running"}, nil
	}

	triggerID := uuid.NewString()
	m.running[name] = triggerID
	m.wg.Add(1)
	go m.run(name, triggerID, runner)
	return Trigger{Job: name, TriggerID: triggerID, Status: "accepted"}, nil
}

// StartTask persists admission before starting a separately audited maintenance task.
// Unlike Trigger, every call is independent; its worker owns cross-instance locking.
func (m *Manager) StartTask(prepare func() error, run func(context.Context)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closing {
		return ErrJobManagerStopped
	}
	if err := prepare(); err != nil {
		return err
	}
	m.wg.Add(1)
	go func() { defer m.wg.Done(); run(m.ctx) }()
	return nil
}

func (m *Manager) run(name, triggerID string, runner Runner) {
	defer m.wg.Done()
	defer func() {
		m.mu.Lock()
		if m.running[name] == triggerID {
			delete(m.running, name)
		}
		m.mu.Unlock()
	}()

	executeRun(m.ctx, name, triggerID, runner.RunOnce)
}

// Shutdown prevents new triggers, cancels all jobs, and waits up to timeout.
func (m *Manager) Shutdown(timeout time.Duration) bool {
	m.mu.Lock()
	if !m.closing {
		m.closing = true
		m.cancel()
	}
	m.mu.Unlock()

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()
	if timeout <= 0 {
		<-done
		return true
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
		return true
	case <-timer.C:
		return false
	}
}
