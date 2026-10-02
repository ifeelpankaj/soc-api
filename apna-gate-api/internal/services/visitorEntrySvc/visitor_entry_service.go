package visitorentrysvc

import (
	"context"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	flatauthz "go-server/internal/services/flatAuthz"
	notificationsvc "go-server/internal/services/notificationSvc"
	shortlinksvc "go-server/internal/services/shortLinkSvc"
)

const defaultInviteDuration = 24 * time.Hour

type VisitorInviteService interface {
	CreateInvite(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, actorUserID int64, expiresAt *time.Time) (*models.VisitorEntryMutationResponse, *models.VisitorInvite, *models.ShortLinkResponse, error)
	CreateStaffInvite(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, staffUserID int64, expiresAt *time.Time) (*models.VisitorEntryMutationResponse, *models.VisitorInvite, *models.ShortLinkResponse, error)
	GetInviteByToken(ctx context.Context, societyID int64, guardUserID int64, rawToken string) (*models.VisitorInvite, error)
	GetPublicInviteByToken(ctx context.Context, rawToken string) (*models.PublicVisitorInvitePageResponse, error)
	GetPublicInviteByID(ctx context.Context, inviteID int64) (*models.PublicVisitorInvitePageResponse, error)
	SubmitInviteForm(ctx context.Context, societyID int64, rawToken string, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error)
	SubmitInviteFormByID(ctx context.Context, inviteID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error)
	CancelInvite(ctx context.Context, societyID int64, inviteID int64, actorUserID int64) error
	ListVisitorInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter models.VisitorInviteHistoryFilter) (*models.VisitorInviteHistoryResult, error)
	GetVisitorInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, inviteID int64, actorUserID int64) (*models.VisitorInviteHistoryItem, error)
	ExpireOldInvites(ctx context.Context) error
}

type VisitorEntryService interface {
	GetEntryOptions(ctx context.Context, societyID int64) (*models.VisitorEntryOptionsResponse, error)
	CreatePublicQREntry(ctx context.Context, societyID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error)
	CreateQuickLinkEntry(ctx context.Context, societyID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error)
	CreateGuardEntry(ctx context.Context, societyID int64, req models.VisitorFormRequest, guardUserID int64) (*models.VisitorEntryMutationResponse, error)
	ApproveEntry(ctx context.Context, societyID int64, entryID int64, actorUserID int64) (*models.VisitorEntryMutationResponse, error)
	RejectEntry(ctx context.Context, societyID int64, entryID int64, reason string, actorUserID int64) (*models.VisitorEntryMutationResponse, error)
	GenerateQR(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntryMutationResponse, error)
	GetEntryForQRScan(ctx context.Context, rawToken string) (*models.VisitorEntry, error)
	GetEntryForQRScanForViewer(ctx context.Context, rawToken string, viewerUserID *int64) (*models.VisitorEntry, error)
	ValidateQR(ctx context.Context, rawToken string) (*models.VisitorEntry, error)
	CheckIn(ctx context.Context, rawToken string, guardUserID int64) (*models.VisitorEntry, error)
	CheckOut(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error)
	AutoCloseExpiredEntries(ctx context.Context) error
	ExpireStaleEntries(ctx context.Context) error
	GetEntry(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error)
	ListEntries(ctx context.Context, filter models.VisitorEntryFilter) ([]*models.VisitorEntry, error)
	ListEntriesPaginated(ctx context.Context, filter models.VisitorEntryFilter) (*models.VisitorEntryListResult, error)
	GetEntryStats(ctx context.Context, societyID int64) (*models.VisitorEntryStatsResponse, error)
	GetEntryStatsInRange(ctx context.Context, societyID int64, from, to time.Time) (*models.VisitorEntryStatsResponse, error)
	GetDailyEntryStats(ctx context.Context, societyID int64, days int32) (*models.VisitorEntryDailyStatsResponse, error)
	GetGuardDeskBootstrap(ctx context.Context, societyID int64) (*models.GuardDeskBootstrapResponse, error)
	ListSocietyPendingApprovals(ctx context.Context, filter models.VisitorPendingFilter) (*models.VisitorPendingListResult, error)
	ListWaitingAtGate(ctx context.Context, filter models.WaitingAtGateFilter) (*models.VisitorEntryListResult, error)
	ListExpectedGuests(ctx context.Context, filter models.ExpectedGuestFilter) (*models.VisitorEntryListResult, error)
	NotifyPendingEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64) error
	GuardApproveEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64, opts GuardApproveOptions) (*models.VisitorEntryMutationResponse, error)
	GuardApproveAndCheckIn(ctx context.Context, societyID int64, entryID int64, guardUserID int64, opts GuardApproveOptions) (*models.VisitorEntry, error)
	CheckInByEntryID(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error)
	GetFlatVisitorContext(ctx context.Context, societyID int64, flatID int64) (*models.FlatVisitorContextResponse, error)
	GetFlatVisitorContextForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) (*models.FlatVisitorContextResponse, error)
	UpdateGuardEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64, req models.UpdateGuardVisitorEntryRequest) (*models.VisitorEntry, error)
	GetMemberVisitorApprovalStats(ctx context.Context, societyID int64, memberID int64) (*models.MemberVisitorApprovalStatsResponse, error)
	ListPendingApprovals(ctx context.Context, societyID int64, flatID int64, actorUserID int64) ([]*models.VisitorEntry, error)
	ListFlatEntriesForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter models.VisitorEntryFilter) (*models.VisitorEntryListResult, error)
	GetFlatEntryForActor(ctx context.Context, societyID int64, flatID int64, entryID int64, actorUserID int64) (*models.VisitorEntry, error)
	ListEvents(ctx context.Context, societyID int64, entryID int64) ([]*models.VisitorEntryEvent, error)
}

