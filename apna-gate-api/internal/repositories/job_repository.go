package repository

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"
	"go-server/pkg/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	visitorInviteResourceType = "visitor_invite"
	memberInviteResourceType  = "member_invite"
)

type JobRepository = contracts.JobRepository

type jobRepository struct {
	db        *database.Database
	txManager TransactionManager
}

func NewJobRepository(database *database.Database, txManager TransactionManager) JobRepository {
	return &jobRepository{db: database, txManager: txManager}
}

func (r *jobRepository) WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error) {
	conn, err := r.db.Pool.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("acquire advisory lock connection: %w", err)
	}
	defer conn.Release()

	var locked bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock(hashtextextended($1, 0))", name).Scan(&locked); err != nil {
		return false, fmt.Errorf("acquire advisory lock %q: %w", name, err)
	}
	if !locked {
		return false, nil
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = conn.Exec(unlockCtx, "SELECT pg_advisory_unlock(hashtextextended($1, 0))", name)
	}()

	return true, fn(ctx)
}

func (r *jobRepository) ExpireWaitingVisitorEntries(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireWaitingVisitorEntriesForJob(ctx)
}

func (r *jobRepository) ExpireApprovedVisitorEntries(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireApprovedVisitorEntriesForJob(ctx)
}

func (r *jobRepository) ExpireVisitorInvites(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireVisitorInvitesForJob(ctx)
}

func (r *jobRepository) ExpireFlatMemberInvites(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireFlatMemberInvitesForJob(ctx)
}

func (r *jobRepository) ExpireSubscriptions(ctx context.Context) (int64, error) {
	return GetQueries(ctx, r.db).ExpireSubscriptionsForJob(ctx)
}

func (r *jobRepository) DeleteNotificationsBatch(ctx context.Context, cutoff time.Time, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		count, err = GetQueries(txCtx, r.db).DeleteNotificationsBatch(txCtx, db.DeleteNotificationsBatchParams{
			Cutoff: timestampParam(cutoff), BatchSize: batchSize,
		})
		return err
	})
	return count, err
}

func (r *jobRepository) DeleteVerificationsBatch(ctx context.Context, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		count, err = GetQueries(txCtx, r.db).DeleteVerificationsBatch(txCtx, batchSize)
		return err
	})
	return count, err
}

func (r *jobRepository) ListVisitorEntryCleanupWindows(ctx context.Context, timezone string, cutoff time.Time) ([]models.VisitorReportCleanupWindow, error) {
	rows, err := GetQueries(ctx, r.db).ListVisitorEntryCleanupWindows(ctx, db.ListVisitorEntryCleanupWindowsParams{
		JobTimezone: timezone, Cutoff: timestampParam(cutoff),
	})
	if err != nil {
		return nil, err
	}
	windows := make([]models.VisitorReportCleanupWindow, 0, len(rows))
	for _, row := range rows {
		windows = append(windows, models.VisitorReportCleanupWindow{
			SocietyID: row.SocietyID, ReportMonth: row.ReportMonth.Time,
		})
	}
	return windows, nil
}

func (r *jobRepository) DeleteVisitorEntriesBatch(ctx context.Context, window models.VisitorReportCleanupWindow, from, to, cutoff time.Time, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		queries := GetQueries(txCtx, r.db)
		ids, queryErr := queries.DeleteVisitorEntriesBatch(txCtx, db.DeleteVisitorEntriesBatchParams{
			SocietyID: window.SocietyID, PeriodStart: timestampParam(from), PeriodEnd: timestampParam(to),
			Cutoff: timestampParam(cutoff), ReportMonth: dateParam(window.ReportMonth), BatchSize: batchSize,
		})
		if queryErr != nil {
			return queryErr
		}
		count = int64(len(ids))
		locked, queryErr := queries.LockCleanupVisitors(txCtx, ids)
		if queryErr != nil {
			return queryErr
		}
		return r.deleteLockedVisitors(txCtx, locked)
	})
	return count, err
}

