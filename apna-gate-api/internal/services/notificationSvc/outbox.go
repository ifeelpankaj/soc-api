package notificationsvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/prometheus/client_golang/prometheus"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	"go-server/internal/requestctx"
	"strconv"
	"time"
)

type outboxStore interface {
	Enqueue(context.Context, models.NotificationCreate, string) error
	contracts.OutboxRepository
}

var outboxPending = prometheus.NewGauge(prometheus.GaugeOpts{Name: "apna_gate_notification_outbox_pending", Help: "Pending notification deliveries."})
var outboxFailures = prometheus.NewCounter(prometheus.CounterOpts{Name: "apna_gate_notification_outbox_failures_total", Help: "Failed notification delivery attempts."})
var reminderDeliveries = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "apna_gate_maintenance_reminder_deliveries_total", Help: "Reminder delivery outcomes."}, []string{"outcome"})

func init() { prometheus.MustRegister(outboxPending, outboxFailures, reminderDeliveries) }

func (s *notificationService) DurableNotifications() bool {
	_, ok := s.notifications.(outboxStore)
	return ok
}

func (s *notificationService) enqueueSpec(ctx context.Context, users []int64, spec notificationSpec) error {
	store := s.notifications.(outboxStore)
	key := spec.Type + ":" + fmt.Sprint(spec.Data["entry_id"])
	if spec.EventKey != nil {
		key = *spec.EventKey
	}
	audience := "resident"
	if spec.Type == EventVisitorApproved || spec.Type == EventVisitorRejected {
		audience = "staff"
	} else if spec.FlatID == nil {
		audience = "society_resident"
	}
	for _, user := range users {
		data := cloneData(spec.Data)
		data["type"] = spec.Type
		data["event"] = spec.Alias
		data["event_key"] = key
		data["category_id"] = spec.CategoryID
		n := models.NotificationCreate{ID: uuid.NewString(), UserID: user, SocietyID: spec.SocietyID, FlatID: spec.FlatID, Type: spec.Type, Title: spec.Title, Body: spec.Body, Data: data, EventKey: &key}
		if err := store.Enqueue(ctx, n, audience); err != nil {
			return err
		}
	}
	return nil
}

func (s *notificationService) DeliverOutbox(ctx context.Context) error {
	_, err := s.deliverOutboxBatch(ctx)
	return err
}

// DrainOutbox processes multiple batches, bounded by the calling job's deadline.
func (s *notificationService) DrainOutbox(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	var failures []error
	for {
		processed, err := s.deliverOutboxBatch(ctx)
		if err != nil {
			failures = append(failures, err)
		}
		if processed < 100 || ctx.Err() != nil {
			return errors.Join(failures...)
		}
	}
}

func (s *notificationService) deliverOutboxBatch(ctx context.Context) (int, error) {
	store, ok := s.notifications.(outboxStore)
	if !ok {
		return 0, nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	repo := store
	defer func() {
		if backlog, err := repo.NotificationBacklog(ctx); err == nil {
			outboxPending.Set(float64(backlog.Pending))
		}
	}()
	var failures []error
	processed := 0
	for i := 0; i < 100; i++ {
		if err := ctx.Err(); err != nil {
			return processed, errors.Join(append(failures, err)...)
		}
		d, err := repo.ClaimNotificationOutbox(ctx, uuid.New())
		if errors.Is(err, contracts.ErrNotFound) {
			break
		}
		if err != nil {
			return processed, err
		}
		processed++
		err = s.deliverOutboxItem(ctx, repo, d)
		if err != nil {
			outboxFailures.Inc()
			failures = append(failures, errors.New("notification delivery failed"))
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			retryErr := repo.RetryNotificationOutbox(cleanup, contracts.RetryNotificationOutboxInput{ID: d.ID, LeaseToken: d.LeaseToken})
			cancel()
			if retryErr != nil {
				return processed, retryErr
			}
		}
	}
	return processed, errors.Join(failures...)
}
func (s *notificationService) deliverOutboxItem(ctx context.Context, repo contracts.OutboxRepository, d contracts.NotificationOutbox) (resultErr error) {
	var n models.NotificationCreate
	if err := json.Unmarshal(d.Payload, &n); err != nil {
		return err
	}
	defer func() {
		if resultErr != nil && n.Type == models.MaintenanceReminderType {
			reminderDeliveries.WithLabelValues("failed").Inc()
		}
	}()
	allowed, err := s.outboxDeliveryAllowed(ctx, repo, d, &n)
	if err != nil {
		return err
	}
	if !allowed {
		return s.skipOutboxDelivery(ctx, repo, d, n.Type)
	}
	n.ID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("outbox/%d/%s", d.UserID, d.EventKey))).String()
	n.UserID = d.UserID
	n.SocietyID = &d.SocietyID
	n.FlatID = d.FlatID
	n.EventKey = &d.EventKey
	if n.Data == nil {
		n.Data = map[string]any{}
	}
	n.Data["notification_id"] = n.ID
	n.Data["event_key"] = d.EventKey
	if pipeline, ok := s.notifications.(contracts.PipelineRepository); ok {
		// The repository commits inbox, per-device work and producer acknowledgement
		// atomically. Push delivery has an independent lifecycle.
		return pipeline.MaterializeOutbox(ctx, d, n, d.PushEnabled, s.enabled, s.enabled && s.webEnabled, s.apnsClient != nil)
	}
	if d.InboxCompletedAt.IsZero() {
		if _, err = s.notifications.Create(ctx, n); err != nil {
			return err
		}
	}
	changed, err := repo.MarkOutboxInbox(ctx, contracts.MarkOutboxInboxInput{ID: d.ID, LeaseToken: d.LeaseToken})
	if err != nil {
		return err
	}
	if changed != 1 {
		return errors.New("notification lease lost")
	}
	inboxID, err := repo.OutboxInboxID(ctx, contracts.OutboxInboxIDInput{UserID: d.UserID, EventKey: &d.EventKey})
	if err != nil {
		return err
	}
	n.Data["notification_id"] = inboxID.String()
	// Access and payment may change while an inbox entry is being persisted.
	allowed, err = s.outboxDeliveryAllowed(ctx, repo, d, &n)
	if err != nil {
		return err
	}
	if !allowed {
		return s.skipOutboxDelivery(ctx, repo, d, n.Type)
	}
	if !d.PushEnabled {
		return repo.CompleteNotificationOutbox(ctx, contracts.CompleteNotificationOutboxInput{ID: d.ID, LeaseToken: d.LeaseToken})
	}
	// A lease is longer than the bounded push attempt. An external send is at-least-once.
	pushCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	category, _ := n.Data["category_id"].(string)
	if err = s.SendToUser(pushCtx, d.UserID, models.NotificationPayload{Title: n.Title, Body: n.Body, Data: stringMap(n.Data), CategoryID: category}); err != nil {
		return err
	}
	err = repo.CompleteNotificationOutbox(ctx, contracts.CompleteNotificationOutboxInput{ID: d.ID, LeaseToken: d.LeaseToken, Delivered: true})
	if err == nil && n.Type == models.MaintenanceReminderType {
		reminderDeliveries.WithLabelValues("completed").Inc()
	}
	return err
}

