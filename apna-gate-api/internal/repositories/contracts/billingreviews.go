package contracts

import (
	"context"
)

type BillingReviewsRepository interface {
	GetMaintenanceReview(ctx context.Context, arg GetMaintenanceReviewInput) (GetMaintenanceReviewRecord, error)
	SaveMaintenanceReview(ctx context.Context, arg SaveMaintenanceReviewInput) error
}
