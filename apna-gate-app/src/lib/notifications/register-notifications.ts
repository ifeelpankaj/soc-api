import * as Device from "expo-device";
import * as Notifications from "expo-notifications";
import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

export type DevicePlatform = "ios" | "android";

const INSTALLATION_ID_KEY = "apna_gate_installation_id";

export const VISITOR_DECISION_CATEGORY = "visitor_decision";
export const NOTIFICATION_INFO_CATEGORY = "notification_info";
export const APPROVE_VISITOR_ACTION = "VISITOR_APPROVE";
export const DECLINE_VISITOR_ACTION = "VISITOR_DECLINE";
export const MARK_NOTIFICATION_READ_ACTION = "NOTIFICATION_MARK_READ";
export const CLOSE_NOTIFICATION_ACTION = "NOTIFICATION_CLOSE";

export function getDevicePlatform(): DevicePlatform {
  return Platform.OS === "ios" ? "ios" : "android";
}

export async function configureNotifications() {
  if (Platform.OS === "web") return;
  Notifications.setNotificationHandler({
    handleNotification: async () => ({
      shouldShowBanner: true,
      shouldShowList: true,
      shouldPlaySound: true,
      shouldSetBadge: false,
    }),
  });

  if (Platform.OS === "android") {
    await Notifications.setNotificationChannelAsync("default", {
      name: "Default",
      importance: Notifications.AndroidImportance.HIGH,
      vibrationPattern: [0, 250, 250, 250],
      lightColor: "#208AEF",
    });
  }

  await Promise.all([
    Notifications.setNotificationCategoryAsync(VISITOR_DECISION_CATEGORY, [
      {
        identifier: APPROVE_VISITOR_ACTION,
        buttonTitle: "Approve",
        options: {
          isAuthenticationRequired: true,
          opensAppToForeground: false,
        },
      },
      {
        identifier: DECLINE_VISITOR_ACTION,
        buttonTitle: "Decline",
        options: {
          isAuthenticationRequired: true,
          isDestructive: true,
          opensAppToForeground: false,
        },
      },
    ]),
    Notifications.setNotificationCategoryAsync(NOTIFICATION_INFO_CATEGORY, [
      {
        identifier: MARK_NOTIFICATION_READ_ACTION,
        buttonTitle: "Mark as read",
        options: { opensAppToForeground: false },
      },
      {
        identifier: CLOSE_NOTIFICATION_ACTION,
        buttonTitle: "Close",
        options: { opensAppToForeground: false },
      },
    ]),
  ]);
}

export async function requestNotificationPermissions() {
  if (Platform.OS === "web" || !Device.isDevice) {
    return null;
  }

  const existingPermission = await Notifications.getPermissionsAsync();
  let finalStatus = existingPermission.status;

  if (existingPermission.status !== "granted") {
    const permission = await Notifications.requestPermissionsAsync();
    finalStatus = permission.status;
  }

  if (finalStatus !== "granted") {
    return null;
  }

  const devicePushToken = await Notifications.getDevicePushTokenAsync();
  return devicePushToken.data;
}

export async function getInstallationDeviceId() {
  if (Platform.OS === "web") throw new Error("Use browser push registration");
  const existing = await SecureStore.getItemAsync(INSTALLATION_ID_KEY);
  if (existing?.trim()) {
    return existing;
  }

  const installationId = createInstallationId();
  await SecureStore.setItemAsync(INSTALLATION_ID_KEY, installationId);
  return installationId;
}

function createInstallationId() {
  const randomUUID = globalThis.crypto?.randomUUID;
  if (typeof randomUUID === "function") {
    return randomUUID.call(globalThis.crypto);
  }
  return `install-${Date.now()}-${Math.random().toString(36).slice(2, 12)}`;
}

export function addNotificationListeners({
  onReceived,
  onResponse,
}: {
  onReceived?: (notification: Notifications.Notification) => void;
  onResponse?: (response: Notifications.NotificationResponse) => void;
}) {
  if (Platform.OS === "web") return () => {};
  const receivedSubscription = Notifications.addNotificationReceivedListener(
    (notification) => {
      onReceived?.(notification);
    },
  );
  const responseSubscription =
    Notifications.addNotificationResponseReceivedListener((response) => {
      onResponse?.(response);
    });

  return () => {
    receivedSubscription.remove();
    responseSubscription.remove();
  };
}
