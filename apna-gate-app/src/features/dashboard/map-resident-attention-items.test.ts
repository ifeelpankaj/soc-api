import assert from "node:assert/strict";
import test from "node:test";

import { mapResidentAttentionItems } from "./map-resident-attention-items";
import type { MaintenanceBill } from "@/lib/api/maintenance-api";

const noop = () => {};

function bill(overrides: Partial<MaintenanceBill>): MaintenanceBill {
  return {
    id: 1,
    society_id: 1,
    flat_id: 1,
    bill_number: "B-1",
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
    ...overrides,
  };
}

const baseInput = {
  pendingCount: 0,
  onReviewVisitors: noop,
  onOpenMaintenanceBill: noop,
  onOpenMaintenancePay: noop,
  onOpenAnnouncement: noop,
};

test("mapResidentAttentionItems skips pending_review maintenance", () => {
  const items = mapResidentAttentionItems({
    ...baseInput,
    currentBill: bill({ payment_claim_status: "pending" }),
  });
  assert.equal(items.length, 0);
});

test("mapResidentAttentionItems prioritizes visitor then maintenance", () => {
  const items = mapResidentAttentionItems({
    ...baseInput,
    pendingCount: 2,
    currentBill: bill({ display_status: "overdue", status: "overdue" }),
  });
  assert.equal(items.length, 2);
  assert.match(items[0].title, /visitors awaiting approval/);
  assert.equal(items[1].title, "Maintenance overdue");
});

test("mapResidentAttentionItems returns empty when nothing needs action", () => {
  assert.deepEqual(
    mapResidentAttentionItems({
      ...baseInput,
      pendingCount: 0,
      currentBill: undefined,
      importantAnnouncement: null,
    }),
    [],
  );
});

test("mapResidentAttentionItems adds important announcement", () => {
  const items = mapResidentAttentionItems({
    ...baseInput,
    importantAnnouncement: {
      channelId: 3,
      post: {
        id: 9,
        post_id: 9,
        body: "Tomorrow · 10 AM–2 PM",
        created_at: "2026-10-01T00:00:00Z",
        is_important: true,
        title: "Water Supply Maintenance",
      },
    },
  });
  assert.equal(items.length, 1);
  assert.equal(items[0].title, "Water Supply Maintenance");
});
