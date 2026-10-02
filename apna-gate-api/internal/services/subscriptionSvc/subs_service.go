package subscriptionsvc

import (
	"context"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
)

type SubscriptionService interface {
	SubscriptionCommandService
	SubscriptionQueryService
	SubscriptionGuardService
	SubscriptionQuotaService
}

type SubscriptionCommandService interface {
	CreatePendingSubscription(ctx context.Context, societyID, planID, createdBy int64) (*models.SocietySubscriptionResponse, error)
	CreateTrialSubscription(ctx context.Context, societyID, planID, createdBy int64, req *models.CreateTrialSubscriptionRequest) (*models.SocietySubscriptionResponse, error)
	ActivateSubscription(ctx context.Context, subscriptionID, activatedBy int64, req *models.ActivateSubscriptionRequest) (*models.SocietySubscriptionResponse, error)
	RenewSubscription(ctx context.Context, subscriptionID, renewedBy int64, req *models.RenewSubscriptionRequest) (*models.SocietySubscriptionResponse, error)
	CancelSubscription(ctx context.Context, subscriptionID, cancelledBy int64, req *models.CancelSubscriptionRequest) (*models.SocietySubscriptionResponse, error)
	ExpireSubscription(ctx context.Context, subscriptionID int64) (*models.SocietySubscriptionResponse, error)
	ExpireDueSubscriptions(ctx context.Context) (int64, error)
	ChangeSubscriptionPlan(ctx context.Context, subscriptionID, newPlanID, changedBy int64) (*models.SocietySubscriptionResponse, error)
}

type SubscriptionQueryService interface {
	GetSubscription(ctx context.Context, filter *models.SubscriptionFilter) (*models.SocietySubscriptionResponse, error)
	GetActiveSubscriptionBySocietyID(ctx context.Context, societyID int64) (*models.SocietySubscriptionResponse, error)
	ListSubscriptions(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscriptionResponse, error)
	GetSubscriptionStats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error)
}

type SubscriptionGuardService interface {
	EnsureSocietyOperational(ctx context.Context, societyID int64) error
	EnsureActiveSubscription(ctx context.Context, societyID int64) error
	EnsureFeatureEnabled(ctx context.Context, societyID int64, feature string) error
}

type SubscriptionQuotaService interface {
	CanAddFlat(ctx context.Context, societyID int64, adding int64) error
	CanAddFlatWithLock(ctx context.Context, societyID int64, adding int64) error
	CanAddAdmin(ctx context.Context, societyID int64, adding int64) error
	CanAddStaff(ctx context.Context, societyID int64, adding int64) error
	CanAddResident(ctx context.Context, societyID int64, adding int64) error
	CanAddResidentWithLock(ctx context.Context, societyID int64, adding int64) error
}

type SubscriptionSvc struct {
	subRepo     repository.SubscriptionRepository
	societyRepo repository.SocietyRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, societyRepo repository.SocietyRepository) SubscriptionService {
	return &SubscriptionSvc{subRepo: subRepo, societyRepo: societyRepo}
}

var _ SubscriptionService = (*SubscriptionSvc)(nil)
