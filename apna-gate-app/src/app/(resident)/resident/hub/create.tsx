import { HubCreatePostScreen } from "@/features/hub/components/hub-create-post-screen";
import { residentHubCommunityRoute } from "@/features/hub/hub-routes";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";

export default function HubCreatePostRoute() {
  return (
    <ResidentSubScreen fallbackHomeRoute={residentHubCommunityRoute()} title="Create post">
      <HubCreatePostScreen showBackHeader={false} />
    </ResidentSubScreen>
  );
}
