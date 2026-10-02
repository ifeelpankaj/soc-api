import { enhancedApi } from "@/lib/api/enhanced-api";

export type AppNotificationData = {
	  schema_version?: number | string;
	  entity_type?: string;
	  entity_id?: string;
  category_id?: "visitor_decision" | "notification_info";
  event?: string;
  entry_id?: string;
  flat_id?: string;
  invite_id?: string;
  invite_type?: "visitor" | "member";
  notification_id?: string;
  occurred_at?: string;
  society_id?: string;
  type?: string;
  [key: string]: unknown;
};

export type AppNotification = {
  id: string;
  user_id: number;
  society_id?: number;
  flat_id?: number;
  type: string;
  domain?: "visitor" | "maintenance" | "hub" | "system";
  title: string;
  body: string;
  data: AppNotificationData;
  event_key?: string;
  read_at?: string;
  created_at: string;
};

type ApiEnvelope<T> = {
  data?: T;
  message?: string;
  success?: boolean;
};

type NotificationsPayload = {
  items?: AppNotification[];
  next_cursor?: string;
};

type UnreadCountPayload = {
  unread_count: number;
};

type ReadNotificationPayload = {
  id: string;
  read_at?: string;
};

export type NotificationPushPreferences = {
  visitor_updates_push: boolean;
  maintenance_push: boolean;
  announcements_push: boolean;
  hub_replies_push: boolean;
};

export const notificationApiExtensions = enhancedApi.injectEndpoints({
  overrideExisting: true,
  endpoints: (build) => ({
    getNotificationPushPreferences: build.query<ApiEnvelope<NotificationPushPreferences>, number>({
      query: (societyId) => ({url: "/v1/me/notifications/preferences",params:{society_id:societyId}}),
      providesTags: ["Notifications"],
    }),
    putNotificationPushPreferences: build.mutation<ApiEnvelope<NotificationPushPreferences>, {societyId:number; preferences:NotificationPushPreferences}>({
      query: ({societyId,preferences}) => ({url: "/v1/me/notifications/preferences",method:"PUT",params:{society_id:societyId},body:preferences}),
      invalidatesTags: ["Notifications"],
    }),
    getOwnedNotification: build.query<ApiEnvelope<AppNotification>, string>({
      query: (id) => ({url: `/v1/me/notifications/${encodeURIComponent(id)}`}),
    }),
    getV1MeNotifications: build.query<
      ApiEnvelope<NotificationsPayload>,
      { limit?: number; cursor?: string } | void
    >({
      query: (arg) => ({
        url: "/v1/me/notifications",
        params: {
          limit: arg?.limit,
          cursor: arg?.cursor,
        },
      }),
      providesTags: ["Notifications"],
    }),
    getV1MeNotificationsUnreadCount: build.query<
      ApiEnvelope<UnreadCountPayload>,
      void
    >({
      query: () => ({ url: "/v1/me/notifications/unread-count" }),
      providesTags: ["Notifications"],
    }),
    patchV1MeNotificationsAndNotificationIdRead: build.mutation<
      ApiEnvelope<ReadNotificationPayload>,
      { notificationId: string }
    >({
      query: ({ notificationId }) => ({
        url: `/v1/me/notifications/${notificationId}/read`,
        method: "PATCH",
      }),
      invalidatesTags: ["Notifications"],
    }),
    postV1MeNotificationsReadAll: build.mutation<
      ApiEnvelope<UnreadCountPayload>,
      void
    >({
      query: () => ({
        url: "/v1/me/notifications/read-all",
        method: "POST",
      }),
      invalidatesTags: ["Notifications"],
    }),
  }),
});

export const {
  useGetNotificationPushPreferencesQuery,
  usePutNotificationPushPreferencesMutation,
  useGetV1MeNotificationsQuery,
  useGetV1MeNotificationsUnreadCountQuery,
  usePatchV1MeNotificationsAndNotificationIdReadMutation,
  usePostV1MeNotificationsReadAllMutation,
} = notificationApiExtensions;
