# Society UPI payments and manual verification

Apply migrations **26–28** after migration 25 before deploying this API. No society
starts with UPI enabled. Billing settings and payment settings are independent.
This module sends no money and integrates no gateway or bank reconciliation API.
Only an owner/admin's explicit verification of a bank credit creates a collection.
Claims, QR scans, app callbacks, and reconciliation reports do not change balances.

All paths below are relative to `/api/v1/societies/{societyId}/maintenance`.
Use the existing bearer authentication, operational society guards, and JSON
`{success,message,data}` response envelope. Admin means an active society owner
or admin, never staff. Resident operations require current active access to the
bill's flat. History is retained when a resident leaves, but their access ends.

## Retry contract

The admin web release adds optional `evidence_reference` to verification and
direct-payment requests: trimmed text of at most 500 Unicode characters. It is
stored immutably with the payment and audit event, visible only in admin payment
responses, and excluded from resident responses and all PDFs. Omission preserves
the pre-extension request hash for retries. Existing payments default to no
evidence; do not fabricate or backfill historical references.

Admin claim/payment lists additionally accept `bill_id`, `flat_id`, and
`billing_month` (`YYYY-MM`). Filters combine with status/reference and cursor
pagination. Admin rows include `bill_number` and the frozen `flat` identity.
`GET /reports/summary` accepts optional `flat_id` with the existing month filter;
the summary remains calculated from server balances, including reversals.

`GET /billing-runs?billing_month=YYYY-MM` returns `status: not_generated`, or
`status: completed` with `run_id` and `bill_count`. Completed zero-bill runs remain
completed. These are read-only admin endpoints and require no retry key.

Supply `Idempotency-Key: <unique-client-generated-key>` for settings updates,
claim submission/cancellation/review, payment recording/reversal, and report
creation/updates. Keys are 1-128 bytes without surrounding whitespace. Keep the
same key when retrying after a timeout. Keys bind society, actor, operation
(including path resource), and normalized request contents. Identical retries
return the original successful response; changed contents return
`IDEMPOTENCY_CONFLICT`. Failed transactions consume no key. A replay still
requires current authorization. Read the resource again for its current status
because a replay returns the original result, even after a later reversal.

Payment-request creation does not need a key: one active request per bill is
reused until superseded or closed. No idempotency records are automatically
expired in this release.

The default CORS headers include `Idempotency-Key`. Deployments overriding
`ALLOW_HEADERS` must also include it for browser admin clients.

## Configure and verify the destination

`GET /payment-settings` returns disabled defaults, version 0, until configured.
`PUT /payment-settings` accepts:

```json
{
  "payment_methods": {
    "upi": {
      "enabled": false,
      "upi_id": "societyname@bank",
      "payee_name": "ABC Housing Society"
    }
  }
}
```

Responses include the saved version, actor, and timestamp. Each actual change
creates an immutable version and audit event and supersedes active requests.
Identical configuration is a no-op. `GET /payment-settings/versions` returns
version history. UPI validation checks syntax only; it cannot verify merchant
status or bank-account ownership. Administrators must supply the society's
bank-issued merchant details. Historical destinations remain available after a
change or disablement.

Keep live UPI disabled until the society validates its actual merchant setup:
scan a controlled QR for the configured URI and verify the displayed payee and
exact amount in the UPI apps it supports. Follow the acquiring bank's merchant
onboarding requirements. App behavior may vary; scanning or callback success is
never confirmation. Enable through the same PUT only after this check.

## Resident request and claim

1. `POST /my/bills/{id}/payment-request` returns `id`, `reference`, `state`,
   `amount_paise`, `settings_version`, `destination`, `upi_uri`, and `qr_url`.
2. Fetch `qr_url` with bearer authentication. It returns a locally generated
   PNG with a white border and `Cache-Control: private, no-store`.
3. Resident transfers to the society, then calls
   `POST /my/bills/{id}/payment-claims`:

```json
{
  "payment_request_id": "uuid-from-payment-request",
  "reference": "001234567890",
  "payment_date": "2026-09-28"
}
```

The request snapshots the bill's full amount and immutable destination version.
The server encodes `pa`, `pn`, `am`, `cu=INR`, `tr`, and `tn` with `net/url`.
Amounts use integer paise, formatted with exactly two decimal places. `tr` is
Apna Gate's request reference, **not** the bank UTR. The claim derives its amount
and destination from the saved request; clients cannot override them.

References are strings of 6-64 letters, digits, slashes or hyphens; trim and
uppercase normalization preserves leading zeros. Dates are `YYYY-MM-DD` and
cannot be in the future in the bill's stored timezone. One pending claim per
bill is allowed. Pending claims reserve the normalized reference within the
society, but do not settle the bill or change its outstanding amount.

`GET /my/payment-claims` lists only the caller's claims at currently accessible
flats. `POST /my/payment-claims/{id}/cancel` accepts `{}` or an optional
`{"reason":"..."}` and releases a pending claim's reservation. Cancellation,
rejection, and resubmission preserve the original records.

