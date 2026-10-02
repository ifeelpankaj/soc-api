package flatsvc

import (
	"context"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	flatauthz "go-server/internal/services/flatAuthz"
	notificationsvc "go-server/internal/services/notificationSvc"
	shortlinksvc "go-server/internal/services/shortLinkSvc"
	societysvc "go-server/internal/services/societySvc"
)

type FlatService interface {
	FlatCommandService
	FlatQueryService
	FlatClaimCommandService
	FlatClaimQueryService
	FlatResidentCommandService
	FlatResidentQueryService
	FlatMemberInviteService
}

type FlatCommandService interface {
	CreateFlat(ctx context.Context, societyID int64, createdBy int64, req *models.CreateFlatRequest) (*models.FlatResponse, error)
	BulkCreateFlats(ctx context.Context, societyID int64, createdBy int64, req *models.BulkCreateFlatsRequest) (*models.BulkCreateFlatsResponse, error)
	GenerateFlats(ctx context.Context, societyID int64, createdBy int64, req *models.GenerateFlatsRequest) (*models.BulkCreateFlatsResponse, error)
	UpdateFlat(ctx context.Context, filter *models.FlatFilter, req *models.UpdateFlatRequest) (*models.FlatResponse, error)
	DeleteFlat(ctx context.Context, filter *models.FlatFilter, deletedBy int64) error
	BlockFlat(ctx context.Context, filter *models.FlatFilter, blockedBy int64, reason string) (*models.FlatResponse, error)
	UnblockFlat(ctx context.Context, filter *models.FlatFilter, unblockedBy int64) (*models.FlatResponse, error)
}

type FlatQueryService interface {
	GetFlat(ctx context.Context, filter *models.FlatFilter) (*models.FlatResponse, error)
	ListFlats(ctx context.Context, filter *models.FlatFilter) ([]*models.FlatResponse, error)
	ListFlatsPaginated(ctx context.Context, filter *models.FlatFilter) (*models.PaginatedFlatsResponse, error)
	GetFlatStats(ctx context.Context, societyID int64) (*models.FlatStatsResponse, error)
}

type FlatClaimCommandService interface {
	SubmitFlatClaim(ctx context.Context, userID int64, req *models.SubmitFlatClaimRequest) (*models.FlatClaimResponse, error)
	ApproveFlatClaim(ctx context.Context, societyID int64, claimID int64, reviewedBy int64) (*models.ApproveFlatClaimResponse, error)
	RejectFlatClaim(ctx context.Context, societyID int64, claimID int64, reviewedBy int64, req *models.RejectFlatClaimRequest) (*models.FlatClaimResponse, error)
	CancelMyFlatClaim(ctx context.Context, claimID int64, userID int64) (*models.FlatClaimResponse, error)
}

type FlatClaimQueryService interface {
	GetFlatClaim(ctx context.Context, filter *models.FlatClaimFilter) (*models.FlatClaimResponse, error)
	ListFlatClaims(ctx context.Context, filter *models.FlatClaimFilter) ([]*models.FlatClaimResponse, error)
	ListMyFlatClaims(ctx context.Context, userID int64, filter *models.FlatClaimFilter) ([]*models.FlatClaimResponse, error)
	GetFlatClaimStats(ctx context.Context, societyID int64) (*models.FlatClaimStatsResponse, error)
}

type FlatResidentCommandService interface {
	AddResidentToFlat(ctx context.Context, societyID int64, flatID int64, userID int64, createdBy int64, req *models.AddFlatResidentRequest) (*models.FlatResidentResponse, error)
	RemoveResidentFromFlat(ctx context.Context, filter *models.FlatResidentFilter, removedBy int64) error
	ChangePrimaryResident(ctx context.Context, societyID int64, flatID int64, residentID int64, changedBy int64) (*models.FlatResidentResponse, error)
	UpdateFlatResidentRole(ctx context.Context, filter *models.FlatResidentFilter, updatedBy int64, req *models.UpdateFlatResidentRoleRequest) (*models.FlatResidentResponse, error)
	MoveOutResident(ctx context.Context, filter *models.FlatResidentFilter, movedOutBy int64) (*models.FlatResidentResponse, error)
}

type FlatResidentQueryService interface {
	GetFlatResident(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResidentResponse, error)
	ListFlatResidents(ctx context.Context, filter *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error)
	ListFlatResidentsForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error)
	ListMyResidences(ctx context.Context, userID int64, filter *models.FlatResidentFilter) ([]*models.FlatResidentResponse, error)
}

