import { ActivityIndicator, Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { DashboardSection } from "@/components/dashboard/dashboard-section";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";
import {
  formatActivityTimestamp,
  getFlatLabel,
  getVisitorName,
  getVisitorStatusMeta,
  titleize,
} from "@/features/guard/guard-utils";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type DashboardActivityFeedEmptyAction = {
  accessibilityLabel?: string;
  label: string;
  onPress: () => void;
};

type DashboardActivityFeedProps = {
  emptyAction?: DashboardActivityFeedEmptyAction;
  hasMore: boolean;
  isLoading: boolean;
  isLoadingMore: boolean;
  items: ModelsVisitorEntry[];
  onItemPress?: (entry: ModelsVisitorEntry) => void;
  onLoadMore: () => void;
  onViewAll: () => void;
  showFlat?: boolean;
  title?: string;
};

function ActivityEmptyIllustration() {
  return (
    <View style={styles.illustrationWrap}>
      <View style={styles.clipboardBody}>
        <SymbolView
          name={{ ios: "list.clipboard", android: "assignment", web: "assignment" }}
          size={32}
          tintColor={colors.brand.orange}
        />
      </View>
      <View style={styles.clockBadge}>
        <SymbolView
          name={{ ios: "clock.fill", android: "schedule", web: "schedule" }}
          size={14}
          tintColor={colors.brand.orange}
        />
      </View>
    </View>
  );
}

function ActivityPreviewCard({
  item,
  onPress,
  showFlat,
}: {
  item: ModelsVisitorEntry;
  onPress?: (entry: ModelsVisitorEntry) => void;
  showFlat: boolean;
}) {
  const name = getVisitorName(item);
  const statusMeta = getVisitorStatusMeta(item.status);
  const timestamp = formatActivityTimestamp(
    item.checked_out_at ?? item.checked_in_at ?? item.expected_at ?? item.created_at,
  );
  const purpose = item.purpose ? titleize(item.purpose) : "Visitor";
  const flatLabel = showFlat ? getFlatLabel(item) : null;
  const canPress = Boolean(onPress && item.id);
  const subtitle = [purpose, flatLabel, timestamp].filter(Boolean).join(" · ");

  const content = (
    <View style={styles.previewCard}>
      <VisitorPhoto entry={item} size={44} />
      <View style={styles.previewCopy}>
        <Text numberOfLines={1} style={styles.previewName}>
          {name}
        </Text>
        <Text numberOfLines={2} style={styles.previewMeta}>
          {subtitle}
        </Text>
      </View>
      <View
        style={[
          styles.previewBadge,
          { backgroundColor: statusMeta.bg, borderColor: statusMeta.border },
        ]}
      >
        <Text style={[styles.previewBadgeText, { color: statusMeta.color }]}>
          {statusMeta.label}
        </Text>
      </View>
      {canPress ? (
        <SymbolView
          name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
          size={14}
          tintColor={colors.guard.textMuted}
        />
      ) : null}
    </View>
  );

  if (!canPress) {
    return content;
  }

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [pressed && styles.previewPressed]}
      onPress={() => onPress?.(item)}
    >
      {content}
    </Pressable>
  );
}

