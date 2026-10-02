import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import type { ModelsFlatResidentResponse } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const CARD_BORDER = "rgba(16, 29, 54, 0.08)";

type MemberCardProps = {
  member: ModelsFlatResidentResponse;
  onPress?: () => void;
};

function getMemberInitials(name: string) {
  const parts = name.split("-").filter(Boolean);

  if (parts.length >= 2 && /^\d+$/.test(parts[parts.length - 1] ?? "")) {
    return `${parts[0]?.charAt(0) ?? ""}${parts[parts.length - 1]}`.toUpperCase();
  }

  const words = name.trim().split(/\s+/);

  if (words.length >= 2) {
    return `${words[0]?.charAt(0) ?? ""}${words[1]?.charAt(0) ?? ""}`.toUpperCase();
  }

  return name.slice(0, 2).toUpperCase();
}

function getAvatarStyle(role?: string | null, isPrimary?: boolean) {
  if (isPrimary || role === "owner") {
    return { backgroundColor: colors.status.successSoft, color: colors.status.success };
  }

  if (role === "tenant") {
    return { backgroundColor: "#EFF6FF", color: "#2563EB" };
  }

  return { backgroundColor: colors.brand.orangeSoft, color: colors.brand.orange };
}

function roleChipLabel(role?: string | null) {
  switch (role) {
    case "owner":
      return "Owner";
    case "tenant":
      return "Tenant";
    default:
      return "Family member";
  }
}

function StatusBadge({ member }: { member: ModelsFlatResidentResponse }) {
  if (member.is_primary) {
    return (
      <View style={[styles.statusBadge, styles.primaryBadge]}>
        <SymbolView
          name={{ ios: "checkmark.shield.fill", android: "verified_user", web: "verified_user" }}
          size={11}
          tintColor={colors.status.success}
        />
        <Text style={[styles.statusText, styles.primaryText]}>Primary</Text>
      </View>
    );
  }

  if (member.role === "family") {
    return (
      <View style={[styles.statusBadge, styles.familyBadge]}>
        <SymbolView
          name={{ ios: "person.2.fill", android: "group", web: "group" }}
          size={11}
          tintColor="#2563EB"
        />
        <Text style={[styles.statusText, styles.familyText]}>Family</Text>
      </View>
    );
  }

  if (member.role === "owner") {
    return (
      <View style={[styles.statusBadge, styles.primaryBadge]}>
        <SymbolView
          name={{ ios: "house.fill", android: "home", web: "home" }}
          size={11}
          tintColor={colors.status.success}
        />
        <Text style={[styles.statusText, styles.primaryText]}>Owner</Text>
      </View>
    );
  }

  if (member.role === "tenant") {
    return (
      <View style={[styles.statusBadge, styles.tenantBadge]}>
        <SymbolView
          name={{ ios: "key.fill", android: "vpn_key", web: "vpn_key" }}
          size={11}
          tintColor="#7C3AED"
        />
        <Text style={[styles.statusText, styles.tenantText]}>Tenant</Text>
      </View>
    );
  }

  return null;
}

function ContactRow({
  icon,
  value,
}: {
  icon: SymbolViewProps["name"];
  value: string;
}) {
  return (
    <View style={styles.contactRow}>
      <SymbolView name={icon} size={12} tintColor={colors.brand.orange} />
      <Text numberOfLines={1} style={styles.contactText}>
        {value}
      </Text>
    </View>
  );
}

export function MemberCard({ member, onPress }: MemberCardProps) {
  const name = member.user_name ?? "Member";
  const avatarStyle = getAvatarStyle(member.role, member.is_primary);

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [styles.card, pressed && styles.cardPressed]}
      onPress={onPress}
    >
      <View style={[styles.avatar, { backgroundColor: avatarStyle.backgroundColor }]}>
        <Text style={[styles.avatarText, { color: avatarStyle.color }]}>
          {getMemberInitials(name)}
        </Text>
      </View>

      <View style={styles.body}>
        <View style={styles.topRow}>
          <Text numberOfLines={1} style={styles.name}>
            {name}
          </Text>
          <StatusBadge member={member} />
        </View>

        {member.user_phone ? (
          <ContactRow
            icon={{ ios: "phone.fill", android: "phone", web: "phone" }}
            value={member.user_phone}
          />
        ) : null}

        {member.user_email ? (
          <ContactRow
            icon={{ ios: "envelope.fill", android: "email", web: "email" }}
            value={member.user_email}
          />
        ) : null}

        <View style={styles.roleChip}>
          <Text style={styles.roleChipText}>{roleChipLabel(member.role)}</Text>
        </View>
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
  avatar: {
    alignItems: "center",
    borderRadius: 999,
    height: 48,
    justifyContent: "center",
    width: 48,
  },
  avatarText: {
    fontSize: 15,
    fontWeight: "800",
  },
  body: {
    flex: 1,
    gap: 4,
    minWidth: 0,
  },
  card: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    ...shadows.sm,
  },
  cardPressed: {
    opacity: 0.92,
  },
  contactRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: 6,
  },
  contactText: {
    color: colors.brand.navy,
    flex: 1,
    fontSize: 13,
    fontWeight: "500",
  },
  familyBadge: {
    backgroundColor: "#EFF6FF",
    borderColor: "#BFDBFE",
  },
  familyText: {
    color: "#2563EB",
  },
  name: {
    color: colors.brand.navy,
    flex: 1,
    fontSize: 16,
    fontWeight: "700",
    minWidth: 0,
  },
  primaryBadge: {
    backgroundColor: colors.status.successSoft,
    borderColor: "#bbf7d0",
  },
  primaryText: {
    color: colors.status.success,
  },
  roleChip: {
    alignSelf: "flex-start",
    backgroundColor: colors.surface.secondary,
    borderRadius: 999,
    marginTop: 2,
    paddingHorizontal: 8,
    paddingVertical: 3,
  },
  roleChipText: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "500",
  },
  statusBadge: {
    alignItems: "center",
    borderRadius: radius["2xl"],
    borderWidth: 1,
    flexDirection: "row",
    flexShrink: 0,
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 3,
  },
  statusText: {
    fontSize: 10,
    fontWeight: "700",
  },
  tenantBadge: {
    backgroundColor: "#F5F3FF",
    borderColor: "#DDD6FE",
  },
  tenantText: {
    color: "#7C3AED",
  },
  topRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
});
