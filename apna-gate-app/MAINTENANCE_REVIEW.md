# Resident Maintenance — phases 1–6 review

## Delivered

Resident Home, Maintenance and Profile tabs retain the central Invite action and its permission checks. Existing Home header and Guard behavior are preserved. Announcements remains deferred until the Maintenance review, as requested in the master specification.

Maintenance uses the existing cream, navy and orange theme with a residence switcher, current-due card, server-provided due countdown, financial summary cards, and a compact bill table. Month search uses `YYYY-MM`; status filters include unpaid, pending, paid, overdue, rejected and outstanding. Summary cards also select the corresponding table filter.

Both pagination modes use backend pages of 20 records:

- Infinite scroll appends successive pages, merges overlapping IDs, and includes an explicit Load More/retry action.
- Pages replaces the visible page and provides numbered page buttons plus Previous/Next.
- Residence, month, status and mode changes reset history. Generation checks discard late responses; data is labeled with its query scope so previous results cannot appear under new filters.
- Pull-to-refresh and screen-focus refresh reload authoritative state. There is no polling or background traversal of financial history.

Current due, outstanding, overdue, total paid, effective bill status, due countdown, and page counts come from the backend. The frontend only formats backend amounts/dates and renders actions. It never derives financial totals or overdue status from fetched pages or the device calendar.

Bill Details provides the actual charge breakdown and authenticated on-demand invoice. Payment creates a server UPI request, shows authenticated QR/open-app/copy actions, and requires reference and editable payment date. Returning from UPI never implies success. Submission snapshots preserve payload/key on uncertain retries; deliberate resubmission creates a new attempt. Pending/rejected/verified states and verified, non-reversed receipts use the same presentation mapper and document handler. The AppState listener exists only on the focused payment screen.

## Changed files

App:

- Navigation: `src/app/(resident)/resident/(tabs)/_layout.tsx`, `maintenance.tsx`, `maintenance/bills/[billId]/index.tsx`, `pay.tsx`; `src/features/resident/components/resident-tab-bar.tsx`, `src/features/resident/resident-routes.ts`, `src/components/layout/notched-tab-bar.tsx`.
- Maintenance: `src/features/resident/maintenance/maintenance-screen.tsx`, `maintenance-components.tsx`, `use-maintenance-history.ts`, `bill-details-screen.tsx`, `bill-route.tsx`, `payment-screen.tsx`, `payment-status-card.tsx`, `maintenance-presentation.ts`, `payment-attempt.ts`, `maintenance.test.ts`.
- API/documents: `src/lib/api/maintenance-api.ts`, `document-api.ts`, `src/features/shared/document-action.tsx`, `open-document.ts`, `open-document.web.ts`, `src/components/ui/info-row.tsx`.
- Existing date formatting adds an optional timezone argument in `src/features/visitors/visitor-utils.ts`. Package/lock files include SDK-compatible Clipboard/Sharing and the Maintenance test command.

Backend follow-up for backend-only calculations, pagination and resident status access:

- `internal/models/maintenance.go`, new `maintenance_presentation.go`, `maintenance_presentation_test.go`.
- `internal/handlers/v1/maintenance_handler.go`, `maintenance_payment_handler.go`, `maintenance_payment_handler_test.go`.
- `internal/services/maintenanceSvc/maintenance_service.go`, `payment_queries.go`, `admin_extensions_test.go`, new `resident_maintenance_integration_test.go`.
- `internal/repositories/maintenance_repository.go`, `maintenance_readiness_repository.go`, `queries/maintenance.sql`, and its generated `internal/db/maintenance.sql.go`.

Existing unrelated app/backend work in the working tree was preserved. No schema migration, backend table, dependency upgrade, announcement feature or notification feature was introduced by this work.

## Reuse and new abstractions

Reused ResidentProvider, residence switch sheet, Resident Stack/Tabs, screen/subscreen shells, central route helpers, NotchedTabBar, enhanced RTK API/session refresh, error normalization, Card/Button/Badge/Input/AppText/EmptyState/LoadingState, FilterChip, SegmentTabs, PaginatedList, Row/Stack, theme tokens, date formatting, pagination item merging, and toast feedback.

New shared primitives are InfoRow and DocumentAction with native/web opening adapters because no authenticated PDF equivalent existed. They use the existing authenticated base query; native uses temporary OS cache files and web uses document blob download. They do not generate PDFs or persist monthly invoices.

