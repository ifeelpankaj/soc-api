import { useMemo, useState } from "react";
import { StyleSheet, Text, View } from "react-native";
import { useRouter } from "expo-router";

import { Stack } from "@/components/layout";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { Button } from "@/components/ui";
import { HubFeedRow } from "@/features/hub/components/hub-feed-row";
import {
  HubChip,
  HubChipRow,
  HubEmptyState,
  HubPageHeader,
} from "@/features/hub/components/ui";
import { HubFeedSkeleton } from "@/features/hub/components/ui/hub-feed-skeleton";
import { hubCategoryName } from "@/features/hub/hub-format";
import {
  residentHubCreatePostRoute,
  residentHubPostRoute,
  residentHubRoute,
} from "@/features/hub/hub-routes";
import { hubTheme } from "@/features/hub/hub-theme";
import { useResident } from "@/features/resident/resident-context";
import {
  findHubChannel,
  useHubCategoriesQuery,
  useHubChannelsQuery,
  useHubPostsQuery,
} from "@/lib/api/hub-api";
import { spacing } from "@/theme/spacing";

type HubChannelFeedScreenProps = {
  channelType: "announcement" | "community";
  showBackHeader?: boolean;
};

export function HubChannelFeedScreen({
  channelType,
  showBackHeader = true,
}: HubChannelFeedScreenProps) {
  const router = useRouter();
  const { societyId } = useResident();
  const [categoryFilter, setCategoryFilter] = useState<number | undefined>();
  const [cursors, setCursors] = useState<(string | undefined)[]>([undefined]);
  const channelsQuery = useHubChannelsQuery(societyId ?? 0, {
    skip: !societyId,
    refetchOnMountOrArgChange: true,
  });
  const categoriesQuery = useHubCategoriesQuery(societyId ?? 0, {
    skip: !societyId || channelType !== "community",
  });
  const channel = findHubChannel(channelsQuery.data, channelType);
  const postsQuery = useHubPostsQuery(
    {
      societyId: societyId ?? 0,
      channelId: channel?.id ?? 0,
      cursor: cursors[cursors.length - 1],
      categoryId: channelType === "community" ? categoryFilter : undefined,
    },
    { skip: !societyId || !channel?.id, refetchOnMountOrArgChange: true },
  );

  const categories = categoriesQuery.data ?? [];

  const items = useMemo(() => {
    const page = postsQuery.data;
    if (!page) {
      return [];
    }
    return Array.from(
      new Map(
        [...(page.pinned ?? []), ...(page.items ?? [])].map((post) => [
          post.post_id ?? post.id,
          post,
        ]),
      ).values(),
    );
  }, [postsQuery.data]);

  const title = channelType === "announcement" ? "Announcements" : "Community";
  const subtitle =
    channelType === "announcement"
      ? "Official society communication"
      : "Posts from residents in your society";

  const openPost = (postId: number) => {
    if (!societyId || !channel?.id) {
      return;
    }
    router.push(
      residentHubPostRoute(
        societyId,
        postId,
        channel.id,
        channelType === "announcement" ? "announcements" : "community",
      ),
    );
  };

  const loading = channelsQuery.isFetching || postsQuery.isFetching;
  const error = channelsQuery.isError || postsQuery.isError;

  return (
    <Stack gap="lg" style={styles.screen}>
      {showBackHeader ? (
        <ScreenBackHeader fallbackHomeRoute={residentHubRoute()} title={title} />
      ) : null}
      <HubPageHeader subtitle={subtitle} />

      {channelType === "community" && categories.length > 0 ? (
        <HubChipRow>
          <HubChip
            label="All"
            selected={categoryFilter === undefined}
            onPress={() => {
              setCategoryFilter(undefined);
              setCursors([undefined]);
            }}
          />
          {categories.map((category) => (
            <HubChip
              key={category.id}
              label={category.name}
              selected={categoryFilter === category.id}
              onPress={() => {
                setCategoryFilter(category.id);
                setCursors([undefined]);
              }}
            />
          ))}
        </HubChipRow>
      ) : null}

      {channelType === "community" ? (
        <Button
          compact
          title="New post"
          variant="secondary"
          onPress={() => router.push(residentHubCreatePostRoute("community"))}
        />
      ) : null}

      {loading && items.length === 0 ? <HubFeedSkeleton count={4} /> : null}
      {error ? (
        <HubEmptyState
          actionLabel="Try again"
          message="Check your connection and try again."
          title={`Unable to load ${title.toLowerCase()}`}
          onAction={() => {
            void channelsQuery.refetch();
            void postsQuery.refetch();
          }}
        />
      ) : null}

      {!error && items.length === 0 && !loading ? (
        <HubEmptyState
          actionLabel={channelType === "community" ? "New post" : undefined}
          message={
            channelType === "announcement"
              ? "Society announcements will appear here."
              : "Be the first to share something with your neighbours."
          }
          title={
            categoryFilter
              ? "No posts in this category"
              : `No ${title.toLowerCase()} yet`
          }
          onAction={
            channelType === "community"
              ? () => router.push(residentHubCreatePostRoute("community"))
              : undefined
          }
        />
      ) : null}

      {items.map((post) => (
        <HubFeedRow
          key={post.post_id ?? post.id}
          categoryLabel={hubCategoryName(categories, post.category_id)}
          post={post}
          variant={channelType}
          onPress={() => openPost(post.post_id ?? post.id)}
        />
      ))}

      {cursors.length > 1 ? (
        <Button
          disabled={loading}
          title="Previous page"
          variant="ghost"
          onPress={() => setCursors((current) => current.slice(0, -1))}
        />
      ) : null}
      {postsQuery.data?.next_cursor ? (
        <Button
          disabled={loading}
          title="Load more"
          variant="ghost"
          onPress={() =>
            setCursors((current) => [...current, postsQuery.data?.next_cursor])
          }
        />
      ) : null}
    </Stack>
  );
}

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    paddingBottom: hubTheme.screenPaddingBottom,
  },
});
