package flatauthz

import (
	"context"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
)

type FlatVisitorAuthz struct {
	memberRepo   repository.SocietyMemberRepository
	residentRepo repository.FlatResidentRepository
	flatRepo     repository.FlatRepository
}

func New(
	memberRepo repository.SocietyMemberRepository,
	residentRepo repository.FlatResidentRepository,
	flatRepo repository.FlatRepository,
) *FlatVisitorAuthz {
	return &FlatVisitorAuthz{
		memberRepo:   memberRepo,
		residentRepo: residentRepo,
		flatRepo:     flatRepo,
	}
}

func (a *FlatVisitorAuthz) IsSocietyOwnerOrAdmin(ctx context.Context, societyID int64, userID int64) (bool, error) {
	if userID <= 0 {
		return false, nil
	}

	active := string(models.SocietyMemberStatusActive)
	member, err := a.memberRepo.Get(ctx, models.GetSocietyMemberFilter{
		SocietyID: &societyID,
		UserID:    &userID,
		Status:    &active,
	})
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}

	return member.Role == models.SocietyMemberRoleOwner || member.Role == models.SocietyMemberRoleAdmin, nil
}

func (a *FlatVisitorAuthz) IsSocietyStaff(ctx context.Context, societyID int64, userID int64) (bool, error) {
	if userID <= 0 {
		return false, nil
	}

	active := string(models.SocietyMemberStatusActive)
	member, err := a.memberRepo.Get(ctx, models.GetSocietyMemberFilter{
		SocietyID: &societyID,
		UserID:    &userID,
		Status:    &active,
	})
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}

	return member.Role == models.SocietyMemberRoleStaff, nil
}

func (a *FlatVisitorAuthz) CanManageFlatVisitors(ctx context.Context, societyID int64, flatID int64, userID int64) error {
	if userID <= 0 {
		return ErrForbidden
	}

	if err := a.ensureFlatInSociety(ctx, societyID, flatID); err != nil {
		return err
	}

	isOwnerOrAdmin, err := a.IsSocietyOwnerOrAdmin(ctx, societyID, userID)
	if err != nil {
		return err
	}
	if isOwnerOrAdmin {
		return nil
	}

	if a.isActiveFlatResident(ctx, societyID, flatID, userID) {
		return nil
	}

	return ErrForbidden
}

func (a *FlatVisitorAuthz) CanManageFlatMembers(ctx context.Context, societyID int64, flatID int64, userID int64) error {
	if userID <= 0 {
		return ErrManageForbidden
	}

	if err := a.ensureFlatInSociety(ctx, societyID, flatID); err != nil {
		return err
	}

	resident, err := a.getActiveFlatResident(ctx, societyID, flatID, userID)
	if err != nil {
		return err
	}
	if resident == nil {
		return ErrManageForbidden
	}
	if resident.IsPrimary || resident.Role == models.FlatResidentRoleOwner {
		return nil
	}

	return ErrManageForbidden
}

// CanViewFlatMemberInvites authorizes invite-history access and reports whether
// the caller may see reusable invite links. Society owners and admins can audit
// invite history, but only the flat owner or primary resident may manage links.
func (a *FlatVisitorAuthz) CanViewFlatMemberInvites(ctx context.Context, societyID int64, flatID int64, userID int64) (bool, error) {
	if userID <= 0 {
		return false, ErrManageForbidden
	}
	if err := a.ensureFlatInSociety(ctx, societyID, flatID); err != nil {
		return false, err
	}

	resident, err := a.getActiveFlatResident(ctx, societyID, flatID, userID)
	if err != nil {
		return false, err
	}
	if resident != nil && (resident.IsPrimary || resident.Role == models.FlatResidentRoleOwner) {
		return true, nil
	}

	isOwnerOrAdmin, err := a.IsSocietyOwnerOrAdmin(ctx, societyID, userID)
	if err != nil {
		return false, err
	}
	if isOwnerOrAdmin {
		return false, nil
	}
	return false, ErrManageForbidden
}

func (a *FlatVisitorAuthz) CanViewFlatVisitors(ctx context.Context, societyID int64, flatID int64, userID int64) error {
	if userID <= 0 {
		return ErrViewForbidden
	}

	if err := a.ensureFlatInSociety(ctx, societyID, flatID); err != nil {
		return err
	}

	isOwnerOrAdmin, err := a.IsSocietyOwnerOrAdmin(ctx, societyID, userID)
	if err != nil {
		return err
	}
	if isOwnerOrAdmin {
		return nil
	}

	isStaff, err := a.IsSocietyStaff(ctx, societyID, userID)
	if err != nil {
		return err
	}
	if isStaff {
		return nil
	}

	if a.isActiveFlatResident(ctx, societyID, flatID, userID) {
		return nil
	}

	return ErrViewForbidden
}

func (a *FlatVisitorAuthz) ensureFlatInSociety(ctx context.Context, societyID int64, flatID int64) error {
	flat, err := a.flatRepo.Get(ctx, &models.FlatFilter{ID: &flatID, SocietyID: &societyID})
	if err != nil {
		return err
	}
	if flat == nil {
		return ErrFlatNotFound
	}
	return nil
}

func (a *FlatVisitorAuthz) isActiveFlatResident(ctx context.Context, societyID int64, flatID int64, userID int64) bool {
	resident, err := a.getActiveFlatResident(ctx, societyID, flatID, userID)
	return err == nil && resident != nil
}

func (a *FlatVisitorAuthz) getActiveFlatResident(ctx context.Context, societyID int64, flatID int64, userID int64) (*models.FlatResident, error) {
	status := string(models.FlatResidentStatusActive)
	return a.residentRepo.Get(ctx, &models.FlatResidentFilter{
		SocietyID: &societyID,
		FlatID:    &flatID,
		UserID:    &userID,
		Status:    &status,
	})
}
