package models

import (
	"testing"
	"time"
)

func TestMaintenancePresentationTimezoneSettlementAndReversal(t *testing.T) {
	for _, tc := range []struct {
		at, claim      string
		paid           int64
		state, message string
	}{
		{"2026-09-30T18:29:59Z", "none", 0, "unpaid", "Due today"},
		{"2026-09-30T18:30:00Z", "none", 0, "overdue", "1 day overdue"},
		{"2026-09-30T18:30:00Z", "pending", 0, "pending_review", "1 day overdue"},
		{"2026-09-30T18:30:00Z", "rejected", 100, "paid", "Payment verified"},
		{"2026-09-30T18:30:00Z", "verified", 0, "overdue", "1 day overdue"},
	} {
		b := MaintenanceBill{DueDate: "2026-09-30", Timezone: "Asia/Kolkata", TotalPaise: 100, PaidAmountPaise: tc.paid, OutstandingAmountPaise: 100 - tc.paid, PaymentClaimStatus: tc.claim}
		now, _ := time.Parse(time.RFC3339, tc.at)
		if err := b.PopulatePresentation(now); err != nil {
			t.Fatal(err)
		}
		if b.DisplayStatus != tc.state || b.DueMessage != tc.message {
			t.Fatalf("%+v: %+v", tc, b)
		}
	}
}
