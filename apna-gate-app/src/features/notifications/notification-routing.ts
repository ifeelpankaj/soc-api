import type { Href } from "expo-router";
import type { SymbolViewProps } from "expo-symbols";

import {
  guardEntriesRoute,
  guardEntryDetailRoute,
} from "@/features/guard/guard-routes";
import { residentHubPostRoute } from "@/features/hub/hub-routes";
import {
  residentEntryDetailRoute,
  residentEntriesRoute,
  residentInviteDetailRoute,
  residentVisitorsRoute,
} from "@/features/resident/resident-routes";
import type {
  AppNotification,
  AppNotificationData,
} from "@/lib/api/notification-api-extensions";

export type NotificationTone = "info" | "success" | "warning";

export function notificationsRoute(): Href {
  return "/notifications" as Href;
}

export function notificationSettingsRoute(): Href {
  return "/notification-settings" as Href;
}

export type NotificationPresentation = {
  body: string;
  icon: SymbolViewProps["name"];
  title: string;
  tone: NotificationTone;
};

const notificationIcons = {
  default: { ios: "bell", android: "notifications", web: "notifications" },
  memberInviteAccepted: {
    ios: "person.2.badge.gearshape",
    android: "group",
    web: "group",
  },
  visitorCheckIn: {
    ios: "checkmark.circle",
    android: "check_circle",
    web: "check_circle",
  },
  visitorCheckOut: {
    ios: "rectangle.portrait.and.arrow.right",
    android: "logout",
    web: "logout",
  },
  visitorApproved: {
    ios: "checkmark.shield",
    android: "verified_user",
    web: "verified_user",
  },
  visitorInviteAccepted: {
    ios: "person.crop.circle.badge.checkmark",
    android: "verified_user",
    web: "person_add_alt_1",
  },
  visitorPending: {
    ios: "bell.badge",
    android: "notifications_active",
    web: "notifications_active",
  },
  visitorRejected: { ios: "xmark.shield", android: "block", web: "block" },
} satisfies Record<string, SymbolViewProps["name"]>;

export function getNotificationType(
  data?: Record<string, unknown> | AppNotificationData,
  fallback?: string,
) {
  const type = stringValue(data?.type) ?? stringValue(data?.event) ?? fallback;
  return type;
}

