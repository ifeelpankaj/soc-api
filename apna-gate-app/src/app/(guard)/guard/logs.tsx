import { Redirect } from "expo-router";

import { guardEntriesRoute } from "@/features/guard/guard-routes";

export default function LegacyGuardLogsRoute() {
  return <Redirect href={guardEntriesRoute("today")} />;
}
