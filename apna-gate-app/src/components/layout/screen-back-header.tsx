import type { ReactNode } from "react";
import type { Href } from "expo-router";
import { useFocusEffect } from "expo-router";
import { useCallback } from "react";
import { BackHandler, Platform, Pressable, StyleSheet, Text } from "react-native";

import { AppIcon } from "@/components/icons";
import { Row } from "@/components/layout/row";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type ScreenBackHeaderProps = {
  fallbackHomeRoute: Href;
  title: string;
  trailing?: ReactNode;
};

export function ScreenBackHeader({
  fallbackHomeRoute,
  title,
  trailing,
}: ScreenBackHeaderProps) {
  const handleBack = useBackAction(fallbackHomeRoute);

  useFocusEffect(
    useCallback(() => {
      if (Platform.OS !== "android") {
        return undefined;
      }
      const subscription = BackHandler.addEventListener("hardwareBackPress", () => {
        handleBack();
        return true;
      });
      return () => subscription.remove();
    }, [handleBack]),
  );

  return (
    <Row align="center" gap="md" style={styles.header}>
      <Pressable
        accessibilityLabel="Go back"
        accessibilityRole="button"
        hitSlop={8}
        style={({ pressed }) => [
          styles.backButton,
          pressed && styles.backButtonPressed,
        ]}
        onPress={handleBack}
      >
        <AppIcon color={colors.guard.text} name="back" size={18} />
      </Pressable>
      <Text style={styles.title}>{title}</Text>
      {trailing}
    </Row>
  );
}

const styles = StyleSheet.create({
  backButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  backButtonPressed: {
    backgroundColor: colors.guard.tealSoft,
  },
  header: {
    justifyContent: "flex-start",
    paddingBottom: spacing.sm,
    paddingTop: spacing.xs,
  },
  title: {
    color: colors.guard.text,
    flex: 1,
    fontSize: 18,
    fontWeight: "700",
  },
});
