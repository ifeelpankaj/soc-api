import { type ReactNode } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type FilterChipProps = {
  label: string;
  selected?: boolean;
  onPress: () => void;
};

export function FilterChip({
  label,
  selected = false,
  onPress,
}: FilterChipProps) {
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityState={{ selected }}
      onPress={onPress}
      style={({ pressed }) => [
        styles.chip,
        selected ? styles.chipSelected : styles.chipDefault,
        pressed ? styles.chipPressed : null,
      ]}
    >
      <Text
        numberOfLines={1}
        style={[
          styles.chipText,
          selected ? styles.chipTextSelected : styles.chipTextDefault,
        ]}
      >
        {selected ? `✓ ${label}` : label}
      </Text>
    </Pressable>
  );
}

type FilterChipGridProps = {
  children: ReactNode;
};

export function FilterChipGrid({ children }: FilterChipGridProps) {
  return <View style={styles.grid}>{children}</View>;
}

type FilterChipCellProps = {
  children: ReactNode;
};

export function FilterChipCell({ children }: FilterChipCellProps) {
  return <View style={styles.cell}>{children}</View>;
}

const styles = StyleSheet.create({
  cell: {
    width: "48%",
  },
  chip: {
    alignItems: "center",
    borderRadius: 12,
    borderWidth: 1.5,
    justifyContent: "center",
    minHeight: 46,
    paddingHorizontal: spacing.md,
    paddingVertical: 10,
  },
  chipDefault: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
  },
  chipPressed: {
    opacity: 0.88,
  },
  chipSelected: {
    backgroundColor: colors.brand.navy,
    borderColor: colors.brand.navy,
  },
  chipText: {
    ...typography.bodySmall,
    textAlign: "center",
  },
  chipTextDefault: {
    color: colors.guard.textMuted,
    fontWeight: "500",
  },
  chipTextSelected: {
    color: colors.text.inverse,
    fontWeight: "700",
  },
  grid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 10,
    justifyContent: "space-between",
  },
});
