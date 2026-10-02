import assert from "node:assert/strict";
import test from "node:test";
import { photoCacheKey, photoRefreshDelay, type PhotoQuery } from "./photo-utils";

test("original photo cache is reusable and isolated from previews and other access", () => {
  const query: PhotoQuery = {
    userId: 1, societyId: 2, entryId: 3, context: "guard",
    photoReference: "/photo.jpg", variant: "original",
  };
  const key = photoCacheKey(query);
  assert.equal(photoCacheKey({ ...query }), key);
  for (const change of [
    { variant: "detail" as const }, { userId: 9 }, { societyId: 9 },
    { entryId: 9 }, { context: "resident:4" as const },
    { photoReference: "/replacement.jpg" },
  ]) {
    assert.notEqual(photoCacheKey({ ...query, ...change }), key);
  }
});

test("valid cached URLs wait to refresh; near-expiry and expired URLs refresh immediately", () => {
  const now = Date.parse("2026-09-07T12:00:00Z");
  assert.equal(photoRefreshDelay("2026-09-07T12:15:00Z", now), 12 * 60_000);
  assert.equal(photoRefreshDelay("2026-09-07T12:02:00Z", now), 0);
  assert.equal(photoRefreshDelay("2026-09-07T11:59:00Z", now), 0);
});
