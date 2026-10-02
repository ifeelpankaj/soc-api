//go:build integration

package maintenancesvc

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	notificationsvc "go-server/internal/services/notificationSvc"
	"go-server/pkg/database"
	"sync"
	"testing"
	"time"
)

func TestCatchUpReviewAndOutstandingIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	ctx := f.ctx
	// New synthetic society, configured in September. Every month remains INR 700.
	sid := f.id(`INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES('Carry forward','CARRY','active',$1,$1,now()) RETURNING id`, f.owner)
	f.exec(`INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'owner','active'),($1,$3,'resident','active')`, sid, f.owner, f.resident)
	flat := f.id(`INSERT INTO flats(society_id,flat_number,status,is_active) VALUES($1,'A-101','occupied',true) RETURNING id`, sid)
	f.exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,status,is_primary) VALUES($1,$2,$3,'active',true)`, sid, flat, f.resident)
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.FixedPaise = 70000
	f.billing.now = func() time.Time { return time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC) }
	if _, err := f.billing.Configure(ctx, sid, f.owner, "settings", settings); err != nil {
		t.Fatal(err)
	}
	first, err := f.billing.GenerateCommand(ctx, sid, f.owner, "september", models.MaintenanceMonthRequest{BillingMonth: "2026-09"})
	if err != nil {
		t.Fatal(err)
	}
	retry, err := f.billing.GenerateCommand(ctx, sid, f.owner, "september", models.MaintenanceMonthRequest{BillingMonth: "2026-09"})
	if err != nil || retry != first {
		t.Fatal("lost response replay", retry, err)
	}
	_, err = f.billing.GenerateCommand(ctx, sid, f.owner, "september", models.MaintenanceMonthRequest{BillingMonth: "2026-10"})
	requirePaymentError(t, err, "IDEMPOTENCY_CONFLICT")
	f.billing.now = func() time.Time { return time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC) }
	req := models.MaintenanceMonthRequest{BillingMonth: "2026-10", CatchUp: true, CurrentPricingAcknowledged: true}
	preview, err := f.billing.PreviewCommand(ctx, sid, f.owner, req)
	if err != nil {
		t.Fatal(err)
	}
	if preview.Bills[0].DueDate != "2026-11-25" {
		t.Fatal(preview)
	}
	req.ReviewToken = preview.ReviewToken
	f.exec(`UPDATE flats SET flat_number='A-102' WHERE id=$1`, flat)
	_, err = f.billing.GenerateCommand(ctx, sid, f.owner, "october", req)
	requirePaymentError(t, err, "MAINTENANCE_PREVIEW_CHANGED")
	preview, err = f.billing.PreviewCommand(ctx, sid, f.owner, req)
	if err != nil {
		t.Fatal(err)
	}
	req.ReviewToken = preview.ReviewToken
	if _, err = f.billing.GenerateCommand(ctx, sid, f.owner, "october", req); err != nil {
		t.Fatal(err)
	}
	if _, err = f.billing.GenerateCommand(ctx, sid, f.owner, "november", models.MaintenanceMonthRequest{BillingMonth: "2026-11"}); err != nil {
		t.Fatal(err)
	}
	summary, err := f.billing.Outstanding(ctx, sid, f.resident, flat)
	if err != nil || summary.TotalOutstandingPaise != 210000 || len(summary.UnpaidBills) != 3 {
		t.Fatal(summary, err)
	}
	_, err = f.billing.Outstanding(ctx, sid, f.staff, flat)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
	_, err = f.billing.Outstanding(ctx, f.society, f.owner, flat)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
	bill := summary.UnpaidBills[2].ID
	f.s.now = f.billing.now
	_, err = f.s.SaveSettings(ctx, sid, f.owner, "upi", models.MaintenancePaymentSettings{PaymentMethods: models.MaintenancePaymentMethods{UPI: models.UPIConfig{Enabled: true, UPIID: "society@bank", PayeeName: "Carry forward"}}})
	if err != nil {
		t.Fatal(err)
	}
	credit := models.UPIDirectPayment{BillID: bill, UPIVerifyCredit: models.UPIVerifyCredit{BankCreditConfirmed: true, AmountPaise: 70000, CreditDate: "2026-11-15", Reference: "NOVEMBER700", SettingsVersion: 1, UPIID: "society@bank", EvidenceReference: "statement row 700"}}
	paid, err := f.s.Record(ctx, sid, f.owner, "november-credit", credit)
	if err != nil {
		t.Fatal(err)
	}
	summary, err = f.billing.Outstanding(ctx, sid, f.resident, flat)
	if err != nil || summary.TotalOutstandingPaise != 140000 || summary.CurrentMonthPaise != 0 || summary.PreviousOutstandingPaise != 140000 {
		t.Fatal(summary, err)
	}
	if _, err = f.s.Reverse(ctx, sid, f.owner, paid.ID, "reverse", models.UPIReason{Reason: "Wrong bank match"}); err != nil {
		t.Fatal(err)
	}
	summary, err = f.billing.Outstanding(ctx, sid, f.resident, flat)
	if err != nil || summary.TotalOutstandingPaise != 210000 {
		t.Fatal(summary, err)
	}
	f.exec(`UPDATE flat_residents SET status='inactive' WHERE society_id=$1 AND user_id=$2`, sid, f.resident)
	_, err = f.billing.Outstanding(ctx, sid, f.resident, flat)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
}

func TestPreviewExpiryConcurrentIssueAndOutboxRollback(t *testing.T) {
	f := newPaymentFixture(t)
	ctx := f.ctx
	f.exec(`UPDATE maintenance_settings SET first_enabled_month='2026-09-01' WHERE society_id=$1`, f.society)
	// first_enabled_month was set in September by the synthetic fixture.
	req := models.MaintenanceMonthRequest{BillingMonth: "2026-09", CatchUp: true, CurrentPricingAcknowledged: true}
	p, err := f.billing.PreviewCommand(ctx, f.society, f.owner, req)
	if err != nil {
		t.Fatal(err)
	}
	req.ReviewToken = p.ReviewToken
	f.exec(`UPDATE maintenance_preview_reviews SET expires_at=now()-interval '1 year' WHERE token=$1`, p.ReviewToken)
	_, err = f.billing.GenerateCommand(ctx, f.society, f.owner, "expired", req)
	requirePaymentError(t, err, "MAINTENANCE_PREVIEW_CHANGED")
	p, err = f.billing.PreviewCommand(ctx, f.society, f.owner, req)
	if err != nil {
		t.Fatal(err)
	}
	req.ReviewToken = p.ReviewToken
	f.exec(`CREATE FUNCTION fail_outbox_fixture() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'fixture outbox failure'; END $$; CREATE TRIGGER fail_outbox BEFORE INSERT ON notification_outbox FOR EACH ROW EXECUTE FUNCTION fail_outbox_fixture()`)
	if _, err = f.billing.GenerateCommand(ctx, f.society, f.owner, "concurrent", req); err == nil {
		t.Fatal("outbox failure committed")
	}
	if n := f.id(`SELECT count(*) FROM maintenance_billing_runs WHERE society_id=$1 AND billing_month='2026-09-01'`, f.society); n != 0 {
		t.Fatal("partial issuance")
	}
	f.exec(`DROP TRIGGER fail_outbox ON notification_outbox`)
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := f.billing.GenerateCommand(ctx, f.society, f.owner, "concurrent", req)
			results <- e
		}()
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	if n := f.id(`SELECT count(*) FROM maintenance_billing_runs WHERE society_id=$1 AND billing_month='2026-09-01'`, f.society); n != 1 {
		t.Fatal("duplicate issuance")
	}
}

func TestSharedOutboxLeaseAndInboxDeduplication(t *testing.T) {
	f := newPaymentFixture(t)
	ctx := f.ctx
	d := &database.Database{Pool: f.pool}
	repo := repository.NewNotificationRepository(d)
	notifier, err := notificationsvc.NewNotificationService(ctx, repository.NewDeviceTokenRepository(d), &config.Config{}, repo)
	if err != nil {
		t.Fatal(err)
	}
	worker := notifier.(interface{ DeliverOutbox(context.Context) error })
	q := db.New(f.pool)
	delivery, err := q.ClaimNotificationOutbox(ctx, pgtype.UUID{Bytes: uuid.New(), Valid: true})
	if err != nil {
		t.Fatal(err)
	}
	stale := delivery.LeaseToken
	f.exec(`UPDATE notification_outbox SET available_at=now()-interval '1 second' WHERE id=$1`, delivery.ID)
	delivery, err = q.ClaimNotificationOutbox(ctx, pgtype.UUID{Bytes: uuid.New(), Valid: true})
	if err != nil {
		t.Fatal(err)
	}
	n, err := q.MarkOutboxInbox(ctx, db.MarkOutboxInboxParams{ID: delivery.ID, LeaseToken: stale})
	if err != nil || n != 0 {
		t.Fatal("stale lease accepted", n, err)
	}
	// A crash after inbox creation but before progress marking must not duplicate it.
	_, err = repo.Create(ctx, models.NotificationCreate{ID: uuid.NewString(), UserID: delivery.UserID, SocietyID: &delivery.SocietyID, FlatID: delivery.FlatID, EventKey: &delivery.EventKey, Type: "maintenance_bill_generated", Title: "Existing inbox", Body: "Existing inbox", Data: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	f.exec(`UPDATE notification_outbox SET available_at=now()-interval '1 second' WHERE id=$1`, delivery.ID)
	if err = worker.DeliverOutbox(ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE user_id=$1 AND event_key=$2`, delivery.UserID, delivery.EventKey); n != 1 {
		t.Fatal("duplicate inbox", n)
	}
	if n := f.id(`SELECT count(*) FROM notification_outbox WHERE completed_at IS NULL`); n != 0 {
		t.Fatal("outbox not drained", n)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE user_id=$1 AND event_key=$2 AND completed_at IS NOT NULL`, delivery.UserID, delivery.EventKey); n != 1 {
		t.Fatal("Maintenance event was not acknowledged at materialization", n)
	}
}

func TestEvidenceAndLegacyReplay(t *testing.T) {
	f := newPaymentFixture(t)
	f.settings("upi", "society@bank", true)
	req := models.UPIDirectPayment{BillID: f.bills[0], UPIVerifyCredit: credit("BANK700", "society@bank", 1)}
	req.EvidenceReference = ""
	_, err := f.s.Record(f.ctx, f.society, f.owner, "missing", req)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	req.EvidenceReference = "bank statement"
	paid, err := f.s.Record(f.ctx, f.society, f.owner, "verified", req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.s.Reverse(f.ctx, f.society, f.owner, paid.ID, "reverse", models.UPIReason{Reason: "Wrong bill"})
	if err != nil {
		t.Fatal(err)
	}
	req.BillID = f.bills[1]
	req.ReversalHistoryAcknowledged = true
	req.CorrectionReason = ""
	_, err = f.s.Record(f.ctx, f.society, f.owner, "correct", req)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	req.CorrectionReason = "Move the original credit to the correct flat"
	_, err = f.s.Record(f.ctx, f.society, f.owner, "correct", req)
	if err != nil {
		t.Fatal(err)
	}
}
