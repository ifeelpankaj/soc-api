# Maintenance billing

This release provides configuration, all four pricing models, monthly bills,
resident history, and notifications. Migration 26 adds [society UPI payments
with manual verification](maintenance-payments.md), JSON receipts, reconciliation
reports, and audited reversals. Migration 27 adds [invoice and receipt PDF
downloads](maintenance-documents.md). The web admin supports billing and collections.
Adjustments and late fees remain outside this release. Bills expose `paid`, `unpaid`, or `overdue`, along with
`paid_amount_paise`, `outstanding_amount_paise`, and `payment_claim_status`.

## Deployment

1. Apply migrations 25–32, including `32_maintenance_reminder_payloads.sql`, with the
   existing migration workflow before deploying the API. Existing societies
   remain disabled for billing and UPI until separately configured.
2. Configure `MAINTENANCE_ELIGIBLE_FLAT_STATUSES=vacant,occupied,blocked` (default).
   This deployment-wide list controls billing for active flats. Inactive flats
   are always excluded; empty or unknown status values fail startup.
3. Configure external cron, e.g. `10 * * * *`, to POST hourly to
   `/api/internal/jobs/maintenance-billing` with
   `Authorization: Bearer <JOB_WEBHOOK_SECRET>`. It returns 202 with a trigger ID.
   The existing job logs/metrics use `maintenance-billing`. No external cron is
   automatically installed by this change. Deploy the migration before the API,
   then deploy the admin web panel.
4. Populate required flat data, save settings, preview, and generate a current-month
   batch. Inspect pending `maintenance_notification_deliveries` and job errors.

All endpoints below are relative to
`/api/v1/societies/{societyId}/maintenance` and use the standard
`{success,message,data}` envelope and existing operational-society guards.
Admin operations require active society owner/admin membership, not staff access.

## Configuration

`GET /settings` returns disabled defaults for an unconfigured society.
`PUT /settings` atomically replaces configuration. Example hybrid request:

```json
{
  "enabled": true,
  "pricing_model": "hybrid",
  "currency": "INR",
  "billing_day": 1,
  "due_day": 10,
  "timezone": "Asia/Kolkata",
  "fixed_paise": 50000,
  "area_rate_paise": 300,
  "type_rates": {}
}
```

- `fixed`: fixed paise per flat.
- `per_sqft`: area multiplied by `area_rate_paise` per square foot.
- `flat_type`: `type_rates`, e.g. `{"1 BHK":150000,"2 BHK":250000}`.
- `hybrid`: fixed plus area charge.

Money uses integer paise. Area multiplication is exact, rounding half-up once to
a paise, with overflow checks. Billing/due days are 1–28; a due day before the
billing day means next month. Currency is INR. `eligible_flat_statuses` is
read-only deployment configuration; submitted values cannot override it.

Existing flat create, bulk-create, PATCH, list and detail APIs support optional
`flat_type` and `area_sqft`, e.g. `{"flat_type":"2 BHK","area_sqft":"1000.25"}`.
Area is a positive decimal **string**, with up to two fractional digits, stored
as integer hundredths. Type keys match exactly and cannot have surrounding spaces.
Omitted/null PATCH fields retain existing values. Generated flats start without
area/type data; populate these through PATCH before area/type billing.

## Billing API and behavior

`POST /bills/preview` and `POST /bills/generate` accept
`{"billing_month":"2026-10"}` using the actual current month in the society timezone.
Manual generation may precede the billing day. Automatic generation runs on or
after that day and recovers missed triggers within the month. Historical backfill
requires the existing catch-up preview and acknowledgment for wholly unissued
past months. Future-month generation is unsupported.

Preview returns `bills`, `issues` (flat IDs/reasons) and `total_paise` for valid
preview rows. With issues, that total is not a completed batch. Previewing an
issued current month returns only missing flats, using the original run's rates,
timezone, billing/due days and eligibility policy, with current flat and resident
details. Generation returns
422 with `error.details.issues` if any eligible flat lacks valid area/type/rate
data, committing no new bills. Already issued flats are excluded from validation.

