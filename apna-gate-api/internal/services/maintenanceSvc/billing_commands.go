package maintenancesvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
	service "go-server/internal/services"
	"strings"
	"time"

	"github.com/google/uuid"
)

type billingPersistence interface {
	contracts.IdempotencyRepository
	contracts.FinancialEventsRepository
	contracts.BillingRemindersRepository
	contracts.BillingReviewsRepository
	MaintenancePendingClaimsCount(context.Context) (int64, error)
}

func billingWrite[T any](s *Service, ctx context.Context, society, user int64, key, operation string, input any, fn func(context.Context, billingPersistence) (T, error)) (T, error) {
	var result T
	if key == "" || len(key) > 128 || strings.TrimSpace(key) != key {
		return result, invalid("Idempotency-Key is required (1–128 bytes, no surrounding whitespace)")
	}
	hash, err := digest(input)
	if err != nil {
		return result, err
	}
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	err = s.store.Locked(ctx, society, func(ctx context.Context) error {
		if err := s.admin(ctx, society, user); err != nil {
			return err
		}
		if err := s.operational.EnsureSocietyOperational(ctx, society); err != nil {
			return err
		}
		store, ok := s.store.(billingPersistence)
		if !ok {
			return errors.New("billing repository does not support financial commands")
		}
		repo := store
		prior, err := repo.GetUPIIdempotency(ctx, contracts.GetUPIIdempotencyInput{SocietyID: society, ActorID: user, Operation: operation, Key: key})
		if err == nil {
			if prior.RequestHash != hash {
				return paymentConflict("IDEMPOTENCY_CONFLICT", "Idempotency key was used with different request contents")
			}
			return json.Unmarshal(prior.Response, &result)
		}
		if !errors.Is(err, contracts.ErrNotFound) {
			return err
		}
		result, err = fn(ctx, repo)
		if err != nil {
			return err
		}
		response, err := json.Marshal(result)
		if err != nil {
			return err
		}
		return repo.InsertUPIIdempotency(ctx, contracts.InsertUPIIdempotencyInput{SocietyID: society, ActorID: user, Operation: operation, Key: key, RequestHash: hash, Response: response})
	})
	observeFinancial(operation, err)
	return result, paymentError(err)
}

func (s *Service) Configure(ctx context.Context, society, user int64, key string, v models.MaintenanceSettings) (models.MaintenanceSettings, error) {
	v.FirstEnabledMonth = ""
	v.EligibleStatuses = nil
	return billingWrite(s, ctx, society, user, key, "billing/settings", v, func(ctx context.Context, repo billingPersistence) (models.MaintenanceSettings, error) {
		result, err := s.SaveSettings(ctx, society, user, v)
		if err != nil {
			return result, err
		}
		if err = paymentAudit(ctx, repo, society, user, 0, "billing_settings_changed", fmt.Sprint(society), v); err != nil {
			return result, err
		}
		return s.Settings(ctx, society, user)
	})
}

func catchUpMonth(month string, now time.Time, v models.MaintenanceSettings) (time.Time, error) {
	loc, err := time.LoadLocation(v.Timezone)
	if err != nil {
		return time.Time{}, err
	}
	m, err := time.ParseInLocation("2006-01", month, loc)
	if err != nil || v.FirstEnabledMonth == "" || month < v.FirstEnabledMonth || month >= now.In(loc).Format("2006-01") {
		return time.Time{}, invalid("Catch-up month must be between the first enabled month and the previous month")
	}
	return m, nil
}

func (s *Service) catchUpPreview(ctx context.Context, id int64, req models.MaintenanceMonthRequest, v models.MaintenanceSettings) (models.MaintenancePreview, error) {
	if !req.CatchUp {
		return s.preview(ctx, id, req.BillingMonth, v)
	}
	if _, err := catchUpMonth(req.BillingMonth, s.now(), v); err != nil {
		return models.MaintenancePreview{}, err
	}
	// Pricing and eligibility are deliberately current; only the billing period changes.
	p, err := s.preview(ctx, id, req.BillingMonth, v, true)
	if err != nil {
		return p, err
	}
	loc, _ := time.LoadLocation(v.Timezone)
	now := s.now().In(loc)
	due := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, 10).Format("2006-01-02")
	p.IsCatchUp = true
	for i := range p.Bills {
		p.Bills[i].DueDate = due
		p.Bills[i].IsCatchUp = true
	}
	return p, nil
}

func reviewHash(v models.MaintenanceSettings, p models.MaintenancePreview) (string, error) {
	return digest(struct {
		Settings models.MaintenanceSettings
		Preview  models.MaintenancePreview
	}{v, p})
}

