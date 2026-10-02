export type WebDevice = "apple" | "android" | "desktop";

export type BrowserDevice = {
  userAgent: string;
  platform: string;
  maxTouchPoints: number;
};

export function isIPhoneOrIPad(device?: BrowserDevice) {
  return !!device && (
    /iPad|iPhone|iPod/i.test(device.userAgent) ||
    (device.platform === "MacIntel" && device.maxTouchPoints > 1)
  );
}

export function isMac(device?: BrowserDevice) {
  return !!device && !isIPhoneOrIPad(device) && (
    /^Mac/i.test(device.platform) || /Macintosh/i.test(device.userAgent)
  );
}

export function isAppleWebBrowser(device?: BrowserDevice) {
  return isIPhoneOrIPad(device) || isMac(device);
}

// Presentation only. Membership authorization always belongs to the API.
// "apple" here means mobile Apple: Macs do not need Home Screen push guidance.
export function detectWebDevice(device?: BrowserDevice): WebDevice {
  if (!device) return "desktop";
  if (isIPhoneOrIPad(device)) return "apple";
  return /Android/i.test(device.userAgent) ? "android" : "desktop";
}

export function isStandalone() {
  return typeof window !== "undefined" &&
    (window.matchMedia("(display-mode: standalone)").matches ||
      (navigator as Navigator & { standalone?: boolean }).standalone === true);
}

export function requiresWebPushInstallation(device: WebDevice, standalone: boolean) {
  return device === "apple" && !standalone;
}

export function needsWebPushInstallation() {
  return typeof navigator !== "undefined" &&
    requiresWebPushInstallation(detectWebDevice(navigator), isStandalone());
}