Request lifecycle:

| State | Behavior |
| --- | --- |
| `active` | Reusable while UPI is enabled, bill is unsettled, and no claim is pending. |
| `superseded` | QR no longer served; an earlier real transfer can still be claimed against the saved destination. |
| `closed` | All requests close on settlement; ordinary claims are blocked. Use reconciliation reports for transfers against these requests. |

Reversal leaves old requests closed; the next request has a fresh reference.
Disabling UPI stops new requests and QR serving but does not block verification
of transfers already made. **Saved QR images can still send money outside Apna
Gate.** The backend cannot revoke the underlying UPI destination.

## Admin verification and direct recording

Use `GET /payment-claims?status=pending` and `GET /payment-claims/{id}` to review
claims. Detail includes the saved amount/destination and previous claims and
payments for the reference, including reversed payments and reversal reasons.
`GET /payment-reference-history?reference=001234567890` also retrieves history.
If a history page has `next_cursor`, continue through `/payment-claims` or
`/payments` with the same reference and that cursor.

Check the actual bank statement for the saved destination, full credited
amount, date, and reference. Then `POST /payment-claims/{id}/verify`:

```json
{
  "bank_credit_confirmed": true,
  "amount_paise": 200000,
  "credit_date": "2026-09-28",
  "reference": "001234567890",
  "settings_version": 1,
  "upi_id": "societyname@bank",
  "reversal_history_acknowledged": false
}
```

Amount must exactly equal the full bill. Reference and version must match the
claim, and UPI ID must match that historical version, not today's configuration.
Prior reversed payments require review and explicit
`reversal_history_acknowledged: true`. A disabled current configuration does not
invalidate a verified historical credit. Reject a pending claim with
`POST /payment-claims/{id}/reject` and `{"reason":"No matching bank credit"}`.

`POST /payments/manual` accepts the same verification fields plus `bill_id`.
Without a pending claim it records a standalone verified full-bill collection.
With a pending claim it also requires `pending_claim_id` and one of:

- `claim_resolution: "verify"` for the matching reference and destination.
- `claim_resolution: "reject"`, a required `reason`, and a **distinct** verified
  reference. Claim rejection and recording happen in the same transaction.

There is no unverified manual payment status. A receipt is issued only after
the verification transaction commits.

## Receipts, balances, and reversal

Admins: `GET /payments`, `GET /payments/{id}`, `GET /payments/{id}/receipt`.
Residents: `GET /my/payments`, `GET /my/payments/{id}/receipt`.

Receipts are JSON with a stable `MPR-<uuid>` number, bill reference, integer
amount, currency, destination/version, credit date, verification timestamp and
status. Current residents see the accessible flat's complete payment history;
another payer's reference, claim ID, payer ID, and admin personal/review details
are redacted. Resident reference-search filters are rejected to prevent hidden
reference discovery. Admins retain complete records. Direct payments without a
resident claim have no attributed payer, so their bank references are admin-only.

The same verified payment can be downloaded as a PDF using
`GET /payments/{id}/receipt/pdf` (admin) or
`GET /my/payments/{id}/receipt/pdf` (resident). Existing `/receipt` routes stay
JSON. See [maintenance documents](maintenance-documents.md) for invoice downloads,
PDF privacy, historical issuer snapshots and visual validation.

`POST /payments/{id}/reverse` requires `{"reason":"Incorrect bank match"}`.
It records the responsible admin/time/reason, adds a negative compensating
ledger entry, releases the reference, and reopens the bill's balance. The
original payment and collection entry remain immutable. Reversed receipts keep
their number and display `status: "reversed"` and `reversed_at`. **Reversal is
an accounting correction; it does not refund or move money.** Only one reversal
is allowed, apart from identical idempotent retries.

Bill responses include `paid_amount_paise`, `outstanding_amount_paise`, and the
latest `payment_claim_status` (`none`, `pending`, `verified`, `rejected`, or
`cancelled`). Claim status describes claim history, not the current payment
balance; a verified claim remains historical evidence after payment reversal.
`paid` requires an active verified settlement. Pending claims leave bills
`unpaid` or `overdue`. Issued bill totals, items, and billed-party snapshots are
unchanged by payment actions.

## Duplicate-transfer reconciliation

A second real transfer must not be silently discarded or applied as an advance.
`BILL_ALREADY_PAID` directs clients to `POST /my/payment-reports`:

```json
{
  "bill_id": 123,
  "reference": "009876543210",
  "amount_paise": 200000,
  "payment_date": "2026-09-28",
  "payment_request_id": "optional-saved-request-uuid",
  "explanation": "Another member of our flat had already paid."
}
```

Omit `payment_request_id` if unavailable. Residents can list their reports via
`GET /my/payment-reports`; admins can create reports using `POST /payment-reports`.
Identical normalized submissions by the same reporter are deduplicated even
with different retry keys. Reports can be made against accessible bills whether
settled or unsettled, including closed saved requests.

