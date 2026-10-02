import { useCallback, useEffect, useMemo, useState } from "react";
import {
  Animated,
  AccessibilityInfo,
  Easing,
  Pressable,
  StyleSheet,
  Text,
  View,
} from "react-native";
import Svg, { Circle, Line, Path } from "react-native-svg";
import ToastMessage from "react-native-toast-message";
import { useSafeAreaInsets } from "react-native-safe-area-context";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

// ── Types ──────────────────────────────────────────────────────────────────

export type ToastVariant = "error" | "success" | "info" | "warning";

type ToastConfigParams = {
  type: string;
  text1?: string;
  text2?: string;
  hide: () => void;
  props?: { duration?: number; id?: number };
};

type ShowToastOptions = {
  title: string;
  message?: string;
  variant?: ToastVariant;
  duration?: number;
};

const DEFAULT_DURATION = 3500;
let toastID = 0;

// Distinct accent per variant — saturated on purpose, kept local to the
// toast so the rest of the app's palette stays untouched.
const VARIANT: Record<ToastVariant, { accent: string }> = {
  error: { accent: "#FF5A6E" },
  success: { accent: "#33DDA8" },
  info: { accent: "#7C9CFF" },
  warning: { accent: "#FFB648" },
};

// ── Icons ──────────────────────────────────────────────────────────────────

function ToastGlyph({
  variant,
  size = 15,
}: {
  variant: ToastVariant;
  size?: number;
}) {
  const stroke = "#0B0B10";

  switch (variant) {
    case "success":
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24" fill="none">
          <Path
            d="M5 12.5L10 17.5L19.5 7"
            stroke={stroke}
            strokeWidth={2.6}
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </Svg>
      );
    case "error":
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24" fill="none">
          <Line
            x1="6"
            y1="6"
            x2="18"
            y2="18"
            stroke={stroke}
            strokeWidth={2.6}
            strokeLinecap="round"
          />
          <Line
            x1="18"
            y1="6"
            x2="6"
            y2="18"
            stroke={stroke}
            strokeWidth={2.6}
            strokeLinecap="round"
          />
        </Svg>
      );
    case "warning":
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24" fill="none">
          <Line
            x1="13.4"
            y1="5.5"
            x2="11.6"
            y2="15"
            stroke={stroke}
            strokeWidth={2.6}
            strokeLinecap="round"
          />
          <Circle cx="10.9" cy="18.4" r="1.5" fill={stroke} />
        </Svg>
      );
    case "info":
    default:
      return (
        <Svg width={size} height={size} viewBox="0 0 24 24" fill="none">
          <Circle cx="14.1" cy="6.6" r="1.6" fill={stroke} />
          <Line
            x1="13.4"
            y1="10.4"
            x2="10.6"
            y2="18.4"
            stroke={stroke}
            strokeWidth={2.6}
            strokeLinecap="round"
          />
        </Svg>
      );
  }
}

// ── Card ───────────────────────────────────────────────────────────────────

function ToastCard({ type, text1, text2, hide, props }: ToastConfigParams) {
  const variant = (
    VARIANT[type as ToastVariant] ? type : "info"
  ) as ToastVariant;
  const duration = props?.duration ?? DEFAULT_DURATION;
  const [progress] = useState(() => new Animated.Value(0));

  // The library can reuse the card; restart even when duration and text match.
  useEffect(() => {
    progress.setValue(0);
    const animation = Animated.timing(progress, {
      toValue: 1,
      duration,
      easing: Easing.linear,
      useNativeDriver: false,
    });
    animation.start();
    return () => animation.stop();
  }, [duration, progress, props?.id]);

  const barWidth = useMemo(
    () =>
      progress.interpolate({
        inputRange: [0, 1],
        outputRange: ["100%", "0%"],
      }),
    [progress],
  );

  return (
    <View style={styles.container}>
      <View style={[styles.toast, { shadowColor: VARIANT[variant].accent }]}>
        <View style={styles.row}>
          <View
            style={[
              styles.iconBadge,
              {
                backgroundColor: VARIANT[variant].accent,
                shadowColor: VARIANT[variant].accent,
              },
            ]}
          >
            <ToastGlyph variant={variant} />
          </View>

          <View style={styles.textWrap}>
            {text1 ? (
              <Text style={styles.title} numberOfLines={2}>
                {text1}
              </Text>
            ) : null}
            {text2 ? (
              <Text style={styles.message} numberOfLines={3}>
                {text2}
              </Text>
            ) : null}
          </View>

          <Pressable
            accessibilityRole="button"
            accessibilityLabel="Dismiss notification"
            onPress={hide}
            hitSlop={10}
            style={styles.closeButton}
          >
            <Svg width={11} height={11} viewBox="0 0 24 24" fill="none">
              <Line
                x1="6"
                y1="6"
                x2="18"
                y2="18"
                stroke={colors.text.ghost}
                strokeWidth={2}
                strokeLinecap="round"
              />
              <Line
                x1="18"
                y1="6"
                x2="6"
                y2="18"
                stroke={colors.text.ghost}
                strokeWidth={2}
                strokeLinecap="round"
              />
            </Svg>
          </Pressable>
        </View>

        <View style={styles.trackWrap}>
          <Animated.View
            style={[
              styles.track,
              { width: barWidth, backgroundColor: VARIANT[variant].accent },
            ]}
          />
        </View>
      </View>
    </View>
  );
}

