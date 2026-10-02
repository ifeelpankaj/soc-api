import { useFocusEffect, useRouter } from "expo-router";
import { SymbolView, type SymbolViewProps } from "expo-symbols";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  AppState,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";

import { EmptyState } from "@/components/ui/empty-state";
import { LoadingState } from "@/components/ui/loading-state";
import { SegmentTabs } from "@/components/ui/segment-tabs";
import { StatusPill } from "@/components/ui/status-pill";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResident } from "@/features/resident/resident-context";
import { residentInviteDetailRoute } from "@/features/resident/resident-routes";
import { useDebouncedValue } from "@/features/shared/use-debounced-value";
import {
  type ModelsFlatMemberInviteStatus,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery,
} from "@/lib/api/resident-api-extensions";
import {
  type ModelsFlatMemberInviteResponse,
  type ModelsVisitorInviteHistoryItem,
  type ModelsVisitorInviteStatus,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesQuery,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type InviteTab = "visitor" | "member";
type StatusFilter =
  "all" | ModelsVisitorInviteStatus | ModelsFlatMemberInviteStatus;
type AppSymbolName = SymbolViewProps["name"];

const INVITE_LIST_SYMBOLS = {
  memberInvite: {
    ios: "person.2.badge.plus",
    android: "group_add",
    web: "group_add",
  },
  search: { ios: "magnifyingglass", android: "search", web: "search" },
  visitorInvite: {
    ios: "person.crop.circle.badge.plus",
    android: "person_add",
    web: "person_add",
  },
} satisfies Record<string, AppSymbolName>;

const PAGE_SIZE = 20;
const TABS: { value: InviteTab; label: string }[] = [
  { value: "visitor", label: "Guest Invites" },
  { value: "member", label: "Member Invites" },
];

const VISITOR_STATUSES: { value: StatusFilter; label: string }[] = [
  { value: "all", label: "All" },
  { value: "active", label: "Active" },
  { value: "used", label: "Used" },
  { value: "expired", label: "Expired" },
  { value: "cancelled", label: "Cancelled" },
];

const MEMBER_STATUSES: { value: StatusFilter; label: string }[] = [
  { value: "all", label: "All" },
  { value: "pending", label: "Pending" },
  { value: "accepted", label: "Accepted" },
  { value: "expired", label: "Expired" },
  { value: "cancelled", label: "Cancelled" },
];

export function ResidentInvitesScreen() {
  const router = useRouter();
  const { canManageFlatMembers, flatId, societyId } = useResident();
  const [tab, setTab] = useState<InviteTab>("visitor");
  const [status, setStatus] = useState<StatusFilter>("all");
  const [search, setSearch] = useState("");
  const [offset, setOffset] = useState(0);
  const [visitorItems, setVisitorItems] = useState<
    ModelsVisitorInviteHistoryItem[]
  >([]);
  const [memberItems, setMemberItems] = useState<
    ModelsFlatMemberInviteResponse[]
  >([]);
  const lastLoadMoreAt = useRef(0);
  const debouncedSearch = useDebouncedValue(search.trim(), 350);

  useEffect(() => {
    /* eslint-disable react-hooks/set-state-in-effect */
    setStatus("all");
    setOffset(0);
    /* eslint-enable react-hooks/set-state-in-effect */
  }, [tab]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setOffset(0);
  }, [debouncedSearch, status]);

  const baseArgs = {
    societyId: societyId ?? 0,
    flatId: flatId ?? 0,
    limit: PAGE_SIZE,
    offset,
    search: debouncedSearch || undefined,
  };
  const canQuery = Boolean(societyId && flatId);

  const visitorQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesQuery(
      {
        ...baseArgs,
        status:
          tab === "visitor" && status !== "all"
            ? (status as ModelsVisitorInviteStatus)
            : undefined,
      },
      { skip: !canQuery || tab !== "visitor" },
    );
  const memberQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery(
      {
        ...baseArgs,
        status:
          tab === "member" && status !== "all"
            ? (status as ModelsFlatMemberInviteStatus)
            : undefined,
      },
      { skip: !canQuery || tab !== "member" || !canManageFlatMembers },
    );

  useEffect(() => {
    const incoming = visitorQuery.data?.data?.invites;
    if (!incoming || tab !== "visitor") {
      return;
    }
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setVisitorItems((current) =>
      offset === 0 ? incoming : mergeById(current, incoming),
    );
  }, [offset, tab, visitorQuery.data?.data?.invites]);

  useEffect(() => {
    const incoming = memberQuery.data?.data?.invites;
    if (!incoming || tab !== "member") {
      return;
    }
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMemberItems((current) =>
      offset === 0 ? incoming : mergeById(current, incoming),
    );
  }, [memberQuery.data?.data?.invites, offset, tab]);

  const activeItems = tab === "visitor" ? visitorItems : memberItems;
  const total =
    tab === "visitor"
      ? visitorQuery.data?.data?.total
      : memberQuery.data?.data?.total;
  const isFetching =
    tab === "visitor" ? visitorQuery.isFetching : memberQuery.isFetching;
  const isLoading =
    tab === "visitor" ? visitorQuery.isLoading : memberQuery.isLoading;
  const statuses = tab === "visitor" ? VISITOR_STATUSES : MEMBER_STATUSES;
  const hasMore =
    typeof total === "number"
      ? activeItems.length < total
      : activeItems.length >= PAGE_SIZE;
  const visitorRefetch = visitorQuery.refetch;
  const visitorIsUninitialized = visitorQuery.isUninitialized;
  const memberRefetch = memberQuery.refetch;
  const memberIsUninitialized = memberQuery.isUninitialized;

  const refetchCurrentTab = useCallback(() => {
    if (tab === "visitor" && !visitorIsUninitialized) {
      void visitorRefetch();
    }
    if (tab === "member" && !memberIsUninitialized) {
      void memberRefetch();
    }
  }, [
    memberIsUninitialized,
    memberRefetch,
    tab,
    visitorIsUninitialized,
    visitorRefetch,
  ]);

  const refresh = useCallback(() => {
    setOffset(0);
    refetchCurrentTab();
  }, [refetchCurrentTab]);

  useFocusEffect(refetchCurrentTab);

  useEffect(() => {
    const subscription = AppState.addEventListener("change", (state) => {
      if (state === "active") {
        refetchCurrentTab();
      }
    });
    return () => subscription.remove();
  }, [refetchCurrentTab]);

  const loadMore = useCallback(() => {
    const now = Date.now();
    if (isFetching || !hasMore || now - lastLoadMoreAt.current < 900) {
      return;
    }
    lastLoadMoreAt.current = now;
    setOffset((value) => value + PAGE_SIZE);
  }, [hasMore, isFetching]);

  const content = useMemo(() => {
    if (tab === "member" && !canManageFlatMembers) {
      return (
        <EmptyState
          message="Only the flat owner or primary resident can manage member invites."
          title="Member invite access needed"
        />
      );
    }

    if (isLoading && activeItems.length === 0) {
      return <LoadingState message="Loading invites" />;
    }

    return (
      <FlatList<ModelsVisitorInviteHistoryItem | ModelsFlatMemberInviteResponse>
        data={activeItems}
        keyExtractor={(item) => `${tab}-${item.id ?? item.created_at}`}
        ListEmptyComponent={
          <EmptyState
            message={
              search
                ? "Try a different name, status, contact, or short code."
                : "Created invites will appear here."
            }
            title={search ? "No matching invites" : "No invites yet"}
          />
        }
        ListFooterComponent={
          isFetching && activeItems.length > 0 ? (
            <Text style={styles.footerText}>Loading more...</Text>
          ) : null
        }
        refreshControl={
          <RefreshControl
            refreshing={isFetching && offset === 0}
            onRefresh={refresh}
          />
        }
        renderItem={({ item }) =>
          tab === "visitor" ? (
            <VisitorInviteCard
              invite={item as ModelsVisitorInviteHistoryItem}
              onPress={() =>
                item.id &&
                router.push(residentInviteDetailRoute("visitor", item.id))
              }
            />
          ) : (
            <MemberInviteCard
              invite={item as ModelsFlatMemberInviteResponse}
              onPress={() =>
                item.id &&
                router.push(residentInviteDetailRoute("member", item.id))
              }
            />
          )
        }
        contentContainerStyle={styles.listContent}
        onEndReached={loadMore}
        onEndReachedThreshold={0.45}
      />
    );
  }, [
    activeItems,
    canManageFlatMembers,
    isFetching,
    isLoading,
    loadMore,
    offset,
    refresh,
    router,
    search,
    tab,
  ]);

  return (
    <ResidentSubScreen title="Manage Invites">
      <View style={styles.container}>
        <SegmentTabs options={TABS} value={tab} onChange={setTab} />
        <View style={styles.searchWrap}>
          <SymbolView
            name={INVITE_LIST_SYMBOLS.search}
            size={17}
            tintColor={colors.text.muted}
          />
          <TextInput
            placeholder="Search invites"
            placeholderTextColor={colors.text.placeholder}
            style={styles.searchInput}
            value={search}
            onChangeText={setSearch}
          />
        </View>
        <View style={styles.filterRow}>
          {statuses.map((option) => (
            <Pressable
              key={option.value}
              style={[
                styles.filterChip,
                status === option.value && styles.filterChipActive,
              ]}
              onPress={() => setStatus(option.value)}
            >
              <Text
                style={[
                  styles.filterText,
                  status === option.value && styles.filterTextActive,
                ]}
              >
                {option.label}
              </Text>
            </Pressable>
          ))}
        </View>
        {content}
      </View>
    </ResidentSubScreen>
  );
}

