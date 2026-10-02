package maintenancesvc

import (
	"go-server/internal/models"
	"testing"
	"time"
)

func TestCatchUpBoundaries(t *testing.T) {
	v := models.DefaultMaintenanceSettings()
	v.FirstEnabledMonth = "2026-07"
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	for _, m := range []string{"2026-07", "2026-08"} {
		if _, err := catchUpMonth(m, now, v); err != nil {
			t.Fatal(m, err)
		}
	}
	for _, m := range []string{"2026-06", "2026-09", "2026-10", "2026-7", ""} {
		if _, err := catchUpMonth(m, now, v); err == nil {
			t.Fatal("accepted", m)
		}
	}
}
func TestReviewHashBindsSettingsAndFinancialSnapshot(t *testing.T) {
	v := models.DefaultMaintenanceSettings()
	p := models.MaintenancePreview{Bills: []models.MaintenanceBill{{FlatID: 1, TotalPaise: 70000, DueDate: "2026-10-08"}}}
	original, err := reviewHash(v, p)
	if err != nil {
		t.Fatal(err)
	}
	p.Bills[0].DueDate = "2026-10-09"
	changed, _ := reviewHash(v, p)
	if changed == original {
		t.Fatal("due date not bound")
	}
	p.Bills[0].DueDate = "2026-10-08"
	v.FixedPaise = 80000
	changed, _ = reviewHash(v, p)
	if changed == original {
		t.Fatal("settings not bound")
	}
}
