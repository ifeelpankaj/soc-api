package subscriptionsvc_test

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/requestctx"
	subscriptionsvc "go-server/internal/services/subscriptionSvc"
)

type lifecycleSubRepo struct {
	sub               *models.SocietySubscription
	latest            []*models.SocietySubscription
	getErr            error
	listErr           error
	countErr          error
	activeForUpdate   *models.SocietySubscription
	activeForUpdateCt int
	counts            map[string]int64
	nilOnRenew        bool
	nilOnExpire       bool
}

func (r *lifecycleSubRepo) CreatePending(context.Context, int64, int64, int64) (*models.SocietySubscription, error) {
	return nil, nil
}

func (r *lifecycleSubRepo) CreateTrial(_ context.Context, societyID, planID, createdBy int64, req *contracts.CreateTrialSubscriptionInput) (*models.SocietySubscription, error) {
	planIDValue := planID
	sub := activeSubscription(societyID)
	sub.ID = 88
	sub.PlanID = &planIDValue
	sub.Status = models.SubscriptionStatusTrial
	sub.StartsAt = &req.StartsAt
	sub.TrialEndsAt = &req.TrialEndsAt
	sub.EndsAt = req.EndsAt
	sub.CreatedBy = &createdBy
	r.sub = sub
	r.latest = []*models.SocietySubscription{sub}
	return sub, nil
}

func (r *lifecycleSubRepo) Activate(_ context.Context, subscriptionID, activatedBy int64, req *contracts.ActivateSubscriptionInput) (*models.SocietySubscription, error) {
	sub := r.ensureSub(subscriptionID, models.SubscriptionStatusActive)
	sub.StartsAt = &req.StartsAt
	sub.EndsAt = &req.EndsAt
	sub.ActivatedBy = &activatedBy
	now := time.Now().UTC()
	sub.ActivatedAt = &now
	return sub, nil
}

func (r *lifecycleSubRepo) Renew(_ context.Context, subscriptionID, renewedBy int64, req *contracts.RenewSubscriptionInput) (*models.SocietySubscription, error) {
	if r.nilOnRenew {
		return nil, nil
	}
	sub := r.ensureSub(subscriptionID, models.SubscriptionStatusActive)
	sub.StartsAt = &req.StartsAt
	sub.EndsAt = &req.EndsAt
	sub.ActivatedBy = &renewedBy
	sub.ExpiredAt = nil
	sub.CancelledAt = nil
	r.latest = []*models.SocietySubscription{sub}
	return sub, nil
}

func (r *lifecycleSubRepo) Cancel(_ context.Context, subscriptionID, cancelledBy int64, req *contracts.CancelSubscriptionInput) (*models.SocietySubscription, error) {
	sub := r.ensureSub(subscriptionID, models.SubscriptionStatusCancelled)
	now := time.Now().UTC()
	sub.CancelledAt = &now
	sub.CancelledBy = &cancelledBy
	sub.CancellationReason = &req.Reason
	return sub, nil
}

func (r *lifecycleSubRepo) Expire(_ context.Context, subscriptionID int64) (*models.SocietySubscription, error) {
	if r.nilOnExpire {
		return nil, nil
	}
	sub := r.ensureSub(subscriptionID, models.SubscriptionStatusExpired)
	now := time.Now().UTC()
	sub.ExpiredAt = &now
	past := now.Add(-time.Minute)
	sub.EndsAt = &past
	r.latest = []*models.SocietySubscription{sub}
	return sub, nil
}

func (r *lifecycleSubRepo) ExpireDue(context.Context) (int64, error) { return 0, nil }

func (r *lifecycleSubRepo) ChangePlan(context.Context, int64, int64) (*models.SocietySubscription, error) {
	return nil, nil
}

func (r *lifecycleSubRepo) Get(_ context.Context, filter *models.SubscriptionFilter) (*models.SocietySubscription, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	if filter != nil && filter.IsActiveOnly != nil && *filter.IsActiveOnly {
		if isActiveNow(r.sub) {
			return r.sub, nil
		}
		return nil, nil
	}
	return r.sub, nil
}