type FlatMemberInviteService interface {
	ListPendingMemberInvites(ctx context.Context, societyID int64, flatID int64, actorUserID int64) ([]*models.FlatMemberInviteResponse, error)
	ListMemberInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter models.FlatMemberInviteHistoryFilter) (*models.FlatMemberInviteHistoryResult, error)
	GetMemberInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, inviteID int64, actorUserID int64) (*models.FlatMemberInviteResponse, error)
	CreateMemberInvite(ctx context.Context, societyID int64, flatID int64, actorUserID int64, req *models.CreateFlatMemberInviteRequest) (*models.FlatMemberInviteTokenResponse, *models.FlatMemberInviteResponse, *models.ShortLinkResponse, error)
	CancelMemberInvite(ctx context.Context, societyID int64, flatID int64, inviteID int64, actorUserID int64) error
	GetPublicMemberInviteByToken(ctx context.Context, rawToken string) (*models.PublicFlatMemberInviteView, error)
	GetPublicMemberInviteByID(ctx context.Context, inviteID int64) (*models.PublicFlatMemberInviteView, error)
	JoinMemberInvite(ctx context.Context, rawToken string, req *models.JoinFlatMemberInviteRequest) (*models.JoinFlatMemberInviteResponse, error)
	JoinMemberInviteByID(ctx context.Context, inviteID int64, req *models.JoinFlatMemberInviteRequest) (*models.JoinFlatMemberInviteResponse, error)
	AcceptMemberInvite(ctx context.Context, rawToken string, userID int64) (*models.AcceptFlatMemberInviteResponse, error)
	AcceptMemberInviteByID(ctx context.Context, inviteID int64, userID int64) (*models.AcceptFlatMemberInviteResponse, error)
	ExpireOldMemberInvites(ctx context.Context) error
}

type FlatSvc struct {
	flatRepo                      repository.FlatRepository
	claimRepo                     repository.FlatClaimRepository
	residentRepo                  repository.FlatResidentRepository
	memberInviteRepo              repository.FlatMemberInviteRepository
	memberRepo                    repository.SocietyMemberRepository
	txManager                     repository.TransactionManager
	societySvc                    societysvc.SocietyService
	subscriptionSvc               flatSubscriptionGuard
	visitorSettingSvc             flatVisitorSettingDefaults
	flatAuthz                     *flatauthz.FlatVisitorAuthz
	memberInviteTTL               time.Duration
	notifier                      notificationsvc.NotificationService
	shortLinks                    shortlinksvc.ShortLinkService
	memberInviteUserCreator       memberInviteUserCreator
	memberInviteUserAuthenticator memberInviteUserAuthenticator
}

type memberInviteUserCreator interface {
	CreateResidentUser(ctx context.Context, req *models.ResidentRegisterRequest) (*models.UserResponse, error)
}

type memberInviteUserAuthenticator interface {
	AuthenticateCredentials(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error)
}

type flatVisitorSettingDefaults interface {
	CreateDefaultFlatSettings(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
}

type flatSubscriptionGuard interface {
	EnsureSocietyOperational(ctx context.Context, societyID int64) error
	CanAddFlatWithLock(ctx context.Context, societyID int64, adding int64) error
	CanAddResident(ctx context.Context, societyID int64, adding int64) error
	CanAddResidentWithLock(ctx context.Context, societyID int64, adding int64) error
}

func NewFlatService(
	flatRepo repository.FlatRepository,
	claimRepo repository.FlatClaimRepository,
	residentRepo repository.FlatResidentRepository,
	memberInviteRepo repository.FlatMemberInviteRepository,
	memberRepo repository.SocietyMemberRepository,
	txManager repository.TransactionManager,
	societySvc societysvc.SocietyService,
	subscriptionSvc flatSubscriptionGuard,
	memberInviteTTL time.Duration,
	deps ...any,
) FlatService {
	svc := &FlatSvc{
		flatRepo: flatRepo, claimRepo: claimRepo, residentRepo: residentRepo,
		memberInviteRepo: memberInviteRepo, memberRepo: memberRepo, txManager: txManager,
		societySvc: societySvc, subscriptionSvc: subscriptionSvc,
		memberInviteTTL: memberInviteTTL,
	}
	if svc.memberInviteTTL <= 0 {
		svc.memberInviteTTL = 168 * time.Hour
	}
	for _, dep := range deps {
		switch value := dep.(type) {
		case flatVisitorSettingDefaults:
			svc.visitorSettingSvc = value
		case *flatauthz.FlatVisitorAuthz:
			svc.flatAuthz = value
		case notificationsvc.NotificationService:
			svc.notifier = value
		case shortlinksvc.ShortLinkService:
			svc.shortLinks = value
		case memberInviteUserCreator:
			svc.memberInviteUserCreator = value
		case memberInviteUserAuthenticator:
			svc.memberInviteUserAuthenticator = value
		}
	}
	return svc
}

var _ FlatService = (*FlatSvc)(nil)
