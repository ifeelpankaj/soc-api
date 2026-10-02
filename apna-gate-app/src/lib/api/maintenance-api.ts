import { enhancedApi } from "@/lib/api/enhanced-api";

// Resident DTOs mirror internal/models/maintenance*.go. The app's generated
// schema predates these endpoints; keep extensions on the existing API client.
export type MaintenanceBill = {
  id: number;
  society_id: number;
  flat_id: number;
  bill_number: string;
  billing_month: string;
  due_date: string;
  timezone: string;
  currency: string;
  total_paise: number;
  paid_amount_paise: number;
  outstanding_amount_paise: number;
  status: string;
  display_status?: string;
  due_message?: string;
  paid_on?: string;
  payment_claim_status: string;
  issued_at: string;
  created_at: string;
  items: { description: string; amount_paise: number }[];
};
export type MaintenanceClaim = {
  id: number;
  bill_id: number;
  status: string;
  reference: string;
  payment_date: string;
  created_at: string;
  reviewed_at?: string;
  reason?: string;
};
export type MaintenancePayment = {
  id: number;
  bill_id: number;
  status: string;
  amount_paise: number;
  verified_at: string;
  credit_date: string;
  reference?: string;
  receipt_number: string;
};
export type PaymentRequest = {
  id: string;
  bill_id: number;
  state: string;
  amount_paise: number;
  destination: { upi_id: string; payee_name: string; enabled: boolean };
  upi_uri: string;
  qr_url: string;
};
export type ClaimSubmission = {
  payment_request_id: string;
  reference: string;
  payment_date: string;
};
export type MaintenanceScope = { societyId: number; flatId: number };
export type BillArgs = MaintenanceScope & { billId: number };
export type CursorPage<T> = { items: T[]; next_cursor?: number };
export type MaintenanceBillsPage = CursorPage<MaintenanceBill> & {
  total_count: number;
  page: number;
  total_pages: number;
  has_more: boolean;
};
type Envelope<T> = { data: T };
type Outstanding = {
  flat_id: number;
  current_month: string;
  current_month_paise: number;
  total_outstanding_paise: number;
  previous_outstanding_paise: number;
  overdue_paise: number;
  total_paid_paise: number;
  current_bill?: MaintenanceBill;
};
export const maintenancePath = (societyId: number) =>
  `/v1/societies/${societyId}/maintenance`;
const api = enhancedApi.enhanceEndpoints({ addTagTypes: ["Maintenance"] });
const tag = ({ societyId, flatId }: MaintenanceScope) => [
  { type: "Maintenance" as const, id: `${societyId}:${flatId}` },
];
export const maintenanceApi = api.injectEndpoints({
  endpoints: (build) => ({
    maintenanceBills: build.query<
      MaintenanceBillsPage,
      MaintenanceScope & { page: number; month?: string; status?: string }
    >({
      query: ({ societyId, flatId, page, month, status }) => ({
        url: `${maintenancePath(societyId)}/my/bills`,
        params: {
          flat_id: flatId,
          page,
          billing_month: month,
          payment_status: status,
          limit: 20,
        },
      }),
      transformResponse: (r: Envelope<MaintenanceBillsPage>) => r.data,
      providesTags: (_r, _e, args) => tag(args),
      keepUnusedDataFor: 0,
    }),
    maintenanceOutstanding: build.query<Outstanding, MaintenanceScope>({
      query: ({ societyId, flatId }) =>
        `${maintenancePath(societyId)}/flats/${flatId}/outstanding`,
      transformResponse: (r: Envelope<Outstanding>) => r.data,
      providesTags: (_r, _e, args) => tag(args),
      keepUnusedDataFor: 0,
    }),
    maintenanceBill: build.query<MaintenanceBill, BillArgs>({
      query: ({ societyId, billId }) =>
        `${maintenancePath(societyId)}/my/bills/${billId}`,
      transformResponse: (r: Envelope<MaintenanceBill>) => r.data,
      providesTags: (_r, _e, args) => tag(args),
      keepUnusedDataFor: 0,
    }),
    maintenanceRecords: build.query<
      {
        claims: CursorPage<MaintenanceClaim>;
        payments: CursorPage<MaintenancePayment>;
      },
      BillArgs & { month: string }
    >({
      // These filters narrow authenticated resident ownership; they never grant access.
      async queryFn(
        { societyId, flatId, billId, month },
        _api,
        _extra,
        baseQuery,
      ) {
        const params = {
          bill_id: billId,
          flat_id: flatId,
          billing_month: month.slice(0, 7),
          limit: 100,
        };
        const claims = await baseQuery({
          url: `${maintenancePath(societyId)}/my/payment-claims`,
          params,
        });
        if (claims.error) return { error: claims.error };
        const payments = await baseQuery({
          url: `${maintenancePath(societyId)}/my/payments`,
          params,
        });
        if (payments.error) return { error: payments.error };
        return {
          data: {
            claims: (claims.data as Envelope<CursorPage<MaintenanceClaim>>)
              .data,
            payments: (
              payments.data as Envelope<CursorPage<MaintenancePayment>>
            ).data,
          },
        };
      },
      providesTags: (_r, _e, args) => tag(args),
      keepUnusedDataFor: 0,
    }),
    createMaintenancePaymentRequest: build.mutation<
      PaymentRequest,
      BillArgs & { key: string }
    >({
      query: ({ societyId, billId, key }) => ({
        url: `${maintenancePath(societyId)}/my/bills/${billId}/payment-request`,
        method: "POST",
        headers: { "Idempotency-Key": key },
      }),
      transformResponse: (r: Envelope<PaymentRequest>) => r.data,
    }),
    submitMaintenanceClaim: build.mutation<
      MaintenanceClaim,
      BillArgs & { key: string; body: ClaimSubmission }
    >({
      query: ({ societyId, billId, key, body }) => ({
        url: `${maintenancePath(societyId)}/my/bills/${billId}/payment-claims`,
        method: "POST",
        headers: { "Idempotency-Key": key },
        body,
      }),
      transformResponse: (r: Envelope<MaintenanceClaim>) => r.data,
      invalidatesTags: (_r, _e, args) => tag(args),
    }),
  }),
});
export const {
  useMaintenanceBillsQuery,
  useLazyMaintenanceBillsQuery,
  useMaintenanceOutstandingQuery,
  useMaintenanceBillQuery,
  useMaintenanceRecordsQuery,
  useCreateMaintenancePaymentRequestMutation,
  useSubmitMaintenanceClaimMutation,
} = maintenanceApi;
