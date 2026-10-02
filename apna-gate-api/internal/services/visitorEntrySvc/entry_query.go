package visitorentrysvc

import (
	"context"
	"sort"
	"time"

	"go-server/internal/models"
	service "go-server/internal/services"
	"go-server/pkg/logger"

	"go.uber.org/zap"
)

func (s *VisitorEntrySvc) GetEntryOptions(ctx context.Context, societyID int64) (*models.VisitorEntryOptionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	active := true
	occupied := string(models.FlatStatusOccupied)
	const entryOptionsFlatLimit int32 = 500
	filter := &models.FlatFilter{SocietyID: &societyID, Status: &occupied, IsActive: &active, Limit: entryOptionsFlatLimit}
	flats, err := s.flatRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.flatRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]models.VisitorEntryOptionsFlat, 0, len(flats))
	grouped := map[string][]models.VisitorEntryOptionsFlat{}
	for _, flat := range flats {
		item := models.VisitorEntryOptionsFlat{ID: flat.ID, Block: flat.Block, Floor: flat.Floor, FlatNumber: flat.FlatNumber}
		items = append(items, item)
		grouped[blockKey(flat.Block)] = append(grouped[blockKey(flat.Block)], item)
	}
	blocks := make([]models.VisitorEntryOptionsBlock, 0, len(grouped))
	for key, flats := range grouped {
		var block *string
		if key != "" {
			value := key
			block = &value
		}
		blocks = append(blocks, models.VisitorEntryOptionsBlock{Block: block, Flats: flats})
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blockKey(blocks[i].Block) < blockKey(blocks[j].Block)
	})
	return &models.VisitorEntryOptionsResponse{
		Purposes: []models.VisitorPurpose{models.VisitorPurposeGuest, models.VisitorPurposeDelivery, models.VisitorPurposeCab, models.VisitorPurposeService, models.VisitorPurposeMaintenance, models.VisitorPurposeStaff, models.VisitorPurposeOther},
		Blocks:   blocks,
		Flats:    items,
		HasMore:  total > int64(len(flats)),
	}, nil
}

// GetEntryForQRScan resolves a scan token for guard preview (validate endpoint).
// Returns the linked entry even when already checked in so repeat scans stay idempotent.
func (s *VisitorEntrySvc) GetEntryForQRScan(ctx context.Context, rawToken string) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	entry, err := s.resolveEntryByScanToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorQRInvalid
	}
	switch entry.Status {
	case models.VisitorStatusApproved:
		if entry.QRExpiresAt == nil {
			return nil, ErrVisitorQRUnavailable
		}
		if time.Now().After(*entry.QRExpiresAt) {
			return nil, ErrVisitorQRExpired
		}
		return entry, nil
	case models.VisitorStatusCheckedIn, models.VisitorStatusCheckedOut, models.VisitorStatusWaitingApproval:
		return entry, nil
	default:
		return nil, ErrVisitorInvalidState
	}
}

func (s *VisitorEntrySvc) GetEntryForQRScanForViewer(ctx context.Context, rawToken string, viewerUserID *int64) (*models.VisitorEntry, error) {
	entry, err := s.GetEntryForQRScan(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if viewerUserID != nil && *viewerUserID > 0 {
		if err := s.ensureStaffActor(ctx, entry.SocietyID, *viewerUserID); err == nil {
			return entry, nil
		}
	}
	return entry.ToScanPreview(), nil
}

func (s *VisitorEntrySvc) ValidateQR(ctx context.Context, rawToken string) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	entry, err := s.resolveEntryByScanToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorQRInvalid
	}
	if entry.Status == models.VisitorStatusCheckedIn {
		return nil, ErrVisitorAlreadyCheckedIn
	}
	if entry.Status != models.VisitorStatusApproved {
		return nil, ErrVisitorInvalidState
	}
	if entry.QRExpiresAt == nil {
		return nil, ErrVisitorQRUnavailable
	}
	if time.Now().After(*entry.QRExpiresAt) {
		return nil, ErrVisitorQRExpired
	}
	return entry, nil
}