export function notificationPresentation(
  notification: Pick<AppNotification, "body" | "title" | "type"> & {
    data?: AppNotificationData;
  },
): NotificationPresentation {
  const type = getNotificationType(notification.data, notification.type);
  switch (type) {
    case "visitor.pending":
    case "visitor.approval_requested":
    case "visitor_pending_approval":
      return {
        title: notificationTitle(
          notification.title,
          personTitle(
            notification.data,
            "Visitor is at the gate",
            "is at the gate",
          ),
          legacyVisitorPendingTitles,
          "Visitor approval needed",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorPendingBody,
          visitorPendingBody,
          legacyVisitorPendingBodies,
          legacyVisitorPendingTitles,
        ),
        icon: notificationIcons.visitorPending,
        tone: "warning",
      };
    case "visitor_invite.accepted":
    case "visitor_invite_accepted":
      return {
        title: notificationTitle(
          notification.title,
          personTitle(
            notification.data,
            "Guest accepted your invite",
            "accepted your invite",
          ),
          legacyVisitorInviteAcceptedTitles,
          "Guest invite accepted",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorInviteAcceptedBody,
          visitorInviteAcceptedBody,
          legacyVisitorInviteAcceptedBodies,
          legacyVisitorInviteAcceptedTitles,
        ),
        icon: notificationIcons.visitorInviteAccepted,
        tone: "success",
      };
    case "member_invite.accepted":
    case "member_invite_accepted":
      return {
        title: preferText(
          notification.title,
          memberJoinedTitle(notification.data),
        ),
        body: preferText(
          notification.body,
          memberJoinedBody(notification.data),
        ),
        icon: notificationIcons.memberInviteAccepted,
        tone: "success",
      };
    case "visitor.checkin":
    case "visitor_checked_in":
      return {
        title: notificationTitle(
          notification.title,
          personTitle(notification.data, "Visitor has arrived", "has arrived"),
          legacyVisitorCheckInTitles,
          "Visitor checked in",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorCheckInBody,
          visitorCheckInBody,
          legacyVisitorCheckInBodies,
          legacyVisitorCheckInTitles,
        ),
        icon: notificationIcons.visitorCheckIn,
        tone: "success",
      };
    case "visitor.checkout":
    case "visitor_checked_out":
      return {
        title: notificationTitle(
          notification.title,
          personTitle(notification.data, "Visitor has left", "has left"),
          legacyVisitorCheckOutTitles,
          "Visitor checked out",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorCheckOutBody,
          visitorCheckOutBody,
          legacyVisitorCheckOutBodies,
          legacyVisitorCheckOutTitles,
        ),
        icon: notificationIcons.visitorCheckOut,
        tone: "info",
      };
    case "visitor.approved":
    case "visitor_approved":
      return {
        title: notificationTitle(
          notification.title,
          "Visitor approved",
          legacyVisitorApprovedTitles,
          "Visitor approved",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorApprovedBody,
          (data) => visitorDecisionBody(data, "approved"),
          legacyVisitorApprovedBodies,
          legacyVisitorApprovedTitles,
        ),
        icon: notificationIcons.visitorApproved,
        tone: "success",
      };
    case "visitor.rejected":
    case "visitor_rejected":
      return {
        title: notificationTitle(
          notification.title,
          "Visitor request declined",
          legacyVisitorRejectedTitles,
          "Visitor request declined",
        ),
        body: notificationBody(
          notification.data,
          notification.title,
          notification.body,
          newVisitorRejectedBody,
          (data) => visitorDecisionBody(data, "declined"),
          legacyVisitorRejectedBodies,
          legacyVisitorRejectedTitles,
        ),
        icon: notificationIcons.visitorRejected,
        tone: "warning",
      };
    default:
      return {
        title: preferText(notification.title, "Notification"),
        body: notification.body,
        icon: notificationIcons.default,
        tone: "info",
      };
  }
}

