package societysvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SocietySvc) GetDashboardBootstrap(ctx context.Context, societyID int64) (*models.SocietyDashboardBootstrapResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrSocietyNotFound
	}

	flatStats, err := s.flatRepo.Stats(ctx, societyID)
	if err != nil {
		return nil, err
	}

	claimStats := &models.FlatClaimStatsResponse{}
	recentPendingClaims := []*models.FlatClaimResponse{}
	if s.flatClaimRepo != nil {
		claimStats, err = s.flatClaimRepo.Stats(ctx, societyID)
		if err != nil {
			return nil, err
		}

		pending := string(models.FlatClaimStatusPending)
		claims, err := s.flatClaimRepo.List(ctx, &models.FlatClaimFilter{
			SocietyID: &societyID,
			Status:    &pending,
			Limit:     5,
		})
		if err != nil {
			return nil, err
		}
		recentPendingClaims = make([]*models.FlatClaimResponse, 0, len(claims))
		for _, claim := range claims {
			recentPendingClaims = append(recentPendingClaims, claim.ToResponse())
		}
	}

	memberStats, err := s.dashboardMemberStats(ctx, societyID)
	if err != nil {
		return nil, err
	}

	currentSubscription, err := s.dashboardCurrentSubscription(ctx, societyID)
	if err != nil {
		return nil, err
	}

	latestSubscription, err := s.dashboardLatestSubscription(ctx, societyID)
	if err != nil {
		return nil, err
	}

	subscriptionHealth := computeSubscriptionHealth(currentSubscription, latestSubscription)

	var subscriptionUsage *models.SocietyDashboardSubscriptionUsageResponse
	if currentSubscription != nil {
		subscriptionUsage = &models.SocietyDashboardSubscriptionUsageResponse{
			Flats:     quotaUsage(flatStats.ActiveFlats, int64(currentSubscription.MaxFlats)),
			Admins:    quotaUsage(memberStats.Owners+memberStats.Admins, int64(currentSubscription.MaxAdmins)),
			Staff:     quotaUsage(memberStats.Staff, int64(currentSubscription.MaxStaff)),
			Residents: quotaUsage(memberStats.Residents, int64(currentSubscription.MaxResidents)),
		}
	}

	planAds, err := s.dashboardPlanAds(ctx, currentSubscription)
	if err != nil {
		return nil, err
	}

	var visitorStats *models.VisitorEntryStatsResponse
	var visitorDaily []models.VisitorDailyCountResponse
	if s.visitorEntrySvc != nil {
		visitorStats, err = s.visitorEntrySvc.GetEntryStats(ctx, societyID)
		if err != nil {
			return nil, err
		}
		dailyStats, dailyErr := s.visitorEntrySvc.GetDailyEntryStats(ctx, societyID, 7)
		if dailyErr != nil {
			return nil, dailyErr
		}
		if dailyStats != nil {
			visitorDaily = dailyStats.Daily
		}
	}

	return &models.SocietyDashboardBootstrapResponse{
		Society:             society.ToResponse(),
		FlatStats:           flatStats,
		ClaimStats:          claimStats,
		RecentPendingClaims: recentPendingClaims,
		MemberStats:         memberStats,
		CurrentSubscription: currentSubscription,
		SubscriptionUsage:   subscriptionUsage,
		SubscriptionHealth:  subscriptionHealth,
		PlanAds:             planAds,
		VisitorStats:        visitorStats,
		VisitorDailyLast7:   visitorDaily,
	}, nil
}

func (s *SocietySvc) dashboardMemberStats(ctx context.Context, societyID int64) (*models.SocietyDashboardMemberStatsResponse, error) {
	active := string(models.SocietyMemberStatusActive)
	owner := string(models.SocietyMemberRoleOwner)
	admin := string(models.SocietyMemberRoleAdmin)
	staff := string(models.SocietyMemberRoleStaff)
	resident := string(models.SocietyMemberRoleResident)

	total, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: societyID, Status: &active})
	if err != nil {
		return nil, err
	}
	owners, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: societyID, Role: &owner, Status: &active})
	if err != nil {
		return nil, err
	}
	admins, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: societyID, Role: &admin, Status: &active})
	if err != nil {
		return nil, err
	}
	staffCount, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: societyID, Role: &staff, Status: &active})
	if err != nil {
		return nil, err
	}
	residents, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{SocietyID: societyID, Role: &resident, Status: &active})
	if err != nil {
		return nil, err
	}

	return &models.SocietyDashboardMemberStatsResponse{
		TotalActiveMembers: total,
		Owners:             owners,
		Admins:             admins,
		Staff:              staffCount,
		Residents:          residents,
	}, nil
}

func (s *SocietySvc) dashboardCurrentSubscription(ctx context.Context, societyID int64) (*models.SocietySubscriptionResponse, error) {
	if s.subscriptionSvc == nil {
		return nil, nil
	}
	activeOnly := true
	items, err := s.subscriptionSvc.ListSubscriptions(ctx, &models.SubscriptionFilter{
		SocietyID:    &societyID,
		IsActiveOnly: &activeOnly,
		Limit:        1,
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}

func (s *SocietySvc) dashboardLatestSubscription(ctx context.Context, societyID int64) (*models.SocietySubscriptionResponse, error) {
	if s.subscriptionSvc == nil {
		return nil, nil
	}
	items, err := s.subscriptionSvc.ListSubscriptions(ctx, &models.SubscriptionFilter{
		SocietyID: &societyID,
		Limit:     1,
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}

func (s *SocietySvc) dashboardPlanAds(ctx context.Context, current *models.SocietySubscriptionResponse) ([]*models.PlanResponse, error) {
	if s.planSvc == nil {
		return []*models.PlanResponse{}, nil
	}
	active := true
	plans, err := s.planSvc.ListPlans(ctx, &models.PlanFilter{IsActive: &active, Limit: 20})
	if err != nil {
		return nil, err
	}
	ads := make([]*models.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		if plan == nil {
			continue
		}
		if current != nil && plan.Code == current.PlanCode {
			continue
		}
		if current != nil && !isUpgradePlan(plan, current) {
			continue
		}
		ads = append(ads, plan)
	}
	return ads, nil
}

func isUpgradePlan(plan *models.PlanResponse, current *models.SocietySubscriptionResponse) bool {
	return plan.PriceAmountPaise > current.PriceAmountPaise ||
		plan.MaxFlats > current.MaxFlats ||
		plan.MaxAdmins > current.MaxAdmins ||
		plan.MaxStaff > current.MaxStaff ||
		plan.MaxResidents > current.MaxResidents
}

func quotaUsage(used int64, limit int64) models.SocietyDashboardQuotaUsageResponse {
	remaining := limit - used
	if remaining < 0 {
		remaining = 0
	}
	percent := int64(0)
	if limit > 0 {
		percent = used * 100 / limit
		if percent > 100 {
			percent = 100
		}
	}
	return models.SocietyDashboardQuotaUsageResponse{
		Used:      used,
		Limit:     limit,
		Remaining: remaining,
		Percent:   percent,
	}
}
