import { useCallback, useEffect, useMemo, useRef, useState } from "react";

type PaginatedPayload<T> = {
  items: T[];
  total: number;
  limit: number;
  offset: number;
};

type UsePaginatedQueryOptions<T> = {
  pageSize?: number;
  skip?: boolean;
  fetchPage: (args: { limit: number; offset: number }) => Promise<PaginatedPayload<T>>;
  getItemKey?: (item: T) => string | number | undefined;
};

export function mergePaginatedItems<T>(
  current: T[],
  incoming: T[],
  getItemKey?: (item: T) => string | number | undefined,
) {
  if (!getItemKey) return [...current, ...incoming];
  const positions = new Map<string | number, number>();
  const merged = [...current];
  current.forEach((item, index) => {
    const key = getItemKey(item);
    if (key !== undefined) positions.set(key, index);
  });
  for (const item of incoming) {
    const key = getItemKey(item);
    const position = key === undefined ? undefined : positions.get(key);
    if (position === undefined) {
      if (key !== undefined) positions.set(key, merged.length);
      merged.push(item);
    } else {
      merged[position] = item;
    }
  }
  return merged;
}

export function usePaginatedQuery<T>({
  pageSize = 20,
  skip = false,
  fetchPage,
  getItemKey,
}: UsePaginatedQueryOptions<T>) {
  const [offset, setOffset] = useState(0);
  const [items, setItems] = useState<T[]>([]);
  const [total, setTotal] = useState(0);
  const [loadedThrough, setLoadedThrough] = useState(0);
  const [isFetchingFirstPage, setIsLoading] = useState(false);
  const [hasSettledFirstPage, setHasSettledFirstPage] = useState(false);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [error, setError] = useState<unknown>();
  const loadMoreInFlight = useRef(false);
  const requestSequence = useRef(0);

  const hasMore = loadedThrough < total;
  const isLoading = !skip && (!hasSettledFirstPage || isFetchingFirstPage);

  const loadPage = useCallback(
    async (nextOffset: number, mode: "replace" | "append") => {
      if (skip) {
        return;
      }

      const setLoading = nextOffset === 0 ? setIsLoading : setIsLoadingMore;
      const requestId = ++requestSequence.current;
      setLoading(true);
      setError(undefined);

      try {
        const page = await fetchPage({ limit: pageSize, offset: nextOffset });
        if (requestId !== requestSequence.current) return;
        setTotal(page.total);
        setOffset(nextOffset);
        setLoadedThrough(nextOffset + page.items.length);
        setItems((current) =>
          mode === "append"
            ? mergePaginatedItems(current, page.items, getItemKey)
            : page.items,
        );
      } catch (nextError) {
        if (requestId === requestSequence.current) setError(nextError);
      } finally {
        if (requestId !== requestSequence.current) return;
        if (nextOffset === 0) {
          setHasSettledFirstPage(true);
        }
        setLoading(false);
        setIsLoadingMore(false);
        setIsRefreshing(false);
      }
    },
    [fetchPage, getItemKey, pageSize, skip],
  );

  useEffect(() => {
    if (skip) {
      requestSequence.current += 1;
      /* eslint-disable react-hooks/set-state-in-effect -- Skip means the current page cache must be cleared immediately. */
      setItems([]);
      setTotal(0);
      setOffset(0);
      setLoadedThrough(0);
      setHasSettledFirstPage(false);
      setIsLoading(false);
      setIsLoadingMore(false);
      setIsRefreshing(false);
      /* eslint-enable react-hooks/set-state-in-effect */
      return;
    }

    void loadPage(0, "replace");
  }, [loadPage, skip]);

  const refresh = useCallback(async () => {
    setIsRefreshing(true);
    await loadPage(0, "replace");
  }, [loadPage]);

  const loadMore = useCallback(async () => {
    if (
      skip ||
      isLoading ||
      isLoadingMore ||
      loadMoreInFlight.current ||
      !hasMore
    ) {
      return;
    }
    loadMoreInFlight.current = true;
    try {
      await loadPage(loadedThrough, "append");
    } finally {
      loadMoreInFlight.current = false;
    }
  }, [hasMore, isLoading, isLoadingMore, loadPage, loadedThrough, skip]);

  return useMemo(
    () => ({
      items,
      total,
      offset,
      isLoading,
      isLoadingMore,
      isRefreshing,
      hasMore,
      error,
      refresh,
      loadMore,
    }),
    [
      error,
      hasMore,
      isLoading,
      isLoadingMore,
      isRefreshing,
      items,
      loadMore,
      offset,
      refresh,
      total,
    ],
  );
}
