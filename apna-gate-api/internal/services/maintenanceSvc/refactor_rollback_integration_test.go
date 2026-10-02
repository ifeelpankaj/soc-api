//go:build integration

package maintenancesvc

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go-server/internal/models"
)

// Snapshot actual persisted rows, rather than relying only on mock call order.
func financialSnapshot(t *testing.T, f *paymentFixture) string {
	t.Helper()
	tables := []string{"maintenance_billing_runs", "maintenance_bills", "maintenance_bill_items", "maintenance_billing_audit", "maintenance_payment_requests", "maintenance_payment_claims", "maintenance_payments", "maintenance_payment_reference_reservations", "maintenance_payment_ledger", "maintenance_payment_audit_events", "maintenance_payment_idempotency", "notification_outbox"}
	var snapshot strings.Builder
	for _, table := range tables {
		var rows string
		if err := f.pool.QueryRow(f.ctx, "SELECT COALESCE(jsonb_agg(to_jsonb(r) ORDER BY to_jsonb(r)::text)::text,'[]') FROM "+table+" r").Scan(&rows); err != nil {
			t.Fatal(err)
		}
		snapshot.WriteString(table + rows)
	}
	return snapshot.String()
}

func TestFinancialWriteFailuresRollbackIntegration(t *testing.T) {
	f := newPaymentFixture(t)
	v := f.settings("rollback-settings", "rollback@bank", true)
	request := f.request(f.bills[0])
	claim := f.claim(f.bills[0], request, "rollback-claim", "ROLLBACK123")
	f.exec(`CREATE FUNCTION refactor_fail_write() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected refactor failure'; END $$`)
	tests := []struct {
		name, table, event string
		deferred, billing  bool
	}{
		{"release claim reference", "maintenance_payment_reference_reservations", "DELETE", false, false},
		{"review claim", "maintenance_payment_claims", "UPDATE", false, false},
		{"create payment", "maintenance_payments", "INSERT", false, false},
		{"reserve paid reference", "maintenance_payment_reference_reservations", "INSERT", false, false},
		{"collection ledger", "maintenance_payment_ledger", "INSERT", false, false},
		{"close payment request", "maintenance_payment_requests", "UPDATE", false, false},
		{"payment audit", "maintenance_payment_audit_events", "INSERT", false, false},
		{"payment outbox", "notification_outbox", "INSERT", false, false},
		{"payment idempotency", "maintenance_payment_idempotency", "INSERT", false, false},
		{"payment commit", "maintenance_payment_idempotency", "INSERT", true, false},
		{"billing run", "maintenance_billing_runs", "INSERT", false, true},
		{"bill", "maintenance_bills", "INSERT", false, true},
		{"bill item", "maintenance_bill_items", "INSERT", false, true},
		{"billing outbox", "notification_outbox", "INSERT", false, true},
		{"billing audit", "maintenance_payment_audit_events", "INSERT", false, true},
		{"billing idempotency", "maintenance_payment_idempotency", "INSERT", false, true},
		{"billing commit", "maintenance_payment_idempotency", "INSERT", true, true},
	}
	f.billing.now = func() time.Time { return time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC) }
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			before := financialSnapshot(t, f)
			kind, timing := "", ""
			if tc.deferred {
				kind, timing = "CONSTRAINT ", "DEFERRABLE INITIALLY DEFERRED"
			}
			f.exec(fmt.Sprintf("CREATE %sTRIGGER refactor_failure AFTER %s ON %s %s FOR EACH ROW EXECUTE FUNCTION refactor_fail_write()", kind, tc.event, tc.table, timing))
			var err error
			if tc.billing {
				_, err = f.billing.GenerateCommand(f.ctx, f.society, f.owner, fmt.Sprintf("billing-failure-%d", i), models.MaintenanceMonthRequest{BillingMonth: "2026-11"})
			} else {
				_, err = f.s.Verify(f.ctx, f.society, f.owner, claim.ID, fmt.Sprintf("payment-failure-%d", i), credit("ROLLBACK123", "rollback@bank", v.Version))
			}
			f.exec("DROP TRIGGER refactor_failure ON " + tc.table)
			if err == nil || !strings.Contains(err.Error(), "injected refactor failure") {
				t.Fatalf("failure point was not reached: %v", err)
			}
			if after := financialSnapshot(t, f); after != before {
				t.Fatal("failed operation retained financial or outbox writes")
			}
		})
	}
}
