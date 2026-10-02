# Apna Gate API

Maintenance: [Billing configuration, API, and scheduled-job setup](docs/maintenance.md).
Payments: [Society UPI configuration, manual verification, and reconciliation](docs/maintenance-payments.md).
Documents: [Maintenance invoice and payment receipt PDF downloads](docs/maintenance-documents.md).

Database backups: [PostgreSQL backup webhook and Google Drive setup](docs/database-backups.md).

## Running The API

Start the API locally with an interactive env-file prompt:

```bash
make dev
```

The selected env file is passed to the app as `ENV_FILE`, while `GO_ENV` stays set to `development`.

```text
GO_ENV=development
ENV_FILE=<selected-file>
```

Use a specific env file non-interactively for CI, scripts, or repeat local runs:

```bash
make dev ENV_FILE=.env.development
make dev ENV_FILE=.env.staging
make dev ENV_FILE=.env.testing
```

`make dev` can run with any real `.env*` file in this folder. Files such as `.env.example`, templates, and samples are excluded from the prompt.

## Configuration Precedence

Environment file loading order:

```text
ENV_FILE
.env.<GO_ENV>
.env
```

Server host precedence:

```text
SERVER_HOST
HOST
default host
```

Server port precedence:

```text
SERVER_PORT
PORT
default port
```

The selected env file is also used by `kill-port`, so the process killed before startup uses the same port that the server will read.

## Runtime Mode Vs Env File

`GO_ENV` controls application behavior, such as development or production checks.

`ENV_FILE` controls actual configuration values, such as server host, server port, database settings, and secrets.

## Background Jobs

Expiry processing starts only after the HTTP port binds successfully, runs immediately, and repeats every 15 minutes. Cleanup and monthly visitor reporting are triggered only through authenticated internal webhooks.

Set a random secret containing at least 32 characters:

```text
JOB_WEBHOOK_SECRET=<random-secret-at-least-32-characters>
```

Configure the external cron provider to use the `Asia/Kolkata` timezone and send `POST` requests with `Authorization: Bearer <JOB_WEBHOOK_SECRET>`:

```text
0 2 * * *   POST https://api.example.com/api/internal/jobs/cleanup
5 * * * *   POST https://api.example.com/api/internal/jobs/monthly-visitor-report
10 * * * *  POST https://api.example.com/api/internal/jobs/maintenance-billing
0 9 * * *   POST https://api.example.com/api/internal/jobs/maintenance-reminders
```

The report trigger runs hourly so failed deliveries can retry. The job itself enforces the 00:05 first-day eligibility boundary and skips society/month deliveries whose `sent_at` is already recorded. For local testing only, `VISITOR_REPORT_ALLOW_RESEND=true` may be enabled in a development environment. Each trigger then resends the latest eligible month with a unique attempt-based provider idempotency key; earlier sent months remain skipped. Configuration validation rejects this flag in production, test, and other environments.

Example requests:

```bash
curl -X POST -H "Authorization: Bearer ${JOB_WEBHOOK_SECRET}" https://api.example.com/api/internal/jobs/cleanup
curl -X POST -H "Authorization: Bearer ${JOB_WEBHOOK_SECRET}" https://api.example.com/api/internal/jobs/monthly-visitor-report
```

Webhook responses use `202 Accepted`; this confirms that execution was accepted, not that it has completed.

Maintenance reminders use each issued bill's due date and timezone: three days
before, on the due date, and seven days after, while outstanding. Repeated
triggers are deduplicated. Failed queued deliveries can retry on later triggers;
missed reminder dates are not backfilled. External scheduler setup is separate.
Apply migration `32_maintenance_reminder_payloads.sql` before deploying the API.
See [maintenance billing and reminders](docs/maintenance.md) for admin missing-bill
generation, delivery behavior, and monitoring.

Visitor cleanup retains the existing entry retention and sent-report requirements. It also deletes a visitor when their last entry is removed, and removes old unreferenced visitors using the visitor-entry retention cutoff. Visitors referenced by any remaining entry are retained across all societies.

Migration `22_visitor_image_cleanup.sql` adds a durable ImageKit deletion queue and must be applied before deploying this cleanup version. Managed photo file IDs are queued atomically with visitor deletion. Cleanup deletes queued images in bounded batches with a five-second timeout per image; missing images count as success. Failed deletions remain queued for the next cleanup trigger, including after restarts. When ImageKit is disabled the queue is retained. Legacy photo URLs without a stored file ID cannot be deleted automatically and produce a log warning. Job errors and `visitor_image_deletions` queue size indicate pending failures.

All seven timestamp columns in newly generated monthly visitor CSV reports display Indian Standard Time, for example `10 sept 2026 17:45 IST`. Missing dates stay blank. Database timestamps, report month boundaries, and previously sent reports are unchanged.
