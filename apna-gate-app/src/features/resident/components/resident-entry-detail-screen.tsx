import { useLocalSearchParams } from "expo-router";
import { useRef, useState } from "react";
import { ScrollView, StyleSheet, View } from "react-native";

import { Button, EmptyState, LoadingState } from "@/components/ui";
import { BottomActions } from "@/components/ui/bottom-actions";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResidentDashboard } from "@/features/resident/hooks/use-resident-dashboard";
import { useResident } from "@/features/resident/resident-context";
import { VisitorDetailsCard } from "@/features/visitors/components/visitor-details-card";
import { useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery } from "@/lib/api/resident-api-extensions";
import { spacing } from "@/theme/spacing";

function firstParam(value?: string | string[]) {
  return Array.isArray(value) ? value[0] : value;
}

export function ResidentEntryDetailScreen() {
  const { entryId: entryIdParam } = useLocalSearchParams<{
    entryId?: string | string[];
  }>();
  const entryId = Number(firstParam(entryIdParam));
  const { flatId, societyId } = useResident();
  const dashboard = useResidentDashboard();
  const actionInFlight = useRef(false);
  const [activeAction, setActiveAction] = useState<"approve" | "reject" | null>(
    null,
  );
  const canQuery = Boolean(
    societyId && flatId && Number.isFinite(entryId) && entryId > 0,
  );
  const query =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0, entryId },
      { skip: !canQuery },
    );
  const entry = query.currentData?.data?.entry;
  const canApprove =
    entry?.status === "waiting_approval" && dashboard.canManageFlatVisitors;

  const handleApprove = async () => {
    if (!entry?.id || !canApprove || actionInFlight.current) {
      return;
    }
    actionInFlight.current = true;
    setActiveAction("approve");
    try {
      await dashboard.handleApprove(entry.id);
      await query.refetch();
    } finally {
      actionInFlight.current = false;
      setActiveAction(null);
    }
  };

  const handleReject = async () => {
    if (!entry?.id || !canApprove || actionInFlight.current) {
      return;
    }
    actionInFlight.current = true;
    setActiveAction("reject");
    try {
      await dashboard.handleReject(entry.id);
      await query.refetch();
    } finally {
      actionInFlight.current = false;
      setActiveAction(null);
    }
  };

  return (
    <ResidentSubScreen
      title="Visitor Details"
      footer={
        canApprove ? (
          <BottomActions>
            <View style={styles.actionRow}>
            <Button
              title="Approve"
              loading={activeAction === "approve"}
              disabled={Boolean(activeAction)}
              onPress={() => void handleApprove()}
              fullWidth={false}
              style={styles.actionButton}
            />
            <Button
              title="Reject"
              variant="secondary"
              loading={activeAction === "reject"}
              disabled={Boolean(activeAction)}
              onPress={() => void handleReject()}
              fullWidth={false}
              style={styles.actionButton}
            />
            </View>
          </BottomActions>
        ) : undefined
      }
    >
      {query.isLoading ? (
        <LoadingState message="Loading visitor details" />
      ) : !entry ? (
        <View style={styles.emptyWrap}>
          <EmptyState
            actionLabel={canQuery ? "Retry" : undefined}
            message="This visitor entry could not be loaded."
            title="Entry not found"
            onAction={canQuery ? () => void query.refetch() : undefined}
          />
        </View>
      ) : (
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          <VisitorDetailsCard entry={entry} variant="resident" />
        </ScrollView>
      )}
    </ResidentSubScreen>
  );
}

const styles = StyleSheet.create({
  actionRow: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  actionButton: { flex: 1, borderRadius: 999 },
  emptyWrap: {
    padding: spacing.lg,
  },
  scrollContent: {
    flexGrow: 1,
    padding: spacing.lg,
  },
});