// Call only after locking the visitor rows. The separate delete statement gets a
// fresh READ COMMITTED snapshot after any concurrent entry insertion completes.
func (r *jobRepository) deleteLockedVisitors(ctx context.Context, ids []int64) error {
	queries := GetQueries(ctx, r.db)
	rows, err := queries.DeleteUnreferencedVisitors(ctx, ids)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row.PhotoImagekitFileID != nil && *row.PhotoImagekitFileID != "" {
			if err := queries.QueueVisitorImageDeletion(ctx, *row.PhotoImagekitFileID); err != nil {
				return err
			}
		} else if row.PhotoUrl != nil && *row.PhotoUrl != "" {
			logger.Warn("deleted visitor has a legacy photo without an ImageKit file ID", zap.Int64("visitor_id", row.ID))
		}
	}
	return nil
}

func (r *jobRepository) DeleteOrphanVisitorsBatch(ctx context.Context, cutoff time.Time, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		ids, queryErr := GetQueries(txCtx, r.db).LockOldOrphanVisitors(txCtx, db.LockOldOrphanVisitorsParams{Cutoff: timestampParam(cutoff), BatchSize: batchSize})
		if queryErr != nil {
			return queryErr
		}
		queryErr = r.deleteLockedVisitors(txCtx, ids)
		count = int64(len(ids))
		return queryErr
	})
	return count, err
}

func (r *jobRepository) ListVisitorImageDeletions(ctx context.Context, afterID int64, batchSize int32) ([]models.PendingImageDeletion, error) {
	rows, err := GetQueries(ctx, r.db).ListVisitorImageDeletions(ctx, db.ListVisitorImageDeletionsParams{AfterID: afterID, BatchSize: batchSize})
	if err != nil {
		return nil, err
	}
	items := make([]models.PendingImageDeletion, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.PendingImageDeletion{ID: row.ID, FileID: row.FileID})
	}
	return items, nil
}

func (r *jobRepository) CompleteVisitorImageDeletion(ctx context.Context, id int64) error {
	return GetQueries(ctx, r.db).CompleteVisitorImageDeletion(ctx, id)
}

func (r *jobRepository) DeleteVisitorInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		queries := GetQueries(txCtx, r.db)
		ids, queryErr := queries.ListVisitorInviteCleanupCandidateIDs(txCtx, db.ListVisitorInviteCleanupCandidateIDsParams{
			Cutoff: timestampParam(cutoff), BatchSize: batchSize,
		})
		if queryErr != nil || len(ids) == 0 {
			return queryErr
		}
		if _, queryErr = queries.DeleteShortLinksForResources(txCtx, db.DeleteShortLinksForResourcesParams{
			ResourceType: visitorInviteResourceType, ResourceIds: ids,
		}); queryErr != nil {
			return queryErr
		}
		count, queryErr = queries.DeleteVisitorInvitesByIDs(txCtx, ids)
		return queryErr
	})
	return count, err
}

func (r *jobRepository) DeleteFlatMemberInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (count int64, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		queries := GetQueries(txCtx, r.db)
		ids, queryErr := queries.ListFlatMemberInviteCleanupCandidateIDs(txCtx, db.ListFlatMemberInviteCleanupCandidateIDsParams{
			Cutoff: timestampParam(cutoff), BatchSize: batchSize,
		})
		if queryErr != nil || len(ids) == 0 {
			return queryErr
		}
		if _, queryErr = queries.DeleteShortLinksForResources(txCtx, db.DeleteShortLinksForResourcesParams{
			ResourceType: memberInviteResourceType, ResourceIds: ids,
		}); queryErr != nil {
			return queryErr
		}
		count, queryErr = queries.DeleteFlatMemberInvitesByIDs(txCtx, ids)
		return queryErr
	})
	return count, err
}

