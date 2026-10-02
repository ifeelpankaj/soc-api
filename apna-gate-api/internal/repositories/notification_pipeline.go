package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// MaterializeOutbox commits the inbox, device work, and producer acknowledgement
// together. Provider sends never run inside this transaction.
func (r *notificationRepository) MaterializeOutbox(ctx context.Context, d contracts.NotificationOutbox, n models.NotificationCreate, push, androidEnabled, webEnabled, iosEnabled bool) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var leased bool
	err = tx.QueryRow(ctx, `SELECT lease_token=$2 AND available_at>now() AND completed_at IS NULL FROM notification_outbox WHERE id=$1 FOR UPDATE`, d.ID, d.LeaseToken).Scan(&leased)
	if err != nil {
		return err
	}
	if !leased {
		return errors.New("notification outbox lease lost")
	}
	if n.Data == nil {
		n.Data = map[string]any{}
	}
	n.Data["notification_id"] = n.ID
	n.Data["event_key"] = d.EventKey
	n.Data["schema_version"] = 1
	domain := notificationDomain(n.Type)
	if domain == "maintenance" {
		n.Data["entity_type"] = "maintenance_bill"
		if v, ok := n.Data["bill_id"]; ok {
			n.Data["entity_id"] = fmt.Sprint(v)
		}
	}
	if domain == "hub" {
		n.Data["entity_type"] = "hub_post"
		if v, ok := n.Data["post_id"]; ok {
			n.Data["entity_id"] = fmt.Sprint(v)
		}
	}
	if domain == "visitor" {
		n.Data["entity_type"] = "visitor_entry"
		if v, ok := n.Data["entry_id"]; ok {
			n.Data["entity_id"] = fmt.Sprint(v)
		}
	}
	data, err := json.Marshal(n.Data)
	if err != nil {
		return err
	}
	var id string
	err = tx.QueryRow(ctx, `WITH inserted AS (
 INSERT INTO notifications(id,user_id,society_id,flat_id,type,title,body,data,event_key,domain)
 VALUES($1::uuid,$2,$3,$4,$5,$6,$7,$8,$9,$10)
 ON CONFLICT(user_id,event_key) DO NOTHING RETURNING id
) SELECT id::text FROM inserted UNION ALL SELECT id::text FROM notifications WHERE user_id=$2 AND event_key=$9 LIMIT 1`,
		n.ID, n.UserID, n.SocietyID, n.FlatID, n.Type, n.Title, n.Body, data, d.EventKey, domain).Scan(&id)
	if err != nil {
		return err
	}
	if push {
		allowed, err := pushPreference(ctx, tx, d.UserID, d.SocietyID, n.Type, n.Data)
		if err != nil {
			return err
		}
		if allowed {
			_, err = tx.Exec(ctx, `INSERT INTO push_deliveries(notification_id,device_token_id,provider,audience)
 SELECT $1::uuid,t.id,CASE WHEN t.platform='ios' THEN 'apns' ELSE 'fcm' END,$3
 FROM device_tokens t JOIN users u ON u.id=t.user_id
 WHERE t.user_id=$2 AND ((t.platform='android' AND $4) OR (t.platform='web' AND $5) OR (t.platform='ios' AND $6))
 AND (t.platform<>'web' OR (t.session_version=u.session_version AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL))
 ON CONFLICT(notification_id,device_token_id) DO NOTHING`, id, d.UserID, d.Audience, androidEnabled, webEnabled, iosEnabled)
			if err != nil {
				return err
			}
		}
	}
	_, err = tx.Exec(ctx, `UPDATE notification_outbox SET inbox_completed_at=coalesce(inbox_completed_at,now()), completed_at=now(),lease_token=NULL,last_error=NULL WHERE id=$1`, d.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func notificationDomain(kind string) string {
	switch {
	case len(kind) >= 7 && kind[:7] == "visitor", len(kind) >= 13 && kind[:13] == "member_invite":
		return "visitor"
	case len(kind) >= 11 && kind[:11] == "maintenance":
		return "maintenance"
	case len(kind) >= 4 && kind[:4] == "hub.":
		return "hub"
	default:
		return "system"
	}
}

func pushPreference(ctx context.Context, tx pgx.Tx, userID, societyID int64, kind string, data map[string]any) (bool, error) {
	if kind == "visitor.pending" {
		return true, nil
	}
	if kind == "hub.announcement" && data["priority"] == "high" {
		return true, nil
	}
	if kind == "hub.reaction" || kind == "hub.removed" {
		return false, nil
	}
	var p contracts.NotificationPreferences
	err := tx.QueryRow(ctx, `SELECT visitor_updates_push,maintenance_push,announcements_push,hub_replies_push
 FROM notification_preferences WHERE user_id=$1 AND society_id=$2`, userID, societyID).Scan(&p.VisitorUpdatesPush, &p.MaintenancePush, &p.AnnouncementsPush, &p.HubRepliesPush)
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	switch notificationDomain(kind) {
	case "visitor":
		return p.VisitorUpdatesPush, nil
	case "maintenance":
		return p.MaintenancePush, nil
	case "hub":
		if kind == "hub.announcement" {
			return p.AnnouncementsPush, nil
		}
		return p.HubRepliesPush, nil
	default:
		return true, nil
	}
}

func (r *notificationRepository) ClaimPushDelivery(ctx context.Context, lease uuid.UUID) (*contracts.PushDelivery, error) {
	var d contracts.PushDelivery
	var raw []byte
	var id string
	var platform string
	err := r.db.Pool.QueryRow(ctx, `WITH claimed AS (
 UPDATE push_deliveries SET status='processing',lease_token=$1,lease_until=now()+interval '60 seconds',attempt_count=attempt_count+1,lifetime_attempt_count=lifetime_attempt_count+1,updated_at=now()
 WHERE id=(SELECT id FROM push_deliveries WHERE
 ((status IN ('pending','retry') AND next_attempt_at<=now()) OR (status='processing' AND lease_until<now()))
 ORDER BY next_attempt_at,id FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING *
 ) SELECT c.id,c.notification_id::text,n.user_id,n.society_id,n.flat_id,c.device_token_id,
 coalesce(t.token,''),coalesce(t.platform::text,''),c.provider,c.audience,n.title,n.body,n.data,c.attempt_count
 FROM claimed c JOIN notifications n ON n.id=c.notification_id LEFT JOIN device_tokens t ON t.id=c.device_token_id`, lease).Scan(
		&d.ID, &id, &d.UserID, &d.SocietyID, &d.FlatID, &d.TokenID, &d.Token, &platform, &d.Provider, &d.Audience, &d.Title, &d.Body, &raw, &d.AttemptCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.NotificationID = id
	d.Platform = models.DevicePlatform(platform)
	d.LeaseToken = lease
	if err = json.Unmarshal(raw, &d.Data); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *notificationRepository) FinishPushDelivery(ctx context.Context, id int64, lease uuid.UUID, messageID string) error {
	tag, err := r.db.Pool.Exec(ctx, `UPDATE push_deliveries SET status='sent',sent_at=now(),provider_message_id=$3,lease_until=NULL,lease_token=NULL,last_error=NULL,last_error_code=NULL,updated_at=now() WHERE id=$1 AND lease_token=$2 AND status='processing' AND lease_until>now()`, id, lease, messageID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("push lease lost")
	}
	return nil
}

func (r *notificationRepository) FailPushDelivery(ctx context.Context, id int64, lease uuid.UUID, code, detail string, permanent bool) error {
	delays := []time.Duration{15 * time.Second, time.Minute, 5 * time.Minute, 30 * time.Minute, 2 * time.Hour}
	var attempts int
	err := r.db.Pool.QueryRow(ctx, `SELECT attempt_count FROM push_deliveries WHERE id=$1 AND lease_token=$2 AND status='processing' AND lease_until>now()`, id, lease).Scan(&attempts)
	if err != nil {
		return err
	}
	status := "retry"
	if permanent || attempts >= 6 {
		status = "dead"
	}
	delay := time.Duration(0)
	if status == "retry" {
		delay = delays[min(max(attempts-1, 0), len(delays)-1)]
	}
	tag, err := r.db.Pool.Exec(ctx, `UPDATE push_deliveries SET status=$3,next_attempt_at=now()+$4::interval,lease_until=NULL,lease_token=NULL,last_error_code=$5,last_error=$6,updated_at=now() WHERE id=$1 AND lease_token=$2 AND status='processing'`, id, lease, status, fmt.Sprintf("%f seconds", delay.Seconds()), code, detail)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("push lease lost")
	}
	return nil
}

func (r *notificationRepository) ReplayPushDelivery(ctx context.Context, id int64) error {
	tag, err := r.db.Pool.Exec(ctx, `UPDATE push_deliveries SET status='pending',attempt_count=0,replay_count=replay_count+1,next_attempt_at=now(),last_error=NULL,last_error_code=NULL,updated_at=now() WHERE id=$1 AND status='dead'`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return contracts.ErrNotFound
	}
	return nil
}

func (r *notificationRepository) GetNotificationPreferences(ctx context.Context, user, society int64) (contracts.NotificationPreferences, error) {
	p := contracts.NotificationPreferences{VisitorUpdatesPush: true, MaintenancePush: true, AnnouncementsPush: true, HubRepliesPush: true}
	var member bool
	if err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM society_members WHERE user_id=$1 AND society_id=$2 AND status='active')`, user, society).Scan(&member); err != nil {
		return p, err
	}
	if !member {
		return p, contracts.ErrNotFound
	}
	err := r.db.Pool.QueryRow(ctx, `SELECT visitor_updates_push,maintenance_push,announcements_push,hub_replies_push FROM notification_preferences WHERE user_id=$1 AND society_id=$2`, user, society).Scan(&p.VisitorUpdatesPush, &p.MaintenancePush, &p.AnnouncementsPush, &p.HubRepliesPush)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, nil
	}
	return p, err
}

func (r *notificationRepository) SetNotificationPreferences(ctx context.Context, user, society int64, p contracts.NotificationPreferences) error {
	tag, err := r.db.Pool.Exec(ctx, `INSERT INTO notification_preferences(user_id,society_id,visitor_updates_push,maintenance_push,announcements_push,hub_replies_push)
 SELECT $1,$2,$3,$4,$5,$6 WHERE EXISTS(SELECT 1 FROM society_members WHERE user_id=$1 AND society_id=$2 AND status='active')
 ON CONFLICT(user_id,society_id) DO UPDATE SET visitor_updates_push=$3,maintenance_push=$4,announcements_push=$5,hub_replies_push=$6,updated_at=now()`, user, society, p.VisitorUpdatesPush, p.MaintenancePush, p.AnnouncementsPush, p.HubRepliesPush)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return contracts.ErrNotFound
	}
	return nil
}

func (r *notificationRepository) ClaimHubFanout(ctx context.Context, lease uuid.UUID) (*contracts.HubFanoutJob, error) {
	var j contracts.HubFanoutJob
	err := r.db.Pool.QueryRow(ctx, `UPDATE hub_notification_fanout_jobs SET lease_token=$1,lease_until=now()+interval '60 seconds'
 WHERE id=(SELECT id FROM hub_notification_fanout_jobs WHERE completed_at IS NULL AND (lease_until IS NULL OR lease_until<now()) ORDER BY id FOR UPDATE SKIP LOCKED LIMIT 1)
 RETURNING id,society_id,post_id,actor_user_id,important`, lease).Scan(&j.ID, &j.SocietyID, &j.PostID, &j.ActorUserID, &j.Important)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	j.LeaseToken = lease
	return &j, nil
}

func (r *notificationRepository) HubContentAvailable(ctx context.Context, society, post int64) (bool, error) {
	var available bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM channel_posts p JOIN societies s ON s.id=p.society_id
	 WHERE p.id=$1 AND p.society_id=$2 AND p.status='active' AND s.status='active' AND s.deleted_at IS NULL)`, post, society).Scan(&available)
	return available, err
}

