import { Pressable, StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type SegmentTabsProps<T extends string> = {
  compact?: boolean;
  options: { value: T; label: string }[];
  value: T;
  variant?: "filled" | "underline";
  onChange: (value: T) => void;
};

export function SegmentTabs<T extends string>({
  compact = false,
  options,
  value,
  variant = "filled",
  onChange,
}: SegmentTabsProps<T>) {
  if (variant === "underline") {
    return (
      <View style={styles.underlineContainer}>
        {options.map((option) => {
          const active = option.value === value;

          return (
            <Pressable
              key={option.value}
              accessibilityRole="button"
              accessibilityState={{ selected: active }}
              style={styles.underlineTab}
              onPress={() => onChange(option.value)}
            >
              <Text
                style={[
                  styles.underlineTabText,
                  compact && styles.underlineTabTextCompact,
                  active && styles.underlineTabTextActive,
                ]}
              >
                {option.label}
              </Text>
              {active ? (
                <View style={styles.underlineIndicator} />
              ) : (
                <View style={styles.underlineSpacer} />
              )}
            </Pressable>
          );
        })}
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {options.map((option) => {
        const active = option.value === value;

        return (
          <Pressable
            key={option.value}
            accessibilityRole="button"
            accessibilityState={{ selected: active }}
            style={[
              styles.tab,
              active && styles.tabActive,
              { paddingVertical: compact ? spacing.sm : 10 },
            ]}
            onPress={() => onChange(option.value)}
          >
            <Text
              style={[
                styles.tabText,
                active && styles.tabTextActive,
                { fontSize: compact ? 11 : 13 },
              ]}
            >
              {option.label}
            </Text>
          </Pressable>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.md,
    borderWidth: 1,
    flexDirection: "row",
    padding: spacing.xs,
  },
  tab: {
    minHeight: 44,
    justifyContent: "center",
    alignItems: "center",
    borderRadius: radius.sm,
    flex: 1,
  },
  tabActive: {
    backgroundColor: colors.brand.navy,
  },
  tabText: {
    textAlign: "center",
    color: colors.guard.textMuted,
    fontWeight: "600",
  },
  tabTextActive: {
    color: colors.text.inverse,
  },
  underlineContainer: {
    borderBottomColor: colors.guard.border,
    borderBottomWidth: 1,
    flexDirection: "row",
  },
  underlineIndicator: {
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 3,
    marginTop: spacing.sm,
    width: "100%",
  },
  underlineSpacer: {
    height: 3,
    marginTop: spacing.sm,
  },
  underlineTab: {
    minHeight: 44,
    justifyContent: "center",
    alignItems: "center",
    flex: 1,
    paddingBottom: spacing.xs,
    paddingTop: spacing.xs,
  },
  underlineTabText: {
    color: colors.text.muted,
    fontSize: 14,
    fontWeight: "600",
  },
  underlineTabTextActive: {
    color: colors.brand.navy,
    fontWeight: "700",
  },
  underlineTabTextCompact: {
    fontSize: 13,
  },
});
