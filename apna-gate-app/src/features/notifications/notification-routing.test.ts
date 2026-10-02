import assert from "node:assert/strict";
import test from "node:test";

import {
  notificationsForHomeRoute,
  notificationPresentation,
  notificationRoute,
  preferText,
  unreadCountForHomeRoute,
} from "./notification-routing";
import type { AppNotification } from "@/lib/api/notification-api-extensions";

function notification(
  type: string,
  data: AppNotification["data"] = {},
  title = "",
  body = "",
): AppNotification {
  return {
    id: "notification-1",
    user_id: 101,
    type,
    title,
    body,
    data,
    created_at: "2026-08-29T18:10:00Z",
  };
}

test("future payloads and Maintenance events fall back to the inbox", () => {
  assert.equal(notificationRoute(notification("visitor.checkin", {entry_id:"42",schema_version:2})), "/notifications");
  assert.deepEqual(
    notificationRoute(
      notification("hub.announcement", {
        post_id: "9",
        society_id: "3",
        schema_version: 1,
      }),
    ),
    {
      pathname: "/resident/hub/posts/[postId]",
      params: { postId: "9", societyId: "3", channelId: undefined, returnTo: "home" },
    },
  );
  assert.equal(notificationRoute(notification("maintenance_bill_generated", {bill_id:"7",schema_version:1})), "/notifications");
});

test("checkout opens entry details only when web routing is requested", () => {
  for (const type of ["visitor.checkout", "visitor_checked_out"]) {
    const item = notification(type, { entry_id: "42" });
    assert.equal(notificationRoute(item, "/resident/dashboard"), null);
    assert.deepEqual(notificationRoute(item, "/resident/dashboard", true),
      notificationRoute(notification("visitor.checkin", { entry_id: "42" }), "/resident/dashboard"));
    for (const entry_id of ["", "bad", "-1", "1.5"]) {
      assert.equal(notificationRoute(notification(type, { entry_id }), "/resident/dashboard", true), null);
    }
  }
});

test("society Service check-in uses informational copy and resident inbox routing", () => {
  const item = notification("visitor.checkin", {
    entry_id: "42", scope: "society", purpose: "service", category_id: "notification_info",
  }, "Newspaper Delivery has arrived", "Rajesh from Newspaper Delivery has entered the society.");
  assert.equal(notificationRoute(item, "/resident/dashboard"), "/notifications");
  assert.equal(notificationRoute(item), "/notifications");
  assert.notEqual(notificationRoute(item, "/guard/dashboard"), "/notifications");
  assert.equal(notificationPresentation(item).title, item.title);
  assert.equal(notificationPresentation(item).body, item.body);
  assert.notEqual(notificationRoute(notification("visitor.checkin", { entry_id: "42" }), "/resident/dashboard"), "/notifications");
});

test("exact unnamed event fallbacks and delivery names", () => {
  const cases = [
    [
      "visitor.pending",
      "Visitor is at the gate",
      "Visitor is waiting for your approval. Approve or decline the request.",
    ],
    [
      "visitor.checkin",
      "Visitor has arrived",
      "The visitor just checked in at the society gate.",
    ],
    [
      "visitor.checkout",
      "Visitor has left",
      "The visitor has checked out and left the society.",
    ],
    [
      "visitor_invite.accepted",
      "Guest accepted your invite",
      "The guest has completed the visitor details. The visitor pass is ready.",
    ],
    [
      "visitor.approved",
      "Visitor approved",
      "The visitor has been approved and can now enter the society.",
    ],
    [
      "visitor.rejected",
      "Visitor request declined",
      "The visitor's entry request was declined.",
    ],
    [
      "member_invite.accepted",
      "A member joined your flat",
      "A member accepted your invitation and is now connected to your flat.",
    ],
  ];
  for (const [type, title, body] of cases) {
    const copy = notificationPresentation(
      notification(type, { visitor_name: "  " }),
    );
    assert.equal(copy.title, title);
    assert.equal(copy.body, body);
  }
  const delivery = notificationPresentation(
    notification("visitor.pending", {
      visitor_name: " ",
      delivery_partner: "  Swiggy  ",
    }),
  );
  assert.equal(delivery.title, "Swiggy is at the gate");
  assert.equal(
    delivery.body,
    "Swiggy is waiting for your approval. Approve or decline the request.",
  );
  const named = notificationPresentation(
    notification("visitor.pending", {
      visitor_name: "  Rahul  ",
      delivery_partner: "Swiggy",
    }),
  );
  assert.equal(named.title, "Rahul is at the gate");
});

