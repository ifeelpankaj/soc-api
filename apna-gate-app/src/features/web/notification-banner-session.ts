const STORAGE_KEY = "apna_alert_banner_dismissed_user";
const listeners = new Set<() => void>();
let dismissedUser: number | null = null;

export function subscribeAlertBanner(listener: () => void) {
  listeners.add(listener);
  return () => { listeners.delete(listener); };
}

export function isAlertBannerDismissed(userID?: number) {
  if (!userID) return false;
  try {
    return sessionStorage.getItem(STORAGE_KEY) === String(userID) || dismissedUser === userID;
  } catch {
    return dismissedUser === userID;
  }
}

export function dismissAlertBanner(userID: number) {
  dismissedUser = userID;
  try { sessionStorage.setItem(STORAGE_KEY, String(userID)); } catch { /* Keep in memory if storage is unavailable. */ }
  listeners.forEach((listener) => listener());
}

export function clearAlertBannerDismissal() {
  dismissedUser = null;
  try { sessionStorage.removeItem(STORAGE_KEY); } catch { /* Storage may be unavailable. */ }
  listeners.forEach((listener) => listener());
}
