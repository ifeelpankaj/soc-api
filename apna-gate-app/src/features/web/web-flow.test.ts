import test from "node:test";
import assert from "node:assert/strict";
import { type BrowserDevice, detectWebDevice, isAppleWebBrowser, isIPhoneOrIPad, isMac, requiresWebPushInstallation } from "./device";
import { getGooglePlayURL, resolveWebAccess } from "./web-access";
import { NotificationTracker, mergeNotificationPages } from "./notification-tracker";
import { getMobileMemberships, getMobileResidences } from "../auth/mobile-access.web";
import { resolveBootstrapRoute } from "../auth/bootstrap-routing.web";
import { clearAlertBannerDismissal, dismissAlertBanner, isAlertBannerDismissed, subscribeAlertBanner } from "./notification-banner-session";
import { notificationWorkspace,validNotificationLink } from "./notification-target";
import type { AppNotification } from "@/lib/api/notification-api-extensions";
import { pendingNotificationRoute, clearPendingNotification } from "./notification-handoff.web";

test("web login handoff survives reads until consumed", () => {
  const original = Object.getOwnPropertyDescriptor(globalThis, "sessionStorage");
  const values = new Map<string, string>();
  Object.defineProperty(globalThis, "sessionStorage", { configurable: true, value: {
    getItem: (key: string) => values.get(key) ?? null,
    removeItem: (key: string) => values.delete(key),
  } });
  try {
    const target = { id: "53bdc1b1-a7bb-41bb-9522-3a17e98c45c8", recipient: "1" };
    values.set("apna_notification_open", JSON.stringify(target));
    const route = { pathname: "/notification-open", params: target };
    assert.deepEqual(pendingNotificationRoute(), route);
    assert.deepEqual(pendingNotificationRoute(), route);
    clearPendingNotification();
    assert.equal(pendingNotificationRoute(), null);
    values.set("apna_notification_open", JSON.stringify({ id: [target.id], recipient: "1" }));
    assert.equal(pendingNotificationRoute(), null);
    values.set("apna_notification_open", "broken json");
    assert.equal(pendingNotificationRoute(), null);
  } finally {
    if (original) Object.defineProperty(globalThis, "sessionStorage", original);
    else Reflect.deleteProperty(globalThis, "sessionStorage");
  }
});

test("device detection is SSR-safe, including desktop-identifying iPad",()=>{
  assert.equal(detectWebDevice(),"desktop");
  assert.equal(detectWebDevice({userAgent:"iPhone",platform:"iPhone",maxTouchPoints:5}),"apple");
  assert.equal(detectWebDevice({userAgent:"Macintosh",platform:"MacIntel",maxTouchPoints:5}),"apple");
  assert.equal(detectWebDevice({userAgent:"Macintosh",platform:"MacIntel",maxTouchPoints:0}),"desktop");
  assert.equal(detectWebDevice({userAgent:"Android",platform:"Linux",maxTouchPoints:5}),"android");
});

