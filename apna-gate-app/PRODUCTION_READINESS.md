# Mobile redesign and release readiness

Reviewed 14 September 2026. The implementation is complete; production release approval remains blocked by the items below. No app-store submission or backend deployment was performed.

## Implemented

| Screen family | Changes and preserved behavior |
| --- | --- |
| Resident approvals | Shared visitor queue cards, complete-response search before local pagination, exact purpose filters/counts, compact invite action, permission-aware decisions, distinct empty/error states, and duplicate-tap protection. |
| Visitor invitations | Descriptive purpose selection, residence context, fixed Create invite action, Guest default, Service/Staff exclusion, preserved share links and navigation. |
| Visitor details and dashboard modal | Fixed resident actions, loading on the submitted decision, both controls disabled while submitting, modal closes only after success. Existing guard status/action mappings retained. |
| Members and invite history | Fixed member-invite submission, one Add Member action instead of a duplicate menu, consistent card styling, no duplicate Copy action in invite details. |
| Guard queues, Create/Edit and Check In | Shared queue presentation and bottom-action implementation; existing permissions and operational actions retained. Manual scanner recovery cannot submit an empty code. |
| Dashboards and history | Warm backgrounds, navy emphasis, quieter shadows, consistent cards and selected filters, wrapping action labels. No route or API changes. |
| Profile and settings | Shared Save profile and Change password footers, quieter secondary controls; existing setting save/reset semantics retained. |
| Authentication and society selection | Shared palette, guarded login/reset submissions, login remains busy through account bootstrap, recoverable session-restoration errors. |
| Notifications | Full title/body text, accessible mark-read target, existing role routing and decision controls retained. |

Shared footers own bottom safe-area spacing; content remains in normal flex layout. Shared sub-screens and footers constrain reading width on tablets. Forms retain one keyboard-aware container. Appearance is explicitly light; portrait orientation is retained.

## Reliability and configuration

- API requests time out after 30 seconds. Network and timeout errors have actionable messages.
- Concurrent automatic refresh still shares one request. Transient refresh/profile/bootstrap failures preserve saved credentials and offer retry; confirmed authentication failures clear the session.
- Automatic recovery replays GET/HEAD only. A mutation rejected for authentication must be retried by the user after recovery.
- Session invalidation clears API and workspace state. Changing residence remounts resident screen state so old lists, edits and invitations are not reused in the new residence.
- Push registration and photo diagnostics no longer log raw exception messages or local file paths. The existing error boundary logs details only in development.
- Expo SDK 57 packages were aligned; native packages changed and require a fresh native build. Required image/browser config plugins are declared.
- Metro, metro-config and metro-transform-worker are pinned to 0.84.5, the patched version used by the installed Expo toolchain, to avoid vulnerable older nested copies. Revisit these overrides with the next SDK update.
- Android SecureStore backup exclusions were inspected in the dependency's Android resources; both cloud backup and device transfer exclude SecureStore. Camera/photo/microphone plugin settings and notification permissions were reviewed.
- `npm run release:check` validates public HTTPS API/web hostnames, appearance, and app identifiers. EAS preview/production builds run it after installation. Development profiles retain local Docker/LAN access.

## Verification

- TypeScript and ESLint pass.
- 96 unit tests pass, including queue filtering beyond the first page, exact historical purposes, invite-purpose presentation, authentication retry policy, notification routing/actions, photo lifecycle, navigation, and guard footer state transitions.
- Android Hermes, iOS Hermes, and web/static production-mode bundles export successfully. This validates bundling, not native signing or device behavior.
- Installed dependencies match the SDK 57 compatibility metadata checked locally. The initial online check identified the updates applied in this change.
- `git diff --check` passes with Windows CRLF handling.
- Release configuration checks reject the current local HTTP API and web URLs, as intended. Actual production service availability and EAS environment values were not verified.

## Release blockers and device acceptance

1. **Production endpoints:** configure real public HTTPS API and web URLs in the production EAS environment, then run `npm run release:check` with those values. Local development environment files were not changed.
2. **Upstream moderate advisory:** npm reports three moderate findings in one chain: `expo-router -> query-string -> decode-uri-component`. The decoder advisory is [GHSA-vcc3-ghjq-m6fr](https://github.com/advisories/GHSA-vcc3-ghjq-m6fr). npm's proposed automatic fix downgrades Expo Router to 5.1.11, incompatible with this SDK. The patched decoder 0.5.0 is ESM whereas the installed consumer expects the older CommonJS API; a blind override was not applied. Resolve through a compatible upstream release or a separately verified patch before release acceptance. No high or critical advisories remained after the dependency fixes.
3. **Native/device acceptance:** no Android device was connected and the browser tool reported no available browser. iOS native builds cannot be run in this Windows environment. Final visual, keyboard, large-font, tablet, permission-dialog, push-delivery and signed-build checks remain unverified; they are not counted as passed.

Device acceptance should cover resident approve/reject success and failure, invite creation/sharing, guard scan/manual recovery and check-in/out, Create/Edit with keyboard open, denied camera/photos/notifications, offline launch/retry/logout, interrupted uploads, session expiry, and residence switching. Confirm footer reachability, no overlap or double keyboard movement, readable long names, and correct Android/iPhone safe areas. Do not use real visitor decisions as test fixtures.

The new light appearance, splash/notification branding and SDK updates require rebuilding the native app; Metro reload alone cannot apply those native changes.