export function notificationRoute(
  notification: AppNotification,
  homeRoute?: Href | null,
  webCheckout = false,
): Href | null {
  const data = notification.data ?? {};
  const version = Number(data.schema_version ?? 1);
  if (!Number.isSafeInteger(version) || version !== 1) {
    return notificationsRoute();
  }
  const type = getNotificationType(data, notification.type);
  const inviteId = stringValue(data.invite_id);
  const entryId = Number(stringValue(data.entry_id));

  if (
    (type === "visitor.checkin" || type === "visitor_checked_in") &&
    data.scope === "society" &&
    data.purpose === "service" &&
    !isGuardHomeRoute(homeRoute)
  ) {
    return notificationsRoute();
  }

  switch (type) {
    case "visitor.checkout":
    case "visitor_checked_out":
      return webCheckout && Number.isSafeInteger(entryId) && entryId > 0
        ? residentEntryDetailRoute(entryId)
        : null;
    case "visitor.pending":
    case "visitor.approval_requested":
    case "visitor_pending_approval":
      if (
        isGuardHomeRoute(homeRoute) &&
        Number.isFinite(entryId) &&
        entryId > 0
      ) {
        return guardEntryDetailRoute(entryId);
      }
      if (isGuardHomeRoute(homeRoute)) {
        return guardEntriesRoute("waiting_at_gate");
      }
      return residentVisitorsRoute();
    case "visitor_invite.accepted":
    case "visitor_invite_accepted":
      return inviteId ? residentInviteDetailRoute("visitor", inviteId) : null;
    case "member_invite.accepted":
    case "member_invite_accepted":
      return inviteId ? residentInviteDetailRoute("member", inviteId) : null;
    case "visitor.checkin":
    case "visitor_checked_in":
      if (
        isGuardHomeRoute(homeRoute) &&
        Number.isFinite(entryId) &&
        entryId > 0
      ) {
        return guardEntryDetailRoute(entryId);
      }
      if (Number.isFinite(entryId) && entryId > 0) {
        return residentEntryDetailRoute(entryId);
      }
      return isGuardHomeRoute(homeRoute)
        ? guardEntriesRoute("inside")
        : residentEntriesRoute("inside");
    case "visitor.approved":
    case "visitor_approved":
      if (
        isGuardHomeRoute(homeRoute) &&
        Number.isFinite(entryId) &&
        entryId > 0
      ) {
        return guardEntryDetailRoute(entryId);
      }
      return isGuardHomeRoute(homeRoute)
        ? guardEntriesRoute("expected")
        : residentEntriesRoute("expected");
    case "visitor.rejected":
    case "visitor_rejected":
      if (
        isGuardHomeRoute(homeRoute) &&
        Number.isFinite(entryId) &&
        entryId > 0
      ) {
        return guardEntryDetailRoute(entryId);
      }
      return isGuardHomeRoute(homeRoute)
        ? guardEntriesRoute("all")
        : residentEntriesRoute("all");
    case "maintenance_bill_generated":
    case "maintenance_bill_reminder":
    case "maintenance_payment_verified":
      return notificationsRoute();
    case "hub.announcement":
    case "hub.reply":
    case "hub.reaction":
    case "hub.removed": {
      const postId = Number(stringValue(data.post_id));
      const societyId = Number(stringValue(data.society_id));
      if (
        Number.isFinite(postId) &&
        postId > 0 &&
        Number.isFinite(societyId) &&
        societyId > 0
      ) {
        return residentHubPostRoute(societyId, postId, undefined, "home");
      }
      return notificationsRoute();
    }
    default:
      return notificationsRoute();
  }
}

export function notificationsForHomeRoute(
  notifications: AppNotification[],
  homeRoute?: Href | null,
) {
  return notifications.filter(
    (notification) => !isHiddenForHomeRoute(notification, homeRoute),
  );
}

export function unreadCountForHomeRoute(
  notifications: AppNotification[],
  homeRoute: Href | null | undefined,
  fallbackUnreadCount: number,
) {
  if (!isGuardHomeRoute(homeRoute)) {
    return fallbackUnreadCount;
  }
  return notificationsForHomeRoute(notifications, homeRoute).filter(
    (item) => !item.read_at,
  ).length;
}

export function isHiddenForHomeRoute(
  notification: AppNotification,
  homeRoute?: Href | null,
) {
  if (!isGuardHomeRoute(homeRoute)) {
    return false;
  }
  const type = getNotificationType(notification.data, notification.type);
  return (
    type === "visitor.checkin" ||
    type === "visitor_checked_in" ||
    type === "visitor.checkout" ||
    type === "visitor_checked_out"
  );
}

export function stringValue(value: unknown) {
  return typeof value === "string" && value.trim() ? value : undefined;
}

export function preferText(value: string | null | undefined, fallback: string) {
  return value?.trim() ? value : fallback;
}

export function preferNonLegacyText(
  value: string | null | undefined,
  fallback: string,
  legacyValue: string,
) {
  const trimmed = value?.trim();
  return trimmed && trimmed !== legacyValue ? trimmed : fallback;
}

const legacyVisitorPendingBodies = new Set([
  "A visitor for your flat needs approval.",
  "Your guest is waiting at the gate.",
  "Waiting for the flat to respond.",
]);
const legacyVisitorPendingTitles = new Set([
  "Visitor approval needed",
  "Visitor approval requested",
]);

const legacyVisitorInviteAcceptedBodies = new Set([
  "A guest accepted your invite.",
  "Your guest invite was accepted successfully.",
  "Your invite was accepted successfully.",
  "Visitor pass is ready.",
]);
const legacyVisitorInviteAcceptedTitles = new Set(["Guest invite accepted"]);

