import {
  type BaseQueryFn,
  createApi,
  type FetchArgs,
  type FetchBaseQueryError,
  fetchBaseQuery,
} from "@reduxjs/toolkit/query/react";

import {
  clearTokens,
  clearWorkspace,
  getStorageGeneration,
  extractAuthSession,
  getAccessToken,
  getRefreshToken,
  saveTokens,
  type AuthSessionPayload,
} from "@/features/auth/auth-storage";
import type { ModelsUserResponse } from "@/lib/api/generated-api";
import { appConfig } from "@/lib/config";
import { clearAuth, setCredentials } from "@/redux/authSlice";
import {
  API_TIMEOUT_MS,
  isAuthenticationFailure,
  isReadRequest,
} from "@/features/auth/request-policy";

const apiTagTypes = [
  "Auth",
  "Bootstrap",
  "Health",
  "Societies",
  "Society Members",
  "Flats",
  "Flat Residents",
  "Flat Claims",
  "Visitor Entries",
  "Visitor Settings",
  "Flat Visitor Context",
  "Plans",
  "Public",
  "Subscriptions",
] as const;

const PUBLIC_AUTH_PATHS = [
  "/v1/auth/login",
  "/v1/auth/refresh",
  "/v1/auth/forgot-password",
  "/v1/auth/forget-password",
  "/v1/auth/register",
  "/v1/auth/resident/register",
  "/v1/auth/reset-password",
  "/v1/auth/verify-otp",
  "/v1/auth/resend-otp",
];

export const baseQuery = fetchBaseQuery({
  baseUrl: appConfig.apiBaseUrl,
  timeout: API_TIMEOUT_MS,
  credentials: "include",
  prepareHeaders: async (headers) => {
    headers.set("accept", "application/json");

    const token = await getAccessToken();
    if (token) {
      headers.set("authorization", `Bearer ${token}`);
    }

    return headers;
  },
});

function getRequestUrl(args: string | FetchArgs) {
  return typeof args === "string" ? args : args.url;
}

function shouldSkipReauth(args: string | FetchArgs) {
  const url = getRequestUrl(args);
  return PUBLIC_AUTH_PATHS.some((path) => url.includes(path));
}

function isAuthSessionRequest(args: string | FetchArgs) {
  const url = getRequestUrl(args);
  return (
    url.includes("/v1/auth/login") ||
    url.includes("/v1/auth/resident/register") ||
    url.includes("/v1/auth/refresh")
  );
}

function getApiErrorCode(result: { error?: FetchBaseQueryError }) {
  const data = result.error?.data;

  if (!data || typeof data !== "object" || !("error" in data)) {
    return undefined;
  }

  const error = data.error;
  if (!error || typeof error !== "object" || !("code" in error)) {
    return undefined;
  }

  return typeof error.code === "string" ? error.code : undefined;
}

function shouldSyncProfileForForbidden(
  args: string | FetchArgs,
  result: { error?: FetchBaseQueryError },
) {
  const url = getRequestUrl(args);

  return (
    result.error?.status === 403 &&
    getApiErrorCode(result) === "FORBIDDEN" &&
    !shouldSkipReauth(args) &&
    !url.includes("/v1/auth/profile")
  );
}

async function clearSessionAfterAuthFailure(
  api: Parameters<
    BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>
  >[1],
) {
  const cleanup = Promise.allSettled([clearTokens(), clearWorkspace()]);
  api.dispatch(clearAuth());
  const results = await cleanup;
  if (results.some((result) => result.status === "rejected")) {
    console.warn("Could not remove saved session");
  }
}

async function syncProfileSession(
  api: Parameters<
    BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>
  >[1],
  extraOptions: Parameters<
    BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError>
  >[2],
) {
  const generation = getStorageGeneration();
  const profileResult = await baseQuery(
    { url: "/v1/auth/profile" },
    api,
    extraOptions,
  );

  if (
    generation !== getStorageGeneration() ||
    profileResult.error ||
    !profileResult.data
  ) {
    return false;
  }

  const data = profileResult.data as {
    data?: { user?: ModelsUserResponse | null };
  };
  const user = data.data?.user ?? null;

  if (user) {
    api.dispatch(setCredentials({ user }));
  }

  return true;
}

