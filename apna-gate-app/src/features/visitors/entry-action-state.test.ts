import assert from "node:assert/strict";
import test from "node:test";

import { actionForEntry } from "./entry-action-state";

test("returns the action only for the entry performing it", () => {
  const state = { entryId: 50, action: "approve" as const };

  assert.equal(actionForEntry(state, 50), "approve");
  assert.equal(actionForEntry(state, 49), undefined);
  assert.equal(actionForEntry(state, 51), undefined);
});

test("returns no loading action when idle or the entry has no id", () => {
  assert.equal(actionForEntry(undefined, 50), undefined);
  assert.equal(actionForEntry({ entryId: 50, action: "notify" }, undefined), undefined);
});
