# Notification action verification

Local verification, 2026-09-09:

- API: `go test ./...` passes in `apna-gate-api`.
- App: `npm run test:unit` passes (76 tests), `tsc --noEmit` and `npm run lint` pass.
- Tests cover strict decision/read/invalidate/dismiss ordering, read failure, aliases, concurrent callbacks, malformed callback/replay data, typed stale conflicts, forbidden/network/server failures, bounded authentication replay, copy fallbacks and historical-copy protection.
- API tests cover member exclusion using the joined resident's user ID, staff-only rejection, strict approval dispatch, durable/push ID and wording equality, multiple recipient devices, category mapping, duplicate invite suppression, idempotent decisions, and photo-free visitor notifications.
- Installed Expo Android source reads direct FCM `categoryId` in `NotificationData.kt`, passes it through `RemoteNotificationContent.kt`, and renders category actions in `ExpoNotificationBuilder.kt`.

## Device acceptance still required before release

Android: `adb devices -l` returned no connected device/emulator. A development-build end-to-end acceptance run was **not performed**. Unit tests and native source inspection do not prove OS rendering or background scheduling.

On an Android development build with a registered FCM token:

1. Deliver resident pending and guard informational pushes. Check exact personalized copy, two appropriate buttons, and no photo.
2. Exercise Approve, Decline, Mark as read and Close in foreground, background and cold start. Actions must never navigate or show confirmation toasts; body taps must preserve resident/guard routing.
3. Verify Close leaves the durable row unread. Other successful actions must complete mark-read before cache invalidation and explicit dismissal.
4. Simulate decision success followed by mark-read failure, failed refresh, delayed login inside/outside five minutes, 403, network/5xx, and typed/unknown 409 responses. Check unread and presented state, authorized stale reconciliation, and no expired replay mutation/dismissal.
5. Repeat callbacks and approve from a second device. Confirm the same durable notification ID and no duplicate QR, events or notifications.
6. Join a flat and confirm only the other active residents receive the member-joined notification. Combined approve/check-in must send only the check-in event.

iOS end-to-end push acceptance is **blocked by the existing server transport limitation**: the server skips iOS tokens. APNs `aps.category` mapping is unit-tested, but enabling iOS delivery is outside this change. Once transport exists, verify response-listener and cold-start last-response handling on a physical iOS development build.
