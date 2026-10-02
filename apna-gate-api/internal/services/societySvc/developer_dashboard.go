package societysvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SocietySvc) GetDeveloperDashboardBootstrap(ctx context.Context) (*models.DeveloperDashboardBootstrapResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	societyStats, err := s.developerSocietyStats(ctx)
	if err != nil {
		return nil, err
	}
	planStats, err := s.developerPlanStats(ctx)
	if err != nil {
		return nil, err
	}
	subscriptionStats, err := s.developerSubscriptionStats(ctx)
	if err != nil {
		return nil, err
	}
	residenceStats, err := s.developerResidenceStats(ctx)
	if err != nil {
		return nil, err
	}
	recentPendingSocieties, err := s.developerRecentPendingSocieties(ctx)
	if err != nil {
		return nil, err
	}
	recentSubscriptions, err := s.developerRecentSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	return &models.DeveloperDashboardBootstrapResponse{
		SocietyStats:           societyStats,
		PlanStats:              planStats,
		SubscriptionStats:      subscriptionStats,
		ResidenceStats:         residenceStats,
		RecentPendingSocieties: recentPendingSocieties,
		RecentSubscriptions:    recentSubscriptions,
	}, nil
}

func (s *SocietySvc) developerSocietyStats(ctx context.Context) (*models.DeveloperDashboardSocietyStatsResponse, error) {
	pending := string(models.SocietyStatusPending)
	active := string(models.SocietyStatusActive)
	suspended := string(models.SocietyStatusSuspended)
	rejected := string(models.SocietyStatusRejected)

	total, err := s.societyRepo.Count(ctx, models.ListSocietiesFilter{})
	if err != nil {
		return nil, err
	}
	pendingCount, err := s.societyRepo.Count(ctx, models.ListSocietiesFilter{Status: &pending})
	if err != nil {
		return nil, err
	}
	activeCount, err := s.societyRepo.Count(ctx, models.ListSocietiesFilter{Status: &active})
	if err != nil {
		return nil, err
	}
	suspendedCount, err := s.societyRepo.Count(ctx, models.ListSocietiesFilter{Status: &suspended})
	if err != nil {
		return nil, err
	}
	rejectedCount, err := s.societyRepo.Count(ctx, models.ListSocietiesFilter{Status: &rejected})
	if err != nil {
		return nil, err
	}

	return &models.DeveloperDashboardSocietyStatsResponse{
		Total: total, Pending: pendingCount, Active: activeCount, Suspended: suspendedCount, Rejected: rejectedCount,
	}, nil
}

func (s *SocietySvc) developerPlanStats(ctx context.Context) (*models.DeveloperDashboardPlanStatsResponse, error) {
	if s.planSvc == nil {
		return &models.DeveloperDashboardPlanStatsResponse{}, nil
	}
	active := true
	inactive := false

	total, err := s.planSvc.CountPlans(ctx, &models.PlanFilter{})
	if err != nil {
		return nil, err
	}
	activeCount, err := s.planSvc.CountPlans(ctx, &models.PlanFilter{IsActive: &active})
	if err != nil {
		return nil, err
	}
	inactiveCount, err := s.planSvc.CountPlans(ctx, &models.PlanFilter{IsActive: &inactive})
	if err != nil {
		return nil, err
	}

	return &models.DeveloperDashboardPlanStatsResponse{Total: total, Active: activeCount, Inactive: inactiveCount}, nil
}

func (s *SocietySvc) developerSubscriptionStats(ctx context.Context) (*models.SubscriptionStatsResponse, error) {
	if s.subscriptionSvc == nil {
		return &models.SubscriptionStatsResponse{}, nil
	}
	return s.subscriptionSvc.GetSubscriptionStats(ctx, &models.SubscriptionFilter{})
}

func (s *SocietySvc) developerResidenceStats(ctx context.Context) (*models.DeveloperDashboardResidenceStatsResponse, error) {
	resident := string(models.SocietyMemberRoleResident)
	active := string(models.SocietyMemberStatusActive)

	total, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{Role: &resident})
	if err != nil {
		return nil, err
	}
	activeCount, err := s.memberRepo.Count(ctx, models.ListSocietyMembersFilter{Role: &resident, Status: &active})
	if err != nil {
		return nil, err
	}
	return &models.DeveloperDashboardResidenceStatsResponse{TotalResidents: total, ActiveResidents: activeCount}, nil
}

func (s *SocietySvc) developerRecentPendingSocieties(ctx context.Context) ([]*models.SocietyResponse, error) {
	pending := string(models.SocietyStatusPending)
	societies, err := s.societyRepo.List(ctx, models.ListSocietiesFilter{Status: &pending, Limit: 5})
	if err != nil {
		return nil, err
	}
	responses := make([]*models.SocietyResponse, 0, len(societies))
	for _, society := range societies {
		responses = append(responses, society.ToResponse())
	}
	return responses, nil
}

func (s *SocietySvc) developerRecentSubscriptions(ctx context.Context) ([]*models.SocietySubscriptionResponse, error) {
	if s.subscriptionSvc == nil {
		return []*models.SocietySubscriptionResponse{}, nil
	}
	return s.subscriptionSvc.ListSubscriptions(ctx, &models.SubscriptionFilter{Limit: 5})
}
