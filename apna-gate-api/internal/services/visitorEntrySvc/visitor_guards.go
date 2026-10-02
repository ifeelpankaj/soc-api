package visitorentrysvc

import (
	"context"
	"errors"

	"go-server/internal/models"
	flatauthz "go-server/internal/services/flatAuthz"
)

func (s *VisitorEntrySvc) ensureApprovalActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	return mapFlatAuthzError(s.flatAuthz.CanManageFlatVisitors(ctx, societyID, flatID, actorUserID))
}

func (s *VisitorEntrySvc) ensureFlatResident(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	return mapFlatAuthzError(s.flatAuthz.CanViewFlatVisitors(ctx, societyID, flatID, actorUserID))
}

func (s *VisitorEntrySvc) ensureFlatVisitorContextActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	return mapFlatAuthzError(s.flatAuthz.CanViewFlatVisitors(ctx, societyID, flatID, actorUserID))
}

func mapFlatAuthzError(err error) error {
	if err == nil {
		return nil
	}

	var appErr *models.AppError
	if !errors.As(err, &appErr) {
		return err
	}

	switch appErr.Code {
	case flatauthz.ErrForbidden.Code, flatauthz.ErrViewForbidden.Code:
		return ErrVisitorForbidden
	case flatauthz.ErrFlatNotFound.Code:
		return ErrVisitorFlatNotFound
	default:
		return err
	}
}

func (s *VisitorEntrySvc) ensureStaffActor(ctx context.Context, societyID int64, actorUserID int64) error {
	if actorUserID <= 0 {
		return ErrVisitorForbidden
	}
	active := string(models.SocietyMemberStatusActive)
	member, err := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &societyID, UserID: &actorUserID, Status: &active})
	if err != nil {
		return err
	}
	if member == nil {
		return ErrVisitorForbidden
	}
	switch member.Role {
	case models.SocietyMemberRoleOwner, models.SocietyMemberRoleAdmin, models.SocietyMemberRoleStaff:
		return nil
	default:
		return ErrVisitorForbidden
	}
}

func (s *VisitorEntrySvc) ensureGuardOnBehalfAllowed(ctx context.Context, societyID int64, entry *models.VisitorEntry) error {
	approvalRequired, err := s.settingSvc.ResolveApprovalRequirement(ctx, societyID, entry.FlatID, entry.Purpose, entry.Source)
	if err != nil {
		return err
	}
	if !approvalRequired {
		return nil
	}
	settings, err := s.settingSvc.GetSocietySettings(ctx, societyID)
	if err != nil {
		return err
	}
	if settings == nil {
		return ErrVisitorSettingsNotFound
	}
	if !settings.AllowGuardOnBehalfApproval {
		return ErrVisitorGuardOnBehalfNotAllowed
	}
	return nil
}
