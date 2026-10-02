package maintenancesvc

import (
	"go-server/internal/models"
	"net/url"
	"testing"
	"time"
)

func TestUPIURIExactAmountAndEncoding(t *testing.T) {
	raw := UPIURI(models.UPIConfig{UPIID: "society@bank", PayeeName: "A & B + Society"}, 123456, "reference-001", "Bill & October")
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if u.Scheme != "upi" || u.Host != "pay" || q.Get("pa") != "society@bank" || q.Get("pn") != "A & B + Society" || q.Get("am") != "1234.56" || q.Get("cu") != "INR" || q.Get("tr") != "reference-001" || q.Get("tn") != "Bill & October" {
		t.Fatal(raw)
	}
	if got := UPIURI(models.UPIConfig{}, 1, "", " "); !containsAmount(got, "0.01") {
		t.Fatal(got)
	}
}
func containsAmount(uri, amount string) bool {
	u, _ := url.Parse(uri)
	return u.Query().Get("am") == amount
}
func TestUPIReferenceNormalization(t *testing.T) {
	got, err := NormalizeUPIReference("  000abc/123  ")
	if err != nil || got != "000ABC/123" {
		t.Fatalf("%q %v", got, err)
	}
	for _, s := range []string{"", "123", "abc xyz", "<script>", "123456\n000"} {
		if _, err := NormalizeUPIReference(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
}
func TestUPICreditDateUsesBillTimezone(t *testing.T) {
	s := &PaymentService{now: func() time.Time { return time.Date(2026, 9, 30, 20, 0, 0, 0, time.UTC) }}
	if _, err := s.paymentDate("2026-10-01", "Asia/Kolkata"); err != nil {
		t.Fatal(err)
	}
	for _, date := range []string{"2026-10-02", "2026-02-30", "2026-1-01"} {
		if _, err := s.paymentDate(date, "Asia/Kolkata"); err == nil {
			t.Fatalf("accepted %q", date)
		}
	}
}
func TestUPIConfigurationValidation(t *testing.T) {
	for _, v := range []string{"merchant@bank", "society.123@okaxis"} {
		if !upiPattern.MatchString(v) {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"", "merchant", "@bank", "x@y", "merchant@bank?am=1", "merchant @bank"} {
		if upiPattern.MatchString(v) {
			t.Fatal(v)
		}
	}
}
