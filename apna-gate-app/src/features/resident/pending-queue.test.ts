import assert from "node:assert/strict";
import test from "node:test";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { filterPendingQueue } from "./pending-queue";
import { residentInvitePurposes } from "@/features/visitors/visitor-utils";
import { invitePurposePresentation } from "./invites/invite-purpose-presentation";

test("search finds requests beyond the first visible page", () => {
  const entries: ModelsVisitorEntry[] = Array.from({ length: 35 }, (_, i) => ({
    id: i + 1,
    purpose: "guest",
    full_name: `Visitor ${i + 1}`,
  }));
  assert.deepEqual(
    filterPendingQueue(entries, "Visitor 35", "all").map((entry) => entry.id),
    [35],
  );
});
test("purpose filters preserve exact historical Service, Staff, Cab and Maintenance identities", () => {
  const entries: ModelsVisitorEntry[] = [
    "service",
    "staff",
    "cab",
    "maintenance",
    "other",
  ].map((purpose, id) => ({
    id,
    purpose: purpose as ModelsVisitorEntry["purpose"],
  }));
  for (const entry of entries)
    assert.deepEqual(filterPendingQueue(entries, "", entry.purpose!), [entry]);
});
test("search normalizes whitespace and case and includes provider information", () => {
  const entry: ModelsVisitorEntry = {
    service_provider: "Newspaper Delivery",
    purpose: "service",
  };
  assert.deepEqual(filterPendingQueue([entry], "  newspaper  ", "service"), [
    entry,
  ]);
  assert.deepEqual(filterPendingQueue([entry], "newspaper", "guest"), []);
});
test("every permitted invite purpose has descriptive presentation without expanding resident access", () => {
  assert.deepEqual(
    Object.keys(invitePurposePresentation),
    residentInvitePurposes,
  );
  for (const purpose of residentInvitePurposes)
    assert.ok(invitePurposePresentation[purpose].description);
});
