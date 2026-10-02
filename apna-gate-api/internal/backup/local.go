package backup

import (
	"bytes"
	"context"
	"crypto/md5" //nolint:gosec // Drive transport checksum only; SHA-256 is the archive integrity identifier.
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shirou/gopsutil/v4/disk"
	"go-server/internal/config"
)

type Local interface {
	Preflight(context.Context, *Run, string) error
	Check(context.Context, string) error
	Dump(context.Context, string) error
	Validate(context.Context, string) error
}
type Command func(context.Context, []string, string, []string, io.Writer) error
type Postgres struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	command   Command
	freeSpace func(string) (uint64, error)
}

func NewPostgres(cfg *config.Config, pool *pgxpool.Pool) *Postgres {
	return &Postgres{cfg: cfg, pool: pool, command: runCommand, freeSpace: func(path string) (uint64, error) {
		u, e := disk.Usage(path)
		if e != nil {
			return 0, e
		}
		return u.Free, nil
	}}
}
func runCommand(ctx context.Context, env []string, name string, args []string, out io.Writer) error {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // Only fixed PostgreSQL tools and server-derived argument arrays; no shell or request inputs.
	cmd.Env = env
	cmd.Stdout = out
	cmd.Stderr = io.Discard
	cmd.WaitDelay = 5 * time.Second
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%s failed; check database permissions, tool availability and server diagnostics", name)
	}
	return nil
}
func escapePass(s string) string { return strings.NewReplacer(`\`, `\\`, ":", `\:`).Replace(s) }
func (p *Postgres) environment(dir string) ([]string, error) {
	values := []string{p.cfg.DBHost, strconv.Itoa(p.cfg.DBPort), p.cfg.DBName, p.cfg.DBUser, p.cfg.DBPassword}
	for i, v := range values {
		if strings.ContainsAny(v, "\r\n\x00") {
			return nil, errors.New("unsupported line break in database configuration")
		}
		values[i] = escapePass(v)
	}
	pass := filepath.Join(dir, "pgpass")
	if err := os.WriteFile(pass, []byte(strings.Join(values, ":")+"\n"), 0600); err != nil {
		return nil, err
	}
	if err := os.Chmod(pass, 0600); err != nil {
		return nil, err
	}
	var env []string
	for _, v := range os.Environ() {
		key, _, _ := strings.Cut(v, "=")
		if !strings.HasPrefix(strings.ToUpper(key), "PG") {
			env = append(env, v)
		}
	}
	return append(env, "PGHOST="+p.cfg.DBHost, "PGPORT="+strconv.Itoa(p.cfg.DBPort), "PGDATABASE="+p.cfg.DBName, "PGUSER="+p.cfg.DBUser, "PGPASSFILE="+pass, "PGSSLMODE="+p.cfg.DBSSLMode, "PGCONNECT_TIMEOUT=10", "PGAPPNAME=apna-gate-backup"), nil
}

var pgVersion = regexp.MustCompile(`PostgreSQL\) (\d+)\.`)

func (p *Postgres) Preflight(ctx context.Context, r *Run, dir string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var size, failures int64
	var version int
	var checksums string
	err := p.pool.QueryRow(ctx, `SELECT pg_database_size(current_database()),current_setting('server_version_num')::int,current_setting('data_checksums'),COALESCE((SELECT checksum_failures FROM pg_stat_database WHERE datname=current_database()),0)`).Scan(&size, &version, &checksums, &failures)
	if err != nil {
		return preflightFailure("database_health", "Database health query failed; check database connectivity and permission to read database size and statistics")
	}
	r.Checks.Database = "passed"
	r.Checks.PageChecksums = "disabled"
	if checksums == "on" {
		r.Checks.PageChecksums = "enabled"
	}
	if failures > 0 {
		r.Checks.Database = "failed"
		return preflightFailure("database_checksums", "Database has recorded checksum failures; investigate PostgreSQL storage diagnostics")
	}
	if size < 0 {
		return preflightFailure("database_size", "Database returned an invalid database size")
	}
	free, err := p.freeSpace(dir)
	if err != nil {
		return preflightFailure("disk_capacity", "Cannot determine free space on BACKUP_TEMP_DIR filesystem")
	}
	if free < uint64(size)+(256<<20) {
		return preflightFailure("disk_space", "Insufficient backup disk space; require database size plus 256 MiB free")
	}
	env, err := p.environment(dir)
	if err != nil {
		return preflightFailure("passfile", "Cannot prepare the PostgreSQL passfile; check database configuration and temporary directory write permissions")
	}
	tools := []string{"pg_dump", "pg_restore"}
	if r.Type == "weekly" {
		tools = append(tools, "pg_amcheck")
	}
	for _, name := range tools {
		var out bytes.Buffer
		if err := p.command(ctx, env, name, []string{"--version"}, &out); err != nil {
			return preflightFailure("tool_unavailable", name+" --version failed; install PostgreSQL 15 client tools and ensure they are executable on the API process PATH")
		}
		match := pgVersion.FindStringSubmatch(out.String())
		if len(match) != 2 {
			return preflightFailure("tool_version", "Cannot identify PostgreSQL client version for "+name)
		}
		major, _ := strconv.Atoi(match[1])
		if major != 15 || version/10000 != major {
			return preflightFailure("version_mismatch", fmt.Sprintf("Backup requires PostgreSQL 15 server and client tools; server major=%d, %s major=%d", version/10000, name, major))
		}
	}
	if r.Type == "weekly" {
		var installed bool
		if err := p.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname='amcheck')`).Scan(&installed); err != nil || !installed {
			return preflightFailure("amcheck", "Weekly backup requires a provisioned amcheck extension and permission to inspect it")
		}
	}
	return nil
}
func (p *Postgres) Check(ctx context.Context, dir string) error {
	env, err := p.environment(dir)
	if err != nil {
		return err
	}
	return p.command(ctx, env, "pg_amcheck", []string{"--jobs=1"}, io.Discard)
}
func (p *Postgres) Dump(ctx context.Context, dir string) (result error) {
	env, err := p.environment(dir)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "database.dump"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600) //nolint:gosec // Exclusive creation inside a private server-generated run directory.
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, f.Close()) }()
	if err := p.command(ctx, env, "pg_dump", []string{"--format=custom", "--compress=6", "--lock-wait-timeout=5000", "--no-password"}, f); err != nil {
		return err
	}
	return f.Sync()
}
func (p *Postgres) Validate(ctx context.Context, dir string) error {
	env, err := p.environment(dir)
	if err != nil {
		return err
	}
	return p.command(ctx, env, "pg_restore", []string{"--file=" + os.DevNull, filepath.Join(dir, "database.dump")}, io.Discard)
}
func hashArchive(ctx context.Context, path string, r *Run) error {
	f, err := os.Open(path) //nolint:gosec // Only the completed archive in the private run directory is hashed.
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	s := sha256.New()
	m := md5.New() //nolint:gosec // Required for comparison with Google Drive, not authentication.
	n, err := io.Copy(io.MultiWriter(s, m), &contextReader{ctx: ctx, r: f})
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("empty backup archive")
	}
	r.ArchiveSize = n
	r.SHA256 = hex.EncodeToString(s.Sum(nil))
	r.LocalMD5 = hex.EncodeToString(m.Sum(nil))
	return nil
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(b []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(b)
}
