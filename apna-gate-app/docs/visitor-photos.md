# Visitor photos

Guards can choose or capture an optional photo when creating or editing an entry.
Creation saves the entry ID before starting the separate photo upload. Failures
appear in toasts and leave the normal form open, with saved fields locked. Save
retries only the pending photo; choosing another image replaces that selection.
Closing keeps the saved entry. Reopening restores normal editing from server data.
Guards and residents can also preview and save camera/gallery profile photos.

The API accepts JPEG/PNG files up to 5 MiB with a 25-megapixel decoded limit.
Visitor uploads are auto-oriented, proportionally resized to at most 1600px,
flattened onto white and JPEG-encoded at quality 80. The stored file retains the
whole composition. Compliant metadata-free JPEGs produced by the app are stored
byte-for-byte; PNG, oversized, metadata-bearing, or unusual JPEGs use the API's
full normalization fallback. Avatar upload behavior is unchanged.

Authorized photo and avatar GET endpoints accept these server-controlled variants:

| Variant | Transformation |
| --- | --- |
| list | w-80,h-80,fo-face |
| avatar | w-300,h-300,fo-face |
| detail | w-600 |
| original (default) | none |

The API signs the transformed URL with a 15-minute expiry. Responses include
`url`, `expires_at`, and `variant`. The app refreshes three minutes before expiry
and keeps the current bitmap visible while the signed URL changes. Memory caches
are separated by user, society, entry, guard/resident-flat context, stored photo
reference and variant; replacement invalidation follows the stored reference so
unrelated visitor photos do not reload. Query and image memory caches are cleared
on logout.

## Automated checks

- API: `go test ./internal/services/imageSvc ./internal/models ./internal/routes ./internal/handlers/v1`
- Optional PostgreSQL concurrency checks (Docker required): `go test -tags=integration ./internal/services/imageSvc`
- App: `npm run lint`, `npx tsc --noEmit`, `npm run test:unit`
- Configured provider smoke test: set `APNA_GATE_IMAGEKIT_SMOKE=1` and run
  `go test ./internal/services/imageSvc -run TestConfiguredImageKitSmoke -count=1 -v`
  in the API project. This explicitly uses development ImageKit configuration,
  uploads/deletes a generated private fixture, checks all signed variants, and
  verifies unsigned/tampered requests are rejected. It does not modify records.

## Device and staging acceptance

These checks require physical Android/iOS devices and a configured ImageKit staging account:

1. Create entries without photos, with gallery JPEG/PNG/HEIC, and with the camera.
   Verify permission denial, settings recovery, cancellation, and oversized images.
2. Interrupt upload after entry creation. Press Save again and confirm the entry
   count remains unchanged. Close and confirm the entry is retained.
3. Add a photo with no other edits; replace an existing photo; interrupt replacement
   after editing fields. The old photo remains until replacement succeeds.
4. Check guard/resident queues, lists, logs and detail screens. Photos replace the
   visitor placeholder. Lists use square face crops; details keep the whole photo.
5. Inspect single-face, group and no-face thumbnails. ImageKit chooses the crop;
   uploads are not rejected based on face detection.
6. Replace a visitor photo from Entry B and reopen Entry A for the same visitor.
   Seeing the new photo is expected: photos belong to visitors, not individual visits.
7. Verify resident flat isolation, logout/login and foreground refresh. Check that
   unsigned, expired and transformation-tampered URLs cannot expose private photos.

## Release order

Deploy the API with ImageKit configured first, then rebuild the native app with
the added Expo modules and updated permission strings. The visitor-photo change
adds no database migration; it uses the existing managed-image schema. Retain the
existing image failure/duration/cleanup metrics when monitoring rollout.
