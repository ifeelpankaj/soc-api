import { useRouter } from "expo-router";
import { SymbolView } from "expo-symbols";

import { NotchedTabBar } from "@/components/layout/notched-tab-bar";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { useResident } from "@/features/resident/resident-context";
import { residentHubRoute } from "@/features/hub/hub-routes";
import {
  residentDashboardRoute,
  residentProfileRoute,
  residentMaintenanceRoute,
  residentVisitorInviteRoute,
} from "@/features/resident/resident-routes";
import { colors } from "@/theme/colors";

type TabBarState = {
  index: number;
  routes: { name: string }[];
};

export function ResidentTabBar({ state }: { state: TabBarState }) {
  const router = useRouter();
  const feedback = useAppFeedback();
  const { canManageFlatVisitors } = useResident();

  const goToInvite = () => {
    if (!canManageFlatVisitors) {
      feedback.showInfo(
        "Permission required",
        "Only active flat residents with visitor access can invite guests.",
      );
      return;
    }

    router.push(residentVisitorInviteRoute());
  };

  return (
    <NotchedTabBar
      fab={{
        accessibilityLabel: "Invite visitor",
        icon: (
          <SymbolView
            name={{
              ios: "person.badge.plus",
              android: "person_add",
              web: "person_add",
            }}
            size={22}
            tintColor={colors.brand.orangeSoft ?? "#1B1F3B"}
          />
        ),
        label: "Invite",
        onPress: goToInvite,
      }}
      leftTab={{
        label: "Home",
        active: state.routes[state.index]?.name === "dashboard",
        onPress: () => router.navigate(residentDashboardRoute()),
        icon: (
          <SymbolView
            name={{ ios: "house", android: "home", web: "home" }}
            size={22}
            tintColor={
              state.routes[state.index]?.name === "dashboard"
                ? colors.brand.orange
                : colors.text.placeholder
            }
          />
        ),
      }}
      additionalLeftTabs={[
        {
          label: "Maintenance",
          active: state.routes[state.index]?.name === "maintenance",
          onPress: () => router.navigate(residentMaintenanceRoute()),
          icon: (
            <SymbolView
              name={{
                ios: "doc.text",
                android: "description",
                web: "description",
              }}
              size={22}
              tintColor={
                state.routes[state.index]?.name === "maintenance"
                  ? colors.brand.orange
                  : colors.text.placeholder
              }
            />
          ),
        },
      ]}
      rightTab={{
        label: "Profile",
        active: state.routes[state.index]?.name === "profile",
        onPress: () => router.navigate(residentProfileRoute()),
        icon: (
          <SymbolView
            name={{
              ios: "person.crop.circle",
              android: "account_circle",
              web: "account_circle",
            }}
            size={22}
            tintColor={
              state.routes[state.index]?.name === "profile"
                ? colors.brand.orange
                : colors.text.placeholder
            }
          />
        ),
      }}
      additionalRightTabs={[
        {
          label: "Hub",
          active: state.routes[state.index]?.name === "hub",
          onPress: () => router.navigate(residentHubRoute()),
          icon: (
            <SymbolView
              name={{ ios: "circle.grid.2x2", android: "hub", web: "hub" }}
              size={22}
              tintColor={
                state.routes[state.index]?.name === "hub"
                  ? colors.brand.orange
                  : colors.text.placeholder
              }
            />
          ),
        },
      ]}
    />
  );
}
