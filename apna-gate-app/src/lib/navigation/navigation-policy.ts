import type { Href } from "expo-router";

const ENTRY_SETUP_ROUTES = new Set(["/", "/login", "/forget", "/select-society"]);
const ROLE_HOME_ROUTES = new Set([
  "/resident/dashboard",
  "/guard/home",
  "/guard/dashboard",
]);
const REDIRECT_ALIAS_ROUTES = new Set(["/resident/logs", "/guard/logs", "/guard/scan"]);

export function routePath(route?: Href | string | null) {
  if (!route) {
    return null;
  }

  if (typeof route === "string") {
    return normalizePath(route);
  }

  if (typeof route === "object" && "pathname" in route) {
    return normalizePath(String(route.pathname));
  }

  return normalizePath(String(route));
}

export function isEntrySetupRoute(route?: Href | string | null) {
  const path = routePath(route);
  return Boolean(path && ENTRY_SETUP_ROUTES.has(path));
}

export function isRoleHomeRoute(route?: Href | string | null) {
  const path = routePath(route);
  return Boolean(path && ROLE_HOME_ROUTES.has(path));
}

export function isRedirectAliasRoute(route?: Href | string | null) {
  const path = routePath(route);
  return Boolean(path && REDIRECT_ALIAS_ROUTES.has(path));
}

export function isSameRoute(route?: Href | string | null, currentRoute?: Href | string | null) {
  const path = routePath(route);
  const currentPath = routePath(currentRoute);
  return Boolean(path && currentPath && path === currentPath);
}

export function isMeaningfulBackRoute(route: Href | string, currentRoute?: Href | string | null) {
  if (isSameRoute(route, currentRoute)) {
    return false;
  }

  if (isEntrySetupRoute(route) || isRedirectAliasRoute(route)) {
    return false;
  }

  if (isRoleHomeRoute(currentRoute) && isRoleHomeRoute(route)) {
    return false;
  }

  return true;
}

export function shouldUseNativeBack(
  currentRoute?: Href | string | null,
  trackedPreviousRoute?: Href | string | null,
) {
  if (isRoleHomeRoute(currentRoute)) {
    return false;
  }

  if (trackedPreviousRoute && !isMeaningfulBackRoute(trackedPreviousRoute, currentRoute)) {
    return false;
  }

  return true;
}

function normalizePath(route: string) {
  const [withoutHash] = route.split("#");
  const [path] = withoutHash.split("?");
  return path.replace(/\/+$/, "") || "/";
}
