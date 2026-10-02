export const SUBSCRIPTION_ERROR_CODES = new Set([
  "SUBSCRIPTION_EXPIRED",
  "SUBSCRIPTION_REQUIRED",
  "SUBSCRIPTION_QUOTA_EXCEEDED",
]);

export const SUBSCRIPTION_EXPIRED_BANNER_MESSAGE =
  "Your society subscription has expired. Please contact your society admin to renew the plan.";

export function getApiMessage(error: unknown, fallback: string) {
  if (error && typeof error === "object" && "status" in error) {
    if (error.status === "TIMEOUT_ERROR") return "This is taking longer than expected. Please try again.";
    if (error.status === "FETCH_ERROR") return "Could not connect. Check your connection and try again.";
  }
  if (!error || typeof error !== "object" || !("data" in error)) {
    return fallback;
  }

  const data = error.data;
  if (!data || typeof data !== "object") {
    return fallback;
  }

  if ("message" in data && typeof data.message === "string") {
    return data.message;
  }

  if ("error" in data && data.error && typeof data.error === "object") {
    const nested = data.error;
    if ("message" in nested && typeof nested.message === "string") {
      return nested.message;
    }
  }

  return fallback;
}

export function getApiErrorCode(error: unknown): string | undefined {
  if (!error || typeof error !== "object" || !("data" in error)) {
    return undefined;
  }

  const data = error.data;
  if (!data || typeof data !== "object" || !("error" in data)) {
    return undefined;
  }

  const nested = data.error;
  if (!nested || typeof nested !== "object" || !("code" in nested)) {
    return undefined;
  }

  return typeof nested.code === "string" ? nested.code : undefined;
}

export function isSubscriptionError(error: unknown): boolean {
  if (!error || typeof error !== "object") {
    return false;
  }

  if ("status" in error && error.status === 402) {
    return true;
  }

  const code = getApiErrorCode(error);
  return code != null && SUBSCRIPTION_ERROR_CODES.has(code);
}

export function getFriendlyApiMessage(error: unknown, fallback: string) {
  if (isSubscriptionError(error)) {
    return SUBSCRIPTION_EXPIRED_BANNER_MESSAGE;
  }

  return getApiMessage(error, fallback);
}

export function getVisitorActionErrorMessage(error: unknown, fallback: string) {
  const code = getApiErrorCode(error);

  if (code === "VISITOR_SOURCE_DISABLED") {
    return "Your admin has not allowed guard entry for this society.";
  }

  if (code === "VISITOR_GUARD_ON_BEHALF_NOT_ALLOWED") {
    return "Your admin has not allowed guard approval on behalf of residents.";
  }

  if (error && typeof error === "object" && "status" in error && error.status === 403) {
    return "Only an active flat resident can perform this action.";
  }

  return getApiMessage(error, fallback);
}
