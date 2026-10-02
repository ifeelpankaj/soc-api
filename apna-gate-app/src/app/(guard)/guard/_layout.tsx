import { Platform } from "react-native";
import { Redirect, Stack } from "expo-router";

import { useAuth } from "@/features/auth/use-auth";
import { GuardSocietyProvider } from "@/features/guard/guard-context";
import { GuardCheckInReadinessProvider } from "@/features/guard/check-in-readiness-provider";

export default function GuardStackLayout() {
  const { status } = useAuth();

  if (Platform.OS === "web") return <Redirect href="/select-society" />;

  if (status === "unauthenticated") {
    return <Redirect href="/login" />;
  }

  return (
    <GuardSocietyProvider>
      <GuardCheckInReadinessProvider>
        <Stack screenOptions={{ headerShown: false }}>
          <Stack.Screen name="(tabs)" />
          <Stack.Screen name="scanner" />
          <Stack.Screen name="check-in" />
          <Stack.Screen name="entries" />
          <Stack.Screen name="entries/[entryId]" />
          <Stack.Screen name="add-entry" />
          <Stack.Screen name="pending" />
          <Stack.Screen name="waiting-at-gate" />
          <Stack.Screen name="dashboard" options={{ animation: "none" }} />
          <Stack.Screen name="scan" options={{ animation: "none" }} />
          <Stack.Screen name="logs" options={{ animation: "none" }} />
        </Stack>
      </GuardCheckInReadinessProvider>
    </GuardSocietyProvider>
  );
}
