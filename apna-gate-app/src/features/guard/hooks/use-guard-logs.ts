import { useCallback, useEffect, useMemo, useState } from "react";

import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import {
  buildExplicitVisitorEntryQueryFilters,
  buildVisitorEntryQueryFilters,
  getHalfOpenRangeIST,
  getVisitorFilterSummary,
  type DateRangePreset,
  type LogsPreset,
  type LogsSegment,
  parseGuardEntriesPreset,
  presetToInitialState,
} from "@/features/guard/guard-routes";
import { usePaginatedQuery } from "@/features/shared/use-paginated-query";
import {
  generatedApi,
  type ModelsVisitorEntry,
  type ModelsVisitorPurpose,
  type ModelsVisitorStatus,
} from "@/lib/api/generated-api";
import { useLazyGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery } from "@/lib/api/guard-api-extensions";
import { useAppDispatch } from "@/redux/hooks";

const PAGE_SIZE = 20;
const SEARCH_DEBOUNCE_MS = 400;
const visitorEntryKey = (entry: ModelsVisitorEntry) => entry.id;

function segmentStatus(segment: LogsSegment): ModelsVisitorStatus | undefined {
  if (segment === "inside") {
    return "checked_in";
  }

  if (segment === "expected") {
    return "approved";
  }

  return undefined;
}

export function useGuardLogs(initialPreset: LogsPreset = "today") {
  const initial = presetToInitialState(initialPreset);
  const { selectedSocietyId } = useGuardScreen();
  const dispatch = useAppDispatch();
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [segment, setSegment] = useState<LogsSegment>(initial.segment);
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
    const next = presetToInitialState(initialPreset);
    /* eslint-disable react-hooks/set-state-in-effect -- Route preset changes intentionally reset the log filters. */
    setSegment(next.segment);
    setDatePreset("today");
    setSheetStatus(undefined);
    setPurpose(undefined);
    setSearchInput("");
    setDebouncedSearch("");
    /* eslint-enable react-hooks/set-state-in-effect */
  }, [initialPreset]);

  const hasSearchQuery = debouncedSearch.length > 0;
  const hasExplicitFilters = Boolean(purpose || sheetStatus || datePreset !== "today");
  const queryFilters = useMemo(() => {
    if (hasSearchQuery || hasExplicitFilters) {
      return buildExplicitVisitorEntryQueryFilters({
        applyDateRange: hasExplicitFilters,
        datePreset,
        purpose,
        status: sheetStatus,
      });
    }

    return buildVisitorEntryQueryFilters({
      segment,
      sheetStatus: initial.sheetStatus,
      purpose: undefined,
      datePreset: initial.datePreset,
      isSearchActive: false,
    });
  }, [
    datePreset,
    hasExplicitFilters,
    hasSearchQuery,
    initial.datePreset,
    initial.sheetStatus,
    purpose,
    segment,
    sheetStatus,
  ]);

  const statsRange = useMemo(() => {
    if (hasSearchQuery) {
      return undefined;
    }

    const preset =
      hasExplicitFilters && datePreset !== "all" ? datePreset : initial.datePreset;
    return getHalfOpenRangeIST(preset === "all" ? "today" : preset);
  }, [datePreset, hasExplicitFilters, hasSearchQuery, initial.datePreset]);

  const [fetchEntries] = useLazyGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery();

  const usesExpectedGuestsApi =
    !hasSearchQuery &&
    !hasExplicitFilters &&
    (segment === "expected" || initial.sheetStatus === "approved");

  const fetchPage = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!selectedSocietyId) {
        return { items: [], total: 0, limit, offset };
      }

      if (usesExpectedGuestsApi) {
        const todayRange = getHalfOpenRangeIST("today");
        const request = dispatch(
          generatedApi.endpoints.getV1SocietiesBySocietyIdVisitorEntriesExpectedGuests.initiate(
            {
              societyId: selectedSocietyId,
              search: debouncedSearch || undefined,
              eventFrom: todayRange?.eventFrom,
              eventTo: todayRange?.eventTo,
              limit,
              offset,
            },
            { forceRefetch: true },
          ),
        );

        const response = await request.unwrap();
        request.unsubscribe();

        return {
          items: response.data?.entries ?? [],
          total: response.data?.total ?? 0,
          limit: response.data?.limit ?? limit,
          offset: response.data?.offset ?? offset,
        };
      }

      const response = await fetchEntries({
        societyId: selectedSocietyId,
        limit,
        offset,
        search: debouncedSearch || undefined,
        status: queryFilters.status,
        purpose: queryFilters.purpose,
        event: queryFilters.event,
        eventFrom: queryFilters.eventFrom,
        eventTo: queryFilters.eventTo,
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
      dispatch,
      fetchEntries,
      queryFilters.event,
      queryFilters.eventFrom,
      queryFilters.eventTo,
      queryFilters.purpose,
      queryFilters.status,
      selectedSocietyId,
      usesExpectedGuestsApi,
    ],
  );

  const pagination = usePaginatedQuery<ModelsVisitorEntry>({
    pageSize: PAGE_SIZE,
    skip: !selectedSocietyId,
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

  const selectSegment = useCallback((next: LogsSegment) => {
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

  const activeStatus =
    hasExplicitFilters || hasSearchQuery
      ? sheetStatus
      : initial.sheetStatus ?? segmentStatus(segment);
  const mode = hasSearchQuery ? "search" : hasExplicitFilters ? "filtered" : "preset";
  const filterSummary = getVisitorFilterSummary({ datePreset, purpose, status: sheetStatus });

  return {
    ...pagination,
    activeFilterCount,
    activeStatus,
    applySheetFilters,
    clearSheetFilters,
    datePreset,
    filterSummary,
    hasExplicitFilters,
    hasSearchQuery: searchInput.trim().length > 0,
    isSearchActive: searchInput.trim().length > 0,
    mode,
    purpose,
    searchInput,
    segment,
    selectSegment,
    setDatePreset,
    setSearchInput,
    sheetStatus,
    statsRange,
  };
}

export function useGuardLogsFromParams(presetParam?: string | string[]) {
  return useGuardLogs(parseGuardEntriesPreset(presetParam));
}
