import { Pressable, ScrollView, StyleSheet, Text, View } from "react-native";

import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubChipProps = {
  label: string;
  selected?: boolean;
  onPress?: () => void;
};

export function HubChip({ label, selected, onPress }: HubChipProps) {
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityState={{ selected: Boolean(selected) }}
      style={({ pressed }) => [
        styles.chip,
        selected && styles.chipSelected,
        pressed && styles.pressed,
      ]}
      onPress={onPress}
    >
      <Text style={[styles.label, selected && styles.labelSelected]}>{label}</Text>
    </Pressable>
  );
}

type HubChipRowProps = {
  children: React.ReactNode;
};

export function HubChipRow({ children }: HubChipRowProps) {
  return (
    <ScrollView
      horizontal
      contentContainerStyle={styles.row}
      showsHorizontalScrollIndicator={false}
    >
      {children}
    </ScrollView>
  );
}

export function HubChipWrap({ children }: HubChipRowProps) {
  return <View style={styles.wrap}>{children}</View>;
}

const styles = StyleSheet.create({
  chip: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.full,
    borderWidth: StyleSheet.hairlineWidth,
    justifyContent: "center",
    minHeight: hubTheme.minTouchTarget,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  chipSelected: {
    backgroundColor: colors.brand.navySoft,
    borderColor: colors.brand.navy,
  },
  label: {
    color: colors.text.secondary,
    fontSize: 13,
    fontWeight: "600",
  },
  labelSelected: {
    color: colors.brand.navy,
  },
  pressed: {
    opacity: hubTheme.pressOpacity,
  },
  row: {
    flexDirection: "row",
    gap: spacing.sm,
    paddingVertical: spacing.xs,
  },
  wrap: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
});
