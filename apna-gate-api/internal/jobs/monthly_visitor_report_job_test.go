package jobs

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-server/internal/models"
)

type reportStoreFake struct {
	societies     []models.VisitorReportSociety
	earliest      map[int64]time.Time
	deliveries    []models.VisitorReportDelivery
	rows          map[int64][]models.MonthlyVisitorReportRow
	recipients    map[int64][]string
	ensured       map[string]bool
	completed     map[string]bool
	queries       []int64
	allowResend   bool
	resendMonth   time.Time
	claimAttempts map[string]int32
}

func (f *reportStoreFake) WithAdvisoryLock(ctx context.Context, _ string, fn func(context.Context) error) (bool, error) {
	return true, fn(ctx)
}
func (f *reportStoreFake) ListActiveReportSocieties(context.Context) ([]models.VisitorReportSociety, error) {
	return f.societies, nil
}
func (f *reportStoreFake) GetEarliestVisitorEntryCreatedAt(_ context.Context, societyID int64) (*time.Time, error) {
	value, ok := f.earliest[societyID]
	if !ok {
		return nil, nil
	}
	return &value, nil
}
func (f *reportStoreFake) EnsureReportDelivery(_ context.Context, societyID int64, month time.Time) error {
	f.ensured[deliveryKey(societyID, month)] = true
	return nil
}
func (f *reportStoreFake) ListClaimableReportDeliveries(_ context.Context, _ int32, allowResend bool, resendMonth time.Time) ([]models.VisitorReportDelivery, error) {
	f.allowResend, f.resendMonth = allowResend, resendMonth
	return f.deliveries, nil
}
func (f *reportStoreFake) ClaimReportDelivery(_ context.Context, societyID int64, month, _ time.Time, _ bool) (*models.VisitorReportDelivery, error) {
	if f.claimAttempts == nil {
		f.claimAttempts = map[string]int32{}
	}
	key := deliveryKey(societyID, month)
	f.claimAttempts[key]++
	return &models.VisitorReportDelivery{SocietyID: societyID, ReportMonth: month, AttemptCount: f.claimAttempts[key]}, nil
}
func (f *reportStoreFake) CompleteReportDelivery(_ context.Context, societyID int64, month time.Time, _ []string, _ string, _ bool) error {
	f.completed[deliveryKey(societyID, month)] = true
	return nil
}
func (f *reportStoreFake) FailReportDelivery(context.Context, int64, time.Time, []string, error, bool) error {
	return nil
}
func (f *reportStoreFake) ListMonthlyVisitorReportRecipients(_ context.Context, societyID int64) ([]string, error) {
	return f.recipients[societyID], nil
}
func (f *reportStoreFake) ListMonthlyVisitorReportRows(_ context.Context, societyID int64, _, _ time.Time) ([]models.MonthlyVisitorReportRow, error) {
	f.queries = append(f.queries, societyID)
	return f.rows[societyID], nil
}

type sentReport struct {
	recipients     []string
	rowNames       []string
	idempotencyKey string
}

type reportEmailFake struct {
	sent map[int64]sentReport
}

func (f *reportEmailFake) SendMonthlyVisitorReport(_ context.Context, recipients []string, _ string, _ time.Time, filename string, content []byte, idempotencyKey string) (string, error) {
	var societyID int64
	if _, err := fmt.Sscanf(filename, "visitor-report-%d-", &societyID); err != nil {
		return "", err
	}
	records, err := csv.NewReader(strings.NewReader(string(content[3:]))).ReadAll()
	if err != nil {
		return "", err
	}
	names := make([]string, 0, len(records)-1)
	for _, record := range records[1:] {
		names = append(names, record[1])
	}
	f.sent[societyID] = sentReport{recipients: append([]string(nil), recipients...), rowNames: names, idempotencyKey: idempotencyKey}
	return fmt.Sprintf("provider-%d", societyID), nil
}

func TestMonthlyVisitorReportDevelopmentResendUsesLatestMonthAndUniqueAttemptKey(t *testing.T) {
	month := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	store := &reportStoreFake{
		deliveries: []models.VisitorReportDelivery{{SocietyID: 1, SocietyName: "A", ReportMonth: month}},
		rows:       map[int64][]models.MonthlyVisitorReportRow{1: reportRows("A", 1)},
		recipients: map[int64][]string{1: {"a@example.com"}}, ensured: map[string]bool{}, completed: map[string]bool{},
	}
	email := &reportEmailFake{sent: map[int64]sentReport{}}
	job := NewMonthlyVisitorReportJob(store, email, MonthlyVisitorReportJobConfig{Location: time.UTC, AllowResend: true})
	job.now = func() time.Time { return time.Date(2026, time.September, 2, 0, 0, 0, 0, time.UTC) }
	if _, err := job.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	firstKey := email.sent[1].idempotencyKey
	if _, err := job.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	secondKey := email.sent[1].idempotencyKey
	if !store.allowResend || !store.resendMonth.Equal(month) {
		t.Fatalf("resend scope allow=%v month=%v", store.allowResend, store.resendMonth)
	}
	if firstKey != "monthly-visitor-report/1/2026-08/attempt-1" || secondKey != "monthly-visitor-report/1/2026-08/attempt-2" {
		t.Fatalf("idempotency keys = %q, %q", firstKey, secondKey)
	}
}

