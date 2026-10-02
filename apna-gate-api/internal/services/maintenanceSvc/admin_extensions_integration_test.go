//go:build integration

package maintenancesvc

import (
	"go-server/internal/models"
	"strings"
	"testing"
)

func TestMaintenanceAdminEvidenceAndReporting(t *testing.T) {
	f := newPaymentFixture(t)
	f.settings("setup", "old@bank", true)
	bill := f.bills[0]
	claim := f.claim(bill, f.request(bill), "claim", "000ADMIN123")
	flat := f.id(`SELECT flat_id FROM maintenance_bills WHERE id=$1`, bill)
	req := credit(claim.Reference, "old@bank", 1)
	req.EvidenceReference = "  Bank statement line 42  "
	payment, err := f.s.Verify(f.ctx, f.society, f.owner, claim.ID, "verify-evidence", req)
	if err != nil || payment.EvidenceReference != "Bank statement line 42" {
		t.Fatalf("%+v %v", payment, err)
	}
	retry, err := f.s.Verify(f.ctx, f.society, f.owner, claim.ID, "verify-evidence", req)
	if err != nil || retry.ID != payment.ID {
		t.Fatalf("retry %+v %v", retry, err)
	}
	resident, err := f.s.Payment(f.ctx, f.society, f.resident, payment.ID, true)
	if err != nil || resident.EvidenceReference != "" {
		t.Fatalf("resident evidence leak: %+v %v", resident, err)
	}
	if _, err = f.pool.Exec(f.ctx, `UPDATE maintenance_payments SET evidence_reference='changed' WHERE id=$1`, payment.ID); err == nil {
		t.Fatal("mutable evidence")
	}
	audit, err := f.s.Audit(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.owner}, bill)
	found := false
	for _, event := range audit.Items {
		if event.Details["evidence_reference"] == "Bank statement line 42" {
			found = true
		}
	}
	if err != nil || !found {
		t.Fatalf("evidence not audited: %+v %v", audit, err)
	}
	summary, err := f.s.Summary(f.ctx, f.society, f.owner, "2026-10", flat)
	if err != nil || summary.BilledCount != 1 || summary.CollectedPaise != 123456 || summary.OutstandingPaise != 0 {
		t.Fatalf("summary %+v %v", summary, err)
	}
	payments, err := f.s.Payments(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.owner, BillID: bill, FlatID: flat, BillingMonth: "2026-10"})
	if err != nil || len(payments.Items) != 1 || payments.Items[0].Flat == nil || payments.Items[0].Flat.ID != flat {
		t.Fatalf("filtered payments %+v %v", payments, err)
	}
	claims, err := f.s.Claims(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.owner, BillID: bill, FlatID: flat, BillingMonth: "2026-10"})
	if err != nil || len(claims.Items) != 1 || claims.Items[0].Flat == nil || claims.Items[0].BillNumber == "" {
		t.Fatalf("filtered claims %+v %v", claims, err)
	}
	empty, err := f.s.Payments(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.owner, BillingMonth: "2026-09"})
	if err != nil || len(empty.Items) != 0 {
		t.Fatalf("month filter %+v %v", empty, err)
	}
	_, err = f.s.Summary(f.ctx, f.society, f.staff, "2026-10", flat)
	requirePaymentError(t, err, "PAYMENT_FORBIDDEN")
	other, err := f.s.Summary(f.ctx, f.otherSociety, f.outsider, "2026-10", flat)
	if err != nil || other.BilledCount != 0 {
		t.Fatalf("cross tenant %+v %v", other, err)
	}
	if _, err = f.s.Reverse(f.ctx, f.society, f.owner, payment.ID, "reverse", models.UPIReason{Reason: "Incorrect match"}); err != nil {
		t.Fatal(err)
	}
	summary, err = f.s.Summary(f.ctx, f.society, f.owner, "2026-10", flat)
	if err != nil || summary.CollectedPaise != 0 || summary.OutstandingPaise != 123456 {
		t.Fatalf("reversal summary %+v %v", summary, err)
	}
	direct := models.UPIDirectPayment{BillID: f.bills[1], UPIVerifyCredit: credit("000DIRECT123", "old@bank", 1)}
	direct.EvidenceReference = strings.Repeat("x", 501)
	_, err = f.s.Record(f.ctx, f.society, f.owner, "long-evidence", direct)
	requirePaymentError(t, err, "MAINTENANCE_INVALID")
	direct.EvidenceReference = "Direct bank statement"
	saved, err := f.s.Record(f.ctx, f.society, f.owner, "direct-evidence", direct)
	if err != nil || saved.EvidenceReference != direct.EvidenceReference {
		t.Fatalf("direct evidence %+v %v", saved, err)
	}
}
