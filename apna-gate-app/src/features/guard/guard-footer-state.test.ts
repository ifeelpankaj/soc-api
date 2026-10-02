import assert from "node:assert/strict";
import test from "node:test";
import { checkInFooterState, visitorDetailsFooterState } from "./guard-footer-state";
import { deriveScanOutcome } from "./guard-check-in-state";

test("visitor details has fixed actions only for approved and checked-in entries", () => {
  assert.equal(visitorDetailsFooterState("approved"), "check_in");
  assert.equal(visitorDetailsFooterState("checked_in"), "check_out");
  for (const status of ["waiting_approval", "rejected", "checked_out", "expired", "cancelled"] as const) {
    assert.equal(visitorDetailsFooterState(status), null);
  }
  assert.equal(visitorDetailsFooterState(), null);
});

test("Check In retains its controls while submitting and removes stale controls on completion", () => {
  assert.equal(checkInFooterState("ready", "approved"), "check_in");
  assert.equal(checkInFooterState("blocked", "approved", true), "check_in");
  assert.equal(checkInFooterState("just_checked_in", "approved", true), "scan");
  assert.equal(checkInFooterState("already_inside", "approved"), "scan");
  assert.equal(checkInFooterState("loading", "checked_in"), "scan");
  for (const outcome of ["loading", "pending_approval", "blocked", "error"] as const) {
    assert.equal(checkInFooterState(outcome, "approved"), null);
  }
  assert.equal(checkInFooterState("ready", "waiting_approval"), null);
});

test("completed mutation wins over a refresh of old approved data", () => {
  const input = {
    entry: { status: "approved" as const }, entryError: null,
    isLoadingEntry: true, justCheckedInThisSession: true,
    canCheckIn: false, alreadyCheckedIn: true, isResolvingEntry: true,
  };
  assert.equal(deriveScanOutcome(input), "just_checked_in");
  assert.equal(deriveScanOutcome({ ...input, justCheckedInThisSession: false }), "already_inside");
  assert.equal(deriveScanOutcome({ ...input, alreadyCheckedIn: false, justCheckedInThisSession: false }), "loading");
});

test("failed check-in keeps the eligible visitor retryable", () => {
  const outcome = deriveScanOutcome({
    entry: { status: "approved" }, entryError: { kind: "network", message: "Offline" },
    isLoadingEntry: false, justCheckedInThisSession: false, canCheckIn: true,
    alreadyCheckedIn: false, isResolvingEntry: false,
  });
  assert.equal(outcome, "ready");
  assert.equal(checkInFooterState(outcome, "approved"), "check_in");
});
