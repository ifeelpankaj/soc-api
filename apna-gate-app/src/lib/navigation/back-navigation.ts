import type { Href } from "expo-router";

import {
  isMeaningfulBackRoute,
  shouldUseNativeBack,
} from "@/lib/navigation/navigation-policy";

type BackRouter = {
  back: () => void;
  canGoBack: () => boolean;
  dismissTo?: (href: Href) => void;
  replace: (href: Href) => void;
};

type BackHistory = {
  consumePreviousRoute?: (currentRoute?: Href | null) => Href | null;
  currentRoute?: Href | null;
  previousRoute?: Href | null;
};

export function handleBack(
  router: BackRouter,
  fallbackRoute: Href,
  history?: BackHistory,
) {
  const currentRoute = history?.currentRoute ?? null;
  const trackedPreviousRoute = history?.previousRoute ?? null;
  const previousRoute =
    history?.consumePreviousRoute?.(currentRoute) ?? trackedPreviousRoute ?? null;

  if (previousRoute && isMeaningfulBackRoute(previousRoute, currentRoute)) {
    if (router.dismissTo) {
      router.dismissTo(previousRoute);
      return;
    }

    router.replace(previousRoute);
    return;
  }

  if (shouldUseNativeBack(currentRoute, trackedPreviousRoute) && router.canGoBack()) {
    router.back();
    return;
  }

  router.replace(fallbackRoute);
}
