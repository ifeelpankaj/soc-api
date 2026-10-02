import { enhancedApi } from "@/lib/api/enhanced-api";
import type {
  ModelsFlatMemberInviteResponse,
  ModelsFlatResidentResponse,
  ModelsShortLinkResponse,
  ModelsVisitorEntry,
} from "@/lib/api/generated-api";

export type { ModelsFlatMemberInviteResponse } from "@/lib/api/generated-api";

export type ModelsFlatMemberInviteRole = "family" | "tenant";
export type ModelsFlatMemberInviteStatus = "pending" | "accepted" | "expired" | "cancelled";

export type ModelsFlatMemberInviteTokenResponse = {
  token?: string;
  expires_at?: string;
};

export type CreateFlatMemberInviteRequest = {
  role: ModelsFlatMemberInviteRole;
  full_name: string;
  phone?: string;
  email?: string;
};

type FlatResidentsPayload = {
  residents?: ModelsFlatResidentResponse[];
};

type FlatMemberInvitesPayload = {
  invites?: ModelsFlatMemberInviteResponse[];
  total?: number;
  limit?: number;
  offset?: number;
};

type FlatMemberInviteCreatePayload = {
  invite?: ModelsFlatMemberInviteResponse;
  token?: ModelsFlatMemberInviteTokenResponse;
  link?: ModelsShortLinkResponse;
  short_link?: ModelsShortLinkResponse;
};

type VisitorEntryPayload = {
  entry?: ModelsVisitorEntry;
};

type ApiEnvelope<T> = {
  data?: T;
  message?: string;
  success?: boolean;
};

export const residentApiExtensions = enhancedApi.injectEndpoints({
  overrideExisting: true,
  endpoints: (build) => ({
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntries: build.query<
      ApiEnvelope<{
        entries?: ModelsVisitorEntry[];
        total?: number;
        limit?: number;
        offset?: number;
      }>,
      {
        societyId: number;
        flatId: number;
        status?: ModelsVisitorEntry["status"];
        purpose?: ModelsVisitorEntry["purpose"];
        createdFrom?: string;
        createdTo?: string;
        event?: "activity" | "checked_in" | "checked_out" | "created" | "expected";
        eventFrom?: string;
        eventTo?: string;
        limit?: number;
        offset?: number;
        search?: string;
      }
    >({
      query: ({
        societyId,
        flatId,
        status,
        purpose,
        createdFrom,
        createdTo,
        event,
        eventFrom,
        eventTo,
        limit,
        offset,
        search,
      }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/visitor-entries`,
        params: {
          status,
          purpose,
          created_from: createdFrom,
          created_to: createdTo,
          event,
          event_from: eventFrom,
          event_to: eventTo,
          limit,
          offset,
          search,
        },
      }),
      providesTags: ["Visitor Entries", "FlatVisitorEntries"],
    }),
    getV1SocietiesBySocietyIdFlatsAndFlatIdMembers: build.query<
      ApiEnvelope<FlatResidentsPayload>,
      {
        societyId: number;
        flatId: number;
        limit?: number;
        offset?: number;
        search?: string;
        role?: string;
        isPrimary?: boolean;
      }
    >({
      query: ({ societyId, flatId, limit, offset, search, role, isPrimary }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/members`,
        params: { limit, offset, search, role, is_primary: isPrimary },
      }),
      providesTags: ["FlatMembers"],
    }),
    getV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvites: build.query<
      ApiEnvelope<FlatMemberInvitesPayload>,
      {
        societyId: number;
        flatId: number;
        status?: ModelsFlatMemberInviteStatus;
        search?: string;
        limit?: number;
        offset?: number;
      }
    >({
      query: ({ societyId, flatId, status, search, limit, offset }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/member-invites`,
        params: { status, search, limit, offset },
      }),
      providesTags: ["FlatMemberInvites"],
    }),
    getV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteId: build.query<
      ApiEnvelope<{ invite?: ModelsFlatMemberInviteResponse }>,
      { societyId: number; flatId: number; inviteId: number }
    >({
      query: ({ societyId, flatId, inviteId }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/member-invites/${inviteId}`,
      }),
      providesTags: (_result, _error, arg) => [
        "FlatMemberInvites",
        { type: "FlatMemberInvites" as const, id: arg.inviteId },
      ],
    }),
    postV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvites: build.mutation<
      ApiEnvelope<FlatMemberInviteCreatePayload>,
      {
        societyId: number;
        flatId: number;
        modelsCreateFlatMemberInviteRequest: CreateFlatMemberInviteRequest;
      }
    >({
      query: ({ societyId, flatId, modelsCreateFlatMemberInviteRequest }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/member-invites`,
        method: "POST",
        body: modelsCreateFlatMemberInviteRequest,
      }),
      invalidatesTags: ["FlatMemberInvites"],
    }),
    postV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesAndInviteIdCancel: build.mutation<
      ApiEnvelope<unknown>,
      { societyId: number; flatId: number; inviteId: number }
    >({
      query: ({ societyId, flatId, inviteId }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/member-invites/${inviteId}/cancel`,
        method: "POST",
      }),
      invalidatesTags: ["FlatMemberInvites"],
    }),
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryId: build.query<
      ApiEnvelope<VisitorEntryPayload>,
      { societyId: number; flatId: number; entryId: number }
    >({
      query: ({ societyId, flatId, entryId }) => ({
        url: `/v1/societies/${societyId}/flats/${flatId}/visitor-entries/${entryId}`,
      }),
      providesTags: (_result, _error, arg) => [
        { type: "FlatVisitorEntries" as const, id: arg.entryId },
      ],
    }),
  }),
});

export const {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMembersQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery,
  useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesAndInviteIdCancelMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery,
  useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryIdQuery,
} = residentApiExtensions;
