import {
  buildExplicitVisitorEntryQueryFilters,
  buildVisitorEntryQueryFilters,
  type DateRangePreset,
  type VisitorEntryQueryFilters,
} from "@/features/visitors/visitor-date-ranges";
import type { ModelsVisitorPurpose, ModelsVisitorStatus } from "@/lib/api/generated-api";

export type ResidentVisibleEntriesSegment = "expected" | "inside" | "all";

export type ResidentEntriesPresetState = {
  datePreset: DateRangePreset;
  segment: ResidentVisibleEntriesSegment;
  sheetStatus?: ModelsVisitorStatus;
};

type ResidentEntriesQueryInput = {
  datePreset: DateRangePreset;
  hasExplicitFilters: boolean;
  hasSearchQuery: boolean;
  initial: ResidentEntriesPresetState;
  purpose?: ModelsVisitorPurpose;
  segment: ResidentVisibleEntriesSegment;
  sheetStatus?: ModelsVisitorStatus;
};

type ResidentEntriesRequestParams = VisitorEntryQueryFilters & {
  limit: number;
  offset: number;
  search?: string;
};

function withoutUndefined<T extends Record<string, unknown>>(value: T) {
  return Object.fromEntries(
    Object.entries(value).filter(([, entry]) => entry !== undefined),
  ) as Partial<T>;
}

export function buildResidentEntriesQueryFilters({
  datePreset,
  hasExplicitFilters,
  hasSearchQuery,
  initial,
  purpose,
  segment,
  sheetStatus,
}: ResidentEntriesQueryInput): VisitorEntryQueryFilters {
  if (hasExplicitFilters) {
    return buildExplicitVisitorEntryQueryFilters({
      applyDateRange: true,
      datePreset,
      purpose,
      status: sheetStatus,
    });
  }

  if (hasSearchQuery) {
    return {};
  }

  return buildVisitorEntryQueryFilters({
    datePreset: initial.datePreset,
    isSearchActive: false,
    purpose: undefined,
    segment,
    sheetStatus: initial.sheetStatus,
  });
}

export function buildResidentEntriesRequestParams(
  input: ResidentEntriesQueryInput & {
    limit: number;
    offset: number;
    search?: string;
  },
): Partial<ResidentEntriesRequestParams> {
  return withoutUndefined({
    ...buildResidentEntriesQueryFilters(input),
    limit: input.limit,
    offset: input.offset,
    search: input.search,
  });
}
