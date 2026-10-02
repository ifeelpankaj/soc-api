package societysvc

import (
	"context"
	"go-server/internal/repositories/contracts"
	"testing"
	"time"

	"go-server/internal/models"
)

func TestGetDashboardBootstrap(t *testing.T) {
	ctx := context.Background()
	societyID := int64(42)
	now := time.Now()

	tests := []struct {
		name                   string
		subscriptions          []*models.SocietySubscriptionResponse
		plans                  []*models.PlanResponse
		claimStats             *models.FlatClaimStatsResponse
		claims                 []*models.FlatClaim
		wantSubscription       bool
		wantUsage              bool
		wantPendingClaims      int
		wantPendingClaimCount  int64
		wantPlanAds            int
		wantTotalActiveMembers int64
	}{
		{
			name: "no subscription still returns summary and plan ads",
			plans: []*models.PlanResponse{
				{ID: 1, Name: "Starter", Code: "STARTER", PriceAmountPaise: 99900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 50, MaxAdmins: 2, MaxStaff: 4, MaxResidents: 100},
				{ID: 2, Name: "Growth", Code: "GROWTH", PriceAmountPaise: 199900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 150, MaxAdmins: 5, MaxStaff: 10, MaxResidents: 300},
			},
			claimStats:             &models.FlatClaimStatsResponse{},
			wantPlanAds:            2,
			wantTotalActiveMembers: 7,
		},
		{
			name: "active subscription includes usage and upgrade plans",
			subscriptions: []*models.SocietySubscriptionResponse{
				{ID: 10, SocietyID: societyID, Status: models.SubscriptionStatusActive, PlanName: "Starter", PlanCode: "STARTER", PriceAmountPaise: 99900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 50, MaxAdmins: 2, MaxStaff: 4, MaxResidents: 100, StartsAt: &now},
			},
			plans: []*models.PlanResponse{
				{ID: 1, Name: "Starter", Code: "STARTER", PriceAmountPaise: 99900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 50, MaxAdmins: 2, MaxStaff: 4, MaxResidents: 100},
				{ID: 2, Name: "Growth", Code: "GROWTH", PriceAmountPaise: 199900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 150, MaxAdmins: 5, MaxStaff: 10, MaxResidents: 300},
			},
			claimStats:             &models.FlatClaimStatsResponse{},
			wantSubscription:       true,
			wantUsage:              true,
			wantPlanAds:            1,
			wantTotalActiveMembers: 7,
		},
		{
			name: "pending claims are counted and returned",
			plans: []*models.PlanResponse{
				{ID: 1, Name: "Starter", Code: "STARTER", PriceAmountPaise: 99900, Currency: "INR", BillingCycle: models.BillingCycleMonthly, MaxFlats: 50, MaxAdmins: 2, MaxStaff: 4, MaxResidents: 100},
			},
			claimStats: &models.FlatClaimStatsResponse{TotalClaims: 3, PendingClaims: 2, ApprovedClaims: 1},
			claims: []*models.FlatClaim{
				{ID: 1, SocietyID: societyID, FlatID: 10, UserID: 100, RequestedRole: models.FlatResidentRoleOwner, Status: models.FlatClaimStatusPending, CreatedAt: now},
				{ID: 2, SocietyID: societyID, FlatID: 11, UserID: 101, RequestedRole: models.FlatResidentRoleTenant, Status: models.FlatClaimStatusPending, CreatedAt: now},
			},
			wantPendingClaims:      2,
			wantPendingClaimCount:  2,
			wantPlanAds:            1,
			wantTotalActiveMembers: 7,
		},
		{
			name:                   "no active plans returns empty plan ads",
			claimStats:             &models.FlatClaimStatsResponse{},
			wantPlanAds:            0,
			wantTotalActiveMembers: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &SocietySvc{
				societyRepo: &dashboardSocietyRepo{society: &models.Society{
					ID: societyID, Name: "Green Heights", SocietyCode: "GH01", Status: models.SocietyStatusActive, TotalFlats: 100, CreatedBy: 1,
				}},
				memberRepo: &dashboardMemberRepo{
					counts: map[string]int64{
						"active":          7,
						"owner:active":    1,
						"admin:active":    1,
						"staff:active":    2,
						"resident:active": 3,
					},
				},
				flatRepo: &dashboardFlatRepo{stats: &models.FlatStatsResponse{
					SocietyID: societyID, TotalFlats: 80, VacantFlats: 20, OccupiedFlats: 55, BlockedFlats: 5, ActiveFlats: 75, InactiveFlats: 5,
				}},
				flatClaimRepo:   &dashboardClaimRepo{stats: tt.claimStats, claims: tt.claims},
				planSvc:         &dashboardPlanSvc{plans: tt.plans},
				subscriptionSvc: &dashboardSubscriptionSvc{subscriptions: tt.subscriptions},
			}

			got, err := svc.GetDashboardBootstrap(ctx, societyID)
			if err != nil {
				t.Fatalf("GetDashboardBootstrap() error = %v", err)
			}
			if got.Society == nil || got.Society.ID != societyID {
				t.Fatalf("society = %#v, want society %d", got.Society, societyID)
			}
			if got.MemberStats.TotalActiveMembers != tt.wantTotalActiveMembers {
				t.Fatalf("total active members = %d, want %d", got.MemberStats.TotalActiveMembers, tt.wantTotalActiveMembers)
			}
			if (got.CurrentSubscription != nil) != tt.wantSubscription {
				t.Fatalf("current subscription set = %v, want %v", got.CurrentSubscription != nil, tt.wantSubscription)
			}
			if (got.SubscriptionUsage != nil) != tt.wantUsage {
				t.Fatalf("subscription usage set = %v, want %v", got.SubscriptionUsage != nil, tt.wantUsage)
			}
			if tt.wantUsage && got.SubscriptionUsage.Residents.Used != 3 {
				t.Fatalf("resident usage = %d, want 3", got.SubscriptionUsage.Residents.Used)
			}
			if got.ClaimStats.PendingClaims != tt.wantPendingClaimCount {
				t.Fatalf("pending claim count = %d, want %d", got.ClaimStats.PendingClaims, tt.wantPendingClaimCount)
			}
			if len(got.RecentPendingClaims) != tt.wantPendingClaims {
				t.Fatalf("recent pending claims = %d, want %d", len(got.RecentPendingClaims), tt.wantPendingClaims)
			}
			if len(got.PlanAds) != tt.wantPlanAds {
				t.Fatalf("plan ads = %d, want %d", len(got.PlanAds), tt.wantPlanAds)
			}
		})
	}
}

