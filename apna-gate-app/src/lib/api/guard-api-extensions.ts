import { enhancedApi } from "@/lib/api/enhanced-api";
import type { FetchBaseQueryError } from "@reduxjs/toolkit/query";
import { loadCheckInReadiness, type CheckInReadiness } from "@/features/guard/check-in-readiness";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPurpose,
  ModelsVisitorStatus,
} from "@/lib/api/generated-api";

type VisitorEntryEvent = "activity" | "checked_in" | "checked_out" | "created" | "expected";

type VisitorEntriesPayload = {
  entries?: ModelsVisitorEntry[];
  total?: number;
  limit?: number;
  offset?: number;
};

type VisitorStatsPayload = {
  today_visitors?: number;
  visitors_inside?: number;
  pending_approvals?: number;
  checked_out_today?: number;
  checked_out_in_range?: number;
  rejected_today?: number;
  auto_closed_today?: number;
};

type ApiEnvelope<T> = {
  data?: T;
  message?: string;
  success?: boolean;
};

export type GuardVisitorEntriesQuery = {
  societyId: number;
  flatId?: number;
  status?: ModelsVisitorStatus;
  source?: ModelsVisitorEntry["source"];
  purpose?: ModelsVisitorPurpose;
  block?: string;
  event?: VisitorEntryEvent;
  eventFrom?: string;
  eventTo?: string;
  createdFrom?: string;
  createdTo?: string;
  search?: string;
  limit?: number;
  offset?: number;
};

export type GuardVisitorStatsQuery = {
  societyId: number;
  eventFrom?: string;
  eventTo?: string;
};

export type GuardDeskVisitorSettings = {
  allow_guard_entry?: boolean;
  allow_guard_on_behalf_approval?: boolean;
};

export type UpdateGuardVisitorEntryBody = {
  companion_details?: Record<string, string>[];
  companions_count?: number;
  email?: string;
  flat_id?: number;
  full_name?: string;
  notes?: string;
  phone_number?: string;
  photo_url?: string;
  vehicle_number?: string;
  vehicle_type?: ModelsVisitorEntry["vehicle_type"];
};

export const guardApiExtensions = enhancedApi.injectEndpoints({
  overrideExisting: true,
  endpoints: (build) => ({
    getGuardCheckInReadiness: build.query<CheckInReadiness, { societyId: number }>({
      async queryFn({ societyId }, _api, _options, baseQuery) {
        try {
          const data = await loadCheckInReadiness(societyId, async (limit, offset) => {
            const response = await baseQuery({
              url: `/v1/societies/${societyId}/visitor-entries/waiting-at-gate`,
              params: { limit, offset },
            });
            if (response.error) throw response.error;
            return (response.data as ApiEnvelope<VisitorEntriesPayload>).data ?? {};
          });
          return { data };
        } catch (error) {
          return {
            error: error instanceof Error
              ? { status: "CUSTOM_ERROR", error: error.message }
              : error as FetchBaseQueryError,
          };
        }
      },
      providesTags: ["Visitor Entries", "VisitorWaitingAtGate"],
    }),
    getV1SocietiesBySocietyIdVisitorEntriesExtended: build.query<
      ApiEnvelope<VisitorEntriesPayload>,
      GuardVisitorEntriesQuery
    >({
      query: ({
        societyId,
        flatId,
        status,
        source,
        purpose,
        block,
        event,
        eventFrom,
        eventTo,
        createdFrom,
        createdTo,
        search,
        limit,
        offset,
      }) => ({
        url: `/v1/societies/${societyId}/visitor-entries`,
        params: {
          flat_id: flatId,
          status,
          source,
          purpose,
          block,
          event,
          event_from: eventFrom,
          event_to: eventTo,
          created_from: createdFrom,
          created_to: createdTo,
          search,
          limit,
          offset,
        },
      }),
      providesTags: ["Visitor Entries"],
    }),
    getV1SocietiesBySocietyIdVisitorEntriesStatsExtended: build.query<
      ApiEnvelope<{ stats?: VisitorStatsPayload }>,
      GuardVisitorStatsQuery
    >({
      query: ({ societyId, eventFrom, eventTo }) => ({
        url: `/v1/societies/${societyId}/visitor-entries/stats`,
        params: {
          event_from: eventFrom,
          event_to: eventTo,
        },
      }),
      providesTags: ["VisitorStats", "Visitor Entries"],
    }),
    patchV1SocietiesBySocietyIdVisitorEntriesAndEntryId: build.mutation<
      ApiEnvelope<{ entry?: ModelsVisitorEntry }>,
      { societyId: number; entryId: number; body: UpdateGuardVisitorEntryBody }
    >({
      query: ({ societyId, entryId, body }) => ({
        url: `/v1/societies/${societyId}/visitor-entries/${entryId}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: ["Visitor Entries"],
    }),
  }),
});

export const {
  useGetGuardCheckInReadinessQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery,
  useLazyGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesStatsExtendedQuery,
  useLazyGetV1SocietiesBySocietyIdVisitorEntriesStatsExtendedQuery,
  usePatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdMutation,
} = guardApiExtensions;
