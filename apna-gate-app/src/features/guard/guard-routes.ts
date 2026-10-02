import type { Href } from "expo-router";

import type { ModelsVisitorStatus } from "@/lib/api/generated-api";
import {
  type DateRangePreset,
  type LogsSegment,
  DATE_RANGE_OPTIONS,
  buildExplicitVisitorEntryQueryFilters,
  buildVisitorEntryQueryFilters,
  getDateRange,
  getDateRangeLabel,
  getVisitorFilterSummary,
  getHalfOpenRangeIST,
  getLogsSegmentSummaryTitle,
  getTodayRange,
  type HalfOpenRange,
  type VisitorEntryEvent,
  type VisitorEntryQueryFilters,
} from "@/features/visitors/visitor-date-ranges";

export type GuardEntriesPreset = "all" | "today" | "expected" | "waiting_at_gate" | "inside" | "checked_out";
export type LogsPreset = GuardEntriesPreset;

export type {
  DateRangePreset,
  HalfOpenRange,
  LogsSegment,
  VisitorEntryEvent,
  VisitorEntryQueryFilters,
};

export {
  DATE_RANGE_OPTIONS,
  buildExplicitVisitorEntryQueryFilters,
  buildVisitorEntryQueryFilters,
  getDateRange,
  getDateRangeLabel,
  getVisitorFilterSummary,
  getHalfOpenRangeIST,
  getLogsSegmentSummaryTitle,
  getTodayRange,
};

export type GuardCheckInInput =
  | { source: "qr"; token: string }
  | { source: "entry"; entryId: number; token?: string };

export function firstParam(value?: string | string[]): string | undefined {
  return Array.isArray(value) ? value[0] : value;
}

export function parseCheckInParams(params: {
  source?: string | string[];
  token?: string | string[];
  entryId?: string | string[];
}): GuardCheckInInput | null {
  const source = firstParam(params.source);
  const token = firstParam(params.token)?.trim();
  const entryId = Number(firstParam(params.entryId));

  if (source === "qr" && token) {
    return { source: "qr", token };
  }

  if (source === "entry" && Number.isFinite(entryId) && entryId > 0) {
    return token
      ? { source: "entry", entryId, token }
      : { source: "entry", entryId };
  }

  return null;
}

export const guardHomeRoute = (): Href =>
  ({
    pathname: "/guard/home",
  }) as unknown as Href;

export const guardProfileRoute = (): Href =>
  ({
    pathname: "/guard/profile",
  }) as unknown as Href;

export const guardScannerRoute = (): Href =>
  ({
    pathname: "/guard/scanner",
  }) as unknown as Href;

export const guardAddEntryRoute = (): Href =>
  ({
    pathname: "/guard/add-entry",
  }) as unknown as Href;

export const guardPendingRoute = (): Href =>
  ({
    pathname: "/guard/pending",
  }) as unknown as Href;

export const guardWaitingAtGateRoute = (): Href =>
  ({
    pathname: "/guard/waiting-at-gate",
  }) as unknown as Href;

export const guardCheckInRoute = (params: GuardCheckInInput): Href =>
  ({
    pathname: "/guard/check-in",
    params:
      params.source === "qr"
        ? { source: "qr", token: params.token }
        : {
            source: "entry",
            entryId: String(params.entryId),
            ...(params.token ? { token: params.token } : {}),
          },
  }) as unknown as Href;

export const guardEntriesRoute = (preset: GuardEntriesPreset = "today"): Href =>
  ({
    pathname: "/guard/entries",
    params: { preset },
  }) as unknown as Href;

export const guardEntryDetailRoute = (entryId: number): Href =>
  ({
    pathname: "/guard/entries/[entryId]",
    params: { entryId: String(entryId) },
  }) as unknown as Href;

/** @deprecated Use guardEntriesRoute instead */
export function guardLogsRoute(preset: GuardEntriesPreset = "today"): Href {
  return guardEntriesRoute(preset);
}

const VALID_PRESETS = new Set<GuardEntriesPreset>([
  "today",
  "all",
  "expected",
  "waiting_at_gate",
  "inside",
  "checked_out",
]);

export function parseGuardEntriesPreset(value?: string | string[]): GuardEntriesPreset {
  const preset = firstParam(value);

  if (preset && VALID_PRESETS.has(preset as GuardEntriesPreset)) {
    return preset as GuardEntriesPreset;
  }

  return "today";
}

/** @deprecated Use parseGuardEntriesPreset instead */
export function parseLogsPreset(value?: string | string[]): GuardEntriesPreset {
  return parseGuardEntriesPreset(value);
}

export type PresetState = {
  datePreset: DateRangePreset;
  segment: LogsSegment;
  sheetStatus?: ModelsVisitorStatus;
};

export const PRESET_CONFIG = {
  today: {
    segment: "today" as LogsSegment,
    datePreset: "today" as DateRangePreset,
    sheetStatus: undefined,
  },
  expected: {
    segment: "expected" as LogsSegment,
    datePreset: "today" as DateRangePreset,
    sheetStatus: undefined,
  },
  waiting_at_gate: {
    segment: "expected" as LogsSegment,
    datePreset: "today" as DateRangePreset,
    sheetStatus: "approved" as ModelsVisitorStatus,
  },
  inside: {
    segment: "inside" as LogsSegment,
    datePreset: "today" as DateRangePreset,
    sheetStatus: undefined,
  },
  checked_out: {
    segment: "today" as LogsSegment,
    datePreset: "today" as DateRangePreset,
    sheetStatus: "checked_out" as ModelsVisitorStatus,
  },
  all: {
    segment: "all" as LogsSegment,
    datePreset: "all" as DateRangePreset,
    sheetStatus: undefined,
  },
} satisfies Record<GuardEntriesPreset, PresetState>;

export function presetToInitialState(preset: GuardEntriesPreset): PresetState {
  return PRESET_CONFIG[preset];
}
