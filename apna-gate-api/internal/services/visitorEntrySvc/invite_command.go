package visitorentrysvc

import (
	"context"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorEntrySvc) CreateInvite(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, actorUserID int64, expiresAt *time.Time) (*models.VisitorEntryMutationResponse, *models.VisitorInvite, *models.ShortLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if societyID <= 0 || flatID <= 0 || actorUserID <= 0 || !purpose.IsValid() {
		return nil, nil, nil, ErrInvalidVisitorRequest
	}
	if purpose == models.VisitorPurposeService || purpose == models.VisitorPurposeStaff {
		return nil, nil, nil, ErrInvalidVisitorRequest
	}
	if err := s.ensureApprovalActor(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, nil, nil, err
	}
	return s.createInvite(ctx, societyID, flatID, purpose, actorUserID, expiresAt)
}

func (s *VisitorEntrySvc) CreateStaffInvite(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, staffUserID int64, expiresAt *time.Time) (*models.VisitorEntryMutationResponse, *models.VisitorInvite, *models.ShortLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if societyID <= 0 || flatID <= 0 || staffUserID <= 0 || !purpose.IsValid() {
		return nil, nil, nil, ErrInvalidVisitorRequest
	}
	if err := s.ensureStaffActor(ctx, societyID, staffUserID); err != nil {
		return nil, nil, nil, err
	}
	flat, err := s.flatRepo.Get(ctx, &models.FlatFilter{ID: &flatID, SocietyID: &societyID})
	if err != nil {
		return nil, nil, nil, err
	}
	if flat == nil {
		return nil, nil, nil, ErrVisitorFlatNotFound
	}
	return s.createInvite(ctx, societyID, flatID, purpose, staffUserID, expiresAt)
}

func (s *VisitorEntrySvc) createInvite(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, actorUserID int64, expiresAt *time.Time) (*models.VisitorEntryMutationResponse, *models.VisitorInvite, *models.ShortLinkResponse, error) {
	if _, err := s.settingSvc.ResolveApprovalRequirement(ctx, societyID, flatID, purpose, models.VisitorEntrySourceResidentLink); err != nil {
		return nil, nil, nil, err
	}
	token, tokenHash, err := newToken()
	if err != nil {
		return nil, nil, nil, err
	}
	expiry := time.Now().Add(defaultInviteDuration)
	if expiresAt != nil {
		expiry = *expiresAt
	}
	invite, err := s.inviteRepo.Create(ctx, societyID, flatID, purpose, tokenHash, expiry, actorUserID)
	if err != nil {
		return nil, nil, nil, err
	}
	var link *models.ShortLinkResponse
	if s.shortLinks != nil {
		link, err = s.shortLinks.Create(ctx, models.ShortLinkResourceVisitorInvite, invite.ID, &invite.ExpiresAt, &actorUserID)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return &models.VisitorEntryMutationResponse{QR: &models.QRTokenResponse{Token: token, ExpiresAt: expiry}}, invite, link, nil
}

func (s *VisitorEntrySvc) SubmitInviteForm(ctx context.Context, societyID int64, rawToken string, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if societyID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}
	if err := req.ValidateInviteSubmit(); err != nil {
		return nil, ErrInvalidVisitorRequest.WithCause(err)
	}
	invite, err := s.getUsableInviteByToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if invite.SocietyID != societyID {
		return nil, ErrVisitorInviteNotFound
	}
	return s.submitInvite(ctx, invite, req)
}

func (s *VisitorEntrySvc) SubmitInviteFormByID(ctx context.Context, inviteID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if inviteID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}
	if err := req.ValidateInviteSubmit(); err != nil {
		return nil, ErrInvalidVisitorRequest.WithCause(err)
	}
	invite, err := s.inviteRepo.GetForUpdate(ctx, inviteID)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrVisitorInviteNotFound
	}
	return s.submitInvite(ctx, invite, req)
}