func TestGetDeveloperDashboardBootstrap(t *testing.T) {
	ctx := context.Background()
	pending := string(models.SocietyStatusPending)
	active := string(models.SocietyStatusActive)
	suspended := string(models.SocietyStatusSuspended)
	rejected := string(models.SocietyStatusRejected)
	resident := string(models.SocietyMemberRoleResident)
	activeMember := string(models.SocietyMemberStatusActive)

	svc := &SocietySvc{
		societyRepo: &dashboardSocietyRepo{
			societies: []*models.Society{
				{ID: 1, Name: "Pending Society", SocietyCode: "PS01", Status: models.SocietyStatusPending, CreatedBy: 1},
			},
			counts: map[string]int64{
				"":        4,
				pending:   1,
				active:    1,
				suspended: 1,
				rejected:  1,
			},
		},
		memberRepo: &dashboardMemberRepo{
			counts: map[string]int64{
				resident:                      12,
				resident + ":" + activeMember: 9,
			},
		},
		planSvc: &dashboardPlanSvc{
			plans: []*models.PlanResponse{
				{ID: 1, IsActive: true, MaxResidents: 100},
				{ID: 2, IsActive: false, MaxResidents: 200},
			},
		},
		subscriptionSvc: &dashboardSubscriptionSvc{
			subscriptions: []*models.SocietySubscriptionResponse{
				{ID: 1, Status: models.SubscriptionStatusActive, MaxResidents: 100},
				{ID: 2, Status: models.SubscriptionStatusExpired, MaxResidents: 200},
			},
		},
	}

	got, err := svc.GetDeveloperDashboardBootstrap(ctx)
	if err != nil {
		t.Fatalf("GetDeveloperDashboardBootstrap() error = %v", err)
	}
	if got.SocietyStats.Total != 4 || got.SocietyStats.Pending != 1 {
		t.Fatalf("society stats = %#v", got.SocietyStats)
	}
	if got.PlanStats.Total != 2 || got.PlanStats.Active != 1 || got.PlanStats.Inactive != 1 {
		t.Fatalf("plan stats = %#v", got.PlanStats)
	}
	if got.SubscriptionStats.TotalSubscriptions != 2 || got.SubscriptionStats.ActiveSubscriptions != 1 {
		t.Fatalf("subscription stats = %#v", got.SubscriptionStats)
	}
	if got.ResidenceStats.TotalResidents != 12 || got.ResidenceStats.ActiveResidents != 9 {
		t.Fatalf("residence stats = %#v", got.ResidenceStats)
	}
	if len(got.RecentPendingSocieties) != 1 || len(got.RecentSubscriptions) != 2 {
		t.Fatalf("recent lists = %d societies, %d subscriptions", len(got.RecentPendingSocieties), len(got.RecentSubscriptions))
	}
}

