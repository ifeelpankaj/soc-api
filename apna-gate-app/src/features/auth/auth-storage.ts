import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";
import { SessionWriteQueue } from "./session-write-queue";
const tokenWrites = new SessionWriteQueue();
export const getStorageGeneration = () => tokenWrites.current();

export const AUTH_ACCESS_TOKEN_KEY = "auth.accessToken";
export const AUTH_REFRESH_TOKEN_KEY = "auth.refreshToken";
export const AUTH_ACCESS_EXPIRES_AT_KEY = "auth.accessExpiresAt";
export const AUTH_REFRESH_EXPIRES_AT_KEY = "auth.refreshExpiresAt";
export const WORKSPACE_FLAT_ID_KEY = "workspace.selectedFlatId";
export const WORKSPACE_SOCIETY_ID_KEY = "workspace.selectedSocietyId";

export type StoredAuthTokens = {
  accessToken: string;
  refreshToken?: string | null;
  accessExpiresAt?: string | null;
  refreshExpiresAt?: string | null;
};

export type StoredWorkspace = {
  flatId?: number | null;
  societyId?: number | null;
};

function getWebStorage() {
  if (typeof window === "undefined") {
    return null;
  }

  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

async function setItem(key: string, value: string) {
  if (Platform.OS === "web") {
    getWebStorage()?.setItem(key, value);
    return;
  }

  await SecureStore.setItemAsync(key, value);
}

async function getItem(key: string) {
  if (Platform.OS === "web") {
    return getWebStorage()?.getItem(key) ?? null;
  }

  return SecureStore.getItemAsync(key);
}

async function deleteItem(key: string) {
  if (Platform.OS === "web") {
    getWebStorage()?.removeItem(key);
    return;
  }

  await SecureStore.deleteItemAsync(key);
}

export function saveTokens(
  tokens: StoredAuthTokens,
  generation = getStorageGeneration(),
) {
  return tokenWrites.save(generation, async () => {
    await setItem(AUTH_ACCESS_TOKEN_KEY, tokens.accessToken);

    if (tokens.refreshToken) {
      await setItem(AUTH_REFRESH_TOKEN_KEY, tokens.refreshToken);
    }

    if (tokens.accessExpiresAt) {
      await setItem(AUTH_ACCESS_EXPIRES_AT_KEY, tokens.accessExpiresAt);
    }

    if (tokens.refreshExpiresAt) {
      await setItem(AUTH_REFRESH_EXPIRES_AT_KEY, tokens.refreshExpiresAt);
    }
  });
}

export async function getAccessToken() {
  return getItem(AUTH_ACCESS_TOKEN_KEY);
}

export async function getRefreshToken() {
  return getItem(AUTH_REFRESH_TOKEN_KEY);
}

export function clearTokens() {
  return tokenWrites.clear(async () => {
    const results = await Promise.allSettled([
      deleteItem(AUTH_ACCESS_TOKEN_KEY),
      deleteItem(AUTH_REFRESH_TOKEN_KEY),
      deleteItem(AUTH_ACCESS_EXPIRES_AT_KEY),
      deleteItem(AUTH_REFRESH_EXPIRES_AT_KEY),
    ]);
    if (results.some((result) => result.status === "rejected"))
      throw new Error("Could not remove saved session");
  });
}

export async function saveWorkspace(workspace: StoredWorkspace) {
  if (workspace.flatId != null) {
    await setItem(WORKSPACE_FLAT_ID_KEY, String(workspace.flatId));
    await deleteItem(WORKSPACE_SOCIETY_ID_KEY);
    return;
  }

  if (workspace.societyId != null) {
    await setItem(WORKSPACE_SOCIETY_ID_KEY, String(workspace.societyId));
    await deleteItem(WORKSPACE_FLAT_ID_KEY);
  }
}

export async function loadWorkspace(): Promise<StoredWorkspace> {
  const flatRaw = await getItem(WORKSPACE_FLAT_ID_KEY);
  const societyRaw = await getItem(WORKSPACE_SOCIETY_ID_KEY);

  const flatId = flatRaw ? Number(flatRaw) : null;
  const societyId = societyRaw ? Number(societyRaw) : null;

  return {
    flatId: Number.isFinite(flatId) ? flatId : null,
    societyId: Number.isFinite(societyId) ? societyId : null,
  };
}

export async function clearWorkspace() {
  await Promise.all([
    deleteItem(WORKSPACE_FLAT_ID_KEY),
    deleteItem(WORKSPACE_SOCIETY_ID_KEY),
  ]);
}

export type AuthSessionPayload = {
  access_token?: string;
  refresh_token?: string;
  access_token_expires_at?: string;
  refresh_token_expires_at?: string;
};

export function extractAuthSession(
  data?: AuthSessionPayload | null,
): StoredAuthTokens | null {
  if (!data?.access_token) {
    return null;
  }

  return {
    accessToken: data.access_token,
    refreshToken: data.refresh_token ?? null,
    accessExpiresAt: data.access_token_expires_at ?? null,
    refreshExpiresAt: data.refresh_token_expires_at ?? null,
  };
}