Admin queue/detail: `GET /payment-reports?status=open`,
`GET /payment-reports/{id}`. Detail includes matching claims/payments.
`PATCH /payment-reports/{id}` with `{"status":"investigating"}` records the
review action. Resolve using `{"status":"resolved","resolution_note":"..."}`
after documenting the society's manual resolution. Resolved reports cannot be
changed. `GET /payment-audit?bill_id=123` exposes immutable action history and
responsible actors/timestamps, including report creation and transitions.

Reports never reserve references, create collections, credits, advances or
refunds, or affect summary totals. Actual duplicate money movement is resolved
by the society and its bank outside this module.

## Summaries, pagination, notifications, and integrity

`GET /reports/summary?billing_month=2026-09` returns billed count/paise, active
verified collections, outstanding, overdue balances and pending-claim count.
Omit month for all bills. This groups by **bill month**, not bank credit date.
List endpoints accept `limit` (1-100, default 25), `cursor` (exclusive before ID;
version for settings history), and applicable status filters. Responses contain
`items` and optional `next_cursor`. Admin claim/payment lists accept `reference`.

All financial commands acquire the existing society advisory lock, then the
bill row lock. Shared reference reservations prevent concurrent pending claims
or active payments reusing a society's reference. Composite foreign keys enforce
tenant-consistent links; partial unique indexes enforce one pending claim, one
active request, and one active settlement per bill. Settings versions, ledger
and audit records are append-only; payment/request/claim triggers protect
snapshots and one-way state transitions.

Payment, allocation (one payment's bill link), ledger, claim resolution,
reservation, request closure, audit, idempotency response, and notification
enqueueing commit together. Any failure rolls back all changes. Ledger entries
are positive collections and negative reversals; no bill amounts are mutated.

The existing `maintenance_notification_deliveries` queue handles submissions,
verification, rejection/cancellation, reversal, and reconciliation transitions.
It uses stable per-event keys and access checks at delivery. The existing
maintenance job drains it, so retain the cron setup in [maintenance.md](maintenance.md).
Push failures retry independently and never roll back committed payments.
Monitor job failures and unfinished deliveries; do not mark a payment unverified
because notification delivery failed.

## Client error handling

| Code | Action |
| --- | --- |
| `MAINTENANCE_INVALID` (400) | Correct fields, dates, exact amount/destination, or missing idempotency key. |
| `PAYMENT_FORBIDDEN` (403) | Active owner/admin permission is required. |
| `PAYMENT_NOT_FOUND` (404) | Resource is missing or outside current tenant/flat access. |
| `UPI_DISABLED` (409) | Do not initiate a new transfer. |
| `BILL_ALREADY_PAID` (409) | A real additional transfer belongs in reconciliation reports. |
| `BILL_NO_OUTSTANDING` (409) | Bill has no positive amount to collect. |
| `PAYMENT_CLAIM_PENDING` (409) | Await review or cancel the caller's pending claim. |
| `PAYMENT_REQUEST_OBSOLETE` (409) | Stop displaying the QR and reload bill/payment state. |
| `PAYMENT_REQUEST_CLOSED` (409) | Report a transfer against this request for reconciliation. |
| `PAYMENT_CONFLICT` (409) | Reference/operation is reserved; inspect existing history. |
| `IDEMPOTENCY_CONFLICT` (409) | Do not change an existing retry key's payload. |
| `CLAIM_ALREADY_REVIEWED` (409) | Refresh the claim. |
| `PENDING_CLAIM_RESOLUTION_REQUIRED`, `PENDING_CLAIM_CHANGED` (409) | Refresh and explicitly resolve the pending claim. |
| `REVERSAL_HISTORY_ACKNOWLEDGMENT_REQUIRED` (409) | Review prior reversals before acknowledging and verifying. |
| `PAYMENT_ALREADY_REVERSED`, `REPORT_RESOLVED` (409) | Read the final record; do not repeat the action with a new key. |

## Validation and release

```powershell
sqlc generate
sqlc vet
swag init -g main.go -d cmd/server,internal/handlers/v1,internal/models,internal/server --parseInternal
go test ./...
go test -tags=integration ./...
```

Integration tests require Docker and use disposable PostgreSQL containers. They
exercise migration/schema guards, old destinations, exact amounts, PNG access,
claims/reservations, idempotency, races, rollback injection, receipt redaction,
reversal accounting, duplicate reports, and notification retries/access rechecks.
Back up production before migration and follow the normal deployment workflow.
Do not down-migrate after collecting real payment records: the down migration
removes this domain's records. Retain financial history and deploy forward fixes.

Gateways, automatic bank reconciliation, partial/advance payments, refunds,
screenshots and client screens are outside this release. Maintenance invoice
PDFs, JSON receipts and payment receipt PDFs are included.
