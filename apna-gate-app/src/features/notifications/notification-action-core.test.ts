import assert from "node:assert/strict";
import test from "node:test";
import {
  createActionExecutor,
  validateAction,
  parseActionResponse,
  REPLAY_TTL,
  type Action,
  type ActionDependencies,
} from "./notification-action-core";

const action: Action = {
  actionIdentifier: "VISITOR_APPROVE",
  notificationId: "a1234567-1234-4567-89ab-123456789abc",
  osNotificationId: "os-1",
  societyId: 1,
  entryId: 2,
  flatId: 3,
  type: "visitor.pending",
};
const stale = {
  status: 409,
  data: { error: { code: "VISITOR_INVALID_STATE" } },
};
function fixture(overrides: Partial<ActionDependencies> = {}) {
  const calls: string[] = [];
  let pending: unknown;
  const deps: ActionDependencies = {
    decide: async () => {
      calls.push("decision");
    },
    markRead: async () => {
      calls.push("read");
    },
    entryState: async () => {
      calls.push("fetch");
      return "approved";
    },
    invalidate: () => {
      calls.push("invalidate");
    },
    dismiss: async () => {
      calls.push("dismiss");
    },
    loadPending: async () => pending,
    savePending: async (value) => {
      pending = value;
      calls.push("save");
    },
    deletePending: async () => {
      pending = undefined;
    },
    forbidden: () => {
      calls.push("forbidden");
    },
    now: () => 1_000_000,
    ...overrides,
  };
  return {
    executor: createActionExecutor(deps),
    deps,
    calls,
    pending: () => pending,
    setPending: (value: unknown) => {
      pending = value;
    },
  };
}
test("decision actions strictly order success, read, invalidation, dismissal", async () => {
  for (const actionIdentifier of [
    "VISITOR_APPROVE",
    "VISITOR_DECLINE",
  ] as const) {
    const f = fixture();
    assert.equal(
      await f.executor.run({ ...action, actionIdentifier }),
      "success",
    );
    assert.deepEqual(f.calls, ["decision", "read", "invalidate", "dismiss"]);
  }
});
test("mark-read failure does not invalidate or dismiss and permits safe decision retry", async () => {
  const f = fixture();
  f.deps.markRead = async () => {
    f.calls.push("read-failed");
    throw { status: 500 };
  };
  assert.equal(await f.executor.run(action), "retryable");
  assert.deepEqual(f.calls, ["decision", "read-failed"]);
  f.deps.markRead = async () => {
    f.calls.push("read");
  };
  assert.equal(await f.executor.run(action), "success");
  assert.deepEqual(f.calls.slice(2), [
    "decision",
    "read",
    "invalidate",
    "dismiss",
  ]);
});
test("informational controls: close only dismisses; mark read orders read before invalidation", async () => {
  const f = fixture();
  await f.executor.run({ ...action, actionIdentifier: "NOTIFICATION_CLOSE" });
  assert.deepEqual(f.calls, ["dismiss"]);
  f.calls.length = 0;
  await f.executor.run({
    ...action,
    actionIdentifier: "NOTIFICATION_MARK_READ",
  });
  assert.deepEqual(f.calls, ["read", "invalidate", "dismiss"]);
});
test("aliases and concurrent callbacks share a promise; successful pairs remain deduplicated", async () => {
  const f = fixture();
  const first = f.executor.run(action);
  const duplicate = f.executor.run({
    ...action,
    actionIdentifier: "approve_visitor",
  });
  assert.equal(first, duplicate);
  await Promise.all([first, duplicate]);
  assert.equal(await f.executor.run(action), "terminal");
  assert.deepEqual(f.calls, ["decision", "read", "invalidate", "dismiss"]);
  for (const [alias, normalized] of [
    ["decline_visitor", "VISITOR_DECLINE"],
    ["mark_notification_read", "NOTIFICATION_MARK_READ"],
    ["close_notification", "NOTIFICATION_CLOSE"],
  ]) {
    assert.equal(
      validateAction({ ...action, actionIdentifier: alias })?.actionIdentifier,
      normalized,
    );
  }
});
test("malformed callbacks and persisted identifiers cannot mutate or dismiss", async () => {
  for (const invalid of [
    null,
    {},
    { ...action, notificationId: "bad" },
    { ...action, societyId: true },
    { ...action, entryId: 1.5 },
    { ...action, entryId: "1e2" },
    { ...action, actionIdentifier: "__proto__" },
    { ...action, osNotificationId: " " },
  ]) {
    const f = fixture();
    assert.equal(await f.executor.run(invalid), "invalid");
    f.setPending({ ...(invalid as object), createdAt: 999_999 });
    await f.executor.replay();
    assert.deepEqual(f.calls, []);
  }
  assert.equal(parseActionResponse({ notification: null }), undefined);
});
test("only a decision's typed conflict and confirmed authorized terminal state reconcile", async () => {
  for (const state of [
    "approved",
    "rejected",
    "checked_in",
    "checked_out",
    "cancelled",
    "expired",
    "auto_closed",
  ]) {
    const f = fixture({
      decide: async () => {
        throw stale;
      },
      entryState: async () => state,
    });
    assert.equal(await f.executor.run(action), "terminal");
    assert.deepEqual(f.calls, ["read", "invalidate", "dismiss"]);
  }
  for (const overrides of [
    {
      decide: async () => {
        throw { status: 409 };
      },
    },
    {
      decide: async () => {
        throw stale;
      },
      entryState: async () => "waiting_approval",
    },
    {
      decide: async () => {
        throw stale;
      },
      entryState: async () => {
        throw { status: 500 };
      },
    },
  ]) {
    const f = fixture(overrides);
    assert.equal(await f.executor.run(action), "retryable");
    assert.deepEqual(f.calls, []);
  }
  const missingFlat = fixture({
    decide: async () => {
      throw stale;
    },
  });
  await missingFlat.executor.run({ ...action, flatId: undefined });
  assert.deepEqual(missingFlat.calls, []);
  const readConflict = fixture({
    markRead: async () => {
      throw stale;
    },
  });
  await readConflict.executor.run(action);
  assert.deepEqual(readConflict.calls, ["decision"]);
});
test("403, network and server failures leave notification presented and unread", async () => {
  for (const status of [403, 500, "FETCH_ERROR"]) {
    const f = fixture({
      decide: async () => {
        throw { status };
      },
    });
    await f.executor.run(action);
    assert.deepEqual(f.calls, status === 403 ? ["forbidden"] : []);
    assert.equal(f.pending(), undefined);
  }
});
test("duplicate forbidden callbacks do not retry the mutation", async () => {
  let mutations = 0;
  const f = fixture({
    decide: async () => {
      mutations++;
      throw { status: 403 };
    },
  });
  await f.executor.run(action);
  await f.executor.run({ ...action, actionIdentifier: "approve_visitor" });
  assert.equal(mutations, 1);
  assert.deepEqual(f.calls, ["forbidden"]);
});
test("failed authentication securely queues a normalized action and replays only once", async () => {
  const f = fixture({
    decide: async () => {
      throw { status: 401 };
    },
  });
  await f.executor.run({ ...action, actionIdentifier: "approve_visitor" });
  assert.equal((f.pending() as Action).actionIdentifier, "VISITOR_APPROVE");
  await Promise.all([f.executor.replay(), f.executor.replay()]);
  assert.equal(f.pending(), undefined);
  assert.deepEqual(f.calls, ["save"]); // Replay failure never requeues.
  const success = fixture();
  success.setPending({
    ...action,
    actionIdentifier: "approve_visitor",
    createdAt: 999_999,
  });
  await Promise.all([success.executor.replay(), success.executor.replay()]);
  assert.deepEqual(success.calls, [
    "decision",
    "read",
    "invalidate",
    "dismiss",
  ]);
});
test("expired, future, malformed replay timestamps are discarded without mutation or dismissal", async () => {
  for (const createdAt of [
    1_000_000 - REPLAY_TTL,
    1_000_001,
    -1,
    NaN,
    Infinity,
    "999999",
    undefined,
  ]) {
    const f = fixture();
    f.setPending({ ...action, createdAt });
    await f.executor.replay();
    assert.equal(f.pending(), undefined);
    assert.deepEqual(f.calls, []);
  }
});
test("repeated unauthorized callbacks do not extend the pending replay window", async () => {
  const f = fixture({
    decide: async () => {
      throw { status: 401 };
    },
  });
  f.setPending({ ...action, createdAt: 800_000 });
  await f.executor.run(action);
  assert.equal((f.pending() as { createdAt: number }).createdAt, 800_000);
});
test("independent devices use the same durable row and safely complete idempotent decisions", async () => {
  const reads: string[] = [];
  for (const osNotificationId of ["phone-one", "phone-two"]) {
    const f = fixture({
      decide: async () => ({ status: "approved" }),
      markRead: async (id) => {
        reads.push(id);
      },
    });
    assert.equal(
      await f.executor.run({ ...action, osNotificationId }),
      "success",
    );
  }
  assert.deepEqual(reads, [action.notificationId, action.notificationId]);
});
