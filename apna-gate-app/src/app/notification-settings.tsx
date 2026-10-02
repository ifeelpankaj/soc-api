import { Redirect } from "expo-router";
import { ScrollView, StyleSheet } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { LoadingState } from "@/components/ui";
import { useAuth } from "@/features/auth/use-auth";
import { NotificationPushPreferencesPanel } from "@/features/notifications/notification-push-preferences";
import { WebNotificationSettings } from "@/features/web/notification-settings";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

export default function NotificationSettingsScreen() {
  const { homeRoute, status } = useAuth();
  if (status === "unauthenticated") return <Redirect href="/login" />;
  if (status !== "authenticated") return <LoadingState message="Loading notification settings" />;

  return (
    <SafeAreaView style={styles.screen}>
      <AppStatusBar />
      <ScrollView contentContainerStyle={styles.content}>
        <ScreenBackHeader fallbackHomeRoute={homeRoute ?? "/"} title="Notification settings" />
        <WebNotificationSettings />
        <NotificationPushPreferencesPanel />
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.guard.screenBg },
  content: {
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: layout.screenPaddingTop,
    paddingBottom: spacing["3xl"],
    gap: spacing.md,
  },
});
