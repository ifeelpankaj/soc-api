import { useCallback, useRef, useState } from "react";
import { useFocusEffect } from "expo-router";
import { mergePaginatedItems } from "@/features/shared/use-paginated-query";
import {
  useLazyMaintenanceBillsQuery,
  type MaintenanceBillsPage,
  type MaintenanceScope,
} from "@/lib/api/maintenance-api";

// Feature state only: RTK owns transport/auth/cache. Residence remounts and
// filter/mode focus effects reset history; generations discard late replies.
export function useMaintenanceHistory(
  args: MaintenanceScope & {
    month?: string;
    status?: string;
    mode: "scroll" | "pages";
  },
) {
  const [fetch] = useLazyMaintenanceBillsQuery();
  const [stored, setStored] = useState<{
    key: string;
    value: MaintenanceBillsPage;
  }>();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>();
  const sequence = useRef(0);
  const flight = useRef(false);
  const failed = useRef({ page: 1, append: false });
  const { societyId, flatId, month, status, mode } = args;
  const queryKey = JSON.stringify([societyId, flatId, month, status, mode]);
  const load = useCallback(
    async (page = 1, append = false) => {
      if (append && flight.current) return;
      const generation = ++sequence.current;
      flight.current = true;
      failed.current = { page, append };
      setBusy(true);
      setError(undefined);
      // Page navigation never labels previous-page rows as the requested page.
      if (!append) setStored(undefined);
      try {
        const result = await fetch(
          { societyId, flatId, month, status, page },
          false,
        ).unwrap();
        if (generation !== sequence.current) return;
        if (
          !Number.isInteger(result.total_count) ||
          !Number.isInteger(result.total_pages) ||
          result.page !== page ||
          typeof result.has_more !== "boolean"
        ) {
          throw new Error(
            "The server does not provide bill pagination metadata.",
          );
        }
        setStored((old) => ({
          key: queryKey,
          value: {
            ...result,
            items:
              append && old?.key === queryKey
                ? mergePaginatedItems(
                    old.value.items,
                    result.items,
                    (bill) => bill.id,
                  )
                : result.items,
          },
        }));
      } catch (e) {
        if (generation === sequence.current) setError(e);
      } finally {
        if (generation === sequence.current) {
          flight.current = false;
          setBusy(false);
        }
      }
    },
    [fetch, societyId, flatId, month, status, queryKey],
  );
  useFocusEffect(
    useCallback(() => {
      void load();
      return () => {
        sequence.current++;
        flight.current = false;
      };
    }, [load]),
  );
  return {
    data: stored?.key === queryKey ? stored.value : undefined,
    busy,
    error,
    load,
    retry: () => load(failed.current.page, failed.current.append),
  };
}
