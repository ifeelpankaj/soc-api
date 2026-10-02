export type ActionId =
  | "VISITOR_APPROVE"
  | "VISITOR_DECLINE"
  | "NOTIFICATION_MARK_READ"
  | "NOTIFICATION_CLOSE";
export type Action = {
  actionIdentifier: ActionId;
  notificationId: string;
  osNotificationId: string;
  societyId?: number;
  entryId?: number;
  flatId?: number;
  type?: string;
};
type Outcome = "success" | "terminal" | "retryable" | "forbidden" | "invalid";
export const REPLAY_TTL = 5 * 60 * 1000;
const aliases: Record<string, ActionId> = {
  VISITOR_APPROVE: "VISITOR_APPROVE",
  approve_visitor: "VISITOR_APPROVE",
  VISITOR_DECLINE: "VISITOR_DECLINE",
  decline_visitor: "VISITOR_DECLINE",
  NOTIFICATION_MARK_READ: "NOTIFICATION_MARK_READ",
  mark_notification_read: "NOTIFICATION_MARK_READ",
  NOTIFICATION_CLOSE: "NOTIFICATION_CLOSE",
  close_notification: "NOTIFICATION_CLOSE",
};
const resolvedStates = new Set([
  "approved",
  "rejected",
  "checked_in",
  "checked_out",
  "cancelled",
  "expired",
  "auto_closed",
]);
function object(value: unknown): Record<string, unknown> {
  return value !== null && typeof value === "object"
    ? (value as Record<string, unknown>)
    : {};
}
function text(value: unknown) {
  return typeof value === "string" && value.trim() ? value.trim() : undefined;
}
function positiveId(value: unknown) {
  if (
    typeof value !== "number" &&
    (typeof value !== "string" || !/^\d+$/.test(value))
  )
    return undefined;
  const n = Number(value);
  return Number.isSafeInteger(n) && n > 0 ? n : undefined;
}
export function validateAction(value: unknown): Action | undefined {
  const data = object(value);
  const identifier = text(data.actionIdentifier);
  const actionIdentifier =
    identifier && Object.hasOwn(aliases, identifier)
      ? aliases[identifier]
      : undefined;
  const notificationId = text(data.notificationId);
  const osNotificationId = text(data.osNotificationId);
  if (
    !actionIdentifier ||
    !notificationId ||
    !/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
      notificationId,
    ) ||
    !osNotificationId
  )
    return;
  const societyId = positiveId(data.societyId);
  const entryId = positiveId(data.entryId);
  const flatId = positiveId(data.flatId);
  if (
    (actionIdentifier === "VISITOR_APPROVE" ||
      actionIdentifier === "VISITOR_DECLINE") &&
    (!societyId || !entryId)
  )
    return;
  return {
    actionIdentifier,
    notificationId: notificationId.toLowerCase(),
    osNotificationId,
    societyId,
    entryId,
    flatId,
    type: text(data.type),
  };
}
export function parseActionResponse(value: unknown) {
  const response = object(value);
  const request = object(object(response.notification).request);
  const data = object(object(request.content).data);
  return validateAction({
    actionIdentifier: response.actionIdentifier,
    notificationId: data.notification_id,
    osNotificationId: request.identifier,
    societyId: data.society_id,
    entryId: data.entry_id,
    flatId: data.flat_id,
    type: data.type ?? data.event,
  });
}
export type ActionDependencies = {
  decide: (action: Action) => Promise<unknown>;
  markRead: (id: string) => Promise<unknown>;
  entryState: (action: Action) => Promise<unknown>;
  invalidate: (action: Action) => void;
  dismiss: (id: string) => Promise<void>;
  loadPending: () => Promise<unknown>;
  savePending: (value: Action & { createdAt: number }) => Promise<void>;
  deletePending: () => Promise<void>;
  forbidden: (action: Action) => void;
  now?: () => number;
};
export function createActionExecutor(deps: ActionDependencies) {
  const inFlight = new Map<string, Promise<Outcome>>();
  const terminal = new Set<string>();
  const forbidden = new Set<string>();
  const now = deps.now ?? Date.now;
  const keyOf = (a: Action) => `${a.notificationId}:${a.actionIdentifier}`;
  let replayInFlight: Promise<void> | undefined;

  async function failure(
    error: unknown,
    action: Action,
    persist: boolean,
  ): Promise<Outcome> {
    const status = object(error).status;
    if (status === 403) {
      deps.forbidden(action);
      return "forbidden";
    }
    if (status === 401 && persist) {
      try {
        const previous = object(await deps.loadPending());
        const previousAction = validateAction(previous);
        // Repeated callbacks must not extend an existing replay window.
        const createdAt =
          previousAction &&
          keyOf(previousAction) === keyOf(action) &&
          typeof previous.createdAt === "number"
            ? previous.createdAt
            : now();
        await deps.savePending({ ...action, createdAt });
      } catch {
        /* Secure storage failure must not hide the notification. */
      }
    }
    return "retryable";
  }
  async function execute(action: Action, persist: boolean): Promise<Outcome> {
    try {
      if (action.actionIdentifier === "NOTIFICATION_CLOSE") {
        await deps.dismiss(action.osNotificationId);
        return "success";
      }
      let stale = false;
      if (
        action.actionIdentifier === "VISITOR_APPROVE" ||
        action.actionIdentifier === "VISITOR_DECLINE"
      ) {
        try {
          await deps.decide(action);
        } catch (error) {
          const e = object(error);
          const code = object(object(e.data).error).code;
          if (
            e.status !== 409 ||
            code !== "VISITOR_INVALID_STATE" ||
            !action.flatId
          )
            return failure(error, action, persist);
          const state = await deps.entryState(action);
          if (typeof state !== "string" || !resolvedStates.has(state))
            return "retryable";
          stale = true;
        }
      }
      // No endpoint used above may automatically invalidate tags.
      await deps.markRead(action.notificationId);
      deps.invalidate(action);
      await deps.dismiss(action.osNotificationId);
      return stale ? "terminal" : "success";
    } catch (error) {
      return failure(error, action, persist);
    }
  }
  function run(value: unknown, persist = true): Promise<Outcome> {
    const action = validateAction(value);
    if (!action) return Promise.resolve("invalid");
    const key = keyOf(action);
    if (terminal.has(key)) return Promise.resolve("terminal");
    if (forbidden.has(key)) return Promise.resolve("forbidden");
    const existing = inFlight.get(key);
    if (existing) return existing;
    const promise = Promise.resolve()
      .then(() => execute(action, persist))
      .then((outcome) => {
        if (outcome === "success" || outcome === "terminal") terminal.add(key);
        if (outcome === "forbidden") forbidden.add(key);
        return outcome;
      })
      .finally(() => inFlight.delete(key));
    inFlight.set(key, promise);
    return promise;
  }
  function replay(): Promise<void> {
    if (replayInFlight) return replayInFlight;
    replayInFlight = (async () => {
      // Early response processing may still be persisting a failed-auth action.
      await Promise.all(inFlight.values());
      const pending = object(await deps.loadPending());
      await deps.deletePending();
      const action = validateAction(pending);
      const timestamp = pending.createdAt;
      if (
        !action ||
        typeof timestamp !== "number" ||
        !Number.isFinite(timestamp) ||
        timestamp <= 0 ||
        timestamp > now() ||
        now() - timestamp >= REPLAY_TTL
      )
        return;
      await run(action, false);
    })()
      .catch(() => undefined)
      .finally(() => {
        replayInFlight = undefined;
      });
    return replayInFlight;
  }
  return {
    run,
    replay,
    handle: (value: unknown) => run(parseActionResponse(value)),
  };
}
