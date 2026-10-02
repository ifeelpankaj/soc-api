import { type BrowserDevice, detectWebDevice, isAppleWebBrowser } from "./device";

export function resolveWebAccess(flag: string | undefined, device?: BrowserDevice) {
  if (flag !== "true") return "allowed";
  if (!device) return "loading";
  if (isAppleWebBrowser(device)) return "allowed";
  return detectWebDevice(device) === "android" ? "android" : "desktop";
}

export function getGooglePlayURL() {
  return process.env.EXPO_PUBLIC_GOOGLE_PLAY_URL?.trim() || undefined;
}
