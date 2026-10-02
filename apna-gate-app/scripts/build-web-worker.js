const { build } = require("esbuild");
require("@expo/env").load(process.cwd());
const keys = ["API_KEY", "AUTH_DOMAIN", "PROJECT_ID", "MESSAGING_SENDER_ID", "APP_ID", "VAPID_KEY"];
const define = Object.fromEntries(keys.map((key) => {
  const name = `EXPO_PUBLIC_FIREBASE_WEB_${key}`;
  return [`process.env.${name}`, JSON.stringify(process.env[name] || "")];
}));
define["process.env.EXPO_PUBLIC_WEB_PUSH_ENABLED"] = JSON.stringify(process.env.EXPO_PUBLIC_WEB_PUSH_ENABLED || "false");
build({ entryPoints: ["src/features/web/push-worker.ts"], outfile: "public/firebase-messaging-sw.js", bundle: true, format: "iife", platform: "browser", target: ["safari16.4"], minify: true, define }).catch(() => process.exit(1));
