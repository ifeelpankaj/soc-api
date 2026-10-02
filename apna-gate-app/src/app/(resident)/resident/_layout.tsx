import { Redirect, Stack } from "expo-router";

import { useAuth } from "@/features/auth/use-auth";
import { ResidentProvider } from "@/features/resident/resident-context";

export default function ResidentStackLayout() {
  const { status } = useAuth();

  if (status === "unauthenticated") {
    return <Redirect href="/login" />;
  }

  return (
    <ResidentProvider>
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="(tabs)" />
        <Stack.Screen name="logs" />
        <Stack.Screen name="entries/index" />
        <Stack.Screen name="entries/[entryId]" />
        <Stack.Screen name="members/index" />
        <Stack.Screen name="members/add" />
        <Stack.Screen name="invites/index" />
        <Stack.Screen name="invites/visitor/[inviteId]" />
        <Stack.Screen name="invites/member/[inviteId]" />
        <Stack.Screen name="visitors/index" />
        <Stack.Screen name="visitors/invite" />
        <Stack.Screen name="visitors/settings" />
        <Stack.Screen name="hub/announcements" />
        <Stack.Screen name="hub/community" />
        <Stack.Screen name="hub/create" />
        <Stack.Screen name="hub/edit/[postId]" />
        <Stack.Screen name="hub/posts/[postId]" />
      </Stack>
    </ResidentProvider>
  );
}
