//go:build integration

package maintenancesvc

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/db"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/testutil"
	"go-server/pkg/database"
	"image/png"
	"sync"
	"testing"
	"time"
)

type paymentFixture struct {
	t                                                               *testing.T
	ctx                                                             context.Context
	pool                                                            *pgxpool.Pool
	s                                                               *PaymentService
	billing                                                         *Service
	push                                                            *integrationPush
	society, otherSociety, owner, resident, second, staff, outsider int64
	bills                                                           []int64
}

func (f *paymentFixture) exec(query string, args ...any) {
	f.t.Helper()
	if _, err := f.pool.Exec(f.ctx, query, args...); err != nil {
		f.t.Fatal(err)
	}
}
func (f *paymentFixture) id(query string, args ...any) int64 {
	f.t.Helper()
	var id int64
	if err := f.pool.QueryRow(f.ctx, query, args...).Scan(&id); err != nil {
		f.t.Fatal(err)
	}
	return id
}
func (f *paymentFixture) balance(bill, paid int64) {
	f.t.Helper()
	b, err := f.billing.Get(f.ctx, models.MaintenanceBillFilter{SocietyID: f.society, UserID: f.owner, ID: bill})
	if err != nil || b.PaidAmountPaise != paid || b.OutstandingAmountPaise != 123456-paid {
		f.t.Fatalf("balance %+v %v", b, err)
	}
	if paid > 0 && b.Status != "paid" {
		f.t.Fatal("verified payment did not settle bill")
	}
}
func newPaymentFixture(t *testing.T) *paymentFixture {
	return newPaymentFixtureAt(t, 0)
}
func newPaymentFixtureAt(t *testing.T, maxMigrations int) *paymentFixture {
	t.Helper()
	ctx := context.Background()
	pg := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	sqlDB, err := sql.Open("pgx", pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if _, err = migrate.ExecMax(sqlDB, "postgres", &migrate.FileMigrationSource{Dir: "../../../migrations"}, migrate.Up, maxMigrations); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	f := &paymentFixture{t: t, ctx: ctx, pool: pool, push: &integrationPush{}}
	user := func(name string) int64 {
		return f.id(`INSERT INTO users(full_name,email,email_verified) VALUES($1,$2,true) RETURNING id`, name, name+"@example.com")
	}
	f.owner = user("owner")
	f.resident = user("resident")
	f.second = user("second")
	f.staff = user("staff")
	f.outsider = user("outsider")
	society := func(name string, owner int64) int64 {
		return f.id(`INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES($1,$1,'active',$2,$2,now()) RETURNING id`, name, owner)
	}
	f.society = society("PAYMENT", f.owner)
	f.otherSociety = society("OTHER", f.outsider)
	f.exec(`INSERT INTO society_members(society_id,user_id,role,status) VALUES($1,$2,'owner','active'),($1,$3,'resident','active'),($1,$4,'resident','active'),($1,$5,'staff','active'),($6,$7,'owner','active')`, f.society, f.owner, f.resident, f.second, f.staff, f.otherSociety, f.outsider)
	for i := 0; i < 3; i++ {
		flat := f.id(`INSERT INTO flats(society_id,flat_number,status) VALUES($1,$2,'occupied') RETURNING id`, f.society, fmt.Sprint(i+101))
		f.exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,status,is_primary) VALUES($1,$2,$3,'active',true)`, f.society, flat, f.resident)
		if i == 0 {
			f.exec(`INSERT INTO flat_residents(society_id,flat_id,user_id,status) VALUES($1,$2,$3,'active')`, f.society, flat, f.second)
		}
	}
	d := &database.Database{Pool: pool}
	tx := repository.NewTransactionManager(d)
	var billingStore Store = repository.NewMaintenanceRepository(d, tx)
	if maxMigrations > 0 && maxMigrations < 29 {
		billingStore = legacyBillingStore{billingStore, db.New(pool)}
	}
	f.billing = New(billingStore, allowOperational{}, []string{"occupied"}, repository.NewNotificationRepository(d), f.push)
	f.s = NewPaymentService(repository.NewMaintenancePaymentRepository(d, tx), allowOperational{})
	now := func() time.Time { return time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC) }
	f.s.now = now
	f.billing.now = now
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.FixedPaise = 123456
	if _, err = f.billing.SaveSettings(ctx, f.society, f.owner, settings); err != nil {
		t.Fatal(err)
	}
	if _, err = f.billing.Generate(ctx, f.society, f.owner, "2026-10"); err != nil {
		t.Fatal(err)
	}
	rows, err := pool.Query(ctx, `SELECT id FROM maintenance_bills WHERE society_id=$1 ORDER BY flat_id`, f.society)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		f.bills = append(f.bills, id)
	}
	rows.Close()
	return f
}
func (f *paymentFixture) settings(key, upi string, enabled bool) models.MaintenancePaymentSettings {
	f.t.Helper()
	v, err := f.s.SaveSettings(f.ctx, f.society, f.owner, key, models.MaintenancePaymentSettings{PaymentMethods: models.MaintenancePaymentMethods{UPI: models.UPIConfig{Enabled: enabled, UPIID: upi, PayeeName: "A & B Society"}}})
	if err != nil {
		f.t.Fatal(err)
	}
	return v
}
func (f *paymentFixture) request(bill int64) models.UPIPaymentRequest {
	f.t.Helper()
	r, err := f.s.Request(f.ctx, f.society, f.resident, bill)
	if err != nil {
		f.t.Fatal(err)
	}
	return r
}
func (f *paymentFixture) claim(bill int64, r models.UPIPaymentRequest, key, reference string) models.UPIClaim {
	f.t.Helper()
	c, err := f.s.SubmitClaim(f.ctx, f.society, f.resident, bill, key, models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: reference, PaymentDate: "2026-10-14"})
	if err != nil {
		f.t.Fatal(err)
	}
	return c
}
func credit(reference, upi string, version int64) models.UPIVerifyCredit {
	return models.UPIVerifyCredit{EvidenceReference: "Synthetic bank statement row", CorrectionReason: "Correct the original bill allocation", BankCreditConfirmed: true, AmountPaise: 123456, CreditDate: "2026-10-14", Reference: reference, UPIID: upi, SettingsVersion: version}
}
func requirePaymentError(t *testing.T, err error, code string) {
	t.Helper()
	var app *models.AppError
	if !errors.As(err, &app) || app.Code != code {
		t.Fatalf("error=%v want %s", err, code)
	}
}

func TestUPIPaymentLifecycleIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	ctx := f.ctx
	s := f.s
	bill := f.bills[0]
	disabled, err := s.Settings(ctx, f.society, f.owner)
	if err != nil || disabled.PaymentMethods.UPI.Enabled {
		t.Fatal("UPI not disabled by default", err)
	}
	_, err = s.Request(ctx, f.society, f.resident, bill)
	requirePaymentError(t, err, "UPI_DISABLED")
	v1 := f.settings("settings-1", "old@bank", true)
	_, err = s.SaveSettings(ctx, f.society, f.staff, "staff-settings", v1)
	requirePaymentError(t, err, "PAYMENT_FORBIDDEN")
	r1 := f.request(bill)
	again := f.request(bill)
	if r1.ID != again.ID || !containsAmount(r1.URI, "1234.56") {
		t.Fatal("request not reused/exact")
	}
	image, err := s.QR(ctx, f.society, f.resident, r1.ID)
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(image))
	if err != nil || img.Bounds().Dx() != 384 {
		t.Fatal("bad QR PNG", err)
	}
	red, green, blue, _ := img.At(0, 0).RGBA()
	if red != 65535 || green != 65535 || blue != 65535 {
		t.Fatal("QR white border missing")
	}
	_, err = s.QR(ctx, f.society, f.outsider, r1.ID)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	_, err = s.QR(ctx, f.otherSociety, f.outsider, r1.ID)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	v2 := f.settings("settings-2", "new@bank", true)
	if v2.Version != v1.Version+1 {
		t.Fatal("version not incremented")
	}
	_, err = s.QR(ctx, f.society, f.resident, r1.ID)
	requirePaymentError(t, err, "PAYMENT_REQUEST_OBSOLETE")
	r2 := f.request(bill)
	if r1.ID == r2.ID {
		t.Fatal("superseded request reused")
	}
	c := f.claim(bill, r1, "claim-1", " 000abc123456 ")
	if c.Reference != "000ABC123456" || c.SettingsVersion != v1.Version {
		t.Fatal("historical destination lost")
	}
	f.balance(bill, 0)
	retry, err := s.SubmitClaim(ctx, f.society, f.resident, bill, "claim-1", models.UPISubmitClaim{PaymentRequestID: r1.ID, Reference: "000ABC123456", PaymentDate: "2026-10-14"})
	if err != nil || retry.ID != c.ID {
		t.Fatal("claim idempotency", err)
	}
	_, err = s.SubmitClaim(ctx, f.society, f.resident, bill, "claim-1", models.UPISubmitClaim{PaymentRequestID: r1.ID, Reference: "000ABC999999", PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "IDEMPOTENCY_CONFLICT")
	_, err = s.QR(ctx, f.society, f.resident, r2.ID)
	requirePaymentError(t, err, "PAYMENT_CLAIM_PENDING")
	_, err = s.Request(ctx, f.society, f.second, bill)
	requirePaymentError(t, err, "PAYMENT_CLAIM_PENDING")
	bad := credit(c.Reference, "new@bank", v2.Version)
	_, err = s.Verify(ctx, f.society, f.owner, c.ID, "wrong-destination", bad)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	good := credit(c.Reference, "old@bank", v1.Version)
	bad = good
	bad.AmountPaise--
	_, err = s.Verify(ctx, f.society, f.owner, c.ID, "bad-amount", bad)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	bad = good
	bad.BankCreditConfirmed = false
	_, err = s.Verify(ctx, f.society, f.owner, c.ID, "unconfirmed", bad)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	// Disable new requests, while permitting verification of the historical credit.
	f.settings("settings-disabled", "new@bank", false)
	p, err := s.Verify(ctx, f.society, f.owner, c.ID, "verify-1", good)
	if err != nil {
		t.Fatal(err)
	}
	f.balance(bill, 123456)
	repeat, err := s.Verify(ctx, f.society, f.owner, c.ID, "verify-1", good)
	if err != nil || repeat.ID != p.ID {
		t.Fatal("verification replay", err)
	}
	if n := f.id(`SELECT count(*) FROM maintenance_payment_ledger WHERE payment_id=$1`, p.ID); n != 1 {
		t.Fatal("duplicate collection ledger")
	}
	_, err = s.SubmitClaim(ctx, f.society, f.second, bill, "extra-claim", models.UPISubmitClaim{PaymentRequestID: r2.ID, Reference: "DUPLICATE001", PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "BILL_ALREADY_PAID")
	reportReq := models.UPICreateReport{BillID: bill, PaymentRequestID: r2.ID, Reference: "DUPLICATE001", PaymentDate: "2026-10-14", AmountPaise: 123456, Explanation: "Second resident transferred again"}
	report, err := s.CreateReport(ctx, f.society, f.second, "report-1", reportReq, false)
	if err != nil {
		t.Fatal(err)
	}
	reportAgain, err := s.CreateReport(ctx, f.society, f.second, "report-retry-other-key", reportReq, false)
	if err != nil || reportAgain.ID != report.ID {
		t.Fatal("report dedupe", err)
	}
	f.balance(bill, 123456)
	if n := f.id(`SELECT count(*) FROM maintenance_payment_reference_reservations WHERE reference='DUPLICATE001'`); n != 0 {
		t.Fatal("report reserved reference")
	}
	if _, err = s.UpdateReport(ctx, f.society, f.owner, report.ID, "investigate", models.UPIUpdateReport{Status: "investigating"}); err != nil {
		t.Fatal(err)
	}
	if _, err = s.UpdateReport(ctx, f.society, f.owner, report.ID, "resolve", models.UPIUpdateReport{Status: "resolved", ResolutionNote: "Society is resolving the extra transfer with its bank"}); err != nil {
		t.Fatal(err)
	}
	f.balance(bill, 123456)
	own, err := s.Payment(ctx, f.society, f.resident, p.ID, true)
	if err != nil || own.Reference == nil {
		t.Fatal("payer receipt", err)
	}
	other, err := s.Payment(ctx, f.society, f.second, p.ID, true)
	if err != nil || other.Reference != nil || other.PayerID != nil || other.VerifiedBy != nil {
		t.Fatal("receipt privacy", err)
	}
	_, err = s.Payments(ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.second, Resident: true, Reference: c.Reference})
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	reversed, err := s.Reverse(ctx, f.society, f.owner, p.ID, "reverse-1", models.UPIReason{Reason: "Incorrect bank credit match"})
	if err != nil || reversed.Status != "reversed" || reversed.ReceiptNumber != p.ReceiptNumber {
		t.Fatal("reversal", err)
	}
	f.balance(bill, 0)
	_, err = s.Reverse(ctx, f.society, f.owner, p.ID, "reverse-again", models.UPIReason{Reason: "Again"})
	requirePaymentError(t, err, "PAYMENT_ALREADY_REVERSED")
	if n := f.id(`SELECT sum(amount_paise)::bigint FROM maintenance_payment_ledger WHERE payment_id=$1`, p.ID); n != 0 {
		t.Fatal("unbalanced reversal")
	}
	f.settings("settings-enable", "new@bank", true)
	fresh := f.request(bill)
	if fresh.ID == r2.ID || fresh.Reference == r1.Reference {
		t.Fatal("old QR reactivated after reversal")
	}
	freshClaim := f.claim(bill, fresh, "fresh-claim", c.Reference)
	freshCredit := credit(c.Reference, "new@bank", fresh.SettingsVersion)
	_, err = s.Verify(ctx, f.society, f.owner, freshClaim.ID, "fresh-verify", freshCredit)
	requirePaymentError(t, err, "REVERSAL_HISTORY_ACKNOWLEDGMENT_REQUIRED")
	history, err := s.Claim(ctx, f.society, f.owner, freshClaim.ID)
	if err != nil || len(history.ReferenceHistory.Payments.Items) != 1 || history.ReferenceHistory.Payments.Items[0].Status != "reversed" {
		t.Fatal("reversal history missing", err)
	}
	freshCredit.ReversalHistoryAcknowledged = true
	if _, err = s.Verify(ctx, f.society, f.owner, freshClaim.ID, "fresh-verify", freshCredit); err != nil {
		t.Fatal(err)
	}
	f.balance(bill, 123456)
	summary, err := s.Summary(ctx, f.society, f.owner, "2026-10")
	if err != nil || summary.BilledPaise != 370368 || summary.CollectedPaise != 123456 || summary.OutstandingPaise != 246912 {
		t.Fatalf("summary %+v %v", summary, err)
	}
	// A delivery provider failure cannot roll back the verified collection.
	f.push.fail = true
	if err = f.billing.Deliver(ctx); err == nil {
		t.Fatal("expected push failure")
	}
	f.balance(bill, 123456)
	f.exec(`UPDATE maintenance_notification_deliveries SET available_at=now() WHERE completed_at IS NULL`)
	f.push.fail = false
	if err = f.billing.Deliver(ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE user_id=$1 AND event_key=$2`, f.resident, fmt.Sprintf("maintenance_payment_verified:%d", p.ID)); n != 1 {
		t.Fatal("duplicate notification")
	}
	// Database guards preserve original records, not only service-level conventions.
	for _, query := range []string{`UPDATE maintenance_payments SET amount_paise=1 WHERE id=$1`, `DELETE FROM maintenance_payments WHERE id=$1`, `UPDATE maintenance_payment_ledger SET amount_paise=1 WHERE payment_id=$1`} {
		if _, err = f.pool.Exec(ctx, query, p.ID); err == nil {
			t.Fatal("financial mutation accepted")
		}
	}
	f.exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE user_id=$1`, f.resident)
	_, err = s.Payment(ctx, f.society, f.resident, p.ID, true)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	_, err = s.SubmitClaim(ctx, f.society, f.resident, bill, "claim-1", models.UPISubmitClaim{PaymentRequestID: r1.ID, Reference: c.Reference, PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
}

func TestUPIConcurrentReservationsAndAtomicity(t *testing.T) {
	f := newPaymentFixture(t)
	s := f.s
	ctx := f.ctx
	f.settings("settings", "society@bank", true)
	a, b := f.bills[0], f.bills[1]
	ra, rb := f.request(a), f.request(b)
	// Same reference across different bills/requests: exactly one pending claim.
	var wg sync.WaitGroup
	results := make(chan models.UPIClaim, 2)
	errs := make(chan error, 2)
	for i, pair := range []struct {
		bill int64
		r    models.UPIPaymentRequest
	}{{a, ra}, {b, rb}} {
		wg.Add(1)
		go func(i int, bill int64, r models.UPIPaymentRequest) {
			defer wg.Done()
			c, e := s.SubmitClaim(ctx, f.society, f.resident, bill, fmt.Sprintf("race-%d", i), models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: "000DUPLICATE", PaymentDate: "2026-10-14"})
			results <- c
			errs <- e
		}(i, pair.bill, pair.r)
	}
	wg.Wait()
	close(results)
	close(errs)
	successes := 0
	for err := range errs {
		if err == nil {
			successes++
		} else {
			requirePaymentError(t, err, "PAYMENT_CONFLICT")
		}
	}
	if successes != 1 {
		t.Fatalf("claim successes=%d", successes)
	}
	var pending models.UPIClaim
	for c := range results {
		if c.ID > 0 {
			pending = c
		}
	}
	otherBill := a
	if pending.BillID == a {
		otherBill = b
	}
	_, err := s.Record(ctx, f.society, f.owner, "direct-collision", models.UPIDirectPayment{BillID: otherBill, UPIVerifyCredit: credit(pending.Reference, "society@bank", 1)})
	requirePaymentError(t, err, "PAYMENT_CONFLICT")
	if n := f.id(`SELECT count(*) FROM maintenance_payments`); n != 0 {
		t.Fatal("failed reservation left payment")
	}
	// Cancellation frees the shared reservation; direct credit can now use it.
	if _, err = s.CloseClaim(ctx, f.society, f.resident, pending.ID, "cancel", models.UPIReason{}, false); err != nil {
		t.Fatal(err)
	}
	paid, err := s.Record(ctx, f.society, f.owner, "direct", models.UPIDirectPayment{BillID: otherBill, UPIVerifyCredit: credit(pending.Reference, "society@bank", 1)})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.SubmitClaim(ctx, f.society, f.resident, pending.BillID, "against-paid-ref", models.UPISubmitClaim{PaymentRequestID: pending.PaymentRequestID, Reference: pending.Reference, PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "PAYMENT_CONFLICT")
	if _, err = s.Reverse(ctx, f.society, f.owner, paid.ID, "undo-direct", models.UPIReason{Reason: "Test correction"}); err != nil {
		t.Fatal(err)
	}
	// Inject failures at each late transaction stage and verify full rollback.
	req := f.request(f.bills[2])
	c := f.claim(f.bills[2], req, "atomic-claim", "000ATOMIC123")
	f.exec(`CREATE FUNCTION reject_payment_write() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected payment failure'; END $$`)
	for _, table := range []string{"maintenance_payment_ledger", "maintenance_payment_audit_events", "maintenance_notification_deliveries"} {
		f.exec(fmt.Sprintf(`CREATE TRIGGER reject_payment_write BEFORE INSERT ON %s FOR EACH ROW EXECUTE FUNCTION reject_payment_write()`, table))
		if _, err = s.Verify(ctx, f.society, f.owner, c.ID, "atomic-verify", credit(c.Reference, "society@bank", 1)); err == nil {
			t.Fatal("injected failure ignored", table)
		}
		f.balance(c.BillID, 0)
		if n := f.id(`SELECT count(*) FROM maintenance_payment_claims WHERE id=$1 AND status='pending'`, c.ID); n != 1 {
			t.Fatal("claim update escaped rollback")
		}
		if n := f.id(`SELECT count(*) FROM maintenance_payment_reference_reservations WHERE claim_id=$1`, c.ID); n != 1 {
			t.Fatal("reference released on rollback")
		}
		f.exec(fmt.Sprintf(`DROP TRIGGER reject_payment_write ON %s`, table))
	}
	// Concurrent approvals produce only one settlement and one collection entry.
	errs = make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Verify(ctx, f.society, f.owner, c.ID, fmt.Sprintf("approve-%d", i), credit(c.Reference, "society@bank", 1))
			errs <- err
		}(i)
	}
	wg.Wait()
	close(errs)
	successes = 0
	for err := range errs {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("approval successes=%d", successes)
	}
	f.balance(c.BillID, 123456)
	if n := f.id(`SELECT count(*) FROM maintenance_payment_ledger WHERE bill_id=$1 AND kind='collection'`, c.BillID); n != 1 {
		t.Fatal("duplicate ledger")
	}
	// Direct recording must explicitly resolve a resident's pending claim.
	r := f.request(pending.BillID)
	c2 := f.claim(pending.BillID, r, "resolve-claim", "000RESOLVE12")
	direct := models.UPIDirectPayment{BillID: pending.BillID, UPIVerifyCredit: credit("000OTHER1234", "society@bank", 1)}
	_, err = s.Record(ctx, f.society, f.owner, "resolve-direct", direct)
	requirePaymentError(t, err, "PENDING_CLAIM_RESOLUTION_REQUIRED")
	direct.PendingClaimID = c2.ID
	direct.ClaimResolution = "reject"
	direct.Reason = "Bank credit matches a different transfer"
	if _, err = s.Record(ctx, f.society, f.owner, "resolve-direct", direct); err != nil {
		t.Fatal(err)
	}
	review, err := s.Claim(ctx, f.society, f.owner, c2.ID)
	if err != nil || review.Status != "rejected" {
		t.Fatal("claim not resolved", err)
	}
}