const legacyVisitorCheckInBodies = new Set([
  "Your visitor has arrived and completed check-in.",
]);
const legacyVisitorCheckInTitles = new Set(["Visitor checked in"]);

const legacyVisitorCheckOutBodies = new Set(["A visitor checked out."]);
const legacyVisitorCheckOutTitles = new Set(["Visitor checked out"]);

const legacyVisitorApprovedBodies = new Set([
  "The visitor request was approved.",
]);
const legacyVisitorApprovedTitles = new Set([
  "Visitor approved",
  "Visitor approved by the flat",
]);

const legacyVisitorRejectedBodies = new Set([
  "A visitor entry was declined.",
  "The visitor request was rejected.",
]);
const legacyVisitorRejectedTitles = new Set([
  "Visitor request declined",
  "Visitor rejected by the flat",
]);

function notificationBody(
  data: AppNotificationData | undefined,
  title: string | null | undefined,
  value: string | null | undefined,
  missingCompose: (data: AppNotificationData | undefined) => string,
  legacyCompose: (data: AppNotificationData | undefined) => string,
  legacyBodies: Set<string>,
  legacyTitles: Set<string>,
) {
  const trimmed = value?.trim();
  if (!trimmed) {
    return missingCompose(data);
  }
  const suppliedTitle = title?.trim();
  if (suppliedTitle && !legacyTitles.has(suppliedTitle)) {
    return value!;
  }
  if (legacyBodies.has(trimmed) || isLegacyPersonalizedBody(trimmed)) {
    return legacyCompose(data);
  }
  return value!;
}

function notificationTitle(
  value: string | null | undefined,
  fallback: string,
  legacyTitles: Set<string>,
  legacyFallback: string,
) {
  const trimmed = value?.trim();
  if (!trimmed) {
    return fallback;
  }
  return legacyTitles.has(trimmed) ? legacyFallback : value!;
}

function personTitle(
  data: AppNotificationData | undefined,
  fallback: string,
  event: string,
) {
  const name = notificationVisitorName(data);
  return name ? `${name} ${event}` : fallback;
}

function newVisitorPendingBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  if (notificationDataText(data?.category_id) === "notification_info") {
    return name
      ? `${name} is waiting for approval from the resident.`
      : "Visitor is waiting for approval from the resident.";
  }
  return name
    ? `${name} is waiting for your approval. Approve or decline the request.`
    : "Visitor is waiting for your approval. Approve or decline the request.";
}

function newVisitorInviteAcceptedBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  return name
    ? `${name} has completed the visitor details. The visitor pass is ready.`
    : "The guest has completed the visitor details. The visitor pass is ready.";
}

function newVisitorCheckInBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  return name
    ? `${name} just checked in at the society gate.`
    : "The visitor just checked in at the society gate.";
}

function newVisitorCheckOutBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  return name
    ? `${name} has checked out and left the society.`
    : "The visitor has checked out and left the society.";
}

function newVisitorApprovedBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  return name
    ? `${name} has been approved and can now enter the society.`
    : "The visitor has been approved and can now enter the society.";
}

function newVisitorRejectedBody(data: AppNotificationData | undefined) {
  const name = notificationVisitorName(data);
  return name
    ? `${name}'s entry request was declined.`
    : "The visitor's entry request was declined.";
}

function memberJoinedBody(data: AppNotificationData | undefined) {
  const name =
    notificationDataText(
      data?.member_name,
      data?.resident_name,
      data?.full_name,
    ) ?? "A member";
  const flat = notificationDataText(data?.flat_number);
  return flat
    ? `${name} accepted your invitation and is now connected to Flat ${flat}.`
    : `${name} accepted your invitation and is now connected to your flat.`;
}

function memberJoinedTitle(data: AppNotificationData | undefined) {
  const name = notificationDataText(
    data?.member_name,
    data?.resident_name,
    data?.full_name,
  );
  return name ? `${name} joined your flat` : "A member joined your flat";
}

