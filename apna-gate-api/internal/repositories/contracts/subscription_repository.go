package contracts

import (
	"context"
	"go-server/internal/models"
)

type SubscriptionRepository interface {
	CreatePending(ctx context.Context, societyID, planID, createdBy int64) (*models.SocietySubscription, error)
	CreateTrial(ctx context.Context, societyID, planID, createdBy int64, req *CreateTrialSubscriptionInput) (*models.SocietySubscription, error)
	Activate(ctx context.Context, subscriptionID, activatedBy int64, req *ActivateSubscriptionInput) (*models.SocietySubscription, error)
	Renew(ctx context.Context, subscriptionID, renewedBy int64, req *RenewSubscriptionInput) (*models.SocietySubscription, error)
	Cancel(ctx context.Context, subscriptionID, cancelledBy int64, req *CancelSubscriptionInput) (*models.SocietySubscription, error)
	Expire(ctx context.Context, subscriptionID int64) (*models.SocietySubscription, error)
	ExpireDue(ctx context.Context) (int64, error)
	ChangePlan(ctx context.Context, subscriptionID, newPlanID int64) (*models.SocietySubscription, error)
	Get(ctx context.Context, filter *models.SubscriptionFilter) (*models.SocietySubscription, error)
	List(ctx context.Context, filter *models.SubscriptionFilter) ([]*models.SocietySubscription, error)
	Stats(ctx context.Context, filter *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error)
	CountActiveFlats(ctx context.Context, societyID int64) (int64, error)
	CountActiveAdmins(ctx context.Context, societyID int64) (int64, error)
	CountActiveStaff(ctx context.Context, societyID int64) (int64, error)
	CountActiveResidents(ctx context.Context, societyID int64) (int64, error)
	GetActiveForUpdate(ctx context.Context, societyID int64) (*models.SocietySubscription, error)
}
