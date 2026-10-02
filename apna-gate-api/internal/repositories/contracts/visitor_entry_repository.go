package contracts

import (
	"context"
	"go-server/internal/models"
	"time"
)

type VisitorRepository interface {
	Create(ctx context.Context, req VisitorFormInput) (*models.Visitor, error)
	Get(ctx context.Context, id int64) (*models.Visitor, error)
	UpdateProfile(ctx context.Context, visitorID int64, req UpdateGuardVisitorEntryInput) (*models.Visitor, error)
}

type VisitorInviteRepository interface {
	Create(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, tokenHash string, expiresAt time.Time, actorUserID int64) (*models.VisitorInvite, error)
	GetByID(ctx context.Context, societyID int64, inviteID int64) (*models.VisitorInvite, error)
	GetHistoryByID(ctx context.Context, societyID int64, flatID int64, inviteID int64, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) (*models.VisitorInviteHistoryItem, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.VisitorInvite, error)
	ListHistory(ctx context.Context, societyID int64, flatID int64, filter models.VisitorInviteHistoryFilter, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) ([]*models.VisitorInviteHistoryItem, error)
	CountHistory(ctx context.Context, societyID int64, flatID int64, filter models.VisitorInviteHistoryFilter) (int64, error)
	MarkUsed(ctx context.Context, inviteID int64) (*models.VisitorInvite, error)
	GetForUpdate(ctx context.Context, inviteID int64) (*models.VisitorInvite, error)
	Cancel(ctx context.Context, societyID int64, inviteID int64) (*models.VisitorInvite, error)
	ExpireOld(ctx context.Context) error
}

type VisitorEntryRepository interface {
	Create(ctx context.Context, req VisitorFormInput, societyID int64, flatID *int64, visitorID int64, inviteID *int64, source models.VisitorEntrySource, purpose models.VisitorPurpose, status models.VisitorStatus, actorUserID *int64, guardUserID *int64, qrHash *string, qrExpiresAt *time.Time) (*models.VisitorEntry, error)
	Get(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error)
	GetForUpdate(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error)
	GetByQRHash(ctx context.Context, qrHash string) (*models.VisitorEntry, error)
	GetByInviteID(ctx context.Context, inviteID int64) (*models.VisitorEntry, error)
	List(ctx context.Context, filter models.VisitorEntryFilter) ([]*models.VisitorEntry, error)
	Count(ctx context.Context, filter models.VisitorEntryFilter) (int64, error)
	ListPending(ctx context.Context, societyID int64, flatID int64) ([]*models.VisitorEntry, error)
	ListSocietyPending(ctx context.Context, filter models.VisitorPendingFilter) ([]*models.VisitorPendingEntry, error)
	CountSocietyPending(ctx context.Context, filter models.VisitorPendingFilter) (int64, error)
	ListRecentByFlat(ctx context.Context, societyID int64, flatID int64, limit int32) ([]*models.VisitorEntry, error)
	GetStats(ctx context.Context, societyID int64) (*models.VisitorEntryStatsResponse, error)
	GetStatsInRange(ctx context.Context, societyID int64, from, to time.Time) (*models.VisitorEntryStatsResponse, error)
	GetDailyStatsCreated(ctx context.Context, societyID int64, days int32) ([]models.VisitorDailyCountResponse, error)
	CountWaitingAtGate(ctx context.Context, societyID int64) (int64, error)
	ListWaitingAtGate(ctx context.Context, filter models.WaitingAtGateFilter) ([]*models.VisitorEntry, error)
	CountWaitingAtGateFiltered(ctx context.Context, filter models.WaitingAtGateFilter) (int64, error)
	CountExpectedGuests(ctx context.Context, societyID int64, fromAt, toAt time.Time) (int64, error)
	ListExpectedGuests(ctx context.Context, filter models.ExpectedGuestFilter) ([]*models.VisitorEntry, error)
	CountExpectedGuestsFiltered(ctx context.Context, filter models.ExpectedGuestFilter) (int64, error)
	CountMemberApprovals(ctx context.Context, societyID int64, userID int64) (*models.MemberVisitorApprovalStatsResponse, error)
	Approve(ctx context.Context, societyID int64, entryID int64, actorUserID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error)
	MergeMetadata(ctx context.Context, societyID int64, entryID int64, metadata map[string]any) (*models.VisitorEntry, error)
	Reject(ctx context.Context, societyID int64, entryID int64, actorUserID int64, reason string) (*models.VisitorEntry, error)
	GenerateQR(ctx context.Context, societyID int64, entryID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error)
	CheckIn(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error)
	CheckOut(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error)
	UpdateDetails(ctx context.Context, societyID int64, entryID int64, req UpdateGuardVisitorEntryInput) (*models.VisitorEntry, error)
	AutoCloseExpired(ctx context.Context) error
	ExpireStaleEntries(ctx context.Context) error
}

type VisitorEntryEventRepository interface {
	Create(ctx context.Context, entryID int64, societyID int64, actorUserID *int64, eventType models.VisitorEventType, message *string, metadata map[string]any) (*models.VisitorEntryEvent, error)
	List(ctx context.Context, societyID int64, entryID int64) ([]*models.VisitorEntryEvent, error)
}
