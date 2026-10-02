# PostgreSQL backups to personal Google Drive

The API accepts authenticated backup requests and runs them asynchronously. It
backs up the database in `DB_NAME` using the existing `DB_*` connection settings.
It does not back up other databases, cluster-wide roles, Docker volumes, uploaded
images, or application files. There is no internal backup scheduler.

## Configure the server

Apply migration `21_backup_runs.sql` through the existing migration workflow,
then rebuild the API image. The runtime includes PostgreSQL 15 client tools and
runs the API as its existing nonroot user; it does not start another database.
Both server and clients must be major version 15. A future database major upgrade
requires updating the client image and compatibility check together.
The Docker build verifies `pg_dump`, `pg_restore`, and `pg_amcheck` are PostgreSQL
15 tools, so a missing or incompatible runtime fails during deployment instead of
waiting for the first backup request.

Add these settings to the API's server env file (for Compose, normally
`apna-gate-api/.env.production`):

```dotenv
BACKUP_ENABLED=true
BACKUP_DEPLOYMENT=production
BACKUP_TEMP_DIR=/var/lib/apna-gate/backups/tmp
BACKUP_TIMEOUT_SECONDS=7200
BACKUP_GDRIVE_FOLDER=https://drive.google.com/drive/folders/YOUR_FOLDER_ID
BACKUP_GDRIVE_CLIENT_ID=YOUR_DESKTOP_OAUTH_CLIENT_ID
BACKUP_GDRIVE_CLIENT_SECRET=YOUR_DESKTOP_OAUTH_CLIENT_SECRET
BACKUP_GDRIVE_REFRESH_TOKEN=YOUR_REFRESH_TOKEN
```

`BACKUP_GDRIVE_FOLDER` also accepts a bare folder ID. The target folder must
already exist and be writable by the authorized Google account. Keep it private.
`BACKUP_DEPLOYMENT` is a stable, unique identifier for this deployment; retain it
across releases and instances so retention recognizes older backups. Use different
folders and deployment identifiers for production and staging.

Keep `BACKUP_ENABLED=false` until configuration and database permissions are ready.
Disabled backups return HTTP 503, while other API jobs continue normally. Credentials
are validated locally at startup; actual Drive access is checked on each run.

The API user must be able to read every table, sequence, and large object in the
database and write `backup_runs`. Existing row-level security must not hide data:
`pg_dump` fails rather than silently omitting rows when full access is unavailable.
New tables must remain accessible to this role. The logical dump preserves object
ownership/ACL metadata; cluster-wide role definitions themselves are not included.

The temporary directory is dedicated to backups, with one private UUID directory
per run. The image provisions the default directory for UID 10001. A custom bind
mount must be writable by that UID (and correctly labeled on SELinux hosts).
Have free space of at least `pg_database_size(DB_NAME) + 256 MiB` there; database
growth can still exhaust space during a run, which fails safely. Dumps stream to
disk and uploads use 8 MiB chunks, without buffering the database in memory.
Archives and temporary password files are removed after each run. Startup recovery
reconciles abandoned runs after a 125-second grace period; subsequent backup
triggers also reconcile expired heartbeats under the maintenance lock.

## Authorize personal Google Drive once

1. Create/select a Google Cloud project and enable the Google Drive API.
2. Configure Google OAuth consent for your account and create a **Desktop app**
   OAuth client. The helper uses an ephemeral `127.0.0.1` callback and PKCE.
3. Set `BACKUP_GDRIVE_CLIENT_ID` and `BACKUP_GDRIVE_CLIENT_SECRET` in your local
   shell. The helper reads shell variables, not the API env file.
4. Run this from `apna-gate-api`:

   ```bash
   go run ./cmd/gdrive-auth
   ```

5. Open the printed URL, authorize the Google account owning the destination,
   and place the printed refresh token in the server env file. The helper exits
   after authorization and is never run in production.

