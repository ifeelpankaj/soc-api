import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import type { SelectedFlat } from "@/features/guard/hooks/use-guard-manual-entry";

export function canEditVisitorEntry(entry?: ModelsVisitorEntry | null) {
  return entry?.status === "waiting_approval" || entry?.status === "approved";
}

export function canEditVisitorFlat(entry?: ModelsVisitorEntry | null) {
  return (
    canEditVisitorEntry(entry) &&
    entry?.source === "guard_entry" &&
    entry?.purpose !== "staff" &&
    !(entry?.purpose === "service" && !entry.flat_id)
  );
}

export function selectedFlatFromEntry(entry?: ModelsVisitorEntry | null): SelectedFlat | null {
  if (!entry?.flat_id) {
    return null;
  }

  return {
    id: entry.flat_id,
    block: entry.flat?.block,
    flat_number: entry.flat?.flat_number,
    floor: entry.flat?.floor,
  };
}

export function getCheckInSessionKey(input: { source: "qr"; token: string } | { source: "entry"; entryId: number; token?: string }) {
  return input.source === "entry" ? `entry:${input.entryId}` : `qr:${input.token}`;
}
