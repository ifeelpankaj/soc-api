import type { SymbolViewProps } from "expo-symbols";
import { StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Stack } from "@/components/layout";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

export type ProfileSymbolName = SymbolViewProps["name"];

type ProfileInfoRowProps = {
  icon: ProfileSymbolName;
  label: string;
  value: string;
  valueMuted?: boolean;
  isLast?: boolean;
};

export function ProfileInfoRow({
  icon,
  label,
  value,
  valueMuted = false,
  isLast = false,
}: ProfileInfoRowProps) {
  return (
    <View style={[styles.row, !isLast && styles.rowBorder]}>
      <View style={styles.icon}>
        <SymbolView name={icon} size={18} tintColor={colors.brand.orange} />
      </View>
      <Stack gap="xs" style={styles.copy}>
        <Text style={styles.label}>{label}</Text>
        <Text numberOfLines={2} style={[styles.value, valueMuted && styles.valueMuted]}>
          {value}
        </Text>
      </Stack>
    </View>
  );
}

const styles = StyleSheet.create({
  copy: {
    flex: 1,
    minWidth: 0,
  },
  icon: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  label: {
    color: colors.text.muted,
    fontSize: 11,
    fontWeight: "600",
    letterSpacing: 1.32,
    textTransform: "uppercase",
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
  value: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "600",
  },
  valueMuted: {
    color: colors.text.muted,
  },
});