func (r *jobRepository) ListActiveReportSocieties(ctx context.Context) ([]models.VisitorReportSociety, error) {
	rows, err := GetQueries(ctx, r.db).ListActiveSocietiesForVisitorReports(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]models.VisitorReportSociety, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.VisitorReportSociety{ID: row.ID, Name: row.Name, Email: row.Email})
	}
	return items, nil
}

func (r *jobRepository) GetEarliestVisitorEntryCreatedAt(ctx context.Context, societyID int64) (*time.Time, error) {
	value, err := GetQueries(ctx, r.db).GetEarliestVisitorEntryCreatedAt(ctx, societyID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return timeFromTimestamp(value), nil
}

func (r *jobRepository) EnsureReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time) error {
	return GetQueries(ctx, r.db).EnsureMonthlyVisitorReportDelivery(ctx, db.EnsureMonthlyVisitorReportDeliveryParams{
		SocietyID: societyID, ReportMonth: dateParam(reportMonth),
	})
}

func (r *jobRepository) ListClaimableReportDeliveries(ctx context.Context, batchSize int32, allowResend bool, resendMonth time.Time) ([]models.VisitorReportDelivery, error) {
	rows, err := GetQueries(ctx, r.db).ListClaimableMonthlyVisitorReportDeliveries(ctx, db.ListClaimableMonthlyVisitorReportDeliveriesParams{AllowResend: allowResend, ResendMonth: dateParam(resendMonth), BatchSize: batchSize})
	if err != nil {
		return nil, err
	}
	items := make([]models.VisitorReportDelivery, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.VisitorReportDelivery{
			SocietyID: row.SocietyID, SocietyName: row.SocietyName, SocietyEmail: row.SocietyEmail,
			ReportMonth: row.ReportMonth.Time, AttemptCount: row.AttemptCount,
			ProcessingUntil: timeFromTimestamp(row.ProcessingUntil), SentAt: timeFromTimestamp(row.SentAt),
		})
	}
	return items, nil
}

func (r *jobRepository) ClaimReportDelivery(ctx context.Context, societyID int64, reportMonth, processingUntil time.Time, allowResend bool) (delivery *models.VisitorReportDelivery, err error) {
	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		queries := GetQueries(txCtx, r.db)
		current, queryErr := queries.GetMonthlyVisitorReportDeliveryForUpdate(txCtx, db.GetMonthlyVisitorReportDeliveryForUpdateParams{
			SocietyID: societyID, ReportMonth: dateParam(reportMonth),
		})
		if queryErr != nil {
			return queryErr
		}
		if (!allowResend && current.SentAt.Valid) || (current.ProcessingUntil.Valid && current.ProcessingUntil.Time.After(time.Now())) {
			return nil
		}
		claimed, queryErr := queries.ClaimMonthlyVisitorReportDelivery(txCtx, db.ClaimMonthlyVisitorReportDeliveryParams{
			SocietyID: societyID, ReportMonth: dateParam(reportMonth), ProcessingUntil: timestampParam(processingUntil),
			AllowResend: allowResend,
		})
		if errors.Is(queryErr, pgx.ErrNoRows) {
			return nil
		}
		if queryErr != nil {
			return queryErr
		}
		delivery = &models.VisitorReportDelivery{
			SocietyID: claimed.SocietyID, ReportMonth: claimed.ReportMonth.Time,
			AttemptCount: claimed.AttemptCount, ProcessingUntil: timeFromTimestamp(claimed.ProcessingUntil),
			SentAt: timeFromTimestamp(claimed.SentAt),
		}
		return nil
	})
	return delivery, err
}

func (r *jobRepository) CompleteReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time, recipients []string, providerMessageID string, allowResend bool) error {
	return r.finishReportDelivery(ctx, societyID, reportMonth, allowResend, func(txCtx context.Context, queries *db.Queries) error {
		count, err := queries.CompleteMonthlyVisitorReportDelivery(txCtx, db.CompleteMonthlyVisitorReportDeliveryParams{
			SocietyID: societyID, ReportMonth: dateParam(reportMonth), Recipients: recipients,
			ProviderMessageID: stringPtrOrNil(providerMessageID),
			AllowResend:       allowResend,
		})
		if err == nil && count != 1 {
			return fmt.Errorf("complete report delivery affected %d rows", count)
		}
		return err
	})
}

