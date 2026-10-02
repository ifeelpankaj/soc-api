import { useCallback, useState } from "react";
import { Platform, StyleSheet, TextInput, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { useRouter } from "expo-router";

import { PaginatedList } from "@/components/ui";
import { GuardConfirmDialog } from "@/features/guard/components/guard-confirm-dialog";
import { GuardEntryEditSheet } from "@/features/guard/components/guard-entry-edit-sheet";
import {
  GuardQueueEntryCard,
  GuardQueueEntryDivider,
} from "@/features/guard/components/guard-queue-entry-card";
import { GuardQueueEmptyState } from "@/features/guard/components/guard-queue-empty-state";
import { GuardSubScreen } from "@/features/guard/components/guard-sub-screen";
import { canEditVisitorEntry } from "@/features/guard/guard-entry-edit";
import { guardEntryDetailRoute } from "@/features/guard/guard-routes";
import { useGuardActions } from "@/features/guard/hooks/use-guard-actions";
import { useGuardDashboard } from "@/features/guard/hooks/use-guard-dashboard";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardPending } from "@/features/guard/hooks/use-guard-pending";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { actionForEntry } from "@/features/visitors/entry-action-state";
import type { ModelsVisitorPendingEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export default function GuardPendingScreen() {
  const router = useRouter();
  const { selectedSocietyId } = useGuardScreen();
  const feedback = useGuardFeedback();
  const [search, setSearch] = useState("");
  const pending = useGuardPending(selectedSocietyId, search);
  const actions = useGuardActions(selectedSocietyId ?? 0);
  const { visitorSettings } = useGuardDashboard();
  const allowOnBehalfApproval = visitorSettings?.allow_guard_on_behalf_approval !== false;
  const [forceEntry, setForceEntry] = useState<ModelsVisitorPendingEntry | null>(null);
  const [editEntry, setEditEntry] = useState<ModelsVisitorPendingEntry | null>(null);

  const handleNotify = useCallback(
    async (entry: ModelsVisitorPendingEntry) => {
      if (!entry.id) {
        return;
      }

      const result = await actions.notifyResident(entry.id);
      feedback.showActionResult(result, {
        successTitle: "Notification sent",
        errorTitle: "Notify failed",
      });
    },
    [actions, feedback],
  );

  const handleApproveAndCheckIn = useCallback(
    async (entry: ModelsVisitorPendingEntry, onBehalf: boolean) => {
      if (!entry.id) {
        return;
      }

      const result = await actions.approveAndCheckInEntry(entry.id, { onBehalf });
      feedback.showActionResult(result, {
        successTitle: "Checked in",
        errorTitle: "Action failed",
      });

      if (result.success) {
        setForceEntry(null);
      }
    },
    [actions, feedback],
  );

  return (
    <GuardSubScreen title="Pending Approvals">
      <PaginatedList
        ItemSeparatorComponent={GuardQueueEntryDivider}
        contentContainerStyle={styles.listContent}
        data={pending.items}
        emptyComponent={
          <GuardQueueEmptyState
            message={
              search
                ? "Try a different name, phone, flat, or vehicle."
                : "Visitors waiting for resident approval will appear here."
            }
            refreshing={pending.isRefreshing}
            title={search ? "No matching approvals" : "No pending approvals"}
            onRefresh={() => {
              void pending.refresh();
            }}
          />
        }
        footer={<View style={styles.footerSpacer} />}
        hasMore={pending.hasMore}
        header={
          <View style={styles.searchWrap}>
            <View style={styles.searchRow}>
              <SymbolView
                name={{ ios: "magnifyingglass", android: "search", web: "search" }}
                size={18}
                tintColor={colors.brand.orange}
              />
              <TextInput
                autoCapitalize="none"
                autoCorrect={false}
                placeholder="Search name, phone, flat, vehicle..."
                placeholderTextColor={colors.guard.textMuted}
                selectionColor={colors.brand.orange}
                style={[styles.searchInput, Platform.OS === "web" ? styles.searchInputWeb : null]}
                value={search}
                onChangeText={setSearch}
              />
            </View>
          </View>
        }
        isLoading={pending.isLoading}
        isLoadingMore={pending.isLoadingMore}
        isRefreshing={pending.isRefreshing}
        keyExtractor={(item) => `pending-${item.id}`}
        renderItem={({ item }) => {
          const isGuardEntry = item.source === "guard_entry";
          const activeAction = actionForEntry(actions.activeAction, item.id);
          return (
            <GuardQueueEntryCard
              disabled={actions.isLoading}
              entry={item}
              loadingAction={
                activeAction === "notify"
                  ? "tertiary"
                  : activeAction === "approve_and_check_in"
                    ? "primary"
                    : undefined
              }
              primaryActionLabel={isGuardEntry ? "Approve & Check In" : "Force Check In"}
              secondaryActionLabel={canEditVisitorEntry(item) ? "Edit Details" : undefined}
              tertiaryActionLabel="Notify Resident"
              onPress={() => {
                if (item.id) {
                  router.push(guardEntryDetailRoute(item.id));
                }
              }}
              onPrimaryAction={() => {
                if (isGuardEntry) {
                  void handleApproveAndCheckIn(item, false);
                  return;
                }
                if (!allowOnBehalfApproval) {
                  feedback.showActionResult(
                    {
                      success: false,
                      message:
                        "Your admin has not allowed guard approval on behalf of residents.",
                    },
                    { errorTitle: "Action not allowed", successTitle: "Checked in" },
                  );
                  return;
                }
                setForceEntry(item);
              }}
              onSecondaryAction={
                canEditVisitorEntry(item) ? () => setEditEntry(item) : undefined
              }
              onTertiaryAction={() => void handleNotify(item)}
            />
          );
        }}
        onLoadMore={pending.loadMore}
        onRefresh={() => {
          void pending.refresh();
        }}
      />

      <GuardConfirmDialog
        confirmLabel="Allow Visitor"
        loading={actions.isLoading}
        message={`You are allowing this visitor on behalf of ${
          forceEntry?.primary_resident_name ?? "the flat owner"
        }. This action will be recorded.`}
        title="Allow on behalf of resident"
        visible={Boolean(forceEntry)}
        onCancel={() => setForceEntry(null)}
        onConfirm={() => forceEntry && void handleApproveAndCheckIn(forceEntry, true)}
      />

      <GuardEntryEditSheet
        entry={editEntry}
        societyId={selectedSocietyId ?? 0}
        visible={Boolean(editEntry)}
        onClose={() => setEditEntry(null)}
        onSaved={() => {
          setEditEntry(null);
          void pending.refresh();
          feedback.showSuccess("Details updated", "Visitor information saved.");
        }}
      />
    </GuardSubScreen>
  );
}

const styles = StyleSheet.create({
  footerSpacer: {
    height: spacing.sm,
  },
  listContent: {
    gap: spacing.xs,
    paddingBottom: spacing["2xl"],
    paddingHorizontal: 0,
    paddingTop: spacing.md,
  },
  searchInput: {
    backgroundColor: "transparent",
    borderWidth: 0,
    color: colors.guard.text,
    flex: 1,
    fontSize: 14,
    marginLeft: spacing.sm,
    minHeight: 40,
    paddingVertical: 0,
  },
  searchInputWeb: {
    outlineWidth: 0,
  },
  searchRow: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: "#F1E4DA",
    borderRadius: radius.xl,
    borderWidth: 1,
    flexDirection: "row",
    minHeight: 48,
    paddingHorizontal: spacing.lg,
    ...shadows.sm,
  },
  searchWrap: {
    paddingHorizontal: layout.screenPaddingHorizontal,
  },
});
