import { Redirect } from "expo-router";

import { residentEntriesRoute } from "@/features/resident/resident-routes";

export default function ResidentLogsRedirect() {
  return <Redirect href={residentEntriesRoute()} />;
}
