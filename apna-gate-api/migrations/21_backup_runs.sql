-- +migrate Up
CREATE TABLE backup_runs (
 id UUID PRIMARY KEY,
 backup_type TEXT NOT NULL CHECK (backup_type IN ('daily','weekly')),
 status TEXT NOT NULL CHECK (status IN ('queued','preflight','checking','dumping','validating','uploading','verifying','retention','completed','failed','interrupted','skipped_concurrent')),
 database_name TEXT NOT NULL,
 deployment TEXT NOT NULL,
 worker_id UUID NOT NULL,
 heartbeat_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
 created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
 started_at TIMESTAMPTZ,
 completed_at TIMESTAMPTZ,
 details JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX backup_runs_active ON backup_runs (heartbeat_at) WHERE completed_at IS NULL;
CREATE INDEX backup_runs_created ON backup_runs (created_at DESC);

-- +migrate Down
DROP TABLE backup_runs;
