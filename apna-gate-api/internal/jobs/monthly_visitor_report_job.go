package jobs

import (
	"context"
	"errors"
	"fmt"
	"go-server/pkg/logger"
	"strings"
	"time"

	"go-server/internal/models"

	"go.uber.org/zap"
)

const monthlyVisitorReportJobLock = "apna-gate/jobs/monthly-visitor-report"

type monthlyVisitorReportStore interface {
	WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error)
	ListActiveReportSocieties(ctx context.Context) ([]models.VisitorReportSociety, error)
	GetEarliestVisitorEntryCreatedAt(ctx context.Context, societyID int64) (*time.Time, error)
	EnsureReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time) error
	ListClaimableReportDeliveries(ctx context.Context, batchSize int32, allowResend bool, resendMonth time.Time) ([]models.VisitorReportDelivery, error)
	ClaimReportDelivery(ctx context.Context, societyID int64, reportMonth, processingUntil time.Time, allowResend bool) (*models.VisitorReportDelivery, error)
	CompleteReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time, recipients []string, providerMessageID string, allowResend bool) error
	FailReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time, recipients []string, cause error, allowResend bool) error
	ListMonthlyVisitorReportRecipients(ctx context.Context, societyID int64) ([]string, error)
	ListMonthlyVisitorReportRows(ctx context.Context, societyID int64, from, to time.Time) ([]models.MonthlyVisitorReportRow, error)
}

type visitorReportEmailSender interface {
	SendMonthlyVisitorReport(ctx context.Context, recipients []string, societyName string, reportMonth time.Time, filename string, content []byte, idempotencyKey string) (string, error)
}

type MonthlyVisitorReportJobConfig struct {
	Location        *time.Location
	BatchSize       int32
	ProcessingLease time.Duration
	AllowResend     bool
}

type MonthlyVisitorReportJob struct {
	repo  monthlyVisitorReportStore
	email visitorReportEmailSender
	cfg   MonthlyVisitorReportJobConfig
	now   func() time.Time
}

func NewMonthlyVisitorReportJob(repo monthlyVisitorReportStore, email visitorReportEmailSender, cfg MonthlyVisitorReportJobConfig) *MonthlyVisitorReportJob {
	if cfg.Location == nil {
		cfg.Location = time.UTC
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 50
	}
	if cfg.ProcessingLease <= 0 {
		cfg.ProcessingLease = 15 * time.Minute
	}
	return &MonthlyVisitorReportJob{repo: repo, email: email, cfg: cfg, now: time.Now}
}

