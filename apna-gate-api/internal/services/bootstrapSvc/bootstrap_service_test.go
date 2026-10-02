package bootstrapsvc

import (
	"context"
	"testing"

	"go-server/internal/models"
	societysvc "go-server/internal/services/societySvc"
)

type bootstrapSocietyService struct {
	societysvc.SocietyService
	onboarding      *models.SocietyOnboardingBootstrapResponse
	onboardingCalls int
}

func (s *bootstrapSocietyService) GetOnboardingBootstrap(context.Context, int64) (*models.SocietyOnboardingBootstrapResponse, error) {
	s.onboardingCalls++
	return s.onboarding, nil
}

func adminSociety(id int64, role models.SocietyMemberRole, memberStatus models.SocietyMemberStatus, societyStatus models.SocietyStatus) *models.MySocietyResponse {
	return &models.MySocietyResponse{
		Member:  &models.SocietyMemberResponse{SocietyID: id, Role: role, Status: memberStatus},
		Society: &models.SocietyResponse{ID: id, Status: societyStatus},
	}
}

func TestResolveDefaultDashboard(t *testing.T) {
	tests := []struct {
		name                string
		user                *models.UserResponse
		societies           []*models.MySocietyResponse
		onboarding          *models.SocietyOnboardingBootstrapResponse
		wantKind            models.DashboardKind
		wantPath            string
		wantSocietyID       int64
		wantSocietyIDSet    bool
		wantOnboardingCalls int
	}{
		{name: "developer gets developer dashboard", user: &models.UserResponse{GlobalRole: models.GlobalRoleDeveloper}, wantKind: models.DashboardKindDeveloper, wantPath: "/developer"},
		{name: "super admin gets developer dashboard", user: &models.UserResponse{GlobalRole: models.GlobalRoleSuperAdmin}, wantKind: models.DashboardKindDeveloper, wantPath: "/developer"},
		{name: "no membership gets onboarding", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser}, wantKind: models.DashboardKindOnboarding, wantPath: "/onboarding"},
		{
			name: "staff membership does not participate in admin routing", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies: []*models.MySocietyResponse{adminSociety(7, models.SocietyMemberRoleStaff, models.SocietyMemberStatusActive, models.SocietyStatusActive)},
			wantKind:  models.DashboardKindOnboarding, wantPath: "/onboarding",
		},
		{
			name: "inactive admin membership does not participate", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies: []*models.MySocietyResponse{adminSociety(8, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusSuspended, models.SocietyStatusActive)},
			wantKind:  models.DashboardKindOnboarding, wantPath: "/onboarding",
		},
		{
			name: "pending owner continues onboarding", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies:  []*models.MySocietyResponse{adminSociety(42, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive, models.SocietyStatusPending)},
			onboarding: &models.SocietyOnboardingBootstrapResponse{NextPath: "/onboarding/soc16"},
			wantKind:   models.DashboardKindOnboarding, wantPath: "/onboarding/soc16", wantSocietyID: 42, wantSocietyIDSet: true, wantOnboardingCalls: 1,
		},
		{
			name: "active admin without subscription continues onboarding", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies:  []*models.MySocietyResponse{adminSociety(42, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive, models.SocietyStatusActive)},
			onboarding: &models.SocietyOnboardingBootstrapResponse{NextPath: "/onboarding/soc16"},
			wantKind:   models.DashboardKindOnboarding, wantPath: "/onboarding/soc16", wantSocietyID: 42, wantSocietyIDSet: true, wantOnboardingCalls: 1,
		},
		{
			name: "configured operational society gets dashboard", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies:  []*models.MySocietyResponse{adminSociety(42, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive, models.SocietyStatusActive)},
			onboarding: &models.SocietyOnboardingBootstrapResponse{IsOnboarded: true, NextPath: "/dashboard/soc16"},
			wantKind:   models.DashboardKindSocietyAdmin, wantPath: "/dashboard/soc16", wantSocietyID: 42, wantSocietyIDSet: true, wantOnboardingCalls: 1,
		},
		{
			name: "single rejected society uses selector", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies: []*models.MySocietyResponse{adminSociety(42, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive, models.SocietyStatusRejected)},
			wantKind:  models.DashboardKindSelectSociety, wantPath: "/select-society",
		},
		{
			name: "single suspended society uses selector", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies: []*models.MySocietyResponse{adminSociety(42, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive, models.SocietyStatusSuspended)},
			wantKind:  models.DashboardKindSelectSociety, wantPath: "/select-society",
		},
		{
			name: "multiple admin societies use selector without readiness queries", user: &models.UserResponse{GlobalRole: models.GlobalRoleUser},
			societies: []*models.MySocietyResponse{
				adminSociety(42, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive, models.SocietyStatusPending),
				adminSociety(84, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive, models.SocietyStatusActive),
			},
			wantKind: models.DashboardKindSelectSociety, wantPath: "/select-society",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			societySvc := &bootstrapSocietyService{onboarding: tt.onboarding}
			svc := &bootstrapService{societySvc: societySvc}
			got, err := svc.resolveDefaultDashboard(context.Background(), tt.user, tt.societies)
			if err != nil {
				t.Fatalf("resolveDefaultDashboard() error = %v", err)
			}
			if got == nil {
				t.Fatal("expected dashboard, got nil")
			}
			if got.Kind != tt.wantKind {
				t.Fatalf("kind = %q, want %q", got.Kind, tt.wantKind)
			}
			if got.Path != tt.wantPath {
				t.Fatalf("path = %q, want %q", got.Path, tt.wantPath)
			}
			if tt.wantSocietyIDSet {
				if got.SocietyID == nil || *got.SocietyID != tt.wantSocietyID {
					t.Fatalf("society_id = %v, want %d", got.SocietyID, tt.wantSocietyID)
				}
			} else if got.SocietyID != nil {
				t.Fatalf("society_id = %d, want nil", *got.SocietyID)
			}
			if societySvc.onboardingCalls != tt.wantOnboardingCalls {
				t.Fatalf("onboarding calls = %d, want %d", societySvc.onboardingCalls, tt.wantOnboardingCalls)
			}
		})
	}
}
