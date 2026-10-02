import type { ReactNode } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";

import { Stack } from "@/components/layout";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type DashboardSectionProps = {
  actionLabel?: string;
  children: ReactNode;
  onAction?: () => void;
  subtitle?: string;
  title: string;
  trailing?: string;
};

export function DashboardSection({
  actionLabel,
  children,
  onAction,
  subtitle,
  title,
  trailing,
}: DashboardSectionProps) {
  return (
    <Stack gap="md">
      <View style={styles.header}>
        <View style={styles.copy}>
          <Text style={styles.title}>{title}</Text>
          {subtitle ? <Text style={styles.subtitle}>{subtitle}</Text> : null}
        </View>
        {actionLabel && onAction ? (
          <Pressable accessibilityRole="button" hitSlop={8} onPress={onAction}>
            <Text style={styles.action}>{actionLabel}</Text>
          </Pressable>
        ) : null}
        {trailing ? <Text style={styles.action}>{trailing}</Text> : null}
      </View>
      {children}
    </Stack>
  );
}

const styles = StyleSheet.create({
  action: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "700",
  },
  copy: {
    flex: 1,
    paddingRight: spacing.sm,
  },
  header: {
    alignItems: "center",
    flexDirection: "row",
    justifyContent: "space-between",
    width: "100%",
  },
  subtitle: {
    color: colors.text.secondary,
    fontSize: 14,
    marginTop: 2,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 18,
    fontWeight: "700",
    letterSpacing: -0.2,
  },
});
