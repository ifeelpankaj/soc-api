import assert from "node:assert/strict";
import test from "node:test";

import { RESIDENT_SEGMENT_OPTIONS } from "./resident-entries-screen-options";

test("resident visitor entry tabs expose only expected, inside, and all", () => {
  assert.deepEqual(RESIDENT_SEGMENT_OPTIONS, [
    { label: "Expected", value: "expected" },
    { label: "Inside", value: "inside" },
    { label: "All", value: "all" },
  ]);
});
