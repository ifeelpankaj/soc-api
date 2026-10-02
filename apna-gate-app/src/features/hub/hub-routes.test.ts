import assert from "node:assert/strict";
import test from "node:test";

import {
  resolveHubPostBackRoute,
  residentHubAnnouncementsRoute,
  residentHubCommunityRoute,
  residentHubRoute,
} from "./hub-routes";
import { residentDashboardRoute } from "@/features/resident/resident-routes";

test("resolveHubPostBackRoute maps returnTo values", () => {
  assert.equal(resolveHubPostBackRoute("home"), residentDashboardRoute());
  assert.equal(resolveHubPostBackRoute("announcements"), residentHubAnnouncementsRoute());
  assert.equal(resolveHubPostBackRoute("community"), residentHubCommunityRoute());
  assert.equal(resolveHubPostBackRoute("hub"), residentHubRoute());
  assert.equal(resolveHubPostBackRoute(undefined), residentHubRoute());
  assert.equal(resolveHubPostBackRoute(["announcements"]), residentHubAnnouncementsRoute());
});
