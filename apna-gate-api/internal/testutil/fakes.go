package testutil

import (
	"context"
	"sync"

	"go-server/internal/models"
)

type TxManager struct {
	Calls int
	Err   error
}

func (m *TxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	m.Calls++
	if m.Err != nil {
		return m.Err
	}
	return fn(ctx)
}

type NotificationRecorder struct {
	mu           sync.Mutex
	UserPayloads map[int64][]models.NotificationPayload
	BatchPayload []NotificationBatch
	Err          error
	Closed       bool
}

type NotificationBatch struct {
	UserIDs []int64
	Payload models.NotificationPayload
}

func NewNotificationRecorder() *NotificationRecorder {
	return &NotificationRecorder{UserPayloads: map[int64][]models.NotificationPayload{}}
}

func (n *NotificationRecorder) RegisterDeviceToken(context.Context, int64, models.RegisterDeviceTokenRequest) (*models.DeviceToken, error) {
	return &models.DeviceToken{}, n.Err
}

func (n *NotificationRecorder) UnregisterDeviceToken(context.Context, int64, string) error {
	return n.Err
}

func (n *NotificationRecorder) ListNotifications(context.Context, int64, int32, string) (*models.NotificationListResult, error) {
	return &models.NotificationListResult{}, n.Err
}

func (n *NotificationRecorder) GetUnreadCount(context.Context, int64) (*models.NotificationUnreadCountResponse, error) {
	return &models.NotificationUnreadCountResponse{}, n.Err
}

func (n *NotificationRecorder) MarkNotificationRead(context.Context, int64, string) (*models.NotificationReadResponse, error) {
	return &models.NotificationReadResponse{}, n.Err
}

func (n *NotificationRecorder) MarkAllNotificationsRead(context.Context, int64) (*models.NotificationUnreadCountResponse, error) {
	return &models.NotificationUnreadCountResponse{}, n.Err
}

func (n *NotificationRecorder) SendVisitorApprovalRequested(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Visitor approval required", Data: map[string]string{"type": "visitor.pending"}})
}

func (n *NotificationRecorder) SendVisitorApproved(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Visitor approved", Data: map[string]string{"type": "visitor.approved"}})
}

func (n *NotificationRecorder) SendVisitorRejected(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Visitor rejected", Data: map[string]string{"type": "visitor.rejected"}})
}

func (n *NotificationRecorder) SendVisitorCheckIn(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Visitor checked in", Data: map[string]string{"type": "visitor.checkin"}})
}

func (n *NotificationRecorder) SendVisitorCheckOut(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Visitor checked out", Data: map[string]string{"type": "visitor.checkout"}})
}

func (n *NotificationRecorder) SendVisitorInviteAccepted(_ context.Context, entry *models.VisitorEntry) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Guest invite accepted", Data: map[string]string{"type": "visitor_invite.accepted"}})
}

func (n *NotificationRecorder) SendMemberInviteAccepted(_ context.Context, invite *models.FlatMemberInvite, flatNumber string, joinedName string, residentID int64) error {
	return n.SendToUsers(context.Background(), []int64{}, models.NotificationPayload{Title: "Member joined your flat", Data: map[string]string{"type": "member_invite.accepted"}})
}

func (n *NotificationRecorder) SendToUser(_ context.Context, userID int64, payload models.NotificationPayload) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.UserPayloads[userID] = append(n.UserPayloads[userID], payload)
	return n.Err
}

func (n *NotificationRecorder) SendToUsers(_ context.Context, userIDs []int64, payload models.NotificationPayload) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	copied := append([]int64(nil), userIDs...)
	n.BatchPayload = append(n.BatchPayload, NotificationBatch{UserIDs: copied, Payload: payload})
	return n.Err
}

func (n *NotificationRecorder) Close() error {
	n.Closed = true
	return n.Err
}
