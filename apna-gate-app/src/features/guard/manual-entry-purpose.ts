import type { ModelsVisitorPurpose } from "@/lib/api/generated-api";

export function isSocietyWidePurpose(purpose: ModelsVisitorPurpose) {
  return purpose === "staff" || purpose === "service";
}

export function flatForManualEntryPurpose<T>(
  purpose: ModelsVisitorPurpose,
  flat: T | null,
): T | null {
  return isSocietyWidePurpose(purpose) ? null : flat;
}

export function manualEntryFlatId(
  purpose: ModelsVisitorPurpose,
  flat: { id: number } | null,
) {
  return flatForManualEntryPurpose(purpose, flat)?.id;
}
