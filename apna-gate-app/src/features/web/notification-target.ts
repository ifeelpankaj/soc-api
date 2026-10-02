import type { ModelsBootstrapData } from "@/lib/api/generated-api";
import type { AppNotification } from "@/lib/api/notification-api-extensions";

export function notificationWorkspace(item: AppNotification, bootstrap: ModelsBootstrapData, userID: number) {
  if (item.user_id !== userID) return null;
  const societyID = Number(item.society_id ?? item.data.society_id);
  const flatID = Number(item.flat_id ?? item.data.flat_id);
  if (!Number.isSafeInteger(societyID) || societyID <= 0 || !Number.isSafeInteger(flatID) || flatID <= 0) return null;
  return bootstrap.residences?.find((r) => r.status === "active" && r.society_id === societyID && r.flat_id === flatID) ?? null;
}
export function validNotificationLink(id?: unknown, recipient?: unknown) {
  return Boolean(typeof id === "string" && /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(id) && typeof recipient === "string" && /^[1-9]\d*$/.test(recipient));
}
