import { useCallback, useEffect, useState } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import Animated, {
  FadeIn,
  FadeOut,
  LinearTransition,
} from "react-native-reanimated";
import { SymbolView } from "expo-symbols";

import { DashboardSection } from "@/components/dashboard/dashboard-section";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

import {
  formatAttentionCount,
  hiddenAttentionCount,
  shouldShowExpandControl,
  shouldShowNeedsAttention,
  visibleAttentionItems,
} from "./needs-attention-logic";
import { NeedsAttentionItem } from "./needs-attention-item";
import type { AttentionItem } from "./types";

type NeedsAttentionSectionProps = {
  items: AttentionItem[];
};

export function NeedsAttentionSection({ items }: NeedsAttentionSectionProps) {
  const [expanded, setExpanded] = useState(false);
  const showExpand = shouldShowExpandControl(items);
  const hiddenCount = hiddenAttentionCount(items, expanded);
  const visibleItems = visibleAttentionItems(items, expanded);

  useEffect(() => {
    if (!showExpand) {
      setExpanded(false);
    }
  }, [showExpand]);

  const toggleExpanded = useCallback(() => {
    setExpanded((current) => !current);
  }, []);

  if (!shouldShowNeedsAttention(items)) {
    return null;
  }

  return (
    <DashboardSection
      title="Needs Attention"
      trailing={formatAttentionCount(items.length)}
    >
      <View style={styles.card}>
        <Animated.View layout={LinearTransition.duration(200)}>
          {visibleItems.map((item, index) => (
            <Animated.View
              key={item.id}
              entering={FadeIn.duration(200)}
              exiting={FadeOut.duration(180)}
              layout={LinearTransition.duration(200)}
            >
              <NeedsAttentionItem
                isLast={index === visibleItems.length - 1 && hiddenCount === 0}
                item={item}
              />
            </Animated.View>
          ))}

          {hiddenCount > 0 ? (
            <Pressable
              accessibilityRole="button"
              accessibilityLabel={`${hiddenCount} more items`}
              style={({ pressed }) => [styles.moreRow, pressed && styles.rowPressed]}
              onPress={toggleExpanded}
            >
              <Text style={styles.moreText}>{hiddenCount} more</Text>
              <SymbolView
                name={{ ios: "chevron.down", android: "expand_more", web: "expand_more" }}
                size={16}
                tintColor={colors.text.muted}
              />
            </Pressable>
          ) : null}

          {showExpand && expanded && hiddenCount === 0 ? (
            <Pressable
              accessibilityRole="button"
              accessibilityLabel="Show less"
              style={({ pressed }) => [styles.moreRow, pressed && styles.rowPressed]}
              onPress={toggleExpanded}
            >
              <Text style={styles.moreText}>Show less</Text>
              <SymbolView
                name={{ ios: "chevron.up", android: "expand_less", web: "expand_less" }}
                size={16}
                tintColor={colors.text.muted}
              />
            </Pressable>
          ) : null}
        </Animated.View>
      </View>
    </DashboardSection>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    overflow: "hidden",
    ...shadows.sm,
  },
  moreRow: {
    alignItems: "center",
    borderTopColor: colors.dashboard.cardBorder,
    borderTopWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    paddingVertical: spacing.md,
  },
  moreText: {
    color: colors.text.secondary,
    fontSize: 14,
    fontWeight: "600",
  },
  rowPressed: {
    backgroundColor: colors.surface.muted,
  },
});
