package backup

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Two positive int32 keys allow exact pg_locks ownership checks on the same session.
const lockClass = 1095781966
const lockObject = 1111573323

type Lease interface {
	Pulse(context.Context, *Run) error
	Recover(context.Context, *Run) ([]string, error)
	Close()
}
type Locker interface {
	Try(context.Context) (Lease, bool, error)
}
type PGLocker struct{ Pool *pgxpool.Pool }
type pgLease struct{ conn *pgx.Conn }

func (l *PGLocker) Try(ctx context.Context) (Lease, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// A standalone connection never returns session-level locks to the application pool.
	conn, err := pgx.ConnectConfig(ctx, l.Pool.Config().ConnConfig.Copy())
	if err != nil {
		return nil, false, err
	}
	var ok bool
	if err = conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1,$2)`, lockClass, lockObject).Scan(&ok); err != nil || !ok {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = conn.Close(closeCtx)
		return nil, false, err
	}
	return &pgLease{conn: conn}, true, nil
}
func (l *pgLease) Pulse(ctx context.Context, r *Run) error {
	tag, err := l.conn.Exec(ctx, `UPDATE backup_runs SET heartbeat_at=clock_timestamp() WHERE id=$1 AND worker_id=$2 AND completed_at IS NULL
	 AND EXISTS(SELECT 1 FROM pg_locks WHERE locktype='advisory' AND pid=pg_backend_pid() AND classid=$3 AND objid=$4 AND objsubid=2 AND granted)`, r.ID, r.Worker, lockClass, lockObject)
	if err == nil && tag.RowsAffected() != 1 {
		return errors.New("backup lease lost")
	}
	return err
}
func (l *pgLease) Recover(ctx context.Context, r *Run) ([]string, error) {
	// The global backup lock proves no other backup worker may execute now.
	// A queued worker can still be starting; require an expired heartbeat as well.
	_, err := l.conn.Exec(ctx, `UPDATE backup_runs SET status='interrupted',completed_at=clock_timestamp(),details=details || '{"error_code":"worker_abandoned","error_message":"Worker stopped before completing the backup"}'::jsonb
	 WHERE id<>$1 AND database_name=$2 AND completed_at IS NULL AND heartbeat_at<clock_timestamp()-interval '2 minutes'`, r.ID, r.Database)
	if err != nil {
		return nil, err
	}
	rows, err := l.conn.Query(ctx, `SELECT id::text FROM backup_runs WHERE completed_at IS NOT NULL AND database_name=$1 AND deployment=$2`, r.Database, r.Deployment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
func (l *pgLease) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = l.conn.Close(ctx)
}
