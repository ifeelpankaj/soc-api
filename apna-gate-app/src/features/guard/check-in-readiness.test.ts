import assert from "node:assert/strict";
import test from "node:test";
import { consumeReadyEntries, loadCheckInReadiness, readyToCheckInMessage } from "./check-in-readiness";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";

function entry(id: number, changes: Partial<ModelsVisitorEntry> = {}): ModelsVisitorEntry {
  return { id, society_id: 10, status: "approved", source: "guard_entry", purpose: "guest", ...changes };
}

test("readiness includes only approved, unchecked-in visitors in the selected society's gate queue", async () => {
  const entries = [entry(1), entry(2, { source: "public_qr" }), entry(3, { status: "waiting_approval" }),
    entry(4, { status: "checked_in" }), entry(5, { checked_in_at: "2026-09-13T10:00:00Z" }),
    entry(6, { source: "resident_link" }), entry(7, { society_id: 20 }), entry(1), entry(0)];
  const result = await loadCheckInReadiness(10, async () => ({ entries, total: entries.length }));
  assert.deepEqual(result, { entryIds: [1, 2], total: 2 });
});

test("readiness reads all pages rather than silently ignoring visitors beyond the first page", async () => {
  const entries = Array.from({ length: 205 }, (_, i) => entry(i + 1));
  const offsets: number[] = [];
  const result = await loadCheckInReadiness(10, async (limit, offset) => {
    offsets.push(offset);
    return { entries: entries.slice(offset, offset + limit), total: entries.length };
  });
  assert.deepEqual(offsets, [0, 100, 200]);
  assert.equal(result.total, 205);
  assert.equal(result.entryIds.at(-1), 205);
});

test("failed or malformed snapshots are errors, not empty queues", async () => {
  await assert.rejects(loadCheckInReadiness(10, async () => { throw new Error("offline"); }), /offline/);
  await assert.rejects(loadCheckInReadiness(10, async () => ({})), /Invalid/);
  await assert.rejects(loadCheckInReadiness(10, async () => ({ entries: [], total: NaN })), /Invalid/);
});

test("vibrate once for the initial queue and for new entries, including equal-count replacements", () => {
  const seen = new Set<number>();
  assert.equal(consumeReadyEntries(seen, []), false);
  assert.equal(consumeReadyEntries(seen, [1, 2]), true);
  assert.equal(consumeReadyEntries(seen, [2, 1]), false);
  assert.equal(consumeReadyEntries(seen, [2, 3]), true);
  assert.equal(consumeReadyEntries(seen, [3]), false);
  assert.equal(consumeReadyEntries(seen, []), false);
  assert.equal(consumeReadyEntries(seen, [1, 2, 3]), false);
  // A new user/society session starts with its own seen set.
  assert.equal(consumeReadyEntries(new Set(), [1]), true);
});

test("ready labels explain the action and distinguish singular/plural", () => {
  assert.equal(readyToCheckInMessage(1), "1 visitor ready to check in");
  assert.equal(readyToCheckInMessage(3), "3 visitors ready to check in");
});
