import { useCallback, useState, type ReactNode } from "react";
import { Pressable, ScrollView, StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useFocusEffect, useRouter } from "expo-router";
import { SymbolView } from "expo-symbols";
import { Stack, Row } from "@/components/layout";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { MonthPicker } from "@/components/ui/month-picker";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import {
  AppText,
  Badge,
  Button,
  Card,
  EmptyState,
  FilterChip,
  LoadingState,
  PaginatedList,
} from "@/components/ui";
import { ResidentScreenShell } from "@/features/resident/components/resident-screen-shell";
import { ResidenceSwitchSheet } from "@/features/profile/components/residence-switch-sheet";
import { useResident } from "@/features/resident/resident-context";
import {
  residentDashboardRoute,
  residentMaintenanceBillRoute,
  residentMaintenancePaymentRoute,
} from "@/features/resident/resident-routes";
import { formatDateOnly } from "@/features/guard/guard-utils";
import { getFriendlyApiMessage } from "@/features/auth/api-error";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";
import {
  useMaintenanceOutstandingQuery,
  type MaintenanceBill,
  type MaintenanceScope,
} from "@/lib/api/maintenance-api";
import { BillListItem } from "./maintenance-components";
import {
  formatBillingMonth,
  formatMoney,
  getMaintenanceBillPresentation,
} from "./maintenance-presentation";
import { useMaintenanceHistory } from "./use-maintenance-history";

