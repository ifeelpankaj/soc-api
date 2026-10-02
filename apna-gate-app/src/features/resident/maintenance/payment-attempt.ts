import { nanoid } from "@reduxjs/toolkit";
import type { ClaimSubmission } from "@/lib/api/maintenance-api";

export type SubmissionAttempt = { key: string; body: ClaimSubmission };
export function createSubmissionAttempt(
  body: ClaimSubmission,
): SubmissionAttempt {
  return { key: nanoid(), body: { ...body } };
}
export function isDefiniteRejection(error: unknown) {
  if (!error || typeof error !== "object" || !("status" in error)) return false;
  return (
    typeof error.status === "number" &&
    error.status >= 400 &&
    error.status < 500 &&
    error.status !== 408 &&
    error.status !== 429
  );
}
