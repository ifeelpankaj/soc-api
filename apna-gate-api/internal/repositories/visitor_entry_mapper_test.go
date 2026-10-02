package repository

import "testing"

func TestFlatSummaryReturnsNilForNullFlatID(t *testing.T) {
	if got := flatSummary(0, "G-01", nil, nil); got != nil {
		t.Fatalf("expected nil flat summary when flat_id is null")
	}
}

func TestFlatSummaryReturnsNilWhenFlatNumberMissing(t *testing.T) {
	if got := flatSummary(10, "", nil, nil); got != nil {
		t.Fatalf("expected nil flat summary when flat number is empty")
	}
}

func TestFlatSummaryReturnsSummaryForValidFlat(t *testing.T) {
	block := "A"
	got := flatSummary(10, "G-01", &block, nil)
	if got == nil || got.ID != 10 || got.FlatNumber != "G-01" {
		t.Fatalf("expected flat summary, got %+v", got)
	}
}
