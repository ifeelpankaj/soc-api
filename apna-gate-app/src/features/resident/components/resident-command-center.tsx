import { useRouter } from "expo-router";
import { useCallback, useMemo, useState } from "react";
import { StyleSheet, View } from "react-native";

import { Stack } from "@/components/layout";
import {
  DashboardActionRow,
  DashboardActivityFeed,
  DashboardBannerCarousel,
  DashboardErrorBanner,
  DashboardHeader,
  DashboardHeroCard,
  DashboardOverviewGrid,
  DashboardSection,
  DashboardSkeleton,
  NeedsAttentionSection,
  SubscriptionExpiredBanner,
  getTimeGreeting,
  type DashboardActionTileConfig,
  type DashboardOverviewStatConfig,
} from "@/components/dashboard";
import { useResidentAttentionItems } from "@/features/dashboard/use-resident-attention-items";
import { residentHubPostRoute, residentHubRoute } from "@/features/hub/hub-routes";
import { notificationsRoute } from "@/features/notifications/notification-routing";
import { ResidentAnnouncementsSnapshot } from "@/features/resident/components/resident-announcements-snapshot";
import { ResidentMaintenanceDashboardCard } from "@/features/resident/components/resident-maintenance-dashboard-card";
import { ResidentScreenShell } from "@/features/resident/components/resident-screen-shell";
import { useResidentActivityFeed } from "@/features/resident/hooks/use-resident-activity-feed";
import { useResidentDashboard } from "@/features/resident/hooks/use-resident-dashboard";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { useResident } from "@/features/resident/resident-context";
import {
  residentEntriesRoute,
  residentInvitesRoute,
  residentMaintenanceRoute,
  residentMembersRoute,
  residentVisitorInviteRoute,
  residentVisitorsRoute,
} from "@/features/resident/resident-routes";
import { VisitorDetailSheet } from "@/features/visitors/components/visitor-detail-sheet";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import {
  findHubChannel,
  useHubChannelsQuery,
  useHubPostsQuery,
} from "@/lib/api/hub-api";
import { useGetV1MeNotificationsUnreadCountQuery } from "@/lib/api/notification-api-extensions";
import { useMaintenanceOutstandingQuery } from "@/lib/api/maintenance-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

