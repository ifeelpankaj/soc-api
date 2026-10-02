package societysvc

import (
	"testing"
	"time"

	"go-server/internal/models"
)

func TestComputeSubscriptionHealthActive(t *testing.T) {
	endsAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	current := &models.SocietySubscriptionResponse{
		Status: models.SubscriptionStatusActive,
		EndsAt: &endsAt,
	}

	health := computeSubscriptionHealth(current, current)

	if !health.IsActive {
		t.Fatalf("expected active subscription")
	}
	if health.IsExpiringSoon {
		t.Fatalf("did not expect expiring soon within 30 days window")
	}
	if health.LifecycleLabel != "active" {
		t.Fatalf("expected lifecycle active, got %q", health.LifecycleLabel)
	}
	if health.DaysUntilExpiry == nil || *health.DaysUntilExpiry < 29 {
		t.Fatalf("expected days until expiry around 30, got %v", health.DaysUntilExpiry)
	}
}

func TestComputeSubscriptionHealthExpiringSoon(t *testing.T) {
	endsAt := time.Now().UTC().Add(5 * 24 * time.Hour)
	current := &models.SocietySubscriptionResponse{
		Status: models.SubscriptionStatusActive,
		EndsAt: &endsAt,
	}

	health := computeSubscriptionHealth(current, current)

	if !health.IsExpiringSoon {
		t.Fatalf("expected expiring soon")
	}
	if health.LifecycleLabel != "expiring_soon" {
		t.Fatalf("expected lifecycle expiring_soon, got %q", health.LifecycleLabel)
	}
}

func TestComputeSubscriptionHealthExpiredByDate(t *testing.T) {
	endsAt := time.Now().UTC().Add(-24 * time.Hour)
	latest := &models.SocietySubscriptionResponse{
		Status: models.SubscriptionStatusActive,
		EndsAt: &endsAt,
	}

	health := computeSubscriptionHealth(nil, latest)

	if health.IsActive {
		t.Fatalf("expected inactive subscription")
	}
	if health.LifecycleLabel != "expired" {
		t.Fatalf("expected lifecycle expired, got %q", health.LifecycleLabel)
	}
}

func TestComputeSubscriptionHealthNone(t *testing.T) {
	health := computeSubscriptionHealth(nil, nil)

	if health.LifecycleLabel != "none" {
		t.Fatalf("expected lifecycle none, got %q", health.LifecycleLabel)
	}
}