func (r *lifecycleSubRepo) List(context.Context, *models.SubscriptionFilter) ([]*models.SocietySubscription, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	if r.latest != nil {
		return r.latest, nil
	}
	if r.sub == nil {
		return nil, nil
	}
	return []*models.SocietySubscription{r.sub}, nil
}

func (r *lifecycleSubRepo) Stats(context.Context, *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	return nil, nil
}

func (r *lifecycleSubRepo) CountActiveFlats(context.Context, int64) (int64, error) {
	return r.count("flats")
}

func (r *lifecycleSubRepo) CountActiveAdmins(context.Context, int64) (int64, error) {
	return r.count("admins")
}

func (r *lifecycleSubRepo) CountActiveStaff(context.Context, int64) (int64, error) {
	return r.count("staff")
}

func (r *lifecycleSubRepo) CountActiveResidents(context.Context, int64) (int64, error) {
	return r.count("residents")
}

func (r *lifecycleSubRepo) GetActiveForUpdate(context.Context, int64) (*models.SocietySubscription, error) {
	r.activeForUpdateCt++
	if r.activeForUpdate != nil {
		return r.activeForUpdate, nil
	}
	if isActiveNow(r.sub) {
		return r.sub, nil
	}
	return nil, nil
}

func (r *lifecycleSubRepo) ensureSub(id int64, status models.SubscriptionStatus) *models.SocietySubscription {
	if r.sub == nil {
		r.sub = activeSubscription(1)
	}
	r.sub.ID = id
	r.sub.Status = status
	return r.sub
}

func (r *lifecycleSubRepo) count(name string) (int64, error) {
	if r.countErr != nil {
		return 0, r.countErr
	}
	if r.counts == nil {
		return 0, nil
	}
	return r.counts[name], nil
}

type lifecycleSocietyRepo struct {
	society *models.Society
	err     error
	calls   int
}

