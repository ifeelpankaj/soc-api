import { Redirect } from "expo-router";

import { residentHubAnnouncementsRoute } from "@/features/hub/hub-routes";

export default function LegacyAnnouncementsRedirect() {
  return <Redirect href={residentHubAnnouncementsRoute()} />;
}
