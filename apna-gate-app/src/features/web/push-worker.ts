/// <reference lib="webworker" />
import { initializeApp } from "firebase/app";
import { getMessaging, onBackgroundMessage } from "firebase/messaging/sw";
import { firebaseWebConfig } from "./firebase-config";
declare const self: ServiceWorkerGlobalScope;

function database() {
  return new Promise<IDBDatabase>((resolve, reject) => {
    const request = indexedDB.open("apna-push", 1);
    request.onupgradeneeded = () => request.result.createObjectStore("state");
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}
async function state(key: string, write?: unknown) {
  const db = await database();
  try { return await new Promise<unknown>((resolve,reject) => {
    const tx = db.transaction("state", write === undefined ? "readonly" : "readwrite");
    const request = write === undefined ? tx.objectStore("state").get(key) : tx.objectStore("state").put(write,key);
    tx.oncomplete = () => resolve(request.result);
    tx.onerror = () => reject(tx.error); tx.onabort = () => reject(tx.error);
  }); } finally { db.close(); }
}

self.addEventListener("install", (event) => event.waitUntil(self.skipWaiting()));
self.addEventListener("activate", (event) => event.waitUntil(self.clients.claim()));
self.addEventListener("message", (event) => {
  if (event.data?.type !== "APNA_PUSH_SESSION") return;
  event.waitUntil(state("user", event.data.userID ?? null).then(() => event.ports[0]?.postMessage({ok:true}), () => event.ports[0]?.postMessage({ok:false})));
});
// Own click handling before initializing Firebase's worker messaging.
self.addEventListener("notificationclick", (event) => {
  event.stopImmediatePropagation(); event.notification.close();
  event.waitUntil((async () => {
    const data = event.notification.data;
    if (!data?.id || String(await state("user")) !== data.recipient) return;
    const url = new URL("/",self.location.origin);
    url.searchParams.set("notification_id",data.id); url.searchParams.set("recipient",data.recipient);
    const windows = await self.clients.matchAll({type:"window",includeUncontrolled:true});
    const existing = windows.find((client) => new URL(client.url).origin === self.location.origin) as WindowClient | undefined;
    if (existing) { await existing.navigate(url.href); await existing.focus(); }
    else await self.clients.openWindow(url.href);
  })());
});

if (firebaseWebConfig.apiKey && firebaseWebConfig.appId) {
  onBackgroundMessage(getMessaging(initializeApp(firebaseWebConfig)), async (payload) => {
    const data = payload.data;
    if (!data?.notification_id || !/^[a-f0-9-]{36}$/i.test(data.notification_id) || !data.user_id) return;
    if (String(await state("user")) !== data.user_id) return;
    const seen = (await state("seen") as string[] | undefined) ?? [];
    if (seen.includes(data.notification_id)) return;
    await self.registration.showNotification(data.title || "Apna Gate", {
      body: data.body || "Open your inbox for an update.", icon:"/icon-192.png",
      tag:data.notification_id, data:{id:data.notification_id,recipient:data.user_id},
    });
    await state("seen", [...seen,data.notification_id].slice(-200));
  });
}
// Deliberately no fetch/cache handler: authenticated responses are never cached here.
