import { useCallback, useMemo, useRef, useState } from "react";

import { getApiMessage, isSubscriptionError } from "@/features/auth/api-error";
import { titleize } from "@/features/guard/guard-utils";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { useResident } from "@/features/resident/resident-context";
import { getTodayRange } from "@/features/visitors/visitor-date-ranges";
import {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery,
  useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery,
} from "@/lib/api/resident-api-extensions";
import {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingQuery,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectMutation,
  type ModelsFlatRecentVisitorSummary,
  type ModelsVisitorEntry,
} from "@/lib/api/generated-api";

export function useResidentDashboard() {
  const {
    canManageFlatMembers,
    canManageFlatVisitors,
    flatId,
    refetch: refetchBootstrap,
    selectedResidence,
    societyId,
    user,
  } = useResident();
  const { showError, showSuccess } = useResidentFeedback();
  const [actionEntryId, setActionEntryId] = useState<number | null>(null);
  const actionInFlight = useRef(false);
  const [fetchEntryDetail, fetchEntryState] =
    useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery();

  const shouldSkip = !societyId || !flatId;
  const queryOpts = { skip: shouldSkip };

  const contextQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0 },
      queryOpts,
    );
  const pendingQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0 },
      queryOpts,
    );
  const expectedRange = getTodayRange();
  const approvedQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery(
      {
        societyId: societyId ?? 0,
        flatId: flatId ?? 0,
        limit: 50,
        offset: 0,
        status: "approved",
        event: expectedRange ? "expected" : undefined,
        eventFrom: expectedRange?.eventFrom,
        eventTo: expectedRange?.eventTo,
      },
      queryOpts,
    );

  const [approveEntry, approveState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveMutation();
  const [rejectEntry, rejectState] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectMutation();

  const queries = [contextQuery, pendingQuery, approvedQuery];
  const failedQuery = queries.find((query) => query.isError);
  const failedError = failedQuery?.error ?? null;
  const isSubscriptionBlocked = isSubscriptionError(failedError);
  const hasError = !!failedQuery && !isSubscriptionBlocked;

  const { refetch: refetchContext } = contextQuery;
  const { refetch: refetchPending } = pendingQuery;
  const { refetch: refetchApproved } = approvedQuery;

  const refetchAll = useCallback(() => {
    refetchBootstrap();

    if (shouldSkip) {
      return;
    }

    if (!contextQuery.isUninitialized) {
      void refetchContext();
    }
    if (!pendingQuery.isUninitialized) {
      void refetchPending();
    }
    if (!approvedQuery.isUninitialized) {
      void refetchApproved();
    }
  }, [
    approvedQuery.isUninitialized,
    contextQuery.isUninitialized,
    pendingQuery.isUninitialized,
    refetchApproved,
    refetchBootstrap,
    refetchContext,
    refetchPending,
    shouldSkip,
  ]);

  const context = contextQuery.data?.data?.context;
  const pendingEntries = pendingQuery.data?.data?.entries ?? [];
  const pendingCount = pendingQuery.data?.data?.total ?? pendingEntries.length;
  const expectedCount =
    approvedQuery.data?.data?.total ??
    approvedQuery.data?.data?.entries?.length ??
    0;

  const flatLabel = useMemo(() => {
    const parts = [
      selectedResidence?.flat_number
        ? `Flat ${selectedResidence.flat_number}`
        : null,
      selectedResidence?.role ? titleize(selectedResidence.role) : null,
    ].filter(Boolean);

    return parts.join(" • ") || "Your flat";
  }, [selectedResidence]);

  const displayName = user?.full_name ?? user?.first_name ?? "Resident";
  const membersCount = context?.total_residents ?? 0;
  const recentVisitors: ModelsFlatRecentVisitorSummary[] =
    context?.recent_visitors ?? [];
  const visitorsCount = recentVisitors.length;

  const fetchEntryDetailById = useCallback(
    async (entryId: number): Promise<ModelsVisitorEntry | null> => {
      if (!societyId || !flatId) {
        return null;
      }

      try {
        const response = await fetchEntryDetail({
          societyId,
          flatId,
          entryId,
        }).unwrap();
        return response.data?.entry ?? null;
      } catch (error) {
        showError("Unable to load entry", error, "Please try again.");
        return null;
      }
    },
    [fetchEntryDetail, flatId, showError, societyId],
  );

  const handleApprove = useCallback(
    async (entryId?: number) => {
      if (
        !entryId ||
        !societyId ||
        !canManageFlatVisitors ||
        actionInFlight.current
      ) {
        return false;
      }
      actionInFlight.current = true;
      setActionEntryId(entryId);

      try {
        await approveEntry({ societyId, entryId }).unwrap();
        showSuccess("Approved", "Visitor can now enter the society.");
        refetchAll();
        return true;
      } catch (error) {
        showError("Approval failed", error, "Please try again.");
        return false;
      } finally {
        actionInFlight.current = false;
        setActionEntryId(null);
      }
    },
    [
      approveEntry,
      canManageFlatVisitors,
      refetchAll,
      showError,
      showSuccess,
      societyId,
    ],
  );

  const handleReject = useCallback(
    async (entryId?: number) => {
      if (
        !entryId ||
        !societyId ||
        !canManageFlatVisitors ||
        actionInFlight.current
      ) {
        return false;
      }
      actionInFlight.current = true;
      setActionEntryId(entryId);

      try {
        await rejectEntry({
          societyId,
          entryId,
          modelsRejectVisitorEntryRequest: { reason: "Declined by resident" },
        }).unwrap();
        showSuccess("Rejected", "Visitor entry was declined.");
        refetchAll();
        return true;
      } catch (error) {
        showError("Rejection failed", error, "Please try again.");
        return false;
      } finally {
        actionInFlight.current = false;
        setActionEntryId(null);
      }
    },
    [
      canManageFlatVisitors,
      refetchAll,
      rejectEntry,
      showError,
      showSuccess,
      societyId,
    ],
  );

  const isInitialLoading =
    shouldSkip || queries.some((query) => query.isLoading && !query.data);

  const approvalMode = context?.society_approval_mode ?? "mandatory";
  const isHybrid = approvalMode === "hybrid";

  return {
    actionEntryId,
    approvalMode,
    displayName,
    errorMessage: hasError
      ? getApiMessage(failedError, "Unable to load resident home data.")
      : null,
    expectedCount,
    fetchEntryDetail: fetchEntryDetailById,
    flatLabel,
    handleApprove,
    handleReject,
    hasError,
    isSubscriptionBlocked,
    isActionLoading: approveState.isLoading || rejectState.isLoading,
    loadingAction: approveState.isLoading
      ? ("primary" as const)
      : rejectState.isLoading
        ? ("secondary" as const)
        : undefined,
    isDetailLoading: fetchEntryState.isFetching,
    isHybrid,
    isInitialLoading,
    isReady: !shouldSkip,
    canManageFlatMembers,
    canManageFlatVisitors,
    isRefreshing: queries.some((query) => query.isFetching && !query.isLoading),
    membersCount,
    pendingCount,
    pendingEntries,
    primaryResidentName: context?.primary_resident?.full_name,
    recentVisitors,
    refetchAll,
    selectedResidence,
    visitorsCount,
  };
}
