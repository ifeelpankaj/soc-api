-- name: HubProvisionChannels :exec
INSERT INTO society_channels(society_id,type,name,description)
VALUES ($1,'announcement','Announcements','Official society updates'),($1,'community','Community','Connect with neighbours')
ON CONFLICT(society_id,type) DO NOTHING;

-- name: HubRole :one
SELECT m.role::text FROM society_members m JOIN societies s ON s.id=m.society_id JOIN users u ON u.id=m.user_id
WHERE m.society_id=$1 AND m.user_id=$2 AND m.status='active' AND m.role IN ('owner','admin','resident')
AND s.status='active' AND s.deleted_at IS NULL AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL
ORDER BY CASE m.role WHEN 'owner' THEN 0 WHEN 'admin' THEN 1 ELSE 2 END LIMIT 1;

-- name: HubActorLock :exec
SELECT pg_advisory_xact_lock(hashtextextended('hub/' || sqlc.arg(society_id)::bigint::text || '/' || sqlc.arg(user_id)::bigint::text,0));

-- name: HubChannels :many
SELECT c.*, (SELECT count(*) FROM channel_posts p WHERE p.channel_id=c.id AND p.status='active'
 AND (r.last_read_post_id IS NULL OR (p.created_at,p.id)>(r.last_read_created_at,r.last_read_post_id)))::bigint AS unread_count,
 COALESCE((SELECT p.id FROM channel_posts p WHERE p.channel_id=c.id AND p.status='active' AND p.is_important ORDER BY p.created_at DESC,p.id DESC LIMIT 1),0)::bigint AS important_post_id
FROM society_channels c LEFT JOIN channel_read_states r ON r.channel_id=c.id AND r.user_id=sqlc.arg(user_id)
WHERE c.society_id=sqlc.arg(society_id) ORDER BY c.type;

-- name: HubChannel :one
SELECT * FROM society_channels WHERE society_id=$1 AND id=$2;

-- name: HubLockChannel :one
SELECT * FROM society_channels WHERE society_id=$1 AND id=$2 FOR UPDATE;

-- name: HubCategories :many
SELECT * FROM community_categories WHERE is_active ORDER BY display_order,id;

-- name: HubCategory :one
SELECT id FROM community_categories WHERE is_active AND (id=sqlc.narg(id)::smallint OR (sqlc.narg(id)::smallint IS NULL AND code='general'));

-- name: HubPost :one
SELECT p.*,c.type AS channel_type FROM channel_posts p JOIN society_channels c ON c.id=p.channel_id WHERE p.society_id=$1 AND p.id=$2;

-- name: HubLockPost :one
SELECT p.*,c.type AS channel_type FROM channel_posts p JOIN society_channels c ON c.id=p.channel_id WHERE p.society_id=$1 AND p.id=$2 FOR UPDATE OF p;

-- name: HubPostIDs :many
SELECT id FROM channel_posts WHERE society_id=sqlc.arg(society_id) AND channel_id=sqlc.arg(channel_id) AND status='active'
AND (sqlc.narg(category_id)::smallint IS NULL OR category_id=sqlc.narg(category_id))
AND (sqlc.narg(before_at)::timestamptz IS NULL OR (created_at,id)<(sqlc.narg(before_at),sqlc.arg(before_id)::bigint))
ORDER BY created_at DESC,id DESC LIMIT sqlc.arg(page_limit);

-- name: HubPinnedIDs :many
SELECT id FROM channel_posts WHERE society_id=$1 AND channel_id=$2 AND status='active' AND is_pinned ORDER BY created_at DESC,id DESC;

-- name: HubCreatePost :one
INSERT INTO channel_posts(society_id,channel_id,author_id,title,body,category_id,comments_enabled,is_pinned,is_important)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;

