import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { Row } from "@/components/layout";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type HubDestinationCardProps = {
  title: string;
  description: string;
  unreadCount?: number;
  icon: SymbolViewProps["name"];
  iconBackground: string;
  iconColor: string;
  onPress: () => void;
};

export function HubDestinationCard({
  title,
  description,
  unreadCount = 0,
  icon,
  iconBackground,
  iconColor,
  onPress,
}: HubDestinationCardProps) {
  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [styles.card, pressed && styles.pressed]}
      onPress={onPress}
    >
      <Row align="center" gap="md" justify="space-between">
        <Row align="center" gap="md" style={styles.leading}>
          <View style={[styles.iconWell, { backgroundColor: iconBackground }]}>
            <SymbolView name={icon} size={22} tintColor={iconColor} />
          </View>
          <View style={styles.copy}>
            <Row align="center" gap="sm">
              <Text style={styles.title}>{title}</Text>
              {unreadCount > 0 ? (
                <View style={styles.badge}>
                  <Text style={styles.badgeText}>
                    {unreadCount > 99 ? "99+" : unreadCount} new
                  </Text>
                </View>
              ) : null}
            </Row>
            <Text style={styles.description}>{description}</Text>
          </View>
        </Row>
        <SymbolView
          name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
          size={16}
          tintColor={colors.text.muted}
        />
      </Row>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  badge: {
    backgroundColor: colors.brand.navySoft,
    borderRadius: radius.full,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
  },
  badgeText: {
    color: colors.brand.navy,
    fontSize: 11,
    fontWeight: "700",
  },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    ...shadows.sm,
  },
  copy: {
    flex: 1,
    minWidth: 0,
  },
  description: {
    color: colors.text.muted,
    fontSize: 13,
    lineHeight: 18,
    marginTop: 4,
  },
  iconWell: {
    alignItems: "center",
    borderRadius: radius.full,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  leading: {
    flex: 1,
    minWidth: 0,
  },
  pressed: {
    opacity: 0.92,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
  },
});
