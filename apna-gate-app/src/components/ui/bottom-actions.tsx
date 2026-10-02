import { type PropsWithChildren, useEffect, useState } from "react";
import { Keyboard, StyleSheet, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

export function BottomActions({ children }: PropsWithChildren) {
  const insets = useSafeAreaInsets();
  const [keyboardVisible, setKeyboardVisible] = useState(Keyboard.isVisible());
  useEffect(() => {
    const show = Keyboard.addListener("keyboardDidShow", () =>
      setKeyboardVisible(true),
    );
    const hide = Keyboard.addListener("keyboardDidHide", () =>
      setKeyboardVisible(false),
    );
    return () => {
      show.remove();
      hide.remove();
    };
  }, []);

  return (
    <View
      style={[
        styles.footer,
        {
          paddingBottom: keyboardVisible
            ? spacing.md
            : Math.max(insets.bottom, spacing.md),
        },
      ]}
    >
      <View style={styles.actions}>{children}</View>
    </View>
  );
}

const styles = StyleSheet.create({
  actions: {
    width: "100%",
    maxWidth: 800,
    alignSelf: "center",
    gap: spacing.sm,
  },
  footer: {
    flexShrink: 0,
    gap: spacing.sm,
    backgroundColor: colors.surface.card,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderTopColor: colors.border.default,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
});
