import { useEffect, useState } from "react";
import { AppState, type AppStateStatus } from "react-native";

export function useAppActivePollingInterval(intervalMs: number) {
  const [isActive, setIsActive] = useState(() => AppState.currentState === "active");

  useEffect(() => {
    const subscription = AppState.addEventListener("change", (nextState: AppStateStatus) => {
      setIsActive(nextState === "active");
    });

    return () => {
      subscription.remove();
    };
  }, []);

  return isActive ? intervalMs : 0;
}
