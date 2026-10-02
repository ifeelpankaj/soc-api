import type { Href } from "expo-router";

// Native navigation never reads browser notification state.
export function pendingNotificationRoute(): Href | null { return null; }
export function clearPendingNotification() {}
