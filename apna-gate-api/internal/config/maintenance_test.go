package config

import "testing"

func TestMaintenanceStatuses(t *testing.T) {
	for _, tc := range []struct {
		input   []string
		count   int
		invalid bool
	}{{nil, 3, false}, {[]string{"occupied"}, 1, false}, {[]string{"occupied", " occupied "}, 1, false}, {[]string{}, 0, true}, {[]string{""}, 0, true}, {[]string{"vacant", "bad"}, 0, true}} {
		c := Config{MaintenanceEligibleFlatStatuses: tc.input}
		got, err := c.MaintenanceStatuses()
		if (err != nil) != tc.invalid || len(got) != tc.count {
			t.Fatalf("%v => %v %v", tc.input, got, err)
		}
	}
	t.Setenv("MAINTENANCE_ELIGIBLE_FLAT_STATUSES", "")
	c := Config{MaintenanceEligibleFlatStatuses: maintenanceStatusesFromEnv()}
	if _, err := c.MaintenanceStatuses(); err == nil {
		t.Fatal("explicit empty env accepted")
	}
}
