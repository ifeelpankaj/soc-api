import { usePathname, useRouter, type Href } from "expo-router";
import type { ReactNode } from "react";

import { NotchedTabBar } from "@/components/layout/notched-tab-bar";
import { GuardTabIcon } from "@/features/guard/components/guard-tab-icons";
import { colors } from "@/theme/colors";

type TabBarState = {
  index: number;
  routes: { name: string }[];
};

type RoleTab = "home" | "profile";

type RoleTabBarProps = {
  fab: {
    accessibilityLabel: string;
    icon: ReactNode;
    label: string;
    onPress: () => void;
  };
  homeRoute: Href;
  profileRoute: Href;
  state: TabBarState;
};

function resolveActiveTab(pathname: string, state: TabBarState): RoleTab {
  if (pathname.includes("/profile")) {
    return "profile";
  }

  const routeName = state.routes[state.index]?.name;
  if (routeName === "profile") {
    return "profile";
  }

  return "home";
}

export function RoleTabBar({ fab, homeRoute, profileRoute, state }: RoleTabBarProps) {
  const router = useRouter();
  const pathname = usePathname();
  const activeTab = resolveActiveTab(pathname, state);

  const goToTab = (tab: RoleTab) => {
    if (tab === activeTab) {
      return;
    }

    router.replace(tab === "profile" ? profileRoute : homeRoute);
  };

  return (
    <NotchedTabBar
      fab={fab}
      leftTab={{
        active: activeTab === "home",
        icon: (
          <GuardTabIcon
            color={activeTab === "home" ? colors.brand.orange : colors.text.placeholder}
            name="dashboard"
          />
        ),
        label: "Home",
        onPress: () => goToTab("home"),
      }}
      rightTab={{
        active: activeTab === "profile",
        icon: (
          <GuardTabIcon
            color={activeTab === "profile" ? colors.brand.orange : colors.text.placeholder}
            name="profile"
          />
        ),
        label: "Profile",
        onPress: () => goToTab("profile"),
      }}
    />
  );
}
