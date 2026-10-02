import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";

import { VisitorTimelineRow } from "@/features/visitors/components/visitor-timeline-row";
import {
  getFlatLabel,
  getVisitorName,
  getVisitorStatusMeta,
  titleize,
} from "@/features/visitors/visitor-utils";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type LogEntryRowProps = {
  entry: ModelsVisitorEntry;
  isCheckingOut?: boolean;
  onCheckOut?: () => void;
  onPress?: () => void;
};

const CARD_BORDER = "rgba(16, 29, 54, 0.11)";

export function LogEntryRow({ entry, isCheckingOut, onCheckOut, onPress }: LogEntryRowProps) {
  const statusMeta = getVisitorStatusMeta(entry.status);
  const purpose = entry.purpose ? titleize(entry.purpose) : "Visitor";
  const visitorName = getVisitorName(entry);
  const hasCheckoutAction = entry.status === "checked_in" && Boolean(onCheckOut);

  return (
    <View style={styles.card}>
      <View style={styles.topSection}>
        <View style={styles.titleRow}>
          <Pressable
            accessibilityRole="button"
            disabled={!onPress}
            style={({ pressed }) => [styles.titlePressable, onPress && pressed && styles.cardPressed]}
            onPress={onPress}
          >
            <VisitorPhoto entry={entry} fallback={<SymbolView
                name={{ ios: "person.fill", android: "person", web: "person" }}
                size={18}
                tintColor={colors.brand.orange}
              />} />

            <View style={styles.titleCopy}>
              <Text numberOfLines={1} style={styles.name}>
                {visitorName}
              </Text>
              <Text numberOfLines={1} style={styles.meta}>
                {purpose} · {getFlatLabel(entry)}
              </Text>
            </View>
          </Pressable>

          <View style={[styles.statusBadge, { backgroundColor: statusMeta.bg, borderColor: statusMeta.border }]}>
            <Text style={[styles.statusText, { color: statusMeta.color }]}>{statusMeta.label}</Text>
          </View>

          <Pressable
            accessibilityLabel="Open visitor details"
            accessibilityRole="button"
            hitSlop={8}
            style={styles.menuButton}
            onPress={onPress}
          >
            <SymbolView
              name={{ ios: "ellipsis", android: "more_vert", web: "more_vert" }}
              size={18}
              tintColor={colors.guard.textMuted}
            />
          </Pressable>
        </View>

        <Pressable
          accessibilityRole="button"
          disabled={!onPress}
          onPress={onPress}
        >
          <VisitorTimelineRow entry={entry} />
        </Pressable>
      </View>

      {hasCheckoutAction ? (
        <Pressable
          accessibilityRole="button"
          disabled={isCheckingOut}
          style={({ pressed }) => [
            styles.checkoutButton,
            pressed && !isCheckingOut && styles.checkoutButtonPressed,
          ]}
          onPress={onCheckOut}
        >
          <SymbolView
            name={{ ios: "arrow.right.square", android: "logout", web: "logout" }}
            size={16}
            tintColor={colors.brand.orange}
          />
          <Text style={styles.checkoutText}>{isCheckingOut ? "Checking out..." : "Check out"}</Text>
        </Pressable>
      ) : null}
    </View>
  );
}

export function LogEntryDivider() {
  return <View style={styles.divider} />;
}

const styles = StyleSheet.create({
  avatar: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.sm,
    height: 36,
    justifyContent: "center",
    width: 36,
  },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.sm,
    borderWidth: 1,
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.md,
  },
  cardPressed: {
    opacity: 0.94,
  },
  checkoutButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderColor: "rgba(255, 106, 26, 0.18)",
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    paddingVertical: spacing.sm,
  },
  checkoutButtonPressed: {
    opacity: 0.88,
  },
  checkoutText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
  },
  divider: {
    height: spacing.xs,
  },
  menuButton: {
    alignItems: "center",
    height: 28,
    justifyContent: "center",
    width: 20,
  },
  meta: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "500",
  },
  name: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
  },
  statusBadge: {
    borderRadius: radius.sm,
    borderWidth: 1,
    paddingHorizontal: 8,
    paddingVertical: 3,
  },
  statusText: {
    fontSize: 10,
    fontWeight: "700",
  },
  titleCopy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  titlePressable: {
    alignItems: "center",
    flex: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minWidth: 0,
  },
  titleRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
  topSection: {
    gap: spacing.sm,
  },
});
