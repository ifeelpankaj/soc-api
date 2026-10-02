import assert from "node:assert/strict";
import test from "node:test";
import type {
  MaintenanceBill,
  MaintenanceClaim,
  MaintenancePayment,
} from "../../../lib/api/maintenance-api";
import {
  dateInZone,
  getMaintenanceBillPresentation,
  resolvePaymentContext,
  validateClaim,
} from "./maintenance-presentation";
import {
  createSubmissionAttempt,
  isDefiniteRejection,
} from "./payment-attempt";

const bill: MaintenanceBill = {
  id: 7,
  society_id: 1,
  flat_id: 2,
  bill_number: "M-7",
  billing_month: "2026-09-01",
  due_date: "2026-09-30",
  timezone: "Asia/Kolkata",
  currency: "INR",
  total_paise: 250000,
  paid_amount_paise: 0,
  outstanding_amount_paise: 250000,
  status: "unpaid",
  payment_claim_status: "none",
  issued_at: "2026-09-01T00:00:00Z",
  created_at: "2026-09-01T00:00:00Z",
  items: [],
};
const now = new Date("2026-09-30T18:30:00Z");
const present = (patch: Partial<MaintenanceBill> = {}) =>
  getMaintenanceBillPresentation({ bill: { ...bill, ...patch }, now });
test("one mapper renders backend settlement and calendar state without recalculating overdue on the device", () => {
  assert.equal(dateInZone(now, bill.timezone), "2026-10-01");
  assert.equal(
    getMaintenanceBillPresentation({
      bill,
      now: new Date("2026-09-30T18:29:59Z"),
    }).state,
    "UNPAID",
  );
  assert.equal(present().state, "UNPAID");
  assert.equal(
    present({ status: "overdue", display_status: "overdue" }).state,
    "OVERDUE",
  );
  assert.equal(present({ payment_claim_status: "pending" }).action, "details");
  assert.equal(
    present({ payment_claim_status: "rejected" }).actionLabel,
    "Resubmit Payment",
  );
  assert.equal(
    present({
      status: "paid",
      outstanding_amount_paise: 0,
      payment_claim_status: "rejected",
    }).state,
    "PAID",
  );
  // A historical verified claim does not settle a bill after payment reversal.
  assert.equal(
    present({ status: "overdue", payment_claim_status: "verified" }).state,
    "OVERDUE",
  );
});
test("context resolves by the backend's effective claim ID, excludes unrelated/reversed payments and refuses truncated context", () => {
  const claim = {
    id: 2,
    bill_id: 7,
    status: "rejected",
    reference: "ABC123",
    payment_date: "2026-09-29",
    created_at: "2026-09-29T00:00:00Z",
  } satisfies MaintenanceClaim;
  const payment = {
    id: 8,
    bill_id: 7,
    status: "verified",
    amount_paise: 250000,
    verified_at: "2026-10-01T00:00:00Z",
    credit_date: "2026-09-29",
    receipt_number: "R-8",
  } satisfies MaintenancePayment;
  const records = {
    claims: {
      items: [
        { ...claim, id: 99, bill_id: 10 },
        claim,
        { ...claim, id: 3, status: "pending" },
      ],
    },
    payments: {
      items: [
        { ...payment, id: 9, status: "reversed" },
        { ...payment, bill_id: 10 },
        payment,
      ],
    },
  };
  const context = resolvePaymentContext(7, records);
  assert.equal(context.claim?.id, 3);
  assert.deepEqual(
    context.payments.map((p) => p.id),
    [8],
  );
  assert.deepEqual(
    resolvePaymentContext(7, {
      ...records,
      claims: { ...records.claims, next_cursor: 2 },
    }),
    { payments: [], complete: false },
  );
  assert.equal(
    getMaintenanceBillPresentation({ bill, paymentContext: context, now })
      .receiptPayments.length,
    0,
  );
  assert.equal(
    getMaintenanceBillPresentation({
      bill: { ...bill, status: "paid", outstanding_amount_paise: 0 },
      paymentContext: context,
      now,
    }).receiptPayments[0].id,
    8,
  );
});
test("submission snapshots keep retry contents fixed; deliberate resubmissions get new keys", () => {
  const body = {
    payment_request_id: "request",
    reference: "ABC123",
    payment_date: "2026-09-29",
  };
  const attempt = createSubmissionAttempt(body);
  const retry = attempt;
  body.reference = "CHANGED";
  assert.equal(retry.key, attempt.key);
  assert.equal(retry.body.reference, "ABC123");
  assert.notEqual(createSubmissionAttempt(attempt.body).key, attempt.key);
  assert.equal(isDefiniteRejection({ status: "FETCH_ERROR" }), false);
  assert.equal(isDefiniteRejection({ status: 500 }), false);
  assert.equal(isDefiniteRejection({ status: 429 }), false);
  assert.equal(isDefiniteRejection({ status: 422 }), true);
  assert.equal(
    validateClaim("abc123", "2026-10-01", bill.timezone, now),
    undefined,
  );
  assert.match(
    validateClaim("ABC123", "2026-10-02", bill.timezone, now)!,
    /future/,
  );
  assert.match(
    validateClaim("ABC123", "2026-02-30", bill.timezone, now)!,
    /valid/,
  );
  assert.match(
    validateClaim("", "2026-09-29", bill.timezone, now)!,
    /Reference/,
  );
});
