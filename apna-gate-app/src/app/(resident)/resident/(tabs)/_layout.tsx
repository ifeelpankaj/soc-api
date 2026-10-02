import { Tabs } from "expo-router";

import { ResidentTabBar } from "@/features/resident/components/resident-tab-bar";

export default function ResidentTabsLayout() {
  return (
    <Tabs
      screenOptions={{ headerShown: false }}
      tabBar={(props) => <ResidentTabBar {...props} />}
    >
      <Tabs.Screen name="dashboard" options={{ title: "Home" }} />
      <Tabs.Screen name="maintenance" options={{ title: "Maintenance" }} />
      <Tabs.Screen name="announcements" options={{ title: "Announcements" }} />
      <Tabs.Screen name="profile" options={{ title: "Profile" }} />
    </Tabs>
  );
}
