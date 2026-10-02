import * as Notifications from "expo-notifications";
import * as SecureStore from "expo-secure-store";
import * as TaskManager from "expo-task-manager";
import { Platform } from "react-native";
import type { FetchArgs } from "@reduxjs/toolkit/query";
import {
  enhancedApi,
  invalidateVisitorNotificationTags,
} from "@/lib/api/enhanced-api";
import { store } from "@/redux/store";
import { createActionExecutor, type Action } from "./notification-action-core";

const TASK = "apna-gate-notification-action";
const PENDING_KEY = "notifications.pendingAction";
// No invalidatesTags: the executor owns ordering, including mark-read failures.
const actionApi = enhancedApi.injectEndpoints({
  endpoints: (builder) => ({
    notificationActionRequest: builder.mutation<unknown, FetchArgs>({
      query: (request) => request,
      extraOptions: { skipForbiddenRetry: true },
    }),
  }),
});
async function request(args: FetchArgs) {
  const operation = store.dispatch(
    actionApi.endpoints.notificationActionRequest.initiate(args),
  );
  try {
    return await operation.unwrap();
  } finally {
    operation.reset();
  }
}
const forbiddenActions = new Map<string, Action>();
async function entryState(action: Action) {
  const result = (await request({
    url: `/v1/societies/${action.societyId}/flats/${action.flatId}/visitor-entries/${action.entryId}`,
    method: "GET",
  })) as { data?: { entry?: { status?: string } } };
  return result.data?.entry?.status;
}
function invalidate(action: Action) {
  if (action.actionIdentifier !== "NOTIFICATION_MARK_READ") {
    store.dispatch(
      enhancedApi.util.invalidateTags(
        invalidateVisitorNotificationTags(action.type),
      ),
    );
  }
  store.dispatch(enhancedApi.util.invalidateTags(["Notifications"]));
}
const executor = createActionExecutor({
  decide: (action) =>
    request({
      url: `/v1/societies/${action.societyId}/visitor-entries/${action.entryId}/${action.actionIdentifier === "VISITOR_APPROVE" ? "approve" : "reject"}`,
      method: "POST",
      ...(action.actionIdentifier === "VISITOR_DECLINE"
        ? { body: { reason: "Declined by resident from notification action" } }
        : {}),
    }),
  markRead: (id) =>
    request({ url: `/v1/me/notifications/${id}/read`, method: "PATCH" }),
  entryState,
  invalidate,
  dismiss: (id) => Notifications.dismissNotificationAsync(id),
  loadPending: async () => {
    if (Platform.OS === "web") return null;
    const value = await SecureStore.getItemAsync(PENDING_KEY);
    if (!value) return null;
    try {
      return JSON.parse(value) as unknown;
    } catch {
      return null;
    }
  },
  savePending: async (value) => {
    if (Platform.OS !== "web")
      await SecureStore.setItemAsync(PENDING_KEY, JSON.stringify(value));
  },
  deletePending: async () => {
    if (Platform.OS !== "web") await SecureStore.deleteItemAsync(PENDING_KEY);
  },
  forbidden: (action) => {
    forbiddenActions.set(action.notificationId, action);
  },
});
export const handleNotificationAction = executor.handle;
export const replayPendingNotificationAction = executor.replay;
export function isNotificationActionResponse(
  response: Notifications.NotificationResponse,
) {
  return response.actionIdentifier !== Notifications.DEFAULT_ACTION_IDENTIFIER;
}
// Reconcile on the next authenticated foreground session, without retrying a 403 mutation.
export async function reconcileForbiddenNotificationActions() {
  for (const [id, action] of forbiddenActions) {
    try {
      if (action.flatId && action.entryId && action.societyId)
        await entryState(action);
      forbiddenActions.delete(id);
    } catch {
      /* Keep for a later foreground session. */
    } finally {
      // The inbox can still be refreshed when access to the entry was revoked.
      invalidate(action);
    }
  }
}
if (Platform.OS !== "web") {
  Notifications.addNotificationResponseReceivedListener((response) => {
    if (isNotificationActionResponse(response)) {
      Notifications.clearLastNotificationResponse();
      void handleNotificationAction(response);
    }
  });
  const lastResponse = Notifications.getLastNotificationResponse();
  if (lastResponse && isNotificationActionResponse(lastResponse)) {
    Notifications.clearLastNotificationResponse();
    void handleNotificationAction(lastResponse);
  }
}
if (Platform.OS === "android") {
  if (!TaskManager.isTaskDefined(TASK)) {
    TaskManager.defineTask<Notifications.NotificationTaskPayload>(
      TASK,
      async ({ data }) => {
        if (data && typeof data === "object" && "actionIdentifier" in data)
          await handleNotificationAction(data);
        return Notifications.BackgroundNotificationTaskResult.NoData;
      },
    );
  }
  void Notifications.registerTaskAsync(TASK).catch(() => undefined);
}
