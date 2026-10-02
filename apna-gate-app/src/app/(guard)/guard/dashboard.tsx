import { Redirect } from "expo-router";

import { guardHomeRoute } from "@/features/guard/guard-routes";

export default function LegacyGuardDashboardRoute() {
  return <Redirect href={guardHomeRoute()} />;
}
