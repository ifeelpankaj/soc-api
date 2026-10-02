import { Pressable, StyleSheet, Text } from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type QuickActionProps = {
  label: string;
  onPress: () => void;
  tone?: "light" | "danger";
};

export function QuickAction({ label, onPress, tone = "light" }: QuickActionProps) {
  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.base,
        tone === "danger" ? styles.danger : styles.light,
        pressed && styles.pressed,
      ]}
      onPress={onPress}
    >
      <Text style={[styles.label, tone === "danger" ? styles.dangerLabel : styles.lightLabel]}>
        {label}
      </Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  base: {
    alignItems: "center",
    borderRadius: radius.lg,
    borderWidth: 1,
    flex: 1,
    justifyContent: "center",
    minHeight: layout.buttonHeightCompact + 20,
    paddingHorizontal: spacing.md,
  },
  danger: {
    backgroundColor: colors.status.errorSoft,
    borderColor: "#fecdd3",
  },
  dangerLabel: {
    color: "#be123c",
  },
  label: {
    ...typography.bodySmall,
    fontWeight: "700",
    textAlign: "center",
  },
  light: {
    backgroundColor: colors.surface.input,
    borderColor: "#fef3c7",
    ...shadows.sm,
  },
  lightLabel: {
    color: colors.text.primary,
  },
  pressed: {
    opacity: 0.8,
  },
});
