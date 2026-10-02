import { useRouter } from "expo-router";
import { SymbolView } from "expo-symbols";
import { useRef, useState } from "react";
import {
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";

import { Button, EmptyState, PaginatedList } from "@/components/ui";
import { getVisitorActionErrorMessage } from "@/features/auth/api-error";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { filterPendingQueue } from "@/features/resident/pending-queue";
import { useResident } from "@/features/resident/resident-context";
import {
  residentEntryDetailRoute,
  residentVisitorInviteRoute,
  residentVisitorSettingsRoute,
} from "@/features/resident/resident-routes";
import { VisitorQueueEntryCard } from "@/features/visitors/components/visitor-queue-entry-card";
import {
  actionForEntry,
  type EntryActionState,
} from "@/features/visitors/entry-action-state";
import { titleize, visitorPurposes } from "@/features/visitors/visitor-utils";
import {
  generatedApi,
  type ModelsVisitorPurpose,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectMutation,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

export default function ResidentVisitorsScreen() {
  const { societyId, flatId } = useResident();
  return <ResidentPendingQueue key={`${societyId}:${flatId}`} />;
}

function ResidentPendingQueue() {
  const router = useRouter();
  const feedback = useResidentFeedback();
  const { societyId, flatId, canManageFlatVisitors } = useResident();
  const [search, setSearch] = useState("");
  const [purpose, setPurpose] = useState<ModelsVisitorPurpose | "all">("all");
  const [visibleCount, setVisibleCount] = useState(15);
  const [activeAction, setActiveAction] =
    useState<EntryActionState<"approve" | "reject">>();
  const actionInFlight = useRef(false);
  const [approve] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveMutation();
  const [reject] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectMutation();
  const query =
    generatedApi.endpoints.getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPending.useQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0 },
      { skip: !societyId || !flatId, refetchOnMountOrArgChange: true },
    );
  const entries = query.currentData?.data?.entries ?? [];
  const matching = filterPendingQueue(entries, search, purpose);
  const searched = filterPendingQueue(entries, search, "all");
  const filters = ["all", ...visitorPurposes] as const;

  const decide = async (
    entryId: number | undefined,
    action: "approve" | "reject",
  ) => {
    if (
      !entryId ||
      !societyId ||
      !canManageFlatVisitors ||
      actionInFlight.current
    )
      return;
    actionInFlight.current = true;
    setActiveAction({ entryId, action });
    try {
      if (action === "approve") await approve({ societyId, entryId }).unwrap();
      else
        await reject({
          societyId,
          entryId,
          modelsRejectVisitorEntryRequest: { reason: "Declined by resident" },
        }).unwrap();
      feedback.showSuccess(
        action === "approve" ? "Visitor approved" : "Visitor declined",
        action === "approve"
          ? "The guard can now check them in."
          : "The visitor request was declined.",
      );
      await query.refetch();
    } catch (error) {
      feedback.showError(
        "Could not update visitor",
        error,
        getVisitorActionErrorMessage(error, "Please try again."),
      );
    } finally {
      actionInFlight.current = false;
      setActiveAction(undefined);
    }
  };

  return (
    <ResidentSubScreen
      title="Pending approvals"
      headerTrailing={
        <Pressable
          accessibilityRole="button"
          accessibilityLabel="Visitor settings"
          style={styles.iconButton}
          onPress={() => router.push(residentVisitorSettingsRoute())}
        >
          <SymbolView
            name={{ ios: "slider.horizontal.3", android: "tune", web: "tune" }}
            size={21}
            tintColor={colors.text.primary}
          />
        </Pressable>
      }
    >
      <PaginatedList
        data={matching.slice(0, visibleCount)}
        contentContainerStyle={styles.list}
        keyExtractor={(entry) => String(entry.id)}
        isLoading={query.isLoading}
        isRefreshing={query.isFetching && !query.isLoading}
        hasMore={matching.length > visibleCount}
        onLoadMore={() => setVisibleCount((count) => count + 15)}
        onRefresh={() => {
          void query.refetch();
        }}
        header={
          <View style={styles.header}>
            <View style={styles.headingRow}>
              <View style={styles.copy}>
                <Text style={styles.count}>{entries.length} waiting</Text>
                <Text style={styles.subtitle}>
                  {canManageFlatVisitors
                    ? "Review requests before visitors enter."
                    : "You have view-only access to this queue."}
                </Text>
              </View>
              {canManageFlatVisitors ? (
                <Button
                  compact
                  fullWidth={false}
                  title="Invite visitor"
                  variant="secondary"
                  onPress={() => router.push(residentVisitorInviteRoute())}
                />
              ) : null}
            </View>
            <View style={styles.search}>
              <SymbolView
                name={{
                  ios: "magnifyingglass",
                  android: "search",
                  web: "search",
                }}
                size={20}
                tintColor={colors.text.muted}
              />
              <TextInput
                accessibilityLabel="Search pending visitors"
                placeholder="Search name, phone or provider"
                placeholderTextColor={colors.text.muted}
                autoCapitalize="none"
                autoCorrect={false}
                style={styles.searchInput}
                value={search}
                onChangeText={(value) => {
                  setSearch(value);
                  setVisibleCount(15);
                }}
              />
            </View>
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              contentContainerStyle={styles.filters}
            >
              {filters.map((option) => {
                const count =
                  option === "all"
                    ? searched.length
                    : searched.filter((entry) => entry.purpose === option)
                        .length;
                if (option !== "all" && count === 0 && purpose !== option)
                  return null;
                return (
                  <Pressable
                    key={option}
                    accessibilityRole="button"
                    accessibilityState={{ selected: purpose === option }}
                    style={[
                      styles.filter,
                      purpose === option && styles.selected,
                    ]}
                    onPress={() => {
                      setPurpose(option);
                      setVisibleCount(15);
                    }}
                  >
                    <Text
                      style={[
                        styles.filterText,
                        purpose === option && styles.selectedText,
                      ]}
                    >
                      {titleize(option)} · {count}
                    </Text>
                  </Pressable>
                );
              })}
            </ScrollView>
            {query.isError && entries.length > 0 ? (
              <Text accessibilityRole="alert" style={styles.error}>
                Could not refresh. Pull down to try again.
              </Text>
            ) : null}
          </View>
        }
        emptyComponent={
          <EmptyState
            title={
              query.isError
                ? "Could not load approvals"
                : search || purpose !== "all"
                  ? "No matching visitors"
                  : "You're all caught up"
            }
            message={
              query.isError
                ? "Check your connection and try again."
                : search || purpose !== "all"
                  ? "Try another name or purpose."
                  : "New visitor requests will appear here."
            }
            actionLabel={
              query.isError
                ? "Try again"
                : search || purpose !== "all"
                  ? "Clear filters"
                  : "Refresh"
            }
            onAction={() => {
              if (!query.isError && (search || purpose !== "all")) {
                setSearch("");
                setPurpose("all");
              } else void query.refetch();
            }}
          />
        }
        renderItem={({ item }) => (
          <VisitorQueueEntryCard
            variant="resident"
            entry={item}
            disabled={Boolean(activeAction)}
            primaryActionLabel={canManageFlatVisitors ? "Approve" : undefined}
            secondaryActionLabel={canManageFlatVisitors ? "Reject" : undefined}
            loadingAction={
              actionForEntry(activeAction, item.id) === "approve"
                ? "primary"
                : actionForEntry(activeAction, item.id) === "reject"
                  ? "secondary"
                  : undefined
            }
            onPrimaryAction={() => void decide(item.id, "approve")}
            onSecondaryAction={() => void decide(item.id, "reject")}
            onPress={
              item.id
                ? () => router.push(residentEntryDetailRoute(item.id!))
                : undefined
            }
          />
        )}
      />
    </ResidentSubScreen>
  );
}

