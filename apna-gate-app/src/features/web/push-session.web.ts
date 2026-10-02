import { getApp, getApps, initializeApp } from "firebase/app";
import { deleteToken, getMessaging, getToken, isSupported, onMessage, type MessagePayload } from "firebase/messaging";
import { firebaseWebConfig, hasWebPushConfig, webVapidKey } from "./firebase-config";
import { needsWebPushInstallation } from "./device";
import { clearAlertBannerDismissal } from "./notification-banner-session";

let generation = 0;
let queue: Promise<unknown> = Promise.resolve();
let registration: ServiceWorkerRegistration | undefined;
function serial<T>(work: () => Promise<T>): Promise<T> {
  const result = queue.then(work, work); queue = result.catch(() => undefined); return result;
}
export function browserPushGeneration() { return generation; }
const messaging = () => getMessaging(getApps().length ? getApp() : initializeApp(firebaseWebConfig));

export async function browserPushAvailability() {
  if (needsWebPushInstallation()) return "install" as const;
  if (!hasWebPushConfig()) return "setup" as const;
  if (!window.isSecureContext || !(await isSupported())) return "unsupported" as const;
  return Notification.permission;
}

async function sendIdentity(userID: number | null) {
  if (!registration) return;
  const worker = registration.active;
  if (!worker) throw new Error("Notification worker is not active yet");
  await new Promise<void>((resolve, reject) => {
    const channel = new MessageChannel();
    const timeout = setTimeout(() => { channel.port1.close(); reject(new Error("Notification worker did not respond")); }, 5000);
    channel.port1.onmessage = (event) => { clearTimeout(timeout); channel.port1.close(); if (event.data?.ok) resolve(); else reject(new Error("Notification storage unavailable")); };
    worker.postMessage({ type: "APNA_PUSH_SESSION", userID }, [channel.port2]);
  });
}

export async function registerBrowserPush(userID: number, epoch: number, register: (token: string, installationID: string) => Promise<void>, unregister: (token: string) => Promise<void>) {
  return serial(async () => {
    if (epoch !== generation || Notification.permission !== "granted") return null;
    registration = await navigator.serviceWorker.register("/firebase-messaging-sw.js", { scope: "/", updateViaCache: "none" });
    registration = await navigator.serviceWorker.ready;
    // Disable display until the authenticated registration has completed.
    await sendIdentity(null);
    const token = await getToken(messaging(), { vapidKey: webVapidKey, serviceWorkerRegistration: registration });
    if (epoch !== generation) { await deleteToken(messaging()); return null; }
    let id = localStorage.getItem("apna_web_installation_id");
    if (!id) { id = crypto.randomUUID(); localStorage.setItem("apna_web_installation_id", id); }
    await register(token, id);
    if (epoch !== generation) {
      await Promise.allSettled([unregister(token), deleteToken(messaging())]); return null;
    }
    await sendIdentity(userID);
    if (epoch !== generation) { await sendIdentity(null); return null; }
    return token;
  });
}

export async function listenBrowserPush(callback: (payload: MessagePayload) => void) {
  if (!hasWebPushConfig() || !(await isSupported())) return () => {};
  return onMessage(messaging(), callback);
}

export function clearBrowserPushSession() {
  clearAlertBannerDismissal();
  generation++;
  // Disabling display is independent from token deletion/network availability.
  const disable = (async () => {
    if (typeof navigator === "undefined" || !navigator.serviceWorker) return;
    registration ??= await navigator.serviceWorker.getRegistration("/");
    await sendIdentity(null);
    const notifications = await registration?.getNotifications();
    notifications?.forEach((item) => item.close());
  })();
  const remove = serial(async () => {
    await disable.catch(() => undefined);
    if (typeof window !== "undefined" && hasWebPushConfig() && await isSupported()) await deleteToken(messaging());
  });
  return Promise.allSettled([disable, remove]).then((results) => {
    if (results.some((item) => item.status === "rejected")) console.warn("Browser notification cleanup could not finish; delivery remains session-bound.");
  });
}
