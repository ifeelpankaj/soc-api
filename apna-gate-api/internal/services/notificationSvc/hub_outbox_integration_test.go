//go:build integration

package notificationsvc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/db"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/repositories/contracts"
	"go-server/internal/requestctx"
	"go-server/internal/testutil"
	"go-server/pkg/database"
)

type hubDeliveryGuard struct{ deny bool }

func (g *hubDeliveryGuard) EnsureSocietyOperational(ctx context.Context, _ int64) error {
	if requestctx.HasDeveloperGuardBypass(ctx) {
		return errors.New("bypass leaked")
	}
	if g.deny {
		return models.NewAppError("INACTIVE", "inactive", 403, nil)
	}
	return nil
}

func TestHubOutboxDelivery(t *testing.T) {
	ctx := context.Background()
	pg := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	sqlDB, err := sql.Open("pgx", pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if _, err = migrate.Exec(sqlDB, "postgres", &migrate.FileMigrationSource{Dir: "../../../migrations"}, migrate.Up); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	id := func(query string, args ...any) int64 {
		t.Helper()
		var n int64
		if e := pool.QueryRow(ctx, query, args...).Scan(&n); e != nil {
			t.Fatal(e)
		}
		return n
	}
	exec := func(query string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, query, args...); e != nil {
			t.Fatal(e)
		}
	}
	owner := id("INSERT INTO users(full_name,email) VALUES('Owner','owner@example.test') RETURNING id")
	resident := id("INSERT INTO users(full_name,email) VALUES('Resident','resident@example.test') RETURNING id")
	society := id("INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Hub','HUB','active',$1,$1,now()) RETURNING id", owner)
	exec("INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'owner','active'),($1,$3,'resident','active')", society, owner, resident)
	channel := id("INSERT INTO society_channels(society_id,type,name) VALUES($1,'community','Community') RETURNING id", society)
	post := id("INSERT INTO channel_posts(society_id,channel_id,author_id,body) VALUES($1,$2,$3,'Post') RETURNING id", society, channel, owner)
	exec("INSERT INTO device_tokens(user_id,token,platform) VALUES($1,'resident-android','android')", resident)
	dbWrap := &database.Database{Pool: pool}
	store := repository.NewNotificationRepository(dbWrap)
	pipeline := store.(contracts.PipelineRepository)
	fcm := &fcmClientFake{}
	guard := &hubDeliveryGuard{}
	s := &notificationService{notifications: store, deviceTokens: &deviceTokenRepoFake{}, fcmClient: fcm, enabled: true, operational: guard}
	q := db.New(pool)
	enqueue := func(key string, push bool) {
		t.Helper()
		n := models.NotificationCreate{UserID: resident, SocietyID: &society, Type: "hub.reply", Title: "Hub reply", Body: "Open Hub", EventKey: &key, Data: map[string]any{"post_id": post, "society_id": society, "type": "hub.reply"}}
		b, e := json.Marshal(n)
		if e != nil {
			t.Fatal(e)
		}
		if e = q.HubEnqueue(ctx, db.HubEnqueueParams{UserID: resident, SocietyID: society, EventKey: key, Payload: b, PushEnabled: push}); e != nil {
			t.Fatal(e)
		}
	}
	enqueue("inbox", false)
	enqueue("inbox", false)
	if err = s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM notifications") != 1 || id("SELECT count(*) FROM push_deliveries") != 0 {
		t.Fatal("inbox-only event created push work or duplicate")
	}
	enqueue("push", true)
	if err = s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if fcm.sendCount != 0 || id("SELECT count(*) FROM push_deliveries WHERE status='pending'") != 1 {
		t.Fatal("provider called during materialization")
	}
	if err = s.runPushBatch(ctx, pipeline); err != nil {
		t.Fatal(err)
	}
	if fcm.sendCount != 1 || id("SELECT count(*) FROM push_deliveries WHERE status='sent'") != 1 {
		t.Fatal("device delivery failed")
	}
	enqueue("retry", true)
	exec("INSERT INTO device_tokens(user_id,token,platform) VALUES($1,'resident-android-2','android')", resident)
	if err = s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	fcm.errorsByToken = map[string]error{"resident-android-2": errors.New("temporary provider failure")}
	if err = s.runPushBatch(ctx, pipeline); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM push_deliveries WHERE status='retry'") != 1 || id("SELECT count(*) FROM push_deliveries WHERE status='sent'") != 2 {
		t.Fatal("partial failure changed successful delivery")
	}
	fcm.errorsByToken = nil
	exec("UPDATE push_deliveries SET next_attempt_at=now() WHERE status='retry'")
	if err = s.runPushBatch(ctx, pipeline); err != nil {
		t.Fatal(err)
	}
	if fcm.sendCount != 4 || id("SELECT count(*) FROM push_deliveries WHERE status='sent'") != 3 {
		t.Fatal("retry resent or lost device work")
	}
	enqueue("revoked", true)
	exec("UPDATE society_members SET status='suspended' WHERE user_id=$1", resident)
	if err = s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM notifications") != 3 {
		t.Fatal("revoked member materialized")
	}
	exec("UPDATE society_members SET status='active' WHERE user_id=$1", resident)
	guard.deny = true
	enqueue("inactive", true)
	if err = s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM notifications") != 3 {
		t.Fatal("inactive society materialized")
	}
	guard.deny = false
	if err := pipeline.SetNotificationPreferences(ctx, resident, society, contracts.NotificationPreferences{VisitorUpdatesPush: true, MaintenancePush: true, AnnouncementsPush: true, HubRepliesPush: false}); err != nil {
		t.Fatal(err)
	}
	enqueue("muted", true)
	if err := s.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM notifications") != 4 || id("SELECT count(*) FROM push_deliveries") != 3 {
		t.Fatal("muted push removed inbox or created device work")
	}

	// Fanout retries keep the same per-user event key.
	job := id("INSERT INTO hub_notification_fanout_jobs(society_id,post_id,actor_user_id) VALUES($1,$2,$3) RETURNING id", society, post, owner)
	claimed, err := pipeline.ClaimHubFanout(ctx, uuid.New())
	if err != nil || claimed == nil || claimed.ID != job {
		t.Fatalf("claim fanout: %v", err)
	}
	if count, e := pipeline.FanoutHubBatch(ctx, *claimed, 200); e != nil || count != 1 {
		t.Fatalf("fanout: %d %v", count, e)
	}
	if id("SELECT count(*) FROM notification_outbox WHERE event_key=$1", "hub.announcement:"+fmt.Sprint(post)) != 1 {
		t.Fatal("announcement fanout duplicated")
	}
	exec("UPDATE hub_notification_fanout_jobs SET completed_at=NULL,last_user_id=0 WHERE id=$1", job)
	claimed, err = pipeline.ClaimHubFanout(ctx, uuid.New())
	if err != nil || claimed == nil {
		t.Fatalf("reclaim fanout: %v", err)
	}
	if count, e := pipeline.FanoutHubBatch(ctx, *claimed, 200); e != nil || count != 1 {
		t.Fatalf("retry fanout: %d %v", count, e)
	}
	if id("SELECT count(*) FROM notification_outbox WHERE event_key=$1", "hub.announcement:"+fmt.Sprint(post)) != 1 {
		t.Fatal("fanout retry duplicated recipient")
	}
	deliveryID := id("SELECT id FROM push_deliveries ORDER BY id LIMIT 1")
	exec("UPDATE push_deliveries SET status='dead',attempt_count=6,lifetime_attempt_count=6 WHERE id=$1", deliveryID)
	if err := pipeline.ReplayPushDelivery(ctx, deliveryID); err != nil {
		t.Fatal(err)
	}
	if id("SELECT count(*) FROM push_deliveries WHERE id=$1 AND status='pending' AND attempt_count=0 AND lifetime_attempt_count=6 AND replay_count=1", deliveryID) != 1 {
		t.Fatal("replay discarded attempt history")
	}
	var wg sync.WaitGroup
	claimedDeliveries := make(chan *contracts.PushDelivery, 2)
	claimErrors := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, e := pipeline.ClaimPushDelivery(ctx, uuid.New())
			claimedDeliveries <- d
			claimErrors <- e
		}()
	}
	wg.Wait()
	close(claimedDeliveries)
	close(claimErrors)
	for e := range claimErrors {
		if e != nil {
			t.Fatal(e)
		}
	}
	var first *contracts.PushDelivery
	for d := range claimedDeliveries {
		if d != nil {
			if first != nil || d.ID != deliveryID {
				t.Fatal("concurrent workers claimed the same delivery")
			}
			first = d
		}
	}
	if first == nil {
		t.Fatal("pending delivery was not claimed")
	}
	exec("UPDATE push_deliveries SET lease_until=now()-interval '1 second' WHERE id=$1", deliveryID)
	second, err := pipeline.ClaimPushDelivery(ctx, uuid.New())
	if err != nil || second == nil || second.ID != deliveryID || second.LeaseToken == first.LeaseToken {
		t.Fatalf("expired lease was not reclaimed: %v", err)
	}
	if err := pipeline.FinishPushDelivery(ctx, first.ID, first.LeaseToken, "stale"); err == nil {
		t.Fatal("stale worker completed reclaimed delivery")
	}
	if err := pipeline.FinishPushDelivery(ctx, second.ID, second.LeaseToken, "accepted"); err != nil {
		t.Fatal(err)
	}
}