func (r *notificationRepository) FanoutHubBatch(ctx context.Context, j contracts.HubFanoutJob, limit int) (int, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var cursor int64
	err = tx.QueryRow(ctx, `SELECT last_user_id FROM hub_notification_fanout_jobs WHERE id=$1 AND lease_token=$2 AND lease_until>now() AND completed_at IS NULL FOR UPDATE`, j.ID, j.LeaseToken).Scan(&cursor)
	if err != nil {
		return 0, err
	}
	rows, err := tx.Query(ctx, `SELECT DISTINCT m.user_id FROM society_members m JOIN users u ON u.id=m.user_id
 JOIN societies s ON s.id=m.society_id
 WHERE m.society_id=$1 AND m.role='resident' AND m.status='active' AND m.user_id>$2 AND m.user_id<>$3
	 AND s.status='active' AND s.deleted_at IS NULL
 AND EXISTS(SELECT 1 FROM channel_posts p WHERE p.id=$5 AND p.society_id=$1 AND p.status='active')
 AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL ORDER BY m.user_id LIMIT $4`, j.SocietyID, cursor, j.ActorUserID, limit, j.PostID)
	if err != nil {
		return 0, err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return 0, err
	}
	key := fmt.Sprintf("hub.announcement:%d", j.PostID)
	priority := "normal"
	title := "New society announcement"
	if j.Important {
		priority = "high"
		title = "Important society announcement"
	}
	for _, user := range ids {
		payload := models.NotificationCreate{UserID: user, SocietyID: &j.SocietyID, Type: "hub.announcement", Title: title, Body: "Open Society Hub to view the update.", EventKey: &key, Data: map[string]any{"type": "hub.announcement", "society_id": fmt.Sprint(j.SocietyID), "post_id": fmt.Sprint(j.PostID), "priority": priority, "event_key": key}}
		b, e := json.Marshal(payload)
		if e != nil {
			return 0, e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO notification_outbox(user_id,society_id,audience,event_key,payload,push_enabled)
 VALUES($1,$2,'hub_member',$3,$4,true) ON CONFLICT(user_id,event_key) DO NOTHING`, user, j.SocietyID, key, b); e != nil {
			return 0, e
		}
	}
	if len(ids) > 0 {
		cursor = ids[len(ids)-1]
	}
	_, err = tx.Exec(ctx, `UPDATE hub_notification_fanout_jobs SET last_user_id=$3,lease_token=NULL,lease_until=NULL,
 completed_at=CASE WHEN $4::int<$5::int THEN now() ELSE completed_at END WHERE id=$1 AND lease_token=$2`, j.ID, j.LeaseToken, cursor, len(ids), limit)
	if err != nil {
		return 0, err
	}
	return len(ids), tx.Commit(ctx)
}

func (r *notificationRepository) PipelineStats(ctx context.Context) (map[string]float64, error) {
	stats := make(map[string]float64)
	for key, query := range map[string]string{
		"outbox_pending":            `SELECT count(*)::double precision FROM notification_outbox WHERE completed_at IS NULL`,
		"outbox_oldest_seconds":     `SELECT greatest(coalesce(extract(epoch FROM now()-min(created_at)),0),0) FROM notification_outbox WHERE completed_at IS NULL`,
		"hub_fanout_pending":        `SELECT count(*)::double precision FROM hub_notification_fanout_jobs WHERE completed_at IS NULL`,
		"hub_fanout_oldest_seconds": `SELECT coalesce(extract(epoch FROM now()-min(created_at)),0) FROM hub_notification_fanout_jobs WHERE completed_at IS NULL`,
		"push_pending":              `SELECT count(*)::double precision FROM push_deliveries WHERE status='pending'`,
		"push_retry":                `SELECT count(*)::double precision FROM push_deliveries WHERE status='retry'`,
		"push_dead":                 `SELECT count(*)::double precision FROM push_deliveries WHERE status='dead'`,
	} {
		var value float64
		if err := r.db.Pool.QueryRow(ctx, query).Scan(&value); err != nil {
			return nil, err
		}
		stats[key] = value
	}
	return stats, nil
}

var _ contracts.PipelineRepository = (*notificationRepository)(nil)
