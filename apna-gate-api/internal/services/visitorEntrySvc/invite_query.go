package visitorentrysvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorEntrySvc) GetInviteByToken(ctx context.Context, societyID int64, guardUserID int64, rawToken string) (*models.VisitorInvite, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if societyID <= 0 || guardUserID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}
	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}
	invite, err := s.getUsableInviteByToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if invite.SocietyID != societyID {
		return nil, ErrVisitorInviteNotFound
	}
	return invite, nil
}

func (s *VisitorEntrySvc) getUsableInviteByToken(ctx context.Context, rawToken string) (*models.VisitorInvite, error) {
	invite, err := s.inviteRepo.GetByTokenHash(ctx, hashToken(rawToken))
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrVisitorInviteNotFound
	}
	if !inviteUsable(invite) {
		return nil, ErrVisitorInviteUnavailable
	}
	return invite, nil
}

func (s *VisitorEntrySvc) GetPublicInviteByToken(ctx context.Context, rawToken string) (*models.PublicVisitorInvitePageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	invite, err := s.inviteRepo.GetByTokenHash(ctx, hashToken(rawToken))
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrVisitorInviteNotFound
	}
	if err := s.ensureSocietyActive(ctx, invite.SocietyID); err != nil {
		return nil, err
	}

	inviteView, err := s.buildPublicInviteView(ctx, invite)
	if err != nil {
		return nil, err
	}

	if inviteUsable(invite) {
		return &models.PublicVisitorInvitePageResponse{
			Invite: inviteView,
			View:   models.PublicVisitorInviteViewForm,
		}, nil
	}

	if invite.Status == models.VisitorInviteStatusUsed {
		return s.publicInvitePageForUsed(ctx, inviteView, invite)
	}

	return nil, ErrVisitorInviteUnavailable
}

func (s *VisitorEntrySvc) GetPublicInviteByID(ctx context.Context, inviteID int64) (*models.PublicVisitorInvitePageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if inviteID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}
	invite, err := s.inviteRepo.GetForUpdate(ctx, inviteID)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrVisitorInviteNotFound
	}
	if err := s.ensureSocietyActive(ctx, invite.SocietyID); err != nil {
		return nil, err
	}
	inviteView, err := s.buildPublicInviteView(ctx, invite)
	if err != nil {
		return nil, err
	}
	if inviteUsable(invite) {
		return &models.PublicVisitorInvitePageResponse{
			Invite: inviteView,
			View:   models.PublicVisitorInviteViewForm,
		}, nil
	}
	if invite.Status == models.VisitorInviteStatusUsed {
		return s.publicInvitePageForUsed(ctx, inviteView, invite)
	}
	return &models.PublicVisitorInvitePageResponse{
		Invite: inviteView,
		View:   models.PublicVisitorInviteViewClosed,
	}, nil
}

func (s *VisitorEntrySvc) buildPublicInviteView(ctx context.Context, invite *models.VisitorInvite) (*models.PublicVisitorInviteView, error) {
	flat, err := s.flatRepo.Get(ctx, &models.FlatFilter{ID: &invite.FlatID, SocietyID: &invite.SocietyID})
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrVisitorFlatNotFound
	}
	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &invite.SocietyID})
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrVisitorFlatNotFound
	}
	return &models.PublicVisitorInviteView{
		ID:          invite.ID,
		Purpose:     invite.Purpose,
		Status:      invite.Status,
		ExpiresAt:   invite.ExpiresAt,
		SocietyName: society.Name,
		FlatNumber:  flat.FlatNumber,
		Block:       flat.Block,
		Floor:       flat.Floor,
	}, nil
}

func (s *VisitorEntrySvc) publicInvitePageForUsed(ctx context.Context, inviteView *models.PublicVisitorInviteView, invite *models.VisitorInvite) (*models.PublicVisitorInvitePageResponse, error) {
	entry, err := s.entryRepo.GetByInviteID(ctx, invite.ID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorInviteUnavailable
	}

	page := &models.PublicVisitorInvitePageResponse{
		Invite: inviteView,
		Entry:  entry,
	}

	switch entry.Status {
	case models.VisitorStatusApproved:
		page.View = models.PublicVisitorInviteViewQR
		qrResult, err := s.resolveExistingQR(ctx, invite.SocietyID, entry)
		if err != nil {
			return nil, err
		}
		enriched, err := s.entryRepo.GetByInviteID(ctx, invite.ID)
		if err != nil {
			return nil, err
		}
		page.Entry = enriched
		page.QR = qrResult
		return page, nil
	case models.VisitorStatusCheckedIn:
		page.View = models.PublicVisitorInviteViewCheckedIn
		return page, nil
	case models.VisitorStatusCheckedOut:
		page.View = models.PublicVisitorInviteViewCheckedOut
		return page, nil
	default:
		page.View = models.PublicVisitorInviteViewClosed
		return page, nil
	}
}
