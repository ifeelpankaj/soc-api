import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";
import Svg, { Defs, LinearGradient as SvgLinearGradient, Path, Rect, Stop } from "react-native-svg";

import { Row } from "@/components/layout";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type DashboardHeroCardProps = {
  accessibilityLabel?: string;
  icon: SymbolViewProps["name"];
  onPress: () => void;
  subtitle: string;
  title: string;
};

function HeroGradientBackground() {
  return (
    <Svg height="100%" preserveAspectRatio="none" style={StyleSheet.absoluteFill} width="100%">
      <Defs>
        <SvgLinearGradient id="heroGradient" x1="0%" x2="100%" y1="0%" y2="0%">
          <Stop offset="0%" stopColor="#172A43" />
          <Stop offset="50%" stopColor="#203550" />
          <Stop offset="100%" stopColor="#29405D" />
        </SvgLinearGradient>
      </Defs>
      <Rect fill="url(#heroGradient)" height="100%" width="100%" x="0" y="0" />
    </Svg>
  );
}

function ArchDecoration() {
  return (
    <Svg height={88} style={styles.decoration} viewBox="0 0 80 88" width={80}>
      <Path
        d="M40 8 C20 8 8 28 8 48 L8 88 L72 88 L72 48 C72 28 60 8 40 8 Z"
        fill="rgba(255,255,255,0.12)"
      />
    </Svg>
  );
}

export function DashboardHeroCard({
  accessibilityLabel,
  icon,
  onPress,
  subtitle,
  title,
}: DashboardHeroCardProps) {
  return (
    <Pressable
      accessibilityLabel={accessibilityLabel ?? title}
      accessibilityRole="button"
      style={({ pressed }) => [styles.pressable, pressed && styles.pressed]}
      onPress={onPress}
    >
      <View style={styles.card}>
        <HeroGradientBackground />
        <ArchDecoration />
        <Row align="center" gap="md" justify="flex-start" style={styles.content}>
          <View style={styles.iconWrap}>
            <SymbolView name={icon} size={28} tintColor={colors.brand.orange} />
          </View>
          <View style={styles.copy}>
            <Text style={styles.title}>{title}</Text>
            <Text style={styles.subtitle}>{subtitle}</Text>
          </View>
          <View style={styles.arrowButton}>
            <SymbolView
              name={{ ios: "arrow.right", android: "arrow_forward", web: "arrow_forward" }}
              size={18}
              tintColor={colors.brand.orange}
            />
          </View>
        </Row>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  arrowButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  card: {
    backgroundColor: colors.brand.navy,
    borderRadius: radius["2xl"],
    minHeight: 88,
    overflow: "hidden",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    position: "relative",
    ...shadows.card,
  },
  content: {
    flex: 1,
    zIndex: 1,
  },
  copy: {
    flex: 1,
    gap: 4,
    minWidth: 0,
  },
  decoration: {
    bottom: -8,
    position: "absolute",
    right: spacing.lg,
    zIndex: 1,
  },
  iconWrap: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: radius.lg,
    height: 52,
    justifyContent: "center",
    width: 52,
  },
  pressable: {
    width: "100%",
  },
  pressed: {
    opacity: 0.92,
  },
  subtitle: {
    color: "rgba(255,255,255,0.88)",
    fontSize: 13,
    fontWeight: "500",
    lineHeight: 18,
  },
  title: {
    color: colors.text.inverse,
    fontSize: 20,
    fontWeight: "800",
    letterSpacing: -0.3,
  },
});
