package societysvc

import (
	"context"
	"strings"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *SocietySvc) AddMember(ctx context.Context, req models.AddSocietyMemberRequest, addedBy int64) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if !req.Role.IsValid() {
		return nil, ErrInvalidMemberRole
	}
	if err := s.EnsureRole(ctx, req.SocietyID, addedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}
	if s.subscriptionSvc != nil {
		if req.Role == models.SocietyMemberRoleAdmin {
			if err := s.subscriptionSvc.CanAddAdmin(ctx, req.SocietyID, 1); err != nil {
				return nil, err
			}
		}
		if req.Role == models.SocietyMemberRoleStaff {
			if err := s.subscriptionSvc.CanAddStaff(ctx, req.SocietyID, 1); err != nil {
				return nil, err
			}
		}
	}
	if req.Role == models.SocietyMemberRoleOwner {
		if err := s.ensureNoActiveOwner(ctx, req.SocietyID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureUserExists(ctx, req.UserID); err != nil {
		return nil, err
	}

	member := &models.SocietyMember{
		SocietyID: req.SocietyID,
		UserID:    req.UserID,
		Role:      req.Role,
		Status:    models.SocietyMemberStatusActive,
		InvitedBy: &addedBy,
		Metadata:  req.Metadata,
	}
	if err := s.memberRepo.Add(ctx, member); err != nil {
		return nil, ErrMemberConflict.WithCause(err)
	}

	loaded, _ := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &req.SocietyID, UserID: &req.UserID})
	if loaded != nil {
		return loaded.ToResponse(), nil
	}
	return member.ToResponse(), nil
}

func (s *SocietySvc) ChangeMemberRole(ctx context.Context, req models.ChangeSocietyMemberRoleRequest, changedBy int64) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if !req.Role.IsValid() {
		return nil, ErrInvalidMemberRole
	}
	if err := s.EnsureRole(ctx, req.SocietyID, changedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}

	current, err := s.requireMember(ctx, req.SocietyID, req.UserID)
	if err != nil {
		return nil, err
	}
	if current.Role == models.SocietyMemberRoleOwner && req.Role != models.SocietyMemberRoleOwner {
		if err := s.ensureNotLastOwner(ctx, req.SocietyID); err != nil {
			return nil, err
		}
	}
	if req.Role == models.SocietyMemberRoleOwner && current.Role != models.SocietyMemberRoleOwner {
		if err := s.ensureNoActiveOwner(ctx, req.SocietyID); err != nil {
			return nil, err
		}
	}

	member, err := s.memberRepo.ChangeRole(ctx, req.SocietyID, req.UserID, req.Role)
	if err != nil {
		return nil, err
	}
	return member.ToResponse(), nil
}

func (s *SocietySvc) SuspendMember(ctx context.Context, req models.SuspendSocietyMemberRequest, suspendedBy int64) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.EnsureRole(ctx, req.SocietyID, suspendedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}

	current, err := s.requireMember(ctx, req.SocietyID, req.UserID)
	if err != nil {
		return nil, err
	}
	if current.Role == models.SocietyMemberRoleOwner {
		if err := s.ensureNotLastOwner(ctx, req.SocietyID); err != nil {
			return nil, err
		}
	}

	member, err := s.memberRepo.Suspend(ctx, req.SocietyID, req.UserID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrMemberInactive
	}
	return member.ToResponse(), nil
}

func (s *SocietySvc) ReactivateMember(ctx context.Context, req models.ReactivateSocietyMemberRequest, reactivatedBy int64) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.EnsureRole(ctx, req.SocietyID, reactivatedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}

	member, err := s.memberRepo.Reactivate(ctx, req.SocietyID, req.UserID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrMemberNotFound
	}
	return member.ToResponse(), nil
}

func (s *SocietySvc) RemoveMember(ctx context.Context, req models.RemoveSocietyMemberRequest, removedBy int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.EnsureRole(ctx, req.SocietyID, removedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return err
	}

	current, err := s.requireMember(ctx, req.SocietyID, req.UserID)
	if err != nil {
		return err
	}
	if current.Role == models.SocietyMemberRoleOwner {
		if err := s.ensureNotLastOwner(ctx, req.SocietyID); err != nil {
			return err
		}
	}

	return s.memberRepo.Remove(ctx, req.SocietyID, req.UserID, removedBy, strings.TrimSpace(req.Reason))
}

func (s *SocietySvc) TransferOwnership(ctx context.Context, societyID int64, newOwnerUserID int64, changedBy int64) (*models.SocietyMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.EnsureRole(ctx, societyID, changedBy, string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)); err != nil {
		return nil, err
	}
	if err := s.ensureUserExists(ctx, newOwnerUserID); err != nil {
		return nil, err
	}

	var promoted *models.SocietyMember
	if err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existing, err := s.memberRepo.Get(txCtx, models.GetSocietyMemberFilter{SocietyID: &societyID, UserID: &newOwnerUserID})
		if err != nil {
			return err
		}
		if existing == nil {
			member := &models.SocietyMember{
				SocietyID: societyID,
				UserID:    newOwnerUserID,
				Role:      models.SocietyMemberRoleAdmin,
				Status:    models.SocietyMemberStatusActive,
				InvitedBy: &changedBy,
			}
			if err := s.memberRepo.Add(txCtx, member); err != nil {
				return ErrMemberConflict.WithCause(err)
			}
		}

		if err := s.memberRepo.DemoteActiveOwners(txCtx, societyID, newOwnerUserID); err != nil {
			return err
		}
		promoted, err = s.memberRepo.PromoteToOwner(txCtx, societyID, newOwnerUserID)
		if err != nil {
			return err
		}

		count, err := s.memberRepo.CountActiveOwners(txCtx, societyID)
		if err != nil {
			return err
		}
		if count != 1 {
			return ErrOwnerProtection
		}
		return nil
	}); err != nil {
		return nil, err
	}

	loaded, _ := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &societyID, UserID: &newOwnerUserID})
	if loaded != nil {
		return loaded.ToResponse(), nil
	}
	return promoted.ToResponse(), nil
}