test("maps notification event copy", () => {
  assert.equal(
    notificationPresentation(notification("visitor.pending")).title,
    "Visitor is at the gate",
  );
  assert.equal(
    notificationPresentation(notification("visitor.checkin")).title,
    "Visitor has arrived",
  );
  assert.equal(
    notificationPresentation(notification("visitor.checkout")).title,
    "Visitor has left",
  );
  assert.equal(
    notificationPresentation(notification("visitor.approved")).title,
    "Visitor approved",
  );
  assert.equal(
    notificationPresentation(notification("visitor.rejected")).title,
    "Visitor request declined",
  );
  assert.equal(
    notificationPresentation(notification("visitor_invite.accepted")).title,
    "Guest accepted your invite",
  );
  assert.equal(
    notificationPresentation(notification("member_invite.accepted")).title,
    "A member joined your flat",
  );
});

test("preserves backend copy for known notification events", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor.checkin",
      {},
      "Visitor checked in",
      "Rahul checked in at the gate for Flat 005.",
    ),
  );

  assert.equal(presentation.title, "Visitor checked in");
  assert.equal(presentation.body, "Rahul checked in at the gate for Flat 005.");
});

test("normalizes legacy visitor approval copy with visitor flat and approver", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor.approved",
      {
        visitor_name: "Pankaj",
        approver_name: "Rahul Sharma",
        block: "A",
        flat_number: " 1203 ",
      },
      "Visitor approved by the flat",
      "The visitor request was approved.",
    ),
  );

  assert.equal(presentation.title, "Visitor approved");
  assert.equal(
    presentation.body,
    "Pankaj was approved for Flat 1203 by Rahul Sharma.",
  );
});

test("uses natural copy when backend copy is missing", () => {
  assert.equal(
    notificationPresentation(
      notification(
        "visitor.approved",
        { visitor_name: "Pankaj", flat_number: "005" },
        "",
        "",
      ),
    ).body,
    "Pankaj has been approved and can now enter the society.",
  );
  assert.equal(
    notificationPresentation(
      notification("visitor.approved", { visitor_name: "Pankaj" }, "", ""),
    ).body,
    "Pankaj has been approved and can now enter the society.",
  );
  assert.equal(
    notificationPresentation(
      notification("visitor.approved", { flat_number: "005" }, "", ""),
    ).body,
    "The visitor has been approved and can now enter the society.",
  );
  assert.equal(
    notificationPresentation(notification("visitor.approved", {}, "", "")).body,
    "The visitor has been approved and can now enter the society.",
  );
});

test("preserves personalized backend copy and uses natural missing-copy fallbacks", () => {
  assert.equal(
    notificationPresentation(
      notification(
        "visitor.checkin",
        { visitor_name: "  Pankaj  ", flat_number: " 005 " },
        "Pankaj has arrived",
        "Pankaj checked in for your flat.",
      ),
    ).body,
    "Pankaj checked in for your flat.",
  );
  assert.equal(
    notificationPresentation(
      notification(
        "visitor.checkout",
        { visitor_name: "Pankaj", flat_number: "005" },
        "",
        "",
      ),
    ).body,
    "Pankaj has checked out and left the society.",
  );
});

test("normalizes legacy rejected copy with flat only", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor.rejected",
      { flat_number: "B-404" },
      "Visitor rejected by the flat",
      "The visitor request was rejected.",
    ),
  );

  assert.equal(presentation.title, "Visitor request declined");
  assert.equal(presentation.body, "The visitor was declined for Flat B-404.");
});

test("normalizes legacy visitor invite accepted copy", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor_invite.accepted",
      {},
      "Guest invite accepted",
      "A guest accepted your invite.",
    ),
  );

  assert.equal(presentation.title, "Guest invite accepted");
  assert.equal(
    presentation.body,
    "Guest accepted your invite. Visitor pass is ready.",
  );
});

test("personalized titles protect even old-looking backend bodies", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor_invite.accepted",
      { visitor_name: "Rahul", flat_number: "005" },
      "Rahul accepted your guest invite",
      "Visitor pass is ready.",
    ),
  );

  assert.equal(presentation.title, "Rahul accepted your guest invite");
  assert.equal(presentation.body, "Visitor pass is ready.");
});

test("preserves valid backend body copy", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor_invite.accepted",
      { visitor_name: "Rahul", flat_number: "005" },
      "Guest invite accepted",
      "Rahul accepted your invite. Visitor pass for Flat 005 is ready.",
    ),
  );

  assert.equal(presentation.title, "Guest invite accepted");
  assert.equal(
    presentation.body,
    "Rahul accepted your invite. Visitor pass for Flat 005 is ready.",
  );
});