export function DashboardActivityFeed({
  emptyAction,
  hasMore,
  isLoading,
  isLoadingMore,
  items,
  onItemPress,
  onLoadMore,
  onViewAll,
  showFlat = false,
  title = "Today's Activity",
}: DashboardActivityFeedProps) {
  return (
    <DashboardSection actionLabel="View all >" title={title} onAction={onViewAll}>
      <View style={styles.panel}>
        {isLoading && items.length === 0 ? (
          <View style={styles.loadingState}>
            <ActivityIndicator color={colors.brand.orange} size="small" />
          </View>
        ) : items.length === 0 ? (
          <View style={styles.emptyState}>
            <ActivityEmptyIllustration />
            <View style={styles.emptyCopy}>
              <Text style={styles.emptyTitle}>No visitor activity yet</Text>
              <Text style={styles.emptyText}>
                Visitor check-ins and movements will appear here.
              </Text>
              {emptyAction ? (
                <Pressable
                  accessibilityLabel={emptyAction.accessibilityLabel ?? emptyAction.label}
                  accessibilityRole="button"
                  style={({ pressed }) => [
                    styles.emptyAction,
                    pressed && styles.emptyActionPressed,
                  ]}
                  onPress={emptyAction.onPress}
                >
                  <SymbolView
                    name={{ ios: "person.badge.plus", android: "person_add", web: "person_add" }}
                    size={16}
                    tintColor={colors.brand.orange}
                  />
                  <Text style={styles.emptyActionText}>{emptyAction.label}</Text>
                </Pressable>
              ) : null}
            </View>
          </View>
        ) : (
          <View style={styles.list}>
            {items.map((item) => (
              <ActivityPreviewCard
                key={`activity-${item.id}`}
                item={item}
                showFlat={showFlat}
                onPress={onItemPress}
              />
            ))}
            {isLoadingMore ? (
              <View style={styles.loadingMore}>
                <ActivityIndicator color={colors.brand.orange} size="small" />
              </View>
            ) : hasMore ? (
              <Pressable
                accessibilityRole="button"
                style={styles.loadMoreButton}
                onPress={onLoadMore}
              >
                <Text style={styles.loadMoreText}>Load more</Text>
              </Pressable>
            ) : null}
          </View>
        )}
      </View>
    </DashboardSection>
  );
}

const styles = StyleSheet.create({
  clipboardBody: {
    alignItems: "center",
    backgroundColor: "#FFF4EB",
    borderColor: "#FFE0C7",
    borderRadius: radius.lg,
    borderWidth: 1,
    height: 80,
    justifyContent: "center",
    width: 72,
  },
  clockBadge: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: "#FFE0C7",
    borderRadius: 999,
    borderWidth: 1,
    bottom: -2,
    height: 30,
    justifyContent: "center",
    position: "absolute",
    right: -4,
    width: 30,
    ...shadows.sm,
  },
  emptyAction: {
    alignItems: "center",
    alignSelf: "stretch",
    backgroundColor: "#FFF8F3",
    borderColor: "#F1E4DA",
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
    marginTop: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  emptyActionPressed: {
    opacity: 0.85,
  },
  emptyActionText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "700",
  },
  emptyCopy: {
    flex: 1,
    gap: spacing.xs,
    justifyContent: "center",
    minWidth: 0,
  },
  emptyState: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.lg,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.xl,
  },
  emptyText: {
    ...typography.bodySmall,
    color: colors.guard.textMuted,
    lineHeight: 20,
  },
  emptyTitle: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
    lineHeight: 22,
  },
  illustrationWrap: {
    height: 84,
    position: "relative",
    width: 76,
  },
  list: {
    gap: spacing.sm,
    paddingBottom: spacing.sm,
    paddingHorizontal: spacing.sm,
    paddingTop: spacing.sm,
  },
  loadingMore: {
    alignItems: "center",
    paddingVertical: spacing.md,
  },
  loadingState: {
    alignItems: "center",
    justifyContent: "center",
    minHeight: 160,
    paddingVertical: spacing["2xl"],
  },
  loadMoreButton: {
    alignItems: "center",
    borderTopColor: colors.border.default,
    borderTopWidth: StyleSheet.hairlineWidth,
    marginTop: spacing.sm,
    paddingVertical: spacing.md,
  },
  loadMoreText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "600",
  },
  panel: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    overflow: "hidden",
    ...shadows.card,
  },
  previewBadge: {
    borderRadius: radius["2xl"],
    borderWidth: 1,
    paddingHorizontal: spacing.sm,
    paddingVertical: 4,
  },
  previewBadgeText: {
    fontSize: 10,
    fontWeight: "700",
  },
  previewCard: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderColor: colors.border.default,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.lg,
  },
  previewCopy: {
    flex: 1,
    gap: 4,
    minWidth: 0,
  },
  previewMeta: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    lineHeight: 18,
  },
  previewName: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "700",
  },
  previewPressed: {
    opacity: 0.88,
  },
});
