import { Pressable, StyleSheet, Text } from "react-native";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type DashboardErrorBannerProps = {
  message: string;
  onRetry: () => void;
};

export function DashboardErrorBanner({ message, onRetry }: DashboardErrorBannerProps) {
  return (
    <Pressable accessibilityRole="button" style={styles.errorBanner} onPress={onRetry}>
      <Text style={styles.errorBannerMessage}>{message}</Text>
      <Text style={styles.errorBannerAction}>Retry</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  errorBanner: {
    alignItems: "center",
    backgroundColor: colors.status.errorSoft,
    borderRadius: 14,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  errorBannerAction: {
    color: colors.status.error,
    fontSize: 13,
    fontWeight: "600",
  },
  errorBannerMessage: {
    color: colors.status.error,
    flex: 1,
    fontSize: 13,
    fontWeight: "500",
    paddingRight: spacing.md,
  },
});
