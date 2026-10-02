import { useRouter } from "expo-router";
import * as Notifications from "expo-notifications";
import { type PropsWithChildren, useEffect, useRef } from "react";
import { AppState, Platform } from "react-native";

import { useAuth } from "@/features/auth/use-auth";
import {
  handleNotificationAction,
  isNotificationActionResponse,
  replayPendingNotificationAction,
  reconcileForbiddenNotificationActions,
} from "@/features/notifications/notification-actions";
import {
  getNotificationType,
  notificationRoute,
  notificationPresentation,
  stringValue,
} from "@/features/notifications/notification-routing";
import {
  enhancedApi,
  invalidateVisitorNotificationTags,
} from "@/lib/api/enhanced-api";
import { addNotificationListeners } from "@/lib/notifications/register-notifications";
import { useAppDispatch } from "@/redux/hooks";
import { useToast } from "@/components/ui";

export function NotificationSyncProvider({ children }: PropsWithChildren) {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const { showToast } = useToast();
  const { homeRoute, status } = useAuth();
  const recentEventsRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    if (Platform.OS === "web" || status !== "authenticated") {
      return;
    }

    void replayPendingNotificationAction();
    void reconcileForbiddenNotificationActions();
    const foregroundSubscription = AppState.addEventListener(
      "change",
      (state) => {
        if (state === "active") {
          void replayPendingNotificationAction();
          void reconcileForbiddenNotificationActions();
        }
      },
    );

    const invalidateFromNotification = (
      data: Record<string, unknown> | undefined,
    ) => {
      const eventType = getNotificationType(data);
      dispatch(
        enhancedApi.util.invalidateTags(
          invalidateVisitorNotificationTags(eventType),
        ),
      );
      dispatch(enhancedApi.util.invalidateTags(["Notifications"]));
    };
    const showForegroundToast = (
      data: Record<string, unknown> | undefined,
      title: string | null | undefined,
      body: string | null | undefined,
    ) => {
      const key = notificationDedupeKey(data);
      if (key && recentEventsRef.current.has(key)) {
        return;
      }
      if (key) {
        recentEventsRef.current.add(key);
        setTimeout(() => recentEventsRef.current.delete(key), 60_000);
      }
      const message = notificationMessage(data, title, body);
      if (message) {
        showToast({
          title: message.title,
          message: message.body,
          variant: message.variant,
        });
      }
    };
    const navigateFromNotification = (
      data: Record<string, unknown> | undefined,
    ) => {
      const route = notificationRoute(
        {
          id:
            stringValue(data?.notification_id) ??
            notificationDedupeKey(data) ??
            "",
          type: getNotificationType(data) ?? "",
          title: "",
          body: "",
          data: data ?? {},
          user_id: 0,
          created_at: new Date().toISOString(),
        },
        homeRoute,
      );
      if (route) {
        router.push(route);
      }
    };

    const lastResponse = Notifications.getLastNotificationResponse();
    if (lastResponse && !isNotificationActionResponse(lastResponse)) {
      Notifications.clearLastNotificationResponse();
      navigateFromNotification(lastResponse.notification.request.content.data);
    }
    const removeListeners = addNotificationListeners({
      onReceived: (notification) => {
        const content = notification.request.content;
        const data = content.data as Record<string, unknown> | undefined;
        invalidateFromNotification(data);
        showForegroundToast(data, content.title, content.body);
      },
      onResponse: (response) => {
        const data = response.notification.request.content.data as
          Record<string, unknown> | undefined;
        if (isNotificationActionResponse(response)) {
          void handleNotificationAction(response);
          return;
        }
        invalidateFromNotification(data);
        navigateFromNotification(data);
      },
    });
    return () => {
      removeListeners();
      foregroundSubscription.remove();
    };
  }, [dispatch, homeRoute, router, showToast, status]);

  return children;
}

function notificationDedupeKey(data: Record<string, unknown> | undefined) {
  const type = getNotificationType(data);
  const inviteID = stringValue(data?.invite_id);
  const entryID = stringValue(data?.entry_id);
  const occurredAt = stringValue(data?.occurred_at);
  const notificationID = stringValue(data?.notification_id);
  if (notificationID) {
    return notificationID;
  }
  if (!type) {
    return undefined;
  }
  return [type, inviteID, entryID, occurredAt].filter(Boolean).join(":");
}

function notificationMessage(
  data: Record<string, unknown> | undefined,
  title: string | null | undefined,
  body: string | null | undefined,
) {
  const presentation = notificationPresentation({
    type: getNotificationType(data) ?? "",
    title: title ?? "",
    body: body ?? "",
    data,
  });
  return {
    title: presentation.title,
    body: presentation.body,
    variant:
      presentation.tone === "warning" ? ("info" as const) : presentation.tone,
  };
}
