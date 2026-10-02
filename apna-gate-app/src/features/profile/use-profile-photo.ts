import { useCallback, useEffect, useState } from "react";
import { AppState } from "react-native";
import { useFocusEffect } from "expo-router";
import { useAppSelector } from "@/redux/hooks";
import { useGetProfilePhotoQuery } from "@/lib/api/photo-api";
import { photoRefreshDelay } from "@/features/visitors/photos/photo-utils";
import { normalizeAvatarReference } from "./avatar-reference";

export function useProfilePhoto(legacyUrl?: string | null) {
  const userId = useAppSelector((state) => state.auth.user?.id);
  const avatarReference = normalizeAvatarReference(legacyUrl);
  const [focused, setFocused] = useState(false);
  const [active, setActive] = useState(AppState.currentState !== "background");
  const [now, setNow] = useState(Date.now);
  useFocusEffect(
    useCallback(() => {
      setFocused(true);
      return () => setFocused(false);
    }, []),
  );
  useEffect(() => {
    const sub = AppState.addEventListener("change", (value) => {
      setNow(Date.now());
      setActive(value === "active");
    });
    return () => sub.remove();
  }, []);
  const enabled = Boolean(userId && avatarReference && focused && active);
  const { currentData, isFetching, isError, refetch, error } =
    useGetProfilePhotoQuery(
      { userId: userId ?? 0, variant: "avatar" },
      { skip: !enabled, refetchOnMountOrArgChange: true },
    );
  useEffect(() => {
    if (!enabled || !currentData?.data || isFetching || isError) return;
    const timer = setTimeout(
      () => void refetch(),
      photoRefreshDelay(currentData.data.expires_at),
    );
    return () => clearTimeout(timer);
  }, [currentData, enabled, isError, isFetching, refetch]);
  useEffect(() => {
    if (!currentData?.data) return;
    const timer = setTimeout(
      () => setNow(Date.now()),
      Math.max(0, Date.parse(currentData.data.expires_at) - Date.now()),
    );
    return () => clearTimeout(timer);
  }, [currentData]);
  if (!userId || !avatarReference) return undefined;
  if (currentData?.data && Date.parse(currentData.data.expires_at) > now)
    return currentData.data.url;
  // The photo endpoint returns 404 for legacy, unmanaged profile images.
  return error && "status" in error && error.status === 404
    ? avatarReference
    : undefined;
}