type VisitorEntrySvc struct {
	visitorRepo  repository.VisitorRepository
	inviteRepo   repository.VisitorInviteRepository
	entryRepo    repository.VisitorEntryRepository
	eventRepo    repository.VisitorEntryEventRepository
	settingSvc   visitorSettingRules
	memberRepo   repository.SocietyMemberRepository
	residentRepo repository.FlatResidentRepository
	flatRepo     repository.FlatRepository
	societyRepo  repository.SocietyRepository
	flatAuthz    *flatauthz.FlatVisitorAuthz
	notifier     notificationsvc.NotificationService
	shortLinks   shortlinksvc.ShortLinkService
	txManager    repository.TransactionManager
}

type visitorSettingRules interface {
	GetSocietySettings(ctx context.Context, societyID int64) (*models.SocietyVisitorSettingsResponse, error)
	GetFlatSettings(ctx context.Context, societyID int64, flatID int64) ([]models.FlatVisitorSettingsResponse, error)
	EnsureDefaultFlatSettingsIfMissing(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
	ResolveApprovalRequirement(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, source models.VisitorEntrySource) (bool, error)
	ResolveVisitDurationMinutes(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose) (int32, error)
}

func NewVisitorService(
	visitorRepo repository.VisitorRepository,
	inviteRepo repository.VisitorInviteRepository,
	entryRepo repository.VisitorEntryRepository,
	eventRepo repository.VisitorEntryEventRepository,
	settingSvc visitorSettingRules,
	memberRepo repository.SocietyMemberRepository,
	residentRepo repository.FlatResidentRepository,
	flatRepo repository.FlatRepository,
	societyRepo repository.SocietyRepository,
	flatAuthz *flatauthz.FlatVisitorAuthz,
	notifier notificationsvc.NotificationService,
	deps ...any,
) (VisitorInviteService, VisitorEntryService) {
	svc := &VisitorEntrySvc{
		visitorRepo: visitorRepo, inviteRepo: inviteRepo, entryRepo: entryRepo, eventRepo: eventRepo,
		settingSvc: settingSvc, memberRepo: memberRepo, residentRepo: residentRepo, flatRepo: flatRepo,
		societyRepo: societyRepo, flatAuthz: flatAuthz, notifier: notifier,
	}
	for _, dep := range deps {
		switch value := dep.(type) {
		case shortlinksvc.ShortLinkService:
			svc.shortLinks = value
		case repository.TransactionManager:
			svc.txManager = value
		}
	}
	return svc, svc
}

var _ VisitorInviteService = (*VisitorEntrySvc)(nil)
var _ VisitorEntryService = (*VisitorEntrySvc)(nil)