function VisitorInviteCard({
  invite,
  onPress,
}: {
  invite: ModelsVisitorInviteHistoryItem;
  onPress: () => void;
}) {
  const title =
    invite.visitor?.full_name ||
    `${titleize(invite.purpose ?? "Guest")} invite`;
  const subtitle = invite.entry?.status
    ? `Entry ${titleize(invite.entry.status)}`
    : invite.short_link?.url
      ? "Share link ready"
      : "Share link unavailable";

  return (
    <InviteCard
      icon={INVITE_LIST_SYMBOLS.visitorInvite}
      status={invite.status}
      subtitle={subtitle}
      title={title}
      date={invite.created_at}
      onPress={onPress}
    />
  );
}

function MemberInviteCard({
  invite,
  onPress,
}: {
  invite: ModelsFlatMemberInviteResponse;
  onPress: () => void;
}) {
  return (
    <InviteCard
      icon={INVITE_LIST_SYMBOLS.memberInvite}
      status={invite.status}
      subtitle={`${titleize(invite.role ?? "member")} invite`}
      title={invite.full_name || "Member invite"}
      date={invite.created_at}
      onPress={onPress}
    />
  );
}

function InviteCard({
  date,
  icon,
  onPress,
  status,
  subtitle,
  title,
}: {
  date?: string;
  icon: AppSymbolName;
  onPress: () => void;
  status?: string;
  subtitle: string;
  title: string;
}) {
  return (
    <Pressable
      style={({ pressed }) => [styles.card, pressed && styles.cardPressed]}
      onPress={onPress}
    >
      <View style={styles.cardIcon}>
        <SymbolView name={icon} size={22} tintColor={colors.brand.orange} />
      </View>
      <View style={styles.cardBody}>
        <Text numberOfLines={1} style={styles.cardTitle}>
          {title}
        </Text>
        <Text numberOfLines={1} style={styles.cardSubtitle}>
          {subtitle}
        </Text>
        <Text style={styles.cardDate}>
          {date ? new Date(date).toLocaleString() : "Created recently"}
        </Text>
      </View>
      {status ? <StatusPill status={status} /> : null}
    </Pressable>
  );
}

