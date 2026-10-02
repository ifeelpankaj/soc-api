package societysvc

import (
	"context"
	"reflect"
	"testing"

	"go-server/internal/models"
)

func TestGetOnboardingBootstrapReadiness(t *testing.T) {
	const societyID int64 = 42
	tests := []struct {
		name                   string
		societyStatus          models.SocietyStatus
		subscriptionStatus     *models.SubscriptionStatus
		activeFlats            int64
		activeStaff            int64
		wantActiveSubscription bool
		wantOnboarded          bool
		wantMissingSteps       []string
		wantPath               string
	}{
		{
			name:          "pending society stays in onboarding despite complete setup",
			societyStatus: models.SocietyStatusPending, subscriptionStatus: subscriptionStatus(models.SubscriptionStatusActive),
			activeFlats: 1, activeStaff: 1, wantActiveSubscription: true,
			wantMissingSteps: []string{}, wantPath: "/onboarding/soc16",
		},
		{
			name:          "active society without subscription stays in onboarding",
			societyStatus: models.SocietyStatusActive, activeFlats: 1, activeStaff: 1,
			wantMissingSteps: []string{}, wantPath: "/onboarding/soc16",
		},
		{
			name:          "expired subscription is not active",
			societyStatus: models.SocietyStatusActive, subscriptionStatus: subscriptionStatus(models.SubscriptionStatusExpired),
			activeFlats: 1, activeStaff: 1, wantMissingSteps: []string{}, wantPath: "/onboarding/soc16",
		},
		{
			name:          "operational society reports setup work only",
			societyStatus: models.SocietyStatusActive, subscriptionStatus: subscriptionStatus(models.SubscriptionStatusActive),
			wantActiveSubscription: true, wantMissingSteps: []string{"flats", "staff"}, wantPath: "/onboarding/soc16",
		},
		{
			name:          "active trial with flats and staff is onboarded",
			societyStatus: models.SocietyStatusActive, subscriptionStatus: subscriptionStatus(models.SubscriptionStatusTrial),
			activeFlats: 2, activeStaff: 3, wantActiveSubscription: true, wantOnboarded: true,
			wantMissingSteps: []string{}, wantPath: "/dashboard/soc16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var subscriptions []*models.SocietySubscriptionResponse
			if tt.subscriptionStatus != nil {
				subscriptions = []*models.SocietySubscriptionResponse{{SocietyID: societyID, Status: *tt.subscriptionStatus}}
			}
			svc := &SocietySvc{
				societyRepo: &dashboardSocietyRepo{society: &models.Society{ID: societyID, Status: tt.societyStatus}},
				flatRepo:    &dashboardFlatRepo{stats: &models.FlatStatsResponse{ActiveFlats: tt.activeFlats}},
				memberRepo: &dashboardMemberRepo{counts: map[string]int64{
					"staff:active": tt.activeStaff,
				}},
				subscriptionSvc: &dashboardSubscriptionSvc{subscriptions: subscriptions},
			}

			got, err := svc.GetOnboardingBootstrap(context.Background(), societyID)
			if err != nil {
				t.Fatalf("GetOnboardingBootstrap() error = %v", err)
			}
			if got.HasActiveSubscription != tt.wantActiveSubscription {
				t.Fatalf("has_active_subscription = %v, want %v", got.HasActiveSubscription, tt.wantActiveSubscription)
			}
			if got.IsOnboarded != tt.wantOnboarded {
				t.Fatalf("is_onboarded = %v, want %v", got.IsOnboarded, tt.wantOnboarded)
			}
			if !reflect.DeepEqual(got.MissingSteps, tt.wantMissingSteps) {
				t.Fatalf("missing_steps = %#v, want %#v", got.MissingSteps, tt.wantMissingSteps)
			}
			if got.NextPath != tt.wantPath {
				t.Fatalf("next_path = %q, want %q", got.NextPath, tt.wantPath)
			}
		})
	}
}

func TestEncodeSocietyRouteID(t *testing.T) {
	tests := map[int64]string{
		1:     "soc1",
		35:    "socz",
		36:    "soc10",
		42:    "soc16",
		12345: "soc9ix",
	}
	for id, want := range tests {
		if got := encodeSocietyRouteID(id); got != want {
			t.Fatalf("encodeSocietyRouteID(%d) = %q, want %q", id, got, want)
		}
	}
}

func subscriptionStatus(status models.SubscriptionStatus) *models.SubscriptionStatus {
	return &status
}
