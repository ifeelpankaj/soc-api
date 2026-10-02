package contracts

import (
	"context"
	"go-server/internal/models"
	"time"
)

type ShortLinkRepository interface {
	Create(ctx context.Context, code string, resourceType models.ShortLinkResourceType, resourceID int64, expiresAt *time.Time, createdBy *int64, metadata map[string]any) (*models.ShortLink, error)
	GetByCode(ctx context.Context, code string) (*models.ShortLink, error)
	GetByResource(ctx context.Context, resourceType models.ShortLinkResourceType, resourceID int64) (*models.ShortLink, error)
	Revoke(ctx context.Context, code string) (*models.ShortLink, error)
}