const browserCases: { name: string; device: BrowserDevice; apple: boolean; mobileApple: boolean; blocked: "android" | "desktop" }[] = [
  { name: "iPhone Safari", device: { userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 Version/18.0 Mobile/15E148 Safari/604.1", platform: "iPhone", maxTouchPoints: 5 }, apple: true, mobileApple: true, blocked: "desktop" },
  { name: "iPad Safari", device: { userAgent: "Mozilla/5.0 (iPad; CPU OS 18_0 like Mac OS X) AppleWebKit/605.1.15 Version/18.0 Mobile/15E148 Safari/604.1", platform: "iPad", maxTouchPoints: 5 }, apple: true, mobileApple: true, blocked: "desktop" },
  { name: "iPad desktop-mode Safari", device: { userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 Version/18.0 Safari/605.1.15", platform: "MacIntel", maxTouchPoints: 5 }, apple: true, mobileApple: true, blocked: "desktop" },
  { name: "Mac Safari", device: { userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/18.0 Safari/605.1.15", platform: "MacIntel", maxTouchPoints: 0 }, apple: true, mobileApple: false, blocked: "desktop" },
  { name: "Mac Chrome", device: { userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/140.0.0.0 Safari/537.36", platform: "MacIntel", maxTouchPoints: 0 }, apple: true, mobileApple: false, blocked: "desktop" },
  { name: "Android Chrome", device: { userAgent: "Mozilla/5.0 (Linux; Android 15; Pixel 9) AppleWebKit/537.36 Chrome/140.0.0.0 Mobile Safari/537.36", platform: "Linux armv8l", maxTouchPoints: 5 }, apple: false, mobileApple: false, blocked: "android" },
  { name: "Windows Chrome", device: { userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/140.0.0.0 Safari/537.36", platform: "Win32", maxTouchPoints: 0 }, apple: false, mobileApple: false, blocked: "desktop" },
  { name: "Linux Firefox", device: { userAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:140.0) Gecko/20100101 Firefox/140.0", platform: "Linux x86_64", maxTouchPoints: 0 }, apple: false, mobileApple: false, blocked: "desktop" },
];

for (const browser of browserCases) {
  test(`web access and mobile push guidance: ${browser.name}`, () => {
    assert.equal(isAppleWebBrowser(browser.device), browser.apple);
    assert.equal(isIPhoneOrIPad(browser.device), browser.mobileApple);
    assert.equal(isMac(browser.device), browser.apple && !browser.mobileApple);
    assert.equal(resolveWebAccess("true", browser.device), browser.apple ? "allowed" : browser.blocked);
    for (const flag of [undefined, "false", "", "TRUE"]) {
      assert.equal(resolveWebAccess(flag, browser.device), "allowed");
    }
    assert.equal(requiresWebPushInstallation(detectWebDevice(browser.device), false), browser.mobileApple);
    assert.equal(requiresWebPushInstallation(detectWebDevice(browser.device), true), false);
  });
}

test("restricted web waits for device detection; unrestricted web needs no browser globals", () => {
  assert.equal(resolveWebAccess("true"), "loading");
  assert.equal(resolveWebAccess("false"), "allowed");
  assert.equal(resolveWebAccess(undefined), "allowed");
  assert.equal(isAppleWebBrowser(), false);
  assert.equal(isMac(), false);
  assert.equal(isIPhoneOrIPad(), false);
});

test("Google Play destination comes only from the environment, with no hardcoded fallback", () => {
  const previous = process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL;
  try {
    const customURL = "https://play.google.com/store/apps/details?id=com.example.resident";
    process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL = customURL;
    assert.equal(getGooglePlayURL(), customURL);
    process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL = `  ${customURL}  `;
    assert.equal(getGooglePlayURL(), customURL);
    for (const missing of ["", "   ", undefined]) {
      if (missing === undefined) delete process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL;
      else process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL = missing;
      assert.equal(getGooglePlayURL(), undefined);
    }
  } finally {
    if (previous === undefined) delete process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL;
    else process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL = previous;
  }
});
test("mixed membership ignores guard defaults and offers resident workspaces",()=>{
  const bootstrap={defaultDashboard:{path:"/guard/home"},memberships:[{role:"staff" as const,status:"active" as const,society_id:1},{role:"resident" as const,status:"active" as const,society_id:2}]};
  assert.deepEqual(getMobileMemberships(bootstrap).map(m=>m.role),["resident"]);
  assert.equal(resolveBootstrapRoute(bootstrap),"/resident/dashboard");
  assert.equal(resolveBootstrapRoute({memberships:[bootstrap.memberships[0]]}),"/select-society");
});

test("browser push installation is required only on Apple mobile browsers", () => {
  assert.equal(requiresWebPushInstallation("desktop", false), false);
  assert.equal(requiresWebPushInstallation("android", false), false);
  assert.equal(requiresWebPushInstallation("apple", false), true);
  assert.equal(requiresWebPushInstallation("apple", true), false);
});

test("web includes only active resident access and deduplicates society memberships", () => {
  const bootstrap = {
    memberships: [
      { role: "staff" as const, status: "active" as const, society_id: 1 },
      { role: "resident" as const, status: "suspended" as const, society_id: 2 },
      { role: "resident" as const, status: "pending" as const, society_id: 3 },
      { role: "resident" as const, status: "active" as const, society_id: 4 },
      { role: "resident" as const, status: "active" as const, society_id: 5 },
    ],
    residences: [
      { status: "active" as const, society_id: 4, flat_id: 40 },
      { status: "inactive" as const, society_id: 2, flat_id: 20 },
    ],
  };
  assert.deepEqual(getMobileMemberships(bootstrap).map(item => item.society_id), [5]);
  assert.deepEqual(getMobileResidences(bootstrap).map(item => item.flat_id), [40]);
  assert.equal(resolveBootstrapRoute(bootstrap), "/select-society");
});

test("native routing continues to honor guard defaults while web selects resident access", async () => {
  const { resolveBootstrapRoute: resolveNativeRoute } = await import("../auth/bootstrap-routing");
  const bootstrap = {
    defaultDashboard: { path: "/guard/home" },
    memberships: [{ role: "staff" as const, status: "active" as const, society_id: 1 }],
    residences: [{ status: "active" as const, society_id: 2, flat_id: 20 }],
  };
  assert.equal(resolveNativeRoute(bootstrap), "/guard/home");
  assert.equal(resolveBootstrapRoute(bootstrap), "/resident/dashboard");
  assert.equal(resolveNativeRoute({ residences: bootstrap.residences }), "/resident/dashboard");
});

test("alert dismissal survives remounts, stays account-scoped and clears on logout", () => {
  clearAlertBannerDismissal();
  let updates = 0;
  const unsubscribe = subscribeAlertBanner(() => updates++);
  assert.equal(isAlertBannerDismissed(1), false);
  dismissAlertBanner(1);
  assert.equal(isAlertBannerDismissed(1), true);
  assert.equal(isAlertBannerDismissed(2), false);
  unsubscribe();
  assert.equal(isAlertBannerDismissed(1), true);
  clearAlertBannerDismissal();
  assert.equal(isAlertBannerDismissed(1), false);
  assert.equal(updates, 1);
});
test("initial polling stays quiet and push/poll delivery deduplicates",()=>{
  const tracker=new NotificationTracker();
  assert.deepEqual(tracker.poll(["old"]),[]);
  assert.equal(tracker.push("new"),true);
  assert.deepEqual(tracker.poll(["new","old"]),[]);
  assert.deepEqual(tracker.poll(["newer","new"]),["newer"]);
  assert.equal(tracker.push("newer"),false);
  assert.deepEqual(new NotificationTracker().poll(["newer"]),[]);
});
test("poll updates preserve older loaded pages and read state",()=>{
  const old=[{id:"a",created_at:"2026-01-02",read_at:undefined as string|undefined},{id:"b",created_at:"2026-01-01",read_at:undefined as string|undefined}];
  const merged=mergeNotificationPages(old,[{id:"a",created_at:"2026-01-02",read_at:"now"},{id:"c",created_at:"2026-01-03",read_at:undefined}]);
  assert.deepEqual(merged.map(i=>i.id),["c","a","b"]);assert.equal(merged[1].read_at,"now");
});
test("tap requires owned notification and current matching residence",()=>{
  const item={user_id:1,society_id:2,flat_id:3,data:{}} as AppNotification;
  const bootstrap={residences:[{user_id:1,society_id:2,flat_id:3,status:"active" as const}]};
  assert.ok(notificationWorkspace(item,bootstrap,1));
  assert.equal(notificationWorkspace(item,bootstrap,9),null);
  assert.equal(notificationWorkspace({...item,flat_id:4},bootstrap,1),null);
  assert.equal(notificationWorkspace(item,{residences:[]},1),null);
  assert.equal(validNotificationLink("../../guard/home","1"),false);
  assert.equal(validNotificationLink("53bdc1b1-a7bb-41bb-9522-3a17e98c45c8","1"),true);
});
