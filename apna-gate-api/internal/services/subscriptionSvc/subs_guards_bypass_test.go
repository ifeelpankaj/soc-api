package subscriptionsvc_test

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"testing"

	"go-server/internal/models"
	"go-server/internal/requestctx"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"
)

type bypassSocietyRepo struct {
	getFn func(context.Context, models.GetSocietyFilter) (*models.Society, error)
}

func (r *bypassSocietyRepo) Create(context.Context, *models.Society) error { panic("unused") }
func (r *bypassSocietyRepo) Get(ctx context.Context, filter models.GetSocietyFilter) (*models.Society, error) {
	if r.getFn != nil {
		return r.getFn(ctx, filter)
	}
	return nil, errors.New("should not be called when bypass is active")
}
func (r *bypassSocietyRepo) List(context.Context, models.ListSocietiesFilter) ([]*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Count(context.Context, models.ListSocietiesFilter) (int64, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Update(context.Context, int64, contracts.UpdateSocietyInput) (*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Approve(context.Context, int64, int64) (*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Reject(context.Context, int64, int64, string) (*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Suspend(context.Context, int64, int64, string) (*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Reactivate(context.Context, int64, int64) (*models.Society, error) {
	panic("unused")
}
func (r *bypassSocietyRepo) Restore(context.Context, int64) (*models.Society, error) { panic("unused") }
func (r *bypassSocietyRepo) SoftDelete(context.Context, int64) error                 { panic("unused") }
func (r *bypassSocietyRepo) CountPendingByCreator(context.Context, int64) (int64, error) {
	panic("unused")
}

type bypassSubscriptionRepo struct{}

func (r *bypassSubscriptionRepo) CreatePending(context.Context, int64, int64, int64) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) CreateTrial(context.Context, int64, int64, int64, *contracts.CreateTrialSubscriptionInput) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) Get(context.Context, *models.SubscriptionFilter) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) List(context.Context, *models.SubscriptionFilter) ([]*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) Stats(context.Context, *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) GetActiveForUpdate(context.Context, int64) (*models.SocietySubscription, error) {
	return nil, errors.New("should not be called when bypass is active")
}
func (r *bypassSubscriptionRepo) Activate(context.Context, int64, int64, *contracts.ActivateSubscriptionInput) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) Renew(context.Context, int64, int64, *contracts.RenewSubscriptionInput) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) Cancel(context.Context, int64, int64, *contracts.CancelSubscriptionInput) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) Expire(context.Context, int64) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) ExpireDue(context.Context) (int64, error) { panic("unused") }
func (r *bypassSubscriptionRepo) ChangePlan(context.Context, int64, int64) (*models.SocietySubscription, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) CountActiveFlats(context.Context, int64) (int64, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) CountActiveAdmins(context.Context, int64) (int64, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) CountActiveStaff(context.Context, int64) (int64, error) {
	panic("unused")
}
func (r *bypassSubscriptionRepo) CountActiveResidents(context.Context, int64) (int64, error) {
	panic("unused")
}

func newBypassTestService() subscriptionsvc.SubscriptionService {
	return subscriptionsvc.NewSubscriptionService(&bypassSubscriptionRepo{}, &bypassSocietyRepo{})
}

func newEnforcingTestService() subscriptionsvc.SubscriptionService {
	return subscriptionsvc.NewSubscriptionService(
		&bypassSubscriptionRepo{},
		&bypassSocietyRepo{
			getFn: func(context.Context, models.GetSocietyFilter) (*models.Society, error) {
				return nil, nil
			},
		},
	)
}

func TestEnsureSocietyOperationalBypassSkipsChecks(t *testing.T) {
	svc := newBypassTestService()
	ctx := requestctx.WithDeveloperGuardBypass(context.Background())

	if err := svc.EnsureSocietyOperational(ctx, 42); err != nil {
		t.Fatalf("expected nil error with bypass, got %v", err)
	}
}

func TestCanAddFlatBypassSkipsQuota(t *testing.T) {
	svc := newBypassTestService()
	ctx := requestctx.WithDeveloperGuardBypass(context.Background())

	if err := svc.CanAddFlat(ctx, 42, 1); err != nil {
		t.Fatalf("expected nil error with bypass, got %v", err)
	}
}

func TestCanAddFlatWithLockBypassSkipsLockAndQuota(t *testing.T) {
	svc := newBypassTestService()
	ctx := requestctx.WithDeveloperGuardBypass(context.Background())

	if err := svc.CanAddFlatWithLock(ctx, 42, 1); err != nil {
		t.Fatalf("expected nil error with bypass, got %v", err)
	}
}

func TestCanAddResidentWithLockBypassSkipsLockAndQuota(t *testing.T) {
	svc := newBypassTestService()
	ctx := requestctx.WithDeveloperGuardBypass(context.Background())

	if err := svc.CanAddResidentWithLock(ctx, 42, 1); err != nil {
		t.Fatalf("expected nil error with bypass, got %v", err)
	}
}

func TestEnsureSocietyOperationalWithoutBypassStillEnforces(t *testing.T) {
	svc := newEnforcingTestService()

	err := svc.EnsureSocietyOperational(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error without bypass")
	}
	if !errors.Is(err, subscriptionsvc.ErrSubscriptionRequired) {
		t.Fatalf("expected ErrSubscriptionRequired, got %v", err)
	}
}
