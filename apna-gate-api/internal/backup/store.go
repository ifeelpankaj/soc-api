package backup

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Checks struct {
	Database      string `json:"database"`
	PageChecksums string `json:"page_checksums"`
	Amcheck       string `json:"amcheck"`
	Archive       string `json:"archive"`
}
type Retention struct {
	Status  string `json:"status"`
	Deleted int    `json:"deleted_count"`
}
type Details struct {
	Checks       Checks    `json:"checks"`
	ArchiveSize  int64     `json:"archive_size_bytes"`
	SHA256       string    `json:"archive_sha256,omitempty"`
	LocalMD5     string    `json:"local_md5,omitempty"`
	DriveFileID  string    `json:"drive_file_id,omitempty"`
	DriveMD5     string    `json:"drive_md5,omitempty"`
	Retention    Retention `json:"retention"`
	ErrorCode    string    `json:"error_code,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
}
type Run struct {
	ID         string     `json:"run_id"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Database   string     `json:"database_name"`
	Deployment string     `json:"-"`
	Worker     string     `json:"-"`
	Heartbeat  time.Time  `json:"-"`
	Created    time.Time  `json:"created_at"`
	Started    *time.Time `json:"started_at"`
	Completed  *time.Time `json:"completed_at"`
	Details
}

var ErrNotFound = errors.New("backup run not found")

type Store interface {
	Create(context.Context, *Run) error
	Save(context.Context, *Run) error
	Get(context.Context, string) (*Run, error)
}
type PGStore struct{ Pool *pgxpool.Pool }

func (s *PGStore) Create(ctx context.Context, r *Run) error {
	b, err := json.Marshal(r.Details)
	if err != nil {
		return err
	}
	return s.Pool.QueryRow(ctx, `INSERT INTO backup_runs(id,backup_type,status,database_name,deployment,worker_id,details)
	 VALUES($1,$2,'queued',$3,$4,$5,$6) RETURNING created_at,heartbeat_at`, r.ID, r.Type, r.Database, r.Deployment, r.Worker, b).Scan(&r.Created, &r.Heartbeat)
}
func (s *PGStore) Save(ctx context.Context, r *Run) error {
	b, err := json.Marshal(r.Details)
	if err != nil {
		return err
	}
	tag, err := s.Pool.Exec(ctx, `UPDATE backup_runs SET status=$3,started_at=$4,completed_at=$5,details=$6
	 WHERE id=$1 AND worker_id=$2 AND completed_at IS NULL`, r.ID, r.Worker, r.Status, r.Started, r.Completed, b)
	if err == nil && tag.RowsAffected() != 1 {
		return errors.New("backup run ownership lost")
	}
	return err
}
func (s *PGStore) Get(ctx context.Context, id string) (*Run, error) {
	r := new(Run)
	var b []byte
	err := s.Pool.QueryRow(ctx, `SELECT id::text,backup_type,status,database_name,deployment,worker_id::text,heartbeat_at,created_at,started_at,completed_at,details FROM backup_runs WHERE id=$1`, id).Scan(&r.ID, &r.Type, &r.Status, &r.Database, &r.Deployment, &r.Worker, &r.Heartbeat, &r.Created, &r.Started, &r.Completed, &b)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &r.Details); err != nil {
		return nil, err
	}
	return r, nil
}
