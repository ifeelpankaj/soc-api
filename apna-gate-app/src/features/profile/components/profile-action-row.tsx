import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Stack } from "@/components/layout";
import type { ProfileSymbolName } from "@/features/profile/components/profile-info-row";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type ProfileActionRowProps = {
  description?: string;
  destructive?: boolean;
  icon: ProfileSymbolName;
  isLast?: boolean;
  label: string;
  onPress: () => void;
  showChevron?: boolean;
};

export function ProfileActionRow({
  description,
  destructive = false,
  icon,
  isLast = false,
  label,
  onPress,
  showChevron = true,
}: ProfileActionRowProps) {
  const accent = destructive ? colors.status.error : colors.brand.orange;
  const iconBg = destructive ? colors.status.errorSoft : colors.brand.orangeSoft;

  return (
    <Pressable
      accessibilityRole="button"
      style={[styles.row, !isLast && styles.rowBorder]}
      onPress={onPress}
    >
      <View style={[styles.icon, { backgroundColor: iconBg }]}>
        <SymbolView name={icon} size={18} tintColor={accent} />
      </View>
      <Stack gap="xs" style={styles.copy}>
        <Text style={[styles.label, destructive && styles.labelDestructive]}>{label}</Text>
        {description ? <Text style={styles.description}>{description}</Text> : null}
      </Stack>
      {showChevron ? (
        <SymbolView
          name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
          size={16}
          tintColor={colors.text.muted}
        />
      ) : null}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  copy: {
    flex: 1,
    minWidth: 0,
  },
  description: {
    color: colors.text.secondary,
    fontSize: 14,
  },
  icon: {
    alignItems: "center",
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  label: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "600",
  },
  labelDestructive: {
    color: colors.status.error,
  },
  row: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: 14,
  },
  rowBorder: {
    borderBottomColor: "#f1f5f9",
    borderBottomWidth: 1,
  },
});
