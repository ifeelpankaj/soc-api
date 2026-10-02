package contracts

import (
	"context"
	"go-server/internal/models"
)

type FlatResidentRepository interface {
	Add(ctx context.Context, resident *models.FlatResident) error
	Get(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error)
	List(ctx context.Context, filter *models.FlatResidentFilter) ([]*models.FlatResident, error)
	Remove(ctx context.Context, filter *models.FlatResidentFilter) error
	MoveOut(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error)
	ClearPrimary(ctx context.Context, societyID int64, flatID int64) error
	SetPrimary(ctx context.Context, societyID int64, flatID int64, residentID int64) (*models.FlatResident, error)
	UpdateRole(ctx context.Context, filter *models.FlatResidentFilter, role models.FlatResidentRole) (*models.FlatResident, error)
	CountActive(ctx context.Context, societyID int64, flatID int64) (int64, error)
	CountPrimary(ctx context.Context, societyID int64, flatID int64) (int64, error)
}
