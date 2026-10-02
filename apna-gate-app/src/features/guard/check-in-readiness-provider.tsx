import { createContext, type PropsWithChildren, useContext, useEffect, useRef } from "react";
import { AppState, Platform, Vibration } from "react-native";

import { useAuth } from "@/features/auth/use-auth";
import { useGuardSociety } from "@/features/guard/guard-context";
import { consumeReadyEntries } from "@/features/guard/check-in-readiness";
import { useAppActivePollingInterval } from "@/features/shared/use-app-active-polling-interval";
import { useGetGuardCheckInReadinessQuery } from "@/lib/api/guard-api-extensions";

const ReadinessContext = createContext<{ readyCount: number; queueKey?: string }>({ readyCount: 0 });

export function GuardCheckInReadinessProvider({ children }: PropsWithChildren) {
  const { selectedSocietyId, user } = useGuardSociety();
  return (
    <ReadinessSession key={`${user?.id}:${selectedSocietyId}`} societyId={selectedSocietyId}>
      {children}
    </ReadinessSession>
  );
}

function ReadinessSession({ children, societyId }: PropsWithChildren<{ societyId?: number }>) {
  const { status } = useAuth();
  const pollingInterval = useAppActivePollingInterval(15_000);
  const seen = useRef(new Set<number>());
  const query = useGetGuardCheckInReadinessQuery(
    { societyId: societyId ?? 0 },
    {
      skip: !societyId || status !== "authenticated" || pollingInterval === 0,
      pollingInterval,
      refetchOnMountOrArgChange: true,
      refetchOnFocus: true,
      refetchOnReconnect: true,
    },
  );

  useEffect(() => {
    if (
      status !== "authenticated" || !pollingInterval ||
      AppState.currentState !== "active" || query.isFetching ||
      query.isError || !query.currentData
    ) return;
    if (consumeReadyEntries(seen.current, query.currentData.entryIds) && Platform.OS !== "web") {
      Vibration.vibrate(200);
    }
  }, [pollingInterval, query.currentData, query.isError, query.isFetching, status]);

  return (
    <ReadinessContext.Provider value={{
      readyCount: query.isError ? 0 : query.currentData?.total ?? 0,
      queueKey: query.isSuccess ? query.currentData?.entryIds.join(",") : undefined,
    }}>
      {children}
    </ReadinessContext.Provider>
  );
}

export function useGuardCheckInReadiness() {
  return useContext(ReadinessContext);
}
