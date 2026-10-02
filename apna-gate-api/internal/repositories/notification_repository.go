package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type NotificationRepository = contracts.NotificationRepository

type notificationRepository struct {
	db *database.Database
}

func NewNotificationRepository(db *database.Database) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification models.NotificationCreate) (*models.Notification, error) {
	data, err := jsonMap(notification.Data)
	if err != nil {
		return nil, err
	}

	row, err := GetQueries(ctx, r.db).CreateNotification(ctx, db.CreateNotificationParams{
		ID:        uuidParam(notification.ID),
		UserID:    notification.UserID,
		SocietyID: notification.SocietyID,
		FlatID:    notification.FlatID,
		Type:      notification.Type,
		Title:     notification.Title,
		Body:      notification.Body,
		Data:      data,
		EventKey:  notification.EventKey,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return notificationFromDB(row), nil
}

func (r *notificationRepository) ListByUser(ctx context.Context, filter models.NotificationListFilter) ([]*models.Notification, error) {
	rows, err := GetQueries(ctx, r.db).ListNotificationsByUser(ctx, db.ListNotificationsByUserParams{
		UserID:          filter.UserID,
		CursorCreatedAt: timePtrToPgTimestamptz(filter.CursorCreatedAt),
		CursorID:        uuidParamPtr(filter.CursorID),
		Limit:           normalizeLimit(filter.Limit),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*models.Notification, 0, len(rows))
	for _, row := range rows {
		items = append(items, notificationFromDB(row))
	}
	return items, nil
}

func (r *notificationRepository) CountUnreadByUser(ctx context.Context, userID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountUnreadNotificationsByUser(ctx, userID)
}

func (r *notificationRepository) MarkRead(ctx context.Context, userID int64, id string) (*models.Notification, error) {
	row, err := GetQueries(ctx, r.db).MarkNotificationRead(ctx, db.MarkNotificationReadParams{
		ID:     uuidParam(id),
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return notificationFromDB(row), nil
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, userID int64) error {
	return GetQueries(ctx, r.db).MarkAllNotificationsRead(ctx, userID)
}

func notificationFromDB(row db.Notification) *models.Notification {
	return &models.Notification{
		ID:        uuidString(row.ID),
		UserID:    row.UserID,
		SocietyID: row.SocietyID,
		FlatID:    row.FlatID,
		Type:      row.Type,
		Domain:    row.Domain,
		Title:     row.Title,
		Body:      row.Body,
		Data:      metadataFromJSON(row.Data),
		EventKey:  row.EventKey,
		ReadAt:    pgTimestamptzToTimePtr(row.ReadAt),
		CreatedAt: pgTimestamptzToTime(row.CreatedAt),
	}
}

func uuidParam(value string) pgtype.UUID {
	var id pgtype.UUID
	_ = id.Scan(value)
	return id
}

func uuidParamPtr(value *string) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{}
	}
	return uuidParam(*value)
}

func uuidString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	id, err := uuid.FromBytes(value.Bytes[:])
	if err != nil {
		return ""
	}
	return id.String()
}

func (r *notificationRepository) GetOwned(ctx context.Context, userID int64, id string) (*models.Notification, error) {
	row, err := GetQueries(ctx, r.db).GetOwnedNotification(ctx, db.GetOwnedNotificationParams{UserID: userID, ID: uuidParam(id)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return notificationFromDB(row), nil
}
