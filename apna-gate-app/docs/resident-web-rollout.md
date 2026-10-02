# Resident browser access

Resident access at `https://resident.apnagate.org` supports an optional Apple-only web mode. With `EXPO_PUBLIC_APPLE_WEB_ENABLED=true`, iPhone, iPad (including desktop-mode Safari), and Mac browsers can enter the resident app. Android browsers see **Get the app**; Windows/Linux browsers see **Get the Android app** with instructions to use an Android phone. The gate runs before authentication, notification providers, and routing, including direct links and browsers with saved sessions. Device detection controls presentation; API membership checks remain authoritative. Only active resident memberships and residences are offered, including for mixed guard/resident accounts. Direct guard routes remain blocked on web.

Missing or false leaves browser access unrestricted. Native apps bypass the gate regardless of the flag.

| Platform | Flag false or missing | Flag true |
| --- | --- | --- |
| iPhone / iPad browsers | Web allowed | Web allowed |
| Mac Safari / Chrome | Web allowed | Web allowed |
| Android browser | Web allowed | Get the app |
| Windows / Linux browsers | Web allowed | Get the Android app |
| Android / iOS native | Current behavior | Current behavior |

Android and iOS native apps skip browser detection and the browser installation banner. Native registration still checks platform, physical-device availability, and notification permission. Android Expo/FCM delivery is unchanged. The backend currently skips native iOS delivery; enabling APNs is separate work.

## Configuration and deployment

1. Apply API migrations 23 (user session versions, if not already deployed) and 24 (`device_tokens.session_version`). Deploy the session-aware backend to **all** instances before enabling web delivery. Old/unbound web registrations remain excluded.
2. In the existing Firebase project, add a **Web app** and enable Cloud Messaging's Web Push certificates. Copy its public configuration and public VAPID key to the app's documented `EXPO_PUBLIC_FIREBASE_WEB_*` variables. Keep the existing Firebase Admin credential file on the API server only. Enable the FCM Registration API for older Firebase projects if necessary. No Firebase credentials or public web configuration have been provisioned by this code change.
3. Host the resident app at `https://resident.apnagate.org/`. Configure API CORS for this exact origin and the API's existing HTTPS public URL. Keep the admin portal and API domains unchanged. `EXPO_PUBLIC_WEB_BASE_URL` is the existing invite website URL; it does not set the resident origin.
4. For Apple-only web, set `EXPO_PUBLIC_APPLE_WEB_ENABLED=true` and `EXPO_PUBLIC_GOOGLE_PLAY_URL=https://play.google.com/store/apps/details?id=com.apnagate` in the web build environment. The download destination comes exclusively from that environment value; no URL is hardcoded in the app. If missing or blank, the button is omitted and the page says "The download link is currently unavailable." Both values are public and embedded into the bundle, so they must not contain secrets. Run `npm ci` then `npm run web:export`; deploy `dist` to the resident origin. Rebuild and redeploy web after changing either value. The export command bundles `public/firebase-messaging-sw.js` with the same public config as the app. `npm run web` also builds the worker. After changing Firebase configuration, rebuild both app and worker.
5. Serve the worker at `/firebase-messaging-sw.js`, with JavaScript MIME type, scope `/`, and `Cache-Control: no-cache, max-age=0, must-revalidate`. Do not rewrite that request to HTML, redirect it across origins, or apply immutable asset caching. `public/_headers` supplies these headers for hosts supporting that format; configure equivalent headers on other hosts. Serve the manifest and icons at their root URLs. Configure extensionless Expo routes, including `/notification-open`, to their exported HTML pages.
6. For staged physical-device QA, set API `FCM_ENABLED=true`, `FCM_WEB_ENABLED=true` and app `EXPO_PUBLIC_WEB_PUSH_ENABLED=true`. Configure the API's existing `FCM_CREDENTIALS_PATH`. Flags default off. Missing public config produces setup guidance; denied/unsupported states retain foreground inbox polling.
7. Verify the platform matrix above with the flag true, false, and missing. Check the configured Google Play destination and missing-link state. Guard-only accounts must not enter guard screens on web, including through direct links.
8. Complete the physical-device checklist below before enabling production push or advertising background delivery. Firebase provisioning/deployment and physical-device push verification remain separate release requirements.

## Behavior and privacy

