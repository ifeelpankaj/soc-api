package notificationsvc

import (
	"context"

	"go-server/internal/models"
)

type NotificationService interface {
	RegisterDeviceToken(ctx context.Context, userID int64, req models.RegisterDeviceTokenRequest) (*models.DeviceToken, error)
	UnregisterDeviceToken(ctx context.Context, userID int64, token string) error
	ListNotifications(ctx context.Context, userID int64, limit int32, cursor string) (*models.NotificationListResult, error)
	GetUnreadCount(ctx context.Context, userID int64) (*models.NotificationUnreadCountResponse, error)
	MarkNotificationRead(ctx context.Context, userID int64, notificationID string) (*models.NotificationReadResponse, error)
	MarkAllNotificationsRead(ctx context.Context, userID int64) (*models.NotificationUnreadCountResponse, error)
	SendVisitorApprovalRequested(ctx context.Context, entry *models.VisitorEntry) error
	SendVisitorApproved(ctx context.Context, entry *models.VisitorEntry) error
	SendVisitorRejected(ctx context.Context, entry *models.VisitorEntry) error
	SendVisitorCheckIn(ctx context.Context, entry *models.VisitorEntry) error
	SendVisitorCheckOut(ctx context.Context, entry *models.VisitorEntry) error
	SendVisitorInviteAccepted(ctx context.Context, entry *models.VisitorEntry) error
	SendMemberInviteAccepted(ctx context.Context, invite *models.FlatMemberInvite, flatNumber string, joinedName string, residentID int64) error
	SendToUser(ctx context.Context, userID int64, payload models.NotificationPayload) error
	SendToUsers(ctx context.Context, userIDs []int64, payload models.NotificationPayload) error
	Close() error
}
