package maintenancesvc

import (
	"context"
	"encoding/json"
	"go-server/internal/models"
	"testing"
)

func TestBillingRunStatus(t *testing.T) {
	store := &billingStore{admin: true}
	service := New(store, nil, nil, nil, nil)
	status, err := service.BillingRun(context.Background(), 1, 2, "2026-10")
	if err != nil || status.Status != "not_generated" {
		t.Fatalf("%+v %v", status, err)
	}
	store.run = &models.MaintenanceRunResult{RunID: 7, Existing: 0}
	status, err = service.BillingRun(context.Background(), 1, 2, "2026-10")
	if err != nil || status.Status != "completed" || status.RunID != 7 || status.BillCount != 0 {
		t.Fatalf("empty run: %+v %v", status, err)
	}
	if _, err = service.BillingRun(context.Background(), 1, 2, "2026-13"); err == nil {
		t.Fatal("invalid month accepted")
	}
	store.admin = false
	if _, err = service.BillingRun(context.Background(), 1, 2, "2026-10"); err == nil {
		t.Fatal("non-admin accepted")
	}
}

func TestEvidenceAbsentPreservesLegacyDigest(t *testing.T) {
	// This exact JSON is the pre-extension request shape and field order.
	legacy := `{"bank_credit_confirmed":true,"amount_paise":100,"credit_date":"2026-09-28","reference":"000123","settings_version":1,"upi_id":"old@bank","reversal_history_acknowledged":false}`
	request := models.UPIVerifyCredit{BankCreditConfirmed: true, AmountPaise: 100, CreditDate: "2026-09-28", Reference: "000123", SettingsVersion: 1, UPIID: "old@bank"}
	encoded, err := json.Marshal(request)
	if err != nil || string(encoded) != legacy {
		t.Fatalf("legacy request changed: %s %v", encoded, err)
	}
	before, _ := digest(request)
	request.EvidenceReference = "Statement line 12"
	after, _ := digest(request)
	if before == after {
		t.Fatal("evidence missing from idempotency digest")
	}
}

func TestAdminPaymentListFilters(t *testing.T) {
	for _, filter := range []models.UPIListFilter{{BillingMonth: "2026-13"}, {BillID: -1}, {FlatID: -1}, {Resident: true, BillingMonth: "2026-13"}} {
		if _, err := paymentListMonth(filter); err == nil {
			t.Fatalf("accepted invalid filter %+v", filter)
		}
	}
	date, err := paymentListMonth(models.UPIListFilter{BillID: 1, FlatID: 2, BillingMonth: "2026-10"})
	if err != nil || date.IsZero() || date.Format("2006-01") != "2026-10" {
		t.Fatalf("%+v %v", date, err)
	}
	date, err = paymentListMonth(models.UPIListFilter{Resident: true, BillID: 1, FlatID: 2, BillingMonth: "2026-10"})
	if err != nil || date.Format("2006-01") != "2026-10" {
		t.Fatalf("resident cannot narrow own records: %v", err)
	}
}