func (s *VisitorEntrySvc) submitInvite(ctx context.Context, invite *models.VisitorInvite, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error) {
	if err := s.ensureSocietyActive(ctx, invite.SocietyID); err != nil {
		return nil, err
	}
	settings, err := s.settingSvc.GetSocietySettings(ctx, invite.SocietyID)
	if err != nil {
		return nil, err
	}
	if settings == nil || !settings.AllowResidentPreApproval || !settings.IsActive {
		return nil, ErrVisitorInviteUnavailable
	}

	var response *models.VisitorEntryMutationResponse
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		locked, lockErr := s.inviteRepo.GetForUpdate(txCtx, invite.ID)
		if lockErr != nil {
			return lockErr
		}
		if locked == nil {
			return ErrVisitorInviteNotFound
		}

		if locked.Status == models.VisitorInviteStatusUsed {
			replay, replayErr := s.idempotentInviteSubmitResponse(txCtx, locked)
			if replayErr != nil {
				return replayErr
			}
			response = replay
			return nil
		}
		if !inviteUsable(locked) {
			return ErrVisitorInviteUnavailable
		}

		qr, qrErr := s.makeQR(txCtx, locked.SocietyID)
		if qrErr != nil {
			return qrErr
		}
		visitor, visitorErr := s.visitorRepo.Create(txCtx, contracts.VisitorFormInput(req))
		if visitorErr != nil {
			return visitorErr
		}
		entryReq := req
		if err := s.applyExpectedCheckout(ctx, locked.SocietyID, &entryReq); err != nil {
			return err
		}
		entry, createErr := s.entryRepo.Create(txCtx, contracts.VisitorFormInput(entryReq), locked.SocietyID, entryFlatIDPtr(locked.FlatID), visitor.ID, &locked.ID, models.VisitorEntrySourceResidentLink, locked.Purpose, models.VisitorStatusApproved, &locked.CreatedBy, nil, &qr.hash, &qr.expiresAt)
		if createErr != nil {
			if isUniqueViolation(createErr) {
				replay, replayErr := s.idempotentInviteSubmitResponse(txCtx, locked)
				if replayErr != nil {
					return replayErr
				}
				response = replay
				return nil
			}
			return createErr
		}
		marked, markErr := s.inviteRepo.MarkUsed(txCtx, locked.ID)
		if markErr != nil {
			return markErr
		}
		if marked == nil {
			replay, replayErr := s.idempotentInviteSubmitResponse(txCtx, locked)
			if replayErr != nil {
				return replayErr
			}
			response = replay
			return nil
		}
		if err := s.recordEvents(txCtx, entry, &locked.CreatedBy, models.VisitorEventTypeCreated, models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated); err != nil {
			return err
		}
		entry, err = s.attachQRDisplayToken(txCtx, locked.SocietyID, entry.ID, qr)
		if err != nil {
			return err
		}
		entry.Visitor = &models.VisitorSummary{
			FullName:    visitor.FullName,
			PhoneNumber: visitor.PhoneNumber,
			Email:       visitor.Email,
			PhotoURL:    visitor.PhotoURL,
		}
		response = &models.VisitorEntryMutationResponse{Entry: entry, QR: qr.response()}
		return nil
	})
	if err != nil {
		return response, err
	}
	if response != nil && response.Entry != nil && !response.IdempotentReplay {
		s.notifyVisitorInviteAccepted(response.Entry)
	}
	return response, err
}

func (s *VisitorEntrySvc) CancelInvite(ctx context.Context, societyID int64, inviteID int64, actorUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	invite, err := s.inviteRepo.GetByID(ctx, societyID, inviteID)
	if err != nil {
		return err
	}
	if invite == nil {
		return ErrVisitorInviteNotFound
	}
	if err := s.ensureApprovalActor(ctx, societyID, invite.FlatID, actorUserID); err != nil {
		return err
	}
	cancelled, err := s.inviteRepo.Cancel(ctx, societyID, inviteID)
	if err != nil {
		return err
	}
	if cancelled == nil {
		return ErrVisitorInviteUnavailable
	}
	return nil
}

func (s *VisitorEntrySvc) ExpireOldInvites(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.inviteRepo.ExpireOld(ctx)
}
