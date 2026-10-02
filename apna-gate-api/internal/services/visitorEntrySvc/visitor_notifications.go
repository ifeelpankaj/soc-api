package visitorentrysvc

import (
	"context"
	"time"

	"go-server/internal/models"
	"go-server/pkg/logger"

	"go.uber.org/zap"
)

const notificationDispatchTimeout = 10 * time.Second

func (s *VisitorEntrySvc) notifyVisitorPending(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorApprovalRequested(ctx, entry)
	})
}

func (s *VisitorEntrySvc) notifyVisitorApproved(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorApproved(ctx, entry)
	})
}

func (s *VisitorEntrySvc) notifyVisitorRejected(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorRejected(ctx, entry)
	})
}

func (s *VisitorEntrySvc) notifyVisitorCheckIn(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorCheckIn(ctx, entry)
	})
}

func (s *VisitorEntrySvc) notifyVisitorCheckOut(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorCheckOut(ctx, entry)
	})
}

func (s *VisitorEntrySvc) notifyVisitorInviteAccepted(entry *models.VisitorEntry) {
	s.dispatchNotification(func(ctx context.Context) error {
		return s.notifier.SendVisitorInviteAccepted(ctx, entry)
	})
}

func (s *VisitorEntrySvc) dispatchNotification(send func(context.Context) error) {
	// Durable notifications were enqueued with the visitor event transaction.
	if durable, ok := s.notifier.(interface{ DurableNotifications() bool }); ok && durable.DurableNotifications() {
		return
	}
	if s.notifier == nil || send == nil {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), notificationDispatchTimeout)
		defer cancel()

		if err := send(ctx); err != nil {
			logger.Warn("failed to dispatch visitor notification", zap.Error(err))
		}
	}()
}
