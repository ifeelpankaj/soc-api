import { AppToast } from "@/components/ui";
import { AppErrorBoundary } from "@/components/layout/app-error-boundary";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { AuthProvider } from "@/features/auth/auth-provider";
import { store } from "@/redux/store";
import { Stack } from "expo-router";
import { StyleSheet } from "react-native";
import { GestureHandlerRootView } from "react-native-gesture-handler";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { Provider } from "react-redux";
import "@/lib/api/enhanced-api";
import "@/lib/api/resident-api-extensions";
import "@/lib/api/guard-api-extensions";
import "@/lib/api/auth-api-extensions";
import "@/lib/api/notification-api-extensions";
import { NotificationRegistrationProvider } from "@/features/notifications/notification-registration-provider";
import { NotificationSyncProvider } from "@/features/notifications/notification-sync-provider";
import { NavigationHistoryProvider } from "@/lib/navigation/navigation-history";
import { WebAccessGate } from "@/features/web/web-access-gate";

export default function RootLayout() {
  return (
    <GestureHandlerRootView style={styles.root}>
      <SafeAreaProvider>
        <Provider store={store}>
          <WebAccessGate>
            <AuthProvider>
              <NotificationRegistrationProvider>
                <NavigationHistoryProvider>
                  <NotificationSyncProvider>
                    <AppErrorBoundary>
                      <AppStatusBar />
                      <Stack screenOptions={{ headerShown: false }} />
                    </AppErrorBoundary>
                  </NotificationSyncProvider>
                </NavigationHistoryProvider>
              </NotificationRegistrationProvider>
            </AuthProvider>
          </WebAccessGate>
          <AppToast />
        </Provider>
      </SafeAreaProvider>
    </GestureHandlerRootView>
  );
}

const styles = StyleSheet.create({
  root: {
    flex: 1,
  },
});
