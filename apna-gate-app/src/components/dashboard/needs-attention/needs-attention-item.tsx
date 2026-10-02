import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Row } from "@/components/layout";
import { dashboardActionToneStyles } from "@/components/dashboard/dashboard-action-tile";
import type { DashboardActionTone } from "@/components/dashboard/dashboard-action-tile";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

import type { AttentionItem } from "./types";

type NeedsAttentionItemProps = {
  item: AttentionItem;
  isLast?: boolean;
};

const toneMap: Record<
  NonNullable<AttentionItem["iconTone"]>,
  DashboardActionTone | "teal"
> = {
  orange: "orange",
  blue: "blue",
  purple: "purple",
  neutral: "neutral",
  teal: "teal",
};

const tealWell = {
  backgroundColor: colors.guard.tealSoft,
  iconColor: colors.guard.teal,
};

export function NeedsAttentionItem({ item, isLast }: NeedsAttentionItemProps) {
  const tone = item.iconTone ?? "orange";
  const wellStyle =
    tone === "teal"
      ? tealWell
      : dashboardActionToneStyles[toneMap[tone] as DashboardActionTone];

  return (
    <Pressable
      accessibilityRole="button"
      accessibilityLabel={`${item.title}. ${item.subtitle}`}
      style={({ pressed }) => [
        styles.row,
        !isLast && styles.rowBorder,
        pressed && styles.rowPressed,
      ]}
      onPress={item.onPress}
    >
      <Row align="center" gap="md" justify="flex-start" style={styles.content}>
        <View style={[styles.iconWell, { backgroundColor: wellStyle.backgroundColor }]}>
          <SymbolView name={item.icon} size={18} tintColor={wellStyle.iconColor} />
        </View>
        <View style={styles.copy}>
          <Row align="center" gap="sm" justify="flex-start">
            <Text numberOfLines={2} style={styles.title}>
              {item.title}
            </Text>
            {item.badge ? (
              <View style={styles.badge}>
                <Text style={styles.badgeText}>{item.badge}</Text>
              </View>
            ) : null}
          </Row>
          <Text numberOfLines={2} style={styles.subtitle}>
            {item.subtitle}
          </Text>
        </View>
        <Row align="center" gap="xs" justify="flex-end">
          {item.actionLabel ? (
            <Text style={styles.action}>{item.actionLabel}</Text>
          ) : null}
          <SymbolView
            name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
            size={14}
            tintColor={colors.text.muted}
          />
        </Row>
      </Row>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  action: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
  },
  badge: {
    backgroundColor: colors.surface.secondary,
    borderRadius: 6,
    paddingHorizontal: spacing.xs,
    paddingVertical: 2,
  },
  badgeText: {
    color: colors.text.secondary,
    fontSize: 11,
    fontWeight: "600",
  },
  content: {
    flex: 1,
  },
  copy: {
    flex: 1,
    minWidth: 0,
    paddingRight: spacing.sm,
  },
  iconWell: {
    alignItems: "center",
    borderRadius: 999,
    height: 36,
    justifyContent: "center",
    width: 36,
  },
  row: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  rowBorder: {
    borderBottomColor: colors.dashboard.cardBorder,
    borderBottomWidth: StyleSheet.hairlineWidth,
  },
  rowPressed: {
    backgroundColor: colors.surface.muted,
  },
  subtitle: {
    color: colors.text.muted,
    fontSize: 13,
    lineHeight: 18,
    marginTop: 2,
  },
  title: {
    color: colors.brand.navy,
    flexShrink: 1,
    fontSize: 15,
    fontWeight: "600",
    lineHeight: 20,
  },
});
