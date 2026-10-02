// Public browser identifiers only. Never put Admin credentials in EXPO_PUBLIC_*.
export const firebaseWebConfig = {
  apiKey: process.env.EXPO_PUBLIC_FIREBASE_WEB_API_KEY,
  authDomain: process.env.EXPO_PUBLIC_FIREBASE_WEB_AUTH_DOMAIN,
  projectId: process.env.EXPO_PUBLIC_FIREBASE_WEB_PROJECT_ID,
  messagingSenderId: process.env.EXPO_PUBLIC_FIREBASE_WEB_MESSAGING_SENDER_ID,
  appId: process.env.EXPO_PUBLIC_FIREBASE_WEB_APP_ID,
};
export const webVapidKey = process.env.EXPO_PUBLIC_FIREBASE_WEB_VAPID_KEY;
export function hasWebPushConfig() {
  return Boolean(process.env.EXPO_PUBLIC_WEB_PUSH_ENABLED === "true" && firebaseWebConfig.apiKey && firebaseWebConfig.projectId &&
    firebaseWebConfig.messagingSenderId && firebaseWebConfig.appId && webVapidKey);
}
