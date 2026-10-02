//go:build integration

package maintenancesvc

import (
	"bytes"
	"database/sql"
	"encoding/json"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/models"
	"go-server/internal/services/maintenanceSvc/documentpdf"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func saveIntegrationPDF(t *testing.T, name string, p models.MaintenancePDF) {
	t.Helper()
	if !bytes.HasPrefix(p.Bytes, []byte("%PDF-")) {
		t.Fatal("not a PDF")
	}
	if dir := os.Getenv("MAINTENANCE_PDF_QA_DIR"); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), p.Bytes, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
func TestMaintenanceDocumentsMigrationAndAccess(t *testing.T) {
	f := newPaymentFixtureAt(t, 26)
	ctx := f.ctx
	bill := f.bills[0]
	// Existing application instances issue bills without any issuer column.
	f.exec(`UPDATE societies SET name='Original Society',address_line1='12 Frozen Road',city='Pune',email='office@example.com' WHERE id=$1`, f.society)
	var before []byte
	if err := f.pool.QueryRow(ctx, `SELECT to_jsonb(b) FROM maintenance_bills b WHERE id=$1`, bill).Scan(&before); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := sql.Open("pgx", f.pool.Config().ConnConfig.ConnString())
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	source := &migrate.FileMigrationSource{Dir: "../../../migrations"}
	migrations, err := source.FindMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if n, err := migrate.Exec(sqlDB, "postgres", source, migrate.Up); err != nil || n != len(migrations)-26 {
		t.Fatalf("migration: %d %v", n, err)
	}
	// The upgraded schema uses the current repository capabilities, including
	// outstanding summaries; the legacy wrapper is only for pre-upgrade billing.
	if legacy, ok := f.billing.store.(legacyBillingStore); ok {
		f.billing.store = legacy.Store
	}
	var after []byte
	if err = f.pool.QueryRow(ctx, `SELECT to_jsonb(b)-'issuer_snapshot' FROM maintenance_bills b WHERE id=$1`, bill).Scan(&after); err != nil {
		t.Fatal(err)
	}
	var oldRow, newRow map[string]any
	_ = json.Unmarshal(before, &oldRow)
	_ = json.Unmarshal(after, &newRow)
	a, _ := json.Marshal(oldRow)
	b, _ := json.Marshal(newRow)
	if !bytes.Equal(a, b) {
		t.Fatal("migration changed original bill fields")
	}
	docs := NewDocumentService(f.billing, f.s, documentpdf.New(), allowOperational{})
	inv, err := docs.Invoice(ctx, f.society, f.resident, bill, true)
	if err != nil {
		t.Fatal(err)
	}
	saveIntegrationPDF(t, "integration-invoice.pdf", inv)
	loaded, err := f.billing.Get(ctx, models.MaintenanceBillFilter{SocietyID: f.society, UserID: f.owner, ID: bill})
	if err != nil || loaded.Issuer.CaptureSource != "rollout" || loaded.Issuer.Name != "Original Society" {
		t.Fatal("legacy snapshot", loaded.Issuer, err)
	}
	_, err = docs.Invoice(ctx, f.society, f.staff, bill, false)
	requirePaymentError(t, err, "MAINTENANCE_FORBIDDEN")
	_, err = docs.Invoice(ctx, f.otherSociety, f.outsider, bill, false)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
	_, err = docs.Invoice(ctx, f.society, f.outsider, bill, true)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
	// Pending claim never yields a receipt. UPI availability does not gate PDFs.
	f.settings("documents-settings", "old@bank", true)
	request := f.request(bill)
	claim := f.claim(bill, request, "documents-claim", "000PRIVATEUTR123")
	_, err = docs.Receipt(ctx, f.society, f.resident, claim.ID, true)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	p, err := f.s.Verify(ctx, f.society, f.owner, claim.ID, "documents-verify", credit(claim.Reference, "old@bank", 1))
	if err != nil {
		t.Fatal(err)
	}
	f.settings("documents-disabled", "new@bank", false)
	f.exec(`UPDATE societies SET name='CHANGED LIVE NAME',address_line1='Changed Road' WHERE id=$1`, f.society)
	f.exec(`UPDATE flats SET flat_number='CHANGED-'||id WHERE society_id=$1`, f.society)
	settings := models.DefaultMaintenanceSettings()
	settings.Enabled = true
	settings.FixedPaise = 999999
	if _, err = f.billing.SaveSettings(ctx, f.society, f.owner, settings); err != nil {
		t.Fatal(err)
	}
	paidInvoice, err := docs.Invoice(ctx, f.society, f.resident, bill, true)
	if err != nil || !bytes.Equal(inv.Bytes, paidInvoice.Bytes) {
		t.Fatal("invoice changed after edits/payment", err)
	}
	payer, err := docs.Receipt(ctx, f.society, f.resident, p.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	saveIntegrationPDF(t, "integration-receipt-payer.pdf", payer)
	other, err := docs.Receipt(ctx, f.society, f.second, p.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	saveIntegrationPDF(t, "integration-receipt-other.pdf", other)
	admin, err := docs.Receipt(ctx, f.society, f.owner, p.ID, false)
	if err != nil {
		t.Fatal(err)
	}
	saveIntegrationPDF(t, "integration-receipt-admin.pdf", admin)
	if bytes.Equal(payer.Bytes, other.Bytes) {
		t.Fatal("receipt privacy filter not applied")
	}
	_, err = docs.Receipt(ctx, f.otherSociety, f.outsider, p.ID, false)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
	_, err = docs.Receipt(ctx, f.society, f.staff, p.ID, false)
	requirePaymentError(t, err, "PAYMENT_FORBIDDEN")
	if _, err = f.s.Reverse(ctx, f.society, f.owner, p.ID, "documents-reverse", models.UPIReason{Reason: "ADMIN PRIVATE REVERSAL NOTE"}); err != nil {
		t.Fatal(err)
	}
	reversed, err := docs.Receipt(ctx, f.society, f.second, p.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	saveIntegrationPDF(t, "integration-receipt-reversed.pdf", reversed)
	reversedInvoice, err := docs.Invoice(ctx, f.society, f.resident, bill, true)
	if err != nil || !bytes.Equal(inv.Bytes, reversedInvoice.Bytes) {
		t.Fatal("invoice changed after reversal", err)
	}
	f.balance(bill, 0)
	if n := f.id(`SELECT count(*) FROM maintenance_payment_ledger WHERE bill_id=$1`, bill); n != 2 {
		t.Fatal("downloads modified ledger")
	}
	// New inserts capture current details automatically, despite original INSERT
	// queries omitting issuer_snapshot. Restore unique flat labels for rendering.
	f.billing.now = func() time.Time { return time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC) }
	if _, err = f.billing.Generate(ctx, f.society, f.owner, "2026-11"); err != nil {
		t.Fatal(err)
	}
	newBill := f.id(`SELECT id FROM maintenance_bills WHERE society_id=$1 AND billing_month='2026-11-01' ORDER BY id LIMIT 1`, f.society)
	future, err := f.billing.Get(ctx, models.MaintenanceBillFilter{SocietyID: f.society, UserID: f.owner, ID: newBill})
	if err != nil || future.Issuer.CaptureSource != "issuance" || future.Issuer.Name != "CHANGED LIVE NAME" {
		t.Fatal("new issuer capture", err)
	}
	for _, query := range []string{`UPDATE maintenance_bills SET issuer_snapshot='{}' WHERE id=$1`, `UPDATE maintenance_bills SET total_paise=1 WHERE id=$1`, `DELETE FROM maintenance_bills WHERE id=$1`} {
		if _, err = f.pool.Exec(ctx, query, bill); err == nil {
			t.Fatal("bill immutability disabled")
		}
	}
	f.exec(`UPDATE flat_residents SET status='moved_out',moved_out_at=now() WHERE user_id=$1`, f.resident)
	_, err = docs.Invoice(ctx, f.society, f.resident, bill, true)
	requirePaymentError(t, err, "MAINTENANCE_BILL_NOT_FOUND")
	_, err = docs.Receipt(ctx, f.society, f.resident, p.ID, true)
	requirePaymentError(t, err, "PAYMENT_NOT_FOUND")
}
