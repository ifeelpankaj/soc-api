package contracts

import (
	"context"
	"go-server/internal/models"
	"time"
)

type JobRepository interface {
	DeleteOrphanVisitorsBatch(context.Context, time.Time, int32) (int64, error)
	ListVisitorImageDeletions(context.Context, int64, int32) ([]models.PendingImageDeletion, error)
	CompleteVisitorImageDeletion(context.Context, int64) error
	WithAdvisoryLock(ctx context.Context, name string, fn func(context.Context) error) (bool, error)
	ExpireWaitingVisitorEntries(ctx context.Context) (int64, error)
	ExpireApprovedVisitorEntries(ctx context.Context) (int64, error)
	ExpireVisitorInvites(ctx context.Context) (int64, error)
	ExpireFlatMemberInvites(ctx context.Context) (int64, error)
	ExpireSubscriptions(ctx context.Context) (int64, error)
	DeleteNotificationsBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
	DeleteVerificationsBatch(ctx context.Context, batchSize int32) (int64, error)
	ListVisitorEntryCleanupWindows(ctx context.Context, timezone string, cutoff time.Time) ([]models.VisitorReportCleanupWindow, error)
	DeleteVisitorEntriesBatch(ctx context.Context, window models.VisitorReportCleanupWindow, from, to, cutoff time.Time, batchSize int32) (int64, error)
	DeleteVisitorInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
	DeleteFlatMemberInvitesBatch(ctx context.Context, cutoff time.Time, batchSize int32) (int64, error)
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
