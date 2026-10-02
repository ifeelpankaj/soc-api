package visitorentrysvc

import (
	"context"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

type GuardApproveOptions struct {
	OnBehalf bool
	Reason   *string
}

func (s *VisitorEntrySvc) NotifyPendingEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return err
	}

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return ErrVisitorEntryNotFound
	}
	if entry.Status != models.VisitorStatusWaitingApproval {
		return ErrVisitorInvalidState
	}

	if durable, ok := s.notifier.(interface{ DurableNotifications() bool }); ok && durable.DurableNotifications() {
		return s.notifier.SendVisitorApprovalRequested(ctx, entry)
	}
	s.notifyVisitorPending(entry)
	return nil
}

func (s *VisitorEntrySvc) GuardApproveEntry(ctx context.Context, societyID int64, entryID int64, guardUserID int64, opts GuardApproveOptions) (*models.VisitorEntryMutationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorEntryNotFound
	}
	if entry.Status != models.VisitorStatusWaitingApproval {
		return nil, ErrVisitorInvalidState
	}
	if opts.OnBehalf {
		if err := s.ensureGuardOnBehalfAllowed(ctx, societyID, entry); err != nil {
			return nil, err
		}
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
		if opts.OnBehalf {
			audit, auditErr := s.buildOnBehalfAudit(txCtx, entry, guardUserID, opts.Reason)
			if auditErr != nil {
				return auditErr
			}
			if _, auditErr = s.entryRepo.MergeMetadata(txCtx, societyID, entryID, audit); auditErr != nil {
				return auditErr
			}
		}

		var approveErr error
		approved, approveErr = s.entryRepo.Approve(txCtx, societyID, entryID, guardUserID, qr.hash, qr.expiresAt)
		if approveErr != nil {
			return approveErr
		}
		if approved == nil {
			return ErrVisitorInvalidState
		}

		events := []models.VisitorEventType{models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated}
		if opts.OnBehalf {
			events = append(events, models.VisitorEventTypeGuardApprovedOnBehalf)
		}
		return s.recordEvents(txCtx, approved, &guardUserID, events...)
	})
	if err != nil {
		return nil, err
	}

	approved = s.enrichNotificationDetails(ctx, approved, entry)
	s.notifyVisitorApproved(approved)
	return &models.VisitorEntryMutationResponse{Entry: approved, QR: qr.response()}, nil
}

func (s *VisitorEntrySvc) GuardApproveAndCheckIn(ctx context.Context, societyID int64, entryID int64, guardUserID int64, opts GuardApproveOptions) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorEntryNotFound
	}
	if entry.Status != models.VisitorStatusWaitingApproval {
		return nil, ErrVisitorInvalidState
	}
	if opts.OnBehalf {
		if err := s.ensureGuardOnBehalfAllowed(ctx, societyID, entry); err != nil {
			return nil, err
		}
	}

	qr, err := s.makeQR(ctx, societyID)
	if err != nil {
		return nil, err
	}

	var checkedIn *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.lockEntryForApproval(txCtx, societyID, entryID); err != nil {
			return err
		}
		if opts.OnBehalf {
			audit, auditErr := s.buildOnBehalfAudit(txCtx, entry, guardUserID, opts.Reason)
			if auditErr != nil {
				return auditErr
			}
			if _, auditErr = s.entryRepo.MergeMetadata(txCtx, societyID, entryID, audit); auditErr != nil {
				return auditErr
			}
		}

		approved, approveErr := s.entryRepo.Approve(txCtx, societyID, entryID, guardUserID, qr.hash, qr.expiresAt)
		if approveErr != nil {
			return approveErr
		}
		if approved == nil {
			return ErrVisitorInvalidState
		}

		events := []models.VisitorEventType{models.VisitorEventTypeApproved, models.VisitorEventTypeQRGenerated}
		if opts.OnBehalf {
			events = append(events, models.VisitorEventTypeGuardApprovedOnBehalf)
		}
		if err := s.recordEvents(txCtx, approved, &guardUserID, events...); err != nil {
			return err
		}

		checkedIn, err = s.entryRepo.CheckIn(txCtx, societyID, entryID, guardUserID)
		if err != nil {
			return err
		}
		if checkedIn == nil {
			return ErrVisitorInvalidState
		}
		copyVisitorSummary(checkedIn, entry)
		return s.recordEvents(txCtx, checkedIn, &guardUserID, models.VisitorEventTypeCheckedIn)
	})
	if err != nil {
		return nil, err
	}

	checkedIn = s.enrichNotificationDetails(ctx, checkedIn, entry)
	s.notifyVisitorCheckIn(checkedIn)
	return checkedIn, nil
}

func (s *VisitorEntrySvc) CheckInByEntryID(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}

	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorEntryNotFound
	}
	if entry.Status != models.VisitorStatusApproved {
		return nil, ErrVisitorInvalidState
	}

	var checkedIn *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var checkInErr error
		checkedIn, checkInErr = s.entryRepo.CheckIn(txCtx, societyID, entryID, guardUserID)
		if checkInErr != nil {
			return checkInErr
		}
		if checkedIn == nil {
			return ErrVisitorInvalidState
		}
		copyVisitorSummary(checkedIn, entry)
		return s.recordEvents(txCtx, checkedIn, &guardUserID, models.VisitorEventTypeCheckedIn)
	})
	if err != nil {
		return nil, err
	}

	checkedIn = s.enrichVisitorSummary(ctx, checkedIn, entry)
	s.notifyVisitorCheckIn(checkedIn)
	return checkedIn, nil
}

func (s *VisitorEntrySvc) ListWaitingAtGate(ctx context.Context, filter models.WaitingAtGateFilter) (*models.VisitorEntryListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	entries, err := s.entryRepo.ListWaitingAtGate(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.entryRepo.CountWaitingAtGateFiltered(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.VisitorEntryListResult{
		Entries: entries,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}

func (s *VisitorEntrySvc) ListExpectedGuests(ctx context.Context, filter models.ExpectedGuestFilter) (*models.VisitorEntryListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.FromAt.IsZero() || filter.ToAt.IsZero() {
		filter.FromAt, filter.ToAt = istDayRange(time.Now())
	}

	entries, err := s.entryRepo.ListExpectedGuests(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.entryRepo.CountExpectedGuestsFiltered(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.VisitorEntryListResult{
		Entries: entries,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}

func (s *VisitorEntrySvc) buildOnBehalfAudit(ctx context.Context, entry *models.VisitorEntry, guardUserID int64, reason *string) (map[string]any, error) {
	audit := map[string]any{
		"approved_on_behalf":   true,
		"approved_by_guard_id": guardUserID,
		"approved_at":          time.Now().UTC().Format(time.RFC3339),
	}

	pending, err := s.entryRepo.ListSocietyPending(ctx, models.VisitorPendingFilter{
		SocietyID: entry.SocietyID,
		Limit:     200,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range pending {
		if item.ID == entry.ID {
			if item.PrimaryResidentID != nil {
				audit["on_behalf_of_resident_id"] = *item.PrimaryResidentID
			}
			if item.PrimaryResidentName != nil {
				audit["on_behalf_of_resident_name"] = *item.PrimaryResidentName
			}
			break
		}
	}

	if reason != nil && *reason != "" {
		audit["on_behalf_reason"] = *reason
	}
	return audit, nil
}
