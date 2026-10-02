import type {
  ModelsVisitorEntry,
  ModelsVisitorPurpose,
} from "@/lib/api/generated-api";
import {
  getFlatLabel,
  getVisitorName,
  getVisitorPhone,
} from "@/features/visitors/visitor-utils";

export function filterPendingQueue(
  entries: ModelsVisitorEntry[],
  search: string,
  purpose: ModelsVisitorPurpose | "all",
) {
  const term = search.trim().toLocaleLowerCase();
  return entries.filter(
    (entry) =>
      (purpose === "all" || entry.purpose === purpose) &&
      (!term ||
        [
          getVisitorName(entry),
          getVisitorPhone(entry),
          getFlatLabel(entry),
          entry.service_provider,
          entry.vehicle_number,
        ]
          .filter(Boolean)
          .join(" ")
          .toLocaleLowerCase()
          .includes(term)),
  );
}
