package contracts

import (
	"context"
	"github.com/google/uuid"
)

type OutboxRepository interface {
	ClaimNotificationOutbox(ctx context.Context, leaseToken uuid.UUID) (NotificationOutbox, error)
	CompleteNotificationOutbox(ctx context.Context, arg CompleteNotificationOutboxInput) error
	GetMaintenanceReminderDelivery(ctx context.Context, arg GetMaintenanceReminderDeliveryInput) (GetMaintenanceReminderDeliveryRecord, error)
	MarkOutboxInbox(ctx context.Context, arg MarkOutboxInboxInput) (int64, error)
	NotificationBacklog(ctx context.Context) (NotificationBacklogRecord, error)
	NotificationOutboxAccess(ctx context.Context, arg NotificationOutboxAccessInput) (bool, error)
	OutboxInboxID(ctx context.Context, arg OutboxInboxIDInput) (uuid.UUID, error)
	RetryNotificationOutbox(ctx context.Context, arg RetryNotificationOutboxInput) error
}
