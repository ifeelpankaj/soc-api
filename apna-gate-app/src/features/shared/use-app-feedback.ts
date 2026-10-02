import { useCallback } from "react";

import { useToast } from "@/components/ui";
import { getFriendlyApiMessage, isSubscriptionError } from "@/features/auth/api-error";

type ActionResult = {
  message: string;
  success: boolean;
};

export function useAppFeedback() {
  const { showToast } = useToast();

  const showError = useCallback(
    (title: string, error: unknown, fallback: string) => {
      showToast({
        title,
        message: getFriendlyApiMessage(error, fallback),
        variant: isSubscriptionError(error) ? "info" : "error",
      });
    },
    [showToast],
  );

  const showSuccess = useCallback(
    (title: string, message: string) => {
      showToast({
        title,
        message,
        variant: "success",
      });
    },
    [showToast],
  );

  const showInfo = useCallback(
    (title: string, message: string) => {
      showToast({
        title,
        message,
        variant: "info",
      });
    },
    [showToast],
  );

  const showActionResult = useCallback(
    (result: ActionResult, { errorTitle, successTitle }: { errorTitle: string; successTitle: string }) => {
      showToast({
        title: result.success ? successTitle : errorTitle,
        message: result.message,
        variant: result.success ? "success" : "error",
      });
    },
    [showToast],
  );

  return {
    showActionResult,
    showError,
    showInfo,
    showSuccess,
  };
}
