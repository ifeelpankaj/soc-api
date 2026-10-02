import assert from "node:assert/strict";
import test from "node:test";

import type { Href } from "expo-router";

import { handleBack } from "./back-navigation";
import {
  isEntrySetupRoute,
  isMeaningfulBackRoute,
  isRedirectAliasRoute,
  isRoleHomeRoute,
  routePath,
} from "./navigation-policy";

function mockRouter(canGoBack: boolean) {
  const calls: string[] = [];
  return {
    router: {
      back: () => calls.push("back"),
      canGoBack: () => canGoBack,
      dismissTo: (href: Href) => calls.push(`dismissTo:${String(href)}`),
      replace: (href: Href) => calls.push(`replace:${String(href)}`),
    },
    calls,
  };
}

test("handleBack dismisses to previous app route when available", () => {
  const { calls, router } = mockRouter(true);

  handleBack(router, "/resident/dashboard" as Href, {
    consumePreviousRoute: () => "/notifications" as Href,
  });

  assert.deepEqual(calls, ["dismissTo:/notifications"]);
});

test("handleBack replaces with previous app route when dismissTo is unavailable", () => {
  const { calls, router } = mockRouter(true);
  const { dismissTo: _dismissTo, ...routerWithoutDismissTo } = router;

  handleBack(routerWithoutDismissTo, "/resident/dashboard" as Href, {
    consumePreviousRoute: () => "/notifications" as Href,
  });

  assert.deepEqual(calls, ["replace:/notifications"]);
});

test("handleBack follows native history when app route history is unavailable", () => {
  const { calls, router } = mockRouter(true);

  handleBack(router, "/resident/dashboard" as Href);

  assert.deepEqual(calls, ["back"]);
});

test("handleBack replaces with fallback when history is unavailable", () => {
  const { calls, router } = mockRouter(false);

  handleBack(router, "/resident/dashboard" as Href);

  assert.deepEqual(calls, ["replace:/resident/dashboard"]);
});

test("route policy classifies setup, root, and alias routes", () => {
  assert.equal(isEntrySetupRoute("/select-society"), true);
  assert.equal(isEntrySetupRoute("/login"), true);
  assert.equal(isRoleHomeRoute("/resident/dashboard"), true);
  assert.equal(isRoleHomeRoute({ pathname: "/guard/home" } as unknown as Href), true);
  assert.equal(isRedirectAliasRoute("/guard/scan?token=abc"), true);
  assert.equal(routePath({ pathname: "/guard/entries", params: { preset: "today" } } as unknown as Href), "/guard/entries");
});

test("handleBack ignores a previous app route matching the current route", () => {
  const { calls, router } = mockRouter(true);

  handleBack(router, "/resident/dashboard" as Href, {
    currentRoute: "/resident/invite" as Href,
    consumePreviousRoute: () => "/resident/invite" as Href,
  });

  assert.deepEqual(calls, ["back"]);
});

test("handleBack falls back when matching previous app route and native history are unavailable", () => {
  const { calls, router } = mockRouter(false);

  handleBack(router, "/resident/dashboard" as Href, {
    currentRoute: "/resident/invite" as Href,
    previousRoute: "/resident/invite" as Href,
  });

  assert.deepEqual(calls, ["replace:/resident/dashboard"]);
});

test("handleBack skips setup routes after workspace selection", () => {
  const { calls, router } = mockRouter(true);

  handleBack(router, "/resident/dashboard" as Href, {
    currentRoute: "/resident/dashboard" as Href,
    previousRoute: "/select-society" as Href,
  });

  assert.deepEqual(calls, ["replace:/resident/dashboard"]);
});

test("handleBack skips setup routes for authenticated sub-screen fallback", () => {
  const { calls, router } = mockRouter(true);

  handleBack(router, "/guard/home" as Href, {
    currentRoute: "/guard/pending" as Href,
    previousRoute: "/login" as Href,
  });

  assert.deepEqual(calls, ["replace:/guard/home"]);
});

test("route policy preserves meaningful app sub-screen order", () => {
  assert.equal(isMeaningfulBackRoute("/guard/home", "/guard/pending"), true);
  assert.equal(isMeaningfulBackRoute("/resident/members", "/resident/members/add"), true);
  assert.equal(isMeaningfulBackRoute("/select-society", "/guard/home"), false);
  assert.equal(isMeaningfulBackRoute("/guard/home", "/guard/home"), false);
});

test("handleBack falls back for resident deep-link sub-screen with no history", () => {
  const { calls, router } = mockRouter(false);

  handleBack(router, "/resident/dashboard" as Href, {
    currentRoute: "/resident/entries/123" as Href,
  });

  assert.deepEqual(calls, ["replace:/resident/dashboard"]);
});

test("handleBack falls back for guard deep-link sub-screen with no history", () => {
  const { calls, router } = mockRouter(false);

  handleBack(router, "/guard/home" as Href, {
    currentRoute: "/guard/entries/123" as Href,
  });

  assert.deepEqual(calls, ["replace:/guard/home"]);
});
