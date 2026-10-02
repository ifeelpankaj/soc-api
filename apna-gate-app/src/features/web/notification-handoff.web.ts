import type { Href } from "expo-router";
import { validNotificationLink } from "./notification-target";

export function pendingNotificationRoute(): Href | null {
  try {
    const raw = sessionStorage.getItem("apna_notification_open");
    if (!raw) return null;
    const target = JSON.parse(raw);
    if (!validNotificationLink(target?.id, target?.recipient)) {
      clearPendingNotification();
      return null;
    }
    return { pathname: "/notification-open", params: { id: target.id, recipient: target.recipient } };
  } catch { return null; }
}

export function clearPendingNotification() {
  try { sessionStorage.removeItem("apna_notification_open"); } catch { /* Storage may be unavailable. */ }
}
