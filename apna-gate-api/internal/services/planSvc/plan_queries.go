package plansvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *PlanSvc) GetPlan(ctx context.Context, filter *models.PlanFilter) (*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	plan, err := s.requirePlan(ctx, filter)
	if err != nil {
		return nil, err
	}
	return plan.ToResponse(), nil
}

func (s *PlanSvc) ListPlans(ctx context.Context, filter *models.PlanFilter) ([]*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	plans, err := s.planRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	responses := make([]*models.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		responses = append(responses, plan.ToResponse())
	}
	return responses, nil
}

func (s *PlanSvc) CountPlans(ctx context.Context, filter *models.PlanFilter) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.planRepo.Count(ctx, filter)
}
