import assert from "node:assert/strict";
import test from "node:test";
import { mergePaginatedItems } from "./use-paginated-query";

type Item = { id: number; value: string };
const itemKey = (item: Item) => item.id;

test("overlapping offset pages contain each entry exactly once", () => {
  const merged = mergePaginatedItems(
    [
      { id: 50, value: "first" },
      { id: 49, value: "older" },
    ],
    [
      { id: 50, value: "fresh" },
      { id: 48, value: "oldest" },
    ],
    itemKey,
  );
  assert.deepEqual(merged, [
    { id: 50, value: "fresh" },
    { id: 49, value: "older" },
    { id: 48, value: "oldest" },
  ]);
});

test("items without an identity are retained instead of incorrectly merged", () => {
  const merged = mergePaginatedItems(
    [{ id: 1, value: "known" }],
    [{ id: 0, value: "unknown" }],
    (item) => item.id || undefined,
  );
  assert.equal(merged.length, 2);
});
