import { enhancedApi } from "@/lib/api/enhanced-api";

type Envelope<T> = { data?: T };

export type HubReactionType = "like" | "love" | "helpful";

export type HubAuthor = { id?: number; name?: string; flat_label?: string };
export type HubContent = {
  id: number;
  post_id?: number;
  channel_id?: number;
  society_id?: number;
  title?: string;
  body: string;
  author?: HubAuthor;
  created_at: string;
  updated_at?: string;
  edited_at?: string;
  category_id?: number;
  attachment_ids?: number[];
  is_pinned?: boolean;
  is_important?: boolean;
  comments_enabled?: boolean;
  reply_count?: number;
  reactions?: Record<string, number>;
  my_reactions?: string[];
  status?: string;
};

export type HubUploadView = {
  id: number;
  filename?: string;
  mime_type?: string;
  size?: number;
  expires_at?: string;
};

export type HubAttachmentView = HubUploadView & {
  url: string;
};

export type HubChannelSummary = {
  id: number;
  type: string;
  name?: string;
  description?: string;
  unread_count?: number;
  important_announcement?: HubContent | null;
};

export type HubPage = {
  items: HubContent[];
  pinned?: HubContent[];
  next_cursor?: string;
};

export type HubCategory = { id: number; name: string; slug?: string; code?: string };

export type HubPostInput = {
  title?: string;
  body: string;
  category_id?: number;
  comments_enabled?: boolean;
  attachment_ids?: number[];
};

export type HubCommentInput = {
  body: string;
  parent_id?: number;
  attachment_ids?: number[];
};

const hubTag = (societyId: number) => [{ type: "Hub" as const, id: societyId }];