func (s *Service) PreviewCommand(ctx context.Context, id, user int64, req models.MaintenanceMonthRequest) (models.MaintenancePreview, error) {
	ctx, cancel := context.WithTimeout(ctx, FinancialTimeout)
	defer cancel()
	if !req.CatchUp {
		return s.Preview(ctx, id, user, req.BillingMonth)
	}
	var p models.MaintenancePreview
	err := s.store.Locked(ctx, id, func(ctx context.Context) error {
		if err := s.admin(ctx, id, user); err != nil {
			return err
		}
		if err := s.operational.EnsureSocietyOperational(ctx, id); err != nil {
			return err
		}
		v, err := s.store.Settings(ctx, id)
		if err != nil {
			return err
		}
		run, err := s.store.Run(ctx, id, req.BillingMonth)
		if err != nil {
			return err
		}
		if run != nil {
			return paymentConflict("MAINTENANCE_ALREADY_ISSUED", "This month is already issued; view its bills")
		}
		p, err = s.catchUpPreview(ctx, id, req, v)
		if err != nil {
			return err
		}
		hash, err := reviewHash(v, p)
		if err != nil {
			return err
		}
		token := uuid.New()
		expiry := s.now().Add(PreviewLifetime)
		m, _ := time.Parse("2006-01", req.BillingMonth)
		store, ok := s.store.(billingPersistence)
		if !ok {
			return errors.New("billing repository does not support preview reviews")
		}
		err = store.SaveMaintenanceReview(ctx, contracts.SaveMaintenanceReviewInput{Token: token, SocietyID: id, ActorID: user, BillingMonth: m, SnapshotHash: hash, ExpiresAt: expiry})
		p.ReviewToken = token.String()
		p.ExpiresAt = &expiry
		return err
	})
	return p, err
}

func (s *Service) GenerateCommand(ctx context.Context, id, user int64, key string, req models.MaintenanceMonthRequest) (models.MaintenanceRunResult, error) {
	return billingWrite(s, ctx, id, user, key, "billing/generate", req, func(ctx context.Context, repo billingPersistence) (models.MaintenanceRunResult, error) {
		var result models.MaintenanceRunResult
		if !req.CatchUp {
			var err error
			result, err = s.Generate(ctx, id, user, req.BillingMonth)
			if err != nil {
				return result, err
			}
		} else {
			if !req.CurrentPricingAcknowledged {
				return result, invalid("Acknowledge current pricing and flat information before issuing missed bills")
			}
			v, err := s.store.Settings(ctx, id)
			if err != nil {
				return result, err
			}
			p, err := s.catchUpPreview(ctx, id, req, v)
			if err != nil {
				return result, err
			}
			token, err := paymentUUID(req.ReviewToken)
			if err != nil {
				return result, previewChanged()
			}
			m, _ := time.Parse("2006-01", req.BillingMonth)
			review, err := repo.GetMaintenanceReview(ctx, contracts.GetMaintenanceReviewInput{Token: token, SocietyID: id, ActorID: user, BillingMonth: m})
			if errors.Is(err, contracts.ErrNotFound) {
				return result, previewChanged()
			}
			if err != nil {
				return result, err
			}
			hash, err := reviewHash(v, p)
			if err != nil {
				return result, err
			}
			if !review.ExpiresAt.After(s.now()) || hash != review.SnapshotHash {
				return result, previewChanged()
			}
			if len(p.Issues) > 0 {
				return result, &ValidationIssues{Issues: p.Issues}
			}
			run, err := s.store.Run(ctx, id, req.BillingMonth)
			if err != nil {
				return result, err
			}
			if run != nil {
				return result, paymentConflict("MAINTENANCE_ALREADY_ISSUED", "This month is already issued; view its bills")
			}
			v.EligibleStatuses = append([]string(nil), s.statuses...)
			result, err = s.store.Issue(ctx, id, user, req.BillingMonth, v, p.Bills)
			if err != nil {
				return result, err
			}
		}
		action := "billing_issued"
		if result.Existing > 0 || result.Created == 0 {
			action = "billing_reconciled"
		}
		return result, paymentAudit(ctx, repo, id, user, 0, action, fmt.Sprint(result.RunID), struct {
			Request models.MaintenanceMonthRequest `json:"request"`
			Result  models.MaintenanceRunResult    `json:"result"`
		}{req, result})
	})
}
func previewChanged() error {
	return paymentConflict("MAINTENANCE_PREVIEW_CHANGED", "Billing information changed or the preview expired. Review a fresh preview.")
}

func (s *Service) Outstanding(ctx context.Context, society, user, flat int64) (models.MaintenanceOutstanding, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	store, ok := s.store.(interface {
		Outstanding(context.Context, int64, int64, int64, time.Time) (models.MaintenanceOutstanding, error)
	})
	if !ok {
		return models.MaintenanceOutstanding{}, errors.New("billing repository does not support outstanding summaries")
	}
	return store.Outstanding(ctx, society, user, flat, s.now())
}

func (s *Service) OutstandingFlats(ctx context.Context, f models.MaintenanceOutstandingFlatFilter) (models.MaintenanceOutstandingFlatList, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	store, ok := s.store.(interface {
		OutstandingFlats(context.Context, models.MaintenanceOutstandingFlatFilter) (models.MaintenanceOutstandingFlatList, error)
	})
	if !ok {
		return models.MaintenanceOutstandingFlatList{}, errors.New("billing repository does not support outstanding flat lists")
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 25
	}
	return store.OutstandingFlats(ctx, f)
}