func TestUPIBillRacesAndIsolation(t *testing.T) {
	f := newPaymentFixture(t)
	s, ctx := f.s, f.ctx
	v := f.settings("settings", "society@bank", true)
	// Settings retries cannot accidentally create versions or overwrite newer input.
	replay, err := s.SaveSettings(ctx, f.society, f.owner, "settings", v)
	if err != nil || replay.Version != v.Version {
		t.Fatal("settings replay", err)
	}
	changed := v
	changed.PaymentMethods.UPI.UPIID = "different@bank"
	_, err = s.SaveSettings(ctx, f.society, f.owner, "settings", changed)
	requirePaymentError(t, err, "IDEMPOTENCY_CONFLICT")
	bill := f.bills[0]
	r := f.request(bill)
	_, err = s.SubmitClaim(ctx, f.society, f.resident, bill, "future", models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: "FUTURE0001", PaymentDate: "2026-10-16"})
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	_, err = s.SubmitClaim(ctx, f.society, f.resident, bill, "", models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: "NOKEY00001", PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	// Two flat residents submit different references simultaneously: only one wins.
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for i, user := range []int64{f.resident, f.second} {
		wg.Add(1)
		go func(i int, user int64) {
			defer wg.Done()
			_, e := s.SubmitClaim(ctx, f.society, user, bill, fmt.Sprintf("same-bill-%d", i), models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: fmt.Sprintf("SAMERACE00%d", i), PaymentDate: "2026-10-14"})
			errs <- e
		}(i, user)
	}
	wg.Wait()
	close(errs)
	wins := 0
	for e := range errs {
		if e == nil {
			wins++
		} else {
			requirePaymentError(t, e, "PAYMENT_CLAIM_PENDING")
		}
	}
	if wins != 1 {
		t.Fatal("multiple claims accepted")
	}
	cID := f.id(`SELECT id FROM maintenance_payment_claims WHERE bill_id=$1 AND status='pending'`, bill)
	detail, err := s.Claim(ctx, f.society, f.owner, cID)
	if err != nil {
		t.Fatal(err)
	}
	// Cancellation and approval race; the reservation follows exactly the winning state.
	errs = make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, e := s.CloseClaim(ctx, f.society, detail.SubmittedBy, cID, "cancel-race", models.UPIReason{}, false)
		errs <- e
	}()
	go func() {
		defer wg.Done()
		_, e := s.Verify(ctx, f.society, f.owner, cID, "approve-race", credit(detail.Reference, "society@bank", 1))
		errs <- e
	}()
	wg.Wait()
	close(errs)
	wins = 0
	for e := range errs {
		if e == nil {
			wins++
		} else {
			requirePaymentError(t, e, "CLAIM_ALREADY_REVIEWED")
		}
	}
	if wins != 1 {
		t.Fatal("cancel/verify both won")
	}
	active := f.id(`SELECT count(*) FROM maintenance_payments WHERE bill_id=$1 AND status='verified'`, bill)
	reserved := f.id(`SELECT count(*) FROM maintenance_payment_reference_reservations WHERE bill_id=$1`, bill)
	if active != reserved {
		t.Fatal("orphan reference after cancellation race")
	}
	f.balance(bill, active*123456)
	// Direct recording and claim approval use the same serialization boundary.
	bill = f.bills[1]
	r = f.request(bill)
	c := f.claim(bill, r, "direct-race-claim", "DIRECTRACE01")
	errs = make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, e := s.Record(ctx, f.society, f.owner, "direct-race", models.UPIDirectPayment{BillID: bill, PendingClaimID: c.ID, ClaimResolution: "verify", UPIVerifyCredit: credit(c.Reference, "society@bank", 1)})
		errs <- e
	}()
	go func() {
		defer wg.Done()
		_, e := s.Verify(ctx, f.society, f.owner, c.ID, "verify-direct-race", credit(c.Reference, "society@bank", 1))
		errs <- e
	}()
	wg.Wait()
	close(errs)
	wins = 0
	for e := range errs {
		if e == nil {
			wins++
		}
	}
	if wins != 1 {
		t.Fatal("direct/verify both won")
	}
	f.balance(bill, 123456)
	pid := f.id(`SELECT id FROM maintenance_payments WHERE bill_id=$1 AND status='verified'`, bill)
	_, err = s.Payment(ctx, f.otherSociety, f.outsider, pid, false)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	_, err = s.CreateReport(ctx, f.society, f.second, "wrong-flat", models.UPICreateReport{BillID: bill, Reference: "REPORT0001", AmountPaise: 123456, PaymentDate: "2026-10-14", Explanation: "Wrong flat"}, false)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	// Explicit rejection releases the reference and permits a fresh claim.
	third := f.bills[2]
	thirdRequest := f.request(third)
	rejected := f.claim(third, thirdRequest, "reject-claim", "REJECT00001")
	_, err = s.CloseClaim(ctx, f.society, f.owner, rejected.ID, "reject", models.UPIReason{Reason: "No matching bank credit"}, true)
	if err != nil {
		t.Fatal(err)
	}
	fresh := f.claim(third, thirdRequest, "resubmit", rejected.Reference)
	if fresh.ID == rejected.ID {
		t.Fatal("rejection history overwritten")
	}
	f.balance(third, 0)
	// Tenant-consistent FKs prevent linking another society's bill to this config.
	if _, err = f.pool.Exec(ctx, `INSERT INTO maintenance_payment_requests(id,society_id,bill_id,settings_version,amount_paise,reference,created_by) VALUES(gen_random_uuid(),$1,$2,1,123456,'cross-tenant',$3)`, f.otherSociety, third, f.outsider); err == nil {
		t.Fatal("cross-tenant foreign key accepted")
	}
	// Two reversal actions post only one compensation.
	errs = make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, e := s.Reverse(ctx, f.society, f.owner, pid, fmt.Sprintf("reverse-race-%d", i), models.UPIReason{Reason: "Correct bank match"})
			errs <- e
		}(i)
	}
	wg.Wait()
	close(errs)
	wins = 0
	for e := range errs {
		if e == nil {
			wins++
		} else {
			requirePaymentError(t, e, "PAYMENT_ALREADY_REVERSED")
		}
	}
	if wins != 1 {
		t.Fatal("double reversal")
	}
	f.balance(bill, 0)
	if n := f.id(`SELECT count(*) FROM maintenance_payment_ledger WHERE payment_id=$1 AND kind='reversal'`, pid); n != 1 {
		t.Fatal("duplicate compensation")
	}
	_, err = s.SubmitClaim(ctx, f.society, f.resident, bill, "closed-request", models.UPISubmitClaim{PaymentRequestID: r.ID, Reference: "CLOSED00001", PaymentDate: "2026-10-14"})
	requirePaymentError(t, err, "PAYMENT_REQUEST_CLOSED")
	// Removed members are also excluded from deferred payment notifications.
	f.exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE user_id=$1`, f.resident)
	if err = f.billing.Deliver(ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.id(`SELECT count(*) FROM notifications WHERE user_id=$1`, f.resident); n != 0 {
		t.Fatal("notification leaked to moved-out resident")
	}
	page, err := s.Payments(ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.resident, Resident: true})
	if err != nil || len(page.Items) != 0 {
		t.Fatal("moved-out history visible", err)
	}
}

// Legacy schema adapter is limited to the migration fixture; production requires all migrations.
type legacyBillingStore struct {
	Store
	q *db.Queries
}

func (s legacyBillingStore) SaveSettings(ctx context.Context, id, user int64, v models.MaintenanceSettings) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.q.SaveMaintenanceSettings(ctx, db.SaveMaintenanceSettingsParams{SocietyID: id, Enabled: v.Enabled, Config: b, UpdatedBy: &user})
}
func (s legacyBillingStore) Settings(ctx context.Context, id int64) (models.MaintenanceSettings, error) {
	var v models.MaintenanceSettings
	b, err := s.q.GetMaintenanceSettings(ctx, id)
	if err != nil {
		return models.DefaultMaintenanceSettings(), nil
	}
	err = json.Unmarshal(b, &v)
	return v, err
}
