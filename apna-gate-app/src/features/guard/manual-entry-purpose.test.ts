import assert from "node:assert/strict";
import test from "node:test";
import { flatForManualEntryPurpose, manualEntryFlatId } from "./manual-entry-purpose";
import { canEditVisitorFlat } from "./guard-entry-edit";
import { pendingNotificationRoute } from "../web/notification-handoff";
import { supportsGuardCompanions, usesCabOptionalCompanions } from "./guard-platform-policy";
import { hasCompanionDetailErrors, serializeCompanionDetails, validateCompanionDetails } from "./guard-companions";

test("cab companions are Android-only and reuse guest validation and serialization", () => {
  assert.equal(pendingNotificationRoute(), null);
  for (const platform of ["android", "ios", "web"]) {
    assert.equal(supportsGuardCompanions(platform, "guest"), true);
    assert.equal(supportsGuardCompanions(platform, "cab"), platform === "android");
    assert.equal(usesCabOptionalCompanions(platform, "cab"), platform === "android");
    for (const purpose of ["guest", "service", "delivery", "staff"]) {
      assert.equal(usesCabOptionalCompanions(platform, purpose), false);
    }
  }
  assert.equal(serializeCompanionDetails([], 0), undefined);
  assert.equal(hasCompanionDetailErrors(validateCompanionDetails([], 1)), true);
  assert.equal(hasCompanionDetailErrors(validateCompanionDetails([{ name: "", phoneNumber: "123" }], 1)), true);
  const passengers = [{ name: " Passenger ", phoneNumber: "" }, { name: "", phoneNumber: "98765 43210" }];
  assert.equal(hasCompanionDetailErrors(validateCompanionDetails(passengers, 2)), false);
  assert.deepEqual(serializeCompanionDetails(passengers, 2), [{ full_name: "Passenger" }, { phone_number: "9876543210" }]);
});

test("Service to Guest to Service clears the flat and omits stale submission values", () => {
  let flat: { id: number } | null = { id: 20 };
  flat = flatForManualEntryPurpose("service", flat);
  assert.equal(flat, null);
  flat = flatForManualEntryPurpose("guest", flat);
  assert.equal(manualEntryFlatId("guest", flat), undefined);
  flat = { id: 30 };
  assert.equal(manualEntryFlatId("guest", flat), 30);
  flat = flatForManualEntryPurpose("service", flat);
  assert.equal(flat, null);
  assert.equal(JSON.stringify({ flat_id: manualEntryFlatId("service", { id: 30 }) }), "{}");
  assert.equal(manualEntryFlatId("staff", { id: 30 }), undefined);
});

test("flat editing preserves historical Service but blocks society-wide Service", () => {
  const entry = { source: "guard_entry", purpose: "service", status: "approved" } as const;
  assert.equal(canEditVisitorFlat({ ...entry, flat_id: 0 }), false);
  assert.equal(canEditVisitorFlat({ ...entry, flat_id: 20 }), true);
  assert.equal(canEditVisitorFlat({ ...entry, purpose: "guest", flat_id: 20 }), true);
});