Feature-local logic centralizes bill presentation, immutable submission attempts and scoped history state. No new API client, residence store or general pagination framework was created. Backend presentation centralizes timezone-aware calendar/status calculation and respects settlement/reversal projections.

## API contracts

Paths below are relative to the existing API base URL, under `/v1/societies/{societyId}/maintenance`:

| Method | Path | Use |
| --- | --- | --- |
| GET | `/my/bills?flat_id=…&page=…&limit=20&billing_month=…&payment_status=…` | Server-ordered, filtered history with total_count/page/total_pages/has_more |
| GET | `/flats/{flatId}/outstanding` | Current month/due, outstanding, overdue, total paid and current_bill |
| GET | `/my/bills/{billId}` | Bill, line items, authoritative display_status/due_message |
| GET | `/my/bills/{billId}/invoice` | Authenticated on-demand invoice PDF |
| POST | `/my/bills/{billId}/payment-request` | Backend UPI amount/destination/URI/QR |
| GET | `/my/payment-requests/{requestId}/qr` | Authenticated QR PNG |
| POST | `/my/bills/{billId}/payment-claims` | Request ID, reference, payment date; stable attempt key |
| GET | `/my/payment-claims?bill_id=…&flat_id=…&billing_month=…&limit=100` | Authenticated resident's own claims |
| GET | `/my/payments?bill_id=…&flat_id=…&billing_month=…&limit=100` | Selected bill's accessible payments, with resident privacy redaction |
| GET | `/my/payments/{paymentId}/receipt/pdf` | Receipt for verified, non-reversed payment |

The existing outstanding response was extended with `overdue_paise`, `total_paid_paise` and `current_bill`. Existing bill responses include `display_status`, `due_message` and verified `paid_on`; list responses support optional page mode alongside legacy cursor mode. Existing server aggregates and balance views remain the financial source of truth. Page/cursor cannot be combined. `payment_status` filters effective display state; the existing `status` filter retains its previous meaning. Ordering remains descending bill ID.

Resident claim/payment requests can now narrow by bill, flat and month. The authenticated user and active ownership checks remain mandatory; supplying a filter does not grant access. This resolves the reported “These filters are restricted to admins” failure for resident payment status. Admin endpoints and admin-only reference searches remain separate.

**Deploy the updated backend with the app:** the old server lacks the required summary/page metadata and rejects the resident bill/flat/month filters. Typed extensions remain on enhancedApi because the app's generated schema predates Maintenance; unrelated generated app endpoints were not regenerated.

## Remaining contract limitations

1. A resident sees only claims they submitted. Another resident's effective pending claim can block duplicate submission through the bill projection without exposing that person's reference/rejection reason. Visible details are labeled “Your latest submission”.
2. Bill-specific claims/payments are capped at 100 per response. If a bill has more records, details explicitly report incomplete context; no unlimited traversal or speculative receipt access occurs.
3. No Due Soon threshold is defined. The UI renders Unpaid and the backend's timezone-aware countdown instead of inventing an interval or late-fee promise.

The former missing-summary and forbidden-resident-filter incompatibilities have been addressed by the targeted backend extensions above.

## Validation

- `npm run typecheck` passed.
- `npm run lint` passed.
- `npm run test:unit`: 125 existing tests passed, including Resident/Guard navigation and permission behavior.
- `npm run test:maintenance`: 3 focused tests passed for authoritative status, deterministic payment context/reversals, and payment validation/idempotency snapshots.
- `go test ./internal/models ./internal/handlers/v1 ./internal/services/maintenanceSvc ./internal/repositories` passed.
- PostgreSQL/Docker integration test `TestResidentMaintenancePagesSummaryAndOwnershipIntegration` passed: server page counts, cross-resident isolation, scoped own claims, pending filters, receipt privacy, verified summary/current bill and reversal exclusion.
- `expo export --platform all` passed for Android, iOS and web production bundles/routes.
- Whitespace/diff checks passed.

No authenticated device/end-to-end validation was performed. A browser fixture preview server was prepared, but browser inventory was empty, so visual QA could not be completed in this session. Before release, test actual small-screen layout, residence switching during requests, UPI external return/retry, live admin verification/rejection, and invoice/receipt opening on supported devices. Native clients need rebuilding for the Clipboard/Sharing modules added in the initial Maintenance implementation.

Stopped after Maintenance (phase 6). Announcements, notifications, Needs Attention and administration remain outside this handoff.
