import { useCallback, useEffect, useRef, useState } from "react";
import { AppState, Linking, ScrollView } from "react-native";
import { Image } from "expo-image";
import { useFocusEffect, useRouter } from "expo-router";
import { setStringAsync } from "expo-clipboard";
import { nanoid } from "@reduxjs/toolkit";
import {
  AppText,
  Badge,
  Button,
  Card,
  EmptyState,
  Input,
  LoadingState,
} from "@/components/ui";
import { InfoRow } from "@/components/ui/info-row";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import {
  residentMaintenanceRoute,
  residentMaintenanceBillRoute,
} from "@/features/resident/resident-routes";
import { getFriendlyApiMessage } from "@/features/auth/api-error";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { useReadDocumentMutation } from "@/lib/api/document-api";
import {
  maintenancePath,
  useMaintenanceBillQuery,
  useCreateMaintenancePaymentRequestMutation,
  useSubmitMaintenanceClaimMutation,
  type BillArgs,
  type MaintenanceBill,
  type PaymentRequest,
} from "@/lib/api/maintenance-api";
import {
  dateInZone,
  formatBillingMonth,
  formatMoney,
  getMaintenanceBillPresentation,
  validateClaim,
} from "./maintenance-presentation";
import {
  createSubmissionAttempt,
  isDefiniteRejection,
  type SubmissionAttempt,
} from "./payment-attempt";
import { maintenanceStyles as styles } from "./maintenance-components";

export function PaymentScreen(args: BillArgs) {
  const query = useMaintenanceBillQuery(args, {
    refetchOnMountOrArgChange: true,
  });
  const refetch = query.refetch;
  useFocusEffect(
    useCallback(() => {
      void refetch();
      const subscription = AppState.addEventListener("change", (state) => {
        if (state === "active") void refetch();
      });
      return () => subscription.remove();
    }, [refetch]),
  );
  const bill = query.currentData;
  const valid =
    bill && bill.flat_id === args.flatId && bill.society_id === args.societyId;
  const lostAccess =
    query.error &&
    "status" in query.error &&
    [401, 403, 404].includes(Number(query.error.status));
  return (
    <ResidentSubScreen
      title="Make Payment"
      fallbackHomeRoute={residentMaintenanceRoute()}
    >
      <ScrollView
        keyboardShouldPersistTaps="handled"
        contentContainerStyle={styles.content}
      >
        {query.isError ? (
          <EmptyState
            title="Payment unavailable"
            message={getFriendlyApiMessage(
              query.error,
              "Unable to load the current bill.",
            )}
            actionLabel="Retry"
            onAction={() => void refetch()}
          />
        ) : null}
        {!bill ? (
          !query.isError && <LoadingState message="Loading bill" />
        ) : !valid ? (
          <EmptyState
            title="Residence changed"
            message="Return to Maintenance to select your bill."
          />
        ) : !lostAccess ? (
          <PaymentForm
            args={args}
            bill={bill}
            refreshing={query.isFetching || query.isError}
          />
        ) : null}
      </ScrollView>
    </ResidentSubScreen>
  );
}

