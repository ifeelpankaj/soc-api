import { useContext, useState, useSyncExternalStore } from "react";
import { useLocalSearchParams, useRouter } from "expo-router";
import { Pressable, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/features/auth/use-auth";
import { BrowserPushContext, type BrowserPushState } from "./notification-context";
import {
  dismissAlertBanner,
  isAlertBannerDismissed,
  subscribeAlertBanner,
} from "./notification-banner-session";

const messages: Record<BrowserPushState, string> = {
  loading: "Checking notification availability…",
  install: "On iPhone or iPad, open Apna Gate in Safari, tap Share, then Add to Home Screen. Open it from that icon and return to notification settings to enable alerts. Background alerts require iOS/iPadOS 16.4 or later.",
  setup: "Browser alerts are not configured yet. Your inbox updates while this app is open. Contact your society administrator for availability.",
  unsupported: "This browser cannot receive background alerts. Try an updated browser that supports notifications. On iPhone or iPad, open Apna Gate from your Home Screen. Your inbox still updates while open.",
  denied: "Notifications are blocked. Allow notifications for this site in your browser settings and, if needed, your device notification settings. Then return here. Your inbox still updates while open.",
  default: "Enable notifications, then choose Allow in your browser's permission prompt. On desktop and Android browsers, installation is optional: use the browser menu's Install app or Add to Home Screen option if available.",
  granted: "Permission is granted. Enable notifications for this signed-in resident account.",
  enabling: "Setting up notifications…",
  error: "Notifications could not be enabled. Check your connection and try again. Your inbox is still available.",
  enabled: "Alerts are enabled for this resident account. You can change permission in your browser or device notification settings. Delivery depends on your device and connection.",
};

export function WebNotificationBanner() {
  const context = useContext(BrowserPushContext);
  const { user } = useAuth();
  const router = useRouter();
  const dismissed = useSyncExternalStore(
    subscribeAlertBanner,
    () => isAlertBannerDismissed(user?.id),
    () => false,
  );

  if (!context || context.state === "loading" || context.state === "enabled" || dismissed) return null;

  return (
    <View style={styles.banner}>
      <Text style={styles.label}>Enable alerts</Text>
      <Pressable
        accessibilityRole="button"
        accessibilityLabel="Set up alerts"
        onPress={() => router.push({ pathname: "/notification-settings", params: { setup: "1" } })}
        style={styles.action}
      >
        <Text style={styles.actionText}>Set up</Text>
      </Pressable>
      <Pressable
        accessibilityRole="button"
        accessibilityLabel="Dismiss alerts reminder for this session"
        onPress={() => { if (user?.id) dismissAlertBanner(user.id); }}
        style={styles.action}
      >
        <Text style={styles.actionText}>×</Text>
      </Pressable>
    </View>
  );
}

export function WebNotificationSettings() {
  const context = useContext(BrowserPushContext);
  const { setup } = useLocalSearchParams<{ setup?: string }>();
  const router = useRouter();
  const [expanded, setExpanded] = useState(false);
  const showDetails = expanded || setup === "1";

  if (!context || context.state === "loading") return null;

  return (
    <View style={styles.settings}>
      <Pressable
        accessibilityRole="button"
        accessibilityState={{ expanded: showDetails }}
        onPress={() => {
          setExpanded(!showDetails);
          if (setup) router.setParams({ setup: undefined });
        }}
        style={styles.settingsToggle}
      >
        <Text style={styles.label}>Notification settings</Text>
        <Text style={styles.actionText}>{showDetails ? "Close" : "Set up"}</Text>
      </Pressable>
      {showDetails && (
        <View style={styles.details}>
          <Text accessibilityRole="header" style={styles.heading}>Notifications on this device</Text>
          <Text selectable accessibilityLiveRegion="polite" style={styles.message}>
            {messages[context.state]}
          </Text>
          {["default", "granted", "error", "enabling"].includes(context.state) && (
            <Button
              title={context.state === "error" ? "Try notifications again" : "Enable notifications"}
              loading={context.state === "enabling"}
              onPress={context.enable}
            />
          )}
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  banner: {
    alignItems: "center",
    backgroundColor: "#fff7ed",
    borderRadius: 12,
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 8,
    padding: 12,
  },
  label: { color: "#212121", flexGrow: 1, flexShrink: 1, fontSize: 14 },
  action: { alignItems: "center", justifyContent: "center", minHeight: 44, minWidth: 44, padding: 8 },
  actionText: { color: "#9a3412", fontSize: 14, fontWeight: "600", flexShrink: 1 },
  settings: { backgroundColor: "#fff7ed", borderRadius: 12, margin: 16, padding: 12 },
  settingsToggle: { alignItems: "center", flexDirection: "row", flexWrap: "wrap", gap: 12, minHeight: 44 },
  details: { gap: 12, paddingTop: 12 },
  heading: { color: "#212121", fontSize: 17, fontWeight: "700" },
  message: { color: "#333333", fontSize: 14, lineHeight: 21 },
});
