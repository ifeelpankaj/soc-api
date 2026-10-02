import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { formatDateOnly, titleize } from "@/features/guard/guard-utils";
import type { ModelsFlatMemberInviteResponse } from "@/lib/api/resident-api-extensions";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const CARD_BORDER = "rgba(16, 29, 54, 0.08)";

type PendingInviteCardProps = {
  invite: ModelsFlatMemberInviteResponse;
  onPress?: () => void;
};

export function PendingInviteCard({ invite, onPress }: PendingInviteCardProps) {
  const roleLabel = titleize(invite.role ?? "family");
  const expiryLabel = invite.expires_at
    ? `Expires on ${formatDateOnly(invite.expires_at)}`
    : "Expires soon";

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [styles.card, pressed && styles.cardPressed]}
      onPress={onPress}
    >
      <View style={styles.iconWrap}>
        <SymbolView
          name={{ ios: "clock.fill", android: "schedule", web: "schedule" }}
          size={18}
          tintColor={colors.brand.orange}
        />
      </View>

      <View style={styles.body}>
        <Text numberOfLines={1} style={styles.name}>
          {invite.full_name ?? "Pending member"}
        </Text>
        <Text numberOfLines={1} style={styles.meta}>
          {roleLabel} • {expiryLabel}
        </Text>
      </View>

      <View style={styles.pendingBadge}>
        <SymbolView
          name={{ ios: "clock", android: "schedule", web: "schedule" }}
          size={10}
          tintColor={colors.brand.orange}
        />
        <Text style={styles.pendingText}>Pending</Text>
      </View>

      <SymbolView
        name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
        size={14}
        tintColor={colors.guard.textMuted}
      />
    </Pressable>
  );
}

const styles = StyleSheet.create({
  body: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  card: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    ...shadows.sm,
  },
  cardPressed: {
    opacity: 0.92,
  },
  iconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  meta: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "500",
  },
  name: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
  },
  pendingBadge: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderColor: "#FED7AA",
    borderRadius: radius["2xl"],
    borderWidth: 1,
    flexDirection: "row",
    flexShrink: 0,
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  pendingText: {
    color: colors.brand.orange,
    fontSize: 10,
    fontWeight: "700",
  },
});