Successful generation returns `run_id`, `created`, and `existing`. Advisory locks
serialize manual and job requests; a database unique society/flat/month constraint
also prevents duplicates. A fresh manual request reconciles an issued current
month and appends bills for newly eligible flats. For 115 existing bills and five
missing flats, it returns `created: 5, existing: 115`; another fresh request returns
`created: 0, existing: 120`. The same `Idempotency-Key` replays its original response;
use a new key for a new action. Automatic billing does not append to an issued run.
An empty initial run can subsequently receive missing bills through the admin action.
Past issued months cannot be reopened. Added bills retain the original due date,
even if it has passed; the admin preview displays this before confirmation.

Pricing, eligible flat IDs, eligibility policy, flat attributes, dates, charges,
and billed-party details are snapshotted. Database triggers reject changes or
deletion of issued bills/items. Old dues remain separate monthly bills and are
never copied into new charges. The original run snapshot also stays immutable;
additional bill IDs, counts and the initiating admin are audited separately.
Disabling maintenance stops generation and keeps
authorized bill history available.

| Endpoint | Access |
|---|---|
| `GET /bills` | Society owner/admin |
| `GET /bills/{id}` | Society owner/admin, with billed-party snapshot |
| `GET /my/bills` | Current resident's accessible flats |
| `GET /my/bills/{id}` | Current resident, without historical personal snapshots |

Lists accept `billing_month`, `flat_id`, `status` (`paid`/`unpaid`/`overdue`), `limit`
(1–100, default 25), and `cursor` (previous `next_cursor`). Results have `items`
and optional `next_cursor`, ordered by descending bill ID. Overdue starts after
the due date in each bill's snapshotted timezone. Sum issued amounts for the
selected period to derive outstanding; no previous-dues line item is added.

Residents require active society membership and active residency in an active
flat. Current residents can see earlier flat bills; moved-out residents lose
access. Missing/inaccessible detail IDs return 404.

## Notifications

Generation queues one delivery per bill/current active resident in its database
transaction. The job drains up to 100 ready deliveries per trigger, rechecks
access, creates an inbox entry and then sends push. Manual generation queues
notifications for the next job trigger. Trigger the job immediately or more
frequently when delivery latency/volume requires it.

The shared outbox uses a two-minute lease, ownership token, separate inbox/push
progress, and exponential backoff up to one hour. Delivery failures cannot roll
back bills. Inbox entries use a unique
bill/user event key. Push is at-least-once and includes the event key for client
deduplication. Notifications contain bill, flat and society IDs rather than PDF attachments.

## Payment reminders

`POST /api/internal/jobs/maintenance-reminders` uses
`Authorization: Bearer <JOB_WEBHOOK_SECRET>` and returns HTTP 202 with a trigger ID.
Configure your external scheduler separately. Daily calls are sufficient for
the reminder policy; more frequent calls improve retry and queue throughput and
are safe because each bill/recipient/milestone/date has one event key.

The job checks enabled, operational societies and outstanding bills across all
billing months. It uses each bill's persisted due date and timezone, sending on
exactly three local dates: three days before, the due date, and seven days after.
Missed scheduling dates are skipped; queued failed deliveries may retry later.
Pending payment claims remain outstanding until verified.

Current active residents of the billed flat receive an inbox and push reminder
with the billing month, outstanding amount and due date. Eligibility, membership
and balance are checked before inbox creation and again before push. Paid bills,
disabled/non-operational societies, and residents who lose access are skipped.
An inbox created before payment remains historical; a later push is suppressed.

The webhook creates reminder events and drains ready outbox batches within a
five-minute job deadline. Pending work and failures remain durable for later
triggers. The billing webhook can also retry queued deliveries, but never creates
reminders. There is no manual reminder action or separate delivery webhook.

Observe job logs under `maintenance-reminders`,
`apna_gate_maintenance_reminder_events_total` (queued/skipped/failed),
`apna_gate_maintenance_reminder_deliveries_total` (completed/skipped/failed), and
the existing outbox backlog/failure metrics. Completion means the delivery
attempt finished, not that the resident has opened the notification.

## Verification

```text
sqlc generate
sqlc vet
go test ./...
go test -tags=integration ./internal/services/maintenanceSvc -v
```

Integration tests create and remove a private PostgreSQL container, never using
the configured app database. They exercise migrations, missing-data/queue-failure
rollback, concurrent retries, snapshots, access checks, pagination and notification
deduplication. Rolling this migration down removes financial records; after real
bills exist, correct schema issues with forward migrations.
