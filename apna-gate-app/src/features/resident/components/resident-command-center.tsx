import { useProfilePhoto } from "@/features/profile/use-profile-photo";
import { useRouter } from "expo-router";
import { useCallback, useMemo, useState } from "react";
import { StyleSheet, View } from "react-native";

import { Stack } from "@/components/layout";
import {
  DashboardActionRow,
  DashboardActivityFeed,
  DashboardAlertBar,
  DashboardBannerCarousel,
  DashboardErrorBanner,
  DashboardHeader,
  DashboardHeroCard,
  DashboardOverviewGrid,
  DashboardSection,
  DashboardSkeleton,
  SubscriptionExpiredBanner,
  getTimeGreeting,
  type DashboardActionTileConfig,
  type DashboardOverviewStatConfig,
} from "@/components/dashboard";
import { VisitorDetailSheet } from "@/features/visitors/components/visitor-detail-sheet";
import { ResidentScreenShell } from "@/features/resident/components/resident-screen-shell";
import { useResidentActivityFeed } from "@/features/resident/hooks/use-resident-activity-feed";
import { useResidentDashboard } from "@/features/resident/hooks/use-resident-dashboard";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { notificationsRoute } from "@/features/notifications/notification-routing";
import { useResident } from "@/features/resident/resident-context";
import {
  residentEntriesRoute,
  residentInvitesRoute,
  residentAnnouncementsRoute,
  residentMaintenanceRoute,
  residentMembersRoute,
  residentProfileRoute,
  residentVisitorInviteRoute,
  residentVisitorsRoute,
} from "@/features/resident/resident-routes";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { useGetV1MeNotificationsUnreadCountQuery } from "@/lib/api/notification-api-extensions";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

export function ResidentCommandCenter() {
  const router = useRouter();
  const { user } = useResident();
  const avatarUrl = useProfilePhoto(user?.avatar_url);
  const feedback = useResidentFeedback();
  const activityFeed = useResidentActivityFeed();
  const dashboard = useResidentDashboard();
  const unreadNotificationsQuery = useGetV1MeNotificationsUnreadCountQuery();
  const { refetchAll } = dashboard;
  const unreadNotifications =
    unreadNotificationsQuery.data?.data?.unread_count ?? 0;
  const [detailEntry, setDetailEntry] = useState<ModelsVisitorEntry | null>(
    null,
  );
  const [detailVisible, setDetailVisible] = useState(false);

  const goApprovals = useCallback(
    () => router.push(residentVisitorsRoute()),
    [router],
  );
  const goEntries = useCallback(
    (preset?: "expected" | "inside" | "all") =>
      router.push(residentEntriesRoute(preset ?? "expected")),
    [router],
  );
  const goProfile = useCallback(
    () => router.push(residentProfileRoute()),
    [router],
  );
  const goNotifications = useCallback(
    () => router.push(notificationsRoute()),
    [router],
  );

  const handleRefresh = useCallback(() => {
    refetchAll();
    void activityFeed.refresh();
  }, [activityFeed, refetchAll]);

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
        id: "announcements",
        title: "Announcements",
        subtitle: "Society news and notices",
        tone: "purple",
        icon: { ios: "megaphone", android: "campaign", web: "campaign" },
        onPress: () => router.navigate(residentAnnouncementsRoute()),
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
        id: "entries",
        title: "Visit History",
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

  return (
    <View style={styles.screen}>
      <ResidentScreenShell
        backgroundColor={colors.guard.screenBg}
        contentPaddingBottom={layout.tabBarHeight + spacing.lg}
        onRefresh={handleRefresh}
        refreshing={dashboard.isRefreshing || activityFeed.isRefreshing}
      >
        {dashboard.isInitialLoading ? (
          <DashboardSkeleton />
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
                profileAvatar={{
                  imageUrl: avatarUrl,
                  name: user?.full_name ?? dashboard.displayName,
                  onPress: goProfile,
                  showOnlineDot: true,
                }}
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

            {dashboard.pendingCount > 0 && !dashboard.isSubscriptionBlocked ? (
              <DashboardAlertBar
                count={dashboard.pendingCount}
                message={`${dashboard.pendingCount} visitor${dashboard.pendingCount === 1 ? "" : "s"} awaiting approval`}
                onPress={goApprovals}
              />
            ) : null}

            {!dashboard.isSubscriptionBlocked ? (
              <>
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

                <DashboardSection
                  actionLabel="View details >"
                  title="Overview"
                  onAction={() => goEntries()}
                >
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
