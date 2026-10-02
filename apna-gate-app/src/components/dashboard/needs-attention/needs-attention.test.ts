import assert from "node:assert/strict";
import test from "node:test";

import type { AttentionItem } from "./types";
import {
  formatAttentionCount,
  hiddenAttentionCount,
  shouldShowExpandControl,
  shouldShowNeedsAttention,
  visibleAttentionItems,
} from "./needs-attention-logic";

function item(id: string): AttentionItem {
  return {
    id,
    icon: { ios: "person", android: "person", web: "person" },
    title: id,
    subtitle: "sub",
    onPress: () => {},
  };
}

test("formatAttentionCount pluralizes correctly", () => {
  assert.equal(formatAttentionCount(1), "1 action");
  assert.equal(formatAttentionCount(4), "4 actions");
  assert.equal(formatAttentionCount(10), "10 actions");
});

test("shouldShowNeedsAttention hides empty lists", () => {
  assert.equal(shouldShowNeedsAttention([]), false);
  assert.equal(shouldShowNeedsAttention([item("a")]), true);
});

test("expand control only appears for 3+ items", () => {
  assert.equal(shouldShowExpandControl([]), false);
  assert.equal(shouldShowExpandControl([item("a")]), false);
  assert.equal(shouldShowExpandControl([item("a"), item("b")]), false);
  assert.equal(shouldShowExpandControl([item("a"), item("b"), item("c")]), true);
  assert.equal(shouldShowExpandControl(Array.from({ length: 10 }, (_, i) => item(String(i)))), true);
});

test("visibleAttentionItems respects collapse rules", () => {
  const three = [item("1"), item("2"), item("3")];
  assert.deepEqual(visibleAttentionItems(three, false).map((v) => v.id), ["1", "2"]);
  assert.deepEqual(visibleAttentionItems(three, true).map((v) => v.id), ["1", "2", "3"]);

  const two = [item("1"), item("2")];
  assert.deepEqual(visibleAttentionItems(two, false).map((v) => v.id), ["1", "2"]);
});

test("hiddenAttentionCount matches visible slice", () => {
  const five = Array.from({ length: 5 }, (_, i) => item(String(i)));
  assert.equal(hiddenAttentionCount(five, false), 3);
  assert.equal(hiddenAttentionCount(five, true), 0);
  assert.equal(hiddenAttentionCount([item("a")], false), 0);
});
