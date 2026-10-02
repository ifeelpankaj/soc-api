//go:build integration

package maintenancesvc

import (
	"context"
	"errors"
	"fmt"
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

type queuedReminders struct{}

func (queuedReminders) SendToUser(context.Context, int64, models.NotificationPayload) error {
	return nil
}
func (queuedReminders) DrainOutbox(context.Context) error { return nil }

type reminderInbox struct {
	repository.NotificationRepository
	queries     *db.Queries
	afterCreate func(models.NotificationCreate)
	fail        bool
}

func (r *reminderInbox) OutboxQueries(context.Context) *db.Queries { return r.queries }
func (r *reminderInbox) Enqueue(ctx context.Context, n models.NotificationCreate, audience string) error {
	return r.NotificationRepository.(interface {
		Enqueue(context.Context, models.NotificationCreate, string) error
	}).Enqueue(ctx, n, audience)
}
func (r *reminderInbox) Create(ctx context.Context, n models.NotificationCreate) (*models.Notification, error) {
	if r.fail && n.Type == models.MaintenanceReminderType {
		return nil, errors.New("test inbox failure")
	}
	result, err := r.NotificationRepository.Create(ctx, n)
	if err == nil && r.afterCreate != nil {
		r.afterCreate(n)
	}
	return result, err
}
func reminderWorker(f *paymentFixture, inbox *reminderInbox) interface{ DrainOutbox(context.Context) error } {
	f.t.Helper()
	d := &database.Database{Pool: f.pool}
	notifier, err := notificationsvc.NewNotificationService(f.ctx, repository.NewDeviceTokenRepository(d), &config.Config{}, inbox, allowOperational{})
	if err != nil {
		f.t.Fatal(err)
	}
	return notifier.(interface{ DrainOutbox(context.Context) error })
}
func newReminderInbox(f *paymentFixture) *reminderInbox {
	return &reminderInbox{NotificationRepository: repository.NewNotificationRepository(&database.Database{Pool: f.pool}), queries: db.New(f.pool)}
}
func setReminderDate(f *paymentFixture, day int) {
	f.billing.now = func() time.Time { return time.Date(2026, 10, day, 12, 0, 0, 0, time.UTC) }
	f.billing.push = queuedReminders{}
}

func TestMissingBillsReconcileIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	var original string
	if err := f.pool.QueryRow(f.ctx, `SELECT md5(string_agg(row_to_json(b)::text,'' ORDER BY id)) FROM maintenance_bills b WHERE society_id=$1`, f.society).Scan(&original); err != nil {
		t.Fatal(err)
	}
	flat := f.id(`INSERT INTO flats(society_id,flat_number,status) VALUES($1,'104','occupied') RETURNING id`, f.society)
	f.exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,status) VALUES($1,$2,$3,'active')`, f.society, flat, f.resident)
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.PricingModel = "per_sqft"
	settings.AreaRatePaise = 500
	settings.DueDay = 28
	if _, err := f.billing.SaveSettings(f.ctx, f.society, f.owner, settings); err != nil {
		t.Fatal(err)
	}
	p, err := f.billing.Preview(f.ctx, f.society, f.owner, "2026-10")
	if err != nil || len(p.Bills) != 1 || p.Bills[0].TotalPaise != 123456 || p.Bills[0].DueDate != "2026-10-10" {
		t.Fatalf("preview %+v %v", p, err)
	}
	// Queue failure must roll back the additional bill and its audit record.
	f.exec(`CREATE FUNCTION fail_missing_fixture() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'queue failure'; END $$; CREATE TRIGGER fail_missing BEFORE INSERT ON notification_outbox FOR EACH ROW EXECUTE FUNCTION fail_missing_fixture()`)
	req := models.MaintenanceMonthRequest{BillingMonth: "2026-10"}
	if _, err = f.billing.GenerateCommand(f.ctx, f.society, f.owner, "failed", req); err == nil {
		t.Fatal("queue failure ignored")
	}
	if n := f.id(`SELECT count(*) FROM maintenance_bills WHERE society_id=$1`, f.society); n != 3 {
		t.Fatal("partial additional issuance", n)
	}
	f.exec(`DROP TRIGGER fail_missing ON notification_outbox`)
	var wg sync.WaitGroup
	results := make(chan models.MaintenanceRunResult, 8)
	failures := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var result models.MaintenanceRunResult
			var err error
			if i%2 == 0 {
				result, err = f.billing.GenerateCommand(f.ctx, f.society, f.owner, fmt.Sprint("reconcile-", i), req)
			} else {
				result, err = f.billing.generate(f.ctx, f.society, 0, "", true)
			}
			results <- result
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
	for result := range results {
		created += result.Created
	}
	if created != 1 {
		t.Fatal("duplicate or missing bill", created)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_billing_runs WHERE society_id=$1`, f.society); n != 1 {
		t.Fatal("new monthly run", n)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_bills WHERE society_id=$1 AND flat_id=$2 AND total_paise=123456 AND due_date='2026-10-10'`, f.society, flat); n != 1 {
		t.Fatal("wrong original terms")
	}
	var unchanged string
	if err := f.pool.QueryRow(f.ctx, `SELECT md5(string_agg(row_to_json(b)::text,'' ORDER BY id)) FROM maintenance_bills b WHERE society_id=$1 AND flat_id<>$2`, f.society, flat).Scan(&unchanged); err != nil {
		t.Fatal(err)
	}
	if unchanged != original {
		t.Fatal("existing bill mutated")
	}
	if n := f.id(`SELECT count(*) FROM maintenance_billing_runs WHERE society_id=$1 AND jsonb_array_length(snapshot->'flat_ids')=3`, f.society); n != 1 {
		t.Fatal("original run snapshot changed")
	}
	if n := f.id(`SELECT count(*) FROM maintenance_billing_audit WHERE society_id=$1 AND event_type='billing_reconciled' AND details->'result'->>'created'='1'`, f.society); n != 1 {
		t.Fatal("missing reconciliation audit")
	}
	first, err := f.billing.GenerateCommand(f.ctx, f.society, f.owner, "fresh", req)
	if err != nil || first.Created != 0 || first.Existing != 4 {
		t.Fatal(first, err)
	}
	// Replaying a completed request cannot pick up newly added flats.
	f.exec(`INSERT INTO flats(society_id,flat_number,status) VALUES($1,'105','occupied')`, f.society)
	replay, err := f.billing.GenerateCommand(f.ctx, f.society, f.owner, "fresh", req)
	if err != nil || replay != first {
		t.Fatal("idempotency replay changed", replay, err)
	}
	fresh, err := f.billing.GenerateCommand(f.ctx, f.society, f.owner, "new-action", req)
	if err != nil || fresh.Created != 1 || fresh.Existing != 4 {
		t.Fatal(fresh, err)
	}
}

func TestMaintenanceReminderDatesAndConcurrentTriggersIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	setReminderDate(f, 6)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder'`); n != 0 {
		t.Fatal("wrong date queued", n)
	}
	// Changing settings must not move previously issued bills' reminders.
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.FixedPaise = 999999
	settings.DueDay = 28
	settings.Timezone = "America/New_York"
	if _, err := f.billing.SaveSettings(f.ctx, f.society, f.owner, settings); err != nil {
		t.Fatal(err)
	}
	for index, day := range []int{7, 10, 17} {
		setReminderDate(f, day)
		var wg sync.WaitGroup
		failures := make(chan error, 4)
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); failures <- f.billing.RunReminders(f.ctx) }()
		}
		wg.Wait()
		close(failures)
		for err := range failures {
			if err != nil {
				t.Fatal(err)
			}
		}
		if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder'`); n != int64((index+1)*4) {
			t.Fatalf("day %d count %d", day, n)
		}
	}
	setReminderDate(f, 18)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder'`); n != 12 {
		t.Fatal("late milestone generated", n)
	}
	if n := f.id(`SELECT count(*) FROM notification_outbox WHERE payload->>'Type'='maintenance_payment_reminder' AND payload->'Data'->>'due_date'='2026-10-10' AND payload->'Data'->>'outstanding_amount_paise'='123456'`); n != 12 {
		t.Fatal("bridge lost reminder details", n)
	}
}