-- name: HubUpdatePost :exec
UPDATE channel_posts SET title=$3,body=$4,category_id=$5,comments_enabled=$6,is_pinned=$7,is_important=$8,edited_at=clock_timestamp(),updated_at=clock_timestamp()
WHERE society_id=$1 AND id=$2;

-- name: HubControlPost :exec
UPDATE channel_posts SET is_pinned=$3,comments_enabled=$4,updated_at=clock_timestamp() WHERE society_id=$1 AND id=$2;

-- name: HubDeletePost :exec
UPDATE channel_posts SET status=$3,deleted_by=$4,deleted_at=clock_timestamp(),removal_reason=$5,updated_at=clock_timestamp() WHERE society_id=$1 AND id=$2;

-- name: HubPostRate :one
SELECT count(*)::bigint AS count,min(created_at)::timestamptz AS oldest FROM channel_posts
WHERE society_id=$1 AND author_id=$2 AND created_at>clock_timestamp()-sqlc.arg(window_seconds)::int*interval '1 second';

-- name: HubCommentRate :one
SELECT count(*)::bigint AS count,min(c.created_at)::timestamptz AS oldest FROM channel_comments c JOIN channel_posts p ON p.id=c.post_id
WHERE p.society_id=$1 AND c.author_id=$2 AND c.created_at>clock_timestamp()-sqlc.arg(window_seconds)::int*interval '1 second';

-- name: HubComment :one
SELECT c.* FROM channel_comments c JOIN channel_posts p ON p.id=c.post_id WHERE p.society_id=$1 AND c.id=$2;

-- name: HubCommentIDs :many
SELECT c.id FROM channel_comments c JOIN channel_posts p ON p.id=c.post_id
WHERE p.society_id=sqlc.arg(society_id) AND p.id=sqlc.arg(post_id)
AND (c.status='active' OR EXISTS(SELECT 1 FROM channel_comments child WHERE child.parent_id=c.id AND child.status='active'))
AND (sqlc.narg(after_at)::timestamptz IS NULL OR (c.created_at,c.id)>(sqlc.narg(after_at),sqlc.arg(after_id)::bigint))
ORDER BY c.created_at,c.id LIMIT sqlc.arg(page_limit);

-- name: HubCreateComment :one
INSERT INTO channel_comments(post_id,author_id,parent_id,body) VALUES($1,$2,$3,$4) RETURNING id;

-- name: HubUpdateComment :exec
UPDATE channel_comments SET body=$3,edited_at=clock_timestamp(),updated_at=clock_timestamp() WHERE id=$1 AND post_id=$2;

-- name: HubDeleteComment :exec
UPDATE channel_comments SET status=$3,deleted_by=$4,deleted_at=clock_timestamp(),removal_reason=$5,updated_at=clock_timestamp() WHERE id=$1 AND post_id=$2;

-- name: HubPostView :one
SELECT jsonb_build_object('id',p.id,'society_id',p.society_id,'channel_id',p.channel_id,'author',jsonb_build_object('id',u.id,'name',u.full_name),
 'title',p.title,'body',p.body,'category_id',p.category_id,'comments_enabled',p.comments_enabled,'is_pinned',p.is_pinned,'is_important',p.is_important,
 'status',p.status,'created_at',p.created_at,'updated_at',p.updated_at,'edited_at',p.edited_at,
 'reply_count',(SELECT count(*) FROM channel_comments c WHERE c.post_id=p.id AND c.status='active'),
 'reactions',COALESCE((SELECT jsonb_object_agg(x.reaction_type,x.n) FROM (SELECT reaction_type,count(*) n FROM channel_reactions WHERE post_id=p.id GROUP BY reaction_type) x),'{}'),
 'my_reactions',COALESCE((SELECT jsonb_agg(reaction_type ORDER BY reaction_type) FROM channel_reactions cr WHERE cr.post_id=p.id AND cr.user_id=sqlc.arg(user_id)),'[]'),
 'attachment_ids',COALESCE((SELECT jsonb_agg(upload_id ORDER BY id) FROM channel_attachments WHERE post_id=p.id),'[]')) AS data
