import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export type DashboardActionTone = "orange" | "blue" | "purple" | "neutral";

export const dashboardActionToneStyles: Record<
  DashboardActionTone,
  { backgroundColor: string; iconColor: string }
> = {
  orange: {
    backgroundColor: colors.dashboard.actionOrangeSoft,
    iconColor: colors.dashboard.actionOrange,
  },
  blue: {
    backgroundColor: colors.dashboard.actionBlueSoft,
    iconColor: colors.dashboard.actionBlue,
  },
  purple: {
    backgroundColor: colors.dashboard.actionPurpleSoft,
    iconColor: colors.dashboard.actionPurple,
  },
  neutral: {
    backgroundColor: colors.dashboard.actionNeutralSoft,
    iconColor: colors.dashboard.actionNeutral,
  },
};

export type DashboardActionTileConfig = {
  badgeCount?: number;
  icon: SymbolViewProps["name"];
  id: string;
  onPress: () => void;
  subtitle: string;
  title: string;
  tone: DashboardActionTone;
};

type DashboardActionTileProps = Omit<DashboardActionTileConfig, "id">;

export function DashboardActionTile({
  badgeCount = 0,
  icon,
  onPress,
  subtitle,
  title,
  tone,
}: DashboardActionTileProps) {
  const toneStyle = dashboardActionToneStyles[tone];

  return (
    <Pressable
      accessibilityLabel={badgeCount > 0 ? `${title}, ${subtitle}` : title}
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.tile,
        {
          opacity: pressed ? 0.88 : 1,
        },
      ]}
      onPress={onPress}
    >
      <View
        style={[
          styles.iconWrap,
          { backgroundColor: toneStyle.backgroundColor },
        ]}
      >
        <SymbolView name={icon} size={22} tintColor={toneStyle.iconColor} />
        {badgeCount > 0 ? (
          <View style={styles.badge}>
            <Text style={styles.badgeText}>{badgeCount > 99 ? "99+" : badgeCount}</Text>
          </View>
        ) : null}
      </View>
      <Text style={styles.title}>
        {title}
      </Text>
      <Text numberOfLines={2} style={styles.subtitle}>
        {subtitle}
      </Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  badge: {
    position: "absolute",
    right: -10,
    top: -5,
    minWidth: 22,
    height: 22,
    borderRadius: 11,
    backgroundColor: colors.brand.orange,
    alignItems: "center",
    justifyContent: "center",
    paddingHorizontal: 4,
  },
  badgeText: { color: "#FFFFFF", fontSize: 11, fontWeight: "700" },
  iconWrap: {
    alignItems: "center",
    borderRadius: 999,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  subtitle: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "500",
    lineHeight: 14,
    textAlign: "center",
  },
  tile: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    flex: 1,
    gap: spacing.sm,
    justifyContent: "center",
    minHeight: layout.actionTileMinHeight,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.md,
    ...shadows.card,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
    textAlign: "center",
  },
});
