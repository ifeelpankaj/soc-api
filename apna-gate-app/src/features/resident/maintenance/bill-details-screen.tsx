import { useCallback } from "react";
import { ScrollView, RefreshControl } from "react-native";
import { useFocusEffect, useRouter } from "expo-router";
import {
  AppText,
  Badge,
  Button,
  Card,
  EmptyState,
  LoadingState,
} from "@/components/ui";
import { InfoRow } from "@/components/ui/info-row";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import {
  residentMaintenanceRoute,
  residentMaintenancePaymentRoute,
} from "@/features/resident/resident-routes";
import { getFriendlyApiMessage } from "@/features/auth/api-error";
import { formatDateOnly } from "@/features/guard/guard-utils";
import { DocumentAction } from "@/features/shared/document-action";
import {
  useMaintenanceBillQuery,
  maintenancePath,
  type BillArgs,
} from "@/lib/api/maintenance-api";
import {
  formatMoney,
  formatBillingMonth,
  getMaintenanceBillPresentation,
} from "./maintenance-presentation";
import { maintenanceStyles as styles } from "./maintenance-components";
import { PaymentStatusCard } from "./payment-status-card";

export function BillDetailsScreen(args: BillArgs) {
  const router = useRouter();
  const query = useMaintenanceBillQuery(args, {
    refetchOnMountOrArgChange: true,
  });
  const refetch = query.refetch;
  useFocusEffect(
    useCallback(() => {
      void refetch();
    }, [refetch]),
  );
  const bill = query.currentData;
  const valid =
    bill && bill.flat_id === args.flatId && bill.society_id === args.societyId;
  const status = valid ? getMaintenanceBillPresentation({ bill }) : undefined;
  return (
    <ResidentSubScreen
      title="Bill Details"
      fallbackHomeRoute={residentMaintenanceRoute()}
    >
      <ScrollView
        contentContainerStyle={styles.content}
        refreshControl={
          <RefreshControl
            refreshing={query.isFetching}
            onRefresh={() => void refetch()}
          />
        }
      >
        {query.isError ? (
          <EmptyState
            title="Bill unavailable"
            message={getFriendlyApiMessage(
              query.error,
              "Unable to load this bill.",
            )}
            actionLabel="Retry"
            onAction={() => void refetch()}
          />
        ) : !bill ? (
          <LoadingState message="Loading bill" />
        ) : !valid ? (
          <EmptyState
            title="Residence changed"
            message="Open a bill belonging to the selected residence."
          />
        ) : (
          <>
            <Card style={styles.card}>
              <AppText style={styles.bold}>
                {formatBillingMonth(bill.billing_month)}
              </AppText>
              <Badge label={status!.label} tone={status!.tone} />
              <AppText variant="titleLarge" color="brand">
                {formatMoney(bill.total_paise)}
              </AppText>
              <AppText color="muted">
                Due {formatDateOnly(bill.due_date, "UTC")}
              </AppText>
            </Card>
            <Card>
              <AppText style={styles.bold}>Bill Information</AppText>
              <InfoRow label="Bill Number" value={bill.bill_number} />
              <InfoRow
                label="Billing Month"
                value={formatBillingMonth(bill.billing_month)}
              />
              <InfoRow
                label="Due Date"
                value={formatDateOnly(bill.due_date, "UTC")}
              />
              <InfoRow
                label="Generated On"
                value={formatDateOnly(
                  bill.issued_at || bill.created_at,
                  bill.timezone,
                )}
              />
            </Card>
            <Card>
              <AppText style={styles.bold}>Bill Breakdown</AppText>
              {bill.items.map((item, i) => (
                <InfoRow
                  key={i}
                  label={item.description}
                  value={formatMoney(item.amount_paise)}
                />
              ))}
              <InfoRow label="Total" value={formatMoney(bill.total_paise)} />
              <InfoRow
                label="Outstanding"
                value={formatMoney(bill.outstanding_amount_paise)}
              />
            </Card>
            <DocumentAction
              path={`${maintenancePath(args.societyId)}/my/bills/${bill.id}/invoice`}
              filename={`bill-${bill.id}.pdf`}
              title="View / Download Bill"
            />
            {status &&
            ["PAID", "PENDING_REVIEW", "REJECTED"].includes(status.state) ? (
              <PaymentStatusCard
                args={args}
                bill={bill}
                refreshBill={() => void refetch()}
              />
            ) : null}
            {status?.action === "pay" ? (
              <Button
                title={status.actionLabel}
                disabled={query.isFetching}
                onPress={() =>
                  router.push(
                    residentMaintenancePaymentRoute(
                      bill.id,
                      args.societyId,
                      args.flatId,
                    ),
                  )
                }
              />
            ) : null}
          </>
        )}
      </ScrollView>
    </ResidentSubScreen>
  );
}
