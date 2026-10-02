import { useRouter } from "expo-router";

import type { ModelsUserResponse } from "@/lib/api/generated-api";

import { useAuth } from "./use-auth";
import type { AuthSessionPayload } from "./auth-storage";

export function useCompleteAuth() {
  const router = useRouter();
  const { completeLogin } = useAuth();

  return async (
    payload?: AuthSessionPayload & { user?: ModelsUserResponse | null },
  ) => {
    const route = await completeLogin(payload);
    router.replace(route);
  };
}
