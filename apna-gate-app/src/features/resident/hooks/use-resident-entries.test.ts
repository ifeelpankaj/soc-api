import assert from "node:assert/strict";
import test from "node:test";

import {
  buildResidentEntriesRequestParams,
  type ResidentVisibleEntriesSegment,
} from "./resident-entries-query";
import { parseResidentEntriesPreset } from "@/features/resident/resident-routes";
import type { DateRangePreset } from "@/features/visitors/visitor-date-ranges";
import type { ModelsVisitorPurpose, ModelsVisitorStatus } from "@/lib/api/generated-api";

type TestInput = {
  datePreset?: DateRangePreset;
  hasExplicitFilters?: boolean;
  hasSearchQuery?: boolean;
  initialDatePreset?: DateRangePreset;
  initialSheetStatus?: ModelsVisitorStatus;
  initialSegment?: ResidentVisibleEntriesSegment;
  limit?: number;
  offset?: number;
  purpose?: ModelsVisitorPurpose;
  search?: string;
  segment?: ResidentVisibleEntriesSegment;
  sheetStatus?: ModelsVisitorStatus;
};

function requestParams(input: TestInput = {}) {
  return buildResidentEntriesRequestParams({
    datePreset: input.datePreset ?? "today",
    hasExplicitFilters: input.hasExplicitFilters ?? false,
    hasSearchQuery: input.hasSearchQuery ?? false,
    initial: {
      datePreset: input.initialDatePreset ?? "today",
      segment: input.initialSegment ?? "expected",
      sheetStatus: input.initialSheetStatus,
    },
    limit: input.limit ?? 20,
    offset: input.offset ?? 0,
    purpose: input.purpose,
    search: input.search,
    segment: input.segment ?? "expected",
    sheetStatus: input.sheetStatus,
  });
}

function assertNoDateOrStatusParams(params: Record<string, unknown>) {
  assert.equal("event" in params, false);
  assert.equal("eventFrom" in params, false);
  assert.equal("eventTo" in params, false);
  assert.equal("status" in params, false);
  assert.equal("createdFrom" in params, false);
  assert.equal("createdTo" in params, false);
}

test("resident default expected preset filters approved entries by expected_at", () => {
  const params = requestParams({ segment: "expected", initialSegment: "expected" });

  assert.equal(params.status, "approved");
  assert.equal(params.event, "expected");
  assert.equal(typeof params.eventFrom, "string");
  assert.equal(typeof params.eventTo, "string");
  assert.equal("createdFrom" in params, false);
  assert.equal("createdTo" in params, false);
  assert.equal(params.limit, 20);
  assert.equal(params.offset, 0);
});

test("resident inside preset only filters currently checked-in visitors", () => {
  const params = requestParams({ segment: "inside", initialSegment: "inside" });

  assert.equal(params.status, "checked_in");
  assert.equal("event" in params, false);
  assert.equal("eventFrom" in params, false);
  assert.equal("eventTo" in params, false);
});

test("resident all preset sends no preset filter fields", () => {
  const params = requestParams({
    initialDatePreset: "all",
    initialSegment: "all",
    segment: "all",
  });

  assertNoDateOrStatusParams(params);
});

test("resident search-only ignores the selected preset", () => {
  const params = requestParams({
    hasSearchQuery: true,
    search: "Rahul",
    segment: "expected",
  });

  assert.equal(params.search, "Rahul");
  assertNoDateOrStatusParams(params);
});

test("resident explicit checked-out date filter uses checkout event time", () => {
  const params = requestParams({
    datePreset: "today",
    hasExplicitFilters: true,
    sheetStatus: "checked_out",
  });

  assert.equal(params.status, "checked_out");
  assert.equal(params.event, "checked_out");
  assert.equal(typeof params.eventFrom, "string");
  assert.equal(typeof params.eventTo, "string");
  assert.equal("createdFrom" in params, false);
  assert.equal("createdTo" in params, false);
});

test("resident search with explicit filters keeps search and explicit params", () => {
  const params = requestParams({
    hasExplicitFilters: true,
    hasSearchQuery: true,
    purpose: "delivery",
    search: "Rahul",
    sheetStatus: "waiting_approval",
  });

  assert.equal(params.search, "Rahul");
  assert.equal(params.purpose, "delivery");
  assert.equal(params.status, "waiting_approval");
  assert.equal(params.event, "created");
  assert.equal(typeof params.eventFrom, "string");
  assert.equal(typeof params.eventTo, "string");
});

test("resident parser keeps only expected, inside, and all as runtime presets", () => {
  assert.equal(parseResidentEntriesPreset("expected"), "expected");
  assert.equal(parseResidentEntriesPreset("inside"), "inside");
  assert.equal(parseResidentEntriesPreset("all"), "all");
});

test("resident legacy today preset falls back to expected", () => {
  assert.equal(parseResidentEntriesPreset("today"), "expected");
});

test("resident legacy recent preset falls back to expected only for stale-link compatibility", () => {
  assert.equal(parseResidentEntriesPreset("recent"), "expected");
});
