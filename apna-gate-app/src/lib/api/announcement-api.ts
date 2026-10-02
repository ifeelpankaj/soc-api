import { enhancedApi } from "@/lib/api/enhanced-api";

type Envelope<T> = { data?: T };
type Channel = { id: number; type: string };
export type Announcement = {
  id: number;
  title?: string;
  body: string;
  author?: { name: string };
  created_at: string;
  is_pinned: boolean;
  is_important: boolean;
};
type Page = { items: Announcement[]; pinned?: Announcement[]; next_cursor?: string };

const announcementApi = enhancedApi.injectEndpoints({
  endpoints: (build) => ({
    residentAnnouncementChannels: build.query<Envelope<Channel[]>, number>({
      query: (societyId) => `/v1/societies/${societyId}/channels`,
    }),
    residentAnnouncements: build.query<Envelope<Page>, { societyId: number; channelId: number; cursor?: string }>({
      query: ({ societyId, channelId, cursor }) => ({
        url: `/v1/societies/${societyId}/channels/${channelId}/posts`,
        params: { limit: 20, cursor },
      }),
    }),
  }),
});

export const { useResidentAnnouncementChannelsQuery, useResidentAnnouncementsQuery } = announcementApi;
