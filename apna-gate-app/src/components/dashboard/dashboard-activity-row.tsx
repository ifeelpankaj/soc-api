import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Row } from "@/components/layout";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";
import {
  ACTIVITY_ROW_HEIGHT,
} from "@/features/shared/activity-feed-config";

type DashboardActivityRowProps = {
  meta: string;
  name: string;
  onPress?: () => void;
  statusBg: string;
  statusBorder: string;
  statusColor: string;
  statusLabel: string;
};

function getInitials(name: string) {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length >= 2) {
    return `${parts[0]?.[0] ?? ""}${parts[1]?.[0] ?? ""}`.toUpperCase();
  }
  return (parts[0]?.slice(0, 2) ?? "?").toUpperCase();
}

export function DashboardActivityRow({
  meta,
  name,
  onPress,
  statusBg,
  statusBorder,
  statusColor,
  statusLabel,
}: DashboardActivityRowProps) {
  const content = (
    <Row align="center" gap="md" justify="flex-start" style={styles.row}>
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>{getInitials(name)}</Text>
      </View>
      <View style={styles.copy}>
        <Text numberOfLines={1} style={styles.name}>
          {name}
        </Text>
        <Text numberOfLines={1} style={styles.meta}>
          {meta}
        </Text>
      </View>
      <View style={[styles.badge, { backgroundColor: statusBg, borderColor: statusBorder }]}>
        <Text style={[styles.badgeText, { color: statusColor }]}>{statusLabel}</Text>
      </View>
      <SymbolView
        name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
        size={14}
        tintColor={colors.guard.textMuted}
      />
    </Row>
  );

  if (!onPress) {
    return content;
  }

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [{ opacity: pressed ? 0.85 : 1 }]}
      onPress={onPress}
    >
      {content}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  avatar: {
    alignItems: "center",
    backgroundColor: colors.guard.tealSoft,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  avatarText: {
    color: colors.guard.teal,
    fontSize: 14,
    fontWeight: "700",
  },
  badge: {
    borderRadius: 999,
    borderWidth: 1,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
  },
  badgeText: {
    fontSize: 10,
    fontWeight: "600",
  },
  copy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  meta: {
    color: colors.guard.textMuted,
    fontSize: 12,
    lineHeight: 16,
  },
  name: {
    color: colors.guard.text,
    fontSize: 14,
    fontWeight: "600",
    lineHeight: 20,
  },
  row: {
    minHeight: ACTIVITY_ROW_HEIGHT,
    paddingHorizontal: spacing.xs,
    paddingVertical: spacing.sm,
  },
});
