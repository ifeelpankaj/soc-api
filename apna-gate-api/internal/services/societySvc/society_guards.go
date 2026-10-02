package societysvc

import (
	"context"

	"go-server/internal/models"
)

// Public guards.

// EnsureActiveSociety checks that the society exists and is active.
func (s *SocietySvc) EnsureActiveSociety(ctx context.Context, societyID int64) error {
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return err
	}
	if society == nil {
		return ErrSocietyNotFound
	}
	if society.Status != models.SocietyStatusActive {
		return ErrSocietyInactive
	}
	return nil
}

// EnsureActiveMember checks that the user is an active member of the society.
func (s *SocietySvc) EnsureActiveMember(ctx context.Context, societyID int64, userID int64) (*models.SocietyMemberResponse, error) {
	member, err := s.requireMember(ctx, societyID, userID)
	if err != nil {
		return nil, err
	}
	if member.Status != models.SocietyMemberStatusActive {
		return nil, ErrMemberInactive
	}
	return member.ToResponse(), nil
}

// EnsureRole checks that the society is active and the user has one of the allowed roles.
func (s *SocietySvc) EnsureRole(ctx context.Context, societyID int64, userID int64, roles ...string) error {
	if err := s.EnsureActiveSociety(ctx, societyID); err != nil {
		return err
	}
	return s.ensureMembershipRole(ctx, societyID, userID, roles...)
}

// EnsureMembershipRole checks that the society exists and the user has an active
// membership with one of the allowed roles, regardless of society lifecycle.
func (s *SocietySvc) EnsureMembershipRole(ctx context.Context, societyID int64, userID int64, roles ...string) error {
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return err
	}
	if society == nil {
		return ErrSocietyNotFound
	}
	return s.ensureMembershipRole(ctx, societyID, userID, roles...)
}

func (s *SocietySvc) ensureMembershipRole(ctx context.Context, societyID int64, userID int64, roles ...string) error {
	member, err := s.requireMember(ctx, societyID, userID)
	if err != nil {
		return err
	}
	if member.Status != models.SocietyMemberStatusActive {
		return ErrMemberInactive
	}

	for _, role := range roles {
		if string(member.Role) == role {
			return nil
		}
	}
	return ErrForbiddenSociety
}

// Private guards.

// requireMember checks that a society member exists for the society and user.
func (s *SocietySvc) requireMember(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error) {
	member, err := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &societyID, UserID: &userID})
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrMemberNotFound
	}
	return member, nil
}

// ensureUserExists checks that the user exists before creating or changing society records.
func (s *SocietySvc) ensureUserExists(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return models.NewAppError(models.ErrCodeNotFound, "user not found", 404, nil)
	}
	return nil
}

// ensureLifecycleStatus checks that the society exists and is currently in the expected status.
func (s *SocietySvc) ensureLifecycleStatus(ctx context.Context, societyID int64, expected models.SocietyStatus) error {
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return err
	}
	if society == nil {
		return ErrSocietyNotFound
	}
	if society.Status != expected {
		return ErrInvalidTransition
	}
	return nil
}

// ensureNoDuplicateSociety checks for duplicate society code, pending creator request, or matching society.
func (s *SocietySvc) ensureNoDuplicateSociety(ctx context.Context, req models.CreateSocietyRequest, requestedBy int64) error {
	if req.SocietyCode != nil {
		existing, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{Code: req.SocietyCode})
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrSocietyConflict
		}
	}

	pending, err := s.societyRepo.CountPendingByCreator(ctx, requestedBy)
	if err != nil {
		return err
	}
	if pending > 0 {
		return ErrSocietyConflict
	}

	candidates, err := s.societyRepo.List(ctx, models.ListSocietiesFilter{
		Name:  req.Name,
		City:  stringValue(req.City),
		State: stringValue(req.State),
		Limit: 1,
	})
	if err != nil {
		return err
	}
	if len(candidates) > 0 {
		return ErrSocietyConflict
	}
	return nil
}

// ensureNoActiveOwner checks that the society does not already have an active owner.
func (s *SocietySvc) ensureNoActiveOwner(ctx context.Context, societyID int64) error {
	count, err := s.memberRepo.CountActiveOwners(ctx, societyID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrOwnerProtection
	}
	return nil
}

// ensureNotLastOwner checks that owner changes will not remove the society's last active owner.
func (s *SocietySvc) ensureNotLastOwner(ctx context.Context, societyID int64) error {
	count, err := s.memberRepo.CountActiveOwners(ctx, societyID)
	if err != nil {
		return err
	}
	if count <= 1 {
		return ErrOwnerProtection
	}
	return nil
}