async function persistAuthResponse(data: unknown, generation: number) {
  if (!data || typeof data !== "object" || !("data" in data)) {
    return;
  }

  const payload = (data as { data?: AuthSessionPayload }).data;
  const session = extractAuthSession(payload ?? null);
  if (session) {
    await saveTokens(session, generation);
  }
}

let refreshPromise: Promise<boolean | FetchBaseQueryError> | null = null;
let refreshFailed = false;
let sessionGeneration = 0;

export const resetRefreshState = () => {
  sessionGeneration += 1;
  refreshPromise = null;
  refreshFailed = false;
};

export const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError,
  { skipForbiddenRetry?: boolean }
> = async (args, api, extraOptions) => {
  const storageGeneration = getStorageGeneration();
  let result = await baseQuery(args, api, extraOptions);
  // A background request can discover revocation before the successful change response arrives.
  if (!result.error && getRequestUrl(args) === "/v1/auth/change-password")
    return result;
  if (storageGeneration !== getStorageGeneration())
    return { error: { status: "CUSTOM_ERROR", error: "Session ended" } };
  const code = getApiErrorCode(result);
  if (code === "SESSION_REVOKED") {
    await clearSessionAfterAuthFailure(api);
    return result;
  }
  if (
    getRequestUrl(args) === "/v1/auth/change-password" &&
    code === "INVALID_CREDENTIALS"
  )
    return result;

  if (!result.error && isAuthSessionRequest(args)) {
    resetRefreshState();
    await persistAuthResponse(result.data, storageGeneration);
  }

  if (
    !extraOptions?.skipForbiddenRetry &&
    isReadRequest(args) &&
    shouldSyncProfileForForbidden(args, result)
  ) {
    const profileSynced = await syncProfileSession(api, extraOptions);

    if (profileSynced) {
      result = await baseQuery(args, api, extraOptions);
    }
  }

  if (result.error?.status === 401 && !shouldSkipReauth(args)) {
    if (refreshFailed) {
      await clearSessionAfterAuthFailure(api);
      return result;
    }

    if (!refreshPromise) {
      const generation = sessionGeneration;
      refreshPromise = (async (): Promise<boolean | FetchBaseQueryError> => {
        const refreshToken = await getRefreshToken();
        const refreshResult = await baseQuery(
          {
            url: "/v1/auth/refresh",
            method: "POST",
            body: refreshToken ? { refresh_token: refreshToken } : undefined,
          },
          api,
          extraOptions,
        );

        if (generation !== sessionGeneration) return false;
        if (refreshResult.error || !refreshResult.data) {
          if (isAuthenticationFailure(refreshResult.error)) {
            refreshFailed = true;
            await clearSessionAfterAuthFailure(api);
            return false;
          }
          return (
            refreshResult.error ?? {
              status: "CUSTOM_ERROR",
              error: "Could not restore your session. Please try again.",
            }
          );
        }

        await persistAuthResponse(refreshResult.data, storageGeneration);
        if (
          generation !== sessionGeneration ||
          storageGeneration !== getStorageGeneration()
        )
          return false;
        refreshFailed = false;
        return true;
      })().finally(() => {
        setTimeout(() => {
          if (generation === sessionGeneration) refreshPromise = null;
        }, 100);
      });
    }

    const refreshSucceeded = await refreshPromise;
    if (refreshSucceeded === true && isReadRequest(args)) {
      result = await baseQuery(args, api, extraOptions);
    } else if (refreshSucceeded && typeof refreshSucceeded === "object") {
      return { error: refreshSucceeded };
    }
  }

  if (storageGeneration !== getStorageGeneration())
    return { error: { status: "CUSTOM_ERROR", error: "Session ended" } };
  return result;
};

export const baseApi = createApi({
  reducerPath: "api",
  baseQuery: baseQueryWithReauth,
  tagTypes: apiTagTypes,
  keepUnusedDataFor: 120,
  refetchOnFocus: true,
  refetchOnReconnect: true,
  endpoints: () => ({}),
});
