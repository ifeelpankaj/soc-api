package models

import "testing"

func TestParseFlatArea(t *testing.T) {
	for input, want := range map[string]int64{"1000": 100000, "1000.5": 100050, "0.01": 1, "92233720368547758.07": 9223372036854775807} {
		got, err := ParseFlatArea(input)
		if err != nil || got != want {
			t.Fatalf("%s: %d %v", input, got, err)
		}
	}
	for _, input := range []string{"", "0", "-1", "+1", "1e3", "1.234", "1.", ".1", " 1", "1.2.3", "92233720368547758.08"} {
		if _, err := ParseFlatArea(input); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
}
func TestMaintenanceSettingsValidation(t *testing.T) {
	s := DefaultMaintenanceSettings()
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Enabled = true
	if err := s.Validate(); err == nil {
		t.Fatal("enabled zero rate accepted")
	}
	s.FixedPaise = 200000
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, modify := range []func(*MaintenanceSettings){func(v *MaintenanceSettings) { v.Currency = "USD" }, func(v *MaintenanceSettings) { v.BillingDay = 29 }, func(v *MaintenanceSettings) { v.DueDay = 0 }, func(v *MaintenanceSettings) { v.Timezone = "bad/zone" }, func(v *MaintenanceSettings) { v.PricingModel = "unknown" }, func(v *MaintenanceSettings) { v.TypeRates = map[string]int64{" 1 BHK": 100} }} {
		v := s
		modify(&v)
		if err := v.Validate(); err == nil {
			t.Fatalf("invalid settings accepted: %+v", v)
		}
	}
}
