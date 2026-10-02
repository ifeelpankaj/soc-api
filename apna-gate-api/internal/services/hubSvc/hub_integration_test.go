//go:build integration

package hubsvc

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/config"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/repositories/contracts"
	"go-server/internal/requestctx"
	"go-server/internal/testutil"
	"go-server/pkg/database"
)

type hubGuard struct{ deny bool }

func (g *hubGuard) EnsureSocietyOperational(ctx context.Context, _ int64) error {
	if requestctx.HasDeveloperGuardBypass(ctx) {
		return errors.New("bypass leaked")
	}
	if g.deny {
		return ErrForbidden
	}
	return nil
}

type hubStorage struct {
	mu      sync.Mutex
	deletes int
	fail    bool
}

func (f *hubStorage) Upload(_ context.Context, in models.ImageUploadInput) (*models.StoredImage, error) {
	if f.fail {
		return nil, errors.New("provider failure")
	}
	return &models.StoredImage{FileID: in.FileName, Path: "/" + in.Folder + "/" + in.FileName, URL: "https://example.invalid/" + in.FileName}, nil
}
func (f *hubStorage) Delete(context.Context, string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail {
		return errors.New("provider failure")
	}
	f.deletes++
	return nil
}
func (f *hubStorage) SignedURL(_ string, t time.Time, _ ...models.ImageVariant) (string, error) {
	return fmt.Sprintf("https://example.invalid/private?expires=%d", t.Unix()), nil
}

