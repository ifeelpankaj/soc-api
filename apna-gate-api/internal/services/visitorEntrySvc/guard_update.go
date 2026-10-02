package visitorentrysvc

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorEntrySvc) UpdateGuardEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64, req models.UpdateGuardVisitorEntryRequest) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := req.Validate(); err != nil {
		return nil, ErrInvalidVisitorRequest.WithCause(err)
	}
	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry.Status != models.VisitorStatusWaitingApproval && entry.Status != models.VisitorStatusApproved {
		return nil, ErrVisitorInvalidState
	}

	if req.FlatID != nil {
		if entry.IsSocietyWideGuardService() {
			return nil, ErrInvalidVisitorRequest.WithCause(errors.New("flat cannot be assigned to society-wide service entries"))
		}
		if entry.Source != models.VisitorEntrySourceGuardEntry {
			return nil, ErrInvalidVisitorRequest.WithCause(errors.New("flat can only be changed for guard-created entries"))
		}
		if entry.Purpose == models.VisitorPurposeStaff {
			return nil, ErrInvalidVisitorRequest.WithCause(errors.New("flat cannot be changed for staff entries"))
		}
		if err := s.ensureEntryFlat(ctx, societyID, *req.FlatID); err != nil {
			return nil, err
		}
	}

	flatChanged := req.FlatID != nil && entry.FlatID != *req.FlatID
	var qrForAutoApprove *qrToken
	if flatChanged && entry.Status == models.VisitorStatusWaitingApproval {
		approvalRequired, resolveErr := s.settingSvc.ResolveApprovalRequirement(ctx, societyID, *req.FlatID, entry.Purpose, entry.Source)
		if resolveErr != nil {
			return nil, resolveErr
		}
		if !approvalRequired {
			qrForAutoApprove, err = s.makeQR(ctx, societyID)
			if err != nil {
				return nil, err
			}
		}
	}

	var updated *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Match image mutations: acquire the entry lock before updating its visitor.
		locked, lockErr := s.entryRepo.GetForUpdate(txCtx, societyID, entryID)
		if lockErr != nil {
			return lockErr
		}
		if locked == nil || (locked.Status != models.VisitorStatusWaitingApproval && locked.Status != models.VisitorStatusApproved) {
			return ErrVisitorInvalidState
		}
		if req.FullName != nil || req.PhoneNumber != nil || req.Email != nil || req.PhotoURL != nil {
			if _, err := s.visitorRepo.UpdateProfile(txCtx, entry.VisitorID, contracts.UpdateGuardVisitorEntryInput(req)); err != nil {
				return err
			}
		}
		if req.FlatID != nil || req.VehicleNumber != nil || req.VehicleType != nil || req.CompanionsCount != nil || req.CompanionDetails != nil || req.Notes != nil {
			var detailErr error
			updated, detailErr = s.entryRepo.UpdateDetails(txCtx, societyID, entryID, contracts.UpdateGuardVisitorEntryInput(req))
			if detailErr != nil {
				return detailErr
			}
		}
		if qrForAutoApprove != nil {
			approved, approveErr := s.entryRepo.Approve(txCtx, societyID, entryID, guardUserID, qrForAutoApprove.hash, qrForAutoApprove.expiresAt)
			if approveErr != nil {
				return approveErr
			}
			if approved == nil {
				return ErrVisitorInvalidState
			}
			approved, approveErr = s.attachQRDisplayToken(txCtx, societyID, entryID, qrForAutoApprove)
			if approveErr != nil {
				return approveErr
			}
			updated = approved
			return s.recordEvents(txCtx, approved, &guardUserID, models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated)
		}
		if updated == nil {
			var getErr error
			updated, getErr = s.entryRepo.Get(txCtx, societyID, entryID)
			return getErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if qrForAutoApprove != nil && updated != nil {
		s.notifyVisitorApproved(updated)
	}
	return s.GetEntry(ctx, societyID, entryID)
}
