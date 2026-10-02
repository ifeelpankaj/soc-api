package plansvc

import (
	"context"
	"go-server/internal/repositories/contracts"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *PlanSvc) CreatePlan(ctx context.Context, req *models.CreatePlanRequest) (*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidPlanRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidPlanRequest.WithCause(err)
	}
	plan := &models.Plan{
		Name: req.Name, Code: req.Code, Description: req.Description,
		PriceAmountPaise: req.PriceAmountPaise, Currency: req.Currency,
		BillingCycle: req.BillingCycle, MaxFlats: req.MaxFlats,
		MaxAdmins: req.MaxAdmins, MaxStaff: req.MaxStaff, MaxResidents: req.MaxResidents, Features: req.Features, IsActive: true,
	}
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, ErrPlanConflict.WithCause(err)
	}
	return plan.ToResponse(), nil
}

func (s *PlanSvc) UpdatePlan(ctx context.Context, planID int64, req *models.UpdatePlanRequest) (*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if req == nil {
		return nil, ErrInvalidPlanRequest
	}
	req.Sanitize()
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidPlanRequest.WithCause(err)
	}
	return planResponseOrNotFound(s.planRepo.Update(ctx, planID, (*contracts.UpdatePlanInput)(req)))
}

func (s *PlanSvc) ActivatePlan(ctx context.Context, planID int64) (*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return planResponseOrNotFound(s.planRepo.Activate(ctx, planID))
}

func (s *PlanSvc) DeactivatePlan(ctx context.Context, planID int64) (*models.PlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return planResponseOrNotFound(s.planRepo.Deactivate(ctx, planID))
}
