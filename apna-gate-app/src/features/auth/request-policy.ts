export const API_TIMEOUT_MS = 30_000;

export function isReadRequest(args: string | { method?: string }) {
  return (
    typeof args === "string" ||
    ["GET", "HEAD"].includes((args.method ?? "GET").toUpperCase())
  );
}

export function isAuthenticationFailure(error: unknown) {
  return Boolean(
    error &&
    typeof error === "object" &&
    "status" in error &&
    (error.status === 401 || error.status === 403),
  );
}
