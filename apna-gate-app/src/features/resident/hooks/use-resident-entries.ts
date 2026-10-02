import { useCallback, useEffect, useMemo, useState } from "react";

import {
  getVisitorFilterSummary,
  type DateRangePreset,
} from "@/features/visitors/visitor-date-ranges";
import {
  buildResidentEntriesRequestParams,
  type ResidentEntriesPresetState,
  type ResidentVisibleEntriesSegment,
} from "@/features/resident/hooks/resident-entries-query";
import { useResident } from "@/features/resident/resident-context";
import {
  parseResidentEntriesPreset,
  type ResidentEntriesPreset,
} from "@/features/resident/resident-routes";
import { usePaginatedQuery } from "@/features/shared/use-paginated-query";
import {
  type ModelsVisitorEntry,
  type ModelsVisitorPurpose,
  type ModelsVisitorStatus,
} from "@/lib/api/generated-api";
import { useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery } from "@/lib/api/resident-api-extensions";

const PAGE_SIZE = 20;
const SEARCH_DEBOUNCE_MS = 400;
const visitorEntryKey = (entry: ModelsVisitorEntry) => entry.id;

export type { ResidentVisibleEntriesSegment } from "@/features/resident/hooks/resident-entries-query";

function segmentStatus(segment: ResidentVisibleEntriesSegment): ModelsVisitorStatus | undefined {
  if (segment === "inside") {
    return "checked_in";
  }

  if (segment === "expected") {
    return "approved";
  }

  return undefined;
}

const PRESET_CONFIG: Record<ResidentEntriesPreset, ResidentEntriesPresetState> = {
  expected: { segment: "expected", datePreset: "today" },
  inside: { segment: "inside", datePreset: "today" },
  all: { segment: "all", datePreset: "all" },
};

export function useResidentEntries(initialPreset: ResidentEntriesPreset = "expected") {
  const initial = PRESET_CONFIG[initialPreset];
  const { flatId, societyId } = useResident();
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [segment, setSegment] = useState<ResidentVisibleEntriesSegment>(initial.segment);
  const [datePreset, setDatePreset] = useState<DateRangePreset>("today");
  const [purpose, setPurpose] = useState<ModelsVisitorPurpose | undefined>();
  const [sheetStatus, setSheetStatus] = useState<ModelsVisitorStatus | undefined>(
    undefined,
  );

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(searchInput.trim()), SEARCH_DEBOUNCE_MS);
    return () => clearTimeout(timer);
  }, [searchInput]);

  useEffect(() => {
    const next = PRESET_CONFIG[initialPreset];
    /* eslint-disable react-hooks/set-state-in-effect -- Route preset changes intentionally reset entry filters. */
    setSegment(next.segment);
    setDatePreset("today");
    setSheetStatus(undefined);
    setPurpose(undefined);
    setSearchInput("");
    setDebouncedSearch("");
    /* eslint-enable react-hooks/set-state-in-effect */
  }, [initialPreset]);

  const hasSearchQuery = debouncedSearch.length > 0;
  const isSearchActive = searchInput.trim().length > 0;
  const hasExplicitFilters = Boolean(purpose || sheetStatus || datePreset !== "today");
  const activeStatus =
    hasExplicitFilters || hasSearchQuery
      ? sheetStatus
      : initial.sheetStatus ?? segmentStatus(segment);
  const shouldSkip = !societyId || !flatId;

  const [fetchEntries] = useLazyGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesQuery();

  const fetchPage = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!societyId || !flatId) {
        return { items: [], total: 0, limit, offset };
      }

      const params = buildResidentEntriesRequestParams({
        datePreset,
        hasExplicitFilters,
        hasSearchQuery,
        initial,
        limit,
        offset,
        purpose,
        search: debouncedSearch || undefined,
        segment,
        sheetStatus,
      });

      const response = await fetchEntries({
        societyId,
        flatId,
        limit,
        offset,
        status: params.status,
        purpose: params.purpose,
        event: params.event,
        eventFrom: params.eventFrom,
        eventTo: params.eventTo,
        search: params.search,
      }).unwrap();

      return {
        items: response.data?.entries ?? [],
        total: response.data?.total ?? 0,
        limit: response.data?.limit ?? limit,
        offset: response.data?.offset ?? offset,
      };
    },
    [
      debouncedSearch,
      fetchEntries,
      flatId,
      datePreset,
      hasExplicitFilters,
      hasSearchQuery,
      initial,
      purpose,
      segment,
      sheetStatus,
      societyId,
    ],
  );

  const pagination = usePaginatedQuery<ModelsVisitorEntry>({
    pageSize: PAGE_SIZE,
    skip: shouldSkip,
    fetchPage,
    getItemKey: visitorEntryKey,
  });

  const activeFilterCount = useMemo(() => {
    let count = 0;
    if (purpose) count += 1;
    if (sheetStatus) count += 1;
    if (datePreset !== "today") count += 1;
    return count;
  }, [datePreset, purpose, sheetStatus]);

  const selectSegment = useCallback((next: ResidentVisibleEntriesSegment) => {
    setSegment(next);
  }, []);

  const applySheetFilters = useCallback(
    (next: {
      datePreset?: DateRangePreset;
      purpose?: ModelsVisitorPurpose;
      status?: ModelsVisitorStatus;
    }) => {
      if (next.datePreset !== undefined) {
        setDatePreset(next.datePreset);
      }
      setPurpose(next.purpose);
      setSheetStatus(next.status);
    },
    [],
  );

  const clearSheetFilters = useCallback(() => {
    setPurpose(undefined);
    setSheetStatus(undefined);
    setDatePreset("today");
  }, []);

  const filterSummary = getVisitorFilterSummary({ datePreset, purpose, status: sheetStatus });
  const mode = hasSearchQuery ? "search" : hasExplicitFilters ? "filtered" : "preset";

  return {
    ...pagination,
    activeFilterCount,
    activeStatus,
    applySheetFilters,
    clearSheetFilters,
    datePreset,
    filterSummary,
    hasExplicitFilters,
    hasSearchQuery: isSearchActive,
    isSearchActive,
    mode,
    purpose,
    searchInput,
    segment,
    selectSegment,
    setDatePreset,
    setSearchInput,
    sheetStatus,
  };
}

export function useResidentEntriesFromParams(presetParam?: string | string[]) {
  return useResidentEntries(parseResidentEntriesPreset(presetParam));
}
