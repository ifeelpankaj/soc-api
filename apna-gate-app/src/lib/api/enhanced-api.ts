import { generatedApi } from "@/lib/api/generated-api";

export const enhancedApi = generatedApi.enhanceEndpoints({
  addTagTypes: [
    "GuardDesk",
    "VisitorStats",
    "VisitorPending",
    "VisitorWaitingAtGate",
    "VisitorExpectedGuests",
    "FlatVisitorContext",
    "FlatMembers",
    "FlatMemberInvites",
    "VisitorInvites",
    "FlatVisitorEntries",
    "Notifications",
  ],
  endpoints: {
    postV1AuthChangePassword: { invalidatesTags: [] },
    getV1SocietiesBySocietyIdGuardDeskBootstrap: {
      providesTags: ["GuardDesk", "VisitorStats", "VisitorPending", "VisitorWaitingAtGate", "VisitorExpectedGuests"],
    },
    getV1SocietiesBySocietyIdVisitorEntriesStats: {
      providesTags: ["VisitorStats", "Visitor Entries"],
    },
    getV1SocietiesBySocietyIdVisitorEntriesPending: {
      providesTags: ["VisitorPending", "Visitor Entries"],
    },
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPending: {
      providesTags: ["VisitorPending", "Visitor Entries"],
    },
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContext: {
      providesTags: ["FlatVisitorContext", "Visitor Entries"],
    },
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntries: {
      providesTags: ["Visitor Entries", "FlatVisitorEntries"],
    },
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvites: {
      providesTags: ["VisitorInvites"],
    },
    getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteId: {
      providesTags: (_result, _error, arg) => [
        "VisitorInvites",
        { type: "VisitorInvites" as const, id: arg.inviteId },
      ],
    },
    postV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvites: {
      invalidatesTags: ["VisitorInvites"],
    },
    getV1SocietiesBySocietyIdVisitorEntries: {
      providesTags: ["Visitor Entries"],
    },
    getV1SocietiesBySocietyIdVisitorEntriesExpectedGuests: {
      providesTags: ["Visitor Entries", "VisitorExpectedGuests"],
    },
    postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApprove: {
      invalidatesTags: [
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorPending",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
        "FlatVisitorContext",
      ],
    },
    postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdReject: {
      invalidatesTags: [
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorPending",
        "FlatVisitorContext",
      ],
    },
    postV1SocietiesBySocietyIdVisitorEntriesCheckIn: {
      invalidatesTags: [
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
      ],
    },
    postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOut: {
      invalidatesTags: [
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
      ],
    },
    postV1SocietiesBySocietyIdVisitorEntriesGuard: {
      invalidatesTags: [
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorPending",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
      ],
    },
  },
});

export type VisitorNotificationCacheTag =
  | "Visitor Entries"
  | "VisitorStats"
  | "VisitorPending"
  | "VisitorWaitingAtGate"
  | "VisitorExpectedGuests"
  | "GuardDesk"
  | "FlatVisitorContext"
  | "FlatMembers"
  | "FlatMemberInvites"
  | "VisitorInvites"
  | "FlatVisitorEntries";

export function invalidateVisitorNotificationTags(type?: string): VisitorNotificationCacheTag[] {
  switch (type) {
    case "visitor.pending":
    case "visitor.approval_requested":
      return ["VisitorInvites", "VisitorPending", "VisitorStats"];
    case "visitor_invite.accepted":
      return ["VisitorInvites", "FlatVisitorEntries", "Visitor Entries"];
    case "visitor.approved":
      return [
        "VisitorInvites",
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
      ];
    case "visitor.rejected":
      return ["VisitorInvites", "VisitorPending", "VisitorStats", "Visitor Entries"];
    case "visitor.checkin":
    case "visitor.checkout":
      return [
        "VisitorInvites",
        "Visitor Entries",
        "VisitorStats",
        "GuardDesk",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
      ];
    case "member_invite.accepted":
      return ["FlatMemberInvites", "FlatMembers"];
    default:
      return [
        "VisitorInvites",
        "Visitor Entries",
        "VisitorStats",
        "VisitorPending",
        "VisitorWaitingAtGate",
        "VisitorExpectedGuests",
        "GuardDesk",
        "FlatVisitorContext",
      ];
  }
}
