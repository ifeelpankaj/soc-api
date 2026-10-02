import { ActivityIndicator, StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { AppText } from "@/components/ui/app-text";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type LoadingStateProps = {
  message?: string;
};

export function LoadingState({ message = "Loading your workspace" }: LoadingStateProps) {
  return (
    <SafeAreaView
      edges={["top", "left", "right", "bottom"]}
      style={styles.safeArea}
    >
      <AppStatusBar />
      <View style={styles.container}>
        <ActivityIndicator color={colors.brand.orange} size="large" />
        <AppText variant="body" color="secondary" style={styles.message}>
          {message}
        </AppText>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    alignItems: "center",
    flex: 1,
    gap: spacing.lg,
    justifyContent: "center",
    paddingHorizontal: spacing["2xl"],
  },
  message: {
    textAlign: "center",
  },
  safeArea: {
    backgroundColor: colors.surface.screen,
    flex: 1,
  },
});