const styles = StyleSheet.create({
  list: {
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
    paddingBottom: spacing.xl,
    gap: spacing.md,
  },
  header: { gap: spacing.lg },
  headingRow: {
    flexDirection: "row",
    flexWrap: "wrap",
    alignItems: "center",
    gap: spacing.md,
  },
  copy: { flex: 1, minWidth: 150, gap: spacing.xs },
  count: { ...typography.title, color: colors.text.primary },
  subtitle: { ...typography.bodySmall, color: colors.text.secondary },
  iconButton: {
    minWidth: 44,
    minHeight: 44,
    alignItems: "center",
    justifyContent: "center",
  },
  search: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderWidth: 1,
    borderRadius: radius.lg,
    paddingHorizontal: spacing.md,
  },
  searchInput: {
    flex: 1,
    minWidth: 0,
    minHeight: 52,
    paddingVertical: spacing.md,
    ...typography.bodySmall,
    color: colors.text.primary,
  },
  filters: { gap: spacing.sm },
  filter: {
    minHeight: 44,
    justifyContent: "center",
    borderWidth: 1,
    borderColor: colors.border.default,
    borderRadius: radius.full,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
    backgroundColor: colors.surface.card,
  },
  selected: {
    backgroundColor: colors.brand.navy,
    borderColor: colors.brand.navy,
  },
  filterText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    fontWeight: "600",
  },
  selectedText: { color: colors.text.inverse },
  error: { ...typography.bodySmall, color: colors.status.error },
});
