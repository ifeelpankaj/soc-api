import { useCallback, useMemo } from "react";

import { useDebouncedValue } from "@/features/shared/use-debounced-value";
import { usePaginatedQuery } from "@/features/shared/use-paginated-query";
import { getFlatLabel, getVisitorName, titleize } from "@/features/guard/guard-utils";
import {
  generatedApi,
  type ModelsVisitorPendingEntry,
} from "@/lib/api/generated-api";
import { useAppDispatch } from "@/redux/hooks";

const SEARCH_DEBOUNCE_MS = 400;
const PAGE_LIMIT = 20;
const pendingEntryKey = (entry: ModelsVisitorPendingEntry) => entry.id;

function matchesPendingSearch(entry: ModelsVisitorPendingEntry, query: string) {
  if (!query) {
    return true;
  }

  const haystack = [
    getVisitorName(entry),
    entry.visitor?.phone_number,
    entry.visitor?.email,
    getFlatLabel(entry),
    entry.flat?.block,
    entry.flat?.flat_number,
    entry.vehicle_number,
    entry.delivery_partner,
    entry.service_provider,
    titleize(entry.purpose),
    titleize(entry.status),
    entry.primary_resident_name,
  ]
    .filter(Boolean)
    .join(" ")
    .toLowerCase();

  return haystack.includes(query.toLowerCase());
}

export function useGuardPending(societyId?: number, searchInput = "") {
  const debouncedSearch = useDebouncedValue(searchInput.trim(), SEARCH_DEBOUNCE_MS);
  const skip = !societyId;
  const dispatch = useAppDispatch();

  const fetchPage = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const request = dispatch(
        generatedApi.endpoints.getV1SocietiesBySocietyIdVisitorEntriesPending.initiate(
          {
            societyId: societyId ?? 0,
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
    },
    [dispatch, societyId],
  );

  const pagination = usePaginatedQuery<ModelsVisitorPendingEntry>({
    fetchPage,
    pageSize: PAGE_LIMIT,
    skip,
    getItemKey: pendingEntryKey,
  });

  const filteredItems = useMemo(
    () => pagination.items.filter((entry) => matchesPendingSearch(entry, debouncedSearch)),
    [debouncedSearch, pagination.items],
  );

  return useMemo(
    () => ({
      ...pagination,
      items: filteredItems,
      search: searchInput,
      debouncedSearch,
    }),
    [debouncedSearch, filteredItems, pagination, searchInput],
  );
}
