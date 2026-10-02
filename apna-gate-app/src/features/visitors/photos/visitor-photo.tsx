import {
  useCallback,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { ActivityIndicator, Animated, AppState, StyleSheet, Text, View } from "react-native";
import { Image } from "expo-image";
import { useFocusEffect } from "expo-router";
import { useGetVisitorPhotoQuery } from "@/lib/api/photo-api";
import { useAppSelector } from "@/redux/hooks";
import { colors } from "@/theme/colors";
import { usePhotoAccess } from "./photo-access";
import {
  photoCacheKey,
  photoRefreshDelay,
  type PhotoVariant,
} from "./photo-utils";

type EntryPhoto = {
  id?: number;
  visitor?: { photo_url?: string; full_name?: string };
};

export function VisitorPhoto({
  entry,
  variant = "list",
  fallback,
  size,
  display = "standard",
}: {
  entry: EntryPhoto;
  variant?: PhotoVariant;
  fallback?: ReactNode;
  size?: number;
  display?: "standard" | "banner" | "fullscreen";
}) {
  const access = usePhotoAccess();
  const [openedAt] = useState(() => Date.now());
  const userId = useAppSelector((state) => state.auth.user?.id);
  const [focused, setFocused] = useState(false);
  useFocusEffect(
    useCallback(() => {
      setFocused(true);
      return () => setFocused(false);
    }, []),
  );
  const [active, setActive] = useState(
    AppState.currentState == null || AppState.currentState === "active",
  );
  const photoReference = entry.visitor?.photo_url?.trim();
  const enabled = Boolean(
    userId &&
    access?.societyId &&
    entry.id &&
    photoReference &&
    focused &&
    active,
  );
  const query = useGetVisitorPhotoQuery(
    {
      userId: userId ?? 0,
      societyId: access?.societyId ?? 0,
      entryId: entry.id ?? 0,
      context: access?.context ?? "guard",
      variant,
      photoReference: photoReference ?? "",
    },
    {
      skip: !enabled,
      refetchOnMountOrArgChange: false,
      refetchOnFocus: false,
      refetchOnReconnect: true,
    },
  );
  const view = query.currentData?.data;

  useEffect(() => {
    const subscription = AppState.addEventListener("change", (state) => {
      setActive(state === "active");
    });
    return () => subscription.remove();
  }, []);

  const { refetch, isFetching, isError } = query;
  useEffect(() => {
    if (!enabled || !view || isFetching || isError) return;
    const timer = setTimeout(
      () => void refetch(),
      photoRefreshDelay(view.expires_at),
    );
    return () => clearTimeout(timer);
  }, [enabled, isError, isFetching, refetch, view]);

  const identity = photoCacheKey({
    userId: userId ?? 0,
    societyId: access?.societyId ?? 0,
    entryId: entry.id ?? 0,
    context: access?.context ?? "guard",
    variant,
    photoReference: photoReference ?? "",
  });
  const accessIdentity = JSON.stringify([
    userId ?? 0,
    access?.societyId ?? 0,
    access?.context ?? "guard",
    entry.id ?? 0,
  ]);
  // currentData remains available during refetches. Keep using it even when a
  // refresh fails so an already-decoded private bitmap never flashes away.
  const expiredFullscreen = display === "fullscreen" && view &&
    Date.parse(view.expires_at) <= openedAt;
  const url = userId && access && view && !expiredFullscreen ? view.url : undefined;
  return (
    <PhotoFrame
      identity={identity}
      accessIdentity={accessIdentity}
      hasPhoto={Boolean(photoReference)}
      url={url}
      loading={enabled && !view && query.isFetching}
      detail={variant === "detail"}
      banner={display === "banner"}
      fullscreen={display === "fullscreen"}
      queryError={query.isError}
      size={size ?? (variant === "avatar" ? 160 : 48)}
      fallback={
        fallback ?? (
          <Text style={styles.initial}>
            {entry.visitor?.full_name?.charAt(0).toUpperCase() || "?"}
          </Text>
        )
      }
    />
  );
}

function PhotoFrame({
  identity,
  accessIdentity,
  hasPhoto,
  url,
  loading,
  detail,
  banner,
  fullscreen,
  queryError,
  fallback,
  size,
}: {
  identity: string;
  accessIdentity: string;
  hasPhoto: boolean;
  url?: string;
  loading: boolean;
  detail: boolean;
  banner: boolean;
  fullscreen: boolean;
  queryError: boolean;
  fallback: ReactNode;
  size: number;
}) {
  const [displayed, setDisplayed] = useState<{
    identity: string;
    accessIdentity: string;
    url: string;
  }>();
  const [failedIdentity, setFailedIdentity] = useState<string>();
  const [aspectRatio, setAspectRatio] = useState(1);
  const [fade] = useState(() => new Animated.Value(0));
  const safeDisplayed =
    hasPhoto && displayed?.accessIdentity === accessIdentity
      ? displayed
      : undefined;
  const loaded = safeDisplayed?.identity === identity;
  const failed = failedIdentity === JSON.stringify([identity, url]);

  useEffect(() => {
    if (!loaded) fade.setValue(0);
    return () => fade.stopAnimation();
  }, [fade, identity, loaded]);

  const rememberLoaded = (nextURL: string) => {
    setFailedIdentity(undefined);
    setDisplayed({ identity, accessIdentity, url: nextURL });
  };

  const revealReplacement = (nextURL: string) => {
    Animated.timing(fade, {
      duration: 150,
      toValue: 1,
      useNativeDriver: true,
    }).start(({ finished }) => {
      if (finished) {
        rememberLoaded(nextURL);
        fade.setValue(0);
      }
    });
  };

  const updateAspectRatio = (width: number, height: number) => {
    if (height) setAspectRatio(width / height);
  };
  return (
    <View
      accessibilityLabel="Visitor photo"
      style={[
        styles.frame,
        { width: size, height: size },
        banner && styles.banner,
        fullscreen && styles.fullscreen,
        detail &&
          !banner && !fullscreen &&
          (url || loading) && { width: "100%", height: undefined, aspectRatio },
      ]}
    >
      {safeDisplayed ? (
        <Image
          source={{
            uri: loaded && url ? url : safeDisplayed.url,
            cacheKey: safeDisplayed.identity,
          }}
          style={StyleSheet.absoluteFill}
          contentFit={fullscreen || (detail && !banner) ? "contain" : "cover"}
          contentPosition={banner ? "top center" : "center"}
          cachePolicy="memory"
          recyclingKey={safeDisplayed.identity}
          onLoad={(event) => {
            if (loaded && url) rememberLoaded(url);
            updateAspectRatio(event.source.width, event.source.height);
          }}
        />
      ) : null}
      {url && !failed && !loaded ? (
        <Animated.View
          pointerEvents="none"
          style={[StyleSheet.absoluteFill, safeDisplayed && { opacity: fade }]}
        >
          <Image
            source={{ uri: url, cacheKey: identity }}
            style={StyleSheet.absoluteFill}
            contentFit={fullscreen || (detail && !banner) ? "contain" : "cover"}
            contentPosition={banner ? "top center" : "center"}
            cachePolicy="memory"
            recyclingKey={identity}
            onError={() => setFailedIdentity(JSON.stringify([identity, url]))}
            onLoad={(event) => {
              updateAspectRatio(event.source.width, event.source.height);
              if (safeDisplayed) revealReplacement(url);
              else rememberLoaded(url);
            }}
          />
        </Animated.View>
      ) : null}
      {!safeDisplayed && (!url || failed || !loaded) ? (
        fullscreen ? (
          failed || (queryError && !url) ? (
            <Text style={styles.viewerMessage}>Couldn’t load photo</Text>
          ) : (
            <ActivityIndicator accessibilityLabel="Loading visitor photo" color="#FFFFFF" />
          )
        ) : fallback
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  frame: {
    width: 48,
    height: 48,
    borderRadius: 12,
    overflow: "hidden",
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.surface.secondary,
  },
  banner: {
    borderRadius: 0,
    height: "100%",
    width: "100%",
  },
  fullscreen: {
    width: "100%",
    height: "100%",
    borderRadius: 0,
    backgroundColor: "#000000",
  },
  viewerMessage: { color: "#FFFFFF", fontSize: 16 },
  initial: { color: colors.brand.orange, fontSize: 20, fontWeight: "700" },
});
