import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import type { GuardScanOutcome, GuardCheckInError } from "@/features/guard/hooks/use-guard-check-in";

export function deriveScanOutcome(input: {
  entry: ModelsVisitorEntry | null;
  entryError: GuardCheckInError | null;
  isLoadingEntry: boolean;
  justCheckedInThisSession: boolean;
  canCheckIn: boolean;
  alreadyCheckedIn: boolean;
  isResolvingEntry: boolean;
}): GuardScanOutcome {
  const {
    entry,
    entryError,
    isLoadingEntry,
    justCheckedInThisSession,
    canCheckIn,
    alreadyCheckedIn,
    isResolvingEntry,
  } = input;

  // A refresh of cached details must never undo a completed check-in in the UI.
  if (justCheckedInThisSession) return "just_checked_in";
  if (alreadyCheckedIn || entry?.status === "checked_in") return "already_inside";

  if (isLoadingEntry || isResolvingEntry) {
    return "loading";
  }

  if (entryError && !entry) {
    return "error";
  }

  if (!entry) {
    return "error";
  }

  if (entry.status === "waiting_approval") {
    return "pending_approval";
  }

  if (canCheckIn) {
    return "ready";
  }

  return "blocked";
}

