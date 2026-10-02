//go:build integration

package maintenancesvc

import (
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/pkg/database"
	"testing"
)

func TestResidentMaintenancePagesSummaryAndOwnershipIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	filter := models.MaintenanceBillFilter{SocietyID: f.society, UserID: f.resident, Resident: true, Page: 1, Limit: 1}
	first, err := f.billing.List(f.ctx, filter)
	if err != nil || first.TotalCount == nil || *first.TotalCount != 3 || first.TotalPages != 3 || !first.HasMore || len(first.Items) != 1 {
		t.Fatalf("first page: %+v %v", first, err)
	}
	filter.Page = 2
	next, err := f.billing.List(f.ctx, filter)
	if err != nil || next.Page != 2 || len(next.Items) != 1 || next.Items[0].ID == first.Items[0].ID {
		t.Fatalf("next page: %+v %v", next, err)
	}
	filter.UserID, filter.Page = f.second, 1
	own, err := f.billing.List(f.ctx, filter)
	if err != nil || *own.TotalCount != 1 || own.Items[0].ID != f.bills[0] || own.HasMore {
		t.Fatalf("second resident leaked flat: %+v %v", own, err)
	}
	filter.UserID = f.outsider
	outside, err := f.billing.List(f.ctx, filter)
	if err != nil || *outside.TotalCount != 0 || outside.TotalPages != 0 || len(outside.Items) != 0 {
		t.Fatalf("outside scope leaked: %+v %v", outside, err)
	}

	settings := f.settings("resident-summary-settings", "society@bank", true)
	claim := f.claim(f.bills[0], f.request(f.bills[0]), "resident-summary-claim", "SUMMARY12345")
	bill, err := f.billing.Get(f.ctx, models.MaintenanceBillFilter{SocietyID: f.society, UserID: f.resident, Resident: true, ID: f.bills[0]})
	if err != nil || bill.DisplayStatus != "pending_review" {
		t.Fatalf("pending bill: %+v %v", bill, err)
	}
	claims, err := f.s.Claims(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.resident, Resident: true, BillID: bill.ID, FlatID: bill.FlatID, BillingMonth: "2026-10"})
	if err != nil || len(claims.Items) != 1 || claims.Items[0].ID != claim.ID {
		t.Fatalf("own filtered claim: %+v %v", claims, err)
	}
	hidden, err := f.s.Claims(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.second, Resident: true, BillID: bill.ID, FlatID: bill.FlatID, BillingMonth: "2026-10"})
	if err != nil || len(hidden.Items) != 0 {
		t.Fatalf("other submitter claim disclosed: %+v %v", hidden, err)
	}
	filter.UserID, filter.FlatID, filter.DisplayStatus = f.resident, bill.FlatID, "pending_review"
	pending, err := f.billing.List(f.ctx, filter)
	if err != nil || *pending.TotalCount != 1 || pending.Items[0].DisplayStatus != "pending_review" {
		t.Fatalf("filtered pending: %+v %v", pending, err)
	}
	payment, err := f.s.Verify(f.ctx, f.society, f.owner, claim.ID, "resident-summary-verify", credit(claim.Reference, "society@bank", settings.Version))
	if err != nil {
		t.Fatal(err)
	}
	d := &database.Database{Pool: f.pool}
	repo := repository.NewMaintenanceRepository(d, repository.NewTransactionManager(d))
	summary, err := repo.Outstanding(f.ctx, f.society, f.resident, bill.FlatID, f.billing.now())
	if err != nil || summary.TotalPaidPaise != bill.TotalPaise || summary.TotalOutstandingPaise != 0 || summary.CurrentBill == nil || summary.CurrentBill.DisplayStatus != "paid" || summary.CurrentBill.PaidOn != "2026-10-14" {
		t.Fatalf("paid summary: %+v %v", summary, err)
	}
	records, err := f.s.Payments(f.ctx, models.UPIListFilter{SocietyID: f.society, UserID: f.second, Resident: true, BillID: bill.ID, FlatID: bill.FlatID, BillingMonth: "2026-10"})
	if err != nil || len(records.Items) != 1 || records.Items[0].Reference != nil {
		t.Fatalf("shared flat receipt privacy: %+v %v", records, err)
	}
	if _, err = f.s.Reverse(f.ctx, f.society, f.owner, payment.ID, "resident-summary-reverse", models.UPIReason{Reason: "Incorrect bank credit"}); err != nil {
		t.Fatal(err)
	}
	summary, err = repo.Outstanding(f.ctx, f.society, f.resident, bill.FlatID, f.billing.now())
	if err != nil || summary.TotalPaidPaise != 0 || summary.TotalOutstandingPaise != bill.TotalPaise || summary.CurrentBill.DisplayStatus == "paid" || summary.CurrentBill.PaidOn != "" {
		t.Fatalf("reversed summary: %+v %v", summary, err)
	}
	if _, err = repo.Outstanding(f.ctx, f.society, f.outsider, bill.FlatID, f.billing.now()); err == nil {
		t.Fatal("outsider read summary")
	}
}
