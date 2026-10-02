# Maintenance invoice and payment receipt PDFs

Apply `27_maintenance_documents.sql` after migrations 25 and 26, before deploying
this API. Downloads run entirely inside the Go backend using pinned gopdf and
embedded Noto Sans fonts (SIL Open Font License). There is no PDF service,
browser, filesystem storage, external font fetch, or PDF table at runtime.

## Download endpoints

All paths below are relative to `/api/v1/societies/{societyId}/maintenance`.

| Document | Owner/admin | Current resident |
| --- | --- | --- |
| Invoice PDF | `GET /bills/{id}/invoice` | `GET /my/bills/{id}/invoice` |
| Receipt PDF | `GET /payments/{id}/receipt/pdf` | `GET /my/payments/{id}/receipt/pdf` |
| Receipt JSON (unchanged) | `GET /payments/{id}/receipt` | `GET /my/payments/{id}/receipt` |

Receipt IDs refer to **payment records**, not claims or bills. A claim alone has
no receipt. Every download uses existing bearer authentication, operational
guards and tenant/flat authorization. Staff cannot use admin routes. Residents
must retain active access to the flat, including for historical documents.

Successful responses contain raw PDF bytes, not the JSON envelope:

```http
Content-Type: application/pdf
Content-Disposition: attachment; filename=maintenance-invoice-123.pdf
Cache-Control: private, no-store
X-Content-Type-Options: nosniff
```

Receipts use `maintenance-receipt-{paymentId}.pdf`. Filenames never interpolate
user-entered text. No retry key is required for these read-only GET requests.
Clients can download the authenticated response as a blob and share the saved
file themselves; the backend creates no public sharing URL. Downloads do not
change any ledger, payment, claim, QR request, or balance. UPI can be disabled.

## Document contents and privacy

The **Maintenance Invoice** shows the frozen issuer details, existing bill
number, issue date, flat/block/type/area, billing month, immutable charges,
original total and due date. It says **Not a payment receipt**. It contains no
UPI QR, live configuration, paid balance, outstanding amount or payment status.
Payment and reversal therefore do not change the invoice contents.

The **Payment Receipt** uses the existing JSON receipt's stable `MPR-...` number,
amount, credit date, verification timestamp, historical UPI destination/version
and status. It also identifies the original bill and flat. Its issuer header
comes from that bill's frozen snapshot, not current society settings.

Receipt PDFs use exactly the same permission-filtered payment data as JSON.
Other residents' bank references, payer IDs, reviewer IDs and reversal reasons
are omitted. Restricted bank references are explicitly labeled. Direct payments
with no attributed resident payer keep the reference admin-only. Historical
billed-party personal details are not printed in either document. PDF metadata
contains only generic document titles, application name and document timestamp.

Reversed receipts keep their number and show **REVERSED**, reversal time and a
notice that they no longer represent an active settlement. Reversal does not
constitute a refund or move money. Previously downloaded PDFs cannot be revoked;
download again to obtain the latest receipt status.

## Issuer history and rollout

Existing bills lacked society name/address snapshots. Migration 27 adds
`maintenance_bills.issuer_snapshot`, captures current society details once, and
marks the capture source `rollout`. Those PDFs disclose the capture date and
that historical issuer details were unavailable. Missing optional addresses or
contact fields are omitted rather than invented.

The migration runs in one transaction, holds an exclusive bill-table lock and
temporarily disables only `maintenance_bill_immutable` during the backfill.
It restores the trigger before commit. Existing bill identifiers, charges,
items, original snapshots and billed parties remain untouched. Schedule this
migration during a maintenance window if the bill table is large; its backfill
blocks bill reads/writes until the transaction completes.

A database BEFORE INSERT trigger always captures society details for new bills,
including inserts from older application instances during rollout. Callers
cannot supply an alternative issuer. The existing immutable bill trigger also
protects this snapshot against later updates/deletes. The bill JSON response
exposes it as `issuer`, including `captured_at` and `capture_source` (`issuance`
or `rollout`). Existing JSON fields and receipt endpoints retain their meanings.

Back up before applying the migration. Do not down-migrate after rollout:
dropping the snapshot loses historical issuer data. Use forward fixes.

## Rendering and error behavior

Documents use A4 pages, embedded regular/bold fonts, wrapped descriptions,
repeated table headers and page numbers. All money is formatted directly from
integer paise as `INR 2,000.00`; monetary values never pass through floating point.
Dates/timestamps use the original bill timezone.

This release supports English/Latin-script content. Unsupported scripts or
missing glyphs return `422 DOCUMENT_UNSUPPORTED_TEXT`; no replacement boxes,
transliterations or silent truncation are used. Unsupported input does not
change the saved bill or payment. `422 DOCUMENT_TOO_LARGE` limits documents to
100 pages. Other rendering failures use the standard server error response.
The complete document is buffered before sending response headers, so failures
return the normal JSON error envelope instead of a partial PDF.

## Tests and visual verification

```powershell
sqlc generate
sqlc vet
swag init -g main.go -d cmd/server,internal/handlers/v1,internal/models,internal/server --parseInternal
go test ./...
go test -tags=integration -p 1 ./...
docker build -t apna-gate-api:maintenance-documents-check .
```

Integration tests use disposable PostgreSQL containers. Document tests verify
legacy backfill without changing original bill fields, new-bill capture,
immutability, UPI independence, payment/reversal stability, access checks and
receipt redaction. Unit/handler tests cover exact int64 amounts, invalid text,
cancelled rendering, document type, attachment headers and JSON error responses.

To export synthetic PDFs for visual inspection, set `MAINTENANCE_PDF_QA_DIR` to
an **absolute ignored temporary directory**, then run:

```powershell
$env:MAINTENANCE_PDF_QA_DIR = Join-Path (Get-Location) 'tmp/pdfs'
go test ./internal/services/maintenanceSvc/documentpdf -run TestDocumentVisualFixtures -count=1
go test -tags=integration ./internal/services/maintenanceSvc -run TestMaintenanceDocumentsMigrationAndAccess -count=1
python scripts/maintenance_pdf_qa.py tmp/pdfs --render-with <path-to-pdftoppm>
```

The QA script requires `pypdf`, and rendering requires Poppler; both are
development tools only. It extracts every page, asserts money, metadata,
redaction and page numbering, and renders all pages. Inspect all rendered PNGs
for overlaps, clipping and glyph issues. Fixtures include ordinary invoices,
large amounts, a 12-page invoice and verified/redacted/reversed receipts.

Gateways, automatic bank reconciliation, partial/advance payments, refunds,
screenshots, and client screens remain outside this release. Invoice PDFs,
JSON receipts and payment receipt PDFs are included.