func (r *jobRepository) FailReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time, recipients []string, cause error, allowResend bool) error {
	message := "unknown report delivery error"
	if cause != nil {
		message = cause.Error()
	}
	return r.finishReportDelivery(ctx, societyID, reportMonth, allowResend, func(txCtx context.Context, queries *db.Queries) error {
		_, err := queries.FailMonthlyVisitorReportDelivery(txCtx, db.FailMonthlyVisitorReportDeliveryParams{
			SocietyID: societyID, ReportMonth: dateParam(reportMonth), Recipients: recipients, LastError: &message, AllowResend: allowResend,
		})
		return err
	})
}

func (r *jobRepository) finishReportDelivery(ctx context.Context, societyID int64, reportMonth time.Time, allowResend bool, fn func(context.Context, *db.Queries) error) error {
	return r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		queries := GetQueries(txCtx, r.db)
		current, err := queries.GetMonthlyVisitorReportDeliveryForUpdate(txCtx, db.GetMonthlyVisitorReportDeliveryForUpdateParams{
			SocietyID: societyID, ReportMonth: dateParam(reportMonth),
		})
		if err != nil {
			return err
		}
		if current.SentAt.Valid && !allowResend {
			return nil
		}
		return fn(txCtx, queries)
	})
}

func (r *jobRepository) ListMonthlyVisitorReportRecipients(ctx context.Context, societyID int64) ([]string, error) {
	return GetQueries(ctx, r.db).ListMonthlyVisitorReportRecipients(ctx, societyID)
}

func (r *jobRepository) ListMonthlyVisitorReportRows(ctx context.Context, societyID int64, from, to time.Time) ([]models.MonthlyVisitorReportRow, error) {
	rows, err := GetQueries(ctx, r.db).ListMonthlyVisitorReportRows(ctx, db.ListMonthlyVisitorReportRowsParams{
		SocietyID: societyID, PeriodStart: timestampParam(from), PeriodEnd: timestampParam(to),
	})
	if err != nil {
		return nil, err
	}
	items := make([]models.MonthlyVisitorReportRow, 0, len(rows))
	for _, row := range rows {
		var vehicleType *string
		if row.VehicleType != nil {
			value := string(*row.VehicleType)
			vehicleType = &value
		}
		items = append(items, models.MonthlyVisitorReportRow{
			ID: row.ID, VisitorName: row.VisitorName, VisitorPhone: row.VisitorPhone, VisitorEmail: row.VisitorEmail,
			FlatNumber: row.FlatNumber, Block: row.Block, Floor: row.Floor, Source: string(row.Source), Purpose: string(row.Purpose),
			Status: string(row.Status), ExpectedAt: timeFromTimestamp(row.ExpectedAt), ExpectedCheckoutAt: timeFromTimestamp(row.ExpectedCheckoutAt),
			CheckedInAt: timeFromTimestamp(row.CheckedInAt), CheckedOutAt: timeFromTimestamp(row.CheckedOutAt), AutoClosedAt: timeFromTimestamp(row.AutoClosedAt),
			VehicleNumber: row.VehicleNumber, VehicleType: vehicleType, CompanionsCount: row.CompanionsCount,
			CompanionDetails: string(row.CompanionDetails),
			ApproverName:     row.ApproverName, GuardName: row.GuardName, RejectionReason: row.RejectionReason, Notes: row.Notes,
			CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return items, nil
}

func timestampParam(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func dateParam(value time.Time) pgtype.Date {
	return pgtype.Date{Time: time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC), Valid: true}
}

func timeFromTimestamp(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