func (s *notificationService) skipOutboxDelivery(ctx context.Context, repo contracts.OutboxRepository, d contracts.NotificationOutbox, kind string) error {
	if _, err := repo.MarkOutboxInbox(ctx, contracts.MarkOutboxInboxInput{ID: d.ID, LeaseToken: d.LeaseToken}); err != nil {
		return err
	}
	err := repo.CompleteNotificationOutbox(ctx, contracts.CompleteNotificationOutboxInput{ID: d.ID, LeaseToken: d.LeaseToken})
	if err == nil && kind == models.MaintenanceReminderType {
		reminderDeliveries.WithLabelValues("skipped").Inc()
	}
	return err
}

func (s *notificationService) outboxDeliveryAllowed(ctx context.Context, repo contracts.OutboxRepository, d contracts.NotificationOutbox, n *models.NotificationCreate) (bool, error) {
	if d.Audience == "hub_member" && (n.Type == "hub.announcement" || n.Type == "hub.reply") {
		postID, parseErr := strconv.ParseInt(fmt.Sprint(n.Data["post_id"]), 10, 64)
		if parseErr != nil || postID <= 0 {
			return false, nil
		}
		if pipeline, ok := s.notifications.(contracts.PipelineRepository); ok {
			available, err := pipeline.HubContentAvailable(ctx, d.SocietyID, postID)
			if err != nil || !available {
				return false, err
			}
		}
	}
	allowed, err := repo.NotificationOutboxAccess(ctx, contracts.NotificationOutboxAccessInput{SocietyID: d.SocietyID, UserID: d.UserID, Audience: d.Audience, FlatID: d.FlatID})
	if err == nil && allowed && d.Audience == "hub_member" {
		if s.operational == nil {
			return false, errors.New("hub delivery requires operational guard")
		}
		if guardErr := s.operational.EnsureSocietyOperational(requestctx.WithoutDeveloperGuardBypass(ctx), d.SocietyID); guardErr != nil {
			var app *models.AppError
			if errors.As(guardErr, &app) && app.StatusCode >= 400 && app.StatusCode < 500 {
				return false, nil
			}
			return false, guardErr
		}
	}
	if err != nil || !allowed || n.Type != models.MaintenanceReminderType {
		return allowed, err
	}
	if s.operational == nil {
		return false, errors.New("maintenance reminder delivery requires an operational guard")
	}
	if err := s.operational.EnsureSocietyOperational(ctx, d.SocietyID); err != nil {
		var app *models.AppError
		if errors.As(err, &app) && app.StatusCode >= 400 && app.StatusCode < 500 {
			return false, nil
		}
		return false, err
	}
	bill, err := repo.GetMaintenanceReminderDelivery(ctx, contracts.GetMaintenanceReminderDeliveryInput{UserID: d.UserID, EventKey: d.EventKey, SocietyID: d.SocietyID})
	if errors.Is(err, contracts.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if n.Data == nil {
		n.Data = map[string]any{}
	}
	stage, _ := n.Data["reminder_milestone"].(string)
	month, due := bill.BillingMonth.Format("2006-01"), bill.DueDate.Format("2006-01-02")
	n.Title, n.Body = models.MaintenanceReminderText(month, bill.BillNumber, due, bill.OutstandingAmountPaise, stage)
	n.Data["billing_month"], n.Data["due_date"] = month, due
	n.Data["outstanding_amount_paise"] = strconv.FormatInt(bill.OutstandingAmountPaise, 10)
	return true, nil
}
