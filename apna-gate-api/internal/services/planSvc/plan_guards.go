package plansvc

import (
	"context"

	"go-server/internal/models"
)

// Private guards.

// requirePlan checks that a plan exists for the supplied filter.
func (s *PlanSvc) requirePlan(ctx context.Context, filter *models.PlanFilter) (*models.Plan, error) {
	plan, err := s.planRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}
	return plan, nil
}

// planResponseOrNotFound converts a nil command result into ErrPlanNotFound.
func planResponseOrNotFound(plan *models.Plan, err error) (*models.PlanResponse, error) {
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}
	return plan.ToResponse(), nil
}