func TestMaintenanceReminderPaymentAndAccessRechecksIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	f.settings("upi", "society@bank", true)
	// An unverified claim remains eligible.
	f.claim(f.bills[0], f.request(f.bills[0]), "claim", "PENDINGREMINDER")
	setReminderDate(f, 7)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder'`); n != 4 {
		t.Fatal("pending claim excluded", n)
	}
	if _, err := f.s.Record(f.ctx, f.society, f.owner, "paid-before-delivery", models.UPIDirectPayment{BillID: f.bills[1], UPIVerifyCredit: credit("PAIDREMINDER", "society@bank", 1)}); err != nil {
		t.Fatal(err)
	}
	f.exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE flat_id=(SELECT flat_id FROM maintenance_bills WHERE id=$1)`, f.bills[2])
	inbox := newReminderInbox(f)
	worker := reminderWorker(f, inbox)
	// Initial delivery failure persists a retry without duplicating the event.
	inbox.fail = true
	if err := worker.DrainOutbox(f.ctx); err == nil {
		t.Fatal("expected delivery failure")
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 0 {
		t.Fatal("paid or inactive resident notified", n)
	}
	inbox.fail = false
	f.exec(`UPDATE notification_outbox SET available_at=now() WHERE completed_at IS NULL`)
	setReminderDate(f, 8) // Retry on a later day, without generating a new milestone.
	if err := worker.DrainOutbox(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 2 {
		t.Fatal("retry count", n)
	}
	if n := f.id(`SELECT count(*) FROM notification_outbox WHERE payload->>'Type'='maintenance_payment_reminder' AND push_completed_at IS NOT NULL`); n != 2 {
		t.Fatal("paid or inaccessible reminder sent", n)
	}
	if err := worker.DrainOutbox(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 2 {
		t.Fatal("duplicate inbox", n)
	}
}

func TestMaintenanceReminderPaidBetweenInboxAndPushIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	f.settings("upi", "society@bank", true)
	// Leave only the first bill's residents eligible.
	f.exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE flat_id<>(SELECT flat_id FROM maintenance_bills WHERE id=$1)`, f.bills[0])
	setReminderDate(f, 10)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	inbox := newReminderInbox(f)
	inbox.afterCreate = func(n models.NotificationCreate) {
		if n.Type != models.MaintenanceReminderType {
			return
		}
		inbox.afterCreate = nil
		if _, err := f.s.Record(f.ctx, f.society, f.owner, "paid-between", models.UPIDirectPayment{BillID: f.bills[0], UPIVerifyCredit: credit("BETWEENREMINDER", "society@bank", 1)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := reminderWorker(f, inbox).DrainOutbox(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 1 {
		t.Fatal("expected only pre-payment inbox", n)
	}
	if n := f.id(`SELECT count(*) FROM notification_outbox WHERE payload->>'Type'='maintenance_payment_reminder' AND push_completed_at IS NOT NULL`); n != 0 {
		t.Fatal("paid bill pushed", n)
	}
}

func TestMaintenanceReminderMultipleBatchesAndDisabledIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	f.exec(`INSERT INTO flats(society_id,flat_number,status) SELECT $1,'bulk-'||n,'occupied' FROM generate_series(1,501) n`, f.society)
	f.exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,status) SELECT society_id,id,$2,'active' FROM flats WHERE society_id=$1 AND flat_number LIKE 'bulk-%'`, f.society, f.resident)
	if _, err := f.billing.Generate(f.ctx, f.society, f.owner, "2026-10"); err != nil {
		t.Fatal(err)
	}
	setReminderDate(f, 17)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder'`); n != 505 {
		t.Fatal("candidate paging lost bills", n)
	}
	worker := reminderWorker(f, newReminderInbox(f))
	if err := worker.DrainOutbox(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notification_outbox WHERE completed_at IS NULL`); n != 0 {
		t.Fatal("outbox batches not drained", n)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 505 {
		t.Fatal("notification batches incomplete", n)
	}
	// Queue another stage, then disable billing before delivery.
	setReminderDate(f, 10)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	f.exec(`UPDATE maintenance_settings SET enabled=false WHERE society_id=$1`, f.society)
	if err := worker.DrainOutbox(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE type='maintenance_payment_reminder'`); n != 505 {
		t.Fatal("disabled society notified", n)
	}
	setReminderDate(f, 7)
	if err := f.billing.RunReminders(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_notification_deliveries WHERE event_type='maintenance_payment_reminder' AND event_data->>'reminder_milestone'='before_3'`); n != 0 {
		t.Fatal("disabled society queued", n)
	}
}
