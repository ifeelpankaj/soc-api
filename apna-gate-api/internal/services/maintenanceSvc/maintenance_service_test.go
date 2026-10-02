package maintenancesvc

import (
	"context"
	"errors"
	"go-server/internal/models"
	"math"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }
func TestExactCharges(t *testing.T) {
	for _, tc := range []struct {
		name, model             string
		area, rate, fixed, want int64
	}{
		{"fixed", "fixed", 0, 0, 200000, 200000},
		{"area", "per_sqft", 100000, 300, 0, 300000},
		{"fraction rounds up", "per_sqft", 101, 50, 0, 51},
		{"fraction rounds down", "per_sqft", 101, 49, 0, 49},
		{"hybrid", "hybrid", 100025, 300, 50000, 350075},
		{"type", "flat_type", 0, 0, 0, 250000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := models.DefaultMaintenanceSettings()
			s.PricingModel = tc.model
			s.FixedPaise = tc.fixed
			s.AreaRatePaise = tc.rate
			s.TypeRates = map[string]int64{"2 BHK": 250000}
			_, got, err := Charges(s, models.MaintenanceFlat{FlatType: ptr("2 BHK"), AreaHundredths: &tc.area})
			if err != nil || got != tc.want {
				t.Fatalf("amount=%d err=%v want=%d", got, err, tc.want)
			}
		})
	}
	if _, err := AreaCharge(math.MaxInt64, math.MaxInt64); err == nil {
		t.Fatal("overflow accepted")
	}
	s := models.DefaultMaintenanceSettings()
	s.PricingModel = "hybrid"
	s.FixedPaise = math.MaxInt64
	s.AreaRatePaise = 100
	if _, _, err := Charges(s, models.MaintenanceFlat{AreaHundredths: ptr(int64(100))}); err == nil {
		t.Fatal("hybrid overflow accepted")
	}
	for _, model := range []string{"per_sqft", "hybrid", "flat_type"} {
		s.PricingModel = model
		if _, _, err := Charges(s, models.MaintenanceFlat{}); err == nil {
			t.Fatalf("missing data accepted for %s", model)
		}
	}
}
func TestBillingCalendar(t *testing.T) {
	now := time.Date(2026, 9, 30, 20, 0, 0, 0, time.UTC)
	month, err := CurrentMonth("2026-10", now, "Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CurrentMonth("2026-09", now, "Asia/Kolkata"); err == nil {
		t.Fatal("previous local month accepted")
	}
	if got := DueDate(month, 20, 10).Format("2006-01-02"); got != "2026-11-10" {
		t.Fatal(got)
	}
	if got := DueDate(time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), 28, 1).Format("2006-01-02"); got != "2027-01-01" {
		t.Fatal(got)
	}
	if got := DueDate(month, 10, 10).Format("2006-01-02"); got != "2026-10-10" {
		t.Fatal(got)
	}
}

type billingStore struct {
	Store
	settings     models.MaintenanceSettings
	flats        []models.MaintenanceFlat
	run          *models.MaintenanceRunResult
	issued       int
	admin        bool
	seenStatuses []string
	issuedTerms  models.MaintenanceSettings
	issuedBills  []models.MaintenanceBill
}

func (r *billingStore) Locked(ctx context.Context, _ int64, fn func(context.Context) error) error {
	return fn(ctx)
}
func (r *billingStore) Admin(context.Context, int64, int64) (bool, error) { return r.admin, nil }
func (r *billingStore) Settings(context.Context, int64) (models.MaintenanceSettings, error) {
	return r.settings, nil
}
func (r *billingStore) Flats(_ context.Context, _ int64, statuses []string) ([]models.MaintenanceFlat, error) {
	r.seenStatuses = statuses
	return r.flats, nil
}
func (r *billingStore) Run(context.Context, int64, string) (*models.MaintenanceRunResult, error) {
	return r.run, nil
}
func (r *billingStore) Issue(_ context.Context, _, _ int64, _ string, terms models.MaintenanceSettings, b []models.MaintenanceBill) (models.MaintenanceRunResult, error) {
	r.issued++
	r.issuedTerms = terms
	r.issuedBills = append([]models.MaintenanceBill(nil), b...)
	r.run = &models.MaintenanceRunResult{RunID: 1, Existing: len(b)}
	return models.MaintenanceRunResult{RunID: 1, Created: len(b)}, nil
}

func (r *billingStore) RunTerms(context.Context, int64, string) (models.MaintenanceSettings, []int64, error) {
	ids := make([]int64, 0, len(r.issuedBills))
	for _, bill := range r.issuedBills {
		ids = append(ids, bill.FlatID)
	}
	return r.issuedTerms, ids, nil
}
func (r *billingStore) IssueMissing(_ context.Context, _, _ int64, previous models.MaintenanceRunResult, bills []models.MaintenanceBill) (models.MaintenanceRunResult, error) {
	r.issuedBills = append(r.issuedBills, bills...)
	r.run = &models.MaintenanceRunResult{RunID: previous.RunID, Existing: previous.Existing + len(bills)}
	return models.MaintenanceRunResult{RunID: previous.RunID, Existing: previous.Existing, Created: len(bills)}, nil
}

