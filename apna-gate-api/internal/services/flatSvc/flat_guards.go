package flatsvc

import (
	"context"

	"go-server/internal/models"
)

// Private guards.

// ensureFlatManager checks that the actor is an owner, admin, or staff member of the society.
func (s *FlatSvc) ensureFlatManager(ctx context.Context, societyID int64, userID int64) error {
	return s.societySvc.EnsureRole(ctx, societyID, userID, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin), string(models.SocietyMemberRoleStaff))
}

// ensureFlatOperational checks that the society is operational when subscription guards are configured.
func (s *FlatSvc) ensureFlatOperational(ctx context.Context, societyID int64) error {
	if s.subscriptionSvc == nil {
		return nil
	}
	return s.subscriptionSvc.EnsureSocietyOperational(ctx, societyID)
}

func (s *FlatSvc) ensureCanAddFlatsWithLock(ctx context.Context, societyID int64, adding int64) error {
	if s.subscriptionSvc == nil {
		return nil
	}
	return s.subscriptionSvc.CanAddFlatWithLock(ctx, societyID, adding)
}

// ensureFlatAssignable checks that the flat exists, is active, and is not blocked.
func (s *FlatSvc) ensureFlatAssignable(ctx context.Context, societyID int64, flatID int64) error {
	flat, err := s.flatRepo.Get(ctx, &models.FlatFilter{ID: &flatID, SocietyID: &societyID})
	if err != nil {
		return err
	}
	if flat == nil {
		return ErrFlatNotFound
	}
	if !flat.IsActive {
		return ErrFlatInactive
	}
	if flat.Status == models.FlatStatusBlocked {
		return ErrFlatBlocked
	}
	return nil
}

// requireResident checks that a resident exists for the supplied filter.
func (s *FlatSvc) requireResident(ctx context.Context, filter *models.FlatResidentFilter) (*models.FlatResident, error) {
	resident, err := s.residentRepo.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	if resident == nil {
		return nil, ErrResidentNotFound
	}
	return resident, nil
}