func TestMonthlyVisitorReportJobKeepsSocietiesIsolated(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	month := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	store := &reportStoreFake{
		societies: []models.VisitorReportSociety{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		earliest:  map[int64]time.Time{1: month, 2: month},
		deliveries: []models.VisitorReportDelivery{
			{SocietyID: 1, SocietyName: "A", ReportMonth: month},
			{SocietyID: 2, SocietyName: "B", ReportMonth: month},
		},
		rows: map[int64][]models.MonthlyVisitorReportRow{
			1: reportRows("A", 5),
			2: reportRows("B", 7),
		},
		recipients: map[int64][]string{1: {"a-owner@example.com"}, 2: {"b-owner@example.com"}},
		ensured:    map[string]bool{}, completed: map[string]bool{},
	}
	email := &reportEmailFake{sent: map[int64]sentReport{}}
	job := NewMonthlyVisitorReportJob(store, email, MonthlyVisitorReportJobConfig{Location: loc, BatchSize: 50})
	job.now = func() time.Time { return time.Date(2026, time.September, 2, 0, 0, 0, 0, loc) }
	if _, err := job.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(email.sent[1].rowNames) != 5 || len(email.sent[2].rowNames) != 7 {
		t.Fatalf("row counts = A:%d B:%d, want A:5 B:7", len(email.sent[1].rowNames), len(email.sent[2].rowNames))
	}
	for _, name := range email.sent[1].rowNames {
		if !strings.HasPrefix(name, "A-") {
			t.Fatalf("society A report contains foreign row %q", name)
		}
	}
	for _, name := range email.sent[2].rowNames {
		if !strings.HasPrefix(name, "B-") {
			t.Fatalf("society B report contains foreign row %q", name)
		}
	}
	if got := email.sent[1].recipients; len(got) != 1 || got[0] != "a-owner@example.com" {
		t.Fatalf("society A recipients = %#v", got)
	}
	if got := email.sent[2].recipients; len(got) != 1 || got[0] != "b-owner@example.com" {
		t.Fatalf("society B recipients = %#v", got)
	}
	if !store.completed[deliveryKey(1, month)] || !store.completed[deliveryKey(2, month)] {
		t.Fatal("delivery completion was not persisted for both societies")
	}
}

func TestMonthlyVisitorReportFallsBackToSocietyEmail(t *testing.T) {
	month := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	fallback := " Society@Example.COM "
	store := &reportStoreFake{
		rows: map[int64][]models.MonthlyVisitorReportRow{}, recipients: map[int64][]string{},
		completed: map[string]bool{}, ensured: map[string]bool{},
	}
	email := &reportEmailFake{sent: map[int64]sentReport{}}
	job := NewMonthlyVisitorReportJob(store, email, MonthlyVisitorReportJobConfig{Location: time.UTC})
	if err := job.deliver(context.Background(), models.VisitorReportDelivery{SocietyID: 3, SocietyName: "C", SocietyEmail: &fallback, ReportMonth: month}); err != nil {
		t.Fatal(err)
	}
	got := email.sent[3].recipients
	if len(got) != 1 || got[0] != "society@example.com" {
		t.Fatalf("fallback recipients = %#v", got)
	}
}

func TestMonthlyVisitorReportCatchUpUsesPerSocietyEarliestMonth(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	earliest := time.Date(2026, time.June, 15, 20, 0, 0, 0, time.UTC)
	store := &reportStoreFake{
		societies: []models.VisitorReportSociety{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		earliest:  map[int64]time.Time{1: earliest},
		rows:      map[int64][]models.MonthlyVisitorReportRow{}, recipients: map[int64][]string{},
		ensured: map[string]bool{}, completed: map[string]bool{},
	}
	job := NewMonthlyVisitorReportJob(store, &reportEmailFake{sent: map[int64]sentReport{}}, MonthlyVisitorReportJobConfig{Location: loc})
	job.now = func() time.Time { return time.Date(2026, time.September, 2, 0, 0, 0, 0, loc) }
	if _, err := job.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}

	for _, month := range []time.Month{time.June, time.July, time.August} {
		key := deliveryKey(1, time.Date(2026, month, 1, 0, 0, 0, 0, time.UTC))
		if !store.ensured[key] {
			t.Fatalf("society A missing catch-up delivery %s", key)
		}
	}
	if len(store.ensured) != 4 || !store.ensured[deliveryKey(2, time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC))] {
		t.Fatalf("delivery rows = %#v, want A June-August and B August only", store.ensured)
	}
}

func reportRows(prefix string, count int) []models.MonthlyVisitorReportRow {
	rows := make([]models.MonthlyVisitorReportRow, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, models.MonthlyVisitorReportRow{
			ID: int64(i + 1), VisitorName: fmt.Sprintf("%s-%d", prefix, i+1), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		})
	}
	return rows
}

func deliveryKey(societyID int64, month time.Time) string {
	return fmt.Sprintf("%d/%s", societyID, month.Format("2006-01"))
}
