import { useCallback, useEffect, useMemo, useState } from "react";

import type { DateRangePreset } from "@/features/visitors/visitor-date-ranges";
import type { ModelsVisitorPurpose, ModelsVisitorStatus } from "@/lib/api/generated-api";

const SEARCH_DEBOUNCE_MS = 400;

type VisitorEntriesFilterState<TSegment extends string> = {
  datePreset: DateRangePreset;
  segment: TSegment;
  sheetStatus?: ModelsVisitorStatus;
};

type UseVisitorEntriesFiltersOptions<TSegment extends string> = {
  initial: VisitorEntriesFilterState<TSegment>;
  onPresetChange?: (presetKey: string) => VisitorEntriesFilterState<TSegment>;
  presetKey?: string;
};

export function useVisitorEntriesFilters<TSegment extends string>({
  initial,
  onPresetChange,
  presetKey,
}: UseVisitorEntriesFiltersOptions<TSegment>) {
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [segment, setSegment] = useState<TSegment>(initial.segment);
  const [datePreset, setDatePreset] = useState<DateRangePreset>(initial.datePreset);
  const [purpose, setPurpose] = useState<ModelsVisitorPurpose | undefined>();
  const [sheetStatus, setSheetStatus] = useState<ModelsVisitorStatus | undefined>(
    initial.sheetStatus,
  );

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(searchInput.trim()), SEARCH_DEBOUNCE_MS);
    return () => clearTimeout(timer);
  }, [searchInput]);

  useEffect(() => {
    if (!onPresetChange || presetKey === undefined) {
      return;
    }

    const next = onPresetChange(presetKey);
    /* eslint-disable react-hooks/set-state-in-effect -- Route preset changes intentionally reset visitor filters. */
    setSegment(next.segment);
    setDatePreset(next.datePreset);
    setSheetStatus(next.sheetStatus);
    setPurpose(undefined);
    setSearchInput("");
    setDebouncedSearch("");
    /* eslint-enable react-hooks/set-state-in-effect */
  }, [onPresetChange, presetKey]);

  const isSearchActive = debouncedSearch.length > 0;

  const selectSegment = useCallback((next: TSegment) => {
    setSegment(next);
    setSheetStatus(undefined);

    if (next === "expected") {
      setDatePreset("today");
    } else if (next === "all") {
      setDatePreset("all");
    } else if (next === "today") {
      setDatePreset("today");
    }
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
    setSegment("today" as TSegment);
  }, []);

  const activeFilterCount = useMemo(() => {
    let count = 0;
    if (purpose) count += 1;
    if (sheetStatus) count += 1;
    if (datePreset !== "today" && segment !== "all") count += 1;
    return count;
  }, [datePreset, purpose, segment, sheetStatus]);

  return {
    activeFilterCount,
    applySheetFilters,
    clearSheetFilters,
    datePreset,
    debouncedSearch,
    isSearchActive,
    purpose,
    searchInput,
    segment,
    selectSegment,
    setDatePreset,
    setSearchInput,
    sheetStatus,
  };
}
