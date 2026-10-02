package notificationsvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrInvalidDeviceToken   = models.NewAppError("INVALID_DEVICE_TOKEN", "invalid device token request", http.StatusBadRequest, nil)
	ErrDeviceTokenNotFound  = models.NewAppError("DEVICE_TOKEN_NOT_FOUND", "device token not found", http.StatusNotFound, nil)
	ErrNotificationDisabled = models.NewAppError("NOTIFICATION_DISABLED", "push notifications are disabled", http.StatusServiceUnavailable, nil)
	ErrNotificationSend     = models.NewAppError("NOTIFICATION_SEND_FAILED", "failed to send push notification", http.StatusInternalServerError, nil)
	ErrNotificationNotFound = models.NewAppError("NOTIFICATION_NOT_FOUND", "notification not found", http.StatusNotFound, nil)
	ErrInvalidNotification  = models.NewAppError("INVALID_NOTIFICATION_REQUEST", "invalid notification request", http.StatusBadRequest, nil)
)