type allowOperational struct{}

func (allowOperational) EnsureSocietyOperational(context.Context, int64) error { return nil }
func billingFixture() (*Service, *billingStore) {
	v := models.DefaultMaintenanceSettings()
	v.Enabled = true
	v.FixedPaise = 200000
	r := &billingStore{settings: v, admin: true, flats: []models.MaintenanceFlat{{ID: 1, FlatNumber: "101"}, {ID: 2, FlatNumber: "102"}}}
	s := New(r, allowOperational{}, []string{"occupied"}, nil, nil)
	s.now = func() time.Time { return time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC) }
	return s, r
}
func TestGenerateMissingDataAndRetry(t *testing.T) {
	s, r := billingFixture()
	r.settings.PricingModel = "per_sqft"
	r.settings.AreaRatePaise = 300
	_, err := s.Generate(context.Background(), 1, 1, "2026-10")
	var issues *ValidationIssues
	if !errors.As(err, &issues) || len(issues.Issues) != 2 || r.issued != 0 {
		t.Fatalf("invalid batch issued or issues lost: %v", err)
	}
	for i := range r.flats {
		r.flats[i].AreaHundredths = ptr(int64(100000))
	}
	result, err := s.Generate(context.Background(), 1, 1, "2026-10")
	if err != nil || result.Created != 2 {
		t.Fatalf("%+v %v", result, err)
	}
	r.settings.AreaRatePaise = 900
	r.flats = append(r.flats, models.MaintenanceFlat{ID: 3})
	result, err = s.Generate(context.Background(), 1, 1, "2026-10")
	if !errors.As(err, &issues) || len(issues.Issues) != 1 || issues.Issues[0].FlatID != 3 || len(r.issuedBills) != 2 {
		t.Fatalf("missing flat validation: %+v %v", result, err)
	}
	r.flats[0].AreaHundredths = nil // Existing flats must not be revalidated.
	r.flats[2].AreaHundredths = ptr(int64(200000))
	r.settings.DueDay = 28
	s.statuses = []string{"vacant"} // Keep the original eligibility policy.
	p, err := s.Preview(context.Background(), 1, 1, "2026-10")
	if err != nil || len(p.Bills) != 1 || p.Bills[0].TotalPaise != 600000 || p.Bills[0].DueDate != "2026-10-10" {
		t.Fatalf("preview %+v %v", p, err)
	}
	result, err = s.Generate(context.Background(), 1, 1, "2026-10")
	if err != nil || result.Existing != 2 || result.Created != 1 || result.RunID != 1 {
		t.Fatalf("reconciliation %+v %v", result, err)
	}
	result, err = s.Generate(context.Background(), 1, 1, "2026-10")
	if err != nil || result.Existing != 3 || result.Created != 0 || r.issued != 1 {
		t.Fatalf("repeat %+v %v", result, err)
	}
	if len(r.seenStatuses) != 1 || r.seenStatuses[0] != "occupied" {
		t.Fatal("eligibility policy not passed to query")
	}
}

func TestMissingBillsEmptyRunAndHistoricalRestriction(t *testing.T) {
	s, r := billingFixture()
	r.flats = nil
	if _, err := s.Generate(context.Background(), 1, 1, "2026-10"); err != nil {
		t.Fatal(err)
	}
	r.flats = []models.MaintenanceFlat{{ID: 4}}
	scheduled, err := s.generate(context.Background(), 1, 0, "", true)
	if err != nil || scheduled.Created != 0 || len(r.issuedBills) != 0 {
		t.Fatal("scheduled run reconciled issued month", scheduled, err)
	}
	result, err := s.Generate(context.Background(), 1, 1, "2026-10")
	if err != nil || result.Created != 1 || result.Existing != 0 {
		t.Fatal(result, err)
	}
	s.now = func() time.Time { return time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC) }
	if _, err = s.Generate(context.Background(), 1, 1, "2026-10"); err == nil {
		t.Fatal("past issued month reopened")
	}
}
func TestGenerationAuthorizationDisableAndScheduling(t *testing.T) {
	s, r := billingFixture()
	r.admin = false
	if _, err := s.Generate(context.Background(), 1, 2, "2026-10"); err == nil {
		t.Fatal("non-admin accepted")
	}
	r.admin = true
	r.settings.Enabled = false
	if _, err := s.Generate(context.Background(), 1, 1, "2026-10"); err == nil {
		t.Fatal("disabled generation accepted")
	}
	r.settings.Enabled = true
	r.settings.BillingDay = 20
	if _, err := s.generate(context.Background(), 1, 0, "", true); err != nil || r.issued != 0 {
		t.Fatal("ran before billing day")
	}
	r.settings.BillingDay = 1
	result, err := s.generate(context.Background(), 1, 0, "", true)
	if err != nil || result.Created != 2 {
		t.Fatalf("missed-day recovery: %+v %v", result, err)
	}
}