type dashboardSocietyRepo struct {
	society   *models.Society
	societies []*models.Society
	counts    map[string]int64
}

func (r *dashboardSocietyRepo) Create(ctx context.Context, society *models.Society) error {
	panic("unused")
}
func (r *dashboardSocietyRepo) Get(ctx context.Context, filter models.GetSocietyFilter) (*models.Society, error) {
	return r.society, nil
}
func (r *dashboardSocietyRepo) List(ctx context.Context, filter models.ListSocietiesFilter) ([]*models.Society, error) {
	return r.societies, nil
}
func (r *dashboardSocietyRepo) Count(ctx context.Context, filter models.ListSocietiesFilter) (int64, error) {
	if r.counts == nil {
		return 0, nil
	}
	if filter.Status == nil {
		return r.counts[""], nil
	}
	return r.counts[*filter.Status], nil
}
func (r *dashboardSocietyRepo) Update(ctx context.Context, societyID int64, req contracts.UpdateSocietyInput) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) Approve(ctx context.Context, societyID int64, approvedBy int64) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) Reject(ctx context.Context, societyID int64, rejectedBy int64, reason string) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) Suspend(ctx context.Context, societyID int64, suspendedBy int64, reason string) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) Reactivate(ctx context.Context, societyID int64, reactivatedBy int64) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) Restore(ctx context.Context, societyID int64) (*models.Society, error) {
	panic("unused")
}
func (r *dashboardSocietyRepo) SoftDelete(ctx context.Context, societyID int64) error {
	panic("unused")
}
func (r *dashboardSocietyRepo) CountPendingByCreator(ctx context.Context, createdBy int64) (int64, error) {
	panic("unused")
}

type dashboardMemberRepo struct {
	counts map[string]int64
}

func (r *dashboardMemberRepo) Add(ctx context.Context, member *models.SocietyMember) error {
	panic("unused")
}
func (r *dashboardMemberRepo) Get(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) List(ctx context.Context, filter models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) ListByUser(ctx context.Context, userID int64) ([]*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) ListMySocietiesByUser(ctx context.Context, userID int64) ([]*models.MySocietyResponse, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) Count(ctx context.Context, filter models.ListSocietyMembersFilter) (int64, error) {
	if filter.Role == nil && filter.Status != nil {
		return r.counts[*filter.Status], nil
	}
	if filter.Role != nil && filter.Status != nil {
		return r.counts[*filter.Role+":"+*filter.Status], nil
	}
	if filter.Role != nil {
		return r.counts[*filter.Role], nil
	}
	return 0, nil
}
func (r *dashboardMemberRepo) ChangeRole(ctx context.Context, societyID int64, userID int64, role models.SocietyMemberRole) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) Suspend(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) Reactivate(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) Remove(ctx context.Context, societyID int64, userID int64, removedBy int64, reason string) error {
	panic("unused")
}
func (r *dashboardMemberRepo) CountActiveOwners(ctx context.Context, societyID int64) (int64, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) DemoteActiveOwners(ctx context.Context, societyID int64, exceptUserID int64) error {
	panic("unused")
}
func (r *dashboardMemberRepo) PromoteToOwner(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *dashboardMemberRepo) UpsertResident(ctx context.Context, societyID int64, userID int64, invitedBy int64) (*models.SocietyMember, error) {
	panic("unused")
}

