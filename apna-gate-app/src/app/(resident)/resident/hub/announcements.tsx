import { HubChannelFeedScreen } from "@/features/hub/components/hub-channel-feed-screen";
import { residentHubRoute } from "@/features/hub/hub-routes";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";

export default function HubAnnouncementsScreen() {
  return (
    <ResidentSubScreen fallbackHomeRoute={residentHubRoute()} title="Announcements">
      <HubChannelFeedScreen channelType="announcement" showBackHeader={false} />
    </ResidentSubScreen>
  );
}
