package contracts

import (
	"context"
	"go-server/internal/models"
)

type FlatRepository interface {
	Create(ctx context.Context, flat *models.Flat) error
	Get(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error)
	List(ctx context.Context, filter *models.FlatFilter) ([]*models.Flat, error)
	Count(ctx context.Context, filter *models.FlatFilter) (int64, error)
	Stats(ctx context.Context, societyID int64) (*models.FlatStatsResponse, error)
	Update(ctx context.Context, filter *models.FlatFilter, req *UpdateFlatInput) (*models.Flat, error)
	Deactivate(ctx context.Context, filter *models.FlatFilter) error
	Block(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error)
	Unblock(ctx context.Context, filter *models.FlatFilter) (*models.Flat, error)
	MarkOccupied(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error)
	MarkVacant(ctx context.Context, societyID int64, flatID int64) (*models.Flat, error)
}
