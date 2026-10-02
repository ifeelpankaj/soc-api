import { useProfilePhoto } from "@/features/profile/use-profile-photo";
import { useRouter } from "expo-router";
import { useCallback, useMemo } from "react";
import { StyleSheet, View } from "react-native";

import { Stack } from "@/components/layout";
import {
  DashboardActionRow,
  DashboardActivityFeed,
  NeedsAttentionSection,
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
import { GuardScreenShell } from "@/features/guard/components/guard-screen-shell";
import { useGuardSociety } from "@/features/guard/guard-context";
import { useGuardCheckInReadiness } from "@/features/guard/check-in-readiness-provider";
import { mapGuardAttentionItems } from "@/features/dashboard/map-guard-attention-items";
import { useGuardActivityFeed } from "@/features/guard/hooks/use-guard-activity-feed";
import { useGuardDashboard } from "@/features/guard/hooks/use-guard-dashboard";
import { notificationsRoute } from "@/features/notifications/notification-routing";
import {
  guardAddEntryRoute,
  guardEntriesRoute,
  guardEntryDetailRoute,
  guardPendingRoute,
  guardProfileRoute,
  guardScannerRoute,
  guardWaitingAtGateRoute,
} from "@/features/guard/guard-routes";
import { useGetV1MeNotificationsUnreadCountQuery } from "@/lib/api/notification-api-extensions";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

export function GuardCommandCenter() {
  const router = useRouter();
  const { user } = useGuardSociety();
  const { readyCount } = useGuardCheckInReadiness();
  const avatarUrl = useProfilePhoto(user?.avatar_url);
  const activityFeed = useGuardActivityFeed();
  const {
    errorMessage,
    expectedGuestsCount,
    hasError,
    isInitialLoading,
    isRefreshing,
    isSubscriptionBlocked,
    refetchAll,
    societyName,
    stats,
  } = useGuardDashboard();
  const unreadNotificationsQuery = useGetV1MeNotificationsUnreadCountQuery();

  const pendingCount = stats?.pending_approvals ?? 0;
  const insideCount = stats?.visitors_inside ?? 0;
  const checkedOutCount = stats?.checked_out_today ?? 0;
  const unreadNotifications = unreadNotificationsQuery.data?.data?.unread_count ?? 0;

  const goPending = () => router.push(guardPendingRoute());
  const goLogs = () => router.push(guardEntriesRoute("today"));
  const goProfile = () => router.push(guardProfileRoute());
  const goNotifications = () => router.push(notificationsRoute());

  const openEntryDetail = useCallback(
    (entry: { id?: number }) => {
      if (!entry.id) {
        return;
      }

      router.push(guardEntryDetailRoute(entry.id));
    },
    [router],
  );

  const handleRefresh = useCallback(() => {
    refetchAll();
    void activityFeed.refresh();
  }, [activityFeed, refetchAll]);

  const guardAttentionItems = useMemo(
    () =>
      mapGuardAttentionItems({
        pendingCount,
        readyCount,
        onReviewPending: goPending,
        onCheckIn: () => router.push(guardWaitingAtGateRoute()),
      }),
    [goPending, pendingCount, readyCount, router],
  );

  const handleStatPress = useCallback(
    (id: string) => {
      switch (id) {
        case "inside":
          router.push(guardEntriesRoute("inside"));
          break;
        case "pending":
          router.push(guardPendingRoute());
          break;
        case "expected-guests":
          router.push(guardEntriesRoute("expected"));
          break;
        case "checked-out":
          router.push(guardEntriesRoute("checked_out"));
          break;
        default:
          break;
      }
    },
    [router],
  );

  const actions = useMemo<DashboardActionTileConfig[]>(
    () => [
      {
        id: "entry",
        title: "Add Visitor",
        subtitle: "Register new visitor",
        tone: "orange",
        icon: {
          ios: "person.badge.plus",
          android: "person_add",
          web: "person_add",
        },
        onPress: () => router.push(guardAddEntryRoute()),
      },
      {
        id: "waiting-at-gate",
        title: "Check In",
        subtitle: readyCount > 0 ? `${readyCount} ready to check in` : "Check-in visitor",
        badgeCount: readyCount,
        tone: "blue",
        icon: {
          ios: "door.left.hand.open",
          android: "meeting_room",
          web: "meeting_room",
        },
        onPress: () => router.push(guardWaitingAtGateRoute()),
      },
      {
        id: "logs",
        title: "View History",
        subtitle: "Recent visitor activity",
        tone: "purple",
        icon: {
          ios: "clock",
          android: "history",
          web: "history",
        },
        onPress: () => router.push(guardEntriesRoute("today")),
      },
    ],
    [router, readyCount],
  );

  const overviewStats = useMemo<DashboardOverviewStatConfig[]>(
    () => [
      {
        id: "inside",
        label: "Inside Visitors",
        subtext: "Currently inside",
        value: insideCount,
        tone: "green",
        icon: { ios: "person.2.fill", android: "groups", web: "groups" },
      },
      {
        id: "pending",
        label: "Awaiting Approval",
        subtext: "Pending approval",
        value: pendingCount,
        tone: "orange",
        icon: {
          ios: "hourglass",
          android: "hourglass_top",
          web: "hourglass_top",
        },
      },
      {
        id: "expected-guests",
        label: "Expected Guests",
        subtext: "Expected today",
        value: expectedGuestsCount,
        tone: "blue",
        icon: {
          ios: "calendar.badge.clock",
          android: "event_available",
          web: "event_available",
        },
      },
      {
        id: "checked-out",
        label: "Checked Out",
        subtext: "Today",
        value: checkedOutCount,
        tone: "neutral",
        icon: { ios: "arrow.right.square", android: "logout", web: "logout" },
      },
    ],
    [checkedOutCount, expectedGuestsCount, insideCount, pendingCount],
  );

  return (
    <View style={styles.screen}>
      <GuardScreenShell
        backgroundColor={colors.guard.screenBg}
        contentPaddingBottom={layout.tabBarHeight + spacing.lg}
        onRefresh={handleRefresh}
        refreshing={isRefreshing || activityFeed.isRefreshing}
      >
        {isInitialLoading ? (
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
                greeting={getTimeGreeting(user?.full_name, "Guard")}
                profileAvatar={{
                  imageUrl: avatarUrl,
                  name: user?.full_name,
                  onPress: goProfile,
                  showOnlineDot: true,
                }}
                statusItems={[
                  {
                    label: isSubscriptionBlocked ? "Limited" : "On duty",
                    live: !hasError && !isSubscriptionBlocked,
                  },
                ]}
                title={societyName}
              />
              <DashboardBannerCarousel />
            </Stack>

            {isSubscriptionBlocked ? <SubscriptionExpiredBanner /> : null}

            {hasError ? (
              <DashboardErrorBanner
                message={errorMessage ?? "Unable to load data."}
                onRetry={handleRefresh}
              />
            ) : null}

            {!isSubscriptionBlocked ? (
              <NeedsAttentionSection items={guardAttentionItems} />
            ) : null}

            {!isSubscriptionBlocked ? (
              <>
                <Stack gap="md">
                  <DashboardHeroCard
                    icon={{
                      ios: "qrcode.viewfinder",
                      android: "qr_code_scanner",
                      web: "qr_code_scanner",
                    }}
                    subtitle="Scan QR code to verify entry"
                    title="Scan Visitor"
                    onPress={() => router.push(guardScannerRoute())}
                  />
                  <DashboardActionRow actions={actions} />
                </Stack>

                <DashboardSection title="Visitor Overview">
                  <DashboardOverviewGrid
                    stats={overviewStats}
                    onStatPress={handleStatPress}
                  />
                </DashboardSection>

                <DashboardActivityFeed
                  emptyAction={{
                    label: "Add Visitor",
                    onPress: () => router.push(guardAddEntryRoute()),
                  }}
                  hasMore={activityFeed.hasMore}
                  isLoading={activityFeed.isLoading}
                  isLoadingMore={activityFeed.isLoadingMore}
                  items={activityFeed.items}
                  showFlat
                  onItemPress={openEntryDetail}
                  onLoadMore={() => {
                    void activityFeed.loadMore();
                  }}
                  onViewAll={goLogs}
                />
              </>
            ) : null}
          </Stack>
        )}
      </GuardScreenShell>
    </View>
  );
}

const styles = StyleSheet.create({
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
});