func (r *lifecycleSocietyRepo) Create(context.Context, *models.Society) error { return nil }
func (r *lifecycleSocietyRepo) Get(context.Context, models.GetSocietyFilter) (*models.Society, error) {
	r.calls++
	return r.society, r.err
}
func (r *lifecycleSocietyRepo) List(context.Context, models.ListSocietiesFilter) ([]*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Count(context.Context, models.ListSocietiesFilter) (int64, error) {
	return 0, nil
}
func (r *lifecycleSocietyRepo) Update(context.Context, int64, contracts.UpdateSocietyInput) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Approve(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Reject(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Suspend(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Reactivate(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) Restore(context.Context, int64) (*models.Society, error) {
	return nil, nil
}
func (r *lifecycleSocietyRepo) SoftDelete(context.Context, int64) error { return nil }
func (r *lifecycleSocietyRepo) CountPendingByCreator(context.Context, int64) (int64, error) {
	return 0, nil
}

func TestSubscriptionExpiryDisablesProtectedOperationsAndRenewalRestoresThem(t *testing.T) {
	societyID := int64(10)
	subRepo := &lifecycleSubRepo{
		sub:    activeSubscription(societyID),
		counts: map[string]int64{"flats": 1, "admins": 1, "staff": 1, "residents": 1},
	}
	svc := subscriptionsvc.NewSubscriptionService(subRepo, &lifecycleSocietyRepo{society: society(societyID, models.SocietyStatusActive)})

	assertNoErr(t, svc.EnsureSocietyOperational(context.Background(), societyID))
	assertNoErr(t, svc.EnsureActiveSubscription(context.Background(), societyID))
	assertNoErr(t, svc.EnsureFeatureEnabled(context.Background(), societyID, "visitor_management"))
	assertNoErr(t, svc.CanAddFlat(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddAdmin(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddStaff(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddResident(context.Background(), societyID, 1))

	if _, err := svc.ExpireSubscription(context.Background(), subRepo.sub.ID); err != nil {
		t.Fatalf("ExpireSubscription returned error: %v", err)
	}

	assertAppCode(t, svc.EnsureSocietyOperational(context.Background(), societyID), subscriptionsvc.ErrSubscriptionExpired.Code)
	assertAppCode(t, svc.EnsureActiveSubscription(context.Background(), societyID), subscriptionsvc.ErrSubscriptionExpired.Code)
	assertAppCode(t, svc.EnsureFeatureEnabled(context.Background(), societyID, "visitor_management"), subscriptionsvc.ErrSubscriptionRequired.Code)
	assertAppCode(t, svc.CanAddFlat(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)
	assertAppCode(t, svc.CanAddAdmin(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)
	assertAppCode(t, svc.CanAddStaff(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)
	assertAppCode(t, svc.CanAddResident(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)

	start := time.Now().UTC()
	end := start.Add(30 * 24 * time.Hour)
	if _, err := svc.RenewSubscription(context.Background(), subRepo.sub.ID, 99, &models.RenewSubscriptionRequest{StartsAt: start, EndsAt: end}); err != nil {
		t.Fatalf("RenewSubscription returned error: %v", err)
	}

	assertNoErr(t, svc.EnsureSocietyOperational(context.Background(), societyID))
	assertNoErr(t, svc.EnsureActiveSubscription(context.Background(), societyID))
	assertNoErr(t, svc.EnsureFeatureEnabled(context.Background(), societyID, "visitor_management"))
	assertNoErr(t, svc.CanAddFlat(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddAdmin(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddStaff(context.Background(), societyID, 1))
	assertNoErr(t, svc.CanAddResident(context.Background(), societyID, 1))
}

func TestEnsureSocietyOperational(t *testing.T) {
	societyID := int64(10)
	tests := []struct {
		name    string
		society *models.Society
		sub     *models.SocietySubscription
		latest  []*models.SocietySubscription
		want    string
	}{
		{name: "active society with active subscription", society: society(societyID, models.SocietyStatusActive), sub: activeSubscription(societyID)},
		{name: "missing society", society: nil, sub: activeSubscription(societyID), want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "pending society", society: society(societyID, models.SocietyStatusPending), sub: activeSubscription(societyID), want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "suspended society", society: society(societyID, models.SocietyStatusSuspended), sub: activeSubscription(societyID), want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "rejected society", society: society(societyID, models.SocietyStatusRejected), sub: activeSubscription(societyID), want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "no subscription", society: society(societyID, models.SocietyStatusActive), want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "expired subscription", society: society(societyID, models.SocietyStatusActive), sub: expiredSubscription(societyID), want: subscriptionsvc.ErrSubscriptionExpired.Code},
		{name: "cancelled subscription", society: society(societyID, models.SocietyStatusActive), sub: cancelledSubscription(societyID), want: subscriptionsvc.ErrSubscriptionExpired.Code},
		{name: "ended latest subscription", society: society(societyID, models.SocietyStatusActive), latest: []*models.SocietySubscription{endedSubscription(societyID)}, want: subscriptionsvc.ErrSubscriptionExpired.Code},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := subscriptionsvc.NewSubscriptionService(
				&lifecycleSubRepo{sub: tt.sub, latest: tt.latest},
				&lifecycleSocietyRepo{society: tt.society},
			)
			err := svc.EnsureSocietyOperational(context.Background(), societyID)
			if tt.want == "" {
				assertNoErr(t, err)
				return
			}
			assertAppCode(t, err, tt.want)
		})
	}
}

func TestEnsureActiveSubscription(t *testing.T) {
	societyID := int64(10)
	repoErr := errors.New("database unavailable")
	tests := []struct {
		name   string
		repo   *lifecycleSubRepo
		want   string
		bypass bool
	}{
		{name: "active subscription", repo: &lifecycleSubRepo{sub: activeSubscription(societyID)}},
		{name: "no subscription", repo: &lifecycleSubRepo{}, want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "latest expired", repo: &lifecycleSubRepo{sub: expiredSubscription(societyID)}, want: subscriptionsvc.ErrSubscriptionExpired.Code},
		{name: "repository error is wrapped", repo: &lifecycleSubRepo{getErr: repoErr}, want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "developer bypass skips repo", repo: &lifecycleSubRepo{getErr: repoErr}, bypass: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := subscriptionsvc.NewSubscriptionService(tt.repo, &lifecycleSocietyRepo{})
			ctx := context.Background()
			if tt.bypass {
				ctx = requestctx.WithDeveloperGuardBypass(ctx)
			}
			err := svc.EnsureActiveSubscription(ctx, societyID)
			if tt.want == "" {
				assertNoErr(t, err)
				return
			}
			assertAppCode(t, err, tt.want)
		})
	}
}

func TestEnsureFeatureEnabled(t *testing.T) {
	societyID := int64(10)
	tests := []struct {
		name     string
		features map[string]any
		sub      *models.SocietySubscription
		feature  string
		want     string
	}{
		{name: "feature true", features: map[string]any{"visitor_management": true}, feature: "visitor_management"},
		{name: "feature false", features: map[string]any{"visitor_management": false}, feature: "visitor_management", want: subscriptionsvc.ErrFeatureDisabled.Code},
		{name: "feature missing", features: map[string]any{}, feature: "visitor_management", want: subscriptionsvc.ErrFeatureDisabled.Code},
		{name: "feature non boolean", features: map[string]any{"visitor_management": "yes"}, feature: "visitor_management", want: subscriptionsvc.ErrFeatureDisabled.Code},
		{name: "expired subscription", sub: expiredSubscription(societyID), feature: "visitor_management", want: subscriptionsvc.ErrSubscriptionRequired.Code},
		{name: "no subscription", sub: nil, feature: "visitor_management", want: subscriptionsvc.ErrSubscriptionRequired.Code},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.sub
			if sub == nil && tt.features != nil {
				sub = activeSubscription(societyID)
				sub.Features = tt.features
			}
			svc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{sub: sub}, &lifecycleSocietyRepo{})
			err := svc.EnsureFeatureEnabled(context.Background(), societyID, tt.feature)
			if tt.want == "" {
				assertNoErr(t, err)
				return
			}
			assertAppCode(t, err, tt.want)
		})
	}

	svc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{getErr: errors.New("should not be called")}, &lifecycleSocietyRepo{})
	assertNoErr(t, svc.EnsureFeatureEnabled(requestctx.WithDeveloperGuardBypass(context.Background()), societyID, "missing"))
}

func TestQuotaGuards(t *testing.T) {
	societyID := int64(10)
	sub := activeSubscription(societyID)
	sub.MaxFlats = 2
	sub.MaxAdmins = 2
	sub.MaxStaff = 2
	sub.MaxResidents = 2

	t.Run("adding zero or negative succeeds", func(t *testing.T) {
		svc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{sub: sub}, &lifecycleSocietyRepo{})
		assertNoErr(t, svc.CanAddFlat(context.Background(), societyID, 0))
		assertNoErr(t, svc.CanAddAdmin(context.Background(), societyID, -1))
		assertNoErr(t, svc.CanAddStaff(context.Background(), societyID, 0))
		assertNoErr(t, svc.CanAddResident(context.Background(), societyID, -1))
	})

	t.Run("within quota succeeds", func(t *testing.T) {
		repo := &lifecycleSubRepo{sub: sub, counts: map[string]int64{"flats": 1, "admins": 1, "staff": 1, "residents": 1}}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertNoErr(t, svc.CanAddFlat(context.Background(), societyID, 1))
		assertNoErr(t, svc.CanAddAdmin(context.Background(), societyID, 1))
		assertNoErr(t, svc.CanAddStaff(context.Background(), societyID, 1))
		assertNoErr(t, svc.CanAddResident(context.Background(), societyID, 1))
	})

	t.Run("over quota fails", func(t *testing.T) {
		repo := &lifecycleSubRepo{sub: sub, counts: map[string]int64{"flats": 2, "admins": 2, "staff": 2, "residents": 2}}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertAppCode(t, svc.CanAddFlat(context.Background(), societyID, 1), subscriptionsvc.ErrQuotaExceeded.Code)
		assertAppCode(t, svc.CanAddAdmin(context.Background(), societyID, 1), subscriptionsvc.ErrQuotaExceeded.Code)
		assertAppCode(t, svc.CanAddStaff(context.Background(), societyID, 1), subscriptionsvc.ErrQuotaExceeded.Code)
		assertAppCode(t, svc.CanAddResident(context.Background(), societyID, 1), subscriptionsvc.ErrQuotaExceeded.Code)
	})

	t.Run("count error is returned", func(t *testing.T) {
		countErr := errors.New("count failed")
		svc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{sub: sub, countErr: countErr}, &lifecycleSocietyRepo{})
		if err := svc.CanAddFlat(context.Background(), societyID, 1); !errors.Is(err, countErr) {
			t.Fatalf("expected count error, got %v", err)
		}
	})

	t.Run("resident lock is taken before quota check", func(t *testing.T) {
		repo := &lifecycleSubRepo{sub: sub, counts: map[string]int64{"residents": 1}}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertNoErr(t, svc.CanAddResidentWithLock(context.Background(), societyID, 1))
		if repo.activeForUpdateCt != 1 {
			t.Fatalf("expected GetActiveForUpdate once, got %d", repo.activeForUpdateCt)
		}
	})

	t.Run("flat lock is taken before quota check", func(t *testing.T) {
		repo := &lifecycleSubRepo{sub: sub, counts: map[string]int64{"flats": 1}}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertNoErr(t, svc.CanAddFlatWithLock(context.Background(), societyID, 1))
		if repo.activeForUpdateCt != 1 {
			t.Fatalf("expected GetActiveForUpdate once, got %d", repo.activeForUpdateCt)
		}
	})

	t.Run("flat lock nil fails", func(t *testing.T) {
		repo := &lifecycleSubRepo{}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertAppCode(t, svc.CanAddFlatWithLock(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)
	})

	t.Run("resident lock nil fails", func(t *testing.T) {
		repo := &lifecycleSubRepo{}
		svc := subscriptionsvc.NewSubscriptionService(repo, &lifecycleSocietyRepo{})
		assertAppCode(t, svc.CanAddResidentWithLock(context.Background(), societyID, 1), subscriptionsvc.ErrSubscriptionRequired.Code)
	})
}

func TestSubscriptionCommandValidation(t *testing.T) {
	societyID := int64(10)
	now := time.Now().UTC()
	future := now.Add(24 * time.Hour)
	past := now.Add(-time.Hour)

	svc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{sub: activeSubscription(societyID)}, &lifecycleSocietyRepo{})

	assertAppCode(t, mustSecond(svc.CreateTrialSubscription(context.Background(), societyID, 1, 1, nil)), subscriptionsvc.ErrInvalidSubscription.Code)
	assertAppCode(t, mustSecond(svc.CreateTrialSubscription(context.Background(), societyID, 1, 1, &models.CreateTrialSubscriptionRequest{StartsAt: now, TrialEndsAt: past})), subscriptionsvc.ErrInvalidSubscription.Code)
	assertAppCode(t, mustSecond(svc.CreateTrialSubscription(context.Background(), societyID, 1, 1, &models.CreateTrialSubscriptionRequest{StartsAt: now, TrialEndsAt: future, EndsAt: &now})), subscriptionsvc.ErrInvalidSubscription.Code)

	assertAppCode(t, mustSecond(svc.ActivateSubscription(context.Background(), 1, 1, nil)), subscriptionsvc.ErrInvalidSubscription.Code)
	assertAppCode(t, mustSecond(svc.ActivateSubscription(context.Background(), 1, 1, &models.ActivateSubscriptionRequest{StartsAt: now, EndsAt: now})), subscriptionsvc.ErrInvalidSubscription.Code)

	assertAppCode(t, mustSecond(svc.RenewSubscription(context.Background(), 1, 1, nil)), subscriptionsvc.ErrInvalidSubscription.Code)
	assertAppCode(t, mustSecond(svc.RenewSubscription(context.Background(), 1, 1, &models.RenewSubscriptionRequest{StartsAt: now, EndsAt: now})), subscriptionsvc.ErrInvalidSubscription.Code)

	assertAppCode(t, mustSecond(svc.CancelSubscription(context.Background(), 1, 1, nil)), subscriptionsvc.ErrInvalidSubscription.Code)
	assertAppCode(t, mustSecond(svc.CancelSubscription(context.Background(), 1, 1, &models.CancelSubscriptionRequest{Reason: "   "})), subscriptionsvc.ErrInvalidSubscription.Code)

	expireNilSvc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{nilOnExpire: true}, &lifecycleSocietyRepo{})
	assertAppCode(t, mustSecond(expireNilSvc.ExpireSubscription(context.Background(), 404)), subscriptionsvc.ErrSubscriptionNotFound.Code)

	renewNilSvc := subscriptionsvc.NewSubscriptionService(&lifecycleSubRepo{nilOnRenew: true}, &lifecycleSocietyRepo{})
	assertAppCode(t, mustSecond(renewNilSvc.RenewSubscription(context.Background(), 404, 1, &models.RenewSubscriptionRequest{StartsAt: now, EndsAt: future})), subscriptionsvc.ErrSubscriptionNotFound.Code)
}

func society(id int64, status models.SocietyStatus) *models.Society {
	return &models.Society{ID: id, Name: "Test Society", SocietyCode: "TEST", Status: status}
}

func activeSubscription(societyID int64) *models.SocietySubscription {
	now := time.Now().UTC()
	end := now.Add(30 * 24 * time.Hour)
	planID := int64(1)
	return &models.SocietySubscription{
		ID:               77,
		SocietyID:        societyID,
		PlanID:           &planID,
		Status:           models.SubscriptionStatusActive,
		StartsAt:         &now,
		EndsAt:           &end,
		PlanName:         "Premium",
		PlanCode:         "PREMIUM",
		Currency:         "INR",
		BillingCycle:     models.BillingCycleMonthly,
		MaxFlats:         2,
		MaxAdmins:        2,
		MaxStaff:         2,
		MaxResidents:     2,
		Features:         map[string]any{"visitor_management": true, "resident_management": true},
		Metadata:         map[string]any{},
		PriceAmountPaise: 99900,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func expiredSubscription(societyID int64) *models.SocietySubscription {
	sub := activeSubscription(societyID)
	now := time.Now().UTC()
	expiredAt := now.Add(-time.Hour)
	sub.Status = models.SubscriptionStatusExpired
	sub.EndsAt = &expiredAt
	sub.ExpiredAt = &expiredAt
	return sub
}

func cancelledSubscription(societyID int64) *models.SocietySubscription {
	sub := activeSubscription(societyID)
	now := time.Now().UTC()
	sub.Status = models.SubscriptionStatusCancelled
	sub.CancelledAt = &now
	return sub
}

func endedSubscription(societyID int64) *models.SocietySubscription {
	sub := activeSubscription(societyID)
	endedAt := time.Now().UTC().Add(-time.Minute)
	sub.EndsAt = &endedAt
	return sub
}

func isActiveNow(sub *models.SocietySubscription) bool {
	return sub != nil &&
		(sub.Status == models.SubscriptionStatusActive || sub.Status == models.SubscriptionStatusTrial) &&
		(sub.EndsAt == nil || sub.EndsAt.After(time.Now().UTC()))
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func assertAppCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected app error code %s, got nil", code)
	}
	var appErr *models.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error code %s, got %T: %v", code, err, err)
	}
	if appErr.Code != code {
		t.Fatalf("expected app error code %s, got %s (%v)", code, appErr.Code, err)
	}
}

func mustSecond[T any](value T, err error) error {
	return err
}

var _ repository.SubscriptionRepository = (*lifecycleSubRepo)(nil)
var _ repository.SocietyRepository = (*lifecycleSocietyRepo)(nil)
