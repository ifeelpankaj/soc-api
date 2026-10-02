# Password change and session revocation

Apply migration `23_user_session_version.sql` before deploying the API. It adds a non-null, zero-default session version without invalidating existing sessions. Deploy the API to **every instance** before shipping the app's all-device sign-out flow. An older API instance does not enforce revocation and must not remain in rotation. Do not roll back to an API without version checks after users have changed passwords; roll forward instead.

Access and refresh tokens carry `session_version`. Legacy tokens without the claim are version zero and work only while the user's database version is zero. All authenticated route guards check the database version. Password change and password reset increment it atomically with the hash update. A conditional update rejects concurrent writes. A refresh rechecks the incoming version against its user snapshot and cannot mint a token for a newer session. Existing requests already authorized before the change may finish; subsequent requests on every device return `401 SESSION_REVOKED`.

## API errors

`POST /api/v1/auth/change-password` keeps the same body (`current_password`, `new_password`, `confirm_password`) and success envelope. Success clears auth cookies; clients must discard their saved tokens and sign in again.

| Status | Code | Meaning |
| --- | --- | --- |
| 400 | CURRENT_PASSWORD_INCORRECT | Current password does not match; keep the user signed in and focus that field. |
| 400 | PASSWORD_NOT_SET | No stored password; offer password reset. |
| 400 | PASSWORD_TOO_LONG | New password exceeds bcrypt's 72 UTF-8 byte limit. |
| 400 | PASSWORD_MISMATCH / PASSWORD_REUSE | Correct the new-password fields. |
| 401 | SESSION_REVOKED | Clear local authentication; do not attempt refresh. |
| 409 | PASSWORD_CHANGE_CONFLICT | Another request changed the password; sign in again. |
| 500 | PASSWORD_VERIFICATION_FAILED | Invalid stored hash; log a sanitized cause and offer reset/support. |

Do not automatically replay the password POST after refresh or transport failure. A timeout may occur after the server committed the change. Offer sign-in/reset without asserting that the password stayed unchanged. During staggered app deployment, the new app also recognizes the old change-password-only `401 INVALID_CREDENTIALS` as a field error.

## Verification before release

- Run `go test ./...` and `go test -tags=integration ./internal/repositories -run TestPasswordSessionAtomicUpdate -count=1 -v` (Docker required). The integration test creates and removes a private PostgreSQL container; it never uses the configured application database.
- Run app `npm run typecheck`, `npm run lint`, and `npm run test:unit`.
- On Android and iOS, verify show/hide, password autofill, keyboard focus, large text, toast placement and announcements, header/hardware/gesture Back, Keep editing/Discard, and disabled navigation during submission.
- With two test devices, change the password on one. Both old access/refresh sessions must fail on their next request; login with the old password must fail and the new password must succeed. Restart the changing device and verify it stays signed out. Offline content on another device cannot be remotely erased.
- Simulate a lost success response, a delayed refresh, and local storage cleanup failure. Never show success for an unconfirmed request or restore protected screens via Back.

Monitor existing request logs by endpoint/status/error code and request ID. Incorrect current passwords should now be 400, not login-style 401. Investigate spikes in `PASSWORD_VERIFICATION_FAILED`, `PASSWORD_CHANGE_CONFLICT`, database lookups, or repeated `SESSION_REVOKED` after fresh login. Never log request passwords, password hashes, or JWTs.
