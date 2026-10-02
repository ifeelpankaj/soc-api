import { baseApi as api } from "../../redux/queries/baseApi";
export const addTagTypes = [
  "Health",
  "Auth",
  "Images",
  "Bootstrap",
  "Developer",
  "Flat Claims",
  "Flat Residents",
  "Flats",
  "Notifications",
  "Plans",
  "Flat Members",
  "Short Links",
  "Public",
  "Visitor Entries",
  "Societies",
  "Visitor Settings",
  "Society Members",
  "Subscriptions",
] as const;
const injectedRtkApi = api
  .enhanceEndpoints({
    addTagTypes,
  })
  .injectEndpoints({
    endpoints: (build) => ({
      getHealth: build.query<GetHealthApiResponse, GetHealthApiArg>({
        query: () => ({ url: `/health` }),
        providesTags: ["Health"],
      }),
      getHealthLive: build.query<GetHealthLiveApiResponse, GetHealthLiveApiArg>(
        {
          query: () => ({ url: `/health/live` }),
          providesTags: ["Health"],
        },
      ),
      getHealthReady: build.query<
        GetHealthReadyApiResponse,
        GetHealthReadyApiArg
      >({
        query: () => ({ url: `/health/ready` }),
        providesTags: ["Health"],
      }),
      postV1AuthChangePassword: build.mutation<
        PostV1AuthChangePasswordApiResponse,
        PostV1AuthChangePasswordApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/change-password`,
          method: "POST",
          body: queryArg.modelsChangePasswordRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthForgotPassword: build.mutation<
        PostV1AuthForgotPasswordApiResponse,
        PostV1AuthForgotPasswordApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/forgot-password`,
          method: "POST",
          body: queryArg.modelsForgotPasswordRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthLogin: build.mutation<
        PostV1AuthLoginApiResponse,
        PostV1AuthLoginApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/login`,
          method: "POST",
          body: queryArg.modelsLoginRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthLogout: build.mutation<
        PostV1AuthLogoutApiResponse,
        PostV1AuthLogoutApiArg
      >({
        query: () => ({ url: `/v1/auth/logout`, method: "POST" }),
        invalidatesTags: ["Auth"],
      }),
      getV1AuthProfile: build.query<
        GetV1AuthProfileApiResponse,
        GetV1AuthProfileApiArg
      >({
        query: () => ({ url: `/v1/auth/profile` }),
        providesTags: ["Auth"],
      }),
      patchV1AuthProfile: build.mutation<
        PatchV1AuthProfileApiResponse,
        PatchV1AuthProfileApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/profile`,
          method: "PATCH",
          body: queryArg.modelsUpdateUserRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      getV1AuthProfileAvatar: build.query<
        GetV1AuthProfileAvatarApiResponse,
        GetV1AuthProfileAvatarApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/profile/avatar`,
          params: {
            variant: queryArg.variant,
          },
        }),
        providesTags: ["Images"],
      }),
      putV1AuthProfileAvatar: build.mutation<
        PutV1AuthProfileAvatarApiResponse,
        PutV1AuthProfileAvatarApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/profile/avatar`,
          method: "PUT",
          body: queryArg.body,
        }),
        invalidatesTags: ["Images"],
      }),
      deleteV1AuthProfileAvatar: build.mutation<
        DeleteV1AuthProfileAvatarApiResponse,
        DeleteV1AuthProfileAvatarApiArg
      >({
        query: () => ({ url: `/v1/auth/profile/avatar`, method: "DELETE" }),
        invalidatesTags: ["Images"],
      }),
      postV1AuthRefresh: build.mutation<
        PostV1AuthRefreshApiResponse,
        PostV1AuthRefreshApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/refresh`,
          method: "POST",
          body: queryArg.modelsRefreshTokenRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthRegister: build.mutation<
        PostV1AuthRegisterApiResponse,
        PostV1AuthRegisterApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/register`,
          method: "POST",
          body: queryArg.modelsRegisterRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthResendOtp: build.mutation<
        PostV1AuthResendOtpApiResponse,
        PostV1AuthResendOtpApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/resend-otp`,
          method: "POST",
          body: queryArg.modelsResendOtpRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthResetPassword: build.mutation<
        PostV1AuthResetPasswordApiResponse,
        PostV1AuthResetPasswordApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/reset-password`,
          method: "POST",
          body: queryArg.modelsResetPasswordRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthResidentRegister: build.mutation<
        PostV1AuthResidentRegisterApiResponse,
        PostV1AuthResidentRegisterApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/resident/register`,
          method: "POST",
          body: queryArg.modelsResidentRegisterRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      postV1AuthVerifyOtp: build.mutation<
        PostV1AuthVerifyOtpApiResponse,
        PostV1AuthVerifyOtpApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/auth/verify-otp`,
          method: "POST",
          body: queryArg.modelsVerifyOtpRequest,
        }),
        invalidatesTags: ["Auth"],
      }),
      getV1Bootstrap: build.query<
        GetV1BootstrapApiResponse,
        GetV1BootstrapApiArg
      >({
        query: () => ({ url: `/v1/bootstrap` }),
        providesTags: ["Bootstrap"],
      }),
      getV1DeveloperDashboardBootstrap: build.query<
        GetV1DeveloperDashboardBootstrapApiResponse,
        GetV1DeveloperDashboardBootstrapApiArg
      >({
        query: () => ({ url: `/v1/developer/dashboard/bootstrap` }),
        providesTags: ["Developer"],
      }),
      getV1FlatClaims: build.query<
        GetV1FlatClaimsApiResponse,
        GetV1FlatClaimsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/flat-claims`,
          params: {
            id: queryArg.id,
            society_id: queryArg.societyId,
            flat_id: queryArg.flatId,
            user_id: queryArg.userId,
            status: queryArg.status,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Claims"],
      }),
      postV1FlatClaims: build.mutation<
        PostV1FlatClaimsApiResponse,
        PostV1FlatClaimsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/flat-claims`,
          method: "POST",
          body: queryArg.modelsSubmitFlatClaimRequest,
        }),
        invalidatesTags: ["Flat Claims"],
      }),
      getV1FlatClaimsByClaimId: build.query<
        GetV1FlatClaimsByClaimIdApiResponse,
        GetV1FlatClaimsByClaimIdApiArg
      >({
        query: (queryArg) => ({ url: `/v1/flat-claims/${queryArg.claimId}` }),
        providesTags: ["Flat Claims"],
      }),
      postV1FlatClaimsByClaimIdCancel: build.mutation<
        PostV1FlatClaimsByClaimIdCancelApiResponse,
        PostV1FlatClaimsByClaimIdCancelApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/flat-claims/${queryArg.claimId}/cancel`,
          method: "POST",
        }),
        invalidatesTags: ["Flat Claims"],
      }),
      getV1FlatResidents: build.query<
        GetV1FlatResidentsApiResponse,
        GetV1FlatResidentsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/flat-residents`,
          params: {
            id: queryArg.id,
            society_id: queryArg.societyId,
            flat_id: queryArg.flatId,
            user_id: queryArg.userId,
            role: queryArg.role,
            status: queryArg.status,
            is_primary: queryArg.isPrimary,
            search: queryArg.search,
            search_mode: queryArg.searchMode,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Residents"],
      }),
      getV1Flats: build.query<GetV1FlatsApiResponse, GetV1FlatsApiArg>({
        query: (queryArg) => ({
          url: `/v1/flats`,
          params: {
            id: queryArg.id,
            society_id: queryArg.societyId,
            block: queryArg.block,
            floor: queryArg.floor,
            flat_number: queryArg.flatNumber,
            status: queryArg.status,
            is_active: queryArg.isActive,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flats"],
      }),
      postV1MeDeviceTokens: build.mutation<
        PostV1MeDeviceTokensApiResponse,
        PostV1MeDeviceTokensApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/device-tokens`,
          method: "POST",
          body: queryArg.modelsRegisterDeviceTokenRequest,
        }),
        invalidatesTags: ["Notifications"],
      }),
      deleteV1MeDeviceTokens: build.mutation<
        DeleteV1MeDeviceTokensApiResponse,
        DeleteV1MeDeviceTokensApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/device-tokens`,
          method: "DELETE",
          body: queryArg.modelsUnregisterDeviceTokenRequest,
        }),
        invalidatesTags: ["Notifications"],
      }),
      getV1MeFlatClaims: build.query<
        GetV1MeFlatClaimsApiResponse,
        GetV1MeFlatClaimsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/flat-claims`,
          params: {
            society_id: queryArg.societyId,
            flat_id: queryArg.flatId,
            status: queryArg.status,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Claims"],
      }),
      getV1MeNotifications: build.query<
        GetV1MeNotificationsApiResponse,
        GetV1MeNotificationsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/notifications`,
          params: {
            limit: queryArg.limit,
            cursor: queryArg.cursor,
          },
        }),
        providesTags: ["Notifications"],
      }),
      postV1MeNotificationsReadAll: build.mutation<
        PostV1MeNotificationsReadAllApiResponse,
        PostV1MeNotificationsReadAllApiArg
      >({
        query: () => ({ url: `/v1/me/notifications/read-all`, method: "POST" }),
        invalidatesTags: ["Notifications"],
      }),
      getV1MeNotificationsUnreadCount: build.query<
        GetV1MeNotificationsUnreadCountApiResponse,
        GetV1MeNotificationsUnreadCountApiArg
      >({
        query: () => ({ url: `/v1/me/notifications/unread-count` }),
        providesTags: ["Notifications"],
      }),
      patchV1MeNotificationsByNotificationIdRead: build.mutation<
        PatchV1MeNotificationsByNotificationIdReadApiResponse,
        PatchV1MeNotificationsByNotificationIdReadApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/notifications/${queryArg.notificationId}/read`,
          method: "PATCH",
        }),
        invalidatesTags: ["Notifications"],
      }),
      getV1MeResidences: build.query<
        GetV1MeResidencesApiResponse,
        GetV1MeResidencesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/residences`,
          params: {
            society_id: queryArg.societyId,
            flat_id: queryArg.flatId,
            role: queryArg.role,
            status: queryArg.status,
            is_primary: queryArg.isPrimary,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Residents"],
      }),
      getV1Plans: build.query<GetV1PlansApiResponse, GetV1PlansApiArg>({
        query: (queryArg) => ({
          url: `/v1/plans`,
          params: {
            code: queryArg.code,
            billing_cycle: queryArg.billingCycle,
            is_active: queryArg.isActive,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Plans"],
      }),
      postV1Plans: build.mutation<PostV1PlansApiResponse, PostV1PlansApiArg>({
        query: (queryArg) => ({
          url: `/v1/plans`,
          method: "POST",
          body: queryArg.modelsCreatePlanRequest,
        }),
        invalidatesTags: ["Plans"],
      }),
      getV1PlansLookup: build.query<
        GetV1PlansLookupApiResponse,
        GetV1PlansLookupApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/plans/lookup`,
          params: {
            id: queryArg.id,
            code: queryArg.code,
            name: queryArg.name,
          },
        }),
        providesTags: ["Plans"],
      }),
      patchV1PlansByPlanId: build.mutation<
        PatchV1PlansByPlanIdApiResponse,
        PatchV1PlansByPlanIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/plans/${queryArg.planId}`,
          method: "PATCH",
          body: queryArg.modelsUpdatePlanRequest,
        }),
        invalidatesTags: ["Plans"],
      }),
      postV1PlansByPlanIdActivate: build.mutation<
        PostV1PlansByPlanIdActivateApiResponse,
        PostV1PlansByPlanIdActivateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/plans/${queryArg.planId}/activate`,
          method: "POST",
        }),
        invalidatesTags: ["Plans"],
      }),
      postV1PlansByPlanIdDeactivate: build.mutation<
        PostV1PlansByPlanIdDeactivateApiResponse,
        PostV1PlansByPlanIdDeactivateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/plans/${queryArg.planId}/deactivate`,
          method: "POST",
        }),
        invalidatesTags: ["Plans"],
      }),
      getV1PublicFlatMemberInvitesByToken: build.query<
        GetV1PublicFlatMemberInvitesByTokenApiResponse,
        GetV1PublicFlatMemberInvitesByTokenApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/flat-member-invites/${queryArg.token}`,
        }),
        providesTags: ["Flat Members"],
      }),
      postV1PublicFlatMemberInvitesByTokenAccept: build.mutation<
        PostV1PublicFlatMemberInvitesByTokenAcceptApiResponse,
        PostV1PublicFlatMemberInvitesByTokenAcceptApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/flat-member-invites/${queryArg.token}/accept`,
          method: "POST",
        }),
        invalidatesTags: ["Flat Members"],
      }),
      postV1PublicFlatMemberInvitesByTokenJoin: build.mutation<
        PostV1PublicFlatMemberInvitesByTokenJoinApiResponse,
        PostV1PublicFlatMemberInvitesByTokenJoinApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/flat-member-invites/${queryArg.token}/join`,
          method: "POST",
          body: queryArg.modelsJoinFlatMemberInviteRequest,
        }),
        invalidatesTags: ["Flat Members"],
      }),
      getV1PublicShortLinksByShortCode: build.query<
        GetV1PublicShortLinksByShortCodeApiResponse,
        GetV1PublicShortLinksByShortCodeApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/short-links/${queryArg.shortCode}`,
        }),
        providesTags: ["Short Links"],
      }),
      postV1PublicShortLinksByShortCodeMemberAccept: build.mutation<
        PostV1PublicShortLinksByShortCodeMemberAcceptApiResponse,
        PostV1PublicShortLinksByShortCodeMemberAcceptApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/short-links/${queryArg.shortCode}/member-accept`,
          method: "POST",
        }),
        invalidatesTags: ["Short Links"],
      }),
      postV1PublicShortLinksByShortCodeMemberJoin: build.mutation<
        PostV1PublicShortLinksByShortCodeMemberJoinApiResponse,
        PostV1PublicShortLinksByShortCodeMemberJoinApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/short-links/${queryArg.shortCode}/member-join`,
          method: "POST",
          body: queryArg.modelsJoinFlatMemberInviteRequest,
        }),
        invalidatesTags: ["Short Links"],
      }),
      postV1PublicShortLinksByShortCodeVisitorSubmit: build.mutation<
        PostV1PublicShortLinksByShortCodeVisitorSubmitApiResponse,
        PostV1PublicShortLinksByShortCodeVisitorSubmitApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/short-links/${queryArg.shortCode}/visitor-submit`,
          method: "POST",
          body: queryArg.modelsVisitorFormRequest,
        }),
        invalidatesTags: ["Short Links"],
      }),
      getV1PublicSocietiesBySocietyCodeClaimOptions: build.query<
        GetV1PublicSocietiesBySocietyCodeClaimOptionsApiResponse,
        GetV1PublicSocietiesBySocietyCodeClaimOptionsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/societies/${queryArg.societyCode}/claim-options`,
        }),
        providesTags: ["Public"],
      }),
      postV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQr: build.mutation<
        PostV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQrApiResponse,
        PostV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQrApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/societies/${queryArg.societyCode}/visitor-entries/public-qr`,
          method: "POST",
          body: queryArg.modelsVisitorFormRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLink: build.mutation<
        PostV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLinkApiResponse,
        PostV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLinkApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/societies/${queryArg.societyCode}/visitor-entries/quick-link`,
          method: "POST",
          body: queryArg.modelsVisitorFormRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1PublicSocietiesBySocietyCodeVisitorEntryOptions: build.query<
        GetV1PublicSocietiesBySocietyCodeVisitorEntryOptionsApiResponse,
        GetV1PublicSocietiesBySocietyCodeVisitorEntryOptionsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/societies/${queryArg.societyCode}/visitor-entry-options`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmit:
        build.mutation<
          PostV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmitApiResponse,
          PostV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmitApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/public/societies/${queryArg.societyCode}/visitor-invites/${queryArg.token}/submit`,
            method: "POST",
            body: queryArg.modelsVisitorFormRequest,
          }),
          invalidatesTags: ["Visitor Entries"],
        }),
      postV1PublicVisitorEntriesQrValidate: build.mutation<
        PostV1PublicVisitorEntriesQrValidateApiResponse,
        PostV1PublicVisitorEntriesQrValidateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/visitor-entries/qr/validate`,
          method: "POST",
          body: queryArg.modelsQrTokenRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1PublicVisitorInvitesByToken: build.query<
        GetV1PublicVisitorInvitesByTokenApiResponse,
        GetV1PublicVisitorInvitesByTokenApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/public/visitor-invites/${queryArg.token}`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1Societies: build.query<
        GetV1SocietiesApiResponse,
        GetV1SocietiesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies`,
          params: {
            status: queryArg.status,
            search: queryArg.search,
            search_mode: queryArg.searchMode,
            name: queryArg.name,
            code: queryArg.code,
            city: queryArg.city,
            state: queryArg.state,
            country: queryArg.country,
            pincode: queryArg.pincode,
            created_by: queryArg.createdBy,
            approved_by: queryArg.approvedBy,
            rejected_by: queryArg.rejectedBy,
            suspended_by: queryArg.suspendedBy,
            created_from: queryArg.createdFrom,
            created_to: queryArg.createdTo,
            limit: queryArg.limit,
            offset: queryArg.offset,
            sort_by: queryArg.sortBy,
            sort_order: queryArg.sortOrder,
          },
        }),
        providesTags: ["Societies"],
      }),
      postV1Societies: build.mutation<
        PostV1SocietiesApiResponse,
        PostV1SocietiesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies`,
          method: "POST",
          body: queryArg.modelsCreateSocietyRequest,
        }),
        invalidatesTags: ["Societies"],
      }),
      getV1SocietiesMy: build.query<
        GetV1SocietiesMyApiResponse,
        GetV1SocietiesMyApiArg
      >({
        query: () => ({ url: `/v1/societies/my` }),
        providesTags: ["Societies"],
      }),
      getV1SocietiesBySocietyId: build.query<
        GetV1SocietiesBySocietyIdApiResponse,
        GetV1SocietiesBySocietyIdApiArg
      >({
        query: (queryArg) => ({ url: `/v1/societies/${queryArg.societyId}` }),
        providesTags: ["Societies"],
      }),
      deleteV1SocietiesBySocietyId: build.mutation<
        DeleteV1SocietiesBySocietyIdApiResponse,
        DeleteV1SocietiesBySocietyIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}`,
          method: "DELETE",
        }),
        invalidatesTags: ["Societies"],
      }),
      patchV1SocietiesBySocietyId: build.mutation<
        PatchV1SocietiesBySocietyIdApiResponse,
        PatchV1SocietiesBySocietyIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}`,
          method: "PATCH",
          body: queryArg.modelsUpdateSocietyRequest,
        }),
        invalidatesTags: ["Societies"],
      }),
      getV1SocietiesBySocietyIdAllmember: build.query<
        GetV1SocietiesBySocietyIdAllmemberApiResponse,
        GetV1SocietiesBySocietyIdAllmemberApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/allmember`,
          params: {
            search: queryArg.search,
            role: queryArg.role,
            status: queryArg.status,
            user_id: queryArg.userId,
            invited_by: queryArg.invitedBy,
            removed_by: queryArg.removedBy,
            joined_from: queryArg.joinedFrom,
            joined_to: queryArg.joinedTo,
            limit: queryArg.limit,
            offset: queryArg.offset,
            sort_by: queryArg.sortBy,
            sort_order: queryArg.sortOrder,
          },
        }),
        providesTags: ["Developer"],
      }),
      postV1SocietiesBySocietyIdApprove: build.mutation<
        PostV1SocietiesBySocietyIdApproveApiResponse,
        PostV1SocietiesBySocietyIdApproveApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/approve`,
          method: "POST",
        }),
        invalidatesTags: ["Societies"],
      }),
      getV1SocietiesBySocietyIdDashboardBootstrap: build.query<
        GetV1SocietiesBySocietyIdDashboardBootstrapApiResponse,
        GetV1SocietiesBySocietyIdDashboardBootstrapApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/dashboard/bootstrap`,
        }),
        providesTags: ["Societies"],
      }),
      getV1SocietiesBySocietyIdFlatClaims: build.query<
        GetV1SocietiesBySocietyIdFlatClaimsApiResponse,
        GetV1SocietiesBySocietyIdFlatClaimsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flat-claims`,
          params: {
            flat_id: queryArg.flatId,
            user_id: queryArg.userId,
            status: queryArg.status,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Claims"],
      }),
      getV1SocietiesBySocietyIdFlatClaimsStats: build.query<
        GetV1SocietiesBySocietyIdFlatClaimsStatsApiResponse,
        GetV1SocietiesBySocietyIdFlatClaimsStatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flat-claims/stats`,
        }),
        providesTags: ["Flat Claims"],
      }),
      getV1SocietiesBySocietyIdFlatClaimsAndClaimId: build.query<
        GetV1SocietiesBySocietyIdFlatClaimsAndClaimIdApiResponse,
        GetV1SocietiesBySocietyIdFlatClaimsAndClaimIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flat-claims/${queryArg.claimId}`,
        }),
        providesTags: ["Flat Claims"],
      }),
      postV1SocietiesBySocietyIdFlatClaimsAndClaimIdApprove: build.mutation<
        PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdApproveApiResponse,
        PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdApproveApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flat-claims/${queryArg.claimId}/approve`,
          method: "POST",
        }),
        invalidatesTags: ["Flat Claims"],
      }),
      postV1SocietiesBySocietyIdFlatClaimsAndClaimIdReject: build.mutation<
        PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdRejectApiResponse,
        PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdRejectApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flat-claims/${queryArg.claimId}/reject`,
          method: "POST",
          body: queryArg.modelsRejectFlatClaimRequest,
        }),
        invalidatesTags: ["Flat Claims"],
      }),
      getV1SocietiesBySocietyIdFlats: build.query<
        GetV1SocietiesBySocietyIdFlatsApiResponse,
        GetV1SocietiesBySocietyIdFlatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats`,
          params: {
            block: queryArg.block,
            floor: queryArg.floor,
            flat_number: queryArg.flatNumber,
            status: queryArg.status,
            is_active: queryArg.isActive,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flats"],
      }),
      postV1SocietiesBySocietyIdFlats: build.mutation<
        PostV1SocietiesBySocietyIdFlatsApiResponse,
        PostV1SocietiesBySocietyIdFlatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats`,
          method: "POST",
          body: queryArg.modelsCreateFlatRequest,
        }),
        invalidatesTags: ["Flats"],
      }),
      postV1SocietiesBySocietyIdFlatsBulk: build.mutation<
        PostV1SocietiesBySocietyIdFlatsBulkApiResponse,
        PostV1SocietiesBySocietyIdFlatsBulkApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/bulk`,
          method: "POST",
          body: queryArg.modelsBulkCreateFlatsRequest,
        }),
        invalidatesTags: ["Flats"],
      }),
      postV1SocietiesBySocietyIdFlatsGenerate: build.mutation<
        PostV1SocietiesBySocietyIdFlatsGenerateApiResponse,
        PostV1SocietiesBySocietyIdFlatsGenerateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/generate`,
          method: "POST",
          body: queryArg.modelsGenerateFlatsRequest,
        }),
        invalidatesTags: ["Flats"],
      }),
      getV1SocietiesBySocietyIdFlatsStats: build.query<
        GetV1SocietiesBySocietyIdFlatsStatsApiResponse,
        GetV1SocietiesBySocietyIdFlatsStatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/stats`,
        }),
        providesTags: ["Flats"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatId: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}`,
        }),
        providesTags: ["Flats"],
      }),
      deleteV1SocietiesBySocietyIdFlatsAndFlatId: build.mutation<
        DeleteV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse,
        DeleteV1SocietiesBySocietyIdFlatsAndFlatIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}`,
          method: "DELETE",
        }),
        invalidatesTags: ["Flats"],
      }),
      patchV1SocietiesBySocietyIdFlatsAndFlatId: build.mutation<
        PatchV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse,
        PatchV1SocietiesBySocietyIdFlatsAndFlatIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}`,
          method: "PATCH",
          body: queryArg.modelsUpdateFlatRequest,
        }),
        invalidatesTags: ["Flats"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdBlock: build.mutation<
        PostV1SocietiesBySocietyIdFlatsAndFlatIdBlockApiResponse,
        PostV1SocietiesBySocietyIdFlatsAndFlatIdBlockApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/block`,
          method: "POST",
        }),
        invalidatesTags: ["Flats"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvites: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/member-invites`,
          params: {
            status: queryArg.status,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Members"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvites: build.mutation<
        PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiResponse,
        PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/member-invites`,
          method: "POST",
          body: queryArg.modelsCreateFlatMemberInviteRequest,
        }),
        invalidatesTags: ["Flat Members"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteId: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/member-invites/${queryArg.inviteId}`,
        }),
        providesTags: ["Flat Members"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancel:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancelApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancelApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/member-invites/${queryArg.inviteId}/cancel`,
            method: "POST",
          }),
          invalidatesTags: ["Flat Members"],
        }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdMembers: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMembersApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdMembersApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/members`,
          params: {
            role: queryArg.role,
            is_primary: queryArg.isPrimary,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Members"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdResidents: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents`,
          params: {
            role: queryArg.role,
            status: queryArg.status,
            is_primary: queryArg.isPrimary,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Flat Residents"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserId:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserIdApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserIdApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/users/${queryArg.userId}`,
            method: "POST",
            body: queryArg.modelsAddFlatResidentRequest,
          }),
          invalidatesTags: ["Flat Residents"],
        }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentId: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/${queryArg.residentId}`,
        }),
        providesTags: ["Flat Residents"],
      }),
      deleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentId:
        build.mutation<
          DeleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiResponse,
          DeleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/${queryArg.residentId}`,
            method: "DELETE",
          }),
          invalidatesTags: ["Flat Residents"],
        }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOut:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOutApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOutApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/${queryArg.residentId}/move-out`,
            method: "POST",
          }),
          invalidatesTags: ["Flat Residents"],
        }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimary:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimaryApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimaryApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/${queryArg.residentId}/primary`,
            method: "POST",
          }),
          invalidatesTags: ["Flat Residents"],
        }),
      patchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRole:
        build.mutation<
          PatchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRoleApiResponse,
          PatchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRoleApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/residents/${queryArg.residentId}/role`,
            method: "PATCH",
            body: queryArg.modelsUpdateFlatResidentRoleRequest,
          }),
          invalidatesTags: ["Flat Residents"],
        }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdUnblock: build.mutation<
        PostV1SocietiesBySocietyIdFlatsAndFlatIdUnblockApiResponse,
        PostV1SocietiesBySocietyIdFlatsAndFlatIdUnblockApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/unblock`,
          method: "POST",
        }),
        invalidatesTags: ["Flats"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContext: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-context`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntries: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-entries`,
          params: {
            status: queryArg.status,
            purpose: queryArg.purpose,
            created_from: queryArg.createdFrom,
            created_to: queryArg.createdTo,
            event: queryArg.event,
            event_from: queryArg.eventFrom,
            event_to: queryArg.eventTo,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPending: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-entries/pending`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryId: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-entries/${queryArg.entryId}`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhoto:
        build.query<
          GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhotoApiResponse,
          GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhotoApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-entries/${queryArg.entryId}/photo`,
            params: {
              variant: queryArg.variant,
            },
          }),
          providesTags: ["Images"],
        }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvites: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-invites`,
          params: {
            status: queryArg.status,
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvites: build.mutation<
        PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiResponse,
        PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-invites`,
          method: "POST",
          body: queryArg.modelsCreateVisitorInviteRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaff:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-invites/staff`,
            method: "POST",
            body: queryArg.modelsCreateVisitorInviteRequest,
          }),
          invalidatesTags: ["Visitor Entries"],
        }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteId:
        build.query<
          GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdApiResponse,
          GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-invites/${queryArg.inviteId}`,
          }),
          providesTags: ["Visitor Entries"],
        }),
      getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettings: build.query<
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsApiResponse,
        GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-settings`,
        }),
        providesTags: ["Visitor Settings"],
      }),
      postV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsReset:
        build.mutation<
          PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetApiResponse,
          PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-settings/reset`,
            method: "POST",
          }),
          invalidatesTags: ["Visitor Settings"],
        }),
      patchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurpose:
        build.mutation<
          PatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeApiResponse,
          PatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/flats/${queryArg.flatId}/visitor-settings/${queryArg.purpose}`,
            method: "PATCH",
            body: queryArg.modelsUpdateFlatVisitorSettingRequest,
          }),
          invalidatesTags: ["Visitor Settings"],
        }),
      getV1SocietiesBySocietyIdGuardDeskBootstrap: build.query<
        GetV1SocietiesBySocietyIdGuardDeskBootstrapApiResponse,
        GetV1SocietiesBySocietyIdGuardDeskBootstrapApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/guard-desk/bootstrap`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdGuards: build.mutation<
        PostV1SocietiesBySocietyIdGuardsApiResponse,
        PostV1SocietiesBySocietyIdGuardsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/guards`,
          method: "POST",
          body: queryArg.modelsCreateGuardRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdMembers: build.query<
        GetV1SocietiesBySocietyIdMembersApiResponse,
        GetV1SocietiesBySocietyIdMembersApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members`,
          params: {
            search: queryArg.search,
            role: queryArg.role,
            status: queryArg.status,
            user_id: queryArg.userId,
            invited_by: queryArg.invitedBy,
            removed_by: queryArg.removedBy,
            joined_from: queryArg.joinedFrom,
            joined_to: queryArg.joinedTo,
            limit: queryArg.limit,
            offset: queryArg.offset,
            sort_by: queryArg.sortBy,
            sort_order: queryArg.sortOrder,
          },
        }),
        providesTags: ["Society Members"],
      }),
      postV1SocietiesBySocietyIdMembers: build.mutation<
        PostV1SocietiesBySocietyIdMembersApiResponse,
        PostV1SocietiesBySocietyIdMembersApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members`,
          method: "POST",
          body: queryArg.modelsAddSocietyMemberRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdMembersSummary: build.query<
        GetV1SocietiesBySocietyIdMembersSummaryApiResponse,
        GetV1SocietiesBySocietyIdMembersSummaryApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/summary`,
        }),
        providesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdMembersAndMemberId: build.query<
        GetV1SocietiesBySocietyIdMembersAndMemberIdApiResponse,
        GetV1SocietiesBySocietyIdMembersAndMemberIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/${queryArg.memberId}`,
        }),
        providesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStats:
        build.query<
          GetV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStatsApiResponse,
          GetV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStatsApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/members/${queryArg.memberId}/visitor-approval-stats`,
          }),
          providesTags: ["Visitor Entries"],
        }),
      deleteV1SocietiesBySocietyIdMembersAndUserId: build.mutation<
        DeleteV1SocietiesBySocietyIdMembersAndUserIdApiResponse,
        DeleteV1SocietiesBySocietyIdMembersAndUserIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/${queryArg.userId}`,
          method: "DELETE",
          body: queryArg.modelsSocietyReasonRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      postV1SocietiesBySocietyIdMembersAndUserIdReactivate: build.mutation<
        PostV1SocietiesBySocietyIdMembersAndUserIdReactivateApiResponse,
        PostV1SocietiesBySocietyIdMembersAndUserIdReactivateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/${queryArg.userId}/reactivate`,
          method: "POST",
        }),
        invalidatesTags: ["Society Members"],
      }),
      patchV1SocietiesBySocietyIdMembersAndUserIdRole: build.mutation<
        PatchV1SocietiesBySocietyIdMembersAndUserIdRoleApiResponse,
        PatchV1SocietiesBySocietyIdMembersAndUserIdRoleApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/${queryArg.userId}/role`,
          method: "PATCH",
          body: queryArg.modelsChangeSocietyMemberRoleRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      postV1SocietiesBySocietyIdMembersAndUserIdSuspend: build.mutation<
        PostV1SocietiesBySocietyIdMembersAndUserIdSuspendApiResponse,
        PostV1SocietiesBySocietyIdMembersAndUserIdSuspendApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/members/${queryArg.userId}/suspend`,
          method: "POST",
          body: queryArg.modelsSocietyReasonRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdOnboardingBootstrap: build.query<
        GetV1SocietiesBySocietyIdOnboardingBootstrapApiResponse,
        GetV1SocietiesBySocietyIdOnboardingBootstrapApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/onboarding/bootstrap`,
        }),
        providesTags: ["Societies"],
      }),
      postV1SocietiesBySocietyIdReactivate: build.mutation<
        PostV1SocietiesBySocietyIdReactivateApiResponse,
        PostV1SocietiesBySocietyIdReactivateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/reactivate`,
          method: "POST",
        }),
        invalidatesTags: ["Societies"],
      }),
      postV1SocietiesBySocietyIdReject: build.mutation<
        PostV1SocietiesBySocietyIdRejectApiResponse,
        PostV1SocietiesBySocietyIdRejectApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/reject`,
          method: "POST",
          body: queryArg.modelsSocietyReasonRequest,
        }),
        invalidatesTags: ["Societies"],
      }),
      postV1SocietiesBySocietyIdRestore: build.mutation<
        PostV1SocietiesBySocietyIdRestoreApiResponse,
        PostV1SocietiesBySocietyIdRestoreApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/restore`,
          method: "POST",
        }),
        invalidatesTags: ["Societies"],
      }),
      postV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPending:
        build.mutation<
          PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPendingApiResponse,
          PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPendingApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/subscriptions/plans/${queryArg.planId}/pending`,
            method: "POST",
          }),
          invalidatesTags: ["Subscriptions"],
        }),
      postV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrial:
        build.mutation<
          PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrialApiResponse,
          PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrialApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/subscriptions/plans/${queryArg.planId}/trial`,
            method: "POST",
            body: queryArg.modelsCreateTrialSubscriptionRequest,
          }),
          invalidatesTags: ["Subscriptions"],
        }),
      postV1SocietiesBySocietyIdSuspend: build.mutation<
        PostV1SocietiesBySocietyIdSuspendApiResponse,
        PostV1SocietiesBySocietyIdSuspendApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/suspend`,
          method: "POST",
          body: queryArg.modelsSocietyReasonRequest,
        }),
        invalidatesTags: ["Societies"],
      }),
      postV1SocietiesBySocietyIdTransferOwnership: build.mutation<
        PostV1SocietiesBySocietyIdTransferOwnershipApiResponse,
        PostV1SocietiesBySocietyIdTransferOwnershipApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/transfer-ownership`,
          method: "POST",
          body: queryArg.modelsTransferOwnershipRequest,
        }),
        invalidatesTags: ["Society Members"],
      }),
      getV1SocietiesBySocietyIdVisitorEntries: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries`,
          params: {
            flat_id: queryArg.flatId,
            status: queryArg.status,
            source: queryArg.source,
            purpose: queryArg.purpose,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesCheckIn: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesCheckInApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesCheckInApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/check-in`,
          method: "POST",
          body: queryArg.modelsQrTokenRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesExpectedGuests: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesExpectedGuestsApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesExpectedGuestsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/expected-guests`,
          params: {
            search: queryArg.search,
            event_from: queryArg.eventFrom,
            event_to: queryArg.eventTo,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesGuard: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesGuardApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesGuardApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/guard`,
          method: "POST",
          body: queryArg.modelsVisitorFormRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesPending: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesPendingApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesPendingApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/pending`,
          params: {
            flat_id: queryArg.flatId,
            block: queryArg.block,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesStats: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesStatsApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesStatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/stats`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesStatsDaily: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesStatsDailyApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesStatsDailyApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/stats/daily`,
          params: {
            days: queryArg.days,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesWaitingAtGate: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesWaitingAtGateApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesWaitingAtGateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/waiting-at-gate`,
          params: {
            search: queryArg.search,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesAndEntryId: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      patchV1SocietiesBySocietyIdVisitorEntriesAndEntryId: build.mutation<
        PatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiResponse,
        PatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}`,
          method: "PATCH",
          body: queryArg.modelsUpdateGuardVisitorEntryRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApprove: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/approve`,
          method: "POST",
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckIn:
        build.mutation<
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInApiResponse,
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/approve-and-check-in`,
            method: "POST",
            body: queryArg.modelsGuardApproveEntryRequest,
          }),
          invalidatesTags: ["Visitor Entries"],
        }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckIn: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/check-in`,
          method: "POST",
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOut:
        build.mutation<
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutApiResponse,
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/check-out`,
            method: "POST",
          }),
          invalidatesTags: ["Visitor Entries"],
        }),
      getV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEvents: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEventsApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEventsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/events`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApprove:
        build.mutation<
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveApiResponse,
          PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveApiArg
        >({
          query: (queryArg) => ({
            url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/guard-approve`,
            method: "POST",
            body: queryArg.modelsGuardApproveEntryRequest,
          }),
          invalidatesTags: ["Visitor Entries"],
        }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotify: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/notify`,
          method: "POST",
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhoto: build.query<
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse,
        GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/photo`,
          params: {
            variant: queryArg.variant,
          },
        }),
        providesTags: ["Images"],
      }),
      putV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhoto: build.mutation<
        PutV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse,
        PutV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/photo`,
          method: "PUT",
          body: queryArg.body,
        }),
        invalidatesTags: ["Images"],
      }),
      deleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhoto: build.mutation<
        DeleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse,
        DeleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/photo`,
          method: "DELETE",
        }),
        invalidatesTags: ["Images"],
      }),
      postV1SocietiesBySocietyIdVisitorEntriesAndEntryIdReject: build.mutation<
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectApiResponse,
        PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-entries/${queryArg.entryId}/reject`,
          method: "POST",
          body: queryArg.modelsRejectVisitorEntryRequest,
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      postV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancel: build.mutation<
        PostV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancelApiResponse,
        PostV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancelApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-invites/${queryArg.inviteId}/cancel`,
          method: "POST",
        }),
        invalidatesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorInvitesAndToken: build.query<
        GetV1SocietiesBySocietyIdVisitorInvitesAndTokenApiResponse,
        GetV1SocietiesBySocietyIdVisitorInvitesAndTokenApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-invites/${queryArg.token}`,
        }),
        providesTags: ["Visitor Entries"],
      }),
      getV1SocietiesBySocietyIdVisitorSettings: build.query<
        GetV1SocietiesBySocietyIdVisitorSettingsApiResponse,
        GetV1SocietiesBySocietyIdVisitorSettingsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-settings`,
        }),
        providesTags: ["Visitor Settings"],
      }),
      patchV1SocietiesBySocietyIdVisitorSettings: build.mutation<
        PatchV1SocietiesBySocietyIdVisitorSettingsApiResponse,
        PatchV1SocietiesBySocietyIdVisitorSettingsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-settings`,
          method: "PATCH",
          body: queryArg.modelsUpdateSocietyVisitorSettingsRequest,
        }),
        invalidatesTags: ["Visitor Settings"],
      }),
      getV1SocietiesBySocietyIdVisitorSettingsFlats: build.query<
        GetV1SocietiesBySocietyIdVisitorSettingsFlatsApiResponse,
        GetV1SocietiesBySocietyIdVisitorSettingsFlatsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/societies/${queryArg.societyId}/visitor-settings/flats`,
          params: {
            flat_id: queryArg.flatId,
            block: queryArg.block,
            purpose: queryArg.purpose,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Visitor Settings"],
      }),
      getV1Subscriptions: build.query<
        GetV1SubscriptionsApiResponse,
        GetV1SubscriptionsApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions`,
          params: {
            society_id: queryArg.societyId,
            plan_id: queryArg.planId,
            status: queryArg.status,
            search: queryArg.search,
            search_mode: queryArg.searchMode,
            is_active_only: queryArg.isActiveOnly,
            expired_only: queryArg.expiredOnly,
            expiring_before: queryArg.expiringBefore,
            expiring_days: queryArg.expiringDays,
            limit: queryArg.limit,
            offset: queryArg.offset,
          },
        }),
        providesTags: ["Subscriptions"],
      }),
      getV1SubscriptionsLookup: build.query<
        GetV1SubscriptionsLookupApiResponse,
        GetV1SubscriptionsLookupApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/lookup`,
          params: {
            id: queryArg.id,
            society_id: queryArg.societyId,
          },
        }),
        providesTags: ["Subscriptions"],
      }),
      getV1SubscriptionsStats: build.query<
        GetV1SubscriptionsStatsApiResponse,
        GetV1SubscriptionsStatsApiArg
      >({
        query: () => ({ url: `/v1/subscriptions/stats` }),
        providesTags: ["Subscriptions"],
      }),
      postV1SubscriptionsBySubscriptionIdActivate: build.mutation<
        PostV1SubscriptionsBySubscriptionIdActivateApiResponse,
        PostV1SubscriptionsBySubscriptionIdActivateApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/${queryArg.subscriptionId}/activate`,
          method: "POST",
          body: queryArg.modelsActivateSubscriptionRequest,
        }),
        invalidatesTags: ["Subscriptions"],
      }),
      postV1SubscriptionsBySubscriptionIdCancel: build.mutation<
        PostV1SubscriptionsBySubscriptionIdCancelApiResponse,
        PostV1SubscriptionsBySubscriptionIdCancelApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/${queryArg.subscriptionId}/cancel`,
          method: "POST",
          body: queryArg.modelsCancelSubscriptionRequest,
        }),
        invalidatesTags: ["Subscriptions"],
      }),
      postV1SubscriptionsBySubscriptionIdExpire: build.mutation<
        PostV1SubscriptionsBySubscriptionIdExpireApiResponse,
        PostV1SubscriptionsBySubscriptionIdExpireApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/${queryArg.subscriptionId}/expire`,
          method: "POST",
        }),
        invalidatesTags: ["Subscriptions"],
      }),
      postV1SubscriptionsBySubscriptionIdPlansAndPlanId: build.mutation<
        PostV1SubscriptionsBySubscriptionIdPlansAndPlanIdApiResponse,
        PostV1SubscriptionsBySubscriptionIdPlansAndPlanIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/${queryArg.subscriptionId}/plans/${queryArg.planId}`,
          method: "POST",
        }),
        invalidatesTags: ["Subscriptions"],
      }),
      postV1SubscriptionsBySubscriptionIdRenew: build.mutation<
        PostV1SubscriptionsBySubscriptionIdRenewApiResponse,
        PostV1SubscriptionsBySubscriptionIdRenewApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/subscriptions/${queryArg.subscriptionId}/renew`,
          method: "POST",
          body: queryArg.modelsRenewSubscriptionRequest,
        }),
        invalidatesTags: ["Subscriptions"],
      }),
      getV1MeNotificationsByNotificationId: build.query<
        GetV1MeNotificationsByNotificationIdApiResponse,
        GetV1MeNotificationsByNotificationIdApiArg
      >({
        query: (queryArg) => ({
          url: `/v1/me/notifications/${queryArg.notificationId}`,
        }),
        providesTags: ["Notifications"],
      }),
    }),
    overrideExisting: true,
  });
export { injectedRtkApi as generatedApi };
export type GetHealthApiResponse =
  /** status 200 API server is healthy */ ModelsHealthCheckResponseDoc;
export type GetHealthApiArg = void;
export type GetHealthLiveApiResponse =
  /** status 200 API process is alive */ ModelsLivenessResponseDoc;
export type GetHealthLiveApiArg = void;
export type GetHealthReadyApiResponse =
  /** status 200 API and database are ready */ ModelsReadinessResponseDoc;
export type GetHealthReadyApiArg = void;
export type PostV1AuthChangePasswordApiResponse =
  /** status 200 Password changed successfully */ ModelsChangePasswordApiResponse;
export type PostV1AuthChangePasswordApiArg = {
  /** Change password payload */
  modelsChangePasswordRequest: ModelsChangePasswordRequest;
};
export type PostV1AuthForgotPasswordApiResponse =
  /** status 200 Password reset instructions response */ ModelsForgotPasswordApiResponse;
export type PostV1AuthForgotPasswordApiArg = {
  /** Forgot password payload */
  modelsForgotPasswordRequest: ModelsForgotPasswordRequest;
};
export type PostV1AuthLoginApiResponse =
  /** status 200 Login successful */ ModelsLoginApiResponse;
export type PostV1AuthLoginApiArg = {
  /** Login payload */
  modelsLoginRequest: ModelsLoginRequest;
};
export type PostV1AuthLogoutApiResponse =
  /** status 200 Logout successful */ ModelsLogoutApiResponse;
export type PostV1AuthLogoutApiArg = void;
export type GetV1AuthProfileApiResponse =
  /** status 200 Profile fetched successfully */ ModelsGetProfileApiResponse;
export type GetV1AuthProfileApiArg = void;
export type PatchV1AuthProfileApiResponse =
  /** status 200 Profile updated successfully */ ModelsGetProfileApiResponse;
export type PatchV1AuthProfileApiArg = {
  /** Profile update payload */
  modelsUpdateUserRequest: ModelsUpdateUserRequest;
};
export type GetV1AuthProfileAvatarApiResponse =
  /** status 200 OK */ ModelsImageViewApiResponse;
export type GetV1AuthProfileAvatarApiArg = {
  /** Server-selected image view */
  variant?: "list" | "avatar" | "detail" | "original";
};
export type PutV1AuthProfileAvatarApiResponse =
  /** status 200 OK */ ModelsImageViewApiResponse;
export type PutV1AuthProfileAvatarApiArg = {
  body: {
    /** JPEG or PNG, maximum 5 MiB and 25 megapixels */
    file: Blob;
  };
};
export type DeleteV1AuthProfileAvatarApiResponse = unknown;
export type DeleteV1AuthProfileAvatarApiArg = void;
export type PostV1AuthRefreshApiResponse =
  /** status 200 OK */ ModelsRefreshTokenApiResponse;
export type PostV1AuthRefreshApiArg = {
  /** Refresh token payload */
  modelsRefreshTokenRequest: ModelsRefreshTokenRequest;
};
export type PostV1AuthRegisterApiResponse =
  /** status 201 Account created successfully */ ModelsRegisterApiResponse;
export type PostV1AuthRegisterApiArg = {
  /** Registration payload */
  modelsRegisterRequest: ModelsRegisterRequest;
};
export type PostV1AuthResendOtpApiResponse =
  /** status 200 OTP sent successfully */ ModelsResendOtpapiResponse;
export type PostV1AuthResendOtpApiArg = {
  /** Resend OTP payload */
  modelsResendOtpRequest: ModelsResendOtpRequest;
};
export type PostV1AuthResetPasswordApiResponse =
  /** status 200 Password reset successfully */ ModelsResetPasswordApiResponse;
export type PostV1AuthResetPasswordApiArg = {
  /** Reset password payload */
  modelsResetPasswordRequest: ModelsResetPasswordRequest;
};
export type PostV1AuthResidentRegisterApiResponse =
  /** status 201 Resident account created successfully */ ModelsLoginApiResponse;
export type PostV1AuthResidentRegisterApiArg = {
  /** Resident registration payload */
  modelsResidentRegisterRequest: ModelsResidentRegisterRequest;
};
export type PostV1AuthVerifyOtpApiResponse =
  /** status 200 Email verified successfully */ ModelsVerifyOtpapiResponse;
export type PostV1AuthVerifyOtpApiArg = {
  /** Email verification payload */
  modelsVerifyOtpRequest: ModelsVerifyOtpRequest;
};
export type GetV1BootstrapApiResponse =
  /** status 200 Bootstrap fetched successfully */ ModelsBootstrapApiResponse;
export type GetV1BootstrapApiArg = void;
export type GetV1DeveloperDashboardBootstrapApiResponse =
  /** status 200 Developer dashboard bootstrap fetched successfully */ ModelsDeveloperDashboardBootstrapApiResponse;
export type GetV1DeveloperDashboardBootstrapApiArg = void;
export type GetV1FlatClaimsApiResponse =
  /** status 200 Flat claims fetched successfully */ ModelsFlatClaimsApiResponse;
export type GetV1FlatClaimsApiArg = {
  /** Claim ID */
  id?: number;
  /** Society ID */
  societyId?: number;
  /** Flat ID */
  flatId?: number;
  /** User ID */
  userId?: number;
  /** Claim status: pending, approved, rejected, cancelled */
  status?: string;
  /** Search user, contact, flat, block, or status */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type PostV1FlatClaimsApiResponse =
  /** status 201 Flat claim submitted successfully */ ModelsFlatClaimApiResponse;
export type PostV1FlatClaimsApiArg = {
  /** Submit flat claim payload */
  modelsSubmitFlatClaimRequest: ModelsSubmitFlatClaimRequest;
};
export type GetV1FlatClaimsByClaimIdApiResponse =
  /** status 200 Flat claim fetched successfully */ ModelsFlatClaimApiResponse;
export type GetV1FlatClaimsByClaimIdApiArg = {
  /** Claim ID */
  claimId: number;
};
export type PostV1FlatClaimsByClaimIdCancelApiResponse =
  /** status 200 Flat claim cancelled successfully */ ModelsFlatClaimApiResponse;
export type PostV1FlatClaimsByClaimIdCancelApiArg = {
  /** Claim ID */
  claimId: number;
};
export type GetV1FlatResidentsApiResponse =
  /** status 200 Residents fetched successfully */ ModelsFlatResidentsApiResponse;
export type GetV1FlatResidentsApiArg = {
  /** Resident ID */
  id?: number;
  /** Society ID */
  societyId?: number;
  /** Flat ID */
  flatId?: number;
  /** User ID */
  userId?: number;
  /** Resident role: owner, tenant, family */
  role?: string;
  /** Resident status: active, inactive, moved_out */
  status?: string;
  /** Primary resident flag */
  isPrimary?: boolean;
  /** Search user, contact, flat, block, role, or status */
  search?: string;
  /** Search mode */
  searchMode?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type GetV1FlatsApiResponse =
  /** status 200 Flats fetched successfully */ ModelsFlatsApiResponse;
export type GetV1FlatsApiArg = {
  /** Flat ID */
  id?: number;
  /** Society ID */
  societyId?: number;
  /** Block */
  block?: string;
  /** Floor */
  floor?: string;
  /** Flat number */
  flatNumber?: string;
  /** Flat status: vacant, occupied, blocked */
  status?: string;
  /** Active state */
  isActive?: boolean;
  /** Search flat number, block, floor, status, society name/code */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type PostV1MeDeviceTokensApiResponse =
  /** status 200 Device token registered successfully */ ModelsRegisterDeviceTokenApiResponse;
export type PostV1MeDeviceTokensApiArg = {
  /** Device token payload */
  modelsRegisterDeviceTokenRequest: ModelsRegisterDeviceTokenRequest;
};
export type DeleteV1MeDeviceTokensApiResponse =
  /** status 200 Device token removed successfully */ ModelsMessageApiResponse;
export type DeleteV1MeDeviceTokensApiArg = {
  /** Device token payload */
  modelsUnregisterDeviceTokenRequest: ModelsUnregisterDeviceTokenRequest;
};
export type GetV1MeFlatClaimsApiResponse =
  /** status 200 My flat claims fetched successfully */ ModelsFlatClaimsApiResponse;
export type GetV1MeFlatClaimsApiArg = {
  /** Society ID */
  societyId?: number;
  /** Flat ID */
  flatId?: number;
  /** Claim status */
  status?: string;
  /** Search text */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type GetV1MeNotificationsApiResponse =
  /** status 200 Notifications fetched successfully */ ModelsNotificationsApiResponse;
export type GetV1MeNotificationsApiArg = {
  /** Page size */
  limit?: number;
  /** Opaque pagination cursor */
  cursor?: string;
};
export type PostV1MeNotificationsReadAllApiResponse =
  /** status 200 Notifications marked read */ ModelsNotificationUnreadCountApiResponse;
export type PostV1MeNotificationsReadAllApiArg = void;
export type GetV1MeNotificationsUnreadCountApiResponse =
  /** status 200 Unread count fetched successfully */ ModelsNotificationUnreadCountApiResponse;
export type GetV1MeNotificationsUnreadCountApiArg = void;
export type PatchV1MeNotificationsByNotificationIdReadApiResponse =
  /** status 200 Notification marked read */ ModelsNotificationReadApiResponse;
export type PatchV1MeNotificationsByNotificationIdReadApiArg = {
  /** Notification ID */
  notificationId: string;
};
export type GetV1MeResidencesApiResponse =
  /** status 200 My residences fetched successfully */ ModelsMyResidencesApiResponse;
export type GetV1MeResidencesApiArg = {
  /** Society ID */
  societyId?: number;
  /** Flat ID */
  flatId?: number;
  /** Resident role */
  role?: string;
  /** Resident status */
  status?: string;
  /** Primary resident flag */
  isPrimary?: boolean;
  /** Search text */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type GetV1PlansApiResponse =
  /** status 200 Plans fetched successfully */ ModelsPlansApiResponse;
export type GetV1PlansApiArg = {
  /** Plan code */
  code?: string;
  /** Billing cycle */
  billingCycle?: string;
  /** Active state */
  isActive?: boolean;
  /** Search text */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type PostV1PlansApiResponse =
  /** status 201 Plan created successfully */ ModelsPlanApiResponse;
export type PostV1PlansApiArg = {
  /** Create plan payload */
  modelsCreatePlanRequest: ModelsCreatePlanRequest;
};
export type GetV1PlansLookupApiResponse =
  /** status 200 Plan fetched successfully */ ModelsPlanApiResponse;
export type GetV1PlansLookupApiArg = {
  /** Plan ID */
  id?: number;
  /** Plan code */
  code?: string;
  /** Plan name */
  name?: string;
};
export type PatchV1PlansByPlanIdApiResponse =
  /** status 200 Plan updated successfully */ ModelsPlanApiResponse;
export type PatchV1PlansByPlanIdApiArg = {
  /** Plan ID */
  planId: number;
  /** Update plan payload */
  modelsUpdatePlanRequest: ModelsUpdatePlanRequest;
};
export type PostV1PlansByPlanIdActivateApiResponse =
  /** status 200 Plan activated successfully */ ModelsPlanApiResponse;
export type PostV1PlansByPlanIdActivateApiArg = {
  /** Plan ID */
  planId: number;
};
export type PostV1PlansByPlanIdDeactivateApiResponse =
  /** status 200 Plan deactivated successfully */ ModelsPlanApiResponse;
export type PostV1PlansByPlanIdDeactivateApiArg = {
  /** Plan ID */
  planId: number;
};
export type GetV1PublicFlatMemberInvitesByTokenApiResponse =
  /** status 200 Member invite fetched successfully */ ModelsPublicFlatMemberInviteApiResponse;
export type GetV1PublicFlatMemberInvitesByTokenApiArg = {
  /** Member invite token */
  token: string;
};
export type PostV1PublicFlatMemberInvitesByTokenAcceptApiResponse =
  /** status 200 Member invite accepted successfully */ ModelsAcceptFlatMemberInviteApiResponse;
export type PostV1PublicFlatMemberInvitesByTokenAcceptApiArg = {
  /** Member invite token */
  token: string;
};
export type PostV1PublicFlatMemberInvitesByTokenJoinApiResponse =
  /** status 200 Member joined successfully */ ModelsJoinFlatMemberInviteApiResponse;
export type PostV1PublicFlatMemberInvitesByTokenJoinApiArg = {
  /** Member invite token */
  token: string;
  /** Join request */
  modelsJoinFlatMemberInviteRequest: ModelsJoinFlatMemberInviteRequest;
};
export type GetV1PublicShortLinksByShortCodeApiResponse =
  /** status 200 Short link resolved successfully */ ModelsShortLinkResolveApiResponse;
export type GetV1PublicShortLinksByShortCodeApiArg = {
  /** Short link code */
  shortCode: string;
};
export type PostV1PublicShortLinksByShortCodeMemberAcceptApiResponse =
  /** status 200 Member invite accepted successfully */ ModelsAcceptFlatMemberInviteApiResponse;
export type PostV1PublicShortLinksByShortCodeMemberAcceptApiArg = {
  /** Short link code */
  shortCode: string;
};
export type PostV1PublicShortLinksByShortCodeMemberJoinApiResponse =
  /** status 200 Member joined successfully */ ModelsJoinFlatMemberInviteApiResponse;
export type PostV1PublicShortLinksByShortCodeMemberJoinApiArg = {
  /** Short link code */
  shortCode: string;
  /** Join request */
  modelsJoinFlatMemberInviteRequest: ModelsJoinFlatMemberInviteRequest;
};
export type PostV1PublicShortLinksByShortCodeVisitorSubmitApiResponse =
  /** status 201 Visitor entry created successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1PublicShortLinksByShortCodeVisitorSubmitApiArg = {
  /** Short link code */
  shortCode: string;
  /** Visitor details */
  modelsVisitorFormRequest: ModelsVisitorFormRequest;
};
export type GetV1PublicSocietiesBySocietyCodeClaimOptionsApiResponse =
  /** status 200 Claim options fetched successfully */ ModelsPublicClaimOptionsApiResponse;
export type GetV1PublicSocietiesBySocietyCodeClaimOptionsApiArg = {
  /** Society code */
  societyCode: string;
};
export type PostV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQrApiResponse =
  /** status 201 Visitor entry created successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQrApiArg = {
  /** Society code */
  societyCode: string;
  /** Visitor details */
  modelsVisitorFormRequest: ModelsVisitorFormRequest;
};
export type PostV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLinkApiResponse =
  /** status 201 Visitor entry created successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLinkApiArg = {
  /** Society code */
  societyCode: string;
  /** Visitor details */
  modelsVisitorFormRequest: ModelsVisitorFormRequest;
};
export type GetV1PublicSocietiesBySocietyCodeVisitorEntryOptionsApiResponse =
  /** status 200 Visitor entry options fetched successfully */ ModelsVisitorEntryOptionsApiResponse;
export type GetV1PublicSocietiesBySocietyCodeVisitorEntryOptionsApiArg = {
  /** Society code */
  societyCode: string;
};
export type PostV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmitApiResponse =
  /** status 201 Visitor entry created successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmitApiArg =
  {
    /** Society code */
    societyCode: string;
    /** Visitor invite token */
    token: string;
    /** Visitor details */
    modelsVisitorFormRequest: ModelsVisitorFormRequest;
  };
export type PostV1PublicVisitorEntriesQrValidateApiResponse =
  /** status 200 Visitor QR validated successfully */ ModelsVisitorEntryApiResponse;
export type PostV1PublicVisitorEntriesQrValidateApiArg = {
  /** Visitor QR token */
  modelsQrTokenRequest: ModelsQrTokenRequest;
};
export type GetV1PublicVisitorInvitesByTokenApiResponse =
  /** status 200 Visitor invite fetched successfully */ ModelsVisitorInviteApiResponse;
export type GetV1PublicVisitorInvitesByTokenApiArg = {
  /** Visitor invite token */
  token: string;
};
export type GetV1SocietiesApiResponse =
  /** status 200 Societies fetched successfully */ ModelsPaginatedSocietiesApiResponse;
export type GetV1SocietiesApiArg = {
  /** Society status */
  status?: string;
  /** Search text */
  search?: string;
  /** Search mode */
  searchMode?: string;
  /** Society name */
  name?: string;
  /** Society code */
  code?: string;
  /** City */
  city?: string;
  /** State */
  state?: string;
  /** Country */
  country?: string;
  /** Pincode */
  pincode?: string;
  /** Created by user ID */
  createdBy?: number;
  /** Approved by user ID */
  approvedBy?: number;
  /** Rejected by user ID */
  rejectedBy?: number;
  /** Suspended by user ID */
  suspendedBy?: number;
  /** Created from RFC3339 timestamp */
  createdFrom?: string;
  /** Created to RFC3339 timestamp */
  createdTo?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
  /** Sort by: created_at, updated_at, name, city, status */
  sortBy?: string;
  /** Sort order: asc, desc */
  sortOrder?: string;
};
export type PostV1SocietiesApiResponse =
  /** status 201 Society request created successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesApiArg = {
  /** Create society payload */
  modelsCreateSocietyRequest: ModelsCreateSocietyRequest;
};
export type GetV1SocietiesMyApiResponse =
  /** status 200 My societies fetched successfully */ ModelsMySocietiesApiResponse;
export type GetV1SocietiesMyApiArg = void;
export type GetV1SocietiesBySocietyIdApiResponse =
  /** status 200 Society fetched successfully */ ModelsSocietyDetailApiResponse;
export type GetV1SocietiesBySocietyIdApiArg = {
  /** Society ID */
  societyId: number;
};
export type DeleteV1SocietiesBySocietyIdApiResponse =
  /** status 200 Society deleted successfully */ ModelsMessageApiResponse;
export type DeleteV1SocietiesBySocietyIdApiArg = {
  /** Society ID */
  societyId: number;
};
export type PatchV1SocietiesBySocietyIdApiResponse =
  /** status 200 Society updated successfully */ ModelsSocietyApiResponse;
export type PatchV1SocietiesBySocietyIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Update society payload */
  modelsUpdateSocietyRequest: ModelsUpdateSocietyRequest;
};
export type GetV1SocietiesBySocietyIdAllmemberApiResponse =
  /** status 200 Members fetched successfully */ ModelsPaginatedMembersApiResponse;
export type GetV1SocietiesBySocietyIdAllmemberApiArg = {
  /** Society ID */
  societyId: number;
  /** Search full name, email, phone, role, or status */
  search?: string;
  /** Member role */
  role?: string;
  /** Member status */
  status?: string;
  /** User ID */
  userId?: number;
  /** Invited by user ID */
  invitedBy?: number;
  /** Removed by user ID */
  removedBy?: number;
  /** Joined from RFC3339 timestamp */
  joinedFrom?: string;
  /** Joined to RFC3339 timestamp */
  joinedTo?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
  /** Sort by: joined_at, role, status */
  sortBy?: string;
  /** Sort order: asc, desc */
  sortOrder?: string;
};
export type PostV1SocietiesBySocietyIdApproveApiResponse =
  /** status 200 Society approved successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesBySocietyIdApproveApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdDashboardBootstrapApiResponse =
  /** status 200 Society dashboard bootstrap fetched successfully */ ModelsSocietyDashboardBootstrapApiResponse;
export type GetV1SocietiesBySocietyIdDashboardBootstrapApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdFlatClaimsApiResponse =
  /** status 200 Flat claims fetched successfully */ ModelsFlatClaimsApiResponse;
export type GetV1SocietiesBySocietyIdFlatClaimsApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId?: number;
  /** User ID */
  userId?: number;
  /** Claim status: pending, approved, rejected, cancelled */
  status?: string;
  /** Search user, contact, flat, block, or status */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type GetV1SocietiesBySocietyIdFlatClaimsStatsApiResponse =
  /** status 200 Flat claim stats fetched successfully */ ModelsFlatClaimStatsApiResponse;
export type GetV1SocietiesBySocietyIdFlatClaimsStatsApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdFlatClaimsAndClaimIdApiResponse =
  /** status 200 Flat claim fetched successfully */ ModelsFlatClaimApiResponse;
export type GetV1SocietiesBySocietyIdFlatClaimsAndClaimIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Claim ID */
  claimId: number;
};
export type PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdApproveApiResponse =
  /** status 200 Flat claim approved successfully */ ModelsFlatApprovalApiResponse;
export type PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdApproveApiArg = {
  /** Society ID */
  societyId: number;
  /** Claim ID */
  claimId: number;
};
export type PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdRejectApiResponse =
  /** status 200 Flat claim rejected successfully */ ModelsFlatClaimApiResponse;
export type PostV1SocietiesBySocietyIdFlatClaimsAndClaimIdRejectApiArg = {
  /** Society ID */
  societyId: number;
  /** Claim ID */
  claimId: number;
  /** Reject flat claim payload */
  modelsRejectFlatClaimRequest: ModelsRejectFlatClaimRequest;
};
export type GetV1SocietiesBySocietyIdFlatsApiResponse =
  /** status 200 Flats fetched successfully */ ModelsPaginatedFlatsApiResponse;
export type GetV1SocietiesBySocietyIdFlatsApiArg = {
  /** Society ID */
  societyId: number;
  /** Block */
  block?: string;
  /** Floor */
  floor?: string;
  /** Flat number */
  flatNumber?: string;
  /** Flat status: vacant, occupied, blocked */
  status?: string;
  /** Active state */
  isActive?: boolean;
  /** Search flat number, block, floor, or status */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdFlatsApiResponse =
  /** status 201 Flat created successfully */ ModelsFlatApiResponse;
export type PostV1SocietiesBySocietyIdFlatsApiArg = {
  /** Society ID */
  societyId: number;
  /** Create flat payload */
  modelsCreateFlatRequest: ModelsCreateFlatRequest;
};
export type PostV1SocietiesBySocietyIdFlatsBulkApiResponse =
  /** status 201 Flats created successfully */ ModelsBulkFlatsApiResponse;
export type PostV1SocietiesBySocietyIdFlatsBulkApiArg = {
  /** Society ID */
  societyId: number;
  /** Bulk create flats payload */
  modelsBulkCreateFlatsRequest: ModelsBulkCreateFlatsRequest;
};
export type PostV1SocietiesBySocietyIdFlatsGenerateApiResponse =
  /** status 201 Flats generated and created successfully */ ModelsBulkFlatsApiResponse;
export type PostV1SocietiesBySocietyIdFlatsGenerateApiArg = {
  /** Society ID */
  societyId: number;
  /** Generate flats payload */
  modelsGenerateFlatsRequest: ModelsGenerateFlatsRequest;
};
export type GetV1SocietiesBySocietyIdFlatsStatsApiResponse =
  /** status 200 Flat stats fetched successfully */ ModelsFlatStatsApiResponse;
export type GetV1SocietiesBySocietyIdFlatsStatsApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse =
  /** status 200 Flat fetched successfully */ ModelsFlatApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type DeleteV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse =
  /** status 200 Flat deleted successfully */ ModelsMessageApiResponse;
export type DeleteV1SocietiesBySocietyIdFlatsAndFlatIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdApiResponse =
  /** status 200 Flat updated successfully */ ModelsFlatApiResponse;
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Update flat payload */
  modelsUpdateFlatRequest: ModelsUpdateFlatRequest;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdBlockApiResponse =
  /** status 200 Flat blocked successfully */ ModelsFlatApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdBlockApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiResponse =
  /** status 200 Member invites fetched successfully */ ModelsFlatMemberInvitesApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Member invite status */
  status?: "pending" | "accepted" | "expired" | "cancelled";
  /** Search name, contact, role, status, or short code */
  search?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiResponse =
  /** status 201 Member invite created successfully */ ModelsFlatMemberInviteTokenApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Member invite request */
  modelsCreateFlatMemberInviteRequest: ModelsCreateFlatMemberInviteRequest;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdApiResponse =
  /** status 200 Member invite fetched successfully */ ModelsFlatMemberInviteApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Member invite ID */
    inviteId: number;
  };
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancelApiResponse =
  /** status 200 Member invite cancelled successfully */ ModelsMessageApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancelApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Member invite ID */
    inviteId: number;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMembersApiResponse =
  /** status 200 Residents fetched successfully */ ModelsFlatResidentsApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdMembersApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Resident role: owner, tenant, family */
  role?: string;
  /** Primary resident flag */
  isPrimary?: boolean;
  /** Search user, contact, flat, block, role, or status */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsApiResponse =
  /** status 200 Residents fetched successfully */ ModelsFlatResidentsApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Resident role: owner, tenant, family */
  role?: string;
  /** Resident status: active, inactive, moved_out */
  status?: string;
  /** Primary resident flag */
  isPrimary?: boolean;
  /** Search user, contact, flat, block, role, or status */
  search?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserIdApiResponse =
  /** status 201 Resident added successfully */ ModelsFlatResidentApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserIdApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** User ID */
    userId: number;
    /** Add flat resident payload */
    modelsAddFlatResidentRequest: ModelsAddFlatResidentRequest;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiResponse =
  /** status 200 Resident fetched successfully */ ModelsFlatResidentApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Resident ID */
  residentId: number;
};
export type DeleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiResponse =
  /** status 200 Resident removed successfully */ ModelsMessageApiResponse;
export type DeleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Resident ID */
    residentId: number;
  };
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOutApiResponse =
  /** status 200 Resident moved out successfully */ ModelsFlatResidentApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOutApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Resident ID */
    residentId: number;
  };
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimaryApiResponse =
  /** status 200 Primary resident changed successfully */ ModelsFlatResidentApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimaryApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Resident ID */
    residentId: number;
  };
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRoleApiResponse =
  /** status 200 Resident role updated successfully */ ModelsFlatResidentApiResponse;
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRoleApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Resident ID */
    residentId: number;
    /** Update flat resident role payload */
    modelsUpdateFlatResidentRoleRequest: ModelsUpdateFlatResidentRoleRequest;
  };
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdUnblockApiResponse =
  /** status 200 Flat unblocked successfully */ ModelsFlatApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdUnblockApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextApiResponse =
  /** status 200 Flat visitor context fetched successfully */ ModelsFlatVisitorContextApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesApiResponse =
  /** status 200 Visitor entries fetched successfully */ ModelsVisitorEntriesApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Visitor status */
  status?:
    | "waiting_approval"
    | "approved"
    | "rejected"
    | "checked_in"
    | "checked_out"
    | "cancelled"
    | "expired"
    | "auto_closed";
  /** Visitor purpose */
  purpose?:
    | "guest"
    | "delivery"
    | "cab"
    | "service"
    | "maintenance"
    | "staff"
    | "other";
  /** Created from (RFC3339) */
  createdFrom?: string;
  /** Created to (RFC3339) */
  createdTo?: string;
  /** Lifecycle timestamp to filter by */
  event?: "created" | "checked_in" | "checked_out" | "expected" | "activity";
  /** Lifecycle filter start (RFC3339) */
  eventFrom?: string;
  /** Lifecycle filter end (RFC3339) */
  eventTo?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingApiResponse =
  /** status 200 Pending visitor approvals fetched successfully */ ModelsVisitorEntriesApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdApiResponse =
  /** status 200 Visitor entry fetched successfully */ ModelsVisitorEntryApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Visitor entry ID */
    entryId: number;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhotoApiResponse =
  /** status 200 OK */ ModelsImageViewApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhotoApiArg =
  {
    /** Server-selected image view */
    variant?: "list" | "avatar" | "detail" | "original";
    /** Society ID */
    societyId: number;
    /** Visitor entry ID */
    entryId: number;
    /** Flat ID */
    flatId: number;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiResponse =
  /** status 200 Visitor invites fetched successfully */ ModelsVisitorInviteHistoryApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Visitor invite status */
  status?: "active" | "used" | "expired" | "cancelled";
  /** Search purpose, status, visitor, contact, or short code */
  search?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiResponse =
  /** status 201 Visitor invite created successfully */ ModelsVisitorInviteTokenApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
  /** Visitor invite request */
  modelsCreateVisitorInviteRequest: ModelsCreateVisitorInviteRequest;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffApiResponse =
  /** status 201 Visitor invite created successfully */ ModelsVisitorInviteTokenApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Visitor invite request */
    modelsCreateVisitorInviteRequest: ModelsCreateVisitorInviteRequest;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdApiResponse =
  /** status 200 Visitor invite fetched successfully */ ModelsVisitorInviteHistoryItemApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Visitor invite ID */
    inviteId: number;
  };
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsApiResponse =
  /** status 200 Flat visitor settings fetched successfully */ ModelsFlatVisitorSettingsApiResponse;
export type GetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId: number;
};
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetApiResponse =
  /** status 200 Flat visitor settings reset successfully */ ModelsMessageApiResponse;
export type PostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
  };
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeApiResponse =
  /** status 200 Flat visitor setting updated successfully */ ModelsFlatVisitorSettingApiResponse;
export type PatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Flat ID */
    flatId: number;
    /** Visitor purpose */
    purpose:
      | "guest"
      | "delivery"
      | "cab"
      | "service"
      | "maintenance"
      | "staff"
      | "other";
    /** Flat visitor purpose setting update payload */
    modelsUpdateFlatVisitorSettingRequest: ModelsUpdateFlatVisitorSettingRequest;
  };
export type GetV1SocietiesBySocietyIdGuardDeskBootstrapApiResponse =
  /** status 200 Guard desk bootstrap fetched successfully */ ModelsGuardDeskBootstrapApiResponse;
export type GetV1SocietiesBySocietyIdGuardDeskBootstrapApiArg = {
  /** Society ID */
  societyId: number;
};
export type PostV1SocietiesBySocietyIdGuardsApiResponse =
  /** status 201 Guard created successfully */ ModelsGuardApiResponse;
export type PostV1SocietiesBySocietyIdGuardsApiArg = {
  /** Society ID */
  societyId: number;
  /** Create guard payload */
  modelsCreateGuardRequest: ModelsCreateGuardRequest;
};
export type GetV1SocietiesBySocietyIdMembersApiResponse =
  /** status 200 Members fetched successfully */ ModelsPaginatedMembersApiResponse;
export type GetV1SocietiesBySocietyIdMembersApiArg = {
  /** Society ID */
  societyId: number;
  /** Search full name, email, phone, role, or status */
  search?: string;
  /** Member role */
  role?: string;
  /** Member status */
  status?: string;
  /** User ID */
  userId?: number;
  /** Invited by user ID */
  invitedBy?: number;
  /** Removed by user ID */
  removedBy?: number;
  /** Joined from RFC3339 timestamp */
  joinedFrom?: string;
  /** Joined to RFC3339 timestamp */
  joinedTo?: string;
  /** Limit */
  limit?: number;
  /** Offset */
  offset?: number;
  /** Sort by: joined_at, role, status */
  sortBy?: string;
  /** Sort order: asc, desc */
  sortOrder?: string;
};
export type PostV1SocietiesBySocietyIdMembersApiResponse =
  /** status 201 Member added successfully */ ModelsSocietyMemberApiResponse;
export type PostV1SocietiesBySocietyIdMembersApiArg = {
  /** Society ID */
  societyId: number;
  /** Add member payload */
  modelsAddSocietyMemberRequest: ModelsAddSocietyMemberRequest;
};
export type GetV1SocietiesBySocietyIdMembersSummaryApiResponse =
  /** status 200 Member summary fetched successfully */ ModelsSocietyMemberSummaryApiResponse;
export type GetV1SocietiesBySocietyIdMembersSummaryApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdMembersAndMemberIdApiResponse =
  /** status 200 Member fetched successfully */ ModelsSocietyMemberApiResponse;
export type GetV1SocietiesBySocietyIdMembersAndMemberIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Member ID */
  memberId: number;
};
export type GetV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStatsApiResponse =
  /** status 200 Member visitor approval stats fetched successfully */ ModelsMemberVisitorApprovalStatsApiResponse;
export type GetV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStatsApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Member ID */
    memberId: number;
  };
export type DeleteV1SocietiesBySocietyIdMembersAndUserIdApiResponse =
  /** status 200 Member removed successfully */ ModelsMessageApiResponse;
export type DeleteV1SocietiesBySocietyIdMembersAndUserIdApiArg = {
  /** Society ID */
  societyId: number;
  /** User ID */
  userId: number;
  /** Removal reason */
  modelsSocietyReasonRequest: ModelsSocietyReasonRequest;
};
export type PostV1SocietiesBySocietyIdMembersAndUserIdReactivateApiResponse =
  /** status 200 Member reactivated successfully */ ModelsSocietyMemberApiResponse;
export type PostV1SocietiesBySocietyIdMembersAndUserIdReactivateApiArg = {
  /** Society ID */
  societyId: number;
  /** User ID */
  userId: number;
};
export type PatchV1SocietiesBySocietyIdMembersAndUserIdRoleApiResponse =
  /** status 200 Member role changed successfully */ ModelsSocietyMemberApiResponse;
export type PatchV1SocietiesBySocietyIdMembersAndUserIdRoleApiArg = {
  /** Society ID */
  societyId: number;
  /** User ID */
  userId: number;
  /** Change role payload */
  modelsChangeSocietyMemberRoleRequest: ModelsChangeSocietyMemberRoleRequest;
};
export type PostV1SocietiesBySocietyIdMembersAndUserIdSuspendApiResponse =
  /** status 200 Member suspended successfully */ ModelsSocietyMemberApiResponse;
export type PostV1SocietiesBySocietyIdMembersAndUserIdSuspendApiArg = {
  /** Society ID */
  societyId: number;
  /** User ID */
  userId: number;
  /** Suspension reason */
  modelsSocietyReasonRequest: ModelsSocietyReasonRequest;
};
export type GetV1SocietiesBySocietyIdOnboardingBootstrapApiResponse =
  /** status 200 Society onboarding bootstrap fetched successfully */ ModelsSocietyOnboardingBootstrapApiResponse;
export type GetV1SocietiesBySocietyIdOnboardingBootstrapApiArg = {
  /** Society ID */
  societyId: number;
};
export type PostV1SocietiesBySocietyIdReactivateApiResponse =
  /** status 200 Society reactivated successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesBySocietyIdReactivateApiArg = {
  /** Society ID */
  societyId: number;
};
export type PostV1SocietiesBySocietyIdRejectApiResponse =
  /** status 200 Society rejected successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesBySocietyIdRejectApiArg = {
  /** Society ID */
  societyId: number;
  /** Rejection reason */
  modelsSocietyReasonRequest: ModelsSocietyReasonRequest;
};
export type PostV1SocietiesBySocietyIdRestoreApiResponse =
  /** status 200 Society restored successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesBySocietyIdRestoreApiArg = {
  /** Society ID */
  societyId: number;
};
export type PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPendingApiResponse =
  /** status 201 Subscription created successfully */ ModelsSubscriptionApiResponse;
export type PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPendingApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Plan ID */
    planId: number;
  };
export type PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrialApiResponse =
  /** status 201 Trial subscription created successfully */ ModelsSubscriptionApiResponse;
export type PostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrialApiArg = {
  /** Society ID */
  societyId: number;
  /** Plan ID */
  planId: number;
  /** Trial subscription payload */
  modelsCreateTrialSubscriptionRequest: ModelsCreateTrialSubscriptionRequest;
};
export type PostV1SocietiesBySocietyIdSuspendApiResponse =
  /** status 200 Society suspended successfully */ ModelsSocietyApiResponse;
export type PostV1SocietiesBySocietyIdSuspendApiArg = {
  /** Society ID */
  societyId: number;
  /** Suspension reason */
  modelsSocietyReasonRequest: ModelsSocietyReasonRequest;
};
export type PostV1SocietiesBySocietyIdTransferOwnershipApiResponse =
  /** status 200 Ownership transferred successfully */ ModelsSocietyMemberApiResponse;
export type PostV1SocietiesBySocietyIdTransferOwnershipApiArg = {
  /** Society ID */
  societyId: number;
  /** Transfer ownership payload */
  modelsTransferOwnershipRequest: ModelsTransferOwnershipRequest;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesApiResponse =
  /** status 200 Visitor entries fetched successfully */ ModelsVisitorEntriesApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId?: number;
  /** Visitor status */
  status?:
    | "waiting_approval"
    | "approved"
    | "rejected"
    | "checked_in"
    | "checked_out"
    | "cancelled"
    | "expired"
    | "auto_closed";
  /** Visitor entry source */
  source?: "resident_link" | "public_qr" | "guard_entry" | "quick_link";
  /** Visitor purpose */
  purpose?:
    | "guest"
    | "delivery"
    | "cab"
    | "service"
    | "maintenance"
    | "staff"
    | "other";
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesCheckInApiResponse =
  /** status 200 Visitor checked in successfully */ ModelsVisitorEntryApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesCheckInApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor QR token */
  modelsQrTokenRequest: ModelsQrTokenRequest;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesExpectedGuestsApiResponse =
  /** status 200 Expected guest entries fetched successfully */ ModelsVisitorEntriesApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesExpectedGuestsApiArg = {
  /** Society ID */
  societyId: number;
  /** Search by name, phone, flat, vehicle, or purpose */
  search?: string;
  /** Expected window start (RFC3339, defaults to IST today start) */
  eventFrom?: string;
  /** Expected window end (RFC3339, defaults to IST tomorrow start) */
  eventTo?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesGuardApiResponse =
  /** status 201 Visitor entry created successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesGuardApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor details */
  modelsVisitorFormRequest: ModelsVisitorFormRequest;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesPendingApiResponse =
  /** status 200 Pending visitor approvals fetched successfully */ ModelsVisitorPendingEntriesApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesPendingApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId?: number;
  /** Block */
  block?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesStatsApiResponse =
  /** status 200 Visitor entry stats fetched successfully */ ModelsVisitorEntryStatsApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesStatsApiArg = {
  /** Society ID */
  societyId: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesStatsDailyApiResponse =
  /** status 200 Visitor entry daily stats fetched successfully */ ModelsVisitorEntryDailyStatsApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesStatsDailyApiArg = {
  /** Society ID */
  societyId: number;
  /** Number of days (default 7, max 90) */
  days?: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesWaitingAtGateApiResponse =
  /** status 200 Waiting at gate entries fetched successfully */ ModelsVisitorEntriesApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesWaitingAtGateApiArg = {
  /** Society ID */
  societyId: number;
  /** Search by name, phone, flat, vehicle, delivery partner, or purpose */
  search?: string;
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiResponse =
  /** status 200 Visitor entry fetched successfully */ ModelsVisitorEntryApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiResponse =
  /** status 200 Visitor entry updated successfully */ ModelsVisitorEntryMutationApiResponse;
export type PatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
  /** Visitor entry update payload */
  modelsUpdateGuardVisitorEntryRequest: ModelsUpdateGuardVisitorEntryRequest;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveApiResponse =
  /** status 200 Visitor entry approved successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInApiResponse =
  /** status 200 Visitor approved and checked in successfully */ ModelsVisitorEntryApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Visitor entry ID */
    entryId: number;
    /** Guard approval options */
    modelsGuardApproveEntryRequest: ModelsGuardApproveEntryRequest;
  };
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInApiResponse =
  /** status 200 Visitor checked in successfully */ ModelsVisitorEntryApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutApiResponse =
  /** status 200 Visitor checked out successfully */ ModelsVisitorEntryApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEventsApiResponse =
  /** status 200 Visitor entry events fetched successfully */ ModelsVisitorEntryEventsApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEventsApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveApiResponse =
  /** status 200 Visitor entry approved successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveApiArg =
  {
    /** Society ID */
    societyId: number;
    /** Visitor entry ID */
    entryId: number;
    /** Guard approval options */
    modelsGuardApproveEntryRequest: ModelsGuardApproveEntryRequest;
  };
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyApiResponse =
  /** status 200 Resident notified successfully */ ModelsApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse =
  /** status 200 OK */ ModelsImageViewApiResponse;
export type GetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg = {
  /** Server-selected image view */
  variant?: "list" | "avatar" | "detail" | "original";
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PutV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse =
  /** status 200 OK */ ModelsImageViewApiResponse;
export type PutV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
  body: {
    /** JPEG or PNG, maximum 5 MiB and 25 megapixels */
    file: Blob;
  };
};
export type DeleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiResponse =
  unknown;
export type DeleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
};
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectApiResponse =
  /** status 200 Visitor entry rejected successfully */ ModelsVisitorEntryMutationApiResponse;
export type PostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor entry ID */
  entryId: number;
  /** Visitor rejection reason */
  modelsRejectVisitorEntryRequest: ModelsRejectVisitorEntryRequest;
};
export type PostV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancelApiResponse =
  /** status 200 Visitor invite cancelled successfully */ ModelsMessageApiResponse;
export type PostV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancelApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor invite ID */
  inviteId: number;
};
export type GetV1SocietiesBySocietyIdVisitorInvitesAndTokenApiResponse =
  /** status 200 Visitor invite fetched successfully */ ModelsVisitorInviteApiResponse;
export type GetV1SocietiesBySocietyIdVisitorInvitesAndTokenApiArg = {
  /** Society ID */
  societyId: number;
  /** Visitor invite token */
  token: string;
};
export type GetV1SocietiesBySocietyIdVisitorSettingsApiResponse =
  /** status 200 Visitor settings fetched successfully */ ModelsSocietyVisitorSettingsApiResponse;
export type GetV1SocietiesBySocietyIdVisitorSettingsApiArg = {
  /** Society ID */
  societyId: number;
};
export type PatchV1SocietiesBySocietyIdVisitorSettingsApiResponse =
  /** status 200 Visitor settings updated successfully */ ModelsSocietyVisitorSettingsApiResponse;
export type PatchV1SocietiesBySocietyIdVisitorSettingsApiArg = {
  /** Society ID */
  societyId: number;
  /** Society visitor settings update payload */
  modelsUpdateSocietyVisitorSettingsRequest: ModelsUpdateSocietyVisitorSettingsRequest;
};
export type GetV1SocietiesBySocietyIdVisitorSettingsFlatsApiResponse =
  /** status 200 Society flat visitor settings fetched successfully */ ModelsSocietyFlatVisitorSettingsApiResponse;
export type GetV1SocietiesBySocietyIdVisitorSettingsFlatsApiArg = {
  /** Society ID */
  societyId: number;
  /** Flat ID */
  flatId?: number;
  /** Block */
  block?: string;
  /** Visitor purpose */
  purpose?:
    | "guest"
    | "delivery"
    | "cab"
    | "service"
    | "maintenance"
    | "staff"
    | "other";
  /** Maximum records to return */
  limit?: number;
  /** Records to skip */
  offset?: number;
};
export type GetV1SubscriptionsApiResponse =
  /** status 200 Subscriptions fetched successfully */ ModelsSubscriptionsApiResponse;
export type GetV1SubscriptionsApiArg = {
  /** Society ID */
  societyId?: number;
  /** Plan ID */
  planId?: number;
  /** Subscription status */
  status?: string;
  /** Search text */
  search?: string;
  /** Search mode */
  searchMode?: string;
  /** Only active subscriptions */
  isActiveOnly?: boolean;
  /** Only expired subscriptions */
  expiredOnly?: boolean;
  /** Subscriptions ending before RFC3339 timestamp */
  expiringBefore?: string;
  /** Subscriptions ending within N days */
  expiringDays?: number;
  /** Page size */
  limit?: number;
  /** Page offset */
  offset?: number;
};
export type GetV1SubscriptionsLookupApiResponse =
  /** status 200 Subscription fetched successfully */ ModelsSubscriptionApiResponse;
export type GetV1SubscriptionsLookupApiArg = {
  /** Subscription ID */
  id?: number;
  /** Society ID */
  societyId?: number;
};
export type GetV1SubscriptionsStatsApiResponse =
  /** status 200 Subscription stats fetched successfully */ ModelsSubscriptionStatsApiResponse;
export type GetV1SubscriptionsStatsApiArg = void;
export type PostV1SubscriptionsBySubscriptionIdActivateApiResponse =
  /** status 200 Subscription activated successfully */ ModelsSubscriptionApiResponse;
export type PostV1SubscriptionsBySubscriptionIdActivateApiArg = {
  /** Subscription ID */
  subscriptionId: number;
  /** Activate subscription payload */
  modelsActivateSubscriptionRequest: ModelsActivateSubscriptionRequest;
};
export type PostV1SubscriptionsBySubscriptionIdCancelApiResponse =
  /** status 200 Subscription cancelled successfully */ ModelsSubscriptionApiResponse;
export type PostV1SubscriptionsBySubscriptionIdCancelApiArg = {
  /** Subscription ID */
  subscriptionId: number;
  /** Cancel subscription payload */
  modelsCancelSubscriptionRequest: ModelsCancelSubscriptionRequest;
};
export type PostV1SubscriptionsBySubscriptionIdExpireApiResponse =
  /** status 200 Subscription expired successfully */ ModelsSubscriptionApiResponse;
export type PostV1SubscriptionsBySubscriptionIdExpireApiArg = {
  /** Subscription ID */
  subscriptionId: number;
};
export type PostV1SubscriptionsBySubscriptionIdPlansAndPlanIdApiResponse =
  /** status 200 Subscription plan changed successfully */ ModelsSubscriptionApiResponse;
export type PostV1SubscriptionsBySubscriptionIdPlansAndPlanIdApiArg = {
  /** Subscription ID */
  subscriptionId: number;
  /** New plan ID */
  planId: number;
};
export type PostV1SubscriptionsBySubscriptionIdRenewApiResponse =
  /** status 200 Subscription renewed successfully */ ModelsSubscriptionApiResponse;
export type PostV1SubscriptionsBySubscriptionIdRenewApiArg = {
  /** Subscription ID */
  subscriptionId: number;
  /** Renew subscription payload */
  modelsRenewSubscriptionRequest: ModelsRenewSubscriptionRequest;
};
export type GetV1MeNotificationsByNotificationIdApiResponse =
  /** status 200 OK */ ModelsNotificationDetailApiResponse;
export type GetV1MeNotificationsByNotificationIdApiArg = {
  /** Notification ID */
  notificationId: string;
};
export type ModelsHealthCheckResponseDoc = {
  app?: string;
  environment?: string;
  status?: string;
  timestamp?: string;
  version?: string;
};
export type ModelsLivenessResponseDoc = {
  status?: string;
  timestamp?: string;
};
export type ModelsReadinessResponseDoc = {
  database?: string;
  status?: string;
  timestamp?: string;
};
export type ModelsHealthErrorResponseDoc = {
  error?: string;
  status?: string;
};
export type ModelsMessageData = {
  message?: string;
};
export type ModelsChangePasswordApiResponse = {
  data?: ModelsMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsErrorData = {
  code?: string;
  details?: {
    [key: string]: any;
  };
  message?: string;
};
export type ModelsErrorResponseDoc = {
  error?: ModelsErrorData;
  success?: boolean;
};
export type ModelsChangePasswordRequest = {
  confirm_password: string;
  current_password: string;
  new_password: string;
};
export type ModelsDevOtpMessageData = {
  dev_otp?: string;
  message?: string;
};
export type ModelsForgotPasswordApiResponse = {
  data?: ModelsDevOtpMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsForgotPasswordRequest = {
  email: string;
};
export type ModelsAuthProvider = "email" | "google" | "apple" | "phone";
export type ModelsGlobalRole = "user" | "developer" | "super_admin";
export type ModelsUserResponse = {
  auth_provider?: ModelsAuthProvider;
  avatar_url?: string;
  blocked_reason?: string;
  created_at?: string;
  date_of_birth?: string;
  email?: string;
  email_verified?: boolean;
  first_name?: string;
  full_name?: string;
  gender?: string;
  global_role?: ModelsGlobalRole;
  id?: number;
  is_active?: boolean;
  is_blocked?: boolean;
  language?: string;
  last_login_at?: string;
  last_name?: string;
  phone_number?: string;
  phone_verified?: boolean;
  timezone?: string;
  updated_at?: string;
};
export type ModelsAuthSessionData = {
  access_token?: string;
  access_token_expires_at?: string;
  refresh_token?: string;
  refresh_token_expires_at?: string;
  user?: ModelsUserResponse;
};
export type ModelsLoginApiResponse = {
  data?: ModelsAuthSessionData;
  message?: string;
  success?: boolean;
};
export type ModelsLoginRequest = {
  email?: string;
  password: string;
  phone_number?: string;
};
export type ModelsLogoutApiResponse = {
  data?: ModelsMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsUserData = {
  user?: ModelsUserResponse;
};
export type ModelsGetProfileApiResponse = {
  data?: ModelsUserData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateUserRequest = {
  avatar_url?: string;
  date_of_birth?: string;
  first_name?: string;
  gender?: "male" | "female" | "other" | "prefer_not_to_say";
  language?: string;
  last_name?: string;
  phone_number?: string;
  timezone?: string;
};
export type ModelsImageVariant = "original" | "list" | "avatar" | "detail";
export type ModelsImageView = {
  expires_at?: string;
  url?: string;
  variant?: ModelsImageVariant;
};
export type ModelsImageViewApiResponse = {
  data?: ModelsImageView;
  message?: string;
  success?: boolean;
};
export type ModelsAuthRefreshData = {
  access_token?: string;
  access_token_expires_at?: string;
  message?: string;
};
export type ModelsRefreshTokenApiResponse = {
  data?: ModelsAuthRefreshData;
  message?: string;
  success?: boolean;
};
export type ModelsRefreshTokenRequest = {
  refresh_token?: string;
};
export type ModelsUserMessageData = {
  dev_otp?: string;
  message?: string;
  user?: ModelsUserResponse;
};
export type ModelsRegisterApiResponse = {
  data?: ModelsUserMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsRegisterRequest = {
  email: string;
  first_name: string;
  last_name?: string;
  password: string;
  phone_number: string;
};
export type ModelsResendOtpapiResponse = {
  data?: ModelsDevOtpMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsResendOtpRequest = {
  email: string;
};
export type ModelsResetPasswordApiResponse = {
  data?: ModelsMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsResetPasswordRequest = {
  confirm_password: string;
  email: string;
  new_password: string;
  otp: string;
};
export type ModelsResidentRegisterRequest = {
  email: string;
  first_name: string;
  last_name?: string;
  password: string;
  phone_number: string;
};
export type ModelsVerifyOtpapiResponse = {
  data?: ModelsUserData;
  message?: string;
  success?: boolean;
};
export type ModelsVerifyOtpRequest = {
  email: string;
  otp: string;
};
export type ModelsDashboardKind =
  "developer" | "society_admin" | "select_society" | "onboarding";
export type ModelsDefaultDashboardResponse = {
  kind?: ModelsDashboardKind;
  path?: string;
  society_id?: number;
};
export type ModelsSocietyMemberRole = "owner" | "admin" | "staff" | "resident";
export type ModelsSocietyMemberStatus =
  "pending" | "active" | "suspended" | "removed";
export type ModelsSocietyMemberResponse = {
  created_at?: string;
  id?: number;
  invited_by?: number;
  joined_at?: string;
  remove_reason?: string;
  removed_at?: string;
  removed_by?: number;
  role?: ModelsSocietyMemberRole;
  society_id?: number;
  status?: ModelsSocietyMemberStatus;
  updated_at?: string;
  user_email?: string;
  user_full_name?: string;
  user_id?: number;
  user_phone?: string;
};
export type ModelsFlatResidentRole = "owner" | "tenant" | "family";
export type ModelsFlatResidentStatus = "active" | "inactive" | "moved_out";
export type ModelsFlatResidentResponse = {
  block?: string;
  created_at?: string;
  created_by?: number;
  flat_id?: number;
  flat_number?: string;
  floor?: string;
  id?: number;
  is_primary?: boolean;
  moved_in_at?: string;
  moved_out_at?: string;
  role?: ModelsFlatResidentRole;
  society_code?: string;
  society_id?: number;
  society_name?: string;
  status?: ModelsFlatResidentStatus;
  updated_at?: string;
  user_email?: string;
  user_id?: number;
  user_name?: string;
  user_phone?: string;
};
export type ModelsBootstrapData = {
  defaultDashboard?: ModelsDefaultDashboardResponse;
  memberships?: ModelsSocietyMemberResponse[];
  residences?: ModelsFlatResidentResponse[];
  user?: ModelsUserResponse;
};
export type ModelsBootstrapApiResponse = {
  data?: ModelsBootstrapData;
  message?: string;
  success?: boolean;
};
export type ModelsDeveloperDashboardPlanStatsResponse = {
  active?: number;
  inactive?: number;
  total?: number;
};
export type ModelsSocietyStatus =
  "pending" | "active" | "suspended" | "rejected";
export type ModelsSocietyResponse = {
  address_line1?: string;
  address_line2?: string;
  approved_at?: string;
  approved_by?: number;
  city?: string;
  country?: string;
  created_at?: string;
  created_by?: number;
  email?: string;
  id?: number;
  landmark?: string;
  name?: string;
  phone_number?: string;
  pincode?: string;
  rejected_at?: string;
  rejected_by?: number;
  rejection_reason?: string;
  society_code?: string;
  state?: string;
  status?: ModelsSocietyStatus;
  suspended_at?: string;
  suspended_by?: number;
  suspension_reason?: string;
  total_blocks?: number;
  total_flats?: number;
  updated_at?: string;
};
export type ModelsBillingCycle = "monthly" | "yearly";
export type ModelsSubscriptionStatus =
  "pending" | "trial" | "active" | "expired" | "cancelled";
export type ModelsSocietySubscriptionResponse = {
  activated_at?: string;
  activated_by?: number;
  billing_cycle?: ModelsBillingCycle;
  cancellation_reason?: string;
  cancelled_at?: string;
  cancelled_by?: number;
  created_at?: string;
  created_by?: number;
  currency?: string;
  current_plan_code?: string;
  current_plan_name?: string;
  ends_at?: string;
  expired_at?: string;
  features?: {
    [key: string]: any;
  };
  id?: number;
  max_admins?: number;
  max_flats?: number;
  max_residents?: number;
  max_staff?: number;
  plan_code?: string;
  plan_id?: number;
  plan_name?: string;
  price_amount_paise?: number;
  society_code?: string;
  society_id?: number;
  society_name?: string;
  starts_at?: string;
  status?: ModelsSubscriptionStatus;
  trial_ends_at?: string;
  updated_at?: string;
};
export type ModelsDeveloperDashboardResidenceStatsResponse = {
  active_residents?: number;
  total_residents?: number;
};
export type ModelsDeveloperDashboardSocietyStatsResponse = {
  active?: number;
  pending?: number;
  rejected?: number;
  suspended?: number;
  total?: number;
};
export type ModelsSubscriptionStatsResponse = {
  active_subscriptions?: number;
  cancelled_subscriptions?: number;
  expired_subscriptions?: number;
  pending_subscriptions?: number;
  total_subscriptions?: number;
  trial_subscriptions?: number;
};
export type ModelsDeveloperDashboardBootstrapResponse = {
  plan_stats?: ModelsDeveloperDashboardPlanStatsResponse;
  recent_pending_societies?: ModelsSocietyResponse[];
  recent_subscriptions?: ModelsSocietySubscriptionResponse[];
  residence_stats?: ModelsDeveloperDashboardResidenceStatsResponse;
  society_stats?: ModelsDeveloperDashboardSocietyStatsResponse;
  subscription_stats?: ModelsSubscriptionStatsResponse;
};
export type ModelsDeveloperDashboardBootstrapData = {
  dashboard?: ModelsDeveloperDashboardBootstrapResponse;
};
export type ModelsDeveloperDashboardBootstrapApiResponse = {
  data?: ModelsDeveloperDashboardBootstrapData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatStatus = "vacant" | "occupied" | "blocked";
export type ModelsFlatClaimStatus =
  "pending" | "approved" | "rejected" | "cancelled";
export type ModelsFlatClaimResponse = {
  block?: string;
  cancelled_at?: string;
  created_at?: string;
  flat_id?: number;
  flat_number?: string;
  flat_status?: ModelsFlatStatus;
  floor?: string;
  id?: number;
  note?: string;
  rejection_reason?: string;
  requested_primary?: boolean;
  requested_role?: ModelsFlatResidentRole;
  reviewed_at?: string;
  reviewed_by?: number;
  reviewer_email?: string;
  reviewer_name?: string;
  reviewer_phone?: string;
  society_code?: string;
  society_id?: number;
  society_name?: string;
  status?: ModelsFlatClaimStatus;
  updated_at?: string;
  user_email?: string;
  user_id?: number;
  user_name?: string;
  user_phone?: string;
};
export type ModelsFlatClaimsData = {
  claims?: ModelsFlatClaimResponse[];
};
export type ModelsFlatClaimsApiResponse = {
  data?: ModelsFlatClaimsData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatClaimData = {
  claim?: ModelsFlatClaimResponse;
};
export type ModelsFlatClaimApiResponse = {
  data?: ModelsFlatClaimData;
  message?: string;
  success?: boolean;
};
export type ModelsSubmitFlatClaimRequest = {
  flat_id: number;
  metadata?: {
    [key: string]: any;
  };
  note?: string;
  requested_primary?: boolean;
  requested_role: ModelsFlatResidentRole;
  society_id: number;
};
export type ModelsFlatResidentsData = {
  residents?: ModelsFlatResidentResponse[];
};
export type ModelsFlatResidentsApiResponse = {
  data?: ModelsFlatResidentsData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatResponse = {
  block?: string;
  created_at?: string;
  created_by?: number;
  flat_number?: string;
  floor?: string;
  id?: number;
  is_active?: boolean;
  primary_resident_name?: string;
  society_code?: string;
  society_id?: number;
  society_name?: string;
  status?: ModelsFlatStatus;
  updated_at?: string;
};
export type ModelsFlatsData = {
  flats?: ModelsFlatResponse[];
};
export type ModelsFlatsApiResponse = {
  data?: ModelsFlatsData;
  message?: string;
  success?: boolean;
};
export type ModelsDevicePlatform = "ios" | "android" | "web";
export type ModelsDeviceToken = {
  created_at?: string;
  device_id?: string;
  id?: number;
  last_seen_at?: string;
  platform?: ModelsDevicePlatform;
  token?: string;
  updated_at?: string;
  user_id?: number;
};
export type ModelsDeviceTokenData = {
  device_token?: ModelsDeviceToken;
};
export type ModelsRegisterDeviceTokenApiResponse = {
  data?: ModelsDeviceTokenData;
  message?: string;
  success?: boolean;
};
export type ModelsRegisterDeviceTokenRequest = {
  device_id?: string;
  platform: ModelsDevicePlatform;
  token: string;
};
export type ModelsMessageApiResponse = {
  data?: ModelsMessageData;
  message?: string;
  success?: boolean;
};
export type ModelsUnregisterDeviceTokenRequest = {
  token: string;
};
export type ModelsNotification = {
  body?: string;
  created_at?: string;
  data?: {
    [key: string]: any;
  };
  event_key?: string;
  flat_id?: number;
  id?: string;
  read_at?: string;
  society_id?: number;
  title?: string;
  type?: string;
  user_id?: number;
};
export type ModelsNotificationListResult = {
  items?: ModelsNotification[];
  next_cursor?: string;
};
export type ModelsNotificationsApiResponse = {
  data?: ModelsNotificationListResult;
  message?: string;
  success?: boolean;
};
export type ModelsNotificationUnreadCountResponse = {
  unread_count?: number;
};
export type ModelsNotificationUnreadCountApiResponse = {
  data?: ModelsNotificationUnreadCountResponse;
  message?: string;
  success?: boolean;
};
export type ModelsNotificationReadResponse = {
  id?: string;
  read_at?: string;
};
export type ModelsNotificationReadApiResponse = {
  data?: ModelsNotificationReadResponse;
  message?: string;
  success?: boolean;
};
export type ModelsMyResidencesData = {
  residences?: ModelsFlatResidentResponse[];
};
export type ModelsMyResidencesApiResponse = {
  data?: ModelsMyResidencesData;
  message?: string;
  success?: boolean;
};
export type ModelsPlanResponse = {
  billing_cycle?: ModelsBillingCycle;
  code?: string;
  created_at?: string;
  currency?: string;
  description?: string;
  features?: {
    [key: string]: any;
  };
  id?: number;
  is_active?: boolean;
  max_admins?: number;
  max_flats?: number;
  max_residents?: number;
  max_staff?: number;
  name?: string;
  price_amount_paise?: number;
  updated_at?: string;
};
export type ModelsPlansData = {
  plans?: ModelsPlanResponse[];
};
export type ModelsPlansApiResponse = {
  data?: ModelsPlansData;
  message?: string;
  success?: boolean;
};
export type ModelsPlanData = {
  plan?: ModelsPlanResponse;
};
export type ModelsPlanApiResponse = {
  data?: ModelsPlanData;
  message?: string;
  success?: boolean;
};
export type ModelsCreatePlanRequest = {
  billing_cycle: ModelsBillingCycle;
  code: string;
  currency: string;
  description?: string;
  features?: {
    [key: string]: any;
  };
  max_admins?: number;
  max_flats: number;
  max_residents: number;
  max_staff?: number;
  name: string;
  price_amount_paise?: number;
};
export type ModelsUpdatePlanRequest = {
  billing_cycle?: ModelsBillingCycle;
  code?: string;
  currency?: string;
  description?: string;
  features?: {
    [key: string]: any;
  };
  max_admins?: number;
  max_flats?: number;
  max_residents?: number;
  max_staff?: number;
  name?: string;
  price_amount_paise?: number;
};
export type ModelsFlatMemberInviteRole = "family" | "tenant";
export type ModelsFlatMemberInviteStatus =
  "pending" | "accepted" | "expired" | "cancelled";
export type ModelsPublicFlatMemberInviteView = {
  block?: string;
  email?: string;
  expires_at?: string;
  flat_number?: string;
  floor?: string;
  full_name?: string;
  id?: number;
  phone?: string;
  role?: ModelsFlatMemberInviteRole;
  society_name?: string;
  status?: ModelsFlatMemberInviteStatus;
};
export type ModelsPublicFlatMemberInviteData = {
  invite?: ModelsPublicFlatMemberInviteView;
};
export type ModelsPublicFlatMemberInviteApiResponse = {
  data?: ModelsPublicFlatMemberInviteData;
  message?: string;
  success?: boolean;
};
export type ModelsShortLinkResponse = {
  code?: string;
  expires_at?: string;
  url?: string;
};
export type ModelsFlatMemberInviteResponse = {
  created_at?: string;
  email?: string;
  expires_at?: string;
  flat_id?: number;
  full_name?: string;
  id?: number;
  invited_by?: number;
  phone?: string;
  role?: ModelsFlatMemberInviteRole;
  short_link?: ModelsShortLinkResponse;
  society_id?: number;
  status?: ModelsFlatMemberInviteStatus;
  updated_at?: string;
};
export type ModelsAcceptFlatMemberInviteResponse = {
  invite?: ModelsFlatMemberInviteResponse;
  resident?: ModelsFlatResidentResponse;
};
export type ModelsAcceptFlatMemberInviteData = {
  acceptance?: ModelsAcceptFlatMemberInviteResponse;
};
export type ModelsAcceptFlatMemberInviteApiResponse = {
  data?: ModelsAcceptFlatMemberInviteData;
  message?: string;
  success?: boolean;
};
export type ModelsJoinFlatMemberInviteResponse = {
  acceptance?: ModelsAcceptFlatMemberInviteResponse;
  user?: ModelsUserResponse;
};
export type ModelsJoinFlatMemberInviteData = {
  join?: ModelsJoinFlatMemberInviteResponse;
};
export type ModelsJoinFlatMemberInviteApiResponse = {
  data?: ModelsJoinFlatMemberInviteData;
  message?: string;
  success?: boolean;
};
export type ModelsJoinFlatMemberInviteRequest = {
  email?: string;
  first_name?: string;
  identifier?: string;
  last_name?: string;
  password: string;
};
export type ModelsVisitorFlatSummary = {
  block?: string;
  flat_number?: string;
  floor?: string;
  id?: number;
};
export type ModelsVisitorPurpose =
  "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
export type ModelsVisitorEntrySource =
  "resident_link" | "public_qr" | "guard_entry" | "quick_link";
export type ModelsVisitorStatus =
  | "waiting_approval"
  | "approved"
  | "rejected"
  | "checked_in"
  | "checked_out"
  | "cancelled"
  | "expired"
  | "auto_closed";
export type ModelsVisitorVehicleType =
  "bike" | "car" | "auto" | "cab" | "truck" | "other";
export type ModelsVisitorSummary = {
  email?: string;
  full_name?: string;
  phone_number?: string;
  photo_url?: string;
};
export type ModelsVisitorEntry = {
  approved_at?: string;
  approved_by?: number;
  approver_name?: string;
  approver_phone_number?: string;
  auto_closed_at?: string;
  checked_in_at?: string;
  checked_out_at?: string;
  companion_details?: {
    [key: string]: any;
  }[];
  companions_count?: number;
  created_at?: string;
  created_by?: number;
  delivery_partner?: string;
  expected_at?: string;
  expected_checkout_at?: string;
  flat?: ModelsVisitorFlatSummary;
  flat_id?: number;
  handled_by_guard_id?: number;
  handled_by_guard_name?: string;
  handled_by_guard_phone?: string;
  id?: number;
  invite_id?: number;
  metadata?: {
    [key: string]: any;
  };
  notes?: string;
  purpose?: ModelsVisitorPurpose;
  qr_expires_at?: string;
  qr_used_at?: string;
  rejected_by?: number;
  rejection_reason?: string;
  service_provider?: string;
  society_id?: number;
  source?: ModelsVisitorEntrySource;
  status?: ModelsVisitorStatus;
  updated_at?: string;
  vehicle_number?: string;
  vehicle_type?: ModelsVisitorVehicleType;
  visitor?: ModelsVisitorSummary;
  visitor_id?: number;
};
export type ModelsQrTokenResponse = {
  expires_at?: string;
  token?: string;
};
export type ModelsShortLinkState =
  | "form"
  | "waiting_approval"
  | "approved"
  | "checked_in"
  | "completed"
  | "join"
  | "accepted"
  | "expired"
  | "cancelled"
  | "unavailable";
export type ModelsShortLinkResourceType = "visitor_invite" | "member_invite";
export type ModelsShortLinkResolveResponse = {
  entry?: ModelsVisitorEntry;
  invite?: any;
  link?: ModelsShortLinkResponse;
  qr?: ModelsQrTokenResponse;
  state?: ModelsShortLinkState;
  type?: ModelsShortLinkResourceType;
  visitor?: ModelsVisitorSummary;
};
export type ModelsShortLinkResolveData = {
  link?: ModelsShortLinkResolveResponse;
};
export type ModelsShortLinkResolveApiResponse = {
  data?: ModelsShortLinkResolveData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntryMutationResponse = {
  entry?: ModelsVisitorEntry;
  qr?: ModelsQrTokenResponse;
};
export type ModelsVisitorEntryMutationApiResponse = {
  data?: ModelsVisitorEntryMutationResponse;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorFormRequest = {
  companion_details?: {
    [key: string]: any;
  }[];
  companions_count?: number;
  delivery_partner?: string;
  email?: string;
  expected_at?: string;
  expected_checkout_at?: string;
  flat_id?: number;
  full_name?: string;
  metadata?: {
    [key: string]: any;
  };
  notes?: string;
  phone_number?: string;
  photo_url?: string;
  purpose?: ModelsVisitorPurpose;
  service_provider?: string;
  vehicle_number?: string;
  vehicle_type?: ModelsVisitorVehicleType;
};
export type ModelsPublicClaimFlatResponse = {
  block?: string;
  flat_number?: string;
  floor?: string;
  id?: number;
  status?: ModelsFlatStatus;
};
export type ModelsPublicClaimSocietyResponse = {
  city?: string;
  country?: string;
  id?: number;
  name?: string;
  pincode?: string;
  society_code?: string;
  state?: string;
  total_flats?: number;
};
export type ModelsPublicClaimOptionsData = {
  flats?: ModelsPublicClaimFlatResponse[];
  society?: ModelsPublicClaimSocietyResponse;
};
export type ModelsPublicClaimOptionsApiResponse = {
  data?: ModelsPublicClaimOptionsData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntryOptionsFlat = {
  block?: string;
  flat_number?: string;
  floor?: string;
  id?: number;
};
export type ModelsVisitorEntryOptionsBlock = {
  block?: string;
  flats?: ModelsVisitorEntryOptionsFlat[];
};
export type ModelsVisitorEntryOptionsResponse = {
  blocks?: ModelsVisitorEntryOptionsBlock[];
  flats?: ModelsVisitorEntryOptionsFlat[];
  has_more?: boolean;
  purposes?: ModelsVisitorPurpose[];
};
export type ModelsVisitorEntryOptionsData = {
  options?: ModelsVisitorEntryOptionsResponse;
};
export type ModelsVisitorEntryOptionsApiResponse = {
  data?: ModelsVisitorEntryOptionsData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntryData = {
  entry?: ModelsVisitorEntry;
};
export type ModelsVisitorEntryApiResponse = {
  data?: ModelsVisitorEntryData;
  message?: string;
  success?: boolean;
};
export type ModelsQrTokenRequest = {
  token?: string;
};
export type ModelsVisitorInviteStatus =
  "active" | "used" | "expired" | "cancelled";
export type ModelsPublicVisitorInviteView = {
  block?: string;
  expires_at?: string;
  flat_number?: string;
  floor?: string;
  id?: number;
  purpose?: ModelsVisitorPurpose;
  society_name?: string;
  status?: ModelsVisitorInviteStatus;
};
export type ModelsVisitorInviteData = {
  entry?: ModelsVisitorEntry;
  invite?: ModelsPublicVisitorInviteView;
  qr?: ModelsQrTokenResponse;
  view?: string;
};
export type ModelsVisitorInviteApiResponse = {
  data?: ModelsVisitorInviteData;
  message?: string;
  success?: boolean;
};
export type ModelsPaginatedSocietiesResponse = {
  items?: ModelsSocietyResponse[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsPaginatedSocietiesData = {
  societies?: ModelsPaginatedSocietiesResponse;
};
export type ModelsPaginatedSocietiesApiResponse = {
  data?: ModelsPaginatedSocietiesData;
  message?: string;
  success?: boolean;
};
export type ModelsSocietyData = {
  society?: ModelsSocietyResponse;
};
export type ModelsSocietyApiResponse = {
  data?: ModelsSocietyData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateSocietyRequest = {
  address_line1?: string;
  address_line2?: string;
  city?: string;
  country?: string;
  email?: string;
  landmark?: string;
  metadata?: {
    [key: string]: any;
  };
  name: string;
  phone_number?: string;
  pincode?: string;
  society_code?: string;
  state?: string;
  total_blocks?: number;
  total_flats?: number;
};
export type ModelsMySocietyResponse = {
  member?: ModelsSocietyMemberResponse;
  society?: ModelsSocietyResponse;
};
export type ModelsMySocietiesData = {
  societies?: ModelsMySocietyResponse[];
};
export type ModelsMySocietiesApiResponse = {
  data?: ModelsMySocietiesData;
  message?: string;
  success?: boolean;
};
export type ModelsSocietyDetailResponse = {
  address_line1?: string;
  address_line2?: string;
  approved_at?: string;
  approved_by?: number;
  city?: string;
  country?: string;
  created_at?: string;
  created_by?: number;
  email?: string;
  id?: number;
  landmark?: string;
  members_count?: number;
  name?: string;
  phone_number?: string;
  pincode?: string;
  rejected_at?: string;
  rejected_by?: number;
  rejection_reason?: string;
  society_code?: string;
  state?: string;
  status?: ModelsSocietyStatus;
  suspended_at?: string;
  suspended_by?: number;
  suspension_reason?: string;
  total_blocks?: number;
  total_flats?: number;
  updated_at?: string;
};
export type ModelsSocietyDetailData = {
  society?: ModelsSocietyDetailResponse;
};
export type ModelsSocietyDetailApiResponse = {
  data?: ModelsSocietyDetailData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateSocietyRequest = {
  address_line1?: string;
  address_line2?: string;
  city?: string;
  country?: string;
  email?: string;
  landmark?: string;
  metadata?: {
    [key: string]: any;
  };
  name?: string;
  phone_number?: string;
  pincode?: string;
  state?: string;
  total_blocks?: number;
  total_flats?: number;
};
export type ModelsPaginatedMembersResponse = {
  items?: ModelsSocietyMemberResponse[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsPaginatedMembersData = {
  members?: ModelsPaginatedMembersResponse;
};
export type ModelsPaginatedMembersApiResponse = {
  data?: ModelsPaginatedMembersData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatClaimStatsResponse = {
  approved_claims?: number;
  cancelled_claims?: number;
  pending_claims?: number;
  rejected_claims?: number;
  total_claims?: number;
};
export type ModelsFlatStatsResponse = {
  active_flats?: number;
  blocked_flats?: number;
  inactive_flats?: number;
  occupied_flats?: number;
  society_id?: number;
  total_flats?: number;
  vacant_flats?: number;
};
export type ModelsSocietyDashboardMemberStatsResponse = {
  admins?: number;
  owners?: number;
  residents?: number;
  staff?: number;
  total_active_members?: number;
};
export type ModelsSocietyDashboardSubscriptionHealthResponse = {
  days_until_expiry?: number;
  is_active?: boolean;
  is_expiring_soon?: boolean;
  lifecycle_label?: string;
};
export type ModelsSocietyDashboardQuotaUsageResponse = {
  limit?: number;
  percent?: number;
  remaining?: number;
  used?: number;
};
export type ModelsSocietyDashboardSubscriptionUsageResponse = {
  admins?: ModelsSocietyDashboardQuotaUsageResponse;
  flats?: ModelsSocietyDashboardQuotaUsageResponse;
  residents?: ModelsSocietyDashboardQuotaUsageResponse;
  staff?: ModelsSocietyDashboardQuotaUsageResponse;
};
export type ModelsVisitorDailyCountResponse = {
  count?: number;
  date?: string;
};
export type ModelsVisitorEntryStatsResponse = {
  auto_closed_today?: number;
  checked_out_in_range?: number;
  checked_out_today?: number;
  pending_approvals?: number;
  rejected_today?: number;
  today_visitors?: number;
  visitors_inside?: number;
};
export type ModelsSocietyDashboardBootstrapResponse = {
  claim_stats?: ModelsFlatClaimStatsResponse;
  current_subscription?: ModelsSocietySubscriptionResponse;
  flat_stats?: ModelsFlatStatsResponse;
  member_stats?: ModelsSocietyDashboardMemberStatsResponse;
  plan_ads?: ModelsPlanResponse[];
  recent_pending_claims?: ModelsFlatClaimResponse[];
  society?: ModelsSocietyResponse;
  subscription_health?: ModelsSocietyDashboardSubscriptionHealthResponse;
  subscription_usage?: ModelsSocietyDashboardSubscriptionUsageResponse;
  visitor_daily_last_7_days?: ModelsVisitorDailyCountResponse[];
  visitor_stats?: ModelsVisitorEntryStatsResponse;
};
export type ModelsSocietyDashboardBootstrapData = {
  dashboard?: ModelsSocietyDashboardBootstrapResponse;
};
export type ModelsSocietyDashboardBootstrapApiResponse = {
  data?: ModelsSocietyDashboardBootstrapData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatClaimStatsData = {
  stats?: ModelsFlatClaimStatsResponse;
};
export type ModelsFlatClaimStatsApiResponse = {
  data?: ModelsFlatClaimStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsApproveFlatClaimResponse = {
  claim?: ModelsFlatClaimResponse;
  flat?: ModelsFlatResponse;
  resident?: ModelsFlatResidentResponse;
};
export type ModelsFlatApprovalData = {
  approval?: ModelsApproveFlatClaimResponse;
};
export type ModelsFlatApprovalApiResponse = {
  data?: ModelsFlatApprovalData;
  message?: string;
  success?: boolean;
};
export type ModelsRejectFlatClaimRequest = {
  reason: string;
};
export type ModelsPaginatedFlatsResponse = {
  items?: ModelsFlatResponse[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsPaginatedFlatsData = {
  flats?: ModelsPaginatedFlatsResponse;
};
export type ModelsPaginatedFlatsApiResponse = {
  data?: ModelsPaginatedFlatsData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatData = {
  flat?: ModelsFlatResponse;
};
export type ModelsFlatApiResponse = {
  data?: ModelsFlatData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateFlatRequest = {
  block?: string;
  flat_number: string;
  floor?: string;
  metadata?: {
    [key: string]: any;
  };
};
export type ModelsBulkCreateFlatsResponse = {
  items?: ModelsFlatResponse[];
  total?: number;
};
export type ModelsBulkFlatsData = {
  flats?: ModelsBulkCreateFlatsResponse;
};
export type ModelsBulkFlatsApiResponse = {
  data?: ModelsBulkFlatsData;
  message?: string;
  success?: boolean;
};
export type ModelsBulkCreateFlatsRequest = {
  flats: ModelsCreateFlatRequest[];
};
export type ModelsFlatNumberingMode = "continuous" | "floor_based";
export type ModelsGenerateFlatsRequest = {
  block: string;
  flats_per_floor: number;
  number_of_floors?: number;
  numbering_mode: ModelsFlatNumberingMode;
  sequence_start?: string;
  start_floor?: number;
  total_flats?: number;
  unit_start?: string;
};
export type ModelsFlatStatsData = {
  stats?: ModelsFlatStatsResponse;
};
export type ModelsFlatStatsApiResponse = {
  data?: ModelsFlatStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateFlatRequest = {
  block?: string;
  flat_number?: string;
  floor?: string;
  is_active?: boolean;
  metadata?: {
    [key: string]: any;
  };
  status?: ModelsFlatStatus;
};
export type ModelsFlatMemberInvitesData = {
  invites?: ModelsFlatMemberInviteResponse[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsFlatMemberInvitesApiResponse = {
  data?: ModelsFlatMemberInvitesData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatMemberInviteTokenResponse = {
  expires_at?: string;
  token?: string;
};
export type ModelsFlatMemberInviteTokenData = {
  invite?: ModelsFlatMemberInviteResponse;
  link?: ModelsShortLinkResponse;
  short_link?: ModelsShortLinkResponse;
  token?: ModelsFlatMemberInviteTokenResponse;
};
export type ModelsFlatMemberInviteTokenApiResponse = {
  data?: ModelsFlatMemberInviteTokenData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateFlatMemberInviteRequest = {
  email?: string;
  full_name: string;
  phone?: string;
  role: ModelsFlatMemberInviteRole;
};
export type ModelsFlatMemberInviteData = {
  invite?: ModelsFlatMemberInviteResponse;
};
export type ModelsFlatMemberInviteApiResponse = {
  data?: ModelsFlatMemberInviteData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatResidentData = {
  resident?: ModelsFlatResidentResponse;
};
export type ModelsFlatResidentApiResponse = {
  data?: ModelsFlatResidentData;
  message?: string;
  success?: boolean;
};
export type ModelsAddFlatResidentRequest = {
  is_primary?: boolean;
  metadata?: {
    [key: string]: any;
  };
  role: ModelsFlatResidentRole;
};
export type ModelsUpdateFlatResidentRoleRequest = {
  role: ModelsFlatResidentRole;
};
export type ModelsFlatVisitorContextResident = {
  full_name?: string;
  id?: number;
  user_id?: number;
};
export type ModelsFlatRecentVisitorSummary = {
  entry_id?: number;
  full_name?: string;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
  visited_on?: string;
};
export type ModelsVisitorApprovalMode = "mandatory" | "optional" | "hybrid";
export type ModelsFlatVisitorSettingsResponse = {
  approval_required?: boolean;
  created_at?: string;
  default_visit_duration_minutes?: number;
  flat_id?: number;
  id?: number;
  is_enabled?: boolean;
  purpose?: ModelsVisitorPurpose;
  society_id?: number;
  updated_at?: string;
  updated_by?: number;
};
export type ModelsFlatVisitorContextResponse = {
  inherits_society_mode?: boolean;
  occupancy_status?: ModelsFlatStatus;
  primary_resident?: ModelsFlatVisitorContextResident;
  recent_visitors?: ModelsFlatRecentVisitorSummary[];
  society_approval_mode?: ModelsVisitorApprovalMode;
  total_residents?: number;
  visitor_settings?: ModelsFlatVisitorSettingsResponse[];
};
export type ModelsFlatVisitorContextData = {
  context?: ModelsFlatVisitorContextResponse;
};
export type ModelsFlatVisitorContextApiResponse = {
  data?: ModelsFlatVisitorContextData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntriesData = {
  entries?: ModelsVisitorEntry[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsVisitorEntriesApiResponse = {
  data?: ModelsVisitorEntriesData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorInviteEntry = {
  checked_in_at?: string;
  checked_out_at?: string;
  id?: number;
  qr_expires_at?: string;
  qr_used_at?: string;
  status?: ModelsVisitorStatus;
};
export type ModelsVisitorInviteHistoryItem = {
  created_at?: string;
  created_by?: number;
  entry?: ModelsVisitorInviteEntry;
  expires_at?: string;
  flat_id?: number;
  id?: number;
  metadata?: {
    [key: string]: any;
  };
  purpose?: ModelsVisitorPurpose;
  short_link?: ModelsShortLinkResponse;
  society_id?: number;
  status?: ModelsVisitorInviteStatus;
  updated_at?: string;
  used_at?: string;
  visitor?: ModelsVisitorSummary;
};
export type ModelsVisitorInviteHistoryData = {
  invites?: ModelsVisitorInviteHistoryItem[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsVisitorInviteHistoryApiResponse = {
  data?: ModelsVisitorInviteHistoryData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorInvite = {
  created_at?: string;
  created_by?: number;
  expires_at?: string;
  flat_id?: number;
  id?: number;
  metadata?: {
    [key: string]: any;
  };
  purpose?: ModelsVisitorPurpose;
  society_id?: number;
  status?: ModelsVisitorInviteStatus;
  updated_at?: string;
  used_at?: string;
};
export type ModelsVisitorInviteTokenData = {
  invite?: ModelsVisitorInvite;
  link?: ModelsShortLinkResponse;
  short_link?: ModelsShortLinkResponse;
  token?: ModelsQrTokenResponse;
};
export type ModelsVisitorInviteTokenApiResponse = {
  data?: ModelsVisitorInviteTokenData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateVisitorInviteRequest = {
  expires_at?: string;
  purpose?: ModelsVisitorPurpose;
};
export type ModelsVisitorInviteHistoryItemData = {
  invite?: ModelsVisitorInviteHistoryItem;
};
export type ModelsVisitorInviteHistoryItemApiResponse = {
  data?: ModelsVisitorInviteHistoryItemData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatVisitorSettingsData = {
  visitor_settings?: ModelsFlatVisitorSettingsResponse[];
};
export type ModelsFlatVisitorSettingsApiResponse = {
  data?: ModelsFlatVisitorSettingsData;
  message?: string;
  success?: boolean;
};
export type ModelsFlatVisitorSettingData = {
  visitor_setting?: ModelsFlatVisitorSettingsResponse;
};
export type ModelsFlatVisitorSettingApiResponse = {
  data?: ModelsFlatVisitorSettingData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateFlatVisitorSettingRequest = {
  approval_required?: boolean;
  default_visit_duration_minutes?: number;
  is_enabled?: boolean;
};
export type ModelsVisitorPendingEntry = {
  approved_at?: string;
  approved_by?: number;
  approver_name?: string;
  approver_phone_number?: string;
  auto_closed_at?: string;
  checked_in_at?: string;
  checked_out_at?: string;
  companion_details?: {
    [key: string]: any;
  }[];
  companions_count?: number;
  created_at?: string;
  created_by?: number;
  delivery_partner?: string;
  expected_at?: string;
  expected_checkout_at?: string;
  flat?: ModelsVisitorFlatSummary;
  flat_id?: number;
  handled_by_guard_id?: number;
  handled_by_guard_name?: string;
  handled_by_guard_phone?: string;
  id?: number;
  invite_id?: number;
  metadata?: {
    [key: string]: any;
  };
  notes?: string;
  primary_resident_id?: number;
  primary_resident_name?: string;
  purpose?: ModelsVisitorPurpose;
  qr_expires_at?: string;
  qr_used_at?: string;
  rejected_by?: number;
  rejection_reason?: string;
  service_provider?: string;
  society_id?: number;
  source?: ModelsVisitorEntrySource;
  status?: ModelsVisitorStatus;
  updated_at?: string;
  vehicle_number?: string;
  vehicle_type?: ModelsVisitorVehicleType;
  visitor?: ModelsVisitorSummary;
  visitor_id?: number;
  waiting_since?: string;
};
export type ModelsGuardDeskVisitorSettingsSummary = {
  allow_guard_entry?: boolean;
  allow_guard_on_behalf_approval?: boolean;
};
export type ModelsGuardDeskBootstrapResponse = {
  expected_guests_count?: number;
  pending_preview?: ModelsVisitorPendingEntry[];
  society?: ModelsSocietyResponse;
  stats?: ModelsVisitorEntryStatsResponse;
  visitor_settings?: ModelsGuardDeskVisitorSettingsSummary;
  waiting_at_gate_count?: number;
};
export type ModelsGuardDeskBootstrapData = {
  desk?: ModelsGuardDeskBootstrapResponse;
};
export type ModelsGuardDeskBootstrapApiResponse = {
  data?: ModelsGuardDeskBootstrapData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateGuardResponse = {
  member?: ModelsSocietyMemberResponse;
  user?: ModelsUserResponse;
};
export type ModelsGuardData = {
  guard?: ModelsCreateGuardResponse;
};
export type ModelsGuardApiResponse = {
  data?: ModelsGuardData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateGuardRequest = {
  email: string;
  first_name: string;
  last_name?: string;
  password: string;
  phone_number: string;
};
export type ModelsSocietyMemberData = {
  member?: ModelsSocietyMemberResponse;
};
export type ModelsSocietyMemberApiResponse = {
  data?: ModelsSocietyMemberData;
  message?: string;
  success?: boolean;
};
export type ModelsAddSocietyMemberRequest = {
  metadata?: {
    [key: string]: any;
  };
  role: ModelsSocietyMemberRole;
  society_id: number;
  user_id: number;
};
export type ModelsSocietyMemberSummaryResponse = {
  active_members?: number;
  admins?: number;
  owners?: number;
  pending_members?: number;
  removed_members?: number;
  residents?: number;
  staff?: number;
  suspended_members?: number;
  total_members?: number;
};
export type ModelsSocietyMemberSummaryData = {
  summary?: ModelsSocietyMemberSummaryResponse;
};
export type ModelsSocietyMemberSummaryApiResponse = {
  data?: ModelsSocietyMemberSummaryData;
  message?: string;
  success?: boolean;
};
export type ModelsMemberVisitorApprovalStatsResponse = {
  approved_count?: number;
  rejected_count?: number;
};
export type ModelsMemberVisitorApprovalStatsData = {
  stats?: ModelsMemberVisitorApprovalStatsResponse;
};
export type ModelsMemberVisitorApprovalStatsApiResponse = {
  data?: ModelsMemberVisitorApprovalStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsSocietyReasonRequest = {
  reason: string;
};
export type ModelsChangeSocietyMemberRoleRequest = {
  role: ModelsSocietyMemberRole;
  society_id: number;
  user_id: number;
};
export type ModelsSocietyOnboardingBootstrapResponse = {
  flat_count?: number;
  has_active_subscription?: boolean;
  has_flats?: boolean;
  has_staff?: boolean;
  is_onboarded?: boolean;
  missing_steps?: string[];
  next_path?: string;
  society?: ModelsSocietyResponse;
  staff_count?: number;
};
export type ModelsSocietyOnboardingBootstrapData = {
  onboarding?: ModelsSocietyOnboardingBootstrapResponse;
};
export type ModelsSocietyOnboardingBootstrapApiResponse = {
  data?: ModelsSocietyOnboardingBootstrapData;
  message?: string;
  success?: boolean;
};
export type ModelsSubscriptionData = {
  subscription?: ModelsSocietySubscriptionResponse;
};
export type ModelsSubscriptionApiResponse = {
  data?: ModelsSubscriptionData;
  message?: string;
  success?: boolean;
};
export type ModelsCreateTrialSubscriptionRequest = {
  ends_at?: string;
  metadata?: {
    [key: string]: any;
  };
  starts_at: string;
  trial_ends_at: string;
};
export type ModelsTransferOwnershipRequest = {
  new_owner_user_id: number;
};
export type ModelsVisitorPendingEntriesData = {
  entries?: ModelsVisitorPendingEntry[];
  limit?: number;
  offset?: number;
  total?: number;
};
export type ModelsVisitorPendingEntriesApiResponse = {
  data?: ModelsVisitorPendingEntriesData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntryStatsData = {
  stats?: ModelsVisitorEntryStatsResponse;
};
export type ModelsVisitorEntryStatsApiResponse = {
  data?: ModelsVisitorEntryStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsVisitorEntryDailyStatsResponse = {
  daily?: ModelsVisitorDailyCountResponse[];
  days?: number;
  metric?: string;
  timezone?: string;
  total?: number;
};
export type ModelsVisitorEntryDailyStatsData = {
  stats?: ModelsVisitorEntryDailyStatsResponse;
};
export type ModelsVisitorEntryDailyStatsApiResponse = {
  data?: ModelsVisitorEntryDailyStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateGuardVisitorEntryRequest = {
  companion_details?: {
    [key: string]: any;
  }[];
  companions_count?: number;
  email?: string;
  flat_id?: number;
  full_name?: string;
  notes?: string;
  phone_number?: string;
  photo_url?: string;
  vehicle_number?: string;
  vehicle_type?: ModelsVisitorVehicleType;
};
export type ModelsGuardApproveEntryRequest = {
  on_behalf?: boolean;
  reason?: string;
};
export type ModelsVisitorEventType =
  | "created"
  | "approved"
  | "rejected"
  | "checked_in"
  | "checked_out"
  | "cancelled"
  | "expired"
  | "auto_closed"
  | "qr_generated"
  | "qr_used"
  | "guard_approved_on_behalf";
export type ModelsVisitorEntryEvent = {
  actor_user_id?: number;
  created_at?: string;
  event_type?: ModelsVisitorEventType;
  id?: number;
  message?: string;
  metadata?: {
    [key: string]: any;
  };
  society_id?: number;
  visitor_entry_id?: number;
};
export type ModelsVisitorEntryEventsData = {
  events?: ModelsVisitorEntryEvent[];
};
export type ModelsVisitorEntryEventsApiResponse = {
  data?: ModelsVisitorEntryEventsData;
  message?: string;
  success?: boolean;
};
export type ModelsApiResponse = {
  data?: any;
  error?: ModelsErrorData;
  message?: string;
  success?: boolean;
};
export type ModelsRejectVisitorEntryRequest = {
  reason?: string;
};
export type ModelsSocietyVisitorSettingsResponse = {
  allow_guard_entry?: boolean;
  allow_guard_on_behalf_approval?: boolean;
  allow_public_qr_entry?: boolean;
  allow_resident_pre_approval?: boolean;
  approval_mode?: ModelsVisitorApprovalMode;
  created_at?: string;
  default_visit_duration_minutes?: number;
  grace_period_minutes?: number;
  id?: number;
  is_active?: boolean;
  qr_expiry_minutes?: number;
  society_id?: number;
  updated_at?: string;
  updated_by?: number;
};
export type ModelsSocietyVisitorSettingsData = {
  visitor_settings?: ModelsSocietyVisitorSettingsResponse;
};
export type ModelsSocietyVisitorSettingsApiResponse = {
  data?: ModelsSocietyVisitorSettingsData;
  message?: string;
  success?: boolean;
};
export type ModelsUpdateSocietyVisitorSettingsRequest = {
  allow_guard_entry?: boolean;
  allow_guard_on_behalf_approval?: boolean;
  allow_public_qr_entry?: boolean;
  allow_resident_pre_approval?: boolean;
  approval_mode?: ModelsVisitorApprovalMode;
  default_visit_duration_minutes?: number;
  grace_period_minutes?: number;
  is_active?: boolean;
  qr_expiry_minutes?: number;
};
export type ModelsSocietyFlatVisitorSettingRow = {
  approval_required?: boolean;
  block?: string;
  default_visit_duration_minutes?: number;
  flat_id?: number;
  flat_number?: string;
  is_enabled?: boolean;
  purpose?: ModelsVisitorPurpose;
};
export type ModelsSocietyFlatVisitorSettingsData = {
  limit?: number;
  offset?: number;
  settings?: ModelsSocietyFlatVisitorSettingRow[];
  total?: number;
};
export type ModelsSocietyFlatVisitorSettingsApiResponse = {
  data?: ModelsSocietyFlatVisitorSettingsData;
  message?: string;
  success?: boolean;
};
export type ModelsSubscriptionsData = {
  subscriptions?: ModelsSocietySubscriptionResponse[];
};
export type ModelsSubscriptionsApiResponse = {
  data?: ModelsSubscriptionsData;
  message?: string;
  success?: boolean;
};
export type ModelsSubscriptionStatsData = {
  stats?: ModelsSubscriptionStatsResponse;
};
export type ModelsSubscriptionStatsApiResponse = {
  data?: ModelsSubscriptionStatsData;
  message?: string;
  success?: boolean;
};
export type ModelsActivateSubscriptionRequest = {
  ends_at: string;
  metadata?: {
    [key: string]: any;
  };
  starts_at: string;
};
export type ModelsCancelSubscriptionRequest = {
  metadata?: {
    [key: string]: any;
  };
  reason: string;
};
export type ModelsRenewSubscriptionRequest = {
  ends_at: string;
  metadata?: {
    [key: string]: any;
  };
  starts_at: string;
};
export type ModelsNotificationDetailApiResponse = {
  data?: ModelsNotification;
  message?: string;
  success?: boolean;
};
export const {
  useGetHealthQuery,
  useGetHealthLiveQuery,
  useGetHealthReadyQuery,
  usePostV1AuthChangePasswordMutation,
  usePostV1AuthForgotPasswordMutation,
  usePostV1AuthLoginMutation,
  usePostV1AuthLogoutMutation,
  useGetV1AuthProfileQuery,
  usePatchV1AuthProfileMutation,
  useGetV1AuthProfileAvatarQuery,
  usePutV1AuthProfileAvatarMutation,
  useDeleteV1AuthProfileAvatarMutation,
  usePostV1AuthRefreshMutation,
  usePostV1AuthRegisterMutation,
  usePostV1AuthResendOtpMutation,
  usePostV1AuthResetPasswordMutation,
  usePostV1AuthResidentRegisterMutation,
  usePostV1AuthVerifyOtpMutation,
  useGetV1BootstrapQuery,
  useGetV1DeveloperDashboardBootstrapQuery,
  useGetV1FlatClaimsQuery,
  usePostV1FlatClaimsMutation,
  useGetV1FlatClaimsByClaimIdQuery,
  usePostV1FlatClaimsByClaimIdCancelMutation,
  useGetV1FlatResidentsQuery,
  useGetV1FlatsQuery,
  usePostV1MeDeviceTokensMutation,
  useDeleteV1MeDeviceTokensMutation,
  useGetV1MeFlatClaimsQuery,
  useGetV1MeNotificationsQuery,
  usePostV1MeNotificationsReadAllMutation,
  useGetV1MeNotificationsUnreadCountQuery,
  usePatchV1MeNotificationsByNotificationIdReadMutation,
  useGetV1MeResidencesQuery,
  useGetV1PlansQuery,
  usePostV1PlansMutation,
  useGetV1PlansLookupQuery,
  usePatchV1PlansByPlanIdMutation,
  usePostV1PlansByPlanIdActivateMutation,
  usePostV1PlansByPlanIdDeactivateMutation,
  useGetV1PublicFlatMemberInvitesByTokenQuery,
  usePostV1PublicFlatMemberInvitesByTokenAcceptMutation,
  usePostV1PublicFlatMemberInvitesByTokenJoinMutation,
  useGetV1PublicShortLinksByShortCodeQuery,
  usePostV1PublicShortLinksByShortCodeMemberAcceptMutation,
  usePostV1PublicShortLinksByShortCodeMemberJoinMutation,
  usePostV1PublicShortLinksByShortCodeVisitorSubmitMutation,
  useGetV1PublicSocietiesBySocietyCodeClaimOptionsQuery,
  usePostV1PublicSocietiesBySocietyCodeVisitorEntriesPublicQrMutation,
  usePostV1PublicSocietiesBySocietyCodeVisitorEntriesQuickLinkMutation,
  useGetV1PublicSocietiesBySocietyCodeVisitorEntryOptionsQuery,
  usePostV1PublicSocietiesBySocietyCodeVisitorInvitesAndTokenSubmitMutation,
  usePostV1PublicVisitorEntriesQrValidateMutation,
  useGetV1PublicVisitorInvitesByTokenQuery,
  useGetV1SocietiesQuery,
  usePostV1SocietiesMutation,
  useGetV1SocietiesMyQuery,
  useGetV1SocietiesBySocietyIdQuery,
  useDeleteV1SocietiesBySocietyIdMutation,
  usePatchV1SocietiesBySocietyIdMutation,
  useGetV1SocietiesBySocietyIdAllmemberQuery,
  usePostV1SocietiesBySocietyIdApproveMutation,
  useGetV1SocietiesBySocietyIdDashboardBootstrapQuery,
  useGetV1SocietiesBySocietyIdFlatClaimsQuery,
  useGetV1SocietiesBySocietyIdFlatClaimsStatsQuery,
  useGetV1SocietiesBySocietyIdFlatClaimsAndClaimIdQuery,
  usePostV1SocietiesBySocietyIdFlatClaimsAndClaimIdApproveMutation,
  usePostV1SocietiesBySocietyIdFlatClaimsAndClaimIdRejectMutation,
  useGetV1SocietiesBySocietyIdFlatsQuery,
  usePostV1SocietiesBySocietyIdFlatsMutation,
  usePostV1SocietiesBySocietyIdFlatsBulkMutation,
  usePostV1SocietiesBySocietyIdFlatsGenerateMutation,
  useGetV1SocietiesBySocietyIdFlatsStatsQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdQuery,
  useDeleteV1SocietiesBySocietyIdFlatsAndFlatIdMutation,
  usePatchV1SocietiesBySocietyIdFlatsAndFlatIdMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdBlockMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdCancelMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMembersQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsUsersUserIdMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdQuery,
  useDeleteV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdMoveOutMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdPrimaryMutation,
  usePatchV1SocietiesBySocietyIdFlatsAndFlatIdResidentsResidentIdRoleMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdUnblockMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesPendingQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesEntryIdPhotoQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffMutation,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetMutation,
  usePatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeMutation,
  useGetV1SocietiesBySocietyIdGuardDeskBootstrapQuery,
  usePostV1SocietiesBySocietyIdGuardsMutation,
  useGetV1SocietiesBySocietyIdMembersQuery,
  usePostV1SocietiesBySocietyIdMembersMutation,
  useGetV1SocietiesBySocietyIdMembersSummaryQuery,
  useGetV1SocietiesBySocietyIdMembersAndMemberIdQuery,
  useGetV1SocietiesBySocietyIdMembersAndMemberIdVisitorApprovalStatsQuery,
  useDeleteV1SocietiesBySocietyIdMembersAndUserIdMutation,
  usePostV1SocietiesBySocietyIdMembersAndUserIdReactivateMutation,
  usePatchV1SocietiesBySocietyIdMembersAndUserIdRoleMutation,
  usePostV1SocietiesBySocietyIdMembersAndUserIdSuspendMutation,
  useGetV1SocietiesBySocietyIdOnboardingBootstrapQuery,
  usePostV1SocietiesBySocietyIdReactivateMutation,
  usePostV1SocietiesBySocietyIdRejectMutation,
  usePostV1SocietiesBySocietyIdRestoreMutation,
  usePostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdPendingMutation,
  usePostV1SocietiesBySocietyIdSubscriptionsPlansAndPlanIdTrialMutation,
  usePostV1SocietiesBySocietyIdSuspendMutation,
  usePostV1SocietiesBySocietyIdTransferOwnershipMutation,
  useGetV1SocietiesBySocietyIdVisitorEntriesQuery,
  usePostV1SocietiesBySocietyIdVisitorEntriesCheckInMutation,
  useGetV1SocietiesBySocietyIdVisitorEntriesExpectedGuestsQuery,
  usePostV1SocietiesBySocietyIdVisitorEntriesGuardMutation,
  useGetV1SocietiesBySocietyIdVisitorEntriesPendingQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesStatsQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesStatsDailyQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesWaitingAtGateQuery,
  useGetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdQuery,
  usePatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdApproveAndCheckInMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckInMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutMutation,
  useGetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdEventsQuery,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdGuardApproveMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdNotifyMutation,
  useGetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoQuery,
  usePutV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoMutation,
  useDeleteV1SocietiesBySocietyIdVisitorEntriesAndEntryIdPhotoMutation,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdRejectMutation,
  usePostV1SocietiesBySocietyIdVisitorInvitesAndInviteIdCancelMutation,
  useGetV1SocietiesBySocietyIdVisitorInvitesAndTokenQuery,
  useGetV1SocietiesBySocietyIdVisitorSettingsQuery,
  usePatchV1SocietiesBySocietyIdVisitorSettingsMutation,
  useGetV1SocietiesBySocietyIdVisitorSettingsFlatsQuery,
  useGetV1SubscriptionsQuery,
  useGetV1SubscriptionsLookupQuery,
  useGetV1SubscriptionsStatsQuery,
  usePostV1SubscriptionsBySubscriptionIdActivateMutation,
  usePostV1SubscriptionsBySubscriptionIdCancelMutation,
  usePostV1SubscriptionsBySubscriptionIdExpireMutation,
  usePostV1SubscriptionsBySubscriptionIdPlansAndPlanIdMutation,
  usePostV1SubscriptionsBySubscriptionIdRenewMutation,
  useGetV1MeNotificationsByNotificationIdQuery,
} = injectedRtkApi;
