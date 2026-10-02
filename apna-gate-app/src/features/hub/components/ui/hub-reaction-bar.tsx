import { StyleSheet, View } from "react-native";

import { HubChip } from "@/features/hub/components/ui/hub-chip";
import type { HubReactionType } from "@/lib/api/hub-api";
import { spacing } from "@/theme/spacing";

const REACTIONS: { type: HubReactionType; label: string }[] = [
  { type: "like", label: "Like" },
  { type: "love", label: "Love" },
  { type: "helpful", label: "Helpful" },
];

type HubReactionBarProps = {
  counts?: Record<string, number>;
  mine?: string[];
  disabled?: boolean;
  onToggle: (type: HubReactionType) => void;
};

export function HubReactionBar({ counts, mine = [], disabled, onToggle }: HubReactionBarProps) {
  return (
    <View style={styles.row}>
      {REACTIONS.map(({ type, label }) => {
        const active = mine.includes(type);
        const count = counts?.[type] ?? 0;
        const text = count > 0 ? `${label} (${count})` : label;
        return (
          <View key={type} pointerEvents={disabled ? "none" : "auto"} style={disabled ? styles.disabled : undefined}>
            <HubChip label={text} selected={active} onPress={() => onToggle(type)} />
          </View>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  disabled: {
    opacity: 0.5,
  },
  row: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
    marginTop: spacing.md,
  },
});
