import { useLocalSearchParams } from "expo-router";

import { HubPostScreen } from "@/features/hub/components/hub-post-screen";
import { resolveHubPostBackRoute } from "@/features/hub/hub-routes";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";

export default function HubPostRoute() {
  const params = useLocalSearchParams<{ returnTo?: string }>();
  const backRoute = resolveHubPostBackRoute(params.returnTo);

  return (
    <ResidentSubScreen fallbackHomeRoute={backRoute} title="Post">
      <HubPostScreen showBackHeader={false} />
    </ResidentSubScreen>
  );
}
