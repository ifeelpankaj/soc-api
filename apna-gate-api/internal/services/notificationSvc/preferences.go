package notificationsvc

import (
	"context"
	"errors"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

func (s *notificationService) GetPreferences(ctx context.Context, user, society int64) (contracts.NotificationPreferences, error) {
	store, ok := s.notifications.(contracts.PipelineRepository)
	if !ok {
		return contracts.NotificationPreferences{}, ErrNotificationDisabled
	}
	p, err := store.GetNotificationPreferences(ctx, user, society)
	if errors.Is(err, contracts.ErrNotFound) {
		return p, models.NewAppError("NOTIFICATION_PREFERENCES_FORBIDDEN", "Active society membership required", 403, nil)
	}
	return p, err
}

func (s *notificationService) SetPreferences(ctx context.Context, user, society int64, p contracts.NotificationPreferences) error {
	store, ok := s.notifications.(contracts.PipelineRepository)
	if !ok {
		return ErrNotificationDisabled
	}
	err := store.SetNotificationPreferences(ctx, user, society, p)
	if errors.Is(err, contracts.ErrNotFound) {
		return models.NewAppError("NOTIFICATION_PREFERENCES_FORBIDDEN", "Active society membership required", 403, nil)
	}
	return err
}
