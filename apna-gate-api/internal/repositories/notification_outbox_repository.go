package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/db"
	"go-server/internal/models"
)

func (r *notificationRepository) Enqueue(ctx context.Context, n models.NotificationCreate, audience string) error {
	if n.SocietyID == nil || n.EventKey == nil {
		return errors.New("outbox requires society and event key")
	}
	payload, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return GetQueries(ctx, r.db).EnqueueNotificationOutbox(ctx, db.EnqueueNotificationOutboxParams{UserID: n.UserID, SocietyID: *n.SocietyID, FlatID: n.FlatID, Audience: audience, EventKey: *n.EventKey, Payload: payload})
}
