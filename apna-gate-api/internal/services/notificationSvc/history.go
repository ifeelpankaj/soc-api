package notificationsvc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

const (
	defaultNotificationLimit = 20
	maxNotificationLimit     = 50
)

type notificationCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func (s *notificationService) ListNotifications(ctx context.Context, userID int64, limit int32, cursor string) (*models.NotificationListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if userID <= 0 || s.notifications == nil {
		return nil, ErrInvalidNotification
	}
	limit = normalizeNotificationLimit(limit)
	cursorValue, err := decodeNotificationCursor(cursor)
	if err != nil {
		return nil, ErrInvalidNotification.WithCause(err)
	}

	queryLimit := limit + 1
	filter := models.NotificationListFilter{UserID: userID, Limit: queryLimit}
	if cursorValue != nil {
		filter.CursorCreatedAt = &cursorValue.CreatedAt
		filter.CursorID = &cursorValue.ID
	}
	items, err := s.notifications.ListByUser(ctx, filter)
	if err != nil {
		return nil, err
	}

	var nextCursor *string
	if len(items) > int(limit) {
		last := items[limit-1]
		items = items[:limit]
		encoded := encodeNotificationCursor(last)
		nextCursor = &encoded
	}

	return &models.NotificationListResult{Items: items, NextCursor: nextCursor}, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (*models.NotificationUnreadCountResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if userID <= 0 || s.notifications == nil {
		return nil, ErrInvalidNotification
	}
	count, err := s.notifications.CountUnreadByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &models.NotificationUnreadCountResponse{UnreadCount: count}, nil
}

func (s *notificationService) MarkNotificationRead(ctx context.Context, userID int64, notificationID string) (*models.NotificationReadResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if userID <= 0 || strings.TrimSpace(notificationID) == "" || s.notifications == nil {
		return nil, ErrInvalidNotification
	}
	item, err := s.notifications.MarkRead(ctx, userID, strings.TrimSpace(notificationID))
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNotificationNotFound
	}
	return &models.NotificationReadResponse{ID: item.ID, ReadAt: item.ReadAt}, nil
}

func (s *notificationService) MarkAllNotificationsRead(ctx context.Context, userID int64) (*models.NotificationUnreadCountResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if userID <= 0 || s.notifications == nil {
		return nil, ErrInvalidNotification
	}
	if err := s.notifications.MarkAllRead(ctx, userID); err != nil {
		return nil, err
	}
	return &models.NotificationUnreadCountResponse{UnreadCount: 0}, nil
}

func normalizeNotificationLimit(limit int32) int32 {
	if limit <= 0 {
		return defaultNotificationLimit
	}
	if limit > maxNotificationLimit {
		return maxNotificationLimit
	}
	return limit
}

func decodeNotificationCursor(raw string) (*notificationCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	bytes, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	var cursor notificationCursor
	if err := json.Unmarshal(bytes, &cursor); err != nil {
		return nil, err
	}
	if cursor.CreatedAt.IsZero() || strings.TrimSpace(cursor.ID) == "" {
		return nil, ErrInvalidNotification
	}
	return &cursor, nil
}

func encodeNotificationCursor(item *models.Notification) string {
	if item == nil {
		return ""
	}
	bytes, _ := json.Marshal(notificationCursor{CreatedAt: item.CreatedAt, ID: item.ID})
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (s *notificationService) GetNotification(ctx context.Context, userID int64, id string) (*models.Notification, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidNotification
	}
	repo, ok := s.notifications.(interface {
		GetOwned(context.Context, int64, string) (*models.Notification, error)
	})
	if !ok {
		return nil, ErrNotificationNotFound
	}
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	item, err := repo.GetOwned(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNotificationNotFound
	}
	return item, nil
}
