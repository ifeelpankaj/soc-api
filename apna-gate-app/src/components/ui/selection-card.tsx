import { Pressable, StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type SelectionCardProps = {
  title: string;
  subtitle: string;
  actionLabel: string;
  badgeLabel?: string;
  badgeTone?: "active" | "pending" | "default";
  selected?: boolean;
  disabled?: boolean;
  onPress: () => void;
};

const badgeStyles = {
  active: {
    backgroundColor: colors.status.successSoft,
    color: colors.status.success,
  },
  pending: {
    backgroundColor: colors.status.warningSoft,
    color: colors.status.warning,
  },
  default: {
    backgroundColor: colors.guard.sectionBg,
    color: colors.text.muted,
  },
} as const;

export function SelectionCard({
  title,
  subtitle,
  actionLabel,
  badgeLabel,
  badgeTone = "default",
  selected = false,
  disabled = false,
  onPress,
}: SelectionCardProps) {
  const badge = badgeStyles[badgeTone];

  return (
    <View style={[styles.card, selected && styles.cardSelected]}>
      <View style={styles.cardHeader}>
        <View style={styles.cardCopy}>
          <Text style={styles.cardTitle}>{title}</Text>
          <Text style={styles.cardSubtitle}>{subtitle}</Text>
        </View>
        {badgeLabel ? (
          <View style={[styles.statusBadge, { backgroundColor: badge.backgroundColor }]}>
            <Text style={[styles.statusText, { color: badge.color }]}>{badgeLabel}</Text>
          </View>
        ) : null}
      </View>
      <Pressable
        accessibilityRole="button"
        disabled={disabled}
        onPress={onPress}
        style={[styles.cardButton, disabled && styles.cardButtonDisabled]}
      >
        <Text style={styles.cardButtonText}>{actionLabel}</Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.selection.cardBg,
    borderColor: colors.selection.cardBorder,
    borderRadius: radius.xl,
    borderWidth: 1,
    gap: 14,
    padding: spacing.lg,
  },
  cardButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 14,
    justifyContent: "center",
    minHeight: 46,
  },
  cardButtonDisabled: {
    opacity: 0.55,
  },
  cardButtonText: {
    ...typography.bodySmall,
    color: colors.text.inverse,
    fontSize: 15,
    fontWeight: "700",
  },
  cardCopy: {
    flex: 1,
    gap: spacing.xs,
  },
  cardHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.md,
  },
  cardSelected: {
    borderColor: colors.guard.teal,
    borderWidth: 2,
  },
  cardSubtitle: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  cardTitle: {
    color: colors.text.primary,
    fontSize: 18,
    fontWeight: "700",
  },
  statusBadge: {
    borderRadius: 999,
    paddingHorizontal: 10,
    paddingVertical: 5,
  },
  statusText: {
    fontSize: 11,
    fontWeight: "700",
  },
});
