package maintenancesvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"

	"strconv"
	"time"
)

var reminderEvents = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "apna_gate_maintenance_reminder_events_total", Help: "Maintenance reminder scheduling outcomes."}, []string{"outcome"})

func init() { prometheus.MustRegister(reminderEvents) }

func reminderMilestone(now time.Time, due time.Time, zone string) (string, string, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return "", "", err
	}
	today := now.In(loc).Format("2006-01-02")
	for _, stage := range []struct {
		offset int
		name   string
	}{{-3, "before_3"}, {0, "due"}, {7, "after_7"}} {
		date := due.AddDate(0, 0, stage.offset).Format("2006-01-02")
		if today == date {
			return stage.name, date, nil
		}
	}
	return "", "", nil
}

// The legacy delivery path has the same eligibility checks as the shared outbox.
func (s *Service) prepareReminderDelivery(ctx context.Context, d *contracts.MaintenanceDelivery) (bool, error) {
	if err := s.operational.EnsureSocietyOperational(ctx, d.Bill.SocietyID); err != nil {
		var app *models.AppError
		if errors.As(err, &app) && app.StatusCode >= 400 && app.StatusCode < 500 {
			return false, nil
		}
		return false, err
	}
	v, err := s.store.Settings(ctx, d.Bill.SocietyID)
	if err != nil || !v.Enabled {
		return false, err
	}
	b, err := s.Get(ctx, models.MaintenanceBillFilter{SocietyID: d.Bill.SocietyID, UserID: d.UserID, ID: d.Bill.ID, Resident: true})
	if err != nil {
		var app *models.AppError
		if errors.As(err, &app) && app.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	if b.OutstandingAmountPaise <= 0 {
		return false, nil
	}
	d.Bill = &b
	if d.EventData == nil {
		d.EventData = map[string]string{}
	}
	d.EventData["billing_month"], d.EventData["due_date"] = b.BillingMonth, b.DueDate
	d.EventData["outstanding_amount_paise"] = strconv.FormatInt(b.OutstandingAmountPaise, 10)
	return true, nil
}

// RunReminders only creates reminder events. Issuance remains a separate job.
func (s *Service) RunReminders(ctx context.Context) error {
	store, ok := s.store.(billingPersistence)
	if !ok {
		return errors.New("billing repository does not support reminders")
	}
	ids, err := s.store.Societies(ctx)
	if err != nil {
		return err
	}
	var failures []error
	for _, society := range ids {
		if err := ctx.Err(); err != nil {
			return errors.Join(append(failures, err)...)
		}
		var after int64
		for {
			count := 0
			var queued, skipped int64
			err := s.store.Locked(ctx, society, func(ctx context.Context) error {
				if err := s.operational.EnsureSocietyOperational(ctx, society); err != nil {
					var app *models.AppError
					if errors.As(err, &app) && app.StatusCode >= 400 && app.StatusCode < 500 {
						skipped++
						return nil
					}
					return err
				}
				repo := store
				now := s.now()
				bills, err := repo.ListMaintenanceReminderCandidates(ctx, contracts.ListMaintenanceReminderCandidatesInput{SocietyID: society, AfterID: after, AsOf: now})
				if err != nil {
					return err
				}
				count = len(bills)
				for _, bill := range bills {
					after = bill.ID
					stage, date, err := reminderMilestone(now, bill.DueDate, bill.Timezone)
					if err != nil {
						return err
					}
					if stage == "" {
						continue
					}
					key := fmt.Sprintf("%s:%d:%s:%s", models.MaintenanceReminderType, bill.ID, stage, date)
					data, err := json.Marshal(map[string]string{"reminder_milestone": stage, "scheduled_date": date})
					if err != nil {
						return err
					}
					n, err := repo.EnqueueMaintenanceReminder(ctx, contracts.EnqueueMaintenanceReminderInput{BillID: bill.ID, SocietyID: society, EventKey: key, EventData: data})
					if err != nil {
						return err
					}
					queued += n
					if n == 0 {
						skipped++
					}
				}
				return nil
			})
			if err != nil {
				reminderEvents.WithLabelValues("failed").Inc()
				failures = append(failures, fmt.Errorf("society %d reminders: %w", society, err))
				break
			}
			reminderEvents.WithLabelValues("queued").Add(float64(queued))
			reminderEvents.WithLabelValues("skipped").Add(float64(skipped))
			if count < 500 {
				break
			}
		}
	}
	return errors.Join(failures...)
}
