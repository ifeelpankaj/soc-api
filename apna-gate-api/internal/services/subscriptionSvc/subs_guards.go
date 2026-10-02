package subscriptionsvc

import (
	"context"
	"errors"
	"fmt"

	"go-server/internal/models"
	"go-server/internal/requestctx"
)

// Public guards.

func skipSubscriptionGuards(ctx context.Context) bool {
	return requestctx.HasDeveloperGuardBypass(ctx)
}

// EnsureSocietyOperational checks that the society is active and has an active subscription.
func (s *SubscriptionSvc) EnsureSocietyOperational(ctx context.Context, societyID int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return err
	}
	if society == nil || society.Status != models.SocietyStatusActive {
		return ErrSubscriptionRequired
	}
	return s.EnsureActiveSubscription(ctx, societyID)
}

// EnsureActiveSubscription checks that the society has an active subscription.
func (s *SubscriptionSvc) EnsureActiveSubscription(ctx context.Context, societyID int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	_, err := s.GetActiveSubscriptionBySocietyID(ctx, societyID)
	if err != nil {
		if errors.Is(err, ErrSubscriptionExpired) {
			return err
		}
		return ErrSubscriptionRequired.WithCause(err)
	}
	return nil
}

// EnsureFeatureEnabled checks that a named subscription feature exists and is enabled.
func (s *SubscriptionSvc) EnsureFeatureEnabled(ctx context.Context, societyID int64, feature string) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	sub, err := s.GetActiveSubscriptionBySocietyID(ctx, societyID)
	if err != nil {
		return ErrSubscriptionRequired.WithCause(err)
	}
	enabled, ok := sub.Features[feature].(bool)
	if !ok || !enabled {
		return ErrFeatureDisabled
	}
	return nil
}

// CanAddFlat checks that adding flats will not exceed the active subscription limit.
func (s *SubscriptionSvc) CanAddFlat(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	return s.checkQuota(ctx, societyID, adding, "flats", func(sub *models.SocietySubscriptionResponse) int64 {
		return int64(sub.MaxFlats)
	}, s.subRepo.CountActiveFlats)
}

// CanAddFlatWithLock serializes flat quota checks by locking the active subscription row in the current transaction.
func (s *SubscriptionSvc) CanAddFlatWithLock(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	if adding <= 0 {
		return nil
	}
	sub, err := s.subRepo.GetActiveForUpdate(ctx, societyID)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionRequired
	}
	return s.CanAddFlat(ctx, societyID, adding)
}

// CanAddAdmin checks that adding admins will not exceed the active subscription limit.
func (s *SubscriptionSvc) CanAddAdmin(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	return s.checkQuota(ctx, societyID, adding, "admins", func(sub *models.SocietySubscriptionResponse) int64 {
		return int64(sub.MaxAdmins)
	}, s.subRepo.CountActiveAdmins)
}

// CanAddStaff checks that adding staff will not exceed the active subscription limit.
func (s *SubscriptionSvc) CanAddStaff(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	return s.checkQuota(ctx, societyID, adding, "staff", func(sub *models.SocietySubscriptionResponse) int64 {
		return int64(sub.MaxStaff)
	}, s.subRepo.CountActiveStaff)
}

// CanAddResident checks that adding residents will not exceed the active subscription limit.
func (s *SubscriptionSvc) CanAddResident(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	return s.checkQuota(ctx, societyID, adding, "residents", func(sub *models.SocietySubscriptionResponse) int64 {
		return int64(sub.MaxResidents)
	}, s.subRepo.CountActiveResidents)
}

// CanAddResidentWithLock serializes quota checks by locking the active subscription row in the current transaction.
func (s *SubscriptionSvc) CanAddResidentWithLock(ctx context.Context, societyID int64, adding int64) error {
	if skipSubscriptionGuards(ctx) {
		return nil
	}
	if adding <= 0 {
		return nil
	}
	sub, err := s.subRepo.GetActiveForUpdate(ctx, societyID)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionRequired
	}
	return s.CanAddResident(ctx, societyID, adding)
}

// Private guards.

// checkQuota checks current usage plus the requested addition against a subscription limit.
func (s *SubscriptionSvc) checkQuota(ctx context.Context, societyID int64, adding int64, name string, limitFn func(*models.SocietySubscriptionResponse) int64, countFn func(context.Context, int64) (int64, error)) error {
	if adding <= 0 {
		return nil
	}
	sub, err := s.GetActiveSubscriptionBySocietyID(ctx, societyID)
	if err != nil {
		return ErrSubscriptionRequired.WithCause(err)
	}
	current, err := countFn(ctx, societyID)
	if err != nil {
		return err
	}
	limit := limitFn(sub)
	if current+adding > limit {
		return ErrQuotaExceeded.WithCause(fmt.Errorf("%s limit %d exceeded: current %d adding %d", name, limit, current, adding))
	}
	return nil
}
