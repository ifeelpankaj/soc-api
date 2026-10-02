import { useCallback } from "react";
import { useFocusEffect } from "expo-router";
import { AppText, Button, Card } from "@/components/ui";
import { InfoRow } from "@/components/ui/info-row";
import { DocumentAction } from "@/features/shared/document-action";
import { formatDateOnly } from "@/features/guard/guard-utils";
import {
  maintenancePath,
  useMaintenanceRecordsQuery,
  type BillArgs,
  type MaintenanceBill,
} from "@/lib/api/maintenance-api";
import {
  getMaintenanceBillPresentation,
  resolvePaymentContext,
} from "./maintenance-presentation";
import { maintenanceStyles as styles } from "./maintenance-components";

export function PaymentStatusCard({
  args,
  bill,
  refreshBill,
}: {
  args: BillArgs;
  bill: MaintenanceBill;
  refreshBill: () => void;
}) {
  const query = useMaintenanceRecordsQuery(
    { ...args, month: bill.billing_month },
    { refetchOnMountOrArgChange: true },
  );
  const refetch = query.refetch;
  useFocusEffect(
    useCallback(() => {
      void refetch();
    }, [refetch]),
  );
  const context = resolvePaymentContext(
    bill.id,
    query.isError ? undefined : query.currentData,
  );
  const presentation = getMaintenanceBillPresentation({
    bill,
    paymentContext: context,
  });
  const claim =
    context.claim?.status === bill.payment_claim_status
      ? context.claim
      : undefined;
  return (
    <Card style={styles.card}>
      <AppText style={styles.bold}>{presentation.label}</AppText>
      {query.isFetching ? (
        <AppText color="muted">Refreshing payment details…</AppText>
      ) : null}
      {context.complete ? (
        <>
          {claim ? (
            <>
              <AppText style={styles.bold}>Your latest submission</AppText>
              <AppText variant="caption" color="muted">
                Bill status includes all residents. These submission details are
                visible to your account only.
              </AppText>
              <InfoRow label="Transaction Reference" value={claim.reference} />
              <InfoRow
                label="Payment Date"
                value={formatDateOnly(claim.payment_date, "UTC")}
              />
              <InfoRow
                label="Submitted On"
                value={formatDateOnly(claim.created_at, bill.timezone)}
              />
              {claim.reviewed_at ? (
                <InfoRow
                  label="Reviewed On"
                  value={formatDateOnly(claim.reviewed_at, bill.timezone)}
                />
              ) : null}
              {claim.reason ? (
                <AppText color="error">{claim.reason}</AppText>
              ) : null}
            </>
          ) : presentation.state !== "PAID" ? (
            <AppText color="muted">
              Payment details may belong to another resident and are not
              available to this account.
            </AppText>
          ) : null}
          {presentation.receiptPayments.map((payment) => (
            <Card key={payment.id} style={styles.card}>
              <InfoRow label="Receipt Number" value={payment.receipt_number} />
              <InfoRow
                label="Paid On"
                value={formatDateOnly(payment.credit_date, "UTC")}
              />
              <DocumentAction
                path={`${maintenancePath(args.societyId)}/my/payments/${payment.id}/receipt/pdf`}
                filename={`receipt-${payment.id}.pdf`}
                title="View Receipt"
              />
            </Card>
          ))}
          {presentation.state === "PAID" &&
          !presentation.receiptPayments.length ? (
            <AppText color="muted">
              The bill is paid, but a verified receipt is not available to this
              account. Contact your society administrator.
            </AppText>
          ) : null}
        </>
      ) : (
        <AppText color="muted">
          {query.isError
            ? "Payment details could not be loaded."
            : query.currentData
              ? "Payment history exceeds the supported page size. Full payment details and receipt lookup are unavailable; contact your society administrator."
              : "Loading payment details…"}
        </AppText>
      )}
      <Button
        title="Refresh Payment Status"
        variant="secondary"
        loading={query.isFetching}
        onPress={() => {
          refreshBill();
          void refetch();
        }}
      />
    </Card>
  );
}
