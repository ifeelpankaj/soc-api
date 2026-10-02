import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import type { GuardScanOutcome } from "@/features/guard/hooks/use-guard-check-in";

export function checkInFooterState(outcome: GuardScanOutcome, status?: ModelsVisitorEntry["status"], isCheckingIn = false) {
  if (outcome === "just_checked_in" || outcome === "already_inside" || status === "checked_in") return "scan";
  return (outcome === "ready" || isCheckingIn) && status === "approved" ? "check_in" : null;
}

export function visitorDetailsFooterState(status?: ModelsVisitorEntry["status"]) {
  if (status === "approved") return "check_in";
  if (status === "checked_in") return "check_out";
  return null;
}
