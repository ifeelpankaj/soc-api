import type { HubCategory } from "@/lib/api/hub-api";

export function formatHubRelativeTime(iso: string) {
  const timestamp = Date.parse(iso);
  if (!Number.isFinite(timestamp)) {
    return "";
  }
  const deltaMs = Date.now() - timestamp;
  const minutes = Math.floor(deltaMs / 60_000);
  if (minutes < 1) {
    return "Just now";
  }
  if (minutes < 60) {
    return `${minutes}m ago`;
  }
  const hours = Math.floor(minutes / 60);
  if (hours < 48) {
    return `${hours}h ago`;
  }
  return new Date(iso).toLocaleDateString("en-IN", {
    day: "numeric",
    month: "short",
  });
}

export function hubCategoryName(
  categories: HubCategory[] | undefined,
  categoryId?: number,
): string | undefined {
  if (!categoryId || !categories?.length) {
    return undefined;
  }
  return categories.find((c) => c.id === categoryId)?.name;
}

export function hubPostPreview(body: string, maxLength = 120) {
  const normalized = body.replace(/\s+/g, " ").trim();
  if (normalized.length <= maxLength) {
    return normalized;
  }
  return `${normalized.slice(0, maxLength - 1)}…`;
}
