import type { Href } from "expo-router";
import { useRouter } from "expo-router";
import { useCallback } from "react";

import { handleBack } from "@/lib/navigation/back-navigation";
import { useNavigationHistory } from "@/lib/navigation/navigation-history";

export function useBackAction(fallbackRoute: Href) {
  const router = useRouter();
  const history = useNavigationHistory();

  return useCallback(() => {
    handleBack(router, fallbackRoute, history);
  }, [fallbackRoute, history, router]);
}
