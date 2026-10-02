package jobs

import (
	"testing"
	"time"
)

func TestReportMonthUsesConfiguredTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	timestamp := time.Date(2026, time.August, 31, 20, 0, 0, 0, time.UTC)
	got := ReportMonth(timestamp, loc)
	want := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("ReportMonth() = %v, want %v", got, want)
	}
}

func TestMonthBoundsUseLocalMidnight(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	month := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	from, to := MonthBounds(month, loc)
	wantFrom := time.Date(2026, time.August, 31, 18, 30, 0, 0, time.UTC)
	wantTo := time.Date(2026, time.September, 30, 18, 30, 0, 0, time.UTC)
	if !from.Equal(wantFrom) || !to.Equal(wantTo) {
		t.Fatalf("MonthBounds() = (%v, %v), want (%v, %v)", from, to, wantFrom, wantTo)
	}
}

func TestLatestEligibleMonthWaitsUntilFiveMinutesPastMidnight(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	job := NewMonthlyVisitorReportJob(nil, nil, MonthlyVisitorReportJobConfig{Location: loc})
	before := time.Date(2026, time.September, 1, 0, 4, 59, 0, loc)
	after := time.Date(2026, time.September, 1, 0, 5, 0, 0, loc)
	if got := job.latestEligibleMonth(before); got.Month() != time.July {
		t.Fatalf("before threshold month = %v, want July", got)
	}
	if got := job.latestEligibleMonth(after); got.Month() != time.August {
		t.Fatalf("at threshold month = %v, want August", got)
	}
}
