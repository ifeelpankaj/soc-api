import { useEffect, useRef } from "react";
import { AppState, Platform, type AppStateStatus } from "react-native";

import {
  usePostV1MeDeviceTokensMutation,
} from "@/lib/api/generated-api";
import {
  configureNotifications,
  getInstallationDeviceId,
  getDevicePlatform,
  requestNotificationPermissions,
} from "@/lib/notifications/register-notifications";
import { setDevicePushToken } from "@/redux/notificationSlice";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";

import { useAuth } from "../auth/use-auth";

async function registerDeviceToken(
  registerToken: ReturnType<typeof usePostV1MeDeviceTokensMutation>[0],
  dispatch: ReturnType<typeof useAppDispatch>,
) {
  await configureNotifications();
  const token = await requestNotificationPermissions();
  if (!token) {
    dispatch(setDevicePushToken(null));
    return;
  }
  const deviceId = await getInstallationDeviceId();

  await registerToken({
    modelsRegisterDeviceTokenRequest: {
      token,
      platform: getDevicePlatform(),
      device_id: deviceId,
    },
  }).unwrap();

  dispatch(setDevicePushToken(token));
}

export function useDeviceTokenRegistration() {
  const dispatch = useAppDispatch();
  const { status } = useAuth();
  const storedToken = useAppSelector((state) => state.notifications.devicePushToken);
  const [registerToken] = usePostV1MeDeviceTokensMutation();
  const appState = useRef(AppState.currentState);

  useEffect(() => {
    if (Platform.OS === "web" || status !== "authenticated") {
      return;
    }

    void registerDeviceToken(registerToken, dispatch).catch(() => {
      console.error("Failed to register push notification token");
      dispatch(setDevicePushToken(null));
    });
  }, [dispatch, registerToken, status]);

  useEffect(() => {
    if (Platform.OS === "web" || status !== "authenticated") {
      return;
    }

    const subscription = AppState.addEventListener("change", (nextState: AppStateStatus) => {
      if (appState.current.match(/inactive|background/) && nextState === "active") {
        void registerDeviceToken(registerToken, dispatch).catch(() => {
          console.error("Failed to refresh push notification token");
        });
      }
      appState.current = nextState;
    });

    return () => {
      subscription.remove();
    };
  }, [dispatch, registerToken, status]);

  return {
    devicePushToken: storedToken,
  };
}
