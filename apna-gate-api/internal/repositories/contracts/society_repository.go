package contracts

import (
	"context"
	"go-server/internal/models"
)

type SocietyRepository interface {
	Create(ctx context.Context, society *models.Society) error
	Get(ctx context.Context, filter models.GetSocietyFilter) (*models.Society, error)
	List(ctx context.Context, filter models.ListSocietiesFilter) ([]*models.Society, error)
	Count(ctx context.Context, filter models.ListSocietiesFilter) (int64, error)
	Update(ctx context.Context, societyID int64, req UpdateSocietyInput) (*models.Society, error)
	Approve(ctx context.Context, societyID int64, approvedBy int64) (*models.Society, error)
	Reject(ctx context.Context, societyID int64, rejectedBy int64, reason string) (*models.Society, error)
	Suspend(ctx context.Context, societyID int64, suspendedBy int64, reason string) (*models.Society, error)
	Reactivate(ctx context.Context, societyID int64, reactivatedBy int64) (*models.Society, error)
	Restore(ctx context.Context, societyID int64) (*models.Society, error)
	SoftDelete(ctx context.Context, societyID int64) error
	CountPendingByCreator(ctx context.Context, createdBy int64) (int64, error)
}