func (s *VisitorEntrySvc) resolveEntryByScanToken(ctx context.Context, rawToken string) (*models.VisitorEntry, error) {
	tokenHash := hashToken(rawToken)

	entry, err := s.entryRepo.GetByQRHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		return entry, nil
	}
	return s.entryByInviteToken(ctx, tokenHash)
}

func (s *VisitorEntrySvc) entryByInviteToken(ctx context.Context, tokenHash string) (*models.VisitorEntry, error) {
	invite, err := s.inviteRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if invite == nil || invite.Status != models.VisitorInviteStatusUsed {
		return nil, nil
	}
	return s.entryRepo.GetByInviteID(ctx, invite.ID)
}

func (s *VisitorEntrySvc) GetEntry(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	entry, err := s.entryRepo.Get(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrVisitorEntryNotFound
	}
	logger.Info("Entry from svc", zap.Any("entry", entry))
	return entry, nil
}

func (s *VisitorEntrySvc) ListEntries(ctx context.Context, filter models.VisitorEntryFilter) ([]*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.entryRepo.List(ctx, filter)
}

func (s *VisitorEntrySvc) ListEntriesPaginated(ctx context.Context, filter models.VisitorEntryFilter) (*models.VisitorEntryListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	entries, err := s.entryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.entryRepo.Count(ctx, filter)
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

func (s *VisitorEntrySvc) GetEntryStats(ctx context.Context, societyID int64) (*models.VisitorEntryStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.entryRepo.GetStats(ctx, societyID)
}

func (s *VisitorEntrySvc) GetEntryStatsInRange(ctx context.Context, societyID int64, from, to time.Time) (*models.VisitorEntryStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	return s.entryRepo.GetStatsInRange(ctx, societyID, from, to)
}

func (s *VisitorEntrySvc) GetFlatVisitorContext(ctx context.Context, societyID int64, flatID int64) (*models.FlatVisitorContextResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	flat, err := s.flatRepo.Get(ctx, &models.FlatFilter{SocietyID: &societyID, ID: &flatID})
	if err != nil {
		return nil, err
	}
	if flat == nil {
		return nil, ErrVisitorFlatNotFound
	}

	societySettings, err := s.settingSvc.GetSocietySettings(ctx, societyID)
	if err != nil {
		return nil, err
	}
	if societySettings == nil {
		return nil, ErrVisitorSettingsNotFound
	}

	flatSettings, err := s.settingSvc.GetFlatSettings(ctx, societyID, flatID)
	if err != nil {
		return nil, err
	}

	totalResidents, err := s.residentRepo.CountActive(ctx, societyID, flatID)
	if err != nil {
		return nil, err
	}

	status := string(models.FlatResidentStatusActive)
	isPrimary := true
	primary, err := s.residentRepo.Get(ctx, &models.FlatResidentFilter{
		SocietyID: &societyID,
		FlatID:    &flatID,
		Status:    &status,
		IsPrimary: &isPrimary,
	})
	if err != nil {
		return nil, err
	}

	recent, err := s.entryRepo.ListRecentByFlat(ctx, societyID, flatID, 10)
	if err != nil {
		return nil, err
	}
	recentSummaries := make([]*models.FlatRecentVisitorSummary, 0, len(recent))
	for _, entry := range recent {
		name := ""
		if entry.Visitor != nil {
			name = entry.Visitor.FullName
		}
		recentSummaries = append(recentSummaries, &models.FlatRecentVisitorSummary{
			EntryID:   entry.ID,
			FullName:  name,
			Purpose:   entry.Purpose,
			Status:    entry.Status,
			VisitedOn: entry.CreatedAt,
		})
	}

	var primaryResident *models.FlatVisitorContextResident
	if primary != nil {
		name := ""
		if primary.UserName != nil {
			name = *primary.UserName
		}
		primaryResident = &models.FlatVisitorContextResident{
			ID:       primary.ID,
			UserID:   primary.UserID,
			FullName: name,
		}
	}

	return &models.FlatVisitorContextResponse{
		OccupancyStatus:     flat.Status,
		PrimaryResident:     primaryResident,
		TotalResidents:      totalResidents,
		InheritsSocietyMode: societySettings.ApprovalMode == models.VisitorApprovalModeHybrid,
		SocietyApprovalMode: societySettings.ApprovalMode,
		VisitorSettings:     flatSettings,
		RecentVisitors:      recentSummaries,
	}, nil
}

func (s *VisitorEntrySvc) GetMemberVisitorApprovalStats(ctx context.Context, societyID int64, memberID int64) (*models.MemberVisitorApprovalStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	member, err := s.memberRepo.Get(ctx, models.GetSocietyMemberFilter{ID: &memberID, SocietyID: &societyID})
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrVisitorForbidden
	}
	return s.entryRepo.CountMemberApprovals(ctx, societyID, member.UserID)
}