func (j *MonthlyVisitorReportJob) RunOnce(ctx context.Context) (RunResult, error) {
	locked, err := j.repo.WithAdvisoryLock(ctx, monthlyVisitorReportJobLock, func(lockCtx context.Context) error {
		latest := j.latestEligibleMonth(j.now())
		if err := j.ensureDeliveries(lockCtx, latest); err != nil {
			return err
		}
		deliveries, err := j.repo.ListClaimableReportDeliveries(lockCtx, j.cfg.BatchSize, j.cfg.AllowResend, latest)
		if err != nil {
			return fmt.Errorf("list claimable report deliveries: %w", err)
		}
		var deliveryErrors []error
		for _, delivery := range deliveries {
			if err := lockCtx.Err(); err != nil {
				return err
			}
			if err := j.deliver(lockCtx, delivery); err != nil {
				if isContextCancellation(err) {
					return err
				}
				deliveryErrors = append(deliveryErrors, err)
			}
		}
		return errors.Join(deliveryErrors...)
	})
	result := RunResult{LockAcquired: locked}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (j *MonthlyVisitorReportJob) latestEligibleMonth(now time.Time) time.Time {
	local := now.In(j.cfg.Location)
	threshold := time.Date(local.Year(), local.Month(), 1, 0, 5, 0, 0, j.cfg.Location)
	monthsBack := -1
	if local.Before(threshold) {
		monthsBack = -2
	}
	return ReportMonth(local.AddDate(0, monthsBack, 0), j.cfg.Location)
}

func (j *MonthlyVisitorReportJob) ensureDeliveries(ctx context.Context, latest time.Time) error {
	societies, err := j.repo.ListActiveReportSocieties(ctx)
	if err != nil {
		return fmt.Errorf("list active report societies: %w", err)
	}
	for _, society := range societies {
		earliest, err := j.repo.GetEarliestVisitorEntryCreatedAt(ctx, society.ID)
		if err != nil {
			return fmt.Errorf("get earliest visitor entry for society %d: %w", society.ID, err)
		}
		start := latest
		if earliest != nil {
			earliestMonth := ReportMonth(*earliest, j.cfg.Location)
			if earliestMonth.Before(latest) {
				start = earliestMonth
			}
		}
		for month := start; !month.After(latest); month = month.AddDate(0, 1, 0) {
			if err := j.repo.EnsureReportDelivery(ctx, society.ID, month); err != nil {
				return fmt.Errorf("ensure visitor report delivery society=%d month=%s: %w", society.ID, month.Format("2006-01"), err)
			}
		}
	}
	return nil
}

func (j *MonthlyVisitorReportJob) deliver(ctx context.Context, delivery models.VisitorReportDelivery) error {
	claimed, err := j.repo.ClaimReportDelivery(ctx, delivery.SocietyID, delivery.ReportMonth, j.now().Add(j.cfg.ProcessingLease), j.cfg.AllowResend)
	if err != nil {
		return fmt.Errorf("claim report delivery society=%d month=%s: %w", delivery.SocietyID, delivery.ReportMonth.Format("2006-01"), err)
	}
	if claimed == nil {
		return nil
	}

	recipients, err := j.repo.ListMonthlyVisitorReportRecipients(ctx, delivery.SocietyID)
	if err == nil && len(recipients) == 0 && delivery.SocietyEmail != nil {
		fallback := strings.ToLower(strings.TrimSpace(*delivery.SocietyEmail))
		if fallback != "" {
			recipients = []string{fallback}
		}
	}
	if err == nil && len(recipients) == 0 {
		err = errors.New("society has no active owner email or society fallback email")
	}
	if err != nil {
		if isContextCancellation(err) {
			return err
		}
		return j.failDelivery(ctx, delivery, recipients, err)
	}

	from, to := MonthBounds(delivery.ReportMonth, j.cfg.Location)
	rows, err := j.repo.ListMonthlyVisitorReportRows(ctx, delivery.SocietyID, from, to)
	if err != nil {
		if isContextCancellation(err) {
			return err
		}
		return j.failDelivery(ctx, delivery, recipients, err)
	}
	csvContent, err := GenerateVisitorReportCSV(rows)
	if err != nil {
		return j.failDelivery(ctx, delivery, recipients, err)
	}

	filename := fmt.Sprintf("visitor-report-%d-%s.csv", delivery.SocietyID, delivery.ReportMonth.Format("2006-01"))
	idempotencyKey := fmt.Sprintf("monthly-visitor-report/%d/%s", delivery.SocietyID, delivery.ReportMonth.Format("2006-01"))
	if j.cfg.AllowResend {
		idempotencyKey = fmt.Sprintf("%s/attempt-%d", idempotencyKey, claimed.AttemptCount)
	}
	providerID, err := j.email.SendMonthlyVisitorReport(
		ctx, recipients, delivery.SocietyName, delivery.ReportMonth, filename, csvContent, idempotencyKey,
	)
	if err != nil {
		if isContextCancellation(err) {
			return err
		}
		return j.failDelivery(ctx, delivery, recipients, err)
	}
	if err := j.repo.CompleteReportDelivery(ctx, delivery.SocietyID, delivery.ReportMonth, recipients, providerID, j.cfg.AllowResend); err != nil {
		return fmt.Errorf("email accepted but report delivery could not be completed society=%d month=%s: %w", delivery.SocietyID, delivery.ReportMonth.Format("2006-01"), err)
	}
	jobLogger(ctx).Info("monthly visitor report sent", reportFields(delivery, nil)...)
	return nil
}

func (j *MonthlyVisitorReportJob) failDelivery(ctx context.Context, delivery models.VisitorReportDelivery, recipients []string, cause error) error {
	if err := j.repo.FailReportDelivery(ctx, delivery.SocietyID, delivery.ReportMonth, recipients, cause, j.cfg.AllowResend); err != nil {
		return errors.Join(cause, fmt.Errorf("record report delivery failure: %w", err))
	}
	jobLogger(ctx).Error("monthly visitor report delivery failed", reportFields(delivery, cause)...)
	return cause
}

func reportFields(delivery models.VisitorReportDelivery, err error) []zap.Field {
	fields := []zap.Field{
		zap.Int64("society_id", delivery.SocietyID),
		zap.String("report_month", delivery.ReportMonth.Format("2006-01")),
	}
	if err != nil {
		fields = append(fields, zap.String("internal_error", logger.Sanitize(err.Error())))
	}
	return fields
}
