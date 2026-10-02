# Society Hub API

Apply migrations **33** and **34** before deploying this API. Migration 33 was an
unapplied draft and has been replaced with a sql-migrate Up/Down migration.
Existing societies receive Announcements and Community. Society creation now
provisions both channels in the same transaction. There is no channel CRUD API.
Rollback removes Hub data; it is intended for disposable testing or a deliberate
rollback before users publish content.

All routes use `/v1/societies/{societyId}` beneath the deployment's API base path.
Swagger documents all 24 operations, their request bodies, and response schemas.
Every request requires a valid access token, an operational subscription, and
active owner/admin/resident membership. Staff are excluded. Developer/super-admin
roles do not bypass either membership or subscription checks.

## Compose and publish

1. `GET /channels` returns the two channel IDs, unread counts, and each channel's
   latest important announcement, if present. `GET /channels/categories` lists
   the seeded active Community categories.
2. Optionally upload each file using `POST /channel-uploads` with multipart field
   `file`. JPEG/PNG is limited to 5 MiB and 25 million pixels; PDF to 10 MiB.
   The server validates image decoding or PDF MIME/header/EOF signatures.
3. Publish with `POST /channels/{channelId}/posts`:

   ```json
   {"title":"Water maintenance","body":"Water will be unavailable tomorrow from 10 AM to noon.","comments_enabled":true,"is_pinned":true,"is_important":true,"attachment_ids":[123]}
   ```

   Announcement publishing requires owner/admin membership. Community requests
   must omit `comments_enabled`, `is_pinned`, and `is_important`; `category_id`
   defaults to General. Bodies are plain text, not HTML. Clients must render
   them as text. Titles are optional (200 characters), post bodies required
   (20,000 characters), and comment bodies required (5,000 characters).
4. Use `GET /channel-attachments/{uploadId}` to obtain a signed URL. Attachment
   IDs in content responses are upload IDs. The URL expires in 15 minutes; its
   `expires_at` is the URL expiry, while upload responses report claim expiry.
   Pending files are visible only to their uploader. Published files are visible
   to authorized members while their content is active. Provider file IDs and
   canonical paths are never accepted from clients or exposed in these responses.

`PATCH /posts/{postId}` and `PATCH /comments/{commentId}` preserve omitted fields.
An explicit empty title clears it. Omitted/null `attachment_ids` preserves the
current set; an array replaces it; `[]` removes all. At most five files may be
attached. Only the uploader can claim pending files; claimed files cannot be
reused on another target. File claims, content, and notifications commit together.
Upload success does not publish content. Client retries after ambiguous publish
responses should refresh the feed before creating another post.

## Discussions, feeds, and moderation

- `GET /channels/{channelId}/posts?limit=20&cursor=...&category_id=...` returns
  newest-first chronological `items` and a separate `pinned` collection.
  Pinned items also remain in chronology; clients can deduplicate presentation by ID.
- `GET /posts/{postId}/comments` returns oldest-first comments and one-level
  replies. Use `parent_id` when creating a reply to a top-level comment. Deeper
  nesting and moving an existing reply to another parent are rejected.
- Cursors are opaque and scoped to the resource/filter. The maximum page size
  is 100. Responses include author summaries, reply counts, edited timestamps,
  reaction totals, current-user reactions, and attachment IDs.
- `POST /channels/{channelId}/read` takes `{"post_id":123}`. It advances a
  timestamp/ID watermark monotonically. Unread counts count active posts, not
  comment activity, and do not reset on edits or pin changes.
- `POST /posts/{postId}/pin` takes `{"is_pinned":true}`. Only announcements can
  be pinned. `POST /posts/{postId}/lock` takes `{"locked":true}`. Both controls
  require owner/admin membership. Locks block new comments and comment edits,
  including those by admins; unlock before replying.
- Add reactions with `POST /posts/{postId}/reactions` or
  `/comments/{commentId}/reactions`, body `{"reaction_type":"helpful"}`.
  Supported types are `like`, `love`, and `helpful`. Delete via
  `/reactions/{reactionType}`. Both operations are idempotent.
- Authors may edit/delete their own content. Admins may remove others' content
  using DELETE with `{"reason":"..."}`; they cannot rewrite its body.
  Soft deletion preserves audit metadata. Deleted parent comments become empty
  placeholders while active replies remain.
- Report via `POST /posts/{postId}/reports` with `{"reason":"..."}`. Repeated
  reports by the same member return their existing report. Admins list reports
  through `GET /channel-reports?status=pending` and resolve using
  `PATCH /channel-reports/{reportId}` with
  `{"status":"dismissed","note":"Reviewed"}` or `status:"actioned"`.
  Actioning removes the post and resolves its pending reports atomically.

## Delivery, cleanup, and configuration

Announcements notify active resident-role members, excluding the publisher.
Replies notify the post author and direct parent-comment author, deduplicated
and excluding the actor. Post reactions and administrative removals are inbox
only. Community publication and comment reactions do not create notifications.
Important announcements carry `data.priority="high"` for UI treatment.

The existing notification outbox worker handles Hub events using `hub_member`
audience and `push_enabled`. It rechecks access and subscription before delivery.
Inbox entries are deduplicated by recipient/event key; external push is
at-least-once and may repeat after an interrupted send. Existing producers retain
push-enabled behavior through the database default. Operational modules do not
automatically create Hub posts.

The existing cleanup job processes expired uploads and removed attachments with
two-minute leases and retryable provider deletion. Pending uploads expire after
24 hours. Soft-deleted posts/comments retain their attached files for audit.
Deleting a pending upload queues cleanup; it does not synchronously delete it.
Monitor `apna_gate_hub_upload_cleanup_failures_total` and existing notification
outbox backlog/failure metrics. Provider/DB uncertainty is logged with file IDs
for reconciliation; content URLs and provider secrets are not logged.

Optional environment configuration (positive integers, maximum 86400):

| Variable | Default |
|---|---:|
| `HUB_POST_LIMIT` | 5 |
| `HUB_POST_WINDOW_SECONDS` | 600 |
| `HUB_COMMENT_LIMIT` | 30 |
| `HUB_COMMENT_WINDOW_SECONDS` | 60 |

Limits are per user/society and serialized in PostgreSQL; deleted content still
counts. A rejection returns `429 HUB_RATE_LIMIT` and `Retry-After` seconds.
ImageKit uses the existing configuration. If disabled, text features remain
available and upload/URL requests return service unavailable.

## Verification

Run `sqlc generate`, `sqlc vet`, the existing Swagger generation command,
`go test ./...`, and `go test -tags=integration ./...`. Integration tests use
disposable PostgreSQL containers. Run race tests for Hub, notifications, and jobs
on a Go installation with a supported C compiler.
