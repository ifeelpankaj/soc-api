import { useRouter } from "expo-router";

import { ScanQrIcon } from "@/components/icons";
import { RoleTabBar } from "@/components/layout/role-tab-bar";
import {
  guardHomeRoute,
  guardProfileRoute,
  guardScannerRoute,
} from "@/features/guard/guard-routes";
import { colors } from "@/theme/colors";

type TabBarState = {
  index: number;
  routes: { name: string }[];
};

export function GuardTabBar({ state }: { state: TabBarState }) {
  const router = useRouter();

  return (
    <RoleTabBar
      fab={{
        accessibilityLabel: "Scan QR",
        icon: <ScanQrIcon color={colors.brand.orangeSoft ?? "#1B1F3B"} size={22} />,
        label: "Scan",
        onPress: () => router.push(guardScannerRoute()),
      }}
      homeRoute={guardHomeRoute()}
      profileRoute={guardProfileRoute()}
      state={state}
    />
  );
}