// ── Config + mount point ────────────────────────────────────────────────────

export const toastConfig = {
  success: (p: ToastConfigParams) => <ToastCard {...p} />,
  error: (p: ToastConfigParams) => <ToastCard {...p} />,
  info: (p: ToastConfigParams) => <ToastCard {...p} />,
  warning: (p: ToastConfigParams) => <ToastCard {...p} />,
};

/** Mount once at the root of the app, e.g. as the last sibling in App.tsx / root layout. */
export function AppToast() {
  const insets = useSafeAreaInsets();
  return (
    <ToastMessage
      config={toastConfig}
      position="top"
      topOffset={insets.top + 12}
    />
  );
}

// ── Hook (keeps the existing call-site API) ────────────────────────────────

export function useToast() {
  const showToast = useCallback((options: ShowToastOptions) => {
    const {
      title,
      message,
      variant = "info",
      duration = DEFAULT_DURATION,
    } = options;

    AccessibilityInfo.announceForAccessibility(
      [title, message].filter(Boolean).join(". "),
    );
    ToastMessage.show({
      type: variant,
      text1: title,
      text2: message,
      position: "top",
      visibilityTime: duration,
      props: { duration, id: ++toastID },
    });
  }, []);

  const hideToast = useCallback(() => ToastMessage.hide(), []);

  return { showToast, hideToast };
}

// ── Styles ───────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  container: {
    width: "92%",
    maxWidth: 420,
    alignSelf: "center",
  },
  toast: {
    borderRadius: radius.lg ?? 18,
    backgroundColor: "rgba(18, 18, 24, 0.94)",
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: "rgba(255, 255, 255, 0.10)",
    overflow: "hidden",
    elevation: 10,
    shadowOffset: { width: 0, height: 14 },
    shadowOpacity: 0.32,
    shadowRadius: 24,
  },
  row: {
    flexDirection: "row",
    alignItems: "flex-start",
    paddingHorizontal: 14,
    paddingVertical: spacing.md,
    gap: 12,
  },
  iconBadge: {
    width: 26,
    height: 26,
    borderRadius: 13,
    alignItems: "center",
    justifyContent: "center",
    marginTop: 1,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.55,
    shadowRadius: 8,
    elevation: 4,
  },
  textWrap: {
    flex: 1,
    gap: 3,
    paddingTop: 1,
  },
  title: {
    ...typography.bodySmall,
    color: "#F5F5F8",
    fontWeight: "700",
    letterSpacing: 0.15,
  },
  message: {
    ...typography.bodySmall,
    color: "rgba(245, 245, 248, 0.62)",
    fontSize: 13,
    lineHeight: 18,
  },
  closeButton: {
    minWidth: 44,
    minHeight: 44,
    alignItems: "center",
    justifyContent: "center",
    paddingHorizontal: 4,
    paddingVertical: 4,
    marginTop: 2,
  },
  trackWrap: {
    height: 2,
    width: "100%",
    backgroundColor: "rgba(255, 255, 255, 0.06)",
  },
  track: {
    height: "100%",
    borderRadius: 2,
  },
});
