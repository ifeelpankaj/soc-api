package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/pkg/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunLifecycleContract(t *testing.T) {
	dir := os.Getenv("JOB_LOG_TEST_DIR")
	if dir == "" {
		//nolint:gosec // Relaunch this test binary with a fixed test selector; no application input.
		cmd := exec.CommandContext(
			context.Background(),
			os.Args[0],
			"-test.run=^TestRunLifecycleContract$",
		)
		cmd.Env = append(os.Environ(), "JOB_LOG_TEST_DIR="+t.TempDir())
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		return
	}
	cfg := logger.DefaultConfig()
	cfg.LogDir = dir
	if err := logger.InitLogger(cfg); err != nil {
		t.Fatal(err)
	}
	for _, outcome := range []string{OutcomeSuccess, OutcomeFailure, OutcomeSkipped, OutcomeInterrupted} {
		executeRun(context.Background(), "test-job", outcome, func(ctx context.Context) (RunResult, error) {
			jobLogger(ctx).Info("Operation")
			switch outcome {
			case OutcomeFailure:
				return RunResult{}, errors.New("token=private")
			case OutcomeInterrupted:
				return RunResult{}, context.Canceled
			case OutcomeSkipped:
				return RunResult{}, nil
			}
			return RunResult{LockAcquired: true}, nil
		})
	}
	_ = logger.Sync()
	data, err := os.ReadFile(filepath.Join(dir, "info.production.log")) //nolint:gosec // Parent test supplies its private temporary directory to the child.
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatal(err)
		}
		run, ok := entry["run_id"].(string)
		if !ok {
			continue
		}
		if entry["job_name"] != "test-job" {
			t.Fatal("lost job context")
		}
		if counts[run] == nil {
			counts[run] = map[string]int{}
		}
		event, _ := entry["event"].(string)
		counts[run][event]++
		if event == "job_completed" && entry["outcome"] != run {
			t.Fatal("incorrect outcome")
		}
	}
	for _, outcome := range []string{OutcomeSuccess, OutcomeFailure, OutcomeSkipped, OutcomeInterrupted} {
		if counts[outcome]["job_started"] != 1 || counts[outcome]["job_completed"] != 1 || counts[outcome][""] != 1 {
			t.Fatalf("incorrect lifecycle: %v", counts)
		}
	}
	if strings.Contains(string(data), "private") {
		t.Fatal("credential leaked")
	}
}
