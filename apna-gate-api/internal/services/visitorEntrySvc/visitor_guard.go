package visitorentrysvc

import (
	"context"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorEntrySvc) GetGuardDeskBootstrap(ctx context.Context, societyID int64) (*models.GuardDeskBootstrapResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if societyID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}

	society, err := s.societyRepo.Get(ctx, models.GetSocietyFilter{ID: &societyID})
	if err != nil {
		return nil, err
	}
	if society == nil {
		return nil, ErrInvalidVisitorRequest
	}

	stats, err := s.entryRepo.GetStats(ctx, societyID)
	if err != nil {
		return nil, err
	}

	fromAt, toAt := istDayRange(time.Now())
	expectedGuestsCount, err := s.entryRepo.CountExpectedGuests(ctx, societyID, fromAt, toAt)
	if err != nil {
		return nil, err
	}

	waitingAtGateCount, err := s.entryRepo.CountWaitingAtGate(ctx, societyID)
	if err != nil {
		return nil, err
	}

	pending, err := s.entryRepo.ListSocietyPending(ctx, models.VisitorPendingFilter{
		SocietyID: societyID,
		Limit:     3,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	visitorSettings, err := s.settingSvc.GetSocietySettings(ctx, societyID)
	if err != nil {
		return nil, err
	}

	return &models.GuardDeskBootstrapResponse{
		Society:             society.ToResponse(),
		Stats:               stats,
		ExpectedGuestsCount: expectedGuestsCount,
		WaitingAtGateCount:  waitingAtGateCount,
		PendingPreview:      pending,
		VisitorSettings: &models.GuardDeskVisitorSettingsSummary{
			AllowGuardEntry:            visitorSettings.AllowGuardEntry,
			AllowGuardOnBehalfApproval: visitorSettings.AllowGuardOnBehalfApproval,
		},
	}, nil
}

func (s *VisitorEntrySvc) CheckIn(ctx context.Context, rawToken string, guardUserID int64) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	entry, err := s.ValidateQR(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if err := s.ensureStaffActor(ctx, entry.SocietyID, guardUserID); err != nil {
		return nil, err
	}
	var checkedIn *models.VisitorEntry
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		checkedIn, err = s.entryRepo.CheckIn(txCtx, entry.SocietyID, entry.ID, guardUserID)
		if err != nil {
			return err
		}
		if checkedIn == nil {
			return ErrVisitorInvalidState
		}
		copyVisitorSummary(checkedIn, entry)
		return s.recordEvents(txCtx, checkedIn, &guardUserID, models.VisitorEventTypeQRUsed, models.VisitorEventTypeCheckedIn)
	})
	if err != nil {
		return checkedIn, err
	}
	checkedIn = s.enrichVisitorSummary(ctx, checkedIn, entry)
	s.notifyVisitorCheckIn(checkedIn)
	return checkedIn, err
}

func (s *VisitorEntrySvc) CheckOut(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureStaffActor(ctx, societyID, guardUserID); err != nil {
		return nil, err
	}
	var checkedOut *models.VisitorEntry
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		checkedOut, err = s.entryRepo.CheckOut(txCtx, societyID, entryID, guardUserID)
		if err != nil {
			return err
		}
		if checkedOut == nil {
			return ErrVisitorInvalidState
		}
		return s.recordEvents(txCtx, checkedOut, &guardUserID, models.VisitorEventTypeCheckedOut)
	})
	if err != nil {
		return checkedOut, err
	}
	checkedOut = s.enrichVisitorSummary(ctx, checkedOut, nil)
	s.notifyVisitorCheckOut(checkedOut)
	return checkedOut, err
}

func (s *VisitorEntrySvc) ListSocietyPendingApprovals(ctx context.Context, filter models.VisitorPendingFilter) (*models.VisitorPendingListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	entries, err := s.entryRepo.ListSocietyPending(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.entryRepo.CountSocietyPending(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &models.VisitorPendingListResult{
		Entries: entries,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}
