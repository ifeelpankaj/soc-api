package contracts

import (
	"context"
	"go-server/internal/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification models.NotificationCreate) (*models.Notification, error)
	ListByUser(ctx context.Context, filter models.NotificationListFilter) ([]*models.Notification, error)
	CountUnreadByUser(ctx context.Context, userID int64) (int64, error)
	MarkRead(ctx context.Context, userID int64, id string) (*models.Notification, error)
	MarkAllRead(ctx context.Context, userID int64) error
}