function mergeById<T extends { id?: number }>(current: T[], incoming: T[]) {
  const map = new Map<number | string, T>();
  for (const item of [...current, ...incoming]) {
    map.set(item.id ?? `${Math.random()}`, item);
  }
  return Array.from(map.values());
}

function titleize(value: string) {
  return value
    .replace(/[_-]+/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    gap: spacing.md,
    padding: spacing.md,
  },
  cardBody: {
    flex: 1,
    gap: 3,
  },
  cardDate: {
    color: colors.text.muted,
    fontSize: 11,
  },
  cardIcon: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 18,
    height: 42,
    justifyContent: "center",
    width: 42,
  },
  cardPressed: {
    opacity: 0.86,
  },
  cardSubtitle: {
    color: colors.text.secondary,
    fontSize: 12,
    fontWeight: "600",
  },
  cardTitle: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "800",
  },
  container: {
    flex: 1,
    gap: spacing.md,
    padding: spacing.md,
  },
  filterChip: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: 999,
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  filterChipActive: {
    backgroundColor: colors.brand.orange,
    borderColor: colors.brand.orange,
  },
  filterRow: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  filterText: {
    color: colors.text.secondary,
    fontSize: 12,
    fontWeight: "700",
  },
  filterTextActive: {
    color: colors.text.inverse,
  },
  footerText: {
    color: colors.text.muted,
    fontSize: 12,
    fontWeight: "600",
    paddingVertical: spacing.md,
    textAlign: "center",
  },
  listContent: {
    gap: spacing.md,
    paddingBottom: spacing["2xl"],
  },
  searchInput: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 15,
    fontWeight: "600",
    paddingVertical: spacing.sm,
  },
  searchWrap: {
    alignItems: "center",
    backgroundColor: colors.surface.input,
    borderColor: colors.border.input,
    borderRadius: radius.md,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
  },
});
