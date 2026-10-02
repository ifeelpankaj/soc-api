package contracts

import (
	"context"
	"go-server/internal/models"
)

type FlatClaimRepository interface {
	Submit(ctx context.Context, claim *models.FlatClaim) error
	Get(ctx context.Context, filter *models.FlatClaimFilter) (*models.FlatClaim, error)
	List(ctx context.Context, filter *models.FlatClaimFilter) ([]*models.FlatClaim, error)
	Stats(ctx context.Context, societyID int64) (*models.FlatClaimStatsResponse, error)
	Approve(ctx context.Context, societyID int64, claimID int64, reviewedBy int64) (*models.FlatClaim, error)
	Reject(ctx context.Context, societyID int64, claimID int64, reviewedBy int64, reason string) (*models.FlatClaim, error)
	Cancel(ctx context.Context, claimID int64, userID int64) (*models.FlatClaim, error)
}