function PaymentForm({
  args,
  bill,
  refreshing,
}: {
  args: BillArgs;
  bill: MaintenanceBill;
  refreshing: boolean;
}) {
  const router = useRouter();
  const feedback = useAppFeedback();
  const [createRequest] = useCreateMaintenancePaymentRequestMutation();
  const [submitClaim] = useSubmitMaintenanceClaimMutation();
  const [readBinary, { reset: resetBinary }] = useReadDocumentMutation();
  const [request, setRequest] = useState<PaymentRequest>();
  const [qr, setQr] = useState<string>();
  const [qrError, setQrError] = useState<string>();
  const [reference, setReference] = useState("");
  const [paymentDate, setPaymentDate] = useState(() =>
    dateInZone(new Date(), bill.timezone),
  );
  const [error, setError] = useState<string>();
  const [busy, setBusy] = useState(false);
  const [uncertain, setUncertain] = useState(false);
  const requestKey = useRef(nanoid());
  const attempt = useRef<SubmissionAttempt | undefined>(undefined);
  const inFlight = useRef(false);
  const active = useRef(false);
  useEffect(() => {
    active.current = true;
    return () => {
      active.current = false;
    };
  }, []);
  const status = getMaintenanceBillPresentation({ bill });
  const goDetails = () =>
    router.dismissTo(
      residentMaintenanceBillRoute(bill.id, args.societyId, args.flatId),
    );
  const loadQR = async (paymentRequest: PaymentRequest) => {
    setQrError(undefined);
    try {
      const bytes = await readBinary({
        path: `${maintenancePath(args.societyId)}/my/payment-requests/${paymentRequest.id}/qr`,
      }).unwrap();
      if (active.current)
        setQr(
          `data:image/png;base64,${btoa(bytes.map((b) => String.fromCharCode(b)).join(""))}`,
        );
    } catch (e) {
      if (active.current)
        setQrError(
          getFriendlyApiMessage(e, "QR unavailable. Retry or use the UPI ID."),
        );
    } finally {
      resetBinary();
    }
  };
  const startPayment = async () => {
    if (inFlight.current || refreshing) return;
    inFlight.current = true;
    setBusy(true);
    setError(undefined);
    try {
      const result = await createRequest({
        ...args,
        key: requestKey.current,
      }).unwrap();
      if (!active.current) return;
      if (result.bill_id !== bill.id || result.state === "closed")
        throw new Error("Request unavailable");
      setRequest(result);
      void loadQR(result);
    } catch (e) {
      if (active.current)
        setError(
          getFriendlyApiMessage(
            e,
            "Could not prepare UPI payment. Please retry or contact your society administrator.",
          ),
        );
    } finally {
      inFlight.current = false;
      if (active.current) setBusy(false);
    }
  };
  const submit = async () => {
    if (!request || inFlight.current || refreshing) return;
    const validation = validateClaim(reference, paymentDate, bill.timezone);
    if (validation) {
      setError(validation);
      return;
    }
    inFlight.current = true;
    setBusy(true);
    setError(undefined);
    attempt.current ??= createSubmissionAttempt({
      payment_request_id: request.id,
      reference: reference.trim().toUpperCase(),
      payment_date: paymentDate,
    });
    try {
      await submitClaim({
        ...args,
        ...attempt.current,
      }).unwrap();
      if (!active.current) return;
      goDetails();
    } catch (e) {
      if (!active.current) return;
      const definite = isDefiniteRejection(e);
      if (definite) attempt.current = undefined;
      setUncertain(!definite);
      setError(
        getFriendlyApiMessage(
          e,
          "Submission could not be confirmed. Retry the same details.",
        ),
      );
    } finally {
      inFlight.current = false;
      if (active.current) setBusy(false);
    }
  };
  const copy = async () => {
    try {
      await setStringAsync(request!.destination.upi_id);
      feedback.showSuccess("Copied", "UPI ID copied.");
    } catch (e) {
      feedback.showError("Copy failed", e, "Select and copy the UPI ID below.");
    }
  };
  const openUPI = async () => {
    try {
      if (!request?.upi_uri.startsWith("upi://pay?"))
        throw new Error("Invalid UPI link");
      await Linking.openURL(request.upi_uri);
    } catch (e) {
      feedback.showError(
        "Could not open UPI",
        e,
        "Use the QR code or copy the UPI ID into your payment app.",
      );
    }
  };
  if (status.action !== "pay")
    return (
      <Card style={styles.card}>
        <Badge label={status.label} tone={status.tone} />
        <AppText>
          Your payment status comes from your society administrator. View bill
          details for the latest update.
        </AppText>
        <Button title="View Bill Details" onPress={goDetails} />
      </Card>
    );
  return (
    <>
      <Card style={styles.card}>
        <AppText variant="title">Pay via UPI</AppText>
        <AppText color="muted">
          {formatBillingMonth(bill.billing_month)}
        </AppText>
        <AppText variant="titleLarge" color="brand">
          {formatMoney(request?.amount_paise ?? bill.outstanding_amount_paise)}
        </AppText>
        <AppText>
          Pay externally, then submit your transaction reference for
          administrator verification.
        </AppText>
      </Card>
      {!request ? (
        <Button
          title="Prepare UPI Payment"
          loading={busy}
          disabled={refreshing}
          onPress={() => void startPayment()}
        />
      ) : (
        <>
          <Card style={styles.card}>
            {qr ? (
              <Image
                source={{ uri: qr }}
                accessibilityLabel="UPI payment QR code"
                style={{ width: 240, height: 240, alignSelf: "center" }}
                contentFit="contain"
              />
            ) : (
              <AppText>{qrError ?? "Loading QR…"}</AppText>
            )}
            {qrError ? (
              <Button
                title="Retry QR"
                variant="secondary"
                onPress={() => void loadQR(request)}
              />
            ) : null}
            <InfoRow label="UPI ID" value={request.destination.upi_id} />
            <InfoRow
              label="Payee Name"
              value={request.destination.payee_name}
            />
            <Button
              title="Copy UPI ID"
              variant="secondary"
              onPress={() => void copy()}
            />
            <Button
              title="Open UPI App"
              disabled={refreshing || uncertain}
              onPress={() => void openUPI()}
            />
          </Card>
          <Card style={styles.card}>
            <Input
              label="Transaction Reference"
              placeholder="Enter UTR / reference number"
              value={reference}
              onChangeText={setReference}
              autoCapitalize="characters"
              autoCorrect={false}
              maxLength={64}
              editable={!busy && !uncertain}
            />
            <Input
              label="Payment Date (YYYY-MM-DD)"
              placeholder="YYYY-MM-DD"
              value={paymentDate}
              onChangeText={setPaymentDate}
              maxLength={10}
              editable={!busy && !uncertain}
            />
            <AppText variant="caption" color="muted">
              Returning from a UPI app does not confirm payment. Submit only
              after paying. Your society administrator will verify the transfer.
            </AppText>
            {uncertain ? (
              <AppText color="warning">
                The last attempt is unconfirmed. Retry with the same details; do
                not pay again.
              </AppText>
            ) : null}
          </Card>
          <Button
            title={uncertain ? "Retry Submission" : "Submit Payment Details"}
            loading={busy}
            disabled={refreshing}
            onPress={() => void submit()}
          />
        </>
      )}
      {error ? <AppText color="error">{error}</AppText> : null}
    </>
  );
}
