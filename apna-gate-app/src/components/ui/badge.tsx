import { StyleSheet, Text, View } from "react-native";

import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type BadgeTone = "slate" | "emerald" | "amber" | "rose" | "blue";

type BadgeProps = {
  label: string;
  tone?: BadgeTone;
};

const toneStyles: Record<BadgeTone, { backgroundColor: string; color: string }> = {
  slate: { backgroundColor: "#f1f5f9", color: "#334155" },
  emerald: { backgroundColor: "#ecfdf5", color: "#047857" },
  amber: { backgroundColor: "#fffbeb", color: "#b45309" },
  rose: { backgroundColor: "#fff1f2", color: "#be123c" },
  blue: { backgroundColor: "#eff6ff", color: "#1d4ed8" },
};

export function Badge({ label, tone = "slate" }: BadgeProps) {
  const toneStyle = toneStyles[tone];

  return (
    <View style={styles.wrapper}>
      <Text style={[styles.badge, { backgroundColor: toneStyle.backgroundColor, color: toneStyle.color }]}>
        {label}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    ...typography.caption,
    borderRadius: radius["2xl"],
    fontWeight: "600",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
  },
  wrapper: {
    alignSelf: "flex-start",
    borderRadius: radius["2xl"],
  },
});
