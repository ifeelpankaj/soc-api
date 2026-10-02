import assert from "node:assert/strict";
import test from "node:test";
import {
  passwordChangeError,
  utf8Length,
  validatePasswordChange,
} from "./password-change";
import { SessionWriteQueue } from "./session-write-queue";

test("password form identifies required fields, mismatch and reuse without trimming", () => {
  assert.deepEqual(
    Object.keys(validatePasswordChange({ current: "", next: "", confirm: "" })),
    ["current", "next", "confirm"],
  );
  assert.ok(
    validatePasswordChange({
      current: "oldpassword",
      next: "oldpassword",
      confirm: "oldpassword",
    }).next,
  );
  assert.ok(
    validatePasswordChange({
      current: "oldpassword",
      next: "newpassword",
      confirm: "different",
    }).confirm,
  );
  assert.deepEqual(
    validatePasswordChange({
      current: "oldpassword",
      next: " newpassword ",
      confirm: " newpassword ",
    }),
    {},
  );
});
test("new passwords respect UTF-8 bcrypt limits", () => {
  assert.equal(utf8Length("😀".repeat(18)), 72);
  for (const next of ["a".repeat(72), "😀".repeat(18)])
    assert.deepEqual(
      validatePasswordChange({ current: "oldpassword", next, confirm: next }),
      {},
    );
  for (const next of ["a".repeat(73), "😀".repeat(19), "short"])
    assert.ok(
      validatePasswordChange({ current: "oldpassword", next, confirm: next })
        .next,
    );
});
test("password API errors select useful fields and recovery", () => {
  for (const code of ["CURRENT_PASSWORD_INCORRECT", "INVALID_CREDENTIALS"])
    assert.equal(
      passwordChangeError({ data: { error: { code } } }).field,
      "current",
    );
  assert.equal(
    passwordChangeError({ data: { error: { code: "PASSWORD_NOT_SET" } } })
      .reset,
    true,
  );
  for (const status of ["TIMEOUT_ERROR", "FETCH_ERROR"])
    assert.equal(passwordChangeError({ status }).uncertain, true);
});
test("sign-out clears a running write and rejects a late refresh write", async () => {
  const queue = new SessionWriteQueue();
  let saved: string | null = null;
  let release!: () => void;
  const pending = new Promise<void>((resolve) => {
    release = resolve;
  });
  const generation = queue.current();
  const first = queue.save(generation, async () => {
    await pending;
    saved = "old-token";
  });
  await Promise.resolve();
  await Promise.resolve();
  const clear = queue.clear(async () => {
    saved = null;
  });
  const late = queue.save(generation, async () => {
    saved = "late-token";
  });
  release();
  await Promise.all([first, clear, late]);
  assert.equal(saved, null);
  await queue.save(queue.current(), async () => {
    saved = "new-login";
  });
  assert.equal(saved, "new-login");
});
test("cleanup failure cannot wedge subsequent session writes", async () => {
  const queue = new SessionWriteQueue();
  await assert.rejects(
    queue.clear(async () => {
      throw new Error("storage failure");
    }),
  );
  let ran = false;
  await queue.save(queue.current(), async () => {
    ran = true;
  });
  assert.equal(ran, true);
});
