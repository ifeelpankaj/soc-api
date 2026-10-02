package documentpdf

import (
	"bytes"
	"context"
	"errors"
	"go-server/internal/models"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func exampleBill() models.MaintenanceBill {
	return models.MaintenanceBill{ID: 123, SocietyID: 1, BillNumber: "AG-ABC001-20260901053000-A-101", BillingMonth: "2026-09", DueDate: "2026-09-30", Timezone: "Asia/Kolkata", Currency: "INR", CreatedAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), TotalPaise: 223456,
		Issuer: models.MaintenanceIssuer{Name: "ABC Housing Society", SocietyCode: "ABC001", AddressLine1: "12 Garden Road", AddressLine2: "Near Central Park", City: "Pune", State: "Maharashtra", Pincode: "411001", Country: "India", Email: "office@example.com", PhoneNumber: "+91 90000 00000", CapturedAt: time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC), CaptureSource: "rollout"},
		Flat:   models.MaintenanceFlat{FlatNumber: "A-101"}, Items: []models.MaintenanceItem{{Description: "Monthly maintenance", AmountPaise: 200000}, {Description: "Shared services", AmountPaise: 23456}},
		BilledParty: map[string]any{"name": "SECRET FORMER RESIDENT"}}
}
func examplePayment(b models.MaintenanceBill) models.UPIPayment {
	ref := "001234567890"
	user := int64(42)
	return models.UPIPayment{ID: 88, SocietyID: b.SocietyID, BillID: b.ID, BillNumber: b.BillNumber, AmountPaise: b.TotalPaise, Currency: "INR", CreditDate: "2026-09-10", ReceiptNumber: "MPR-11111111-2222-3333-4444-555555555555", SettingsVersion: 1, Destination: models.UPIConfig{UPIID: "society@bank", PayeeName: "ABC Housing Society"}, Status: "verified", VerifiedAt: time.Date(2026, 9, 11, 10, 30, 0, 0, time.UTC), Reference: &ref, PayerID: &user, VerifiedBy: &user, VerifiedByName: "Society Admin"}
}
func TestMoneyExact(t *testing.T) {
	for _, tc := range []struct {
		n    int64
		want string
	}{{1, "INR 0.01"}, {200000, "INR 2,000.00"}, {math.MaxInt64, "INR 92,233,720,368,547,758.07"}, {math.MinInt64, "INR -92,233,720,368,547,758.08"}} {
		if got := Money(tc.n); got != tc.want {
			t.Fatalf("%d: %s", tc.n, got)
		}
	}
}
func TestInvoiceIgnoresFinancialStateAndPersonalData(t *testing.T) {
	r := New()
	b := exampleBill()
	one, err := r.Invoice(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	b.PaidAmountPaise = b.TotalPaise
	b.OutstandingAmountPaise = 0
	b.Status = "paid"
	b.PaymentClaimStatus = "verified"
	b.BilledParty = map[string]any{"name": "DIFFERENT PERSON"}
	two, err := r.Invoice(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(one, two) {
		t.Fatal("invoice changed due to financial/personal fields")
	}
	if !bytes.HasPrefix(one, []byte("%PDF-")) {
		t.Fatal("invalid PDF")
	}
}
func TestUnsupportedTextAndCancellation(t *testing.T) {
	r := New()
	b := exampleBill()
	b.Issuer.Name = "सोसायटी"
	_, err := r.Invoice(context.Background(), b)
	var app *models.AppError
	if !errors.As(err, &app) || app.Code != "DOCUMENT_UNSUPPORTED_TEXT" {
		t.Fatalf("%v", err)
	}
	b = exampleBill()
	b.Issuer.Name = "Société Résidence"
	if _, err = r.Invoice(context.Background(), b); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = r.Invoice(ctx, b)
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
func TestReceiptRequiresPaymentAndMatchingBill(t *testing.T) {
	b := exampleBill()
	p := examplePayment(b)
	p.Status = "pending"
	if _, err := New().Receipt(context.Background(), b, p); err == nil {
		t.Fatal("pending receipt rendered")
	}
	p.Status = "verified"
	p.BillID++
	if _, err := New().Receipt(context.Background(), b, p); err == nil {
		t.Fatal("mismatched bill rendered")
	}
}

// Set MAINTENANCE_PDF_QA_DIR to an ignored temporary directory to export visual
// fixtures. Normal unit tests generate PDFs only in memory.
func TestDocumentVisualFixtures(t *testing.T) {
	b := exampleBill()
	b.Outstanding = &models.MaintenanceOutstanding{CalculatedAt: time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC), CurrentMonthPaise: b.TotalPaise, PreviousOutstandingPaise: 70000, TotalOutstandingPaise: b.TotalPaise + 70000}
	p := examplePayment(b)
	ctx := context.Background()
	r := New()
	save := func(name string, data []byte, err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		if len(data) < 1000 {
			t.Fatal("empty PDF")
		}
		if dir := os.Getenv("MAINTENANCE_PDF_QA_DIR"); dir != "" {
			// #nosec G703 -- QA output directory is explicitly supplied by the test runner.
			if err = os.MkdirAll(dir, 0700); err != nil {
				t.Fatal(err)
			}
			// #nosec G703 -- QA output directory is explicitly supplied by the test runner; names below are fixed literals.
			if err = os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	data, err := r.Invoice(ctx, b)
	save("invoice.pdf", data, err)
	data, err = r.Receipt(ctx, b, p)
	save("receipt-admin.pdf", data, err)
	p.Reference = nil
	p.PayerID = nil
	p.VerifiedBy = nil
	data, err = r.Receipt(ctx, b, p)
	save("receipt-resident.pdf", data, err)
	p.Status = "reversed"
	when := p.VerifiedAt.Add(24 * time.Hour)
	p.ReversedAt = &when
	data, err = r.Receipt(ctx, b, p)
	save("receipt-reversed.pdf", data, err)
	b.Issuer.AddressLine1 = strings.Repeat("Long address segment & building name, ", 12)
	b.Items = nil
	b.TotalPaise = 0
	for i := 0; i < 90; i++ {
		b.Items = append(b.Items, models.MaintenanceItem{Description: "Shared services: " + strings.Repeat("description with long wrapped text ", 6), AmountPaise: 12345})
		b.TotalPaise += 12345
	}
	data, err = r.Invoice(ctx, b)
	save("invoice-long.pdf", data, err)
	b = exampleBill()
	b.Items = []models.MaintenanceItem{{Description: "Maximum supported exact amount", AmountPaise: math.MaxInt64}}
	b.TotalPaise = math.MaxInt64
	data, err = r.Invoice(ctx, b)
	save("invoice-large-amount.pdf", data, err)
}
