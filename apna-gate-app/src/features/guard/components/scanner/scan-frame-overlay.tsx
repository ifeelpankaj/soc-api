import { useEffect } from "react";
import { Pressable, StyleSheet, Text, View, type ViewStyle } from "react-native";
import { SymbolView } from "expo-symbols";
import Animated, {
  Easing,
  useAnimatedStyle,
  useSharedValue,
  withRepeat,
  withTiming,
} from "react-native-reanimated";
import Svg, { Defs, LinearGradient, Rect, Stop } from "react-native-svg";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

const FRAME_SIZE = 248;
const CORNER_SIZE = 28;
const CORNER_STROKE = 4;

type ScanFrameOverlayProps = {
  flashEnabled: boolean;
  onToggleFlash: () => void;
};

function ScanCorner({ style }: { style: ViewStyle }) {
  return <View style={[styles.corner, style]} />;
}

function ScanCorners() {
  const cornerStyle = {
    borderColor: colors.text.inverse,
    width: CORNER_SIZE,
    height: CORNER_SIZE,
  };

  return (
    <>
      <ScanCorner
        style={{
          ...cornerStyle,
          borderTopWidth: CORNER_STROKE,
          borderLeftWidth: CORNER_STROKE,
          borderTopLeftRadius: 12,
          top: 0,
          left: 0,
        }}
      />
      <ScanCorner
        style={{
          ...cornerStyle,
          borderTopWidth: CORNER_STROKE,
          borderRightWidth: CORNER_STROKE,
          borderTopRightRadius: 12,
          top: 0,
          right: 0,
        }}
      />
      <ScanCorner
        style={{
          ...cornerStyle,
          borderBottomWidth: CORNER_STROKE,
          borderLeftWidth: CORNER_STROKE,
          borderBottomLeftRadius: 12,
          bottom: 0,
          left: 0,
        }}
      />
      <ScanCorner
        style={{
          ...cornerStyle,
          borderBottomWidth: CORNER_STROKE,
          borderRightWidth: CORNER_STROKE,
          borderBottomRightRadius: 12,
          bottom: 0,
          right: 0,
        }}
      />
    </>
  );
}

function ScanLine() {
  const translateY = useSharedValue(12);

  useEffect(() => {
    translateY.value = withRepeat(
      withTiming(FRAME_SIZE - 20, {
        duration: 2200,
        easing: Easing.inOut(Easing.quad),
      }),
      -1,
      true,
    );
  }, [translateY]);

  const animatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateY: translateY.value }],
  }));

  return (
    <Animated.View style={[styles.scanLineWrap, animatedStyle]}>
      <Svg height={3} width={FRAME_SIZE - 32}>
        <Defs>
          <LinearGradient id="scanLine" x1="0%" x2="100%" y1="0%" y2="0%">
            <Stop offset="0%" stopColor="rgba(255,106,26,0)" />
            <Stop offset="50%" stopColor="#FF6A1A" />
            <Stop offset="100%" stopColor="rgba(255,106,26,0)" />
          </LinearGradient>
        </Defs>
        <Rect fill="url(#scanLine)" height={3} rx={1.5} width={FRAME_SIZE - 32} x="0" y="0" />
      </Svg>
    </Animated.View>
  );
}

export function ScanFrameOverlay({ flashEnabled, onToggleFlash }: ScanFrameOverlayProps) {
  return (
    <View pointerEvents="box-none" style={styles.overlay}>
      <View style={styles.dimTop} />
      <View style={styles.middleRow}>
        <View style={styles.dimSide} />
        <View style={styles.frameWindow}>
          <ScanCorners />
          <ScanLine />
        </View>
        <View style={styles.dimSide} />
      </View>
      <View style={styles.dimBottom} />

      <View pointerEvents="box-none" style={styles.topInstructionWrap}>
        <View style={styles.instructionPill}>
          <SymbolView
            name={{ ios: "qrcode", android: "qr_code", web: "qr_code" }}
            size={16}
            tintColor={colors.brand.orange}
          />
          <Text style={styles.instructionText}>Align QR code within the frame</Text>
        </View>
      </View>

      <View pointerEvents="box-none" style={styles.bottomFlashWrap}>
        <Pressable
          accessibilityLabel={flashEnabled ? "Turn flash off" : "Turn flash on"}
          accessibilityRole="button"
          style={({ pressed }) => [styles.flashPill, pressed && styles.flashPillPressed]}
          onPress={onToggleFlash}
        >
          <SymbolView
            name={{ ios: "bolt.fill", android: "flash_on", web: "flash_on" }}
            size={16}
            tintColor={colors.brand.orange}
          />
          <Text style={styles.flashPillText}>{flashEnabled ? "Flash On" : "Flash Off"}</Text>
        </Pressable>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  bottomFlashWrap: {
    alignItems: "center",
    bottom: spacing.xl,
    left: 0,
    position: "absolute",
    right: 0,
  },
  corner: {
    position: "absolute",
  },
  dimBottom: {
    backgroundColor: "rgba(10, 10, 10, 0.72)",
    flex: 1,
  },
  dimSide: {
    backgroundColor: "rgba(10, 10, 10, 0.72)",
    flex: 1,
  },
  dimTop: {
    backgroundColor: "rgba(10, 10, 10, 0.72)",
    flex: 1,
  },
  flashPill: {
    alignItems: "center",
    backgroundColor: "rgba(20, 20, 20, 0.88)",
    borderRadius: 999,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
  },
  flashPillPressed: {
    opacity: 0.85,
  },
  flashPillText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "700",
  },
  frameWindow: {
    height: FRAME_SIZE,
    overflow: "hidden",
    position: "relative",
    width: FRAME_SIZE,
  },
  instructionPill: {
    alignItems: "center",
    backgroundColor: "rgba(20, 20, 20, 0.88)",
    borderRadius: 999,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
  },
  instructionText: {
    color: colors.text.inverse,
    fontSize: 13,
    fontWeight: "600",
  },
  middleRow: {
    flexDirection: "row",
    height: FRAME_SIZE,
  },
  overlay: {
    ...StyleSheet.absoluteFill,
  },
  scanLineWrap: {
    left: 16,
    position: "absolute",
    right: 16,
  },
  topInstructionWrap: {
    alignItems: "center",
    left: 0,
    position: "absolute",
    right: 0,
    top: spacing.xl,
  },
});