export function ResidentCommandCenter() {
  const router = useRouter();
  const { societyId, flatId, user } = useResident();
  const feedback = useResidentFeedback();
  const activityFeed = useResidentActivityFeed();
  const dashboard = useResidentDashboard();
  const attentionItems = useResidentAttentionItems();
  const unreadNotificationsQuery = useGetV1MeNotificationsUnreadCountQuery();
  const maintenanceQuery = useMaintenanceOutstandingQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0 },
    { skip: !societyId || !flatId, refetchOnMountOrArgChange: true },
  );
  const channelsQuery = useHubChannelsQuery(societyId ?? 0, {
    skip: !societyId,
    refetchOnMountOrArgChange: true,
  });
  const announcementChannel = findHubChannel(channelsQuery.data, "announcement");
  const announcementPostsQuery = useHubPostsQuery(
    {
      societyId: societyId ?? 0,
      channelId: announcementChannel?.id ?? 0,
    },
    { skip: !societyId || !announcementChannel?.id, refetchOnMountOrArgChange: true },
  );

  const { refetchAll } = dashboard;
  const unreadNotifications =
    unreadNotificationsQuery.data?.data?.unread_count ?? 0;
  const [detailEntry, setDetailEntry] = useState<ModelsVisitorEntry | null>(
    null,
  );
  const [detailVisible, setDetailVisible] = useState(false);

  const goEntries = useCallback(
    (preset?: "expected" | "inside" | "all") =>
      router.push(residentEntriesRoute(preset ?? "expected")),
    [router],
  );
  const goNotifications = useCallback(
    () => router.push(notificationsRoute()),
    [router],
  );

  const handleRefresh = useCallback(() => {
    refetchAll();
    void activityFeed.refresh();
    void maintenanceQuery.refetch();
    void channelsQuery.refetch();
    void announcementPostsQuery.refetch();
  }, [
    activityFeed,
    announcementPostsQuery,
    channelsQuery,
    maintenanceQuery,
    refetchAll,
  ]);

  const openEntryDetail = useCallback((entry: ModelsVisitorEntry) => {
    setDetailEntry(entry);
    setDetailVisible(true);
  }, []);

  const goInvite = useCallback(() => {
    if (!dashboard.canManageFlatVisitors) {
      feedback.showInfo(
        "Permission required",
        "Only active flat residents with visitor access can invite guests.",
      );
      return;
    }
    router.push(residentVisitorInviteRoute());
  }, [dashboard.canManageFlatVisitors, feedback, router]);

  const handleStatPress = useCallback(
    (id: string) => {
      switch (id) {
        case "pending":
          router.push(residentVisitorsRoute());
          break;
        case "expected":
          router.push(residentEntriesRoute("expected"));
          break;
        case "visitors":
          router.push(residentEntriesRoute("all"));
          break;
        case "members":
          router.push(residentMembersRoute());
          break;
        default:
          break;
      }
    },
    [router],
  );

  const announcementPreview = useMemo(() => {
    const page = announcementPostsQuery.data;
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
  }, [announcementPostsQuery.data]);

  const actions = useMemo(() => {
    const tiles: DashboardActionTileConfig[] = [
      {
        id: "invites",
        title: "Manage Invites",
        subtitle: "Guest and member links",
        tone: "blue",
        icon: {
          ios: "person.crop.circle.badge.plus",
          android: "person_add",
          web: "person_add",
        },
        onPress: () => router.push(residentInvitesRoute()),
      },
      {
        id: "maintenance",
        title: "Maintenance",
        subtitle: "Bills and payments",
        tone: "orange",
        icon: { ios: "doc.text", android: "description", web: "description" },
        onPress: () => router.navigate(residentMaintenanceRoute()),
      },
      {
        id: "hub",
        title: "Society Hub",
        subtitle: "News and community",
        tone: "purple",
        icon: { ios: "circle.grid.2x2", android: "hub", web: "hub" },
        onPress: () => router.navigate(residentHubRoute()),
      },
      {
        id: "entries",
        title: "History",
        subtitle: "Recent visitor activity",
        tone: "neutral",
        icon: { ios: "clock", android: "history", web: "history" },
        onPress: () => goEntries("all"),
      },
    ];

    return tiles;
  }, [goEntries, router]);

  const overviewStats = useMemo<DashboardOverviewStatConfig[]>(
    () => [
      {
        id: "pending",
        label: "Pending Approval",
        subtext: "Awaiting your action",
        value: dashboard.pendingCount,
        tone: "orange",
        icon: {
          ios: "hourglass",
          android: "hourglass_top",
          web: "hourglass_top",
        },
      },
      {
        id: "expected",
        label: "Expected Today",
        subtext: "Expected today",
        value: dashboard.expectedCount,
        tone: "blue",
        icon: {
          ios: "calendar",
          android: "calendar_today",
          web: "calendar_today",
        },
      },
      {
        id: "visitors",
        label: "Recent Visitors",
        subtext: "Recently visited",
        value: dashboard.visitorsCount,
        tone: "green",
        icon: { ios: "person.2.fill", android: "groups", web: "groups" },
      },
      {
        id: "members",
        label: "Flat Members",
        subtext: "In your flat",
        value: dashboard.membersCount,
        tone: "neutral",
        icon: { ios: "person.2.fill", android: "group", web: "group" },
      },
    ],
    [
      dashboard.expectedCount,
      dashboard.membersCount,
      dashboard.pendingCount,
      dashboard.visitorsCount,
    ],
  );

  const openAnnouncementPost = useCallback(
    (postId: number) => {
      if (!societyId) {
        return;
      }
      router.push(
        residentHubPostRoute(societyId, postId, announcementChannel?.id, "home"),
      );
    },
    [announcementChannel?.id, router, societyId],
  );

  return (
    <View style={styles.screen}>
      <ResidentScreenShell
        backgroundColor={colors.guard.screenBg}
        contentPaddingBottom={layout.tabBarHeight + spacing.lg}
        onRefresh={handleRefresh}
        refreshing={
          dashboard.isRefreshing ||
          activityFeed.isRefreshing ||
          maintenanceQuery.isFetching
        }
      >
        {dashboard.isInitialLoading ? (
          <DashboardSkeleton sections={["header", "hero", "actions", "stats", "attention"]} />
        ) : (
          <Stack gap="2xl">
            <Stack gap="lg">
              <DashboardHeader
                actions={[
                  {
                    accessibilityLabel:
                      unreadNotifications > 0
                        ? `${unreadNotifications} unread notifications`
                        : "Notifications",
                    icon: {
                      ios: "bell",
                      android: "notifications_none",
                      web: "notifications_none",
                    },
                    notificationCount:
                      unreadNotifications > 0 ? unreadNotifications : undefined,
                    onPress: goNotifications,
                  },
                ]}
                greeting={getTimeGreeting(
                  user?.full_name ?? dashboard.displayName,
                )}
                showBrand
                statusItems={[
                  {
                    label: dashboard.isSubscriptionBlocked ? "Limited" : "Live",
                    live:
                      !dashboard.hasError && !dashboard.isSubscriptionBlocked,
                  },
                ]}
                title={dashboard.flatLabel}
              />
              <DashboardBannerCarousel />
            </Stack>

            {dashboard.isSubscriptionBlocked ? (
              <SubscriptionExpiredBanner />
            ) : null}

            {dashboard.hasError ? (
              <DashboardErrorBanner
                message={dashboard.errorMessage ?? "Unable to load data."}
                onRetry={handleRefresh}
              />
            ) : null}

            {!dashboard.isSubscriptionBlocked ? (
              <>
                <NeedsAttentionSection items={attentionItems} />

                <Stack gap="md">
                  <DashboardHeroCard
                    icon={{
                      ios: "person.badge.plus",
                      android: "person_add",
                      web: "person_add",
                    }}
                    subtitle="Pre-approve a guest for entry"
                    title="Invite Visitor"
                    onPress={goInvite}
                  />
                  <DashboardActionRow actions={actions} columns={2} />
                </Stack>

                <ResidentMaintenanceDashboardCard
                  bill={maintenanceQuery.data?.current_bill}
                  loading={maintenanceQuery.isLoading && !maintenanceQuery.data}
                />

                <ResidentAnnouncementsSnapshot
                  items={announcementPreview}
                  loading={
                    announcementPostsQuery.isLoading && !announcementPostsQuery.data
                  }
                  onItemPress={openAnnouncementPost}
                />

                <DashboardSection title="Overview">
                  <DashboardOverviewGrid
                    stats={overviewStats}
                    onStatPress={handleStatPress}
                  />
                </DashboardSection>

                <DashboardActivityFeed
                  emptyAction={{
                    label: "Invite Visitor",
                    onPress: goInvite,
                  }}
                  hasMore={activityFeed.hasMore}
                  isLoading={activityFeed.isLoading}
                  isLoadingMore={activityFeed.isLoadingMore}
                  items={activityFeed.items}
                  onItemPress={openEntryDetail}
                  onLoadMore={() => {
                    void activityFeed.loadMore();
                  }}
                  onViewAll={() => goEntries()}
                />
              </>
            ) : null}
          </Stack>
        )}
      </ResidentScreenShell>

      <VisitorDetailSheet
        variant="resident"
        entry={detailEntry}
        loading={dashboard.isActionLoading}
        loadingAction={dashboard.loadingAction}
        primaryActionLabel={
          detailEntry?.status === "waiting_approval" &&
          dashboard.canManageFlatVisitors
            ? "Approve"
            : undefined
        }
        secondaryActionLabel={
          detailEntry?.status === "waiting_approval" &&
          dashboard.canManageFlatVisitors
            ? "Reject"
            : undefined
        }
        visible={detailVisible}
        onClose={() => {
          setDetailVisible(false);
          setDetailEntry(null);
        }}
        onPrimaryAction={
          detailEntry?.status === "waiting_approval"
            ? async () => {
                if (await dashboard.handleApprove(detailEntry.id)) {
                  setDetailVisible(false);
                  setDetailEntry(null);
                }
              }
            : undefined
        }
        onSecondaryAction={
          detailEntry?.status === "waiting_approval"
            ? async () => {
                if (await dashboard.handleReject(detailEntry.id)) {
                  setDetailVisible(false);
                  setDetailEntry(null);
                }
              }
            : undefined
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
});
