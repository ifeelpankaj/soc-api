import type { ModelsVisitorEntry } from "@/lib/api/generated-api";

export type CheckInReadiness = { entryIds: number[]; total: number };
type QueuePage = { entries?: ModelsVisitorEntry[]; total?: number };

// Use the same queue as manual check-in; resident invitations remain Expected.
export async function loadCheckInReadiness(
  societyId: number,
  fetchPage: (limit: number, offset: number) => Promise<QueuePage>,
): Promise<CheckInReadiness> {
  const ids = new Set<number>();
  const limit = 100;
  let offset = 0;
  while (true) {
    const page = await fetchPage(limit, offset);
    if (
      !Array.isArray(page.entries) || typeof page.total !== "number" ||
      !Number.isFinite(page.total) || page.total < 0
    ) {
      throw new Error("Invalid check-in queue response");
    }
    for (const entry of page.entries) {
      if (
        typeof entry.id === "number" && entry.id > 0 &&
        entry.society_id === societyId && entry.status === "approved" &&
        !entry.checked_in_at && entry.source !== "resident_link"
      ) {
        ids.add(entry.id);
      }
    }
    offset += page.entries.length;
    if (offset >= page.total || page.entries.length === 0) break;
  }
  return { entryIds: [...ids], total: ids.size };
}

export function consumeReadyEntries(seen: Set<number>, entryIds: number[]) {
  const hasNewEntries = entryIds.some((id) => !seen.has(id));
  for (const id of entryIds) seen.add(id);
  return hasNewEntries;
}

export function readyToCheckInMessage(count: number) {
  return `${count} visitor${count === 1 ? "" : "s"} ready to check in`;
}
