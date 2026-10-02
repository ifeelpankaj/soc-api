package maintenancesvc

import (
	"context"
	"go-server/internal/models"
)

type Documents interface {
	Invoice(ctx context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error)
	Receipt(ctx context.Context, society, user, id int64, resident bool) (models.MaintenancePDF, error)
}

var _ Documents = (*DocumentService)(nil)
