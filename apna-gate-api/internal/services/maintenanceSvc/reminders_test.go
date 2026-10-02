package maintenancesvc

import (
	"testing"
	"time"
)

func TestReminderCalendar(t *testing.T) {
	for _, tc := range []struct{ now, due, zone, stage, date string }{
		{"2026-09-27T20:00:00Z", "2026-10-01", "Asia/Kolkata", "before_3", "2026-09-28"},
		{"2026-09-30T20:00:00Z", "2026-10-01", "Asia/Kolkata", "due", "2026-10-01"},
		{"2026-10-07T20:00:00Z", "2026-10-01", "Asia/Kolkata", "after_7", "2026-10-08"},
		{"2026-10-02T00:00:00Z", "2026-10-01", "Asia/Kolkata", "", ""},
		{"2026-10-01T01:00:00Z", "2026-10-01", "America/New_York", "", ""},
		{"2026-12-29T12:00:00Z", "2027-01-01", "UTC", "before_3", "2026-12-29"},
		{"2028-02-29T12:00:00Z", "2028-03-03", "UTC", "before_3", "2028-02-29"},
		{"2026-03-08T07:00:00Z", "2026-03-11", "America/New_York", "before_3", "2026-03-08"},
	} {
		now, _ := time.Parse(time.RFC3339, tc.now)
		due, _ := time.Parse("2006-01-02", tc.due)
		stage, date, err := reminderMilestone(now, due, tc.zone)
		if err != nil || stage != tc.stage || date != tc.date {
			t.Fatalf("%+v: %s %s %v", tc, stage, date, err)
		}
	}
	if _, _, err := reminderMilestone(time.Now(), time.Now(), "invalid"); err == nil {
		t.Fatal("invalid timezone accepted")
	}
}
