import { useCallback, useEffect, useMemo, useRef } from "react";
import { useGuardCheckInReadiness } from "@/features/guard/check-in-readiness-provider";

import { useDebouncedValue } from "@/features/shared/use-debounced-value";
import { usePaginatedQuery } from "@/features/shared/use-paginated-query";
import { generatedApi, type ModelsVisitorEntry } from "@/lib/api/generated-api";
import { useAppDispatch } from "@/redux/hooks";

const SEARCH_DEBOUNCE_MS = 400;
const PAGE_LIMIT = 20;
const visitorEntryKey = (entry: ModelsVisitorEntry) => entry.id;

export function useGuardWaitingAtGate(societyId?: number, searchInput = "") {
  const { queueKey } = useGuardCheckInReadiness();
  const previousQueueKey = useRef<string | undefined>(undefined);
  const debouncedSearch = useDebouncedValue(searchInput.trim(), SEARCH_DEBOUNCE_MS);
  const skip = !societyId;
  const dispatch = useAppDispatch();

  const fetchPage = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      const request = dispatch(
        generatedApi.endpoints.getV1SocietiesBySocietyIdVisitorEntriesWaitingAtGate.initiate(
          {
            societyId: societyId ?? 0,
            search: debouncedSearch || undefined,
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
    [debouncedSearch, dispatch, societyId],
  );

  const pagination = usePaginatedQuery<ModelsVisitorEntry>({
    fetchPage,
    pageSize: PAGE_LIMIT,
    skip,
    getItemKey: visitorEntryKey,
  });

  const { refresh } = pagination;
  useEffect(() => {
    if (queueKey === undefined || skip) return;
    const previous = previousQueueKey.current;
    previousQueueKey.current = queueKey;
    if (previous !== queueKey) void refresh();
  }, [queueKey, refresh, skip]);

  return useMemo(
    () => ({
      ...pagination,
      search: searchInput,
      debouncedSearch,
    }),
    [debouncedSearch, pagination, searchInput],
  );
}
