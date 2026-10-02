import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Row, Stack } from "@/components/layout";
import { Card, StatusPill } from "@/components/ui";
import { titleize } from "@/features/guard/guard-utils";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type SocietyAccessCardProps = {
  membershipCount: number;
  onSwitchPress?: () => void;
  role?: string | null;
  societyName: string;
  status?: string | null;
};

export function SocietyAccessCard({
  membershipCount,
  onSwitchPress,
  role,
  societyName,
  status,
}: SocietyAccessCardProps) {
  const roleLabel = role ? `${titleize(role)} Access` : "Staff Access";
  const statusLabel = status ? titleize(status) : "Active";

  return (
    <Card style={styles.card}>
      <Row align="flex-start" gap="md" justify="flex-start">
        <View style={styles.icon}>
          <SymbolView
            name={{ ios: "building.2.fill", android: "business", web: "business" }}
            size={20}
            tintColor={colors.brand.orange}
          />
        </View>
        <Stack gap="xs" style={styles.copy}>
          <Text style={styles.title}>{societyName}</Text>
          <Text style={styles.subtitle}>
            {roleLabel} • {statusLabel}
          </Text>
        </Stack>
        {status ? <StatusPill status={status} /> : null}
      </Row>

      {membershipCount > 1 && onSwitchPress ? (
        <Pressable
          accessibilityRole="button"
          style={({ pressed }) => [styles.switchButton, pressed && styles.switchButtonPressed]}
          onPress={onSwitchPress}
        >
          <Text style={styles.switchLabel}>Switch society</Text>
          <SymbolView
            name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
            size={14}
            tintColor={colors.brand.orange}
          />
        </Pressable>
      ) : null}
    </Card>
  );
}

const styles = StyleSheet.create({
  card: {
    gap: spacing.md,
  },
  copy: {
    flex: 1,
    minWidth: 0,
  },
  icon: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius["2xl"],
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  subtitle: {
    color: colors.text.secondary,
    fontSize: 14,
  },
  switchButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.xl,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.md,
    paddingVertical: 10,
  },
  switchButtonPressed: {
    opacity: 0.85,
  },
  switchLabel: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "600",
  },
  title: {
    color: colors.text.primary,
    fontSize: 18,
    fontWeight: "700",
  },
});