func TestHubIntegration(t *testing.T) {
	ctx := context.Background()
	pg := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	sqlDB, e := sql.Open("pgx", pg.DSN)
	if e != nil {
		t.Fatal(e)
	}
	defer sqlDB.Close()
	source := &migrate.FileMigrationSource{Dir: "../../../migrations"}
	if _, e = migrate.ExecMax(sqlDB, "postgres", source, migrate.Up, 32); e != nil {
		t.Fatal(e)
	}
	pool, e := pgxpool.New(ctx, pg.DSN)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, query, args...); e != nil {
			t.Fatal(e)
		}
	}
	id := func(query string, args ...any) int64 {
		t.Helper()
		var n int64
		if e := pool.QueryRow(ctx, query, args...).Scan(&n); e != nil {
			t.Fatal(e)
		}
		return n
	}
	user := func(name string) int64 {
		return id("INSERT INTO users(full_name,email,email_verified) VALUES($1,$2,true) RETURNING id", name, name+"@example.test")
	}
	owner, resident, second, staff, outsider, adminUser := user("owner"), user("resident"), user("second"), user("staff"), user("developer"), user("admin")
	exec("UPDATE users SET global_role='developer' WHERE id=$1", outsider)
	society := id("INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Hub','HUB','active',$1,$1,now()) RETURNING id", owner)
	other := id("INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Other','OTHER','active',$1,$1,now()) RETURNING id", owner)
	exec("INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'owner','active'),($1,$3,'resident','active'),($1,$4,'resident','active'),($1,$5,'staff','active'),($1,$6,'admin','active'),($7,$2,'owner','active')", society, owner, resident, second, staff, adminUser, other)
	t.Run("migration and backfill", func(t *testing.T) {
		if n, e := migrate.Exec(sqlDB, "postgres", source, migrate.Up); e != nil || n != 4 {
			t.Fatalf("up: %d %v", n, e)
		}
		if n := id("SELECT count(*) FROM society_channels"); n != 4 {
			t.Fatalf("backfill %d", n)
		}
		if n, e := migrate.ExecMax(sqlDB, "postgres", source, migrate.Down, 2); e != nil || n != 2 {
			t.Fatalf("down: %d %v", n, e)
		}
		if n, e := migrate.Exec(sqlDB, "postgres", source, migrate.Up); e != nil || n != 2 {
			t.Fatalf("up again: %d %v", n, e)
		}
	})
	d := &database.Database{Pool: pool}
	tx := repository.NewTransactionManager(d)
	repo := repository.NewHubRepository(d, tx)
	guard := &hubGuard{}
	storage := &hubStorage{}
	s := New(repo, guard, storage, "test", config.HubConfig{PostsPerWindow: 1000, PostWindowSeconds: 600, CommentsPerWindow: 1000, CommentWindowSeconds: 60, UploadMaxBytes: 10485760, ImageMaxWidth: 1600, ImageMaxHeight: 1600})
	announcement := id("SELECT id FROM society_channels WHERE society_id=$1 AND type='announcement'", society)
	community := id("SELECT id FROM society_channels WHERE society_id=$1 AND type='community'", society)
	otherChannel := id("SELECT id FROM society_channels WHERE society_id=$1 AND type='community'", other)
	post := func(actor, ch int64) models.HubContent {
		t.Helper()
		p, e := s.SavePost(ctx, society, actor, ch, 0, models.HubPostRequest{Body: ptr("Hello neighbours")})
		if e != nil {
			t.Fatal(e)
		}
		return p
	}
	comment := func(actor, p int64, parent *int64) models.HubContent {
		t.Helper()
		c, e := s.SaveComment(ctx, society, actor, p, 0, models.HubCommentRequest{Body: ptr("A reply"), ParentID: parent})
		if e != nil {
			t.Fatal(e)
		}
		return c
	}
	upload := func(actor int64) models.HubUploadView {
		t.Helper()
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		part, e := w.CreateFormFile("file", "notice.pdf")
		if e != nil {
			t.Fatal(e)
		}
		_, _ = part.Write([]byte("%PDF-1.7\n1 0 obj\n<<>>\nendobj\n%%EOF"))
		_ = w.Close()
		r := httptest.NewRequest("POST", "/", &b)
		r.Header.Set("Content-Type", w.FormDataContentType())
		u, e := s.Upload(ctx, society, actor, httptest.NewRecorder(), r)
		if e != nil {
			t.Fatal(e)
		}
		return u
	}
	t.Run("provisioning atomic", func(t *testing.T) {
		r := repository.NewSocietyRepository(d)
		v := &models.Society{Name: "New", SocietyCode: "NEW", Country: "India", CreatedBy: owner}
		if e := r.Create(ctx, v); e != nil {
			t.Fatal(e)
		}
		if id("SELECT count(*) FROM society_channels WHERE society_id=$1", v.ID) != 2 {
			t.Fatal("missing channels")
		}
		v = &models.Society{Name: "Rollback", SocietyCode: "ROLLBACK", Country: "India", CreatedBy: owner}
		_ = tx.WithTransaction(ctx, func(c context.Context) error {
			if e := r.Create(c, v); e != nil {
				t.Fatal(e)
			}
			return errors.New("rollback")
		})
		if id("SELECT count(*) FROM society_channels WHERE society_id=$1", v.ID) != 0 {
			t.Fatal("channels survived rollback")
		}
	})
	t.Run("access matrix", func(t *testing.T) {
		for _, actor := range []int64{owner, adminUser, resident, second} {
			if _, e := s.Channels(ctx, society, actor); e != nil {
				t.Fatal(e)
			}
		}
		for _, actor := range []int64{staff, outsider} {
			if _, e := s.Channels(requestctx.WithDeveloperGuardBypass(ctx), society, actor); !errors.Is(e, ErrForbidden) {
				t.Fatalf("access %d: %v", actor, e)
			}
		}
		exec("UPDATE society_members SET status='suspended' WHERE user_id=$1 AND society_id=$2", second, society)
		if _, e := s.Channels(ctx, society, second); !errors.Is(e, ErrForbidden) {
			t.Fatal("suspended member allowed")
		}
		exec("UPDATE society_members SET status='active' WHERE user_id=$1 AND society_id=$2", second, society)
		guard.deny = true
		if _, e := s.Channels(requestctx.WithDeveloperGuardBypass(ctx), society, owner); !errors.Is(e, ErrForbidden) {
			t.Fatal("subscription bypass")
		}
		guard.deny = false
		if _, e := s.SavePost(ctx, society, resident, announcement, 0, models.HubPostRequest{Body: ptr("no")}); !errors.Is(e, ErrForbidden) {
			t.Fatal("resident announced", e)
		}
		if _, e := s.SavePost(ctx, society, owner, otherChannel, 0, models.HubPostRequest{Body: ptr("no")}); !errors.Is(e, ErrNotFound) {
			t.Fatal("cross society channel", e)
		}
	})
	t.Run("announcements notifications and controls", func(t *testing.T) {
		p, e := s.SavePost(ctx, society, owner, announcement, 0, models.HubPostRequest{Body: ptr("Notice"), IsImportant: ptr(true), IsPinned: ptr(true)})
		if e != nil {
			t.Fatal(e)
		}
		pipeline := repository.NewNotificationRepository(d).(contracts.PipelineRepository)
		job, e := pipeline.ClaimHubFanout(ctx, uuid.New())
		if e != nil || job == nil || job.PostID != p.ID {
			t.Fatal("announcement fanout job missing", e)
		}
		if count, e := pipeline.FanoutHubBatch(ctx, *job, 200); e != nil || count != 2 {
			t.Fatalf("announcement fanout: %d %v", count, e)
		}
		if n := id("SELECT count(*) FROM notification_outbox WHERE event_key=$1 AND push_enabled AND payload->'Data'->>'priority'='high'", fmt.Sprintf("hub.announcement:%d", p.ID)); n != 2 {
			t.Fatalf("resident deliveries %d", n)
		}
		channels, e := s.Channels(ctx, society, resident)
		if e != nil || channels[0].ImportantAnnouncement == nil {
			t.Fatal("missing banner", e)
		}
		feed, e := s.Feed(ctx, society, resident, announcement, "", 20, nil, false)
		if e != nil || len(feed.Pinned) != 1 {
			t.Fatal("pinned feed", e)
		}
		if e = s.Control(ctx, society, resident, p.ID, nil, ptr(true)); !errors.Is(e, ErrForbidden) {
			t.Fatal("resident locked")
		}
		if e = s.Control(ctx, society, adminUser, p.ID, nil, ptr(true)); e != nil {
			t.Fatal(e)
		}
		if _, e = s.SaveComment(ctx, society, owner, p.ID, 0, models.HubCommentRequest{Body: ptr("no")}); !errors.Is(e, ErrConflict) {
			t.Fatal("admin bypassed lock", e)
		}
		if _, e = s.SavePost(ctx, society, adminUser, 0, p.ID, models.HubPostRequest{Body: ptr("rewrite")}); !errors.Is(e, ErrForbidden) {
			t.Fatal("admin rewrote author", e)
		}
	})
	t.Run("threads and soft deletion", func(t *testing.T) {
		p := post(resident, community)
		c := comment(second, p.ID, nil)
		reply := comment(owner, p.ID, &c.ID)
		if _, e := s.SaveComment(ctx, society, resident, p.ID, 0, models.HubCommentRequest{Body: ptr("deep"), ParentID: &reply.ID}); !errors.Is(e, ErrInvalid) {
			t.Fatal("deep nesting accepted", e)
		}
		otherPost := post(owner, community)
		if _, e := s.SaveComment(ctx, society, resident, otherPost.ID, 0, models.HubCommentRequest{Body: ptr("wrong parent"), ParentID: &c.ID}); !errors.Is(e, ErrInvalid) {
			t.Fatal("cross post parent", e)
		}
		if e := s.Delete(ctx, society, second, c.ID, true, ""); e != nil {
			t.Fatal(e)
		}
		page, e := s.Feed(ctx, society, resident, p.ID, "", 20, nil, true)
		if e != nil || len(page.Items) != 2 || page.Items[0].Body != "" || page.Items[0].Author != nil {
			t.Fatalf("placeholder %+v %v", page, e)
		}
		if _, e := s.SaveComment(ctx, society, second, 0, c.ID, models.HubCommentRequest{Body: ptr("resurrect")}); !errors.Is(e, ErrNotFound) {
			t.Fatal("edited deleted comment", e)
		}
		if n := id("SELECT count(*) FROM notification_outbox WHERE event_key=$1", fmt.Sprintf("hub.reply:%d", reply.ID)); n != 2 {
			t.Fatalf("reply recipients %d", n)
		}
	})
	t.Run("read pagination reactions reports", func(t *testing.T) {
		p1, p2 := post(resident, community), post(resident, community)
		exec("UPDATE channel_posts SET created_at='2020-01-01' WHERE id IN ($1,$2)", p1.ID, p2.ID)
		page, e := s.Feed(ctx, society, resident, community, "", 1, nil, false)
		if e != nil || page.NextCursor == "" {
			t.Fatal("cursor", e)
		}
		if _, e = s.Feed(ctx, society, resident, announcement, page.NextCursor, 1, nil, false); !errors.Is(e, ErrInvalid) {
			t.Fatal("cross channel cursor")
		}
		if e = s.Read(ctx, society, resident, community, p2.ID); e != nil {
			t.Fatal(e)
		}
		if e = s.Read(ctx, society, resident, community, p1.ID); e != nil {
			t.Fatal(e)
		}
		if n := id("SELECT last_read_post_id FROM channel_read_states WHERE user_id=$1 AND channel_id=$2", resident, community); n != p2.ID {
			t.Fatal("read regressed")
		}
		if e = s.Read(ctx, society, resident, announcement, p1.ID); !errors.Is(e, ErrNotFound) {
			t.Fatal("cross channel read")
		}
		for i := 0; i < 2; i++ {
			if e = s.Reaction(ctx, society, second, p1.ID, false, false, "helpful"); e != nil {
				t.Fatal(e)
			}
		}
		if n := id("SELECT count(*) FROM channel_reactions WHERE post_id=$1", p1.ID); n != 1 {
			t.Fatal("duplicate reaction")
		}
		if n := id("SELECT count(*) FROM notification_outbox WHERE payload->'Data'->>'post_id'=$1 AND event_key LIKE 'hub.reaction:%' AND NOT push_enabled", fmt.Sprint(p1.ID)); n != 1 {
			t.Fatal("reaction notification", n)
		}
		if e = s.Reaction(ctx, society, second, p1.ID, false, true, "helpful"); e != nil {
			t.Fatal(e)
		}
		report, e := s.Report(ctx, society, second, p1.ID, "Inappropriate")
		if e != nil {
			t.Fatal(e)
		}
		again, e := s.Report(ctx, society, second, p1.ID, "Again")
		if e != nil || again.ID != report.ID {
			t.Fatal("duplicate report", e)
		}
		if e = s.Resolve(ctx, society, resident, report.ID, models.HubResolveRequest{Status: "actioned", Note: "remove"}); !errors.Is(e, ErrForbidden) {
			t.Fatal("resident moderated")
		}
		if e = s.Resolve(ctx, society, owner, report.ID, models.HubResolveRequest{Status: "actioned", Note: "remove"}); e != nil {
			t.Fatal(e)
		}
		if _, e = s.Post(ctx, society, resident, p1.ID); !errors.Is(e, ErrNotFound) {
			t.Fatal("removed post visible")
		}
		reports, e := s.Reports(ctx, society, owner, "actioned", "", 20)
		if e != nil || len(reports.Items) != 1 {
			t.Fatal("resolved reports", e)
		}
	})
	t.Run("upload ownership claim races cleanup and rollback", func(t *testing.T) {
		u := upload(resident)
		if _, e := s.Attachment(ctx, society, second, u.ID); !errors.Is(e, ErrNotFound) {
			t.Fatal("pending upload leaked")
		}
		if _, e := s.SavePost(ctx, society, second, community, 0, models.HubPostRequest{Body: ptr("steal"), AttachmentIDs: ptr([]int64{u.ID})}); !errors.Is(e, ErrForbidden) {
			t.Fatal("stole upload", e)
		}
		results := make(chan error, 2)
		var wg sync.WaitGroup
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := s.SavePost(ctx, society, resident, community, 0, models.HubPostRequest{Body: ptr("attachment"), AttachmentIDs: ptr([]int64{u.ID})})
				results <- e
			}()
		}
		wg.Wait()
		close(results)
		successes := 0
		for e := range results {
			if e == nil {
				successes++
			} else if !errors.Is(e, ErrConflict) {
				t.Fatal(e)
			}
		}
		if successes != 1 {
			t.Fatalf("claims %d", successes)
		}
		if _, e := s.Attachment(ctx, society, second, u.ID); e != nil {
			t.Fatal(e)
		}
		if e := s.DeleteUpload(ctx, society, resident, u.ID); !errors.Is(e, ErrConflict) {
			t.Fatal("deleted claimed upload")
		}
		expired := upload(resident)
		exec("UPDATE channel_uploads SET expires_at=now()-interval '1 second' WHERE id=$1", expired.ID)
		if _, e := s.SavePost(ctx, society, resident, community, 0, models.HubPostRequest{Body: ptr("expired"), AttachmentIDs: ptr([]int64{expired.ID})}); !errors.Is(e, ErrConflict) {
			t.Fatal("claimed expired", e)
		}
		storage.fail = true
		if e := s.CleanupUploads(ctx); e == nil {
			t.Fatal("provider failure swallowed")
		}
		storage.fail = false
		exec("UPDATE channel_uploads SET lease_until=now()-interval '1 second' WHERE id=$1", expired.ID)
		if e := s.CleanupUploads(ctx); e != nil {
			t.Fatal(e)
		}
		if _, e := s.Attachment(ctx, society, resident, expired.ID); !errors.Is(e, ErrNotFound) {
			t.Fatal("expired file visible")
		}
		count := id("SELECT count(*) FROM channel_posts")
		exec(`CREATE FUNCTION fail_hub_outbox() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'test fanout failure'; END $$; CREATE TRIGGER fail_hub_outbox BEFORE INSERT ON hub_notification_fanout_jobs FOR EACH ROW EXECUTE FUNCTION fail_hub_outbox()`)
		if _, e := s.SavePost(ctx, society, owner, announcement, 0, models.HubPostRequest{Body: ptr("rollback")}); e == nil {
			t.Fatal("expected fanout job failure")
		}
		if n := id("SELECT count(*) FROM channel_posts"); n != count {
			t.Fatal("post committed without fanout job")
		}
		exec("DROP TRIGGER fail_hub_outbox ON hub_notification_fanout_jobs; DROP FUNCTION fail_hub_outbox()")
	})
	t.Run("concurrent posting limits", func(t *testing.T) {
		limited := New(repo, guard, storage, "test", config.HubConfig{PostsPerWindow: 1, PostWindowSeconds: 3600, CommentsPerWindow: 1, CommentWindowSeconds: 3600, UploadMaxBytes: 10485760, ImageMaxWidth: 1600, ImageMaxHeight: 1600})
		exec("UPDATE channel_posts SET created_at=now()-interval '2 hours' WHERE author_id=$1", resident)
		errs := make(chan error, 2)
		var wg sync.WaitGroup
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := limited.SavePost(ctx, society, resident, community, 0, models.HubPostRequest{Body: ptr("limited")})
				errs <- e
			}()
		}
		wg.Wait()
		close(errs)
		ok, denied := 0, 0
		for e := range errs {
			var rate *RateError
			if e == nil {
				ok++
			} else if errors.As(e, &rate) {
				denied++
			} else {
				t.Fatal(e)
			}
		}
		if ok != 1 || denied != 1 {
			t.Fatalf("limits %d/%d", ok, denied)
		}
	})
	t.Run("database tenant constraints", func(t *testing.T) {
		p := post(owner, community)
		c := comment(owner, p.ID, nil)
		p2 := post(owner, community)
		if _, e := pool.Exec(ctx, "INSERT INTO channel_comments(post_id,author_id,parent_id,body) VALUES($1,$2,$3,'bad')", p2.ID, owner, c.ID); e == nil {
			t.Fatal("cross post parent stored")
		}
		if _, e := pool.Exec(ctx, "INSERT INTO channel_posts(society_id,channel_id,author_id,body) VALUES($1,$2,$3,'bad')", society, otherChannel, owner); e == nil {
			t.Fatal("cross society post stored")
		}
		if _, e := pool.Exec(ctx, "INSERT INTO channel_reactions(user_id,post_id,comment_id,reaction_type) VALUES($1,$2,$3,'like')", owner, p.ID, c.ID); e == nil {
			t.Fatal("ambiguous reaction stored")
		}
		if e := repo.HubRead(ctx, contracts.HubReadInput{SocietyID: society, UserID: owner, ChannelID: community, PostID: p.ID}); e != nil {
			t.Fatal(e)
		}
	})
}
