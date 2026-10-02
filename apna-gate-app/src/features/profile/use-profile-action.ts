import { useCallback } from "react";

import { useToast } from "@/components/ui";

export function useProfileAction() {
  const { showToast } = useToast();

  const showComingSoon = useCallback(
    (feature?: string) => {
      showToast({
        title: "Coming soon",
        message: feature
          ? `${feature} is not available yet.`
          : "This feature is not available yet.",
        variant: "info",
      });
    },
    [showToast],
  );

  return { showComingSoon };
}
