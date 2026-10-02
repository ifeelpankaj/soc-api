import { useMemo } from "react";
import { StyleSheet, Text, View } from "react-native";
import { useRouter } from "expo-router";

import { Stack } from "@/components/layout";
import { ResidentScreenShell } from "@/features/resident/components/resident-screen-shell";
import { HubDestinationCard } from "@/features/hub/components/hub-destination-card";
import { HubFeedRow } from "@/features/hub/components/hub-feed-row";
import { HubEmptyState } from "@/features/hub/components/ui";
import { HubFeedSkeleton } from "@/features/hub/components/ui/hub-feed-skeleton";
import { hubCategoryName } from "@/features/hub/hub-format";
import {
  residentHubAnnouncementsRoute,
  residentHubCommunityRoute,
  residentHubPostRoute,
} from "@/features/hub/hub-routes";
import { hubTheme } from "@/features/hub/hub-theme";
import { useResident } from "@/features/resident/resident-context";
import {
  findHubChannel,
  useHubCategoriesQuery,
  useHubChannelsQuery,
  useHubPostsQuery,
  type HubContent,
} from "@/lib/api/hub-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

function mergePosts(pinned: HubContent[] = [], items: HubContent[] = []) {
  return Array.from(
    new Map([...pinned, ...items].map((post) => [post.post_id ?? post.id, post])).values(),
  );
}

export function HubHomeScreen() {
  const router = useRouter();
  const { societyId } = useResident();
  const channelsQuery = useHubChannelsQuery(societyId ?? 0, {
    skip: !societyId,
    refetchOnMountOrArgChange: true,
  });
  const categoriesQuery = useHubCategoriesQuery(societyId ?? 0, { skip: !societyId });
  const announcementChannel = findHubChannel(channelsQuery.data, "announcement");
  const communityChannel = findHubChannel(channelsQuery.data, "community");

  const announcementPosts = useHubPostsQuery(
    { societyId: societyId ?? 0, channelId: announcementChannel?.id ?? 0, cursor: undefined },
    { skip: !societyId || !announcementChannel?.id, refetchOnMountOrArgChange: true },
  );
  const communityPosts = useHubPostsQuery(
    { societyId: societyId ?? 0, channelId: communityChannel?.id ?? 0, cursor: undefined },
    { skip: !societyId || !communityChannel?.id, refetchOnMountOrArgChange: true },
  );

  const categories = categoriesQuery.data ?? [];

  const recentItems = useMemo(() => {
    const announcement = mergePosts(
      announcementPosts.data?.pinned,
      announcementPosts.data?.items,
    ).slice(0, 2);
    const community = mergePosts(
      communityPosts.data?.pinned,
      communityPosts.data?.items,
    ).slice(0, 2);
    return [
      ...announcement.map((post) => ({ post, variant: "announcement" as const })),
      ...community.map((post) => ({ post, variant: "community" as const })),
    ]
      .sort(
        (a, b) =>
          Date.parse(b.post.created_at) - Date.parse(a.post.created_at),
      )
      .slice(0, 3);
  }, [announcementPosts.data, communityPosts.data]);

  const loading =
    channelsQuery.isLoading ||
    announcementPosts.isLoading ||
    communityPosts.isLoading;
  const error = channelsQuery.isError;

  const openPost = (postId: number, channelId?: number) => {
    if (!societyId) {
      return;
    }
    router.push(residentHubPostRoute(societyId, postId, channelId, "hub"));
  };

  return (
    <ResidentScreenShell
      backgroundColor={colors.guard.screenBg}
      contentPaddingBottom={layout.tabBarHeight + spacing.lg}
      onRefresh={() => {
        void channelsQuery.refetch();
        void announcementPosts.refetch();
        void communityPosts.refetch();
      }}
      refreshing={channelsQuery.isFetching}
    >
      <Stack gap="xl">
        <View>
          <Text style={styles.heading}>Society Hub</Text>
          <Text style={styles.subheading}>Official updates and community conversations</Text>
        </View>

        {loading ? <HubFeedSkeleton count={2} /> : null}
        {error ? (
          <HubEmptyState
            actionLabel="Try again"
            message="Check your connection and try again."
            title="Unable to load Society Hub"
            onAction={() => void channelsQuery.refetch()}
          />
        ) : null}

        {!loading && !error ? (
          <>
            <Stack gap="md">
              <HubDestinationCard
                description="Official updates from your society"
                icon={{ ios: "megaphone.fill", android: "campaign", web: "campaign" }}
                iconBackground={colors.dashboard.actionPurpleSoft}
                iconColor={colors.dashboard.actionPurple}
                title="Announcements"
                unreadCount={announcementChannel?.unread_count ?? 0}
                onPress={() => router.push(residentHubAnnouncementsRoute())}
              />
              <HubDestinationCard
                description="Connect with your neighbours"
                icon={{ ios: "person.2.fill", android: "groups", web: "groups" }}
                iconBackground={colors.dashboard.actionBlueSoft}
                iconColor={colors.dashboard.actionBlue}
                title="Community"
                unreadCount={communityChannel?.unread_count ?? 0}
                onPress={() => router.push(residentHubCommunityRoute())}
              />
            </Stack>

            {recentItems.length > 0 ? (
              <Stack gap="md">
                <Text style={styles.sectionLabel}>Recent</Text>
                {recentItems.map(({ post, variant }) => {
                  const postId = post.post_id ?? post.id;
                  const channelId =
                    variant === "announcement"
                      ? announcementChannel?.id
                      : communityChannel?.id;
                  return (
                    <HubFeedRow
                      key={`${variant}-${postId}`}
                      categoryLabel={hubCategoryName(categories, post.category_id)}
                      post={post}
                      variant={variant}
                      onPress={() => openPost(postId, channelId)}
                    />
                  );
                })}
              </Stack>
            ) : null}
          </>
        ) : null}
      </Stack>
    </ResidentScreenShell>
  );
}

const styles = StyleSheet.create({
  heading: {
    ...hubTheme.typography.pageTitle,
  },
  sectionLabel: {
    ...hubTheme.typography.sectionEyebrow,
  },
  subheading: {
    ...hubTheme.typography.pageSubtitle,
    marginTop: spacing.xs,
  },
});
