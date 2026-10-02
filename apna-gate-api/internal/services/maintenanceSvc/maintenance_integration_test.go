//go:build integration

package maintenancesvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/repositories/contracts"
	"go-server/internal/testutil"
	"go-server/pkg/database"
	"sync"
	"testing"
	"time"
)

type integrationPush struct {
	calls int
	fail  bool
}

func (p *integrationPush) SendToUser(context.Context, int64, models.NotificationPayload) error {
	p.calls++
	if p.fail {
		return errors.New("temporary provider failure")
	}
	return nil
}

func TestMaintenanceBillingIntegration(t *testing.T) {
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
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			t.Fatal(err)
		}
	}
	id := func(query string, args ...any) int64 {
		t.Helper()
		var n int64
		if err := pool.QueryRow(ctx, query, args...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	count := func(table string, want int64) {
		t.Helper()
		got := id("SELECT count(*) FROM " + table)
		if got != want {
			t.Fatalf("%s count=%d want=%d", table, got, want)
		}
	}
	owner := id(`INSERT INTO users(full_name,email,email_verified) VALUES('Billing Owner','billing-owner@example.com',true) RETURNING id`)
	resident := id(`INSERT INTO users(full_name,email,email_verified) VALUES('Billing Resident','billing-resident@example.com',true) RETURNING id`)
	outsider := id(`INSERT INTO users(full_name,email,email_verified) VALUES('Other Resident','other-resident@example.com',true) RETURNING id`)
	staff := id(`INSERT INTO users(full_name,email,email_verified) VALUES('Staff','staff@example.com',true) RETURNING id`)
	society := id(`INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Maintenance','MAINT','active',$1,$1,now()) RETURNING id`, owner)
	otherSociety := id(`INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Other','OTHER','active',$1,$1,now()) RETURNING id`, outsider)
	exec(`INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'owner','active'),($1,$3,'resident','active'),($1,$4,'staff','active'),($5,$6,'owner','active')`, society, owner, resident, staff, otherSociety, outsider)
	flat1 := id(`INSERT INTO flats(society_id,flat_number,status,area_sqft_hundredths,flat_type) VALUES($1,'101','occupied',100000,'2 BHK') RETURNING id`, society)
	flat2 := id(`INSERT INTO flats(society_id,flat_number,status) VALUES($1,'102','blocked') RETURNING id`, society)
	exec(`INSERT INTO flats(society_id,flat_number,is_active) VALUES($1,'inactive',false)`, society)
	exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,role,status,is_primary) VALUES($1,$2,$3,'owner','active',true)`, society, flat1, resident)
	d := &database.Database{Pool: pool}
	repo := repository.NewMaintenanceRepository(d, repository.NewTransactionManager(d))
	push := &integrationPush{fail: true}
	s := New(repo, allowOperational{}, []string{"vacant", "occupied", "blocked"}, repository.NewNotificationRepository(d), push)
	s.now = func() time.Time { return time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC) }
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.PricingModel = "hybrid"
	settings.FixedPaise = 50000
	settings.AreaRatePaise = 300
	if _, err = s.SaveSettings(ctx, society, owner, settings); err != nil {
		t.Fatal(err)
	}
	if _, err = s.SaveSettings(ctx, society, staff, settings); err == nil {
		t.Fatal("staff changed settings")
	}
	preview, err := s.Preview(ctx, society, owner, "2026-10")
	if err != nil || len(preview.Issues) != 1 || preview.Issues[0].FlatID != flat2 {
		t.Fatalf("preview: %+v %v", preview, err)
	}
	if _, err = s.Generate(ctx, society, owner, "2026-10"); err == nil {
		t.Fatal("missing area did not block")
	}
	count("maintenance_bills", 0)
	count("maintenance_billing_runs", 0)
	// Flat API repository round-trip uses typed decimal strings and exact storage.
	area := "1000.25"
	kind := "2 BHK"
	flatRepo := repository.NewFlatRepository(d)
	updated, err := flatRepo.Update(ctx, &models.FlatFilter{ID: &flat2, SocietyID: &society}, &contracts.UpdateFlatInput{AreaSqft: &area, FlatType: &kind})
	if err != nil || updated.AreaSqft == nil || *updated.AreaSqft != area {
		t.Fatalf("flat update: %+v %v", updated, err)
	}
	// A queue failure must roll back bills, items and the successful-run marker.
	exec(`CREATE FUNCTION reject_maintenance_queue() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected queue failure'; END $$`)
	exec(`CREATE TRIGGER reject_maintenance_queue BEFORE INSERT ON maintenance_notification_deliveries FOR EACH ROW EXECUTE FUNCTION reject_maintenance_queue()`)
	if _, err = s.Generate(ctx, society, owner, "2026-10"); err == nil {
		t.Fatal("injected failure ignored")
	}
	count("maintenance_bills", 0)
	count("maintenance_billing_runs", 0)
	count("maintenance_bill_items", 0)
	exec(`DROP TRIGGER reject_maintenance_queue ON maintenance_notification_deliveries`)
	// Race manual callers and a scheduled caller on the same society/month.
	var wg sync.WaitGroup
	results := make(chan models.MaintenanceRunResult, 8)
	failures := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var r models.MaintenanceRunResult
			var err error
			if i%2 == 0 {
				r, err = s.Generate(ctx, society, owner, "2026-10")
			} else {
				r, err = s.generate(ctx, society, 0, "", true)
			}
			results <- r
			failures <- err
		}(i)
	}
	wg.Wait()
	close(results)
	close(failures)
	for err := range failures {
		if err != nil {
			t.Fatal(err)
		}
	}
	created := 0
	var runID int64
	for r := range results {
		created += r.Created
		if runID != 0 && runID != r.RunID {
			t.Fatal("different successful run IDs")
		}
		runID = r.RunID
	}
	if created != 2 {
		t.Fatalf("created=%d", created)
	}
	count("maintenance_bills", 2)
	count("maintenance_bill_items", 4)
	count("maintenance_notification_deliveries", 1)
	list, err := s.List(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: owner, Limit: 1})
	if err != nil || len(list.Items) != 1 || list.NextCursor == nil {
		t.Fatalf("pagination %+v %v", list, err)
	}
	next, err := s.List(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: owner, BeforeID: *list.NextCursor, Limit: 1})
	if err != nil || len(next.Items) != 1 || next.Items[0].ID == list.Items[0].ID {
		t.Fatal("bad second page")
	}
	billID := id(`SELECT id FROM maintenance_bills WHERE flat_id=$1`, flat1)
	b, err := s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: resident, ID: billID, Resident: true})
	if err != nil || b.TotalPaise != 350000 || b.BilledParty != nil || b.Flat.BilledParty != nil {
		t.Fatalf("resident bill %+v %v", b, err)
	}
	adminBill, err := s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: owner, ID: billID})
	if err != nil || len(adminBill.BilledParty) == 0 {
		t.Fatal("admin snapshot missing")
	}
	for _, f := range []models.MaintenanceBillFilter{{SocietyID: otherSociety, UserID: outsider, ID: billID, Resident: true}, {SocietyID: society, UserID: outsider, ID: billID, Resident: true}, {SocietyID: society, UserID: staff, ID: billID}} {
		if _, err = s.Get(ctx, f); err == nil {
			t.Fatalf("unauthorized bill leaked: %+v", f)
		}
	}
	// Delivery errors do not change bills; retries deduplicate the inbox.
	oldLease, err := repo.ClaimDelivery(ctx)
	if err != nil || oldLease == nil {
		t.Fatalf("claim: %v", err)
	}
	exec(`UPDATE maintenance_notification_deliveries SET available_at=now()`)
	newLease, err := repo.ClaimDelivery(ctx)
	if err != nil || newLease == nil {
		t.Fatalf("reclaim: %v", err)
	}
	if err := repo.FinishDelivery(ctx, oldLease, nil); err != nil {
		t.Fatal(err)
	}
	if got := id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE completed_at IS NOT NULL`); got != 0 {
		t.Fatal("expired lease completed a newer delivery")
	}
	if err := repo.FinishDelivery(ctx, newLease, errors.New("retry")); err != nil {
		t.Fatal(err)
	}
	exec(`UPDATE maintenance_notification_deliveries SET available_at=now()`)
	if err = s.Deliver(ctx); err == nil {
		t.Fatal("push failure ignored")
	}
	count("notifications", 1)
	count("maintenance_bills", 2)
	exec(`UPDATE maintenance_notification_deliveries SET available_at=now()`)
	push.fail = false
	if err = s.Deliver(ctx); err != nil {
		t.Fatal(err)
	}
	count("notifications", 1)
	if push.calls != 2 {
		t.Fatalf("push calls=%d", push.calls)
	}
	// New settings and flat edits must not rewrite the issued bill or add a new flat to a completed run.
	settings.FixedPaise = 100000
	if _, err = s.SaveSettings(ctx, society, owner, settings); err != nil {
		t.Fatal(err)
	}
	exec(`UPDATE flats SET flat_number='changed',area_sqft_hundredths=200000 WHERE id=$1`, flat1)
	repeat, err := s.Generate(ctx, society, owner, "2026-10")
	if err != nil || repeat.Existing != 2 || repeat.Created != 0 {
		t.Fatal("repeat changed run")
	}
	b, err = s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: owner, ID: billID})
	if err != nil || b.TotalPaise != 350000 || b.Flat.FlatNumber != "101" {
		t.Fatal("snapshot changed")
	}
	s.now = func() time.Time { return time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC) }
	if _, err = s.Generate(ctx, society, owner, "2026-11"); err != nil {
		t.Fatal(err)
	}
	count("maintenance_bills", 4)
	exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE user_id=$1`, resident)
	if _, err = s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: resident, ID: billID, Resident: true}); err == nil {
		t.Fatal("moved out resident retained access")
	}
	exec(`INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'resident','active')`, society, outsider)
	exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,role,status) VALUES($1,$2,$3,'tenant','active')`, society, flat1, outsider)
	if _, err = s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: outsider, ID: billID, Resident: true}); err != nil {
		t.Fatal("new resident cannot see flat history", err)
	}
	if err = s.Deliver(ctx); err != nil {
		t.Fatal(err)
	}
	if push.calls != 2 {
		t.Fatal("notified a moved-out resident")
	}
	settings.Enabled = false
	if _, err = s.SaveSettings(ctx, society, owner, settings); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Generate(ctx, society, owner, "2026-11"); err == nil {
		t.Fatal("disabled generation allowed")
	}
	if _, err = s.Get(ctx, models.MaintenanceBillFilter{SocietyID: society, UserID: owner, ID: billID}); err != nil {
		t.Fatal("disabled history hidden")
	}
	if _, err = pool.Exec(ctx, `DELETE FROM flats WHERE id=$1`, flat1); err == nil {
		t.Fatal("billed flat hard-deleted")
	}
	for _, query := range []string{`UPDATE maintenance_bills SET total_paise=1 WHERE id=$1`, `DELETE FROM maintenance_bills WHERE id=$1`, `UPDATE maintenance_bill_items SET amount_paise=1 WHERE bill_id=$1`} {
		if _, err := pool.Exec(ctx, query, billID); err == nil {
			t.Fatal("issued financial data mutated")
		}
	}
	// Cross-society reference must be rejected by the database, independently of services.
	if _, err = pool.Exec(ctx, `INSERT INTO maintenance_bills(run_id,society_id,flat_id,billing_month,bill_number,due_date,timezone,total_paise,snapshot) VALUES($1,$2,$3,'2026-10-01','bad','2026-10-10','UTC',1,'{}')`, runID, otherSociety, flat1); err == nil {
		t.Fatal("cross-society foreign key accepted")
	}
	t.Log(fmt.Sprintf("verified atomic billing, concurrent retries, immutable snapshots and resident isolation for society %d", society))
}
