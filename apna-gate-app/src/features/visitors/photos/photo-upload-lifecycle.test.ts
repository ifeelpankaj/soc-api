import assert from "node:assert/strict";
import test from "node:test";
import { runPhotoUploadLifecycle } from "./photo-upload-lifecycle";

test("marks uploading synchronously and cleans up after success", async () => {
  const events: string[] = [];
  let completeUpload: (() => void) | undefined;
  const pending = runPhotoUploadLifecycle({
    start: () => events.push("start"),
    upload: () =>
      new Promise<void>((resolve) => {
        completeUpload = resolve;
      }),
    onFailure: () => events.push("failure"),
    cleanup: () => events.push("cleanup"),
    finish: () => events.push("finish"),
  });

  assert.deepEqual(events, ["start"]);
  completeUpload?.();
  await pending;
  assert.deepEqual(events, ["start", "cleanup", "finish"]);
});

test("reports failure before cleaning up and finishing", async () => {
  const events: string[] = [];
  await runPhotoUploadLifecycle({
    start: () => events.push("start"),
    upload: async () => {
      throw new Error("offline");
    },
    onFailure: () => events.push("failure"),
    cleanup: () => events.push("cleanup"),
    finish: () => events.push("finish"),
  });

  assert.deepEqual(events, ["start", "failure", "cleanup", "finish"]);
});