function visitorPendingBody(data: AppNotificationData | undefined) {
  const visitorName = notificationVisitorName(data);
  const flatLabel = notificationFlatLabel(data);
  if (visitorName && flatLabel) {
    return `${visitorName} is waiting for approval for ${flatLabel}.`;
  }
  if (visitorName) {
    return `${visitorName} is waiting for approval.`;
  }
  if (flatLabel) {
    return `A visitor is waiting for approval for ${flatLabel}.`;
  }
  return "A visitor is waiting for approval.";
}

function visitorInviteAcceptedBody(data: AppNotificationData | undefined) {
  const visitorName = notificationVisitorName(data);
  const flatLabel = notificationFlatLabel(data);
  if (visitorName && flatLabel) {
    return `${visitorName} accepted your invite. Visitor pass for ${flatLabel} is ready.`;
  }
  if (visitorName) {
    return `${visitorName} accepted your invite. Visitor pass is ready.`;
  }
  if (flatLabel) {
    return `Guest accepted your invite. Visitor pass for ${flatLabel} is ready.`;
  }
  return "Guest accepted your invite. Visitor pass is ready.";
}

function visitorCheckInBody(data: AppNotificationData | undefined) {
  const visitorName = notificationVisitorName(data);
  const flatLabel = notificationFlatLabel(data);
  if (visitorName && flatLabel) {
    return `${visitorName} checked in at the gate for ${flatLabel}.`;
  }
  if (visitorName) {
    return `${visitorName} checked in at the gate.`;
  }
  if (flatLabel) {
    return `A visitor checked in at the gate for ${flatLabel}.`;
  }
  return "The visitor checked in at the gate.";
}

function visitorCheckOutBody(data: AppNotificationData | undefined) {
  const visitorName = notificationVisitorName(data);
  const flatLabel = notificationFlatLabel(data);
  if (visitorName && flatLabel) {
    return `${visitorName} has left ${flatLabel}.`;
  }
  if (visitorName) {
    return `${visitorName} has left.`;
  }
  if (flatLabel) {
    return `A visitor has left ${flatLabel}.`;
  }
  return "The visitor has left.";
}

function visitorDecisionBody(
  data: AppNotificationData | undefined,
  decision: "approved" | "declined",
) {
  const visitorName = notificationVisitorName(data);
  const flatLabel = notificationFlatLabel(data);
  const approverName = notificationDataText(
    data?.approver_name,
    data?.approved_by_name,
    data?.rejected_by_name,
    data?.resident_name,
    data?.primary_resident_name,
  );

  if (visitorName && flatLabel && approverName) {
    return `${visitorName} was ${decision} for ${flatLabel} by ${approverName}.`;
  }
  if (visitorName && flatLabel) {
    return `${visitorName} was ${decision} for ${flatLabel}.`;
  }
  if (visitorName) {
    return `${visitorName} was ${decision}.`;
  }
  if (flatLabel) {
    return `The visitor was ${decision} for ${flatLabel}.`;
  }
  return `The visitor was ${decision}.`;
}

function isLegacyPersonalizedBody(value: string) {
  return (
    /^.+ checked in for your flat\.$/.test(value) ||
    /^Flat .+ (approved|rejected) the visitor request\.$/.test(value)
  );
}

function notificationVisitorName(data: AppNotificationData | undefined) {
  return notificationDataText(
    data?.visitor_name,
    data?.visitorName,
    data?.full_name,
    data?.fullName,
    data?.delivery_partner,
  );
}

function notificationFlatLabel(data: AppNotificationData | undefined) {
  const flatNumber = notificationDataText(data?.flat_number);
  if (flatNumber) {
    return `Flat ${flatNumber}`;
  }

  return undefined;
}

function notificationDataText(...values: unknown[]) {
  for (const value of values) {
    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
    if (typeof value === "number" && Number.isFinite(value)) {
      return String(value);
    }
  }

  return undefined;
}

function isGuardHomeRoute(homeRoute?: Href | null) {
  return Boolean(homeRoute && String(homeRoute).includes("/guard"));
}
