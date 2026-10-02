package contracts

import (
	"context"
	"go-server/internal/models"
)

type DeviceTokenRepository interface {
	Upsert(ctx context.Context, userID int64, token string, platform models.DevicePlatform, deviceID *string) (*models.DeviceToken, error)
	Delete(ctx context.Context, userID int64, token string) error
	ListByUserID(ctx context.Context, userID int64) ([]*models.DeviceToken, error)
	DeleteByToken(ctx context.Context, token string) error
}
