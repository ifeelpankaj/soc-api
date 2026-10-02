import { StyleSheet, View } from "react-native";

import { AppText } from "@/components/ui/app-text";
import { SUBSCRIPTION_EXPIRED_BANNER_MESSAGE } from "@/features/auth/api-error";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type SubscriptionExpiredBannerProps = {
  message?: string;
  title?: string;
};

export function SubscriptionExpiredBanner({
  message = SUBSCRIPTION_EXPIRED_BANNER_MESSAGE,
  title = "Subscription expired",
}: SubscriptionExpiredBannerProps) {
  return (
    <View style={styles.container}>
      <AppText variant="caption" color="warning" style={styles.title}>
        {title}
      </AppText>
      <AppText variant="bodySmall" style={styles.message}>
        {message}
      </AppText>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.status.warningSoft,
    borderColor: "rgba(217, 119, 6, 0.25)",
    borderRadius: radius.xl,
    borderWidth: 1,
    paddingHorizontal: spacing.lg,
    paddingVertical: 14,
  },
  message: {
    color: colors.text.secondary,
    fontSize: 15,
    lineHeight: 22,
    marginTop: spacing.xs,
  },
  title: {
    color: colors.status.warning,
    fontWeight: "700",
    letterSpacing: 0.8,
    textTransform: "uppercase",
  },
});