test("uses blank-aware notification fallbacks", () => {
  assert.equal(preferText("   ", "Fallback"), "Fallback");
  assert.equal(
    preferText("Rahul has arrived", "Fallback"),
    "Rahul has arrived",
  );

  const presentation = notificationPresentation(
    notification("visitor_invite.accepted", {}, "   ", "   "),
  );
  assert.equal(presentation.title, "Guest accepted your invite");
  assert.equal(
    presentation.body,
    "The guest has completed the visitor details. The visitor pass is ready.",
  );
});

test("missing guard waiting copy does not claim the guard can decide", () => {
  const presentation = notificationPresentation(
    notification(
      "visitor.pending",
      { category_id: "notification_info", visitor_name: "Rahul" },
      "",
      "",
    ),
  );
  assert.equal(presentation.title, "Rahul is at the gate");
  assert.equal(
    presentation.body,
    "Rahul is waiting for approval from the resident.",
  );
});

test("routes invite accepted notifications to invite details", () => {
  assert.equal(
    notificationRoute(
      notification("visitor_invite.accepted", { invite_id: "81" }),
    ),
    "/resident/invites/visitor/81",
  );
  assert.equal(
    notificationRoute(
      notification("member_invite.accepted", { invite_id: "82" }),
    ),
    "/resident/invites/member/82",
  );
});

test("routes visitor approval and check-in notifications", () => {
  assert.equal(
    notificationRoute(notification("visitor.pending")),
    "/resident/visitors",
  );
  assert.deepEqual(
    notificationRoute(
      notification("visitor.pending", { entry_id: "143" }),
      "/guard/dashboard",
    ),
    {
      pathname: "/guard/entries/[entryId]",
      params: { entryId: "143" },
    },
  );
  assert.deepEqual(
    notificationRoute(notification("visitor.pending"), "/guard/dashboard"),
    {
      pathname: "/guard/entries",
      params: { preset: "waiting_at_gate" },
    },
  );
  assert.deepEqual(
    notificationRoute(
      notification("visitor.checkin", { entry_id: "143" }),
      "/guard/dashboard",
    ),
    {
      pathname: "/guard/entries/[entryId]",
      params: { entryId: "143" },
    },
  );
  assert.deepEqual(
    notificationRoute(notification("visitor.checkin"), "/guard/dashboard"),
    {
      pathname: "/guard/entries",
      params: { preset: "inside" },
    },
  );
  assert.equal(
    notificationRoute(
      notification("visitor.checkin", { entry_id: "143" }),
      "/resident/dashboard",
    ),
    "/resident/entries/143",
  );
  assert.deepEqual(
    notificationRoute(notification("visitor.checkin"), "/resident/dashboard"),
    {
      pathname: "/resident/entries",
      params: { preset: "inside" },
    },
  );
  assert.deepEqual(
    notificationRoute(
      notification("visitor.approved", { entry_id: "143" }),
      "/guard/dashboard",
    ),
    {
      pathname: "/guard/entries/[entryId]",
      params: { entryId: "143" },
    },
  );
  assert.deepEqual(
    notificationRoute(
      notification("visitor.rejected", { entry_id: "143" }),
      "/guard/dashboard",
    ),
    {
      pathname: "/guard/entries/[entryId]",
      params: { entryId: "143" },
    },
  );
});

test("hides guard check-in notifications from guard notification list", () => {
  const checkIn = notification("visitor.checkin", { entry_id: "143" });
  const checkOut = notification("visitor.checkout", { entry_id: "146" });
  const approved = notification("visitor.approved", { entry_id: "144" });

  assert.deepEqual(
    notificationsForHomeRoute(
      [checkIn, checkOut, approved],
      "/guard/dashboard",
    ),
    [approved],
  );
  assert.deepEqual(
    notificationsForHomeRoute(
      [checkIn, checkOut, approved],
      "/resident/dashboard",
    ),
    [checkIn, checkOut, approved],
  );
});

test("counts only visible unread notifications for guard notification list", () => {
  const checkIn = notification("visitor.checkin", { entry_id: "143" });
  const approved = notification("visitor.approved", { entry_id: "144" });
  const readRejected = {
    ...notification("visitor.rejected", { entry_id: "145" }),
    read_at: "2026-08-29T18:15:00Z",
  };

  assert.equal(
    unreadCountForHomeRoute(
      [checkIn, approved, readRejected],
      "/guard/dashboard",
      3,
    ),
    1,
  );
  assert.equal(
    unreadCountForHomeRoute(
      [checkIn, approved, readRejected],
      "/resident/dashboard",
      3,
    ),
    3,
  );
});
