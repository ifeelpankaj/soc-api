import {
  createContext,
  type PropsWithChildren,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import type { Href } from "expo-router";
import { useRouter } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { SafeAreaView } from "react-native-safe-area-context";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { colors } from "@/theme/colors";
import { View } from "react-native";

import type { ModelsUserResponse } from "@/lib/api/generated-api";
import { clearWorkspaceSelection } from "@/redux/appSlice";

import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { baseApi, resetRefreshState } from "@/redux/queries/baseApi";

import { completeAuthenticatedSession, restoreSession } from "./auth-bootstrap";
import {
  clearTokens,
  clearWorkspace,
  extractAuthSession,
  saveTokens,
  type AuthSessionPayload,
} from "./auth-storage";
import { setDevicePushToken } from "@/redux/notificationSlice";
import { clearAuth } from "@/redux/authSlice";
import { pendingNotificationRoute } from "@/features/web/notification-handoff";

export type AuthStatus = "loading" | "authenticated" | "unauthenticated";

type AuthContextValue = {
  status: AuthStatus;
  user: ModelsUserResponse | null;
  homeRoute: Href | null;
  completeLogin: (
    payload?: AuthSessionPayload & { user?: ModelsUserResponse | null },
  ) => Promise<Href>;
  signOutLocal: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

void SplashScreen.preventAutoHideAsync().catch(() => undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.auth.user);
  const [status, setStatus] = useState<AuthStatus>("loading");
  const [homeRoute, setHomeRoute] = useState<Href | null>(null);
  const [restoreError, setRestoreError] = useState(false);
  const [restoreAttempt, setRestoreAttempt] = useState(0);

  const signOutLocal = useCallback(async () => {
    resetRefreshState();
    const cleanup = Promise.allSettled([clearTokens(), clearWorkspace()]);
    dispatch(clearAuth());
    dispatch(setDevicePushToken(null));
    dispatch(clearWorkspaceSelection());
    dispatch(baseApi.util.resetApiState());
    setHomeRoute(null);
    setStatus("unauthenticated");
    const results = await cleanup;
    if (results.some((result) => result.status === "rejected")) {
      console.warn("Session storage cleanup failed after local sign-out.");
    }
  }, [dispatch]);

  const completeLogin = useCallback(
    async (
      payload?: AuthSessionPayload & { user?: ModelsUserResponse | null },
    ) => {
      const session = extractAuthSession(payload ?? null);
      if (session) {
        await saveTokens(session);
      }

      const result = await completeAuthenticatedSession(
        dispatch,
        payload?.user ?? null,
      );
      setHomeRoute(result.route);
      setStatus("authenticated");
      return pendingNotificationRoute() ?? result.route;
    },
    [dispatch],
  );

  useEffect(() => {
    let active = true;

    (async () => {
      try {
        const result = await restoreSession(dispatch);

        if (!active) {
          return;
        }

        if (result.status === "authenticated") {
          setHomeRoute(result.route);
          setStatus("authenticated");
        } else {
          setHomeRoute(null);
          setStatus("unauthenticated");
        }
      } catch {
        if (active) setRestoreError(true);
      } finally {
        await SplashScreen.hideAsync().catch(() => undefined);
      }
    })();

    return () => {
      active = false;
    };
  }, [dispatch, restoreAttempt]);

  const value = useMemo<AuthContextValue>(
    () => ({
      status: status === "authenticated" && !user ? "unauthenticated" : status,
      user,
      homeRoute,
      completeLogin,
      signOutLocal,
    }),
    [completeLogin, homeRoute, signOutLocal, status, user],
  );

  return (
    <AuthContext.Provider value={value}>
      {restoreError ? (
        <SafeAreaView
          style={{ flex: 1, backgroundColor: colors.surface.screen }}
        >
          <View
            style={{ flex: 1, justifyContent: "center", padding: 24, gap: 16 }}
          >
            <EmptyState
              title="Let's reconnect"
              message="We couldn't load your account. Check your connection and try again. Your saved session is still on this device."
            />
            <Button
              title="Try again"
              onPress={() => {
                setRestoreError(false);
                setRestoreAttempt((attempt) => attempt + 1);
              }}
            />
          </View>
        </SafeAreaView>
      ) : (
        children
      )}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return context;
}

export function useAuthRedirect(route: Href) {
  const router = useRouter();
  const { homeRoute, status } = useAuth();

  useEffect(() => {
    if (status === "authenticated" && homeRoute) {
      router.replace(pendingNotificationRoute() ?? homeRoute ?? route);
    }
  }, [homeRoute, route, router, status]);
}
