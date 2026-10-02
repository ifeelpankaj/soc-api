import { useCallback, useState } from "react";
import { Platform, StyleSheet, TextInput, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { useRouter } from "expo-router";

import { PaginatedList } from "@/components/ui";
import { GuardEntryEditSheet } from "@/features/guard/components/guard-entry-edit-sheet";
import {
  GuardQueueEntryCard,
  GuardQueueEntryDivider,
} from "@/features/guard/components/guard-queue-entry-card";
import { GuardQueueEmptyState } from "@/features/guard/components/guard-queue-empty-state";
import { GuardSubScreen } from "@/features/guard/components/guard-sub-screen";
import { canEditVisitorEntry } from "@/features/guard/guard-entry-edit";
import { guardCheckInRoute } from "@/features/guard/guard-routes";
import { getWaitingDuration } from "@/features/guard/guard-utils";
import { useGuardActions } from "@/features/guard/hooks/use-guard-actions";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { useGuardWaitingAtGate } from "@/features/guard/hooks/use-guard-waiting-at-gate";
import { actionForEntry } from "@/features/visitors/entry-action-state";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export default function GuardWaitingAtGateScreen() {
  const router = useRouter();
  const { selectedSocietyId } = useGuardScreen();
  const feedback = useGuardFeedback();
  const [search, setSearch] = useState("");
  const queue = useGuardWaitingAtGate(selectedSocietyId, search);
  const actions = useGuardActions(selectedSocietyId ?? 0);
  const [editEntry, setEditEntry] = useState<ModelsVisitorEntry | null>(null);

  const handleCheckIn = useCallback(
    async (entry: ModelsVisitorEntry) => {
      if (!entry.id) {
        return;
      }

      const result = await actions.checkInEntry(entry.id);
      feedback.showActionResult(result, {
        successTitle: "Checked in",
        errorTitle: "Check-in failed",
      });

      if (result.success) {
        void queue.refresh();
      }
    },
    [actions, feedback, queue],
  );

  return (
    <GuardSubScreen title="Waiting at Gate">
      <PaginatedList
        ItemSeparatorComponent={GuardQueueEntryDivider}
        contentContainerStyle={styles.listContent}
        data={queue.items}
        emptyComponent={
          <GuardQueueEmptyState
            message={
              search
                ? "Try a different name, phone, flat, or vehicle."
                : "No approved visitors waiting for manual check-in (walk-ins and guard-approved entries only)."
            }
            refreshing={queue.isRefreshing}
            title={search ? "No matching visitors" : "No visitors waiting at gate"}
            onRefresh={() => {
              void queue.refresh();
            }}
          />
        }
        footer={<View style={styles.footerSpacer} />}
        hasMore={queue.hasMore}
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
        isLoading={queue.isLoading}
        isLoadingMore={queue.isLoadingMore}
        isRefreshing={queue.isRefreshing}
        keyExtractor={(item) => `waiting-${item.id}`}
        renderItem={({ item }) => {
          const waiting = getWaitingDuration(item.approved_at);
          const activeAction = actionForEntry(actions.activeAction, item.id);
          return (
            <GuardQueueEntryCard
              disabled={actions.isLoading}
              entry={item}
              loadingAction={activeAction === "check_in" ? "primary" : undefined}
              primaryActionLabel="Check In"
              secondaryActionLabel={canEditVisitorEntry(item) ? "Edit Details" : undefined}
              waitingLabel={waiting.label}
              waitingTone={waiting.tone}
              onPress={() => {
                if (item.id) {
                  router.push(guardCheckInRoute({ source: "entry", entryId: item.id }));
                }
              }}
              onPrimaryAction={() => void handleCheckIn(item)}
              onSecondaryAction={
                canEditVisitorEntry(item) ? () => setEditEntry(item) : undefined
              }
            />
          );
        }}
        onLoadMore={queue.loadMore}
        onRefresh={() => {
          void queue.refresh();
        }}
      />

      <GuardEntryEditSheet
        entry={editEntry}
        societyId={selectedSocietyId ?? 0}
        visible={Boolean(editEntry)}
        onClose={() => setEditEntry(null)}
        onSaved={() => {
          setEditEntry(null);
          void queue.refresh();
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
