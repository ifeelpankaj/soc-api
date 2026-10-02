import type { PropsWithChildren, ReactNode } from "react";
import { useCallback } from "react";
import { useFocusEffect, type Href } from "expo-router";
import {
  ActivityIndicator,
  BackHandler,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

type RoleSubScreenProps = PropsWithChildren<{
  footer?: ReactNode;
  fallbackHomeRoute: Href;
  gate: ReactNode;
  headerExtra?: ReactNode;
  headerTrailing?: ReactNode;
  isLoading?: boolean;
  isReady?: boolean;
  showGate?: boolean;
  title: string;
}>;

export function RoleSubScreen({
  children,
  footer,
  fallbackHomeRoute,
  gate,
  headerExtra,
  headerTrailing,
  isLoading = false,
  isReady = true,
  showGate = false,
  title,
}: RoleSubScreenProps) {
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

  if (showGate) {
    return <>{gate}</>;
  }

  return (
    <SafeAreaView
      edges={footer ? ["top", "left", "right"] : ["top", "left", "right", "bottom"]}
      style={styles.screen}
    >
      <AppStatusBar />
      <View style={styles.header}>
        <ScreenBackHeader
          fallbackHomeRoute={fallbackHomeRoute}
          title={title}
          trailing={headerTrailing}
        />
        {headerExtra}
      </View>
      {isLoading || !isReady ? (
        <View style={styles.inlineLoading}>
          <ActivityIndicator color={colors.guard.teal} size="small" />
        </View>
      ) : (
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          style={styles.content}
        >
          <View style={[styles.content, styles.readingWidth]}>{children}</View>
          {footer}
        </KeyboardAvoidingView>
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  readingWidth: {
    width: "100%",
    maxWidth: 840,
    alignSelf: "center",
  },
  content: {
    flex: 1,
  },
  header: {
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: layout.screenPaddingTop,
  },
  inlineLoading: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
    paddingVertical: spacing["3xl"],
  },
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
});
