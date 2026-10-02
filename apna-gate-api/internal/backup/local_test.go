package backup

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-server/internal/config"
)

func TestPostgresCommandsAndPassfile(t *testing.T) {
	cfg := &config.Config{DBHost: "localhost", DBPort: 5432, DBName: "main_db", DBUser: "backup", DBPassword: `a:b\c`, DBSSLMode: "require"}
	p := NewPostgres(cfg, nil)
	dir := t.TempDir()
	calls := []string{}
	p.command = func(_ context.Context, env []string, name string, args []string, out io.Writer) error {
		calls = append(calls, name)
		all := strings.Join(args, " ")
		if strings.Contains(all, cfg.DBPassword) || strings.Contains(strings.Join(env, "\n"), "PGPASSWORD=") {
			t.Fatal("password exposed")
		}
		found := false
		for _, v := range env {
			if v == "PGPASSFILE="+filepath.Join(dir, "pgpass") {
				found = true
			}
		}
		if !found {
			t.Fatal("missing passfile")
		}
		switch name {
		case "pg_dump":
			if strings.Contains(all, "--jobs") || !strings.Contains(all, "--compress=6") || !strings.Contains(all, "--lock-wait-timeout=5000") {
				t.Fatal(all)
			}
			_, err := io.WriteString(out, "archive")
			return err
		case "pg_amcheck":
			if all != "--jobs=1" {
				t.Fatal(all)
			}
		case "pg_restore":
			if !strings.Contains(all, "--file="+os.DevNull) {
				t.Fatal(all)
			}
		}
		return nil
	}
	if err := p.Dump(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	if err := p.Check(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	if err := p.Validate(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "pgpass")) //nolint:gosec // Test-owned private temporary directory.
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "localhost:5432:main_db:backup:a\\:b\\\\c\n" {
		t.Fatalf("bad escaping: %q", b)
	}
	if len(calls) != 3 {
		t.Fatal(calls)
	}
	cfg.DBPassword = "bad\npassword"
	if _, err := p.environment(dir); err == nil {
		t.Fatal("accepted invalid pgpass input")
	}
}
func TestCommandCancellationAndRedaction(t *testing.T) {
	// Reinvoke the test process instead of relying on a platform-specific shell.
	if os.Getenv("BACKUP_TEST_CHILD") == "1" {
		time.Sleep(time.Minute)
		os.Exit(0)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := runCommand(ctx, append(os.Environ(), "BACKUP_TEST_CHILD=1"), os.Args[0], []string{"-test.run=TestCommandCancellationAndRedaction"}, io.Discard)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("cancellation: %v", err)
	}
	err = runCommand(context.Background(), nil, "nonexistent-backup-tool-secret", nil, io.Discard)
	if err == nil {
		t.Fatal("missing executable succeeded")
	}
}
func TestHashArchiveRejectsEmptyAndCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dump")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if hashArchive(context.Background(), path, &Run{}) == nil {
		t.Fatal("accepted empty archive")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(hashArchive(ctx, path, &Run{}), context.Canceled) {
		t.Fatal("hash ignored cancellation")
	}
}
