import assert from "node:assert/strict";
import test from "node:test";

import {
  defaultCategoryId,
  HUB_MAX_ATTACHMENTS,
  normalizeAttachmentIds,
  validateHubImageMime,
  validateHubImageSize,
} from "./hub-upload-validation";

test("normalizeAttachmentIds dedupes and filters", () => {
  assert.deepEqual(normalizeAttachmentIds([1, 2, 2, 3]), [1, 2, 3]);
  assert.deepEqual(normalizeAttachmentIds([]), []);
});

test("normalizeAttachmentIds rejects too many", () => {
  const ids = Array.from({ length: HUB_MAX_ATTACHMENTS + 1 }, (_, i) => i + 1);
  assert.throws(() => normalizeAttachmentIds(ids), /up to 5/);
});

test("validateHubImageMime", () => {
  assert.equal(validateHubImageMime("image/jpeg"), true);
  assert.equal(validateHubImageMime("image/png"), true);
  assert.equal(validateHubImageMime("application/pdf"), false);
});

test("validateHubImageSize", () => {
  assert.equal(validateHubImageSize(1024), true);
  assert.equal(validateHubImageSize(0), false);
  assert.equal(validateHubImageSize(6 * 1024 * 1024), false);
});

test("defaultCategoryId prefers general", () => {
  assert.equal(
    defaultCategoryId([
      { id: 10, code: "events", name: "Events" },
      { id: 2, code: "general", name: "General" },
    ]),
    2,
  );
  assert.equal(defaultCategoryId([{ id: 5, name: "General" }]), 5);
});