FROM channel_posts p JOIN users u ON u.id=p.author_id WHERE p.society_id=sqlc.arg(society_id) AND p.id=sqlc.arg(id) AND p.status='active';

-- name: HubCommentView :one
SELECT jsonb_build_object('id',c.id,'post_id',c.post_id,'parent_id',c.parent_id,'status',c.status,
 'author',CASE WHEN c.status='active' THEN jsonb_build_object('id',u.id,'name',u.full_name) ELSE NULL END,
 'body',CASE WHEN c.status='active' THEN c.body ELSE '' END,'created_at',c.created_at,'updated_at',c.updated_at,'edited_at',c.edited_at,
 'reply_count',(SELECT count(*) FROM channel_comments child WHERE child.parent_id=c.id AND child.status='active'),
 'reactions',CASE WHEN c.status='active' THEN COALESCE((SELECT jsonb_object_agg(x.reaction_type,x.n) FROM (SELECT reaction_type,count(*) n FROM channel_reactions WHERE comment_id=c.id GROUP BY reaction_type) x),'{}') ELSE '{}'::jsonb END,
 'my_reactions',CASE WHEN c.status='active' THEN COALESCE((SELECT jsonb_agg(reaction_type ORDER BY reaction_type) FROM channel_reactions cr WHERE cr.comment_id=c.id AND cr.user_id=sqlc.arg(user_id)),'[]') ELSE '[]'::jsonb END,
 'attachment_ids',CASE WHEN c.status='active' THEN COALESCE((SELECT jsonb_agg(upload_id ORDER BY id) FROM channel_attachments WHERE comment_id=c.id),'[]') ELSE '[]'::jsonb END) AS data
FROM channel_comments c JOIN channel_posts p ON p.id=c.post_id JOIN users u ON u.id=c.author_id
WHERE p.society_id=sqlc.arg(society_id) AND c.id=sqlc.arg(id) AND p.status='active';

-- name: HubAddReaction :one
INSERT INTO channel_reactions(user_id,post_id,comment_id,reaction_type) VALUES($1,$2,$3,$4) ON CONFLICT DO NOTHING RETURNING id;

-- name: HubRemoveReaction :exec
DELETE FROM channel_reactions WHERE user_id=sqlc.arg(user_id) AND post_id IS NOT DISTINCT FROM sqlc.narg(post_id)::bigint
AND comment_id IS NOT DISTINCT FROM sqlc.narg(comment_id)::bigint AND reaction_type=sqlc.arg(reaction_type);

-- name: HubRead :exec
INSERT INTO channel_read_states(user_id,channel_id,last_read_post_id,last_read_created_at,last_read_at)
SELECT sqlc.arg(user_id),p.channel_id,p.id,p.created_at,clock_timestamp() FROM channel_posts p
WHERE p.society_id=sqlc.arg(society_id) AND p.channel_id=sqlc.arg(channel_id) AND p.id=sqlc.arg(post_id)
ON CONFLICT(user_id,channel_id) DO UPDATE SET last_read_post_id=excluded.last_read_post_id,last_read_created_at=excluded.last_read_created_at,last_read_at=excluded.last_read_at,updated_at=clock_timestamp()
WHERE (excluded.last_read_created_at,excluded.last_read_post_id)>(channel_read_states.last_read_created_at,channel_read_states.last_read_post_id);

-- name: HubResidents :many
SELECT DISTINCT m.user_id FROM society_members m JOIN users u ON u.id=m.user_id
WHERE m.society_id=$1 AND m.role='resident' AND m.status='active' AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL;

-- name: HubEnqueue :exec
INSERT INTO notification_outbox(user_id,society_id,audience,event_key,payload,push_enabled)
VALUES($1,$2,'hub_member',$3,$4,$5) ON CONFLICT(user_id,event_key) DO NOTHING;

