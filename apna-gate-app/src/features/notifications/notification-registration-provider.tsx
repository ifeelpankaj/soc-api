import type { PropsWithChildren } from "react";

import { useDeviceTokenRegistration } from "@/features/notifications/use-device-token-registration";

export function NotificationRegistrationProvider({ children }: PropsWithChildren) {
  useDeviceTokenRegistration();
  return children;
}
