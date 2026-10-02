import { Redirect, useFocusEffect, useRouter } from "expo-router";
import { SymbolView } from "expo-symbols";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  BackHandler,
  FlatList,
  Platform,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { EmptyState } from "@/components/ui/empty-state";
import { LoadingState } from "@/components/ui/loading-state";
import { useAuth } from "@/features/auth/use-auth";
import {
  notificationsForHomeRoute,
  notificationPresentation,
  notificationRoute,
  unreadCountForHomeRoute,
} from "@/features/notifications/notification-routing";
import {
  type AppNotification,
  useGetV1MeNotificationsQuery,
  useGetV1MeNotificationsUnreadCountQuery,
  usePatchV1MeNotificationsAndNotificationIdReadMutation,
  usePostV1MeNotificationsReadAllMutation,
} from "@/lib/api/notification-api-extensions";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

import { mergeNotificationPages } from "@/features/web/notification-tracker";
const PAGE_SIZE = 20;

export default function NotificationsScreen() {
  const router = useRouter();
  const { homeRoute, status, user } = useAuth();
  const fallbackHomeRoute = homeRoute ?? "/";
  const handleBack = useBackAction(fallbackHomeRoute);
  const [manualRefreshing, setManualRefreshing] = useState(false);
  const [cursor, setCursor] = useState<string | undefined>();
  const [items, setItems] = useState<AppNotification[]>([]);
  const loadingCursorRef = useRef<string | undefined>(undefined);
  const refetchFirstPageRef = useRef(false);
  const query = useGetV1MeNotificationsQuery({ limit: PAGE_SIZE, cursor },{skip:status !== "authenticated"});
  const firstPage = useGetV1MeNotificationsQuery({limit: PAGE_SIZE}, {skip:Platform.OS !== "web" || status !== "authenticated"});
  const unreadQuery = useGetV1MeNotificationsUnreadCountQuery(undefined,{skip:status !== "authenticated"});
  const [markRead] = usePatchV1MeNotificationsAndNotificationIdReadMutation();
  const [markAllRead, markAllState] = usePostV1MeNotificationsReadAllMutation();

  useEffect(() => {
    // Clear component-local pages too when the authenticated account changes.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setItems([]);
    setCursor(undefined);
  },[user?.id]);

  const incoming = useMemo(
    () => query.data?.data?.items ?? [],
    [query.data?.data?.items],
  );
  const nextCursor = query.data?.data?.next_cursor;
  const hasMore = Boolean(nextCursor);
  const isLoadingMore = query.isFetching && Boolean(cursor);
  const unreadCount = unreadQuery.data?.data?.unread_count ?? 0;
  const visibleItems = useMemo(
    () => notificationsForHomeRoute(items, homeRoute),
    [homeRoute, items],
  );
  const visibleUnreadCount = useMemo(
    () => unreadCountForHomeRoute(items, homeRoute, unreadCount),
    [homeRoute, items, unreadCount],
  );

  useEffect(() => {
    if (!query.data?.data?.items) {
      return;
    }
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setItems((current) =>
      Platform.OS === "web" ? mergeNotificationPages(current,incoming) : cursor ? mergeNotifications(current, incoming) : incoming,
    );
    loadingCursorRef.current = undefined;
  }, [cursor, incoming, query.data?.data?.items]);

  useEffect(() => {
    if (cursor || !refetchFirstPageRef.current) {
      return;
    }
    refetchFirstPageRef.current = false;
    void query.refetch();
  }, [cursor, query]);

  useEffect(() => {
    if (query.isError) {
      loadingCursorRef.current = undefined;
    }
  }, [query.isError]);

  useEffect(() => {
    if (Platform.OS === "web" && firstPage.data?.data?.items) {
      // Preserve loaded pages and scroll during foreground polling.
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setItems((current) => mergeNotificationPages(current,firstPage.data!.data!.items!));
    }
  }, [firstPage.data]);

  const refresh = useCallback(() => {
    if (Platform.OS === "web") {
      setManualRefreshing(true);
      void Promise.allSettled([firstPage.refetch(),unreadQuery.refetch()]).finally(() => setManualRefreshing(false));
      return;
    }
    loadingCursorRef.current = undefined;
    refetchFirstPageRef.current = true;
    setCursor(undefined);
    if (!cursor) {
      refetchFirstPageRef.current = false;
      void query.refetch();
    }
    void unreadQuery.refetch();
  }, [cursor, firstPage, query, unreadQuery]);

  const loadMore = useCallback(() => {
    if (
      hasMore &&
      !query.isFetching &&
      nextCursor &&
      loadingCursorRef.current !== nextCursor
    ) {
      loadingCursorRef.current = nextCursor;
      setCursor((current) => (current === nextCursor ? current : nextCursor));
    }
  }, [hasMore, nextCursor, query.isFetching]);

  useEffect(() => {
    if (
      visibleItems.length === 0 &&
      items.length > 0 &&
      hasMore &&
      !query.isFetching
    ) {
      loadMore();
    }
  }, [hasMore, items.length, loadMore, query.isFetching, visibleItems.length]);

  const handlePress = useCallback(
    async (item: AppNotification) => {
      if (!item.read_at) {
        await markRead({ notificationId: item.id })
          .unwrap()
          .catch(() => undefined);
        setItems((current) =>
          current.map((notification) =>
            notification.id === item.id
              ? { ...notification, read_at: new Date().toISOString() }
              : notification,
          ),
        );
        void unreadQuery.refetch();
      }
      if (Platform.OS === "web") {
        router.push({pathname:"/notification-open",params:{id:item.id,recipient:String(item.user_id)}});
        return;
      }
      const route = notificationRoute(item, homeRoute);
      if (route) {
        router.push(route);
      }
    },
    [homeRoute, markRead, router, unreadQuery],
  );

  const handleMarkAllRead = useCallback(async () => {
    await markAllRead()
      .unwrap()
      .catch(() => undefined);
    const readAt = new Date().toISOString();
    setItems((current) =>
      current.map((item) => ({ ...item, read_at: item.read_at ?? readAt })),
    );
    loadingCursorRef.current = undefined;
    refetchFirstPageRef.current = true;
    setCursor(undefined);
    if (!cursor) {
      refetchFirstPageRef.current = false;
      void query.refetch();
    }
    void unreadQuery.refetch();
  }, [cursor, markAllRead, query, unreadQuery]);

  useFocusEffect(
    useCallback(() => {
      if (Platform.OS !== "android") {
        return undefined;
      }

      const subscription = BackHandler.addEventListener(
        "hardwareBackPress",
        () => {
          handleBack();
          return true;
        },
      );

      return () => subscription.remove();
    }, [handleBack]),
  );

  const trailing = useMemo(
    () =>
      visibleUnreadCount > 0 ? (
        <Pressable
          accessibilityRole="button"
          disabled={markAllState.isLoading}
          hitSlop={8}
          style={({ pressed }) => [
            styles.markAllButton,
            pressed && styles.pressed,
          ]}
          onPress={handleMarkAllRead}
        >
          <Text style={styles.markAllText}>Read all</Text>
        </Pressable>
      ) : null,
    [handleMarkAllRead, markAllState.isLoading, visibleUnreadCount],
  );

  if (status === "unauthenticated") return <Redirect href="/login" />;

  return (
    <SafeAreaView style={styles.screen}>
      <View style={styles.content}>
        <ScreenBackHeader
          fallbackHomeRoute={fallbackHomeRoute}
          title="Notifications"
          trailing={trailing}
        />
        <View style={styles.summaryRow}>
          <Text style={styles.summaryText}>
            {visibleUnreadCount > 0
              ? `${visibleUnreadCount} unread`
              : "All caught up"}
          </Text>
        </View>
        {query.isLoading && visibleItems.length === 0 ? (
          <LoadingState message="Loading notifications" />
        ) : (
          <FlatList
            data={visibleItems}
            keyExtractor={(item) => item.id}
            ListEmptyComponent={
              <EmptyState
                message="New approvals, invite updates, and visitor activity will appear here."
                title="No notifications yet"
              />
            }
            ListFooterComponent={
              isLoadingMore && items.length > 0 ? (
                <Text style={styles.footerText}>Loading more...</Text>
              ) : null
            }
            refreshControl={
              <RefreshControl
                refreshing={Platform.OS === "web" ? manualRefreshing : query.isFetching && !cursor}
                onRefresh={refresh}
              />
            }
            renderItem={({ item }) => (
              <NotificationRow
                item={item}
                onPress={() => void handlePress(item)}
              />
            )}
            contentContainerStyle={styles.listContent}
            onEndReached={loadMore}
            onEndReachedThreshold={0.5}
          />
        )}
      </View>
    </SafeAreaView>
  );
}

