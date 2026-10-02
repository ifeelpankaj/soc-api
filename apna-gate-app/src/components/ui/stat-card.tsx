import { StyleSheet, Text, View, type ViewStyle } from "react-native";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type StatCardProps = {
  label: string;
  style?: ViewStyle;
  tone?: "default" | "teal" | "warning" | "success";
  value: string | number;
};

const toneColors: Record<NonNullable<StatCardProps["tone"]>, string> = {
  default: colors.text.primary,
  success: colors.status.success,
  teal: colors.guard.teal,
  warning: colors.status.warning,
};

export function StatCard({ label, style, tone = "default", value }: StatCardProps) {
  return (
    <View style={[styles.card, style]}>
      <Text style={styles.label}>{label}</Text>
      <Text style={[styles.value, { color: toneColors[tone] }]}>{value}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface.card,
    borderRadius: radius.xl,
    flex: 1,
    gap: spacing.sm,
    minWidth: "46%",
    padding: spacing.lg,
    ...shadows.sm,
  },
  label: {
    ...typography.caption,
    color: colors.text.muted,
    fontSize: 11,
    fontWeight: "500",
    letterSpacing: 0.88,
    textTransform: "uppercase",
  },
  value: {
    fontSize: 30,
    fontWeight: "700",
    letterSpacing: -0.5,
    lineHeight: 30,
  },
});
