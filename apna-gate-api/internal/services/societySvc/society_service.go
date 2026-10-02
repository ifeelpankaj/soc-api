package societysvc

import (
	"context"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
)

type SocietyService interface {
	CreateSocietyRequest(ctx context.Context, req models.CreateSocietyRequest, requestedBy int64) (*models.SocietyResponse, error)
	ApproveSociety(ctx context.Context, societyID int64, approvedBy int64) (*models.SocietyResponse, error)
	RejectSociety(ctx context.Context, societyID int64, rejectedBy int64, reason string) (*models.SocietyResponse, error)
	SuspendSociety(ctx context.Context, societyID int64, suspendedBy int64, reason string) (*models.SocietyResponse, error)
	ReactivateSociety(ctx context.Context, societyID int64, reactivatedBy int64) (*models.SocietyResponse, error)
	RestoreSociety(ctx context.Context, societyID int64, restoredBy int64) (*models.SocietyResponse, error)
	UpdateSociety(ctx context.Context, societyID int64, req models.UpdateSocietyRequest, updatedBy int64) (*models.SocietyResponse, error)
	DeleteSociety(ctx context.Context, societyID int64, deletedBy int64) error
	GetSociety(ctx context.Context, filter models.GetSocietyFilter) (*models.SocietyDetailResponse, error)
	GetPublicClaimOptions(ctx context.Context, societyCode string) (*models.PublicClaimOptionsResponse, error)
	ResolveActiveSocietyIDByCode(ctx context.Context, societyCode string) (int64, error)
	ListSocieties(ctx context.Context, filter models.ListSocietiesFilter) (*models.PaginatedSocietiesResponse, error)
	ListMySocieties(ctx context.Context, userID int64) ([]*models.MySocietyResponse, error)
	ListMyMemberships(ctx context.Context, userID int64) ([]*models.SocietyMemberResponse, error)
	GetOnboardingBootstrap(ctx context.Context, societyID int64) (*models.SocietyOnboardingBootstrapResponse, error)
	GetDashboardBootstrap(ctx context.Context, societyID int64) (*models.SocietyDashboardBootstrapResponse, error)
	GetDeveloperDashboardBootstrap(ctx context.Context) (*models.DeveloperDashboardBootstrapResponse, error)
	CreateGuard(ctx context.Context, societyID int64, actorID int64, req models.CreateGuardRequest) (*models.CreateGuardResponse, error)
	AddMember(ctx context.Context, req models.AddSocietyMemberRequest, addedBy int64) (*models.SocietyMemberResponse, error)
	ChangeMemberRole(ctx context.Context, req models.ChangeSocietyMemberRoleRequest, changedBy int64) (*models.SocietyMemberResponse, error)
	SuspendMember(ctx context.Context, req models.SuspendSocietyMemberRequest, suspendedBy int64) (*models.SocietyMemberResponse, error)
	ReactivateMember(ctx context.Context, req models.ReactivateSocietyMemberRequest, reactivatedBy int64) (*models.SocietyMemberResponse, error)
	RemoveMember(ctx context.Context, req models.RemoveSocietyMemberRequest, removedBy int64) error
	TransferOwnership(ctx context.Context, societyID int64, newOwnerUserID int64, changedBy int64) (*models.SocietyMemberResponse, error)
	GetSocietyMember(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMemberResponse, error)
	GetSocietyMemberDetail(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMemberDetailResponse, error)
	GetSocietyMemberSummary(ctx context.Context, societyID int64) (*models.SocietyMemberSummaryResponse, error)
	ListSocietyMembers(ctx context.Context, filter models.ListSocietyMembersFilter) (*models.PaginatedMembersResponse, error)
	ListAllSocietyMember(ctx context.Context, filter models.ListSocietyMembersFilter) (*models.PaginatedMembersResponse, error)
	EnsureActiveSociety(ctx context.Context, societyID int64) error
	EnsureActiveMember(ctx context.Context, societyID int64, userID int64) (*models.SocietyMemberResponse, error)
	EnsureMembershipRole(ctx context.Context, societyID int64, userID int64, roles ...string) error
	EnsureRole(ctx context.Context, societyID int64, userID int64, roles ...string) error
}

type SocietySvc struct {
	societyRepo       repository.SocietyRepository
	memberRepo        repository.SocietyMemberRepository
	flatRepo          repository.FlatRepository
	flatResidentRepo  repository.FlatResidentRepository
	flatClaimRepo     repository.FlatClaimRepository
	userRepo          repository.UserRepository
	txManager         repository.TransactionManager
	planSvc           societyPlanQuery
	subscriptionSvc   societySubscriptionDashboard
	visitorSettingSvc societyVisitorSettingDefaults
	visitorEntrySvc   societyVisitorDashboard
}

type societyVisitorSettingDefaults interface {
	CreateDefaultSocietySettings(ctx context.Context, societyID int64, actorUserID int64) error
}

type societyVisitorDashboard interface {
	GetEntryStats(ctx context.Context, societyID int64) (*models.VisitorEntryStatsResponse, error)
	GetDailyEntryStats(ctx context.Context, societyID int64, days int32) (*models.VisitorEntryDailyStatsResponse, error)
}

type societySubscriptionQuota interface {
	CanAddAdmin(ctx context.Context, societyID int64, adding int64) error
	CanAddStaff(ctx context.Context, societyID int64, adding int64) error
}

type societySubscriptionDashboard interface {
	societySubscriptionQuota
	ListSubscriptions(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscriptionResponse, error)
	GetSubscriptionStats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error)
}

type societyPlanQuery interface {
	ListPlans(ctx context.Context, filter *models.PlanFilter) ([]*models.PlanResponse, error)
	CountPlans(ctx context.Context, filter *models.PlanFilter) (int64, error)
}

func NewSocietyService(
	societyRepo repository.SocietyRepository,
	memberRepo repository.SocietyMemberRepository,
	flatRepo repository.FlatRepository,
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
	deps ...any,
) SocietyService {
	svc := &SocietySvc{societyRepo: societyRepo, memberRepo: memberRepo, flatRepo: flatRepo, userRepo: userRepo, txManager: txManager}
	for _, dep := range deps {
		switch value := dep.(type) {
		case repository.FlatClaimRepository:
			svc.flatClaimRepo = value
		case repository.FlatResidentRepository:
			svc.flatResidentRepo = value
		case societyPlanQuery:
			svc.planSvc = value
		case societySubscriptionDashboard:
			svc.subscriptionSvc = value
		case societyVisitorSettingDefaults:
			svc.visitorSettingSvc = value
		case societyVisitorDashboard:
			svc.visitorEntrySvc = value
		}
	}
	return svc
}

func NewSocietySvc(
	societyRepo repository.SocietyRepository,
	memberRepo repository.SocietyMemberRepository,
	flatRepo repository.FlatRepository,
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
) SocietyService {
	return NewSocietyService(societyRepo, memberRepo, flatRepo, userRepo, txManager)
}

var _ SocietyService = (*SocietySvc)(nil)
