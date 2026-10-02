# Notification pipeline

Visitor and Hub write `notification_outbox` in their business transaction. Maintenance writes the legacy-named `maintenance_notification_deliveries` domain event staging table; its trigger bridges to the same shared outbox. Do not use that Maintenance table as a push queue.

`cmd/worker` runs independent materializer, Hub fanout, and push delivery loops. Materialization atomically creates the user inbox row, snapshots eligible device work, and acknowledges the outbox row. Maintenance's completion bridge now acknowledges at this step. Provider success and failure live only in `push_deliveries`.

The worker uses FCM for Android and web and direct APNs for the native iOS tokens registered by Expo. Set `APNS_ENABLED`, `APNS_KEY_PATH`, `APNS_KEY_ID`, `APNS_TEAM_ID`, `APNS_TOPIC`, and `APNS_PRODUCTION` for the environment serving iOS builds. Mount the `.p8` signing key outside the image. Do not send APNs tokens to FCM.

## Rollout

1. Apply migration 36 and deploy the API image containing both `server` and `worker` binaries.
2. Start `apna-gate-worker` and confirm `/health/ready` on port 9091 plus the three loop heartbeats in Prometheus.
3. Check old pending `notification_outbox` rows reach `completed_at`, then verify new inbox rows and `push_deliveries` separately. Existing completed rows are not resent.
4. Keep the Maintenance producer and bridge. Billing and reminder jobs no longer drain the outbox; the worker owns delivery.
5. Configure APNs credentials and validate a physical iOS device in the matching sandbox or production environment.

The `notifications` inbox remains available when a user mutes push or has no device token. Essential visitor approval and important announcements bypass mutable push preferences. Unknown app payload versions open the inbox.

Authenticated clients use `GET /api/v1/me/notifications/preferences?society_id=<id>` and `PUT` to the same URL with all four Boolean fields: `visitor_updates_push`, `maintenance_push`, `announcements_push`, and `hub_replies_push`. Preferences require active membership in that society. Essential categories have no mutable field.

## Operations

`push_deliveries` retries transient failures at 15 seconds, 1 minute, 5 minutes, 30 minutes, and 2 hours, then enters `dead` after the sixth failed attempt. A provider may accept a push immediately before a worker crash; that push can appear twice, while the inbox stays unique by `(user_id,event_key)`.

Inspect a dead delivery and its `last_error_code` before replay. Requeue one delivery with `./worker -replay-push-id <id>` inside the worker container. Replay preserves the row, `lifetime_attempt_count`, and increments `replay_count`; `attempt_count` starts a new retry cycle.

Prometheus metrics include `apna_gate_notification_worker_heartbeat_seconds` by loop, `apna_gate_notification_pipeline_backlog` by state, `apna_gate_push_delivery_latency_seconds`, and `apna_gate_push_provider_failures_total`. Alert separately on stale loop heartbeats, growing event/fanout age, and nonzero dead deliveries.