-- name: HubCreateUpload :one
INSERT INTO channel_uploads(society_id,uploader_id,file_id,file_path,original_filename,mime_type,file_size,width,height)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING *;

-- name: HubUpload :one
SELECT * FROM channel_uploads WHERE society_id=$1 AND id=$2;

-- name: HubUploadByFile :one
SELECT * FROM channel_uploads WHERE file_id=$1;

-- name: HubLockUpload :one
SELECT * FROM channel_uploads WHERE society_id=$1 AND id=$2 FOR UPDATE;

-- name: HubClaimUpload :execrows
UPDATE channel_uploads SET status='claimed' WHERE id=$1 AND society_id=$2 AND uploader_id=$3 AND status='pending' AND expires_at>clock_timestamp();

-- name: HubAttach :exec
INSERT INTO channel_attachments(upload_id,post_id,comment_id) VALUES($1,$2,$3);

-- name: HubAttachments :many
SELECT upload_id FROM channel_attachments WHERE post_id IS NOT DISTINCT FROM sqlc.narg(post_id)::bigint AND comment_id IS NOT DISTINCT FROM sqlc.narg(comment_id)::bigint ORDER BY upload_id;

-- name: HubDetach :exec
DELETE FROM channel_attachments WHERE upload_id=$1;

-- name: HubQueueUploadCleanup :exec
UPDATE channel_uploads SET status='cleanup',expires_at=clock_timestamp() WHERE id=$1 AND status IN ('pending','claimed');

-- name: HubAttachmentTarget :one
SELECT a.post_id,a.comment_id,COALESCE(p.society_id,cp.society_id)::bigint AS society_id,
COALESCE(p.status,cp.status)::text AS post_status,c.status AS comment_status
FROM channel_attachments a LEFT JOIN channel_posts p ON p.id=a.post_id LEFT JOIN channel_comments c ON c.id=a.comment_id LEFT JOIN channel_posts cp ON cp.id=c.post_id
WHERE a.upload_id=$1;

-- name: HubClaimCleanup :one
UPDATE channel_uploads SET status='cleanup',lease_until=clock_timestamp()+interval '2 minutes',lease_token=sqlc.arg(lease_token)
WHERE id=(SELECT id FROM channel_uploads WHERE status IN ('pending','cleanup') AND expires_at<=clock_timestamp()
 AND (lease_until IS NULL OR lease_until<=clock_timestamp()) ORDER BY id FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING *;

-- name: HubFinishCleanup :exec
UPDATE channel_uploads SET status='deleted',lease_until=NULL,lease_token=NULL WHERE id=$1 AND lease_token=$2 AND status='cleanup';

-- name: HubReport :one
INSERT INTO channel_reports(post_id,reporter_id,reason) VALUES($1,$2,$3)
ON CONFLICT(post_id,reporter_id) DO UPDATE SET reason=channel_reports.reason RETURNING *;

-- name: HubReports :many
SELECT r.* FROM channel_reports r JOIN channel_posts p ON p.id=r.post_id WHERE p.society_id=$1
AND (sqlc.arg(status)::text='' OR r.status=sqlc.arg(status)) AND r.id>sqlc.arg(after_id)
ORDER BY r.id LIMIT sqlc.arg(page_limit);

-- name: HubReportByID :one
SELECT r.* FROM channel_reports r JOIN channel_posts p ON p.id=r.post_id WHERE p.society_id=$1 AND r.id=$2;

-- name: HubResolveReport :exec
UPDATE channel_reports SET status=$2,reviewer_id=$3,reviewed_at=clock_timestamp(),resolution_note=$4 WHERE id=$1 AND status='pending';

-- name: HubResolvePostReports :exec
UPDATE channel_reports SET status='actioned',reviewer_id=$2,reviewed_at=clock_timestamp(),resolution_note=$3 WHERE post_id=$1 AND status='pending';
