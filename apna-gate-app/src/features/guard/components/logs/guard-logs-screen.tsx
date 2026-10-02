import { useLocalSearchParams, useRouter } from "expo-router";

import { GuardSocietyGate } from "@/features/guard/components/guard-society-gate";
import { guardEntryDetailRoute, getLogsSegmentSummaryTitle } from "@/features/guard/guard-routes";
import { useGuardActions } from "@/features/guard/hooks/use-guard-actions";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { useGuardLogsFromParams } from "@/features/guard/hooks/use-guard-logs";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { LogEntryRow } from "@/features/visitors/components/log-entry-row";
import { VisitorEntriesListScreen } from "@/features/visitors/components/visitor-entries-list-screen";
import { VisitorLogSummaryStrip } from "@/features/visitors/components/visitor-log-summary-strip";
import { useGetV1SocietiesBySocietyIdVisitorEntriesStatsExtendedQuery } from "@/lib/api/guard-api-extensions";

export function GuardLogsScreen() {
  const router = useRouter();
  const { preset: presetParam } = useLocalSearchParams<{ preset?: string | string[] }>();
  const { isLoading, isReady, memberships, requiresSelection, selectedSocietyId } = useGuardScreen();
  const feedback = useAppFeedback();
  const actions = useGuardActions(selectedSocietyId ?? 0);
  const logs = useGuardLogsFromParams(presetParam);
  const statsQuery = useGetV1SocietiesBySocietyIdVisitorEntriesStatsExtendedQuery(
    {
      societyId: selectedSocietyId ?? 0,
      eventFrom: logs.statsRange?.eventFrom,
      eventTo: logs.statsRange?.eventTo,
    },
    { skip: !selectedSocietyId },
  );

  if (!isLoading && (memberships.length === 0 || requiresSelection || !selectedSocietyId)) {
    return <GuardSocietyGate />;
  }

  const stats = statsQuery.data?.data?.stats;
  const checkedOutCount = stats?.checked_out_in_range ?? stats?.checked_out_today ?? 0;
  const listLoading = logs.isLoading || (isLoading && !isReady);

  const handleCheckOut = async (entryId?: number) => {
    if (!entryId || !selectedSocietyId) {
      feedback.showActionResult(
        { success: false, message: "This visitor entry cannot be checked out." },
        { errorTitle: "Missing entry", successTitle: "Checked out" },
      );
      return;
    }

    const result = await actions.checkOutEntry(entryId);
    feedback.showActionResult(result, {
      successTitle: "Checked out",
      errorTitle: "Checkout failed",
    });

    if (result.success) {
      void logs.refresh();
      void statsQuery.refetch();
    }
  };

  const showingCustomResults = logs.isSearchActive || logs.hasExplicitFilters;
  const emptyTitle = showingCustomResults ? "No matching visitors" : "No visitor logs yet";
  const emptyMessage = showingCustomResults
    ? "Try a different name, phone, flat, or block."
    : "Visitor movement will appear here as entries are created.";

  return (
    <VisitorEntriesListScreen
      controls={{
        ...logs,
        isLoading: listLoading,
        refresh: () => {
          void logs.refresh();
          void statsQuery.refetch();
        },
      }}
      emptyMessage={emptyMessage}
      emptyTitle={emptyTitle}
      headerExtra={
        !showingCustomResults ? (
          <VisitorLogSummaryStrip
            checkedIn={stats?.visitors_inside ?? 0}
            checkedOut={checkedOutCount}
            pending={stats?.pending_approvals ?? 0}
            title={getLogsSegmentSummaryTitle(logs.segment)}
          />
        ) : null
      }
      keyPrefix="log"
      renderRow={(item) => (
        <LogEntryRow
          entry={item}
          isCheckingOut={actions.activeEntryId === item.id}
          onPress={() => {
            if (item.id) {
              router.push(guardEntryDetailRoute(item.id));
            }
          }}
          onCheckOut={
            item.status === "checked_in" ? () => void handleCheckOut(item.id) : undefined
          }
        />
      )}
    />
  );
}
