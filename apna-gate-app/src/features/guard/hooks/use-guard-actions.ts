import { useCallback, useState } from "react";

import { getFriendlyApiMessage } from "@/features/auth/api-error";
import type { EntryActionState } from "@/features/visitors/entry-action-state";
import {
  type ModelsGuardApproveEntryRequest,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyMutation,
} from "@/lib/api/generated-api";

export type GuardActionOptions = {
  onBehalf?: boolean;
  reason?: string;
};

export type GuardAction =
  | "notify"
  | "approve"
  | "approve_and_check_in"
  | "check_in"
  | "check_out";

function toApproveRequest(opts?: GuardActionOptions): ModelsGuardApproveEntryRequest {
  if (!opts?.onBehalf && !opts?.reason) {
    return {};
  }
  return { on_behalf: opts.onBehalf, reason: opts.reason };
}

export function useGuardActions(societyId: number) {
  const [activeAction, setActiveAction] = useState<EntryActionState<GuardAction>>();
  const [notify, notifyState] = usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyMutation();
  const [approve, approveState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveMutation();
  const [approveAndCheckIn, approveAndCheckInState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInMutation();
  const [checkIn, checkInState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInMutation();
  const [checkOut, checkOutState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutMutation();

  const notifyResident = useCallback(
    async (entryId: number) => {
      setActiveAction({ entryId, action: "notify" });
      try {
        await notify({ societyId, entryId }).unwrap();
        return { success: true as const, message: "Notification sent to resident." };
      } catch (error) {
        return { success: false as const, message: getFriendlyApiMessage(error, "Could not notify resident") };
      } finally {
        setActiveAction(undefined);
      }
    },
    [notify, societyId],
  );

  const approveEntry = useCallback(
    async (entryId: number, opts?: GuardActionOptions) => {
      setActiveAction({ entryId, action: "approve" });
      try {
        const response = await approve({
          societyId,
          entryId,
          modelsGuardApproveEntryRequest: toApproveRequest(opts),
        }).unwrap();
        return {
          success: true as const,
          message: response.message ?? "Visitor approved",
          entry: response.data?.entry,
        };
      } catch (error) {
        return { success: false as const, message: getFriendlyApiMessage(error, "Could not approve visitor") };
      } finally {
        setActiveAction(undefined);
      }
    },
    [approve, societyId],
  );

  const approveAndCheckInEntry = useCallback(
    async (entryId: number, opts?: GuardActionOptions) => {
      setActiveAction({ entryId, action: "approve_and_check_in" });
      try {
        const response = await approveAndCheckIn({
          societyId,
          entryId,
          modelsGuardApproveEntryRequest: toApproveRequest(opts),
        }).unwrap();
        return {
          success: true as const,
          message: response.message ?? "Visitor checked in",
          entry: response.data?.entry,
        };
      } catch (error) {
        return {
          success: false as const,
          message: getFriendlyApiMessage(error, "Could not approve and check in visitor"),
        };
      } finally {
        setActiveAction(undefined);
      }
    },
    [approveAndCheckIn, societyId],
  );

  const checkInEntry = useCallback(
    async (entryId: number) => {
      setActiveAction({ entryId, action: "check_in" });
      try {
        const response = await checkIn({ societyId, entryId }).unwrap();
        return {
          success: true as const,
          message: response.message ?? "Visitor checked in",
          entry: response.data?.entry,
        };
      } catch (error) {
        return { success: false as const, message: getFriendlyApiMessage(error, "Could not check in visitor") };
      } finally {
        setActiveAction(undefined);
      }
    },
    [checkIn, societyId],
  );

  const checkOutEntry = useCallback(
    async (entryId: number) => {
      setActiveAction({ entryId, action: "check_out" });
      try {
        const response = await checkOut({ societyId, entryId }).unwrap();
        return {
          success: true as const,
          message: response.message ?? "Visitor checked out",
          entry: response.data?.entry,
        };
      } catch (error) {
        return { success: false as const, message: getFriendlyApiMessage(error, "Could not check out visitor") };
      } finally {
        setActiveAction(undefined);
      }
    },
    [checkOut, societyId],
  );

  const isLoading =
    notifyState.isLoading ||
    approveState.isLoading ||
    approveAndCheckInState.isLoading ||
    checkInState.isLoading ||
    checkOutState.isLoading;

  return {
    activeAction,
    activeEntryId: activeAction?.entryId,
    approveAndCheckInEntry,
    approveEntry,
    checkInEntry,
    checkOutEntry,
    isLoading,
    notifyResident,
  };
}
