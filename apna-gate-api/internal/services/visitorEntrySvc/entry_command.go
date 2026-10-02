package visitorentrysvc

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorEntrySvc) CreatePublicQREntry(ctx context.Context, societyID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error) {
	return s.createEntryFromForm(ctx, societyID, req, models.VisitorEntrySourcePublicQR, nil)
}

func (s *VisitorEntrySvc) CreateQuickLinkEntry(ctx context.Context, societyID int64, req models.VisitorFormRequest) (*models.VisitorEntryMutationResponse, error) {
	return s.createEntryFromForm(ctx, societyID, req, models.VisitorEntrySourceQuickLink, nil)
}

func (s *VisitorEntrySvc) CreateGuardEntry(ctx context.Context, societyID int64, req models.VisitorFormRequest, guardUserID int64) (*models.VisitorEntryMutationResponse, error) {
	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}
	return s.createEntryFromForm(ctx, societyID, req, models.VisitorEntrySourceGuardEntry, &guardUserID)
}

func (s *VisitorEntrySvc) createEntryFromForm(ctx context.Context, societyID int64, req models.VisitorFormRequest, source models.VisitorEntrySource, actorUserID *int64) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	isStaffEntry := req.Purpose == models.VisitorPurposeStaff
	isSocietyService := req.Purpose == models.VisitorPurposeService && source == models.VisitorEntrySourceGuardEntry
	if isSocietyService {
		req.FlatID = 0
	}
	if err := req.Validate(!isSocietyService); err != nil {
		return nil, ErrInvalidVisitorRequest.WithCause(err)
	}
	if err := req.ValidateForPurpose(); err != nil {
		return nil, ErrInvalidVisitorRequest.WithCause(err)
	}
	if !isStaffEntry && !isSocietyService {
		if err := s.ensureEntryFlat(ctx, societyID, req.FlatID); err != nil {
			return nil, err
		}
	}
	if err := s.ensureSocietyActive(ctx, societyID); err != nil {
		return nil, err
	}
	if err := s.applyExpectedCheckout(ctx, societyID, &req); err != nil {
		return nil, err
	}

	status := models.VisitorStatusApproved
	if !isStaffEntry && !isSocietyService {
		approvalRequired, err := s.settingSvc.ResolveApprovalRequirement(ctx, societyID, req.FlatID, req.Purpose, source)
		if err != nil {
			return nil, err
		}
		if approvalRequired {
			status = models.VisitorStatusWaitingApproval
		}
	}
	var response *models.VisitorEntryMutationResponse
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var qr *qrToken
		if status == models.VisitorStatusApproved {
			var err error
			qr, err = s.makeQR(ctx, societyID)
			if err != nil {
				return err
			}
		}
		visitor, err := s.visitorRepo.Create(txCtx, contracts.VisitorFormInput(req))
		if err != nil {
			return err
		}
		var qrHash *string
		var qrExpiresAt *time.Time
		if qr != nil {
			qrHash = &qr.hash
			qrExpiresAt = &qr.expiresAt
		}
		entry, err := s.entryRepo.Create(txCtx, contracts.VisitorFormInput(req), societyID, entryFlatIDPtr(contracts.VisitorFormInput(req).FlatID), visitor.ID, nil, source, contracts.VisitorFormInput(req).Purpose, status, actorUserID, guardActor(source, actorUserID), qrHash, qrExpiresAt)
		if err != nil {
			return err
		}
		events := []models.VisitorEventType{models.VisitorEventTypeCreated}
		if status == models.VisitorStatusApproved {
			events = append(events, models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated)
		}
		if err := s.recordEvents(txCtx, entry, actorUserID, events...); err != nil {
			return err
		}
		if qr != nil {
			entry, err = s.attachQRDisplayToken(txCtx, societyID, entry.ID, qr)
			if err != nil {
				return err
			}
		}
		entry.Visitor = &models.VisitorSummary{
			FullName:    visitor.FullName,
			PhoneNumber: visitor.PhoneNumber,
			Email:       visitor.Email,
			PhotoURL:    visitor.PhotoURL,
		}
		response = &models.VisitorEntryMutationResponse{Entry: entry}
		if qr != nil {
			response.QR = qr.response()
		}
		return nil
	})
	if err != nil {
		return response, err
	}
	if response != nil && response.Entry != nil {
		if response.Entry.Status == models.VisitorStatusWaitingApproval {
			s.notifyVisitorPending(response.Entry)
		} else {
			s.notifyVisitorApproved(response.Entry)
		}
	}
	return response, err
}

