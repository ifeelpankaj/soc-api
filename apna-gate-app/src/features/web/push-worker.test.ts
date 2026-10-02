import assert from "node:assert/strict";
import test from "node:test";
import vm from "node:vm";
import { build } from "esbuild";

const worker = build({
  entryPoints: ["src/features/web/push-worker.ts"], bundle: true, write: false,
  format: "iife", platform: "browser",
  plugins: [{
    name: "isolate-push-click",
    setup(builder) {
      builder.onResolve({ filter: /^(firebase\/|\.\/firebase-config)/ }, (args) => ({ path: args.path, namespace: "stub" }));
      builder.onLoad({ filter: /.*/, namespace: "stub" }, () => ({
        contents: "export const firebaseWebConfig = {}; export const initializeApp = () => {}; export const getMessaging = () => {}; export const onBackgroundMessage = () => {};",
      }));
    },
  }],
});

test("notification clicks use root handoff for closed and open web apps", async () => {
  const code = (await worker).outputFiles[0].text;
  for (const warm of [false, true]) {
    for (const user of ["1", "2", null]) {
      const calls: string[] = [];
      const handlers = new Map<string, (event: unknown) => void>();
      const origin = "https://app.example.test";
      vm.runInNewContext(code, {
        URL,
        self: {
          location: { origin },
          addEventListener: (name: string, handler: (event: unknown) => void) => handlers.set(name, handler),
          clients: {
            matchAll: async () => warm ? [{
              url: `${origin}/resident/dashboard`,
              navigate: async (url: string) => { calls.push(url); },
              focus: async () => { calls.push("focus"); },
            }] : [],
            openWindow: async (url: string) => { calls.push(url); },
          },
        },
        indexedDB: {
          open: () => {
            const request = { result: {
              close() {},
              transaction() {
                const tx = { oncomplete: () => {}, objectStore: () => ({ get: () => ({ result: user }) }) };
                queueMicrotask(() => tx.oncomplete());
                return tx;
              },
            }, onsuccess: () => {} };
            queueMicrotask(() => request.onsuccess());
            return request;
          },
        },
      });
      let work: Promise<void> | undefined;
      const id = "53bdc1b1-a7bb-41bb-9522-3a17e98c45c8";
      handlers.get("notificationclick")!({
        stopImmediatePropagation() {},
        notification: { data: { id, recipient: "1" }, close() {} },
        waitUntil: (promise: Promise<void>) => { work = promise; },
      });
      await work;
      assert.deepEqual(calls, user === "1"
        ? [`${origin}/?notification_id=${id}&recipient=1`, ...(warm ? ["focus"] : [])]
        : []);
    }
  }
});