function NotificationRow({
  item,
  onPress,
}: {
  item: AppNotification;
  onPress: () => void;
}) {
  const presentation = notificationPresentation(item);
  const unread = !item.read_at;
  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.row,
        unread && styles.rowUnread,
        pressed && styles.pressed,
      ]}
      onPress={onPress}
    >
      <View style={[styles.iconWrap, toneStyle[presentation.tone]]}>
        <SymbolView
          name={presentation.icon}
          size={20}
          tintColor={colors.text.inverse}
        />
      </View>
      <View style={styles.rowBody}>
        <View style={styles.rowTitleLine}>
          <Text style={styles.rowTitle}>{presentation.title}</Text>
          {unread ? <View style={styles.unreadDot} /> : null}
        </View>
        <Text style={styles.rowMessage}>{presentation.body}</Text>
        <Text style={styles.rowDate}>
          {formatNotificationTime(item.created_at)}
        </Text>
      </View>
    </Pressable>
  );
}

function mergeNotifications(
  current: AppNotification[],
  incoming: AppNotification[],
) {
  const byID = new Map(current.map((item) => [item.id, item]));
  for (const item of incoming) {
    byID.set(item.id, item);
  }
  return Array.from(byID.values()).sort((a, b) =>
    b.created_at.localeCompare(a.created_at),
  );
}

