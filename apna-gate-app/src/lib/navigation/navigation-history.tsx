import { useAuth } from "@/features/auth/use-auth";
import type { PropsWithChildren } from "react";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { type Href, useUnstableGlobalHref } from "expo-router";

import { isMeaningfulBackRoute } from "@/lib/navigation/navigation-policy";

type NavigationHistoryValue = {
  consumePreviousRoute: (currentRoute?: Href | null) => Href | null;
  currentRoute: Href | null;
  previousRoute: Href | null;
};

const NavigationHistoryContext = createContext<NavigationHistoryValue>({
  consumePreviousRoute: () => null,
  currentRoute: null,
  previousRoute: null,
});

const MAX_HISTORY = 30;

export function NavigationHistoryProvider({ children }: PropsWithChildren) {
  const href = useUnstableGlobalHref();
  const { status } = useAuth();
  const historyRef = useRef<string[]>([]);
  const [history, setHistory] = useState<string[]>([]);

  useEffect(() => {
    if (status !== "authenticated") {
      historyRef.current = [];
      return;
    }
    if (!href) {
      return;
    }
    const current = historyRef.current;
    if (current[current.length - 1] === href) {
      return;
    }
    const next = [...current, href].slice(-MAX_HISTORY);
    historyRef.current = next;
    setHistory(next);
  }, [href, status]);

  const consumePreviousRoute = useCallback((currentRoute?: Href | null) => {
    const current = historyRef.current;
    if (current.length < 2) {
      return null;
    }

    const activeRoute = String(currentRoute ?? current[current.length - 1]);
    const previousIndex = findPreviousRouteIndex(current, activeRoute);
    if (previousIndex < 0) {
      return null;
    }

    const previous = current[previousIndex];
    const next = current.slice(0, previousIndex + 1);
    historyRef.current = next;
    setHistory(next);
    return previous as Href;
  }, []);

  const value = useMemo<NavigationHistoryValue>(
    () => ({
      consumePreviousRoute,
      currentRoute:
        status === "authenticated" && history.length > 0
          ? (history[history.length - 1] as Href)
          : null,
      previousRoute:
        status === "authenticated" && history.length >= 2
          ? (history[history.length - 2] as Href)
          : null,
    }),
    [consumePreviousRoute, history, status],
  );

  return (
    <NavigationHistoryContext.Provider value={value}>
      {children}
    </NavigationHistoryContext.Provider>
  );
}

export function useNavigationHistory() {
  return useContext(NavigationHistoryContext);
}

function findPreviousRouteIndex(history: string[], currentRoute: string) {
  for (let index = history.length - 2; index >= 0; index -= 1) {
    if (isMeaningfulBackRoute(history[index], currentRoute)) {
      return index;
    }
  }

  return -1;
}
