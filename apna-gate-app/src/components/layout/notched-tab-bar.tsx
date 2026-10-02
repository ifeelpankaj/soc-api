import type { ReactNode } from "react";
import { Platform, Pressable, StyleSheet, Text, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import Svg, { Path } from "react-native-svg";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

const BAR_HEIGHT = 64;
const FAB_SIZE = 58;
const NOTCH_RADIUS = FAB_SIZE / 2 + 10;
const CORNER_RADIUS = 28;

type TabSlot = {
  compactLabel?: boolean;
  accessibilityLabel?: string;
  active: boolean;
  icon: ReactNode;
  label: string;
  onPress: () => void;
};

type NotchedTabBarProps = {
  fab: {
    accessibilityLabel: string;
    icon: ReactNode;
    label: string;
    onPress: () => void;
  };
  leftTab: TabSlot;
  rightTab: TabSlot;
  additionalLeftTabs?: TabSlot[];
  additionalRightTabs?: TabSlot[];
};

function buildBarPath(width: number, height: number) {
  const cx = width / 2;
  const r = CORNER_RADIUS;
  const nr = NOTCH_RADIUS;

  return `
    M0,${r}
    Q0,0 ${r},0
    L${cx - nr - 14},0
    C${cx - nr + 4},0 ${cx - nr + 6},${nr * 0.85} ${cx - nr * 0.5},${nr * 0.85}
    C${cx - nr * 0.15},${nr} ${cx + nr * 0.15},${nr} ${cx + nr * 0.5},${nr * 0.85}
    C${cx + nr - 6},${nr * 0.85} ${cx + nr - 4},0 ${cx + nr + 14},0
    L${width - r},0
    Q${width},0 ${width},${r}
    L${width},${height - r}
    Q${width},${height} ${width - r},${height}
    L${r},${height}
    Q0,${height} 0,${height - r}
    Z
  `;
}

function TabButton({
  compactLabel,
  accessibilityLabel,
  active,
  icon,
  label,
  onPress,
}: TabSlot) {
  return (
    <Pressable
      accessibilityLabel={accessibilityLabel}
      accessibilityRole="tab"
      accessibilityState={{ selected: active }}
      style={styles.tab}
      onPress={onPress}
    >
      {icon}
      <Text
        numberOfLines={compactLabel ? 1 : undefined}
        adjustsFontSizeToFit={compactLabel}
        minimumFontScale={0.8}
        style={[
          styles.tabLabel,
          active ? styles.tabLabelActive : styles.tabLabelInactive,
        ]}
      >
        {label}
      </Text>
    </Pressable>
  );
}

export function NotchedTabBar({
  fab,
  leftTab,
  rightTab,
  additionalLeftTabs,
  additionalRightTabs,
}: NotchedTabBarProps) {
  const insets = useSafeAreaInsets();

  return (
    <View
      style={[
        styles.wrapper,
        { paddingBottom: Math.max(insets.bottom, spacing.sm) },
      ]}
    >
      <View style={styles.barContainer}>
        <Svg
          width="100%"
          height={BAR_HEIGHT}
          viewBox={`0 0 340 ${BAR_HEIGHT}`}
          style={StyleSheet.absoluteFill}
        >
          <Path
            d={buildBarPath(340, BAR_HEIGHT)}
            fill={colors.brand.orangeSoft ?? "#1B1F3B"}
          />
        </Svg>

        <View style={styles.row}>
          {additionalLeftTabs ? (
            <View style={{ flex: 1, flexDirection: "row" }}>
              <TabButton {...leftTab} compactLabel />
              {additionalLeftTabs.map((tab) => (
                <TabButton key={tab.label} {...tab} compactLabel />
              ))}
            </View>
          ) : (
            <TabButton {...leftTab} />
          )}
          <View style={styles.fabSpacer} />
          {additionalRightTabs ? (
            <View style={{ flex: 1, flexDirection: "row" }}>
              {additionalRightTabs.map((tab) => (
                <TabButton key={tab.label} {...tab} compactLabel />
              ))}
              <TabButton {...rightTab} compactLabel />
            </View>
          ) : (
            <TabButton {...rightTab} />
          )}
        </View>

        <Pressable
          accessibilityLabel={fab.accessibilityLabel}
          accessibilityRole="button"
          style={({ pressed }) => [styles.fab, pressed && styles.fabPressed]}
          onPress={fab.onPress}
        >
          {fab.icon}
        </Pressable>
        <Text style={styles.fabLabel}>{fab.label}</Text>
      </View>
    </View>
  );
}

export const notchedTabBarHeight = BAR_HEIGHT;

const styles = StyleSheet.create({
  barContainer: {
    alignSelf: "stretch",
    height: BAR_HEIGHT,
    width: "100%",
  },
  fab: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: FAB_SIZE / 2,
    height: FAB_SIZE,
    justifyContent: "center",
    left: "50%",
    marginLeft: -FAB_SIZE / 2,
    position: "absolute",
    top: -FAB_SIZE / 2 + 6,
    width: FAB_SIZE,
    ...(Platform.OS === "web"
      ? { boxShadow: "0 4px 8px rgba(0, 0, 0, 0.15)" }
      : {
          elevation: 4,
          shadowColor: "#000",
          shadowOffset: { height: 4, width: 0 },
          shadowOpacity: 0.15,
          shadowRadius: 8,
        }),
  },
  fabLabel: {
    color: colors.text.placeholder,
    fontSize: 11,
    fontWeight: "600",
    marginTop: 2,
    position: "absolute",
    textAlign: "center",
    top: BAR_HEIGHT - 18,
    width: "100%",
  },
  fabPressed: {
    opacity: 0.9,
    transform: [{ scale: 0.96 }],
  },
  fabSpacer: {
    width: FAB_SIZE + 20,
  },
  row: {
    flexDirection: "row",
    height: "100%",
  },
  tab: {
    alignItems: "center",
    flex: 1,
    gap: 4,
    justifyContent: "center",
  },
  tabLabel: {
    fontSize: 11,
    fontWeight: "600",
  },
  tabLabelActive: {
    color: colors.brand.orange,
  },
  tabLabelInactive: {
    color: colors.text.placeholder,
  },
  wrapper: {
    alignItems: "center",
    backgroundColor: "transparent",
    paddingHorizontal: spacing.md,
    paddingTop: spacing.sm,
  },
});