func (s *VisitorEntrySvc) ApproveEntry(ctx context.Context, societyID int64, entryID int64, actorUserID int64) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureApprovalActor(ctx, societyID, entry.FlatID, actorUserID); err != nil {
		return nil, err
	}
	switch entry.Status {
	case models.VisitorStatusApproved, models.VisitorStatusCheckedIn, models.VisitorStatusCheckedOut:
		return &models.VisitorEntryMutationResponse{Entry: entry}, nil
	case models.VisitorStatusWaitingApproval:
		// Continue with the state transition below.
	default:
		return nil, ErrVisitorInvalidState
	}
	qr, err := s.makeQR(ctx, societyID)
	if err != nil {
		return nil, err
	}
	var approved *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.lockEntryForApproval(txCtx, societyID, entryID); err != nil {
			return err
		}
		var err error
		approved, err = s.entryRepo.Approve(txCtx, societyID, entryID, actorUserID, qr.hash, qr.expiresAt)
		if err != nil {
			return err
		}
		if approved == nil {
			return ErrVisitorInvalidState
		}
		return s.recordEvents(txCtx, approved, &actorUserID, models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated)
	})
	if err != nil {
		if errors.Is(err, ErrVisitorInvalidState) {
			current, getErr := s.GetEntry(ctx, societyID, entryID)
			if getErr == nil && current != nil {
				switch current.Status {
				case models.VisitorStatusApproved, models.VisitorStatusCheckedIn, models.VisitorStatusCheckedOut:
					return &models.VisitorEntryMutationResponse{Entry: current}, nil
				}
			}
		}
		return nil, err
	}
	approved = s.enrichNotificationDetails(ctx, approved, entry)
	s.notifyVisitorApproved(approved)
	return &models.VisitorEntryMutationResponse{Entry: approved, QR: qr.response()}, nil
}

func (s *VisitorEntrySvc) RejectEntry(ctx context.Context, societyID int64, entryID int64, reason string, actorUserID int64) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureApprovalActor(ctx, societyID, entry.FlatID, actorUserID); err != nil {
		return nil, err
	}
	switch entry.Status {
	case models.VisitorStatusRejected:
		return &models.VisitorEntryMutationResponse{Entry: entry}, nil
	case models.VisitorStatusWaitingApproval:
		// Continue with the state transition below.
	default:
		return nil, ErrVisitorInvalidState
	}
	var rejected *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		rejected, err = s.entryRepo.Reject(txCtx, societyID, entryID, actorUserID, reason)
		if err != nil {
			return err
		}
		if rejected == nil {
			return ErrVisitorInvalidState
		}
		return s.recordEvents(txCtx, rejected, &actorUserID, models.VisitorEventTypeRejected)
	})
	if err != nil {
		if errors.Is(err, ErrVisitorInvalidState) {
			current, getErr := s.GetEntry(ctx, societyID, entryID)
			if getErr == nil && current != nil && current.Status == models.VisitorStatusRejected {
				return &models.VisitorEntryMutationResponse{Entry: current}, nil
			}
		}
		return nil, err
	}
	rejected = s.enrichNotificationDetails(ctx, rejected, entry)
	s.notifyVisitorRejected(rejected)
	return &models.VisitorEntryMutationResponse{Entry: rejected}, nil
}

func (s *VisitorEntrySvc) GenerateQR(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	qr, err := s.makeQR(ctx, societyID)
	if err != nil {
		return nil, err
	}
	var entry *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		entry, err = s.entryRepo.GenerateQR(txCtx, societyID, entryID, qr.hash, qr.expiresAt)
		if err != nil {
			return err
		}
		if entry == nil {
			return ErrVisitorInvalidState
		}
		return s.recordEvents(txCtx, entry, nil, models.VisitorEventTypeQRGenerated)
	})
	if err != nil {
		return nil, err
	}
	return &models.VisitorEntryMutationResponse{Entry: entry, QR: qr.response()}, nil
}

func (s *VisitorEntrySvc) AutoCloseExpiredEntries(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.entryRepo.AutoCloseExpired(ctx)
}

func (s *VisitorEntrySvc) ExpireStaleEntries(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.entryRepo.ExpireStaleEntries(ctx)
}
