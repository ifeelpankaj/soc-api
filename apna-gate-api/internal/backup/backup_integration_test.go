//go:build integration

package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go-server/internal/config"
	"go-server/internal/jobs"
)

func integrationDB(t *testing.T) (*postgres.PostgresContainer, *pgxpool.Pool, *config.Config) {
	t.Helper()
	ctx := t.Context()
	c, err := postgres.Run(ctx, "postgres:15-alpine", postgres.WithDatabase("backup_source"), postgres.WithUsername("postgres"), postgres.WithPassword("backup-test-password"), postgres.BasicWaitStrategies())
	if err != nil {
		t.Fatalf("PostgreSQL 15 integration environment unavailable: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(c) })
	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	b, err := os.ReadFile(filepath.Join("..", "..", "migrations", "21_backup_runs.sql"))
	if err != nil {
		t.Fatal(err)
	}
	up, _, _ := strings.Cut(string(b), "-- +migrate Down")
	if _, err := pool.Exec(ctx, up); err != nil {
		t.Fatal(err)
	}
	conn := pool.Config().ConnConfig
	cfg := &config.Config{DBHost: conn.Host, DBPort: int(conn.Port), DBUser: conn.User, DBPassword: conn.Password, DBName: conn.Database, DBSSLMode: "disable"}
	return c, pool, cfg
}

// Execute actual PG15 tools in the disposable server container, retaining the
// production argument construction and file validation on the test host.
func containerCommand(t *testing.T, c *postgres.PostgresContainer) Command {
	t.Helper()
	return func(ctx context.Context, env []string, name string, args []string, out io.Writer) error {
		var childEnv []string
		for _, v := range env {
			if strings.HasPrefix(v, "PG") && !strings.HasPrefix(v, "PGHOST=") && !strings.HasPrefix(v, "PGPORT=") && !strings.HasPrefix(v, "PGPASSFILE=") {
				childEnv = append(childEnv, v)
			}
		}
		childEnv = append(childEnv, "PGHOST=127.0.0.1", "PGPORT=5432", "PGPASSFILE=/tmp/backup-test-pgpass")
		if err := c.CopyToContainer(ctx, []byte("127.0.0.1:5432:backup_source:postgres:backup-test-password\n"), "/tmp/backup-test-pgpass", 0600); err != nil {
			return err
		}
		converted := append([]string(nil), args...)
		for i, arg := range converted {
			if arg == "--file="+os.DevNull {
				converted[i] = "--file=/dev/null"
			}
			if filepath.Base(arg) == "database.dump" {
				if err := c.CopyFileToContainer(ctx, arg, "/tmp/backup-test.dump", 0600); err != nil {
					return err
				}
				converted[i] = "/tmp/backup-test.dump"
			}
		}
		code, reader, err := c.Exec(ctx, append([]string{name}, converted...), tcexec.WithEnv(childEnv), tcexec.Multiplexed())
		if err != nil {
			return err
		}
		if code != 0 {
			b, _ := io.ReadAll(reader)
			return fmt.Errorf("%s exit %d: %s", name, code, b)
		}
		_, err = io.Copy(out, reader)
		return err
	}
}

type capturingCloud struct {
	fakeCloud
	archive []byte
}

func (c *capturingCloud) Upload(ctx context.Context, r *Run, path string) (string, error) {
	var err error
	c.archive, err = os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return c.fakeCloud.Upload(ctx, r, path)
}

func TestIntegrationBackupRoundTrip(t *testing.T) {
	c, pool, cfg := integrationDB(t)
	ctx := t.Context()
	seed := `CREATE EXTENSION amcheck;
CREATE SCHEMA demo;
CREATE TABLE demo.societies(id bigserial PRIMARY KEY,name text NOT NULL UNIQUE);
CREATE TABLE demo.visitors(id bigserial PRIMARY KEY,society_id bigint REFERENCES demo.societies(id),name text NOT NULL,notes text,photo bytea);
INSERT INTO demo.societies(name) VALUES ('Apna Gate');
INSERT INTO demo.visitors(society_id,name,notes,photo) VALUES (1,'नमस्ते',repeat(md5('backup'),5000),decode('001122aaff','hex'));
CREATE INDEX visitors_name_idx ON demo.visitors(name);
CREATE VIEW demo.visitor_names AS SELECT name FROM demo.visitors;
CREATE MATERIALIZED VIEW demo.visitor_count AS SELECT count(*) FROM demo.visitors;`
	if _, err := pool.Exec(ctx, seed); err != nil {
		t.Fatal(err)
	}
	store := &PGStore{Pool: pool}
	local := NewPostgres(cfg, pool)
	local.command = containerCommand(t, c)
	cloud := &capturingCloud{}
	manager := jobs.NewManager(ctx, nil, nil, nil)
	defer manager.Shutdown(5 * time.Second)
	s := NewService(config.BackupConfig{Enabled: true, Deployment: "integration", TempDir: t.TempDir(), Timeout: 2 * time.Minute}, cfg.DBName, store, &PGLocker{Pool: pool}, local, cloud, manager)
	r, err := s.Trigger(ctx, "weekly")
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.After(2 * time.Minute)
	for {
		r, err = s.Get(ctx, r.ID)
		if err != nil {
			t.Fatal(err)
		}
		if r.Completed != nil {
			break
		}
		select {
		case <-deadline:
			t.Fatal("backup timed out")
		case <-time.After(20 * time.Millisecond):
		}
	}
	manager.Shutdown(5 * time.Second)
	if r.Status != "completed" || r.Checks.Amcheck != "passed" || r.Checks.Archive != "passed" || r.Checks.PageChecksums != "disabled" {
		t.Fatalf("backup failed: %+v", r)
	}
	if r.ArchiveSize != int64(len(cloud.archive)) || len(r.SHA256) != 64 {
		t.Fatal("archive metadata mismatch")
	}
	if _, err := pool.Exec(ctx, "CREATE DATABASE backup_restored"); err != nil {
		t.Fatal(err)
	}
	if err := c.CopyToContainer(ctx, cloud.archive, "/tmp/restore-verified.dump", 0600); err != nil {
		t.Fatal(err)
	}
	code, output, err := c.Exec(ctx, []string{"pg_restore", "--exit-on-error", "--username=postgres", "--dbname=backup_restored", "/tmp/restore-verified.dump"}, tcexec.Multiplexed())
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		b, _ := io.ReadAll(output)
		t.Fatalf("restore failed: %s", b)
	}
	restoredCfg := pool.Config().Copy()
	restoredCfg.ConnConfig.Database = "backup_restored"
	restored, err := pgxpool.NewWithConfig(ctx, restoredCfg)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	var name, hexPhoto string
	var notesLen, count int
	if err := restored.QueryRow(ctx, `SELECT name,length(notes),encode(photo,'hex') FROM demo.visitors WHERE society_id=1`).Scan(&name, &notesLen, &hexPhoto); err != nil {
		t.Fatal(err)
	}
	if name != "नमस्ते" || notesLen != 160000 || hexPhoto != "001122aaff" {
		t.Fatal("restored data differs")
	}
	if err := restored.QueryRow(ctx, `SELECT count FROM demo.visitor_count`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("materialized view: %d %v", count, err)
	}
	var nextID int
	if err := restored.QueryRow(ctx, `INSERT INTO demo.societies(name) VALUES ('Restored') RETURNING id`).Scan(&nextID); err != nil || nextID != 2 {
		t.Fatalf("sequence: %d %v", nextID, err)
	}
	if _, err := restored.Exec(ctx, `INSERT INTO demo.visitors(society_id,name) VALUES (999,'invalid')`); err == nil {
		t.Fatal("foreign key was not restored")
	}
	if err := restored.QueryRow(ctx, `SELECT count(*) FROM demo.visitor_names`).Scan(&count); err != nil || count != 1 {
		t.Fatal("view was not restored")
	}
	entries, _ := os.ReadDir(s.cfg.TempDir)
	if len(entries) != 0 {
		t.Fatal("temporary data leaked")
	}
}

func TestIntegrationLockRecoveryAndHealth(t *testing.T) {
	c, pool, cfg := integrationDB(t)
	ctx := t.Context()
	store := &PGStore{Pool: pool}
	locker := &PGLocker{Pool: pool}
	makeRun := func() *Run {
		return &Run{ID: uuid.NewString(), Type: "daily", Status: "queued", Database: cfg.DBName, Deployment: "test", Worker: uuid.NewString()}
	}
	active, stale, fresh := makeRun(), makeRun(), makeRun()
	for _, r := range []*Run{active, stale, fresh} {
		if err := store.Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `UPDATE backup_runs SET heartbeat_at=clock_timestamp()-interval '3 minutes' WHERE id=$1`, stale.ID); err != nil {
		t.Fatal(err)
	}
	lease, ok, err := locker.Try(ctx)
	if err != nil || !ok {
		t.Fatalf("lock: %v", err)
	}
	defer lease.Close()
	if other, ok, err := locker.Try(ctx); err != nil || ok {
		if other != nil {
			other.Close()
		}
		t.Fatalf("nonblocking lock overlap: %v %v", ok, err)
	}
	if err := lease.Pulse(ctx, active); err != nil {
		t.Fatal(err)
	}
	if _, err := lease.Recover(ctx, active); err != nil {
		t.Fatal(err)
	}
	a, _ := store.Get(ctx, stale.ID)
	b, _ := store.Get(ctx, fresh.ID)
	if a.Status != "interrupted" || b.Status != "queued" {
		t.Fatal("recovery touched live queued work or missed stale work")
	}
	// The same session remains connected, but losing its advisory lock fails Pulse.
	pg := lease.(*pgLease)
	if _, err := pg.conn.Exec(ctx, `SELECT pg_advisory_unlock($1,$2)`, lockClass, lockObject); err != nil {
		t.Fatal(err)
	}
	if lease.Pulse(ctx, active) == nil {
		t.Fatal("lost session lock was not detected")
	}
	local := NewPostgres(cfg, pool)
	local.command = containerCommand(t, c)
	dir := t.TempDir()
	r := makeRun()
	if err := local.Preflight(ctx, r, dir); err != nil {
		t.Fatal(err)
	}
	if r.Checks.PageChecksums != "disabled" {
		t.Fatal("expected disabled checksums to continue")
	}
	r.Type = "weekly"
	if local.Preflight(ctx, r, dir) == nil {
		t.Fatal("weekly accepted missing amcheck")
	}
	r.Type = "daily"
	local.freeSpace = func(string) (uint64, error) { return 0, nil }
	if local.Preflight(ctx, r, dir) == nil {
		t.Fatal("accepted full disk")
	}
	local.freeSpace = func(string) (uint64, error) { return ^uint64(0), nil }
	local.command = func(_ context.Context, _ []string, _ string, _ []string, out io.Writer) error {
		_, err := io.WriteString(out, "pg_dump (PostgreSQL) 16.0")
		return err
	}
	if local.Preflight(ctx, r, dir) == nil {
		t.Fatal("accepted client mismatch")
	}
	local.command = containerCommand(t, c)
	if err := os.WriteFile(filepath.Join(dir, "database.dump"), []byte("broken archive"), 0600); err != nil {
		t.Fatal(err)
	}
	if local.Validate(ctx, dir) == nil {
		t.Fatal("accepted malformed PostgreSQL archive")
	}
}
