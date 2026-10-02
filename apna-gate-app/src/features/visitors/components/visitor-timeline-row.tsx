import { StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import {
  getVisitorTimelineBlocks,
  type VisitorTimelineBlock,
} from "@/features/guard/guard-utils";
import type { ModelsVisitorEntry, ModelsVisitorPendingEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type VisitorTimelineRowProps = {
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry;
  blocks?: VisitorTimelineBlock[];
};

function TimelineBlock({ block }: { block: VisitorTimelineBlock }) {
  if (!block.date && !block.time) {
    return null;
  }

  return (
    <View style={styles.block}>
      <Text style={styles.blockLabel}>{block.label}</Text>
      {block.date ? <Text style={styles.blockDate}>{block.date}</Text> : null}
      {block.time ? (
        <View style={styles.timeRow}>
          <SymbolView
            name={{ ios: "clock", android: "schedule", web: "schedule" }}
            size={12}
            tintColor={colors.guard.textMuted}
          />
          <Text style={styles.blockTime}>{block.time}</Text>
        </View>
      ) : null}
    </View>
  );
}

export function VisitorTimelineRow({ blocks, entry }: VisitorTimelineRowProps) {
  const timeline = blocks ?? getVisitorTimelineBlocks(entry);

  if (timeline.length === 0) {
    return (
      <View style={styles.emptyState}>
        <SymbolView
          name={{ ios: "doc.text", android: "description", web: "description" }}
          size={16}
          tintColor={colors.guard.textMuted}
        />
        <Text style={styles.emptyText}>No visit times recorded yet.</Text>
      </View>
    );
  }

  return (
    <View style={styles.panel}>
      {timeline.map((block, index) => (
        <View key={`${block.label}-${index}`} style={styles.blockColumn}>
          {index > 0 ? <View style={styles.divider} /> : null}
          <TimelineBlock block={block} />
        </View>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  block: {
    alignItems: "center",
    flex: 1,
    gap: 2,
    minWidth: 0,
    paddingHorizontal: spacing.xs,
    paddingVertical: spacing.sm,
  },
  blockColumn: {
    alignItems: "stretch",
    flex: 1,
    flexDirection: "row",
    minWidth: 0,
  },
  blockDate: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
    textAlign: "center",
  },
  blockLabel: {
    color: colors.guard.textMuted,
    fontSize: 9,
    fontWeight: "700",
    letterSpacing: 0.5,
    textAlign: "center",
    textTransform: "uppercase",
  },
  blockTime: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
  },
  divider: {
    alignSelf: "stretch",
    backgroundColor: "rgba(16, 29, 54, 0.08)",
    marginVertical: spacing.sm,
    width: 1,
  },
  emptyState: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderColor: "rgba(16, 29, 54, 0.08)",
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    paddingVertical: spacing.sm,
  },
  emptyText: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "500",
  },
  panel: {
    backgroundColor: colors.surface.secondary,
    borderColor: "rgba(16, 29, 54, 0.08)",
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    overflow: "hidden",
  },
  timeRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: 4,
  },
});
