import type {
  MaintenanceBill,
  MaintenanceClaim,
  MaintenancePayment,
  CursorPage,
} from "@/lib/api/maintenance-api";

export function dateInZone(now: Date, timeZone: string) {
  const parts = new Intl.DateTimeFormat("en-CA", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).formatToParts(now);
  const value = (type: string) => parts.find((p) => p.type === type)!.value;
  return `${value("year")}-${value("month")}-${value("day")}`;
}
export function formatMoney(paise: number) {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: paise % 100 ? 2 : 0,
  }).format(paise / 100);
}
export function formatBillingMonth(month: string) {
  return new Date(`${month.slice(0, 7)}-01T12:00:00Z`).toLocaleDateString(
    "en-IN",
    { month: "long", year: "numeric", timeZone: "UTC" },
  );
}
export type PaymentContext = {
  claim?: MaintenanceClaim;
  payments: MaintenancePayment[];
  complete: boolean;
};
export function resolvePaymentContext(
  billId: number,
  records?: {
    claims: CursorPage<MaintenanceClaim>;
    payments: CursorPage<MaintenancePayment>;
  },
): PaymentContext {
  if (!records || records.claims.next_cursor || records.payments.next_cursor)
    return { payments: [], complete: false };
  // The server defines claim recency by ID. Resident lists expose only this
  // user's claims, so this is the latest visible submission, not necessarily
  // the bill's effective claim. Only bill.payment_claim_status drives the CTA.
  const claims = records.claims.items
    .filter((c) => c.bill_id === billId)
    .sort((a, b) => b.id - a.id);
  const payments = records.payments.items
    .filter((p) => p.bill_id === billId && p.status === "verified")
    .sort(
      (a, b) =>
        Date.parse(b.verified_at) - Date.parse(a.verified_at) || b.id - a.id,
    );
  return { claim: claims[0], payments, complete: true };
}
const presentations = {
  PAID: {
    label: "Paid",
    tone: "emerald",
    action: "receipt",
    actionLabel: "View Receipt",
  },
  PENDING_REVIEW: {
    label: "Pending Review",
    tone: "amber",
    action: "details",
    actionLabel: "View Payment Details",
  },
  REJECTED: {
    label: "Payment Rejected",
    tone: "rose",
    action: "pay",
    actionLabel: "Resubmit Payment",
  },
  OVERDUE: {
    label: "Overdue",
    tone: "rose",
    action: "pay",
    actionLabel: "Pay Now",
  },
  UNPAID: {
    label: "Unpaid",
    tone: "amber",
    action: "pay",
    actionLabel: "Pay Now",
  },
} as const;
export function getMaintenanceBillPresentation({
  bill,
  paymentContext,
}: {
  bill: MaintenanceBill;
  paymentContext?: PaymentContext;
  now?: Date;
}) {
  // Bill settlement and effective claim status are server projections. Historical
  // payment arrays must never override an unpaid/reversed bill projection.
  const state =
    bill.status === "paid" && bill.outstanding_amount_paise === 0
      ? "PAID"
      : bill.payment_claim_status === "pending"
        ? "PENDING_REVIEW"
        : bill.payment_claim_status === "rejected"
          ? "REJECTED"
          : bill.display_status === "overdue" || bill.status === "overdue"
            ? "OVERDUE"
            : "UNPAID";
  return {
    state,
    ...presentations[state],
    receiptPayments: state === "PAID" ? (paymentContext?.payments ?? []) : [],
  };
}
export function validateClaim(
  reference: string,
  date: string,
  timeZone: string,
  now = new Date(),
) {
  if (!/^[A-Z0-9/-]{6,64}$/.test(reference.trim().toUpperCase()))
    return "Reference must be 6–64 letters, digits, slashes or hyphens.";
  const parsed = new Date(`${date}T00:00:00Z`);
  if (
    !/^\d{4}-\d{2}-\d{2}$/.test(date) ||
    !Number.isFinite(parsed.getTime()) ||
    parsed.toISOString().slice(0, 10) !== date
  )
    return "Enter a valid payment date (YYYY-MM-DD).";
  if (date > dateInZone(now, timeZone))
    return "Payment date cannot be in the future.";
}
