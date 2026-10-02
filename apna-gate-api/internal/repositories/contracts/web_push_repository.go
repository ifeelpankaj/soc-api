package contracts

import (
	"context"
	"go-server/internal/models"
)

type WebDeviceTokenRepository interface {
	UpsertWeb(context.Context, int64, string, string, int64) (*models.DeviceToken, error)
	ListEligibleWeb(context.Context, int64, int64, int64) ([]*models.DeviceToken, error)
}
