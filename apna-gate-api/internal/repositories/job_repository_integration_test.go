//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestJobRepositoryCleanupDependenciesAndSocietyIsolation(t *testing.T) {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("apna_gate_jobs"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Skipf("PostgreSQL test container unavailable: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	applyJobIntegrationMigrations(t, connectionString)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	database := &database.Database{Pool: pool}
	repo := NewJobRepository(database, NewTransactionManager(database))

	seed := seedJobIntegrationData(t, pool)
	if _, err := pool.Exec(ctx, `UPDATE visitors SET photo_url='https://example.com/photo', photo_imagekit_file_id='photo-' || id, photo_imagekit_file_path='/visitors/' || id`); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2025, time.August, 31, 18, 30, 0, 0, time.UTC)
	to := time.Date(2025, time.September, 30, 18, 30, 0, 0, time.UTC)

	rowsA, err := repo.ListMonthlyVisitorReportRows(ctx, seed.societyA, from, to)
	if err != nil {
		t.Fatal(err)
	}
	rowsB, err := repo.ListMonthlyVisitorReportRows(ctx, seed.societyB, from, to)
	if err != nil {
		t.Fatal(err)
	}
	if len(rowsA) != 5 || len(rowsB) != 7 {
		t.Fatalf("isolated report counts = A:%d B:%d, want A:5 B:7", len(rowsA), len(rowsB))
	}
	for _, row := range rowsA {
		if row.VisitorName != "Society A Visitor" {
			t.Fatalf("society A report leaked row for %q", row.VisitorName)
		}
	}
	for _, row := range rowsB {
		if row.VisitorName != "Society B Visitor" {
			t.Fatalf("society B report leaked row for %q", row.VisitorName)
		}
	}

	// An invite cannot be removed while any surviving entry still references it,
	// and its logical short link must remain untouched as well.
	blockedInvites, err := repo.DeleteVisitorInvitesBatch(ctx, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), 100)
	if err != nil {
		t.Fatal(err)
	}
	if blockedInvites != 0 {
		t.Fatalf("deleted referenced visitor invites = %d, want 0", blockedInvites)
	}
	assertTableCount(t, pool, "short_links", 2)

	month := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
	cutoff := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Queue insertion failure must restore the entries and visitor together.
	if _, err := pool.Exec(ctx, `CREATE FUNCTION reject_image_queue() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'forced queue failure'; END $$;
	CREATE TRIGGER reject_image_queue BEFORE INSERT ON visitor_image_deletions FOR EACH ROW EXECUTE FUNCTION reject_image_queue()`); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{SocietyID: seed.societyA, ReportMonth: month}, from, to, cutoff, 100); err == nil {
		t.Fatal("expected queue failure")
	}
	assertTableCount(t, pool, "visitor_entries", 12)
	assertTableCount(t, pool, "visitors", 2)
	assertTableCount(t, pool, "visitor_image_deletions", 0)
	if _, err := pool.Exec(ctx, `DROP TRIGGER reject_image_queue ON visitor_image_deletions; DROP FUNCTION reject_image_queue()`); err != nil {
		t.Fatal(err)
	}
	// Partial entry deletion must preserve a visitor still used by other entries.
	if count, err := repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{SocietyID: seed.societyA, ReportMonth: month}, from, to, cutoff, 2); err != nil || count != 2 {
		t.Fatalf("partial cleanup=%d, %v", count, err)
	}
	assertTableCount(t, pool, "visitors", 2)
	assertTableCount(t, pool, "visitor_image_deletions", 0)
	deletedA, err := repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{
		SocietyID: seed.societyA, ReportMonth: month,
	}, from, to, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	deletedB, err := repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{
		SocietyID: seed.societyB, ReportMonth: month,
	}, from, to, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deletedA != 3 || deletedB != 0 {
		t.Fatalf("report-gated cleanup counts = A:%d B:%d, want A:3 B:0", deletedA, deletedB)
	}

	// The observable status never authorizes cleanup without sent_at.
	if _, err := pool.Exec(ctx, `INSERT INTO monthly_visitor_report_deliveries (society_id, report_month, status) VALUES ($1, DATE '2025-09-01', 'pending')`, seed.societyB); err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"pending", "failed", "processing", "sent"} {
		if _, err := pool.Exec(ctx, `UPDATE monthly_visitor_report_deliveries SET status = $2, sent_at = NULL WHERE society_id = $1 AND report_month = DATE '2025-09-01'`, seed.societyB, status); err != nil {
			t.Fatal(err)
		}
		deleted, err := repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{
			SocietyID: seed.societyB, ReportMonth: month,
		}, from, to, cutoff, 100)
		if err != nil {
			t.Fatal(err)
		}
		if deleted != 0 {
			t.Fatalf("status %q without sent_at deleted %d entries, want 0", status, deleted)
		}
	}
	if _, err := pool.Exec(ctx, `UPDATE monthly_visitor_report_deliveries SET status = 'sent', sent_at = NOW() WHERE society_id = $1 AND report_month = DATE '2025-09-01'`, seed.societyB); err != nil {
		t.Fatal(err)
	}
	deletedB, err = repo.DeleteVisitorEntriesBatch(ctx, models.VisitorReportCleanupWindow{
		SocietyID: seed.societyB, ReportMonth: month,
	}, from, to, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deletedB != 7 {
		t.Fatalf("deleted society B entries after sent_at = %d, want 7", deletedB)
	}

	// Normal mode excludes sent deliveries. Development resend mode exposes only
	// the requested latest eligible month and can safely claim/complete it again.
	normal, err := repo.ListClaimableReportDeliveries(ctx, 100, false, month)
	if err != nil {
		t.Fatal(err)
	}
	if len(normal) != 0 {
		t.Fatalf("normal mode returned sent deliveries: %+v", normal)
	}
	resends, err := repo.ListClaimableReportDeliveries(ctx, 100, true, month)
	if err != nil {
		t.Fatal(err)
	}
	if len(resends) != 2 {
		t.Fatalf("development resend deliveries=%d, want 2", len(resends))
	}
	claimedResend, err := repo.ClaimReportDelivery(ctx, seed.societyA, month, time.Now().Add(15*time.Minute), true)
	if err != nil || claimedResend == nil {
		t.Fatalf("claim development resend=%+v, %v", claimedResend, err)
	}
	if err := repo.CompleteReportDelivery(ctx, seed.societyA, month, []string{"owner-a@example.com"}, "resend-provider-id", true); err != nil {
		t.Fatal(err)
	}
	otherMonth := time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	resends, err = repo.ListClaimableReportDeliveries(ctx, 100, true, otherMonth)
	if err != nil {
		t.Fatal(err)
	}
	if len(resends) != 0 {
		t.Fatalf("resend leaked other months: %+v", resends)
	}

	assertTableCount(t, pool, "visitor_entry_events", 0)
	assertTableCount(t, pool, "societies", 2)
	assertTableCount(t, pool, "flats", 2)
	assertTableCount(t, pool, "visitors", 0)
	assertTableCount(t, pool, "visitor_image_deletions", 2)
	t.Run("orphan cleanup and rollback", func(t *testing.T) { testOrphanVisitorCleanup(t, pool, repo) })
	assertTableCount(t, pool, "users", 2)
	assertTableCount(t, pool, "visitor_invites", 1)

	// A failure after short-link deletion must roll back the whole invite batch.
	if _, err := pool.Exec(ctx, `
		CREATE FUNCTION reject_test_visitor_invite_delete() RETURNS trigger LANGUAGE plpgsql AS $$
		BEGIN RAISE EXCEPTION 'forced visitor invite delete failure'; END $$;
		CREATE TRIGGER reject_test_visitor_invite_delete BEFORE DELETE ON visitor_invites
		FOR EACH ROW EXECUTE FUNCTION reject_test_visitor_invite_delete();
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.DeleteVisitorInvitesBatch(ctx, cutoff, 100); err == nil {
		t.Fatal("visitor invite cleanup succeeded despite forced delete failure")
	}
	assertTableCount(t, pool, "visitor_invites", 1)
	assertTableCount(t, pool, "short_links", 2)
	if _, err := pool.Exec(ctx, `DROP TRIGGER reject_test_visitor_invite_delete ON visitor_invites; DROP FUNCTION reject_test_visitor_invite_delete()`); err != nil {
		t.Fatal(err)
	}

	deletedInvites, err := repo.DeleteVisitorInvitesBatch(ctx, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deletedInvites != 1 {
		t.Fatalf("deleted visitor invites = %d, want 1", deletedInvites)
	}
	assertTableCount(t, pool, "short_links", 1)

	deletedFlatInvites, err := repo.DeleteFlatMemberInvitesBatch(ctx, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deletedFlatInvites != 1 {
		t.Fatalf("deleted flat member invites = %d, want 1", deletedFlatInvites)
	}
	assertTableCount(t, pool, "short_links", 0)
	assertTableCount(t, pool, "flats", 2)
	assertTableCount(t, pool, "users", 2)

	deletedNotifications, err := repo.DeleteNotificationsBatch(ctx, cutoff, 100)
	if err != nil {
		t.Fatal(err)
	}
	if deletedNotifications != 1 {
		t.Fatalf("deleted notifications = %d, want 1", deletedNotifications)
	}
	assertTableCount(t, pool, "users", 2)

	// A live processing lease blocks claims; an expired one is recoverable.
	recoveryMonth := time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC)
	if _, err := pool.Exec(ctx, `INSERT INTO monthly_visitor_report_deliveries (society_id, report_month, status, attempt_count, processing_until) VALUES ($1, $2, 'processing', 1, NOW() + INTERVAL '1 hour')`, seed.societyA, recoveryMonth); err != nil {
		t.Fatal(err)
	}
	claimed, err := repo.ClaimReportDelivery(ctx, seed.societyA, recoveryMonth, time.Now().Add(15*time.Minute), false)
	if err != nil {
		t.Fatal(err)
	}
	if claimed != nil {
		t.Fatal("delivery with live processing lease was claimed")
	}
	if _, err := pool.Exec(ctx, `UPDATE monthly_visitor_report_deliveries SET processing_until = NOW() - INTERVAL '1 minute' WHERE society_id = $1 AND report_month = $2`, seed.societyA, recoveryMonth); err != nil {
		t.Fatal(err)
	}
	claimed, err = repo.ClaimReportDelivery(ctx, seed.societyA, recoveryMonth, time.Now().Add(15*time.Minute), false)
	if err != nil {
		t.Fatal(err)
	}
	if claimed == nil || claimed.AttemptCount != 2 {
		t.Fatalf("recovered claim = %#v, want attempt 2", claimed)
	}
}

type jobIntegrationSeed struct {
	societyA int64
	societyB int64
}

func seedJobIntegrationData(t *testing.T, pool *pgxpool.Pool) jobIntegrationSeed {
	t.Helper()
	ctx := context.Background()
	old := time.Date(2025, time.September, 15, 12, 0, 0, 0, time.UTC)
	userA := insertReturningID(t, pool, `INSERT INTO users (full_name, email, email_verified) VALUES ('Owner A', 'owner-a@example.com', TRUE) RETURNING id`)
	userB := insertReturningID(t, pool, `INSERT INTO users (full_name, email, email_verified) VALUES ('Owner B', 'owner-b@example.com', TRUE) RETURNING id`)
	societyA := insertReturningID(t, pool, `INSERT INTO societies (name, society_code, status, created_by, approved_by, approved_at) VALUES ('Society A', 'A', 'active', $1, $1, NOW()) RETURNING id`, userA)
	societyB := insertReturningID(t, pool, `INSERT INTO societies (name, society_code, status, created_by, approved_by, approved_at) VALUES ('Society B', 'B', 'active', $1, $1, NOW()) RETURNING id`, userB)
	if _, err := pool.Exec(ctx, `INSERT INTO society_members (society_id, user_id, role, status) VALUES ($1, $2, 'owner', 'active'), ($3, $4, 'owner', 'active')`, societyA, userA, societyB, userB); err != nil {
		t.Fatal(err)
	}
	flatA := insertReturningID(t, pool, `INSERT INTO flats (society_id, flat_number, created_by) VALUES ($1, 'A-101', $2) RETURNING id`, societyA, userA)
	flatB := insertReturningID(t, pool, `INSERT INTO flats (society_id, flat_number, created_by) VALUES ($1, 'B-101', $2) RETURNING id`, societyB, userB)
	visitorA := insertReturningID(t, pool, `INSERT INTO visitors (full_name, phone_number) VALUES ('Society A Visitor', '+910000000001') RETURNING id`)
	visitorB := insertReturningID(t, pool, `INSERT INTO visitors (full_name, phone_number) VALUES ('Society B Visitor', '+910000000002') RETURNING id`)
	inviteA := insertReturningID(t, pool, `INSERT INTO visitor_invites (society_id, flat_id, created_by, purpose, token_hash, status, expires_at, used_at, created_at, updated_at) VALUES ($1, $2, $3, 'guest', 'visitor-token-a', 'used', $4, $4, $4, $4) RETURNING id`, societyA, flatA, userA, old)

	var firstEntryA int64
	for i := 0; i < 5; i++ {
		var inviteID any
		if i == 0 {
			inviteID = inviteA
		}
		entryID := insertReturningID(t, pool, `INSERT INTO visitor_entries (society_id, flat_id, visitor_id, invite_id, source, purpose, status, checked_in_at, checked_out_at, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, 'resident_link', 'guest', 'checked_out', $5, $5, $6, $5, $5) RETURNING id`, societyA, flatA, visitorA, inviteID, old.Add(time.Duration(i)*time.Minute), userA)
		if i == 0 {
			firstEntryA = entryID
		}
	}
	for i := 0; i < 7; i++ {
		insertReturningID(t, pool, `INSERT INTO visitor_entries (society_id, flat_id, visitor_id, source, purpose, status, checked_in_at, checked_out_at, created_by, created_at, updated_at) VALUES ($1, $2, $3, 'guard_entry', 'guest', 'checked_out', $4, $4, $5, $4, $4) RETURNING id`, societyB, flatB, visitorB, old.Add(time.Duration(i)*time.Minute), userB)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO visitor_entry_events (visitor_entry_id, society_id, event_type) VALUES ($1, $2, 'checked_out')`, firstEntryA, societyA); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO monthly_visitor_report_deliveries (society_id, report_month, status, sent_at) VALUES ($1, DATE '2025-09-01', 'sent', NOW())`, societyA); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO short_links (short_code, resource_type, resource_id) VALUES ('visitor-link-a', 'visitor_invite', $1)`, inviteA); err != nil {
		t.Fatal(err)
	}
	flatInvite := insertReturningID(t, pool, `INSERT INTO flat_member_invites (society_id, flat_id, invited_by, role, email, full_name, token_hash, expires_at, status, created_at, updated_at) VALUES ($1, $2, $3, 'tenant', 'tenant@example.com', 'Tenant', 'flat-token-a', $4, 'expired', $4, $4) RETURNING id`, societyA, flatA, userA, old)
	if _, err := pool.Exec(ctx, `INSERT INTO short_links (short_code, resource_type, resource_id) VALUES ('flat-link-a', 'member_invite', $1)`, flatInvite); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO notifications (id, user_id, society_id, type, title, body, created_at) VALUES ('11111111-1111-1111-1111-111111111111', $1, $2, 'visitor.test', 'Test', 'Test', $3)`, userA, societyA, old); err != nil {
		t.Fatal(err)
	}
	return jobIntegrationSeed{societyA: societyA, societyB: societyB}
}

func applyJobIntegrationMigrations(t *testing.T, connectionString string) {
	t.Helper()
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	source := &migrate.FileMigrationSource{Dir: filepath.Join("..", "..", "migrations")}
	if _, err := migrate.Exec(sqlDB, "postgres", source, migrate.Up); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
}

func insertReturningID(t *testing.T, pool *pgxpool.Pool, query string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(), query, args...).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func assertTableCount(t *testing.T, pool *pgxpool.Pool, table string, want int64) {
	t.Helper()
	var got int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table) // table names are fixed test constants.
	if err := pool.QueryRow(context.Background(), query).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("%s count = %d, want %d", table, got, want)
	}
}

func testOrphanVisitorCleanup(t *testing.T, pool *pgxpool.Pool, repo JobRepository) {
	ctx := context.Background()
	cutoff := time.Now().Add(-90 * 24 * time.Hour)
	oldPhoto := insertReturningID(t, pool, `INSERT INTO visitors(full_name, phone_number, updated_at, photo_url, photo_imagekit_file_id, photo_imagekit_file_path) VALUES ('Old photo','123', $1, 'https://example.com/old', 'old-photo', '/old-photo') RETURNING id`, cutoff.Add(-time.Hour))
	oldPlain := insertReturningID(t, pool, `INSERT INTO visitors(full_name, phone_number, updated_at) VALUES ('Old plain','124', $1) RETURNING id`, cutoff.Add(-time.Hour))
	recent := insertReturningID(t, pool, `INSERT INTO visitors(full_name, phone_number) VALUES ('Recent','125') RETURNING id`)
	if count, err := repo.DeleteOrphanVisitorsBatch(ctx, cutoff, 1); err != nil || count != 1 {
		t.Fatalf("first batch=%d, %v", count, err)
	}
	if count, err := repo.DeleteOrphanVisitorsBatch(ctx, cutoff, 1); err != nil || count != 1 {
		t.Fatalf("second batch=%d, %v", count, err)
	}
	if count, err := repo.DeleteOrphanVisitorsBatch(ctx, cutoff, 1); err != nil || count != 0 {
		t.Fatalf("recent cleanup=%d, %v", count, err)
	}
	var remaining int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM visitors WHERE id=ANY($1)`, []int64{oldPhoto, oldPlain}).Scan(&remaining); err != nil || remaining != 0 {
		t.Fatalf("old remaining=%d, %v", remaining, err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM visitors WHERE id=$1`, recent).Scan(&remaining); err != nil || remaining != 1 {
		t.Fatalf("recent remaining=%d, %v", remaining, err)
	}
	items, err := repo.ListVisitorImageDeletions(ctx, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range items {
		if item.FileID == "old-photo" {
			found = true
		}
	}
	if !found {
		t.Fatal("old photo was not queued")
	}

	// An in-flight entry INSERT takes a foreign-key lock on its visitor. Orphan
	// cleanup skips that row, and preserves it after the insert commits as well.
	concurrent := insertReturningID(t, pool, `INSERT INTO visitors(full_name, phone_number, updated_at) VALUES ('Concurrent','126', $1) RETURNING id`, cutoff.Add(-time.Hour))
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	var entryID int64
	if err := tx.QueryRow(ctx, `INSERT INTO visitor_entries(society_id, flat_id, visitor_id, source, purpose, status, created_by)
	SELECT society_id, id, $1, 'guard_entry', 'guest', 'waiting_approval', created_by FROM flats ORDER BY id LIMIT 1 RETURNING id`, concurrent).Scan(&entryID); err != nil {
		t.Fatal(err)
	}
	boundedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if count, err := repo.DeleteOrphanVisitorsBatch(boundedCtx, cutoff, 10); err != nil || count != 0 {
		t.Fatalf("concurrent cleanup=%d, %v", count, err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if count, err := repo.DeleteOrphanVisitorsBatch(ctx, cutoff, 10); err != nil || count != 0 {
		t.Fatalf("referenced cleanup=%d, %v", count, err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM visitor_entries WHERE id=$1`, entryID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM visitors WHERE id=ANY($1)`, []int64{recent, concurrent}); err != nil {
		t.Fatal(err)
	}
}
