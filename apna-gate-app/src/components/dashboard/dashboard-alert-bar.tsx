import { Pressable, StyleSheet, Text } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { Row } from "@/components/layout";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type DashboardAlertBarProps = {
  actionLabel?: string;
  count?: number;
  icon?: SymbolViewProps["name"];
  message?: string;
  onPress: () => void;
};

export function DashboardAlertBar({
  actionLabel = "Review →",
  count,
  icon = { ios: "person.fill", android: "person", web: "person" },
  message,
  onPress,
}: DashboardAlertBarProps) {
  const resolvedMessage =
    message ??
    (count !== undefined
      ? `${count} visitor${count === 1 ? "" : "s"} waiting for approval`
      : "Action required");

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [styles.container, { opacity: pressed ? 0.88 : 1 }]}
      onPress={onPress}
    >
      <Row align="center" gap="md" justify="space-between">
        <Row align="center" gap="sm" justify="flex-start" style={styles.content}>
          <SymbolView name={icon} size={16} tintColor={colors.brand.orange} />
          <Text style={styles.message}>{resolvedMessage}</Text>
        </Row>
        <Text style={styles.action}>{actionLabel}</Text>
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
  container: {
    backgroundColor: colors.surface.card,
    borderRadius: radius.lg,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    ...shadows.sm,
  },
  content: {
    flex: 1,
    paddingRight: spacing.md,
  },
  message: {
    color: colors.text.secondary,
    flex: 1,
    fontSize: 14,
    fontWeight: "500",
  },
});