On allowed browsers, resident web use needs no installation. Apple access detection includes Macs, but Home Screen installation guidance is limited to iPhone/iPad; macOS uses desktop notification setup. A compact **Enable alerts / Set up / dismiss** row appears below the dashboard header with 12px padding, 14px text, wrapping, and at least 44px touch targets. **Set up** opens expanded notification settings in the inbox; settings can also be opened from the inbox itself. Dismissal or successful registration hides the dashboard reminder for the browser tab's session, including navigation/reloads. Signing out clears the dismissal.

On iPhone/iPad, setup explains Safari → Share → Add to Home Screen. Supported iOS/iPadOS 16.4+ Home Screen apps can receive push; desktop and Android browsers do not require installation. Permission is requested only by **Enable notifications**. An already-granted browser re-registers for a fresh authenticated resident session without prompting.

Inbox and unread count poll every 30 seconds only while authenticated, resident-eligible, visible and online, and refetch immediately on return/reconnect. The first response is a quiet baseline. Push and polling share a bounded ID deduplicator. New first-page items merge into loaded pages. Resident API caches are invalidated for related events. No authenticated responses are cached by the worker.

The server binds web subscriptions to its authenticated session version and checks current resident eligibility and event scope at delivery. Password change/reset makes old web subscriptions ineligible. The worker has only a recipient account ID (not access/refresh credentials) in IndexedDB, and displays data-only FCM messages in one background path. Logout disables display, closes displayed notifications, and attempts Firebase token deletion independently of remote API cleanup. In-flight registrations cannot restore an old local notification session. Token rotation/account switching replaces the browser installation's previous server registration. Use one signed-in account per browser installation.

Taps contain only notification ID and recipient ID. After login, the app loads the owned notification and fresh bootstrap, verifies the residence and relevant visitor entry, then uses the existing resident route. Unavailable, resolved approval, expired, wrong-account or inaccessible entries fall back to the inbox. Decisions happen inside the authenticated app. A push already accepted by FCM cannot be recalled; revoked devices stop server access on their next request. Delivery is not guaranteed instant.

## Release checks (physical devices required; not automated claims)

- iPhone **and** iPad: ordinary Safari login; Home Screen installation; allow, deny, dismiss and changed permission; foreground, background and locked-phone delivery; one notification per event; taps cold/warm, expired login and wrong account.
- Resident-only, guard-only, mixed staff/resident accounts; saved guard workspace; direct `/guard/home`, scanner and entry routes; inactive residence and multiple flats.
- Two installed devices: change/reset password; old access/refresh and web delivery excluded; fresh login registers again. Logout offline, delete-token failure, delayed registration and account switch.
- Poll baseline, repeated push ID, >20 inbox items, pagination/scroll, read/read-all, tab hidden, airplane mode, reconnect, revoked membership and app restart. Confirm no toast flood.
- Apple-only mode: iPhone, iPad desktop-mode Safari, and Mac Safari/Chrome access; Android/Windows/Linux download pages at the root and direct resident links, including saved sessions. Confirm auth restoration and notification polling do not start on blocked browsers. Check the loading state during hydration and a custom/missing Google Play URL.
- Unrestricted mode (flag false or missing): desktop and Android browser access without installation; notification setup, permission, foreground polling and reconnect. Android native foreground/background push, cold-start taps and approval actions. Native Android/iOS routing is unchanged in both modes, with no browser detection or installation banner.
- Small screens, landscape, larger text and VoiceOver. HTTPS worker scope/MIME/cache headers, manifest icons and no API response caches in browser storage.

Automated checks: `npm run typecheck`, `npm run lint`, `npm run test:unit`, `npm run web:export`; API `go test ./...` and `go test -tags=integration ./internal/repositories -run 'TestWebPushSessionIsolation|TestPasswordSessionAtomicUpdate'` with Docker. Monitor setup errors, FCM delivery errors and session revocations without logging credentials or full push tokens.

Official requirements: [Apple Web Push](https://developer.apple.com/documentation/usernotifications/sending-web-push-notifications-in-web-apps-and-browsers), [WebKit iOS requirements](https://webkit.org/blog/13878/web-push-for-web-apps-on-ios-and-ipados/), [Firebase setup](https://firebase.google.com/docs/cloud-messaging/web/get-started), [Firebase receive messages](https://firebase.google.com/docs/cloud-messaging/web/receive-messages), [Expo SDK 57 notifications](https://docs.expo.dev/versions/v57.0.0/sdk/notifications/).
