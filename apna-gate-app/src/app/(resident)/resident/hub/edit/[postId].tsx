import { useLocalSearchParams } from "expo-router";

import { HubCreatePostScreen } from "@/features/hub/components/hub-create-post-screen";
import { residentHubCommunityRoute } from "@/features/hub/hub-routes";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";

export default function HubEditPostRoute() {
  const params = useLocalSearchParams<{ postId?: string }>();
  const postId = Number(params.postId);

  return (
    <ResidentSubScreen fallbackHomeRoute={residentHubCommunityRoute()} title="Edit post">
      <HubCreatePostScreen editPostId={postId} showBackHeader={false} />
    </ResidentSubScreen>
  );
}