function formatNotificationTime(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "";
  }
  return date.toLocaleString(undefined, {
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
    month: "short",
  });
}

const toneStyle = StyleSheet.create({
  info: {
    backgroundColor: colors.brand.navy,
  },
  success: {
    backgroundColor: colors.status.success,
  },
  warning: {
    backgroundColor: colors.brand.orange,
  },
});

const styles = StyleSheet.create({
  content: {
    flex: 1,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: layout.screenPaddingTop,
  },
  footerText: {
    color: colors.text.muted,
    paddingVertical: spacing.md,
    textAlign: "center",
  },
  iconWrap: {
    alignItems: "center",
    borderRadius: radius.xl,
    height: 42,
    justifyContent: "center",
    width: 42,
  },
  listContent: {
    gap: spacing.sm,
    paddingBottom: spacing["3xl"],
  },
  markAllButton: {
    minHeight: 44,
    justifyContent: "center",
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
  },
  markAllText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
  },
  pressed: {
    opacity: 0.72,
  },
  row: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.xl,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    padding: spacing.md,
  },
  rowBody: {
    flex: 1,
    gap: 4,
  },
  rowDate: {
    color: colors.text.muted,
    fontSize: 12,
  },
  rowMessage: {
    color: colors.text.secondary,
    fontSize: 13,
    lineHeight: 18,
  },
  rowTitle: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 15,
    fontWeight: "700",
  },
  rowTitleLine: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
  rowUnread: {
    borderColor: colors.brand.orange,
  },
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
  summaryRow: {
    marginBottom: spacing.md,
  },
  summaryText: {
    ...typography.body,
    color: colors.text.muted,
  },
  unreadDot: {
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 8,
    width: 8,
  },
});
