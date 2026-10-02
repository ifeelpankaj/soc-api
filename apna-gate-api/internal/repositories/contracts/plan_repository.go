package contracts

import (
	"context"
	"go-server/internal/models"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *models.Plan) error
	Get(ctx context.Context, filter *models.PlanFilter) (*models.Plan, error)
	List(ctx context.Context, filter *models.PlanFilter) ([]*models.Plan, error)
	Count(ctx context.Context, filter *models.PlanFilter) (int64, error)
	Update(ctx context.Context, planID int64, req *UpdatePlanInput) (*models.Plan, error)
	Activate(ctx context.Context, planID int64) (*models.Plan, error)
	Deactivate(ctx context.Context, planID int64) (*models.Plan, error)
}
