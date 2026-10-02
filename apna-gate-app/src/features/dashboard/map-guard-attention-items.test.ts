import assert from "node:assert/strict";
import test from "node:test";

import { mapGuardAttentionItems } from "./map-guard-attention-items";

test("mapGuardAttentionItems builds pending and ready rows", () => {
  let pending = false;
  let checkIn = false;
  const items = mapGuardAttentionItems({
    pendingCount: 2,
    readyCount: 1,
    onReviewPending: () => {
      pending = true;
    },
    onCheckIn: () => {
      checkIn = true;
    },
  });

  assert.equal(items.length, 2);
  assert.match(items[0].title, /2 visitors waiting/);
  assert.match(items[1].title, /1 visitor ready/);
  items[0].onPress();
  items[1].onPress();
  assert.equal(pending, true);
  assert.equal(checkIn, true);
});

test("mapGuardAttentionItems returns empty when counts are zero", () => {
  assert.deepEqual(
    mapGuardAttentionItems({
      pendingCount: 0,
      readyCount: 0,
      onReviewPending: () => {},
      onCheckIn: () => {},
    }),
    [],
  );
});
