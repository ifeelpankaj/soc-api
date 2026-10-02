import { Image } from "expo-image";
import { ActivityIndicator, Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import Svg, { Path, Rect } from "react-native-svg";

import { Stack } from "@/components/layout";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type VisitorQueueEmptyStateProps = {
  message: string;
  onRefresh?: () => void;
  refreshing?: boolean;
  title: string;
};

function GateEmptyIllustration() {
  return (
    <View style={styles.illustrationWrap}>
      <Svg height={120} viewBox="0 0 160 120" width={160}>
        <Rect fill="#E8EAEF" height={28} rx={4} width={16} x={18} y="24" />
        <Rect fill="#E8EAEF" height={28} rx={4} width={16} x="126" y="24" />
        <Path
          d="M34 52 H126 V88 C126 94 121 98 115 98 H45 C39 98 34 94 34 88 Z"
          fill="#DDE1E8"
        />
        <Path
          d="M42 52 H118 V82 C118 86 114 90 110 90 H50 C46 90 42 86 42 82 Z"
          fill="#ECEFF4"
        />
        <Path d="M52 52 H108 V72 C108 74 106 76 104 76 H56 C54 76 52 74 52 72 Z" fill="#F3F5F8" />
        <Rect fill="#DDE1E8" height={10} rx={2} width={10} x="75" y="16" />
      </Svg>
      <View style={styles.personBadge}>
        <SymbolView
          name={{ ios: "person.fill", android: "person", web: "person" }}
          size={28}
          tintColor={colors.brand.orange}
        />
      </View>
    </View>
  );
}

export function VisitorQueueEmptyState({
  message,
  onRefresh,
  refreshing = false,
  title,
}: VisitorQueueEmptyStateProps) {
  return (
    <View style={styles.wrap}>
      <Image
        accessibilityElementsHidden
        importantForAccessibility="no-hide-descendants"
        pointerEvents="none"
        source={require("@/assets/images/public/logo.png")}
        style={styles.backgroundLogo}
        contentFit="contain"
      />

      <Stack align="center" gap="lg" style={styles.content}>
        <GateEmptyIllustration />

        <Stack align="center" gap="sm" style={styles.copy}>
          <Text style={styles.title}>{title}</Text>
          <Text style={styles.message}>{message}</Text>
        </Stack>

        {onRefresh ? (
          <Pressable
            accessibilityRole="button"
            disabled={refreshing}
            style={({ pressed }) => [
              styles.refreshButton,
              pressed && !refreshing && styles.refreshButtonPressed,
              refreshing && styles.refreshButtonDisabled,
            ]}
            onPress={onRefresh}
          >
            {refreshing ? (
              <ActivityIndicator color={colors.text.inverse} size="small" />
            ) : (
              <>
                <SymbolView
                  name={{ ios: "arrow.clockwise", android: "refresh", web: "refresh" }}
                  size={18}
                  tintColor={colors.text.inverse}
                />
                <Text style={styles.refreshButtonText}>Refresh</Text>
              </>
            )}
          </Pressable>
        ) : null}
      </Stack>
    </View>
  );
}

const styles = StyleSheet.create({
  backgroundLogo: {
    alignSelf: "center",
    height: 280,
    opacity: 0.07,
    position: "absolute",
    top: 24,
    transform: [{ scaleX: 1.18 }, { scaleY: 1.12 }],
    width: "115%",
  },
  content: {
    alignItems: "center",
    alignSelf: "stretch",
    width: "100%",
  },
  copy: {
    width: "100%",
  },
  illustrationWrap: {
    alignItems: "center",
    height: 132,
    justifyContent: "center",
    width: 180,
  },
  message: {
    color: colors.guard.textMuted,
    fontSize: 14,
    lineHeight: 21,
    textAlign: "center",
  },
  personBadge: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderColor: "#F9DCC8",
    borderRadius: 999,
    borderWidth: 2,
    height: 56,
    justifyContent: "center",
    position: "absolute",
    top: 34,
    width: 56,
    ...shadows.sm,
  },
  refreshButton: {
    alignItems: "center",
    alignSelf: "stretch",
    backgroundColor: colors.brand.orange,
    borderRadius: radius.xl,
    flexDirection: "row",
    gap: spacing.sm,
    height: layout.buttonHeight,
    justifyContent: "center",
    marginTop: spacing.sm,
    ...shadows.cta,
  },
  refreshButtonDisabled: {
    opacity: 0.85,
  },
  refreshButtonPressed: {
    opacity: 0.92,
  },
  refreshButtonText: {
    color: colors.text.inverse,
    fontSize: 16,
    fontWeight: "700",
  },
  title: {
    color: colors.brand.navy,
    fontSize: 22,
    fontWeight: "700",
    letterSpacing: -0.3,
    textAlign: "center",
  },
  wrap: {
    alignItems: "center",
    alignSelf: "stretch",
    flexGrow: 1,
    justifyContent: "center",
    minHeight: 420,
    overflow: "hidden",
    paddingBottom: spacing["3xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing["2xl"],
    position: "relative",
    width: "100%",
  },
});