func (s *VisitorEntrySvc) ListPendingApprovals(ctx context.Context, societyID int64, flatID int64, actorUserID int64) ([]*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureApprovalActor(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	return s.entryRepo.ListPending(ctx, societyID, flatID)
}

func (s *VisitorEntrySvc) ListVisitorInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter models.VisitorInviteHistoryFilter) (*models.VisitorInviteHistoryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatResident(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	filter = normalizeVisitorInviteHistoryFilter(filter)
	items, err := s.inviteRepo.ListHistory(ctx, societyID, flatID, filter, s.shortLinkResponseFromCode)
	if err != nil {
		return nil, err
	}
	total, err := s.inviteRepo.CountHistory(ctx, societyID, flatID, filter)
	if err != nil {
		return nil, err
	}
	return &models.VisitorInviteHistoryResult{
		Invites: items,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	}, nil
}

func (s *VisitorEntrySvc) GetVisitorInviteHistoryForActor(ctx context.Context, societyID int64, flatID int64, inviteID int64, actorUserID int64) (*models.VisitorInviteHistoryItem, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatResident(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	if inviteID <= 0 {
		return nil, ErrInvalidVisitorRequest
	}
	item, err := s.inviteRepo.GetHistoryByID(ctx, societyID, flatID, inviteID, s.shortLinkResponseFromCode)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrVisitorInviteNotFound
	}
	return item, nil
}

func (s *VisitorEntrySvc) ListFlatEntriesForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64, filter models.VisitorEntryFilter) (*models.VisitorEntryListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatResident(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}

	filter.SocietyID = societyID
	filter.FlatID = &flatID
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	return s.ListEntriesPaginated(ctx, filter)
}

func (s *VisitorEntrySvc) GetFlatEntryForActor(ctx context.Context, societyID int64, flatID int64, entryID int64, actorUserID int64) (*models.VisitorEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatResident(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	entry, err := s.GetEntry(ctx, societyID, entryID)
	if err != nil {
		return nil, err
	}
	if entry.FlatID != flatID {
		return nil, ErrVisitorEntryNotFound
	}
	return entry, nil
}

func (s *VisitorEntrySvc) GetFlatVisitorContextForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) (*models.FlatVisitorContextResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatVisitorContextActor(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}

	if err := s.flatAuthz.CanManageFlatVisitors(ctx, societyID, flatID, actorUserID); err == nil {
		if err := s.settingSvc.EnsureDefaultFlatSettingsIfMissing(ctx, societyID, flatID, actorUserID); err != nil {
			return nil, err
		}
	}

	return s.GetFlatVisitorContext(ctx, societyID, flatID)
}

func (s *VisitorEntrySvc) ListEvents(ctx context.Context, societyID int64, entryID int64) ([]*models.VisitorEntryEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	if _, err := s.GetEntry(ctx, societyID, entryID); err != nil {
		return nil, err
	}
	return s.eventRepo.List(ctx, societyID, entryID)
}
