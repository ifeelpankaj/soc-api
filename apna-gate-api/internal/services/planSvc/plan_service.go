package plansvc

import (
	"context"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
)

type PlanService interface {
	PlanCommandService
	PlanQueryService
}

type PlanCommandService interface {
	CreatePlan(ctx context.Context, req *models.CreatePlanRequest) (*models.PlanResponse, error)
	UpdatePlan(ctx context.Context, planID int64, req *models.UpdatePlanRequest) (*models.PlanResponse, error)
	ActivatePlan(ctx context.Context, planID int64) (*models.PlanResponse, error)
	DeactivatePlan(ctx context.Context, planID int64) (*models.PlanResponse, error)
}

type PlanQueryService interface {
	GetPlan(ctx context.Context, filter *models.PlanFilter) (*models.PlanResponse, error)
	ListPlans(ctx context.Context, filter *models.PlanFilter) ([]*models.PlanResponse, error)
	CountPlans(ctx context.Context, filter *models.PlanFilter) (int64, error)
}

type PlanSvc struct {
	planRepo repository.PlanRepository
}

func NewPlanService(planRepo repository.PlanRepository) PlanService {
	return &PlanSvc{planRepo: planRepo}
}

var _ PlanService = (*PlanSvc)(nil)
