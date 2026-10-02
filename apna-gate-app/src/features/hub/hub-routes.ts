import type { Href } from "expo-router";

import { residentDashboardRoute } from "@/features/resident/resident-routes";

export type HubPostReturnTo = "home" | "announcements" | "community" | "hub";

export function residentHubRoute(): Href {
  return "/resident/hub" as Href;
}

export function residentHubAnnouncementsRoute(): Href {
  return "/resident/hub/announcements" as Href;
}

export function residentHubCommunityRoute(): Href {
  return "/resident/hub/community" as Href;
}

export function residentHubCreatePostRoute(channelType: "community"): Href {
  return {
    pathname: "/resident/hub/create",
    params: { channelType },
  } as Href;
}

export function residentHubEditPostRoute(postId: number, returnTo?: HubPostReturnTo): Href {
  return {
    pathname: "/resident/hub/edit/[postId]",
    params: {
      postId: String(postId),
      returnTo,
    },
  } as Href;
}

export function residentHubPostRoute(
  societyId: number,
  postId: number,
  channelId?: number,
  returnTo?: HubPostReturnTo,
): Href {
  return {
    pathname: "/resident/hub/posts/[postId]",
    params: {
      postId: String(postId),
      societyId: String(societyId),
      channelId: channelId ? String(channelId) : undefined,
      returnTo,
    },
  } as Href;
}

export function resolveHubPostBackRoute(returnTo?: string | string[]): Href {
  const value = Array.isArray(returnTo) ? returnTo[0] : returnTo;
  switch (value) {
    case "home":
      return residentDashboardRoute();
    case "announcements":
      return residentHubAnnouncementsRoute();
    case "community":
      return residentHubCommunityRoute();
    case "hub":
    default:
      return residentHubRoute();
  }
}