Because the destination is an arbitrary existing folder pasted into env, the
helper requests the Drive scope. Application operations are restricted to the
configured folder, but the OAuth credential itself has broader Drive access.
Protect all three credentials; never commit them. The existing Firebase/Google
application sign-in does not supply this separate backup authorization.

Offline access permits uploads when you are absent. OAuth apps in external
**Testing** mode generally receive refresh tokens that expire after seven days
for Drive access. Configure publishing/consent for unattended use before relying
on backups. Tokens can also be revoked or invalidated later; reauthorize if that
happens. See [Google's OAuth lifecycle documentation](https://developers.google.com/identity/protocols/oauth2)
and [offline access](https://developers.google.com/identity/protocols/oauth2/web-server).

## Provision weekly corruption checks

Run `scripts/backup-amcheck.sql` once against the configured database as an admin.
Supply the API database role as the psql variable `backup_role`. For example, with
your admin `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`, and protected `PGPASSFILE` set:

```bash
psql --set=ON_ERROR_STOP=1 --set=backup_role=app_user --file=scripts/backup-amcheck.sql
```

This installs `amcheck` in `public` and grants only the checking functions used
by the weekly command. If an existing installation uses another schema, adapt
the deployment grants to that schema. The webhook never installs extensions or
grants privileges. Missing installation or permissions fail the weekly run.

Daily runs check connectivity, available disk space, matching tool/server
versions, checksum state and recorded checksum failures, and Drive access.
Disabled checksums are reported as `disabled`, with limited coverage, and do
not fail a backup. Recorded checksum failures do fail preflight: investigate
the database instead of automatically resetting the statistics.

Weekly runs add single-worker `pg_amcheck`, without aggressive parent or
heap-all-indexed options. The dump itself is a single `pg_dump -Fc -Z 6` process
with a five-second table lock-wait timeout. Ordinary reads and writes continue;
there is still database I/O and CPU cost, and schema changes can conflict with
dump locks. No minimum compression ratio is guaranteed.

Every archive is fully read through `pg_restore --file=/dev/null`, then hashed
with SHA-256 and MD5. This validates archive readability. Weekly checks cover
supported relation types; neither these checks nor an archive read prove that
all data is corruption-free or that a complete restore will succeed. Interpret
passing checks as **no corruption detected by performed checks**. See
[PostgreSQL amcheck coverage](https://www.postgresql.org/docs/15/amcheck.html).

## Trigger and inspect backups

The only permitted request body field is `type`. Unknown fields, duplicate keys,
invalid types, and trailing JSON are rejected. The request cannot select a
database, folder, credentials, or retention interval.

```bash
curl --fail-with-body -X POST 'https://api.apnagate.org/api/internal/jobs/db-backup' \
  -H "Authorization: Bearer $JOB_WEBHOOK_SECRET" \
  -H 'Content-Type: application/json' \
  --data '{"type":"daily"}'
```

For a weekly backup send `{"type":"weekly"}`. Configure your external scheduler
for two daily calls and one weekly call. There is no automatic scheduling in the
API and no dependency on the existing expiry job.

A successful admission is HTTP 202:

```json
{"run_id":"UUID","type":"daily","status":"queued"}
```

This means the request was persisted, not that a backup already exists. Poll:

```bash
curl --fail-with-body 'https://api.apnagate.org/api/internal/jobs/db-backup/RUN_ID' \
  -H "Authorization: Bearer $JOB_WEBHOOK_SECRET"
```

Progress states are `queued`, `preflight`, `checking` (weekly only), `dumping`,
`validating`, `uploading`, `verifying`, and `retention`. Terminal states are
`completed`, `failed`, `interrupted`, and `skipped_concurrent`.

Each accepted request has its own UUID and audit record. A nonblocking database
advisory lock permits one backup across daily/weekly requests and API instances;
overlaps become `skipped_concurrent`. A worker pulses its dedicated lock session
every ten seconds and cancels work if it loses ownership. A restart never
automatically retries an interrupted backup; submit a new request.

The status response includes check results, timestamps, archive byte size,
`archive_sha256`, `local_md5`, `drive_file_id`, `drive_md5`, retention status/count,
and sanitized error fields. Grafana/Loki can filter API logs by `run_id` and
`database backup`. Subprocess stderr and provider responses are intentionally
excluded from API responses and logs to avoid credential/data exposure. Consult
PostgreSQL logs and deployment permissions when diagnosing a failed stage.

Preflight failures now log and persist a specific safe `error_code` and
`error_message`, for example `preflight_tool_unavailable`,
`preflight_version_mismatch`, `preflight_disk_space`, or
`preflight_drive_folder_missing`. Failed/interrupted runs log at ERROR level.
For missing tools, ensure PostgreSQL 15 `pg_dump` and `pg_restore` (plus
`pg_amcheck` for weekly runs) are on the API process PATH. A version mismatch
reports the numeric server/client major versions. Temporary directory failures
identify permissions/setup checks, while Drive failures distinguish authentication,
access denial, missing folders, and non-writable destinations. Older run records
containing only `preflight_failed` cannot recover the discarded underlying reason;
trigger a new run after deploying diagnostics. `retention.status=pending` on a
preflight failure means upload and old-backup deletion were never reached.

Your external scheduler should poll until a terminal status and alert on failed
or interrupted runs, skipped calls that need retry, and retention warnings. A 202
response alone is not a successful-backup signal.

## Upload verification and retention

Uploads use resumable sessions, bounded retries, and an allocated Drive file ID
that remains stable across upload retries. After upload, remote size and MD5
must match the archive; its SHA-256 remains in both the run record and Drive
metadata. No public sharing permissions are created.

Only verified uploads activate retention. Daily archives expire after seven days;
weekly archives expire after 84 days, measured from Drive creation time. The
newest verified backup of each type is always retained. Expired archives are
**permanently deleted**, not moved to trash, so they no longer consume quota.

Deletion requires the configured parent folder and matching `appProperties` for
the backup system, deployment, database, and type, plus a valid run UUID,
verification timestamp, and SHA-256 metadata. Metadata is reread before deletion.
Filenames are never used as deletion authority. Unverified uploads are excluded;
review them manually if a run failed after creating a remote file.

If upload or verification fails, no old backup is deleted. If retention fails
after a verified upload, the run is `completed` with `retention.status=warning`
and its successful deletion count. The next successful run retries retention.
Changing folders/deployment identifiers intentionally stops cleanup of old scope.

## Restore and test

Download the archive using its Drive file ID. Compare its SHA-256 to the run
record (`sha256sum file.dump` on Linux, `Get-FileHash file.dump -Algorithm SHA256`
on PowerShell). Use an isolated PostgreSQL 15 instance for a restore rehearsal.
With that instance's connection variables and protected passfile set:

```bash
createdb apna_gate_restore_test
pg_restore --exit-on-error --no-owner --no-acl \
  --dbname=apna_gate_restore_test file.dump
```

`--no-owner --no-acl` makes an isolated test independent of production role names.
For an actual recovery requiring the original ownership/grants, provision those
roles separately and omit those options. Check tables, sequences, constraints,
views, and representative application queries after restore. Never point these
test commands at production. Restored `backup_runs` are historical snapshot data;
keep backups disabled in restored test environments.

Development checks:

```bash
go test ./internal/backup ./internal/config ./internal/jobs ./internal/routes ./cmd/gdrive-auth
go test -race ./internal/backup ./internal/jobs ./internal/routes
go test -tags integration ./internal/backup -run TestIntegration -count=1 -v -timeout 5m
```

The integration suite requires Docker and uses disposable PostgreSQL 15
containers. It performs an actual compressed dump, weekly check, persistent run,
restore, and schema/data checks; Google Drive is faked, so it never uploads test
data to your account. A monthly isolated restore-verification job remains future
work, separate from ordinary daily and weekly backups.