const hubApi = enhancedApi.enhanceEndpoints({ addTagTypes: ["Hub"] }).injectEndpoints({
  endpoints: (build) => ({
    hubChannels: build.query<HubChannelSummary[], number>({
      query: (societyId) => `/v1/societies/${societyId}/channels`,
      transformResponse: (response: Envelope<HubChannelSummary[]>) =>
        response.data ?? [],
      providesTags: (_result, _error, societyId) => hubTag(societyId),
    }),
    hubPosts: build.query<
      HubPage,
      { societyId: number; channelId: number; cursor?: string; categoryId?: number }
    >({
      query: ({ societyId, channelId, cursor, categoryId }) => ({
        url: `/v1/societies/${societyId}/channels/${channelId}/posts`,
        params: {
          limit: 20,
          cursor,
          category_id: categoryId,
        },
      }),
      transformResponse: (response: Envelope<HubPage>) =>
        response.data ?? { items: [] },
      providesTags: (_result, _error, arg) => [
        ...hubTag(arg.societyId),
        { type: "Hub" as const, id: `${arg.societyId}:${arg.channelId}` },
      ],
    }),
    hubPost: build.query<HubContent, { societyId: number; postId: number }>({
      query: ({ societyId, postId }) =>
        `/v1/societies/${societyId}/posts/${postId}`,
      transformResponse: (response: Envelope<HubContent>) => {
        if (!response.data) {
          throw new Error("Post not found");
        }
        return response.data;
      },
      providesTags: (_result, _error, arg) => [
        ...hubTag(arg.societyId),
        { type: "Hub" as const, id: `post:${arg.postId}` },
      ],
    }),
    hubComments: build.query<
      HubPage,
      { societyId: number; postId: number; cursor?: string }
    >({
      query: ({ societyId, postId, cursor }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}/comments`,
        params: { limit: 30, cursor },
      }),
      transformResponse: (response: Envelope<HubPage>) =>
        response.data ?? { items: [] },
      providesTags: (_result, _error, arg) => [
        { type: "Hub" as const, id: `comments:${arg.postId}` },
      ],
    }),
    hubCategories: build.query<HubCategory[], number>({
      query: (societyId) => `/v1/societies/${societyId}/channels/categories`,
      transformResponse: (response: Envelope<HubCategory[]>) => response.data ?? [],
    }),
    hubMarkRead: build.mutation<
      void,
      { societyId: number; channelId: number; postId: number }
    >({
      query: ({ societyId, channelId, postId }) => ({
        url: `/v1/societies/${societyId}/channels/${channelId}/read`,
        method: "POST",
        body: { post_id: postId },
      }),
      invalidatesTags: (_result, _error, arg) => hubTag(arg.societyId),
    }),
    hubCreatePost: build.mutation<
      HubContent,
      { societyId: number; channelId: number; body: HubPostInput }
    >({
      query: ({ societyId, channelId, body }) => ({
        url: `/v1/societies/${societyId}/channels/${channelId}/posts`,
        method: "POST",
        body,
      }),
      transformResponse: (response: Envelope<HubContent>) => {
        if (!response.data) {
          throw new Error("Unable to publish post");
        }
        return response.data;
      },
      invalidatesTags: (_result, _error, arg) => [
        ...hubTag(arg.societyId),
        { type: "Hub" as const, id: `${arg.societyId}:${arg.channelId}` },
      ],
    }),
    hubCreateComment: build.mutation<
      HubContent,
      { societyId: number; postId: number; body: HubCommentInput }
    >({
      query: ({ societyId, postId, body }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}/comments`,
        method: "POST",
        body,
      }),
      transformResponse: (response: Envelope<HubContent>) => {
        if (!response.data) {
          throw new Error("Unable to post comment");
        }
        return response.data;
      },
      invalidatesTags: (_result, _error, arg) => [
        { type: "Hub" as const, id: `comments:${arg.postId}` },
        { type: "Hub" as const, id: `post:${arg.postId}` },
      ],
    }),
    hubToggleReaction: build.mutation<
      void,
      { societyId: number; postId: number; reactionType: HubReactionType }
    >({
      query: ({ societyId, postId, reactionType }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}/reactions`,
        method: "POST",
        body: { reaction_type: reactionType },
      }),
      invalidatesTags: (_result, _error, arg) => [
        { type: "Hub" as const, id: `post:${arg.postId}` },
        ...hubTag(arg.societyId),
      ],
    }),
    hubRemoveReaction: build.mutation<
      void,
      { societyId: number; postId: number; reactionType: HubReactionType }
    >({
      query: ({ societyId, postId, reactionType }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}/reactions/${reactionType}`,
        method: "DELETE",
      }),
      invalidatesTags: (_result, _error, arg) => [
        { type: "Hub" as const, id: `post:${arg.postId}` },
        ...hubTag(arg.societyId),
      ],
    }),
    hubUploadAttachment: build.mutation<
      HubUploadView,
      { societyId: number; body: FormData }
    >({
      query: ({ societyId, body }) => ({
        url: `/v1/societies/${societyId}/channel-uploads`,
        method: "POST",
        body,
      }),
      transformResponse: (response: Envelope<HubUploadView>) => {
        if (!response.data?.id) {
          throw new Error("Upload failed");
        }
        return response.data;
      },
    }),
    hubDeleteUpload: build.mutation<void, { societyId: number; uploadId: number }>({
      query: ({ societyId, uploadId }) => ({
        url: `/v1/societies/${societyId}/channel-uploads/${uploadId}`,
        method: "DELETE",
      }),
    }),
    hubAttachment: build.query<
      HubAttachmentView,
      { societyId: number; uploadId: number }
    >({
      query: ({ societyId, uploadId }) =>
        `/v1/societies/${societyId}/channel-attachments/${uploadId}`,
      transformResponse: (response: Envelope<HubAttachmentView>) => {
        if (!response.data?.url) {
          throw new Error("Attachment unavailable");
        }
        return response.data;
      },
      keepUnusedDataFor: 60 * 14,
    }),
    hubUpdatePost: build.mutation<
      HubContent,
      { societyId: number; postId: number; body: HubPostInput }
    >({
      query: ({ societyId, postId, body }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}`,
        method: "PATCH",
        body,
      }),
      transformResponse: (response: Envelope<HubContent>) => {
        if (!response.data) {
          throw new Error("Unable to update post");
        }
        return response.data;
      },
      invalidatesTags: (_result, _error, arg) => [
        ...hubTag(arg.societyId),
        { type: "Hub" as const, id: `post:${arg.postId}` },
      ],
    }),
    hubDeletePost: build.mutation<
      void,
      { societyId: number; postId: number; reason?: string }
    >({
      query: ({ societyId, postId, reason }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}`,
        method: "DELETE",
        body: reason ? { reason } : undefined,
      }),
      invalidatesTags: (_result, _error, arg) => [
        ...hubTag(arg.societyId),
        { type: "Hub" as const, id: `post:${arg.postId}` },
      ],
    }),
    hubReportPost: build.mutation<
      void,
      { societyId: number; postId: number; reason: string }
    >({
      query: ({ societyId, postId, reason }) => ({
        url: `/v1/societies/${societyId}/posts/${postId}/reports`,
        method: "POST",
        body: { reason },
      }),
    }),
  }),
});

export const {
  useHubChannelsQuery,
  useHubPostsQuery,
  useHubPostQuery,
  useHubCommentsQuery,
  useHubCategoriesQuery,
  useHubMarkReadMutation,
  useHubCreatePostMutation,
  useHubCreateCommentMutation,
  useHubToggleReactionMutation,
  useHubRemoveReactionMutation,
  useHubUploadAttachmentMutation,
  useHubDeleteUploadMutation,
  useHubAttachmentQuery,
  useHubUpdatePostMutation,
  useHubDeletePostMutation,
  useHubReportPostMutation,
} = hubApi;

export function findHubChannel(
  channels: HubChannelSummary[] | undefined,
  type: "announcement" | "community",
) {
  return channels?.find((channel) => channel.type === type);
}

export function reactionTotal(reactions?: Record<string, number>) {
  if (!reactions) {
    return 0;
  }
  return Object.values(reactions).reduce((sum, value) => sum + value, 0);
}
