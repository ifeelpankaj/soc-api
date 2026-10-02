import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export type DashboardOverviewTone = "green" | "orange" | "blue" | "purple" | "neutral";

const toneStyles: Record<
  DashboardOverviewTone,
  { backgroundColor: string; iconColor: string; valueColor: string }
> = {
  green: {
    backgroundColor: colors.dashboard.statGreenSoft,
    iconColor: colors.dashboard.statGreen,
    valueColor: colors.dashboard.statGreen,
  },
  orange: {
    backgroundColor: colors.dashboard.statOrangeSoft,
    iconColor: colors.dashboard.statOrange,
    valueColor: colors.dashboard.statOrange,
  },
  blue: {
    backgroundColor: colors.dashboard.statBlueSoft,
    iconColor: colors.dashboard.statBlue,
    valueColor: colors.dashboard.statBlue,
  },
  purple: {
    backgroundColor: colors.dashboard.statPurpleSoft,
    iconColor: colors.dashboard.statPurple,
    valueColor: colors.dashboard.statPurple,
  },
  neutral: {
    backgroundColor: colors.brand.navySoft,
    iconColor: colors.brand.navy,
    valueColor: colors.brand.navy,
  },
};

export type DashboardOverviewStatConfig = {
  icon: SymbolViewProps["name"];
  id: string;
  label: string;
  subtext?: string;
  tone: DashboardOverviewTone;
  value: number | string;
};

type DashboardOverviewStatProps = Omit<DashboardOverviewStatConfig, "id"> & {
  onPress?: () => void;
};

export function DashboardOverviewStat({
  icon,
  label,
  onPress,
  subtext,
  tone,
  value,
}: DashboardOverviewStatProps) {
  const toneStyle = toneStyles[tone] ?? toneStyles.neutral;

  const content = (
    <View style={styles.card}>
      <View style={styles.row}>
        <View style={[styles.iconWrap, { backgroundColor: toneStyle.backgroundColor }]}>
          <SymbolView name={icon} size={26} tintColor={toneStyle.iconColor} />
        </View>
        <View style={styles.copy}>
          <Text style={[styles.value, { color: toneStyle.valueColor }]}>{value}</Text>
          <Text numberOfLines={2} style={styles.label}>
            {label}
          </Text>
          {subtext ? (
            <Text numberOfLines={1} style={styles.subtext}>
              {subtext}
            </Text>
          ) : null}
        </View>
      </View>
    </View>
  );

  if (!onPress) {
    return content;
  }

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [styles.pressable, pressed && styles.pressed]}
      onPress={onPress}
    >
      {content}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    flex: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    width: "100%",
    ...shadows.card,
  },
  copy: {
    flex: 1,
    gap: 1,
    justifyContent: "center",
    minWidth: 0,
  },
  iconWrap: {
    alignItems: "center",
    borderRadius: radius.lg,
    height: 50,
    justifyContent: "center",
    width: 50,
  },
  label: {
    color: colors.brand.navy,
    fontSize: 12,
    fontWeight: "700",
    lineHeight: 16,
  },
  pressed: {
    opacity: 0.88,
  },
  pressable: {
    flex: 1,
    width: "100%",
  },
  row: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.md,
  },
  subtext: {
    color: colors.guard.textMuted,
    fontSize: 10,
    fontWeight: "500",
    lineHeight: 14,
  },
  value: {
    fontSize: 22,
    fontWeight: "800",
    letterSpacing: -0.4,
    lineHeight: 26,
  },
});
