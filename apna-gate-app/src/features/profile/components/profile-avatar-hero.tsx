import { StyleSheet, Text, View } from "react-native";

import { Stack } from "@/components/layout";
import { UserAvatar } from "@/components/ui";
import { formatPhoneDisplay } from "@/features/profile/profile-formatters";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type ProfileAvatarHeroProps = {
  avatarUrl?: string | null;
  email?: string | null;
  isActive?: boolean;
  name: string;
  onChangePhoto?: () => void;
  phone?: string | null;
};

export function ProfileAvatarHero({
  avatarUrl,
  email,
  isActive = true,
  name,
  onChangePhoto,
  phone,
}: ProfileAvatarHeroProps) {
  const phoneDisplay = formatPhoneDisplay(phone);

  return (
    <Stack align="center" gap="md" style={styles.hero}>
      <UserAvatar
        imageUrl={avatarUrl}
        name={name}
        showCameraBadge={Boolean(onChangePhoto)}
        size={96}
        onPress={onChangePhoto}
      />

      <Stack align="center" gap="xs">
        <Text style={styles.name}>{name}</Text>
        {email ? <Text style={styles.detail}>{email}</Text> : null}
        {phone && phoneDisplay !== "—" ? <Text style={styles.detail}>{phoneDisplay}</Text> : null}
      </Stack>

      <View
        style={[
          styles.statusBadge,
          { backgroundColor: isActive ? colors.status.successSoft : colors.surface.muted },
        ]}
      >
        <View
          style={[
            styles.statusDot,
            { backgroundColor: isActive ? colors.status.success : colors.text.muted },
          ]}
        />
        <Text
          style={[
            styles.statusText,
            { color: isActive ? colors.status.success : colors.text.muted },
          ]}
        >
          {isActive ? "Active Account" : "Inactive Account"}
        </Text>
      </View>
    </Stack>
  );
}

const styles = StyleSheet.create({
  detail: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  hero: {
    paddingVertical: spacing.sm,
  },
  name: {
    ...typography.title,
    color: colors.text.primary,
    letterSpacing: -0.5,
    textAlign: "center",
  },
  statusBadge: {
    alignItems: "center",
    borderRadius: 999,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: 6,
  },
  statusDot: {
    borderRadius: 999,
    height: 8,
    width: 8,
  },
  statusText: {
    fontSize: 12,
    fontWeight: "600",
  },
});
