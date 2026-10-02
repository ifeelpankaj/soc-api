import type { Href } from "expo-router";

export type ResidentEntriesPreset = "expected" | "inside" | "all";
export type LegacyResidentEntriesPreset =
  ResidentEntriesPreset | "today" | "recent";

export function residentDashboardRoute(): Href {
  return "/resident/dashboard" as Href;
}

export function residentProfileRoute(): Href {
  return "/resident/profile" as Href;
}

export function residentMaintenanceRoute(): Href {
  return "/resident/maintenance" as Href;
}

export function residentAnnouncementsRoute(): Href {
  return "/resident/announcements" as Href;
}

export function residentMaintenanceBillRoute(
  billId: number,
  societyId: number,
  flatId: number,
): Href {
  return {
    pathname: `/resident/maintenance/bills/${billId}`,
    params: { societyId, flatId },
  } as Href;
}

export function residentMaintenancePaymentRoute(
  billId: number,
  societyId: number,
  flatId: number,
): Href {
  return {
    pathname: `/resident/maintenance/bills/${billId}/pay`,
    params: { societyId, flatId },
  } as Href;
}

export function residentLogsRoute(): Href {
  return "/resident/logs" as Href;
}

export function residentEntriesRoute(
  preset: ResidentEntriesPreset = "expected",
): Href {
  if (preset === "expected") {
    return "/resident/entries" as Href;
  }
  return {
    pathname: "/resident/entries",
    params: { preset },
  } as unknown as Href;
}

export function residentEntryDetailRoute(entryId: number | string): Href {
  return `/resident/entries/${entryId}` as Href;
}

export function residentVisitorsRoute(): Href {
  return "/resident/visitors" as Href;
}

export function residentVisitorSettingsRoute(): Href {
  return "/resident/visitors/settings" as Href;
}

export function residentVisitorInviteRoute(): Href {
  return "/resident/visitors/invite" as Href;
}

export function residentInvitesRoute(): Href {
  return "/resident/invites" as Href;
}

export function residentInviteDetailRoute(
  type: "member" | "visitor",
  inviteId: number | string,
): Href {
  return `/resident/invites/${type}/${inviteId}` as Href;
}

export function residentMembersRoute(): Href {
  return "/resident/members" as Href;
}

export function residentMembersAddRoute(): Href {
  return "/resident/members/add" as Href;
}

export function parseResidentEntriesPreset(
  value?: string | string[],
): ResidentEntriesPreset {
  const raw = Array.isArray(value) ? value[0] : value;
  if (raw === "expected" || raw === "inside" || raw === "all") {
    return raw;
  }
  // "today" and "recent" are stale resident entry presets. They fall back to
  // Expected for compatibility only; they are not resident runtime tabs.
  return "expected";
}
