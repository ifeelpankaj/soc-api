//go:build integration

package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go-server/pkg/database"
	"testing"
)

func TestWebPushSessionIsolation(t *testing.T) {
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:16-alpine", postgres.WithDatabase("web_push"), postgres.WithUsername("postgres"), postgres.WithPassword("postgres"), postgres.BasicWaitStrategies())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })
	connection, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	applyJobIntegrationMigrations(t, connection)
	pool, err := pgxpool.New(ctx, connection)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	repo := NewDeviceTokenRepository(&database.Database{Pool: pool}).(WebDeviceTokenRepository)
	user := insertReturningID(t, pool, `INSERT INTO users(full_name,email) VALUES ('Resident','web@example.com') RETURNING id`)
	other := insertReturningID(t, pool, `INSERT INTO users(full_name,email) VALUES ('Other','other@example.com') RETURNING id`)
	society := insertReturningID(t, pool, `INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES ('Home','WEB','active',$1,$1,NOW()) RETURNING id`, user)
	flat := insertReturningID(t, pool, `INSERT INTO flats(society_id,flat_number) VALUES ($1,'101') RETURNING id`, society)
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO flat_residents(society_id,flat_id,user_id) VALUES ($1,$2,$3),($1,$2,$4)`, society, flat, user, other)
	check := func(uid, sid, fid int64, want int) {
		t.Helper()
		rows, err := repo.ListEligibleWeb(ctx, uid, sid, fid)
		if err != nil || len(rows) != want {
			t.Fatalf("tokens=%d want=%d err=%v", len(rows), want, err)
		}
	}
	for _, device := range []string{"one", "two"} {
		if _, err := repo.UpsertWeb(ctx, user, "token-"+device, device, 0); err != nil {
			t.Fatal(err)
		}
	}
	check(user, society, flat, 2)
	check(user, society, flat+1, 0)
	// Password change/reset revokes both installations atomically through user version.
	exec(`UPDATE users SET session_version=session_version+1 WHERE id=$1`, user)
	check(user, society, flat, 0)
	if _, err := repo.UpsertWeb(ctx, user, "late", "one", 0); err == nil {
		t.Fatal("late registration bypassed revocation")
	}
	if _, err := repo.UpsertWeb(ctx, user, "fresh", "one", 1); err != nil {
		t.Fatal(err)
	}
	check(user, society, flat, 1)
	// Same browser transferred to another account must not receive the first account's alerts.
	if _, err := repo.UpsertWeb(ctx, other, "fresh", "one", 0); err != nil {
		t.Fatal(err)
	}
	check(user, society, flat, 0)
	check(other, society, flat, 1)
	exec(`UPDATE flat_residents SET status='inactive' WHERE user_id=$1`, other)
	check(other, society, flat, 0)
	if _, err := repo.UpsertWeb(ctx, other, "ineligible", "one", 0); err == nil {
		t.Fatal("inactive resident registered")
	}
	// Failed transaction does not delete a previously registered installation.
	exec(`UPDATE flat_residents SET status='active' WHERE user_id=$1`, other)
	exec(`CREATE FUNCTION reject_web_token() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'forced failure'; END $$; CREATE TRIGGER reject_web_token BEFORE INSERT ON device_tokens FOR EACH ROW EXECUTE FUNCTION reject_web_token()`)
	if _, err := repo.UpsertWeb(ctx, other, "replacement", "one", 0); err == nil {
		t.Fatal("expected insert failure")
	}
	check(other, society, flat, 1)
}
