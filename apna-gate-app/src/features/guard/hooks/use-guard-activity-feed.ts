import { useCallback } from "react";

import { getTodayRange } from "@/features/visitors/visitor-date-ranges";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { useVisitorActivityFeed } from "@/features/visitors/hooks/use-visitor-activity-feed";
import { useLazyGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery } from "@/lib/api/guard-api-extensions";

export function useGuardActivityFeed() {
  const { selectedSocietyId } = useGuardScreen();
  const [fetchEntries] = useLazyGetV1SocietiesBySocietyIdVisitorEntriesExtendedQuery();

  const fetchPage = useCallback(
    async ({ limit, offset }: { limit: number; offset: number }) => {
      if (!selectedSocietyId) {
        return { items: [], total: 0, limit, offset };
      }

      const todayRange = getTodayRange();
      const response = await fetchEntries({
        societyId: selectedSocietyId,
        event: todayRange ? "activity" : undefined,
        eventFrom: todayRange?.eventFrom,
        eventTo: todayRange?.eventTo,
        limit,
        offset,
      }).unwrap();

      return {
        items: response.data?.entries ?? [],
        total: response.data?.total ?? 0,
        limit: response.data?.limit ?? limit,
        offset: response.data?.offset ?? offset,
      };
    },
    [fetchEntries, selectedSocietyId],
  );

  return useVisitorActivityFeed({
    skip: !selectedSocietyId,
    fetchPage,
  });
}

export {
  ACTIVITY_LIST_HEIGHT,
  ACTIVITY_ROW_HEIGHT,
  ACTIVITY_VISIBLE_ROWS,
} from "@/features/visitors/hooks/use-visitor-activity-feed";