const filters = [
  ["", "All"],
  ["unpaid", "Unpaid"],
  ["pending_review", "Pending"],
  ["paid", "Paid"],
  ["overdue", "Overdue"],
  ["rejected", "Rejected"],
  ["outstanding", "Outstanding"],
] as const;
export default function MaintenanceScreen() {
  const { societyId, flatId } = useResident();
  if (!societyId || !flatId) return <ResidentScreenShell />;
  return (
    <MaintenanceOverview
      key={`${societyId}:${flatId}`}
      societyId={societyId}
      flatId={flatId}
    />
  );
}
function MaintenanceOverview(scope: MaintenanceScope) {
  const router = useRouter();
  const resident = useResident();
  const [switchOpen, setSwitchOpen] = useState(false);
  const [month, setMonth] = useState("");
  const [status, setStatus] = useState("");
  const summary = useMaintenanceOutstandingQuery(scope, {
    refetchOnMountOrArgChange: true,
  });
  const refetch = summary.refetch;
  useFocusEffect(
    useCallback(() => {
      void refetch();
    }, [refetch]),
  );
  const current = summary.currentData;
  const bill = current?.current_bill;
  const presentation = bill
    ? getMaintenanceBillPresentation({ bill })
    : undefined;
  const openBill = (item: MaintenanceBill) =>
    router.push(
      residentMaintenanceBillRoute(item.id, scope.societyId, scope.flatId),
    );
  const header = (
    <Stack gap="lg">
      <ScreenBackHeader fallbackHomeRoute={residentDashboardRoute()} title="" trailing={
        <Pressable
          accessibilityRole="button"
          accessibilityLabel="Switch residence"
          onPress={() => setSwitchOpen(true)}
          style={[styles.headerButton, styles.flatButton]}
        >
          <SymbolView
            name={{ ios: "house", android: "home", web: "home" }}
            tintColor={colors.brand.navy}
            size={20}
          />
          <AppText numberOfLines={1} style={[styles.bold, { flexShrink: 1 }]}>
            {[
              resident.selectedResidence?.block,
              resident.selectedResidence?.flat_number,
            ]
              .filter(Boolean)
              .join("-") || "Your flat"}{" "}
            ⌄
          </AppText>
        </Pressable>
      } />
      <Stack gap="xs">
        <AppText variant="titleLarge" style={styles.title}>
          Maintenance
        </AppText>
        <AppText color="muted">View and pay your maintenance bills</AppText>
      </Stack>
      {current && !summary.isError ? (
        <Card style={styles.currentCard}>
          <Row justify="space-between" style={{ gap: spacing.sm, flexWrap: "wrap" }}>
            <AppText style={styles.bold}>Current Due</AppText>
            {presentation ? <Badge label={presentation.label} tone={presentation.tone} /> : null}
          </Row>
          <Stack gap="sm">
            <AppText variant="titleLarge" style={styles.dueAmount}>
              {formatMoney(current.current_month_paise)}
            </AppText>
            <AppText style={styles.bold}>{formatBillingMonth(current.current_month)}</AppText>
            <AppText color="muted" variant="bodySmall">
              {bill ? `Due ${formatDateOnly(bill.due_date, "UTC")}` : "No bill generated for this month"}
            </AppText>
          </Stack>
          {bill?.due_message ? (
            <Row justify="space-between" style={styles.dueNote}>
              <SymbolView
                name={{
                  ios: "calendar",
                  android: "calendar_month",
                  web: "calendar_month",
                }}
                tintColor={colors.brand.orange}
                size={20}
              />
              <AppText style={styles.bold}>{bill.due_message}</AppText>
            </Row>
          ) : null}
          {bill ? (
            <Button
              title={
                presentation?.action === "pay"
                  ? presentation.actionLabel
                  : "View Current Bill"
              }
              onPress={() =>
                presentation?.action === "pay"
                  ? router.push(
                      residentMaintenancePaymentRoute(
                        bill.id,
                        scope.societyId,
                        scope.flatId,
                      ),
                    )
                  : openBill(bill)
              }
            />
          ) : null}
        </Card>
      ) : summary.isError ? (
        <EmptyState
          title="Summary unavailable"
          message={getFriendlyApiMessage(summary.error, "Please try again.")}
          actionLabel="Retry"
          onAction={() => void refetch()}
        />
      ) : (
        <LoadingState message="Loading summary" />
      )}
      <Row style={{ gap: spacing.sm, alignItems: "stretch" }}>
        {[
          {
            label: "Outstanding",
            value: current?.total_outstanding_paise,
            filter: "outstanding",
            color: colors.brand.orange,
            background: colors.brand.orangeSoft,
          },
          {
            label: "Overdue",
            value: current?.overdue_paise,
            filter: "overdue",
            color: colors.status.error,
            background: colors.status.errorSoft,
          },
          {
            label: "Total Paid",
            value: current?.total_paid_paise,
            filter: "paid",
            color: colors.status.success,
            background: colors.status.successSoft,
          },
        ].map((stat) => (
          <Pressable
            key={stat.label}
            accessibilityRole="button"
            accessibilityLabel={`Filter ${stat.label} bills`}
            onPress={() => setStatus(stat.filter)}
            style={styles.statCard}
          >
            <View
              style={[styles.statIcon, { backgroundColor: stat.background }]}
            >
              <AppText
                style={{ color: stat.color, fontSize: 20, fontWeight: "800" }}
              >
                {stat.filter === "paid"
                  ? "✓"
                  : stat.filter === "overdue"
                    ? "!"
                    : "₹"}
              </AppText>
            </View>
            <AppText numberOfLines={1} adjustsFontSizeToFit minimumFontScale={0.65} style={styles.statValue}>
              {!summary.isError && typeof stat.value === "number"
                ? formatMoney(stat.value)
                : "—"}
            </AppText>
            <AppText variant="caption" color="muted" style={{ fontSize: 11 }}>
              {stat.label}
            </AppText>
          </Pressable>
        ))}
      </Row>
      <Stack gap="md" style={styles.tableHeader}>
        <AppText variant="title">Bills</AppText>
        <MonthPicker value={month} onChange={setMonth} />
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ gap: spacing.sm }}
        >
          {filters.map(([value, label]) => (
            <FilterChip
              key={value}
              label={label}
              selected={status === value}
              onPress={() => setStatus(value)}
            />
          ))}
        </ScrollView>
      </Stack>
    </Stack>
  );
  return (
    <SafeAreaView edges={["top", "left", "right"]} style={styles.screen}>
      <AppStatusBar />
      <BillHistory
        scope={scope}
        month={month}
        status={status}
        header={header}
        onOpen={openBill}
        refreshSummary={() => void refetch()}
      />
      <ResidenceSwitchSheet
        visible={switchOpen}
        onClose={() => setSwitchOpen(false)}
        onSelect={resident.selectResidence}
        selectedFlatId={scope.flatId}
        residences={resident.residences}
      />
    </SafeAreaView>
  );
}
function BillHistory({
  scope,
  month,
  status,
  header,
  onOpen,
  refreshSummary,
}: {
  scope: MaintenanceScope;
  month: string;
  status: string;
  header: ReactNode;
  onOpen: (bill: MaintenanceBill) => void;
  refreshSummary: () => void;
}) {
  const history = useMaintenanceHistory({
    ...scope,
    month: month || undefined,
    status: status || undefined,
    mode: "scroll",
  });
  const data = history.data;
  const page = data?.page ?? 1;
  return (
    <PaginatedList
      data={data?.items ?? []}
      keyExtractor={(bill) => String(bill.id)}
      renderItem={({ item }) => (
        <BillListItem bill={item} onPress={() => onOpen(item)} />
      )}
      header={header}
      isLoading={history.busy}
      isRefreshing={history.busy && !data}
      isLoadingMore={history.busy && Boolean(data)}
      hasMore={data?.has_more && !history.error}
      onLoadMore={() => void history.load(page + 1, true)}
      onRefresh={() => {
        void history.load();
        refreshSummary();
      }}
      contentContainerStyle={styles.list}
      emptyTitle={month || status ? "No bills found for this filter." : "No maintenance bills yet"}
      emptyMessage={month || status ? "Try another month or status." : "Maintenance bills from your society will appear here."}
      emptyComponent={history.error ? <View /> : undefined}
      footer={
        <Stack gap="md" style={styles.footer}>
          {history.error ? (
            <EmptyState
              title="Bills unavailable"
              message={getFriendlyApiMessage(
                history.error,
                "Could not load this page. Please retry.",
              )}
              actionLabel="Retry"
              onAction={() => void history.retry()}
            />
          ) : null}
          {data?.has_more && !history.error ? (
            <Button
              title="Load more bills"
              variant="secondary"
              loading={history.busy}
              onPress={() => void history.load(page + 1, true)}
            />
          ) : null}
        </Stack>
      }
    />
  );
}
const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.surface.screen },
  list: {
    width: "100%",
    maxWidth: 840,
    alignSelf: "center",
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    gap: spacing.md,
    paddingBottom: spacing["3xl"] + spacing.lg,
  },
  headerButton: {
    minHeight: 48,
    minWidth: 48,
    padding: spacing.sm,
    backgroundColor: colors.surface.card,
    borderRadius: 16,
    alignItems: "center",
    justifyContent: "center",
    borderWidth: 1,
    borderColor: colors.border.default,
  },
  flatButton: {
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    maxWidth: "75%",
  },
  title: {
    color: colors.brand.navy,
    fontWeight: "800",
    fontSize: 28,
    lineHeight: 36,
  },
  bold: { fontWeight: "700", color: colors.brand.navy },
  currentCard: {
    backgroundColor: colors.surface.card,
    shadowOpacity: 0,
    elevation: 0,
    borderRadius: 24,
    borderColor: colors.border.default,
    padding: spacing.lg,
    gap: spacing.md,
  },
  dueAmount: {
    color: colors.brand.orange,
    fontSize: 36,
    lineHeight: 44,
    fontWeight: "800",
  },
  dueNote: { gap: spacing.sm, paddingVertical: spacing.xs },
  statCard: {
    flex: 1,
    minWidth: 0,
    borderRadius: 18,
    padding: spacing.sm,
    gap: spacing.xs,
    backgroundColor: colors.surface.card,
    borderWidth: 1,
    borderColor: colors.border.default,
  },
  statIcon: {
    width: 36,
    height: 36,
    borderRadius: 18,
    justifyContent: "center",
    alignItems: "center",
  },
  statValue: {
    fontSize: 17,
    lineHeight: 24,
    fontWeight: "800",
    color: colors.brand.navy,
  },
  tableHeader: { gap: spacing.md },
  footer: { paddingTop: spacing.sm },
});
