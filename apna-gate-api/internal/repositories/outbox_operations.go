package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *notificationRepository) ClaimNotificationOutbox(ctx context.Context, leaseToken uuid.UUID) (contracts.NotificationOutbox, error) {
	value, err := GetQueries(ctx, r.db).ClaimNotificationOutbox(ctx, pgtype.UUID{Bytes: leaseToken, Valid: leaseToken != uuid.Nil})
	return mapFinancialNotificationOutbox(value), persistenceError(err)
}

func (r *notificationRepository) CompleteNotificationOutbox(ctx context.Context, arg contracts.CompleteNotificationOutboxInput) error {
	return persistenceError(GetQueries(ctx, r.db).CompleteNotificationOutbox(ctx, db.CompleteNotificationOutboxParams{Delivered: arg.Delivered, ID: arg.ID, LeaseToken: pgtype.UUID{Bytes: arg.LeaseToken, Valid: arg.LeaseToken != uuid.Nil}}))
}

func (r *notificationRepository) GetMaintenanceReminderDelivery(ctx context.Context, arg contracts.GetMaintenanceReminderDeliveryInput) (contracts.GetMaintenanceReminderDeliveryRecord, error) {
	value, err := GetQueries(ctx, r.db).GetMaintenanceReminderDelivery(ctx, db.GetMaintenanceReminderDeliveryParams{UserID: arg.UserID, EventKey: arg.EventKey, SocietyID: arg.SocietyID})
	return mapFinancialGetMaintenanceReminderDeliveryRow(value), persistenceError(err)
}

func (r *notificationRepository) MarkOutboxInbox(ctx context.Context, arg contracts.MarkOutboxInboxInput) (int64, error) {
	value, err := GetQueries(ctx, r.db).MarkOutboxInbox(ctx, db.MarkOutboxInboxParams{ID: arg.ID, LeaseToken: pgtype.UUID{Bytes: arg.LeaseToken, Valid: arg.LeaseToken != uuid.Nil}})
	return value, persistenceError(err)
}

func (r *notificationRepository) NotificationBacklog(ctx context.Context) (contracts.NotificationBacklogRecord, error) {
	value, err := GetQueries(ctx, r.db).NotificationBacklog(ctx)
	return mapFinancialNotificationBacklogRow(value), persistenceError(err)
}

func (r *notificationRepository) NotificationOutboxAccess(ctx context.Context, arg contracts.NotificationOutboxAccessInput) (bool, error) {
	value, err := GetQueries(ctx, r.db).NotificationOutboxAccess(ctx, db.NotificationOutboxAccessParams{SocietyID: arg.SocietyID, UserID: arg.UserID, Audience: arg.Audience, FlatID: arg.FlatID})
	return value, persistenceError(err)
}

func (r *notificationRepository) OutboxInboxID(ctx context.Context, arg contracts.OutboxInboxIDInput) (uuid.UUID, error) {
	value, err := GetQueries(ctx, r.db).OutboxInboxID(ctx, db.OutboxInboxIDParams{UserID: arg.UserID, EventKey: arg.EventKey})
	return uuid.UUID(value.Bytes), persistenceError(err)
}

func (r *notificationRepository) RetryNotificationOutbox(ctx context.Context, arg contracts.RetryNotificationOutboxInput) error {
	return persistenceError(GetQueries(ctx, r.db).RetryNotificationOutbox(ctx, db.RetryNotificationOutboxParams{ID: arg.ID, LeaseToken: pgtype.UUID{Bytes: arg.LeaseToken, Valid: arg.LeaseToken != uuid.Nil}}))
}

var _ contracts.OutboxRepository = (*notificationRepository)(nil)
