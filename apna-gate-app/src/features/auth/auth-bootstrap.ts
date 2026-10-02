import type {
  ModelsBootstrapData,
  ModelsUserResponse,
} from "@/lib/api/generated-api";
import { generatedApi } from "@/lib/api/generated-api";
import { clearAuth, setCredentials } from "@/redux/authSlice";
import type { AppDispatch } from "@/redux/store";
import { baseApi } from "@/redux/queries/baseApi";

import {
  clearTokens,
  getStorageGeneration,
  clearWorkspace,
  extractAuthSession,
  getAccessToken,
  getRefreshToken,
  saveTokens,
  type AuthSessionPayload,
} from "./auth-storage";
import {
  resolveAuthenticatedRoute,
  type RestoreSessionResult,
} from "./auth-session";
import { isAuthenticationFailure } from "./request-policy";

async function fetchProfile(dispatch: AppDispatch) {
  try {
    const response = await dispatch(
      generatedApi.endpoints.getV1AuthProfile.initiate(undefined, {
        forceRefetch: true,
      }),
    ).unwrap();
    const user = response.data?.user ?? null;

    if (!user) {
      throw new Error("Could not load your profile. Please try again.");
    }

    dispatch(setCredentials({ user }));
    return user;
  } catch (error) {
    if (!isAuthenticationFailure(error)) throw error;
    return null;
  }
}

async function refreshAccessToken(dispatch: AppDispatch) {
  const generation = getStorageGeneration();
  const refreshToken = await getRefreshToken();
  if (!refreshToken) {
    return false;
  }

  try {
    const response = await dispatch(
      generatedApi.endpoints.postV1AuthRefresh.initiate({
        modelsRefreshTokenRequest: { refresh_token: refreshToken },
      }),
    ).unwrap();
    const session = extractAuthSession(
      response.data as AuthSessionPayload | undefined,
    );
    if (!session) {
      throw new Error("Could not restore your session. Please try again.");
    }

    await saveTokens(session, generation);
    if (generation !== getStorageGeneration()) return false;
    return true;
  } catch (error) {
    if (!isAuthenticationFailure(error)) throw error;
    return false;
  }
}

async function fetchBootstrap(dispatch: AppDispatch) {
  const response = await dispatch(
    generatedApi.endpoints.getV1Bootstrap.initiate(undefined, {
      forceRefetch: true,
    }),
  ).unwrap();
  return response.data ?? null;
}

export async function restoreSession(
  dispatch: AppDispatch,
): Promise<RestoreSessionResult> {
  const accessToken = await getAccessToken();
  const refreshToken = await getRefreshToken();

  if (!accessToken && !refreshToken) {
    return { status: "unauthenticated" };
  }

  let user: ModelsUserResponse | null = null;

  if (accessToken) {
    user = await fetchProfile(dispatch);
  }

  if (!user) {
    const refreshed = await refreshAccessToken(dispatch);
    if (!refreshed) {
      await clearTokens();
      await clearWorkspace();
      dispatch(clearAuth());
      dispatch(baseApi.util.resetApiState());
      return { status: "unauthenticated" };
    }

    user = await fetchProfile(dispatch);
    if (!user) {
      await clearTokens();
      await clearWorkspace();
      dispatch(clearAuth());
      dispatch(baseApi.util.resetApiState());
      return { status: "unauthenticated" };
    }
  }

  const bootstrap: ModelsBootstrapData | null = await fetchBootstrap(dispatch);

  if (!bootstrap) {
    throw new Error("Could not load your society. Please try again.");
  }

  dispatch(setCredentials({ user: bootstrap.user ?? user }));
  const route = await resolveAuthenticatedRoute(dispatch, bootstrap);

  return {
    status: "authenticated",
    user: bootstrap.user ?? user,
    bootstrap,
    route,
  };
}

export async function completeAuthenticatedSession(
  dispatch: AppDispatch,
  user?: ModelsUserResponse | null,
) {
  if (user) {
    dispatch(setCredentials({ user }));
  }

  const bootstrap = await fetchBootstrap(dispatch);
  if (!bootstrap) {
    throw new Error("Unable to load bootstrap data.");
  }

  dispatch(setCredentials({ user: bootstrap.user ?? user ?? null }));
  const route = await resolveAuthenticatedRoute(dispatch, bootstrap);

  return { bootstrap, route };
}
