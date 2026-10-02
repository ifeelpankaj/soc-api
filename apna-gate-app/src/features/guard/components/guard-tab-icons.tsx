import { SymbolView } from "expo-symbols";
import type { ColorValue } from "react-native";

import { theme } from "@/theme";

type GuardTabName = "dashboard" | "add-entry" | "scan" | "logs" | "profile";

const TAB_ICONS = {
  dashboard: {
    ios: "house.fill",
    android: "home",
    web: "home",
  },
  "add-entry": {
    ios: "person.badge.plus",
    android: "person_add",
    web: "person_add",
  },
  scan: {
    ios: "qrcode.viewfinder",
    android: "qr_code_scanner",
    web: "qr_code_scanner",
  },
  logs: {
    ios: "list.bullet.rectangle",
    android: "list_alt",
    web: "list_alt",
  },
  profile: {
    ios: "person.crop.circle",
    android: "account_circle",
    web: "account_circle",
  },
} as const;

export function GuardTabIcon({
  color,
  name,
}: {
  color: ColorValue;
  name: GuardTabName;
}) {
  return (
    <SymbolView
      name={TAB_ICONS[name]}
      size={22}
      tintColor={color}
      style={{ width: 24, height: 24 }}
    />
  );
}

export const guardTabBarOptions = {
  headerShown: false,
  tabBarActiveTintColor: theme.brand.orange,
  tabBarInactiveTintColor: theme.text.placeholder,
  tabBarLabelStyle: {
    fontSize: 11,
    fontWeight: "600" as const,
  },
  tabBarStyle: {
    backgroundColor: theme.surface.screen,
    borderTopColor: theme.border.default,
    height: 68,
    paddingBottom: 10,
    paddingTop: 8,
  },
};