type dashboardFlatRepo struct {
	stats *models.FlatStatsResponse
}

func (r *dashboardFlatRepo) Create(ctx context.Context, flat *models.Flat) error { panic("unused") }
func (r *dashboardFlatRepo) Get(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) List(ctx context.Context, filter *models.FlatFilter) ([]*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) Count(ctx context.Context, filter *models.FlatFilter) (int64, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) Stats(ctx context.Context, societyID int64) (*models.FlatStatsResponse, error) {
	return r.stats, nil
}
func (r *dashboardFlatRepo) Update(ctx context.Context, filter *models.FlatFilter, req *contracts.UpdateFlatInput) (*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) Deactivate(ctx context.Context, filter *models.FlatFilter) error {
	panic("unused")
}
func (r *dashboardFlatRepo) Block(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) Unblock(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) MarkOccupied(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error) {
	panic("unused")
}
func (r *dashboardFlatRepo) MarkVacant(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error) {
	panic("unused")
}

type dashboardClaimRepo struct {
	stats  *models.FlatClaimStatsResponse
	claims []*models.FlatClaim
}

func (r *dashboardClaimRepo) Submit(ctx context.Context, claim *models.FlatClaim) error {
	panic("unused")
}
func (r *dashboardClaimRepo) Get(ctx context.Context, filter *models.FlatClaimFilter) (*models.FlatClaim, error) {
	panic("unused")
}
func (r *dashboardClaimRepo) List(ctx context.Context, filter *models.FlatClaimFilter) ([]*models.FlatClaim, error) {
	return r.claims, nil
}
func (r *dashboardClaimRepo) Stats(ctx context.Context, societyID int64) (*models.FlatClaimStatsResponse, error) {
	if r.stats == nil {
		return &models.FlatClaimStatsResponse{}, nil
	}
	return r.stats, nil
}
func (r *dashboardClaimRepo) Approve(ctx context.Context, societyID int64, claimID int64, reviewedBy int64) (*models.FlatClaim, error) {
	panic("unused")
}
func (r *dashboardClaimRepo) Reject(ctx context.Context, societyID int64, claimID int64, reviewedBy int64, reason string) (*models.FlatClaim, error) {
	panic("unused")
}
func (r *dashboardClaimRepo) Cancel(ctx context.Context, claimID int64, userID int64) (*models.FlatClaim, error) {
	panic("unused")
}

type dashboardPlanSvc struct {
	plans []*models.PlanResponse
}

func (s *dashboardPlanSvc) ListPlans(ctx context.Context, filter *models.PlanFilter) ([]*models.PlanResponse, error) {
	return s.plans, nil
}
func (s *dashboardPlanSvc) CountPlans(ctx context.Context, filter *models.PlanFilter) (int64, error) {
	if filter == nil || filter.IsActive == nil {
		return int64(len(s.plans)), nil
	}
	count := int64(0)
	for _, plan := range s.plans {
		if plan.IsActive == *filter.IsActive {
			count++
		}
	}
	return count, nil
}

type dashboardSubscriptionSvc struct {
	subscriptions []*models.SocietySubscriptionResponse
}

func (s *dashboardSubscriptionSvc) CanAddAdmin(ctx context.Context, societyID int64, adding int64) error {
	return nil
}
func (s *dashboardSubscriptionSvc) CanAddStaff(ctx context.Context, societyID int64, adding int64) error {
	return nil
}
func (s *dashboardSubscriptionSvc) ListSubscriptions(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscriptionResponse, error) {
	return s.subscriptions, nil
}
func (s *dashboardSubscriptionSvc) GetSubscriptionStats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	stats := &models.SubscriptionStatsResponse{TotalSubscriptions: int64(len(s.subscriptions))}
	for _, subscription := range s.subscriptions {
		switch subscription.Status {
		case models.SubscriptionStatusPending:
			stats.PendingSubscriptions++
		case models.SubscriptionStatusTrial:
			stats.TrialSubscriptions++
		case models.SubscriptionStatusActive:
			stats.ActiveSubscriptions++
		case models.SubscriptionStatusExpired:
			stats.ExpiredSubscriptions++
		case models.SubscriptionStatusCancelled:
			stats.CancelledSubscriptions++
		}
	}
	return stats, nil
}
