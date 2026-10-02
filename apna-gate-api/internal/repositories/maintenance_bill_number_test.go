package repository

import (
	"go-server/internal/models"
	"testing"
	"time"
)

func TestMaintenanceBillNumberUsesSocietyTimeAndFlat(t *testing.T) {
	block := "west-side"
	issuedAt := time.Date(2026, 9, 30, 16, 22, 28, 0, time.FixedZone("IST", 5*3600+30*60))
	flat := models.MaintenanceFlat{Block: &block, FlatNumber: "007"}
	if got, want := maintenanceBillNumber("SOCI122H", flat, issuedAt), "AG-SOCI122H-20260930162228-west-side-007"; got != want {
		t.Fatalf("bill number = %q, want %q", got, want)
	}
}
