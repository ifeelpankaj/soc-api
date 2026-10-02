package repository

import (
	"context"
	"encoding/json"
	"errors"
	"go-server/internal/repositories/contracts"
	"time"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VisitorRepository = contracts.VisitorRepository

type VisitorInviteRepository = contracts.VisitorInviteRepository

type VisitorEntryRepository = contracts.VisitorEntryRepository

type VisitorEntryEventRepository = contracts.VisitorEntryEventRepository

type visitorRepository struct {
	db *database.Database
}

type visitorInviteRepository struct {
	db *database.Database
}

type visitorEntryRepository struct {
	db *database.Database
}

type visitorEntryEventRepository struct {
	db *database.Database
}

func NewVisitorRepository(db *database.Database) VisitorRepository {
	return &visitorRepository{db: db}
}

func NewVisitorInviteRepository(db *database.Database) VisitorInviteRepository {
	return &visitorInviteRepository{db: db}
}

func NewVisitorEntryRepository(db *database.Database) VisitorEntryRepository {
	return &visitorEntryRepository{db: db}
}

func NewVisitorEntryEventRepository(db *database.Database) VisitorEntryEventRepository {
	return &visitorEntryEventRepository{db: db}
}

func (r *visitorRepository) Create(ctx context.Context, req contracts.VisitorFormInput) (*models.Visitor, error) {
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CreateVisitor(ctx, db.CreateVisitorParams{
		FullName: req.FullName, PhoneNumber: req.PhoneNumber, Email: req.Email, PhotoUrl: req.PhotoURL, Metadata: metadata,
	})
	return visitorFromDBNoRows(row, err)
}

func (r *visitorRepository) Get(ctx context.Context, id int64) (*models.Visitor, error) {
	row, err := GetQueries(ctx, r.db).GetVisitor(ctx, id)
	return visitorFromDBNoRows(row, err)
}

func (r *visitorRepository) UpdateProfile(ctx context.Context, visitorID int64, req contracts.UpdateGuardVisitorEntryInput) (*models.Visitor, error) {
	row, err := GetQueries(ctx, r.db).UpdateVisitorProfile(ctx, db.UpdateVisitorProfileParams{
		ID:          visitorID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		PhotoUrl:    req.PhotoURL,
	})
	return visitorFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) Create(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, tokenHash string, expiresAt time.Time, actorUserID int64) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).CreateVisitorInvite(ctx, db.CreateVisitorInviteParams{
		SocietyID: societyID,
		FlatID:    flatID,
		CreatedBy: actorUserID,
		Purpose:   db.VisitorPurpose(purpose),
		TokenHash: tokenHash,
		ExpiresAt: visitorTimeToPgTimestamptz(expiresAt),
		Metadata:  []byte("{}"),
	})
	return visitorInviteFromDBNoRows(row, persistenceError(err))
}

func (r *visitorInviteRepository) GetByID(ctx context.Context, societyID int64, inviteID int64) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorInviteByID(ctx, db.GetVisitorInviteByIDParams{ID: inviteID, SocietyID: societyID})
	return visitorInviteFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) GetHistoryByID(ctx context.Context, societyID int64, flatID int64, inviteID int64, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) (*models.VisitorInviteHistoryItem, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorInviteHistoryByID(ctx, db.GetVisitorInviteHistoryByIDParams{SocietyID: societyID, FlatID: flatID, ID: inviteID})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorInviteHistoryByIDFromDB(row, linkBuilder), nil
}

func (r *visitorInviteRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorInviteByTokenHash(ctx, tokenHash)
	return visitorInviteFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) ListHistory(ctx context.Context, societyID int64, flatID int64, filter models.VisitorInviteHistoryFilter, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) ([]*models.VisitorInviteHistoryItem, error) {
	var status *db.VisitorInviteStatus
	if filter.Status != nil {
		value := db.VisitorInviteStatus(*filter.Status)
		status = &value
	}
	rows, err := GetQueries(ctx, r.db).ListVisitorInviteHistory(ctx, db.ListVisitorInviteHistoryParams{
		SocietyID: societyID,
		FlatID:    flatID,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
		Status:    status,
		Search:    filter.Search,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorInviteHistoryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorInviteHistoryFromDB(row, linkBuilder))
	}
	return items, nil
}

func (r *visitorInviteRepository) CountHistory(ctx context.Context, societyID int64, flatID int64, filter models.VisitorInviteHistoryFilter) (int64, error) {
	var status *db.VisitorInviteStatus
	if filter.Status != nil {
		value := db.VisitorInviteStatus(*filter.Status)
		status = &value
	}
	return GetQueries(ctx, r.db).CountVisitorInviteHistory(ctx, db.CountVisitorInviteHistoryParams{
		SocietyID: societyID,
		FlatID:    flatID,
		Status:    status,
		Search:    filter.Search,
	})
}

func (r *visitorInviteRepository) MarkUsed(ctx context.Context, inviteID int64) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).MarkVisitorInviteUsed(ctx, inviteID)
	return visitorInviteFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) GetForUpdate(ctx context.Context, inviteID int64) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorInviteForUpdate(ctx, inviteID)
	return visitorInviteFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) Cancel(ctx context.Context, societyID int64, inviteID int64) (*models.VisitorInvite, error) {
	row, err := GetQueries(ctx, r.db).CancelVisitorInvite(ctx, db.CancelVisitorInviteParams{ID: inviteID, SocietyID: societyID})
	return visitorInviteFromDBNoRows(row, err)
}

func (r *visitorInviteRepository) ExpireOld(ctx context.Context) error {
	return GetQueries(ctx, r.db).ExpireOldVisitorInvites(ctx)
}

func (r *visitorEntryRepository) Create(ctx context.Context, req contracts.VisitorFormInput, societyID int64, flatID *int64, visitorID int64, inviteID *int64, source models.VisitorEntrySource, purpose models.VisitorPurpose, status models.VisitorStatus, actorUserID *int64, guardUserID *int64, qrHash *string, qrExpiresAt *time.Time) (*models.VisitorEntry, error) {
	companionDetails, err := json.Marshal(req.CompanionDetails)
	if err != nil {
		return nil, err
	}
	if string(companionDetails) == "null" {
		companionDetails = []byte("[]")
	}
	metadata, err := jsonMap(req.Metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CreateVisitorEntry(ctx, db.CreateVisitorEntryParams{
		SocietyID:          societyID,
		FlatID:             flatID,
		VisitorID:          visitorID,
		Source:             db.VisitorSource(source),
		Purpose:            db.VisitorPurpose(purpose),
		Status:             db.VisitorStatus(status),
		CompanionsCount:    req.CompanionsCount,
		InviteID:           inviteID,
		VehicleNumber:      req.VehicleNumber,
		VehicleType:        dbVisitorVehicleTypePtr(req.VehicleType),
		CompanionDetails:   companionDetails,
		ExpectedAt:         timePtrToPgTimestamptz(req.ExpectedAt),
		ExpectedCheckoutAt: timePtrToPgTimestamptz(req.ExpectedCheckoutAt),
		ApprovedBy:         approvedByForCreate(status, actorUserID),
		ApprovedAt:         approvedAtForCreate(status),
		DeliveryPartner:    req.DeliveryPartner,
		ServiceProvider:    req.ServiceProvider,
		HandledByGuardID:   guardUserID,
		CreatedBy:          actorUserID,
		QrTokenHash:        qrHash,
		QrExpiresAt:        timePtrToPgTimestamptz(qrExpiresAt),
		Notes:              req.Notes, Metadata: metadata,
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) Get(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntry(ctx, db.GetVisitorEntryParams{ID: entryID, SocietyID: societyID})
	return visitorEntryFromGetNoRows(row, err)
}

func (r *visitorEntryRepository) GetForUpdate(ctx context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntryForUpdate(ctx, db.GetVisitorEntryForUpdateParams{ID: entryID, SocietyID: societyID})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) GetByQRHash(ctx context.Context, qrHash string) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntryByQRHash(ctx, &qrHash)
	return visitorEntryFromQRNoRows(row, err)
}

func (r *visitorEntryRepository) GetByInviteID(ctx context.Context, inviteID int64) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntryByInviteID(ctx, &inviteID)
	return visitorEntryFromInviteNoRows(row, err)
}

func (r *visitorEntryRepository) List(ctx context.Context, filter models.VisitorEntryFilter) ([]*models.VisitorEntry, error) {
	params := visitorEntryListParams(filter)
	rows, err := GetQueries(ctx, r.db).ListVisitorEntries(ctx, params)
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryFromList(row))
	}
	return items, nil
}

func (r *visitorEntryRepository) Count(ctx context.Context, filter models.VisitorEntryFilter) (int64, error) {
	params := visitorEntryCountParams(filter)
	return GetQueries(ctx, r.db).CountVisitorEntries(ctx, params)
}

func (r *visitorEntryRepository) ListSocietyPending(ctx context.Context, filter models.VisitorPendingFilter) ([]*models.VisitorPendingEntry, error) {
	rows, err := GetQueries(ctx, r.db).ListSocietyPendingVisitorApprovals(ctx, db.ListSocietyPendingVisitorApprovalsParams{
		SocietyID: filter.SocietyID,
		FlatID:    filter.FlatID,
		Block:     filter.Block,
		Limit:     normalizeVisitorLimit(filter.Limit),
		Offset:    normalizeOffset(filter.Offset),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorPendingEntry, 0, len(rows))
	for _, row := range rows {
		entry := visitorEntryFromSocietyPending(row)
		items = append(items, &models.VisitorPendingEntry{
			VisitorEntry:        entry,
			WaitingSince:        entry.CreatedAt,
			PrimaryResidentName: row.PrimaryResidentName,
			PrimaryResidentID:   row.PrimaryResidentID,
		})
	}
	return items, nil
}

func (r *visitorEntryRepository) CountSocietyPending(ctx context.Context, filter models.VisitorPendingFilter) (int64, error) {
	return GetQueries(ctx, r.db).CountSocietyPendingVisitorApprovals(ctx, db.CountSocietyPendingVisitorApprovalsParams{
		SocietyID: filter.SocietyID,
		FlatID:    filter.FlatID,
		Block:     filter.Block,
	})
}

func (r *visitorEntryRepository) ListRecentByFlat(ctx context.Context, societyID int64, flatID int64, limit int32) ([]*models.VisitorEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := GetQueries(ctx, r.db).ListRecentVisitorEntriesByFlat(ctx, db.ListRecentVisitorEntriesByFlatParams{
		SocietyID: societyID,
		FlatID:    flatIDPtr(flatID),
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryFromRecent(row))
	}
	return items, nil
}

func (r *visitorEntryRepository) GetStats(ctx context.Context, societyID int64) (*models.VisitorEntryStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntryStats(ctx, societyID)
	if err != nil {
		return nil, err
	}
	return &models.VisitorEntryStatsResponse{
		TodayVisitors:    row.TodayVisitors,
		VisitorsInside:   row.VisitorsInside,
		PendingApprovals: row.PendingApprovals,
		CheckedOutToday:  row.CheckedOutToday,
		RejectedToday:    row.RejectedToday,
		AutoClosedToday:  row.AutoClosedToday,
	}, nil
}

func (r *visitorEntryRepository) GetDailyStatsCreated(ctx context.Context, societyID int64, days int32) ([]models.VisitorDailyCountResponse, error) {
	if days <= 0 {
		days = 7
	}
	rows, err := GetQueries(ctx, r.db).GetVisitorEntryDailyStatsCreated(ctx, db.GetVisitorEntryDailyStatsCreatedParams{
		SocietyID: societyID,
		Days:      days,
	})
	if err != nil {
		return nil, err
	}
	items := make([]models.VisitorDailyCountResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.VisitorDailyCountResponse{
			Date:  row.StatDate,
			Count: row.Count,
		})
	}
	return items, nil
}

func (r *visitorEntryRepository) GetStatsInRange(ctx context.Context, societyID int64, from, to time.Time) (*models.VisitorEntryStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).GetVisitorEntryStatsInRange(ctx, db.GetVisitorEntryStatsInRangeParams{
		SocietyID: societyID,
		EventFrom: timePtrToPgTimestamptz(&from),
		EventTo:   timePtrToPgTimestamptz(&to),
	})
	if err != nil {
		return nil, err
	}
	return &models.VisitorEntryStatsResponse{
		VisitorsInside:    row.VisitorsInside,
		PendingApprovals:  row.PendingApprovals,
		CheckedOutInRange: row.CheckedOutInRange,
	}, nil
}

func (r *visitorEntryRepository) CountWaitingAtGate(ctx context.Context, societyID int64) (int64, error) {
	return GetQueries(ctx, r.db).CountWaitingAtGateVisitorEntries(ctx, societyID)
}

func (r *visitorEntryRepository) ListWaitingAtGate(ctx context.Context, filter models.WaitingAtGateFilter) ([]*models.VisitorEntry, error) {
	rows, err := GetQueries(ctx, r.db).ListWaitingAtGateVisitorEntries(ctx, db.ListWaitingAtGateVisitorEntriesParams{
		SocietyID: filter.SocietyID,
		Search:    filter.Search,
		Limit:     normalizeVisitorLimit(filter.Limit),
		Offset:    normalizeOffset(filter.Offset),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryFromWaitingAtGate(row))
	}
	return items, nil
}

func (r *visitorEntryRepository) CountWaitingAtGateFiltered(ctx context.Context, filter models.WaitingAtGateFilter) (int64, error) {
	return GetQueries(ctx, r.db).CountWaitingAtGateVisitorEntriesFiltered(ctx, db.CountWaitingAtGateVisitorEntriesFilteredParams{
		SocietyID: filter.SocietyID,
		Search:    filter.Search,
	})
}

func (r *visitorEntryRepository) CountExpectedGuests(ctx context.Context, societyID int64, fromAt, toAt time.Time) (int64, error) {
	return GetQueries(ctx, r.db).CountExpectedGuestEntries(ctx, db.CountExpectedGuestEntriesParams{
		SocietyID: societyID,
		FromAt:    visitorTimeToPgTimestamptz(fromAt),
		ToAt:      visitorTimeToPgTimestamptz(toAt),
	})
}

func (r *visitorEntryRepository) ListExpectedGuests(ctx context.Context, filter models.ExpectedGuestFilter) ([]*models.VisitorEntry, error) {
	rows, err := GetQueries(ctx, r.db).ListExpectedGuestEntries(ctx, db.ListExpectedGuestEntriesParams{
		SocietyID: filter.SocietyID,
		FromAt:    visitorTimeToPgTimestamptz(filter.FromAt),
		ToAt:      visitorTimeToPgTimestamptz(filter.ToAt),
		Search:    filter.Search,
		Limit:     normalizeVisitorLimit(filter.Limit),
		Offset:    normalizeOffset(filter.Offset),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryFromExpectedGuest(row))
	}
	return items, nil
}

func (r *visitorEntryRepository) CountExpectedGuestsFiltered(ctx context.Context, filter models.ExpectedGuestFilter) (int64, error) {
	return GetQueries(ctx, r.db).CountExpectedGuestEntriesFiltered(ctx, db.CountExpectedGuestEntriesFilteredParams{
		SocietyID: filter.SocietyID,
		FromAt:    visitorTimeToPgTimestamptz(filter.FromAt),
		ToAt:      visitorTimeToPgTimestamptz(filter.ToAt),
		Search:    filter.Search,
	})
}

func (r *visitorEntryRepository) MergeMetadata(ctx context.Context, societyID int64, entryID int64, metadata map[string]any) (*models.VisitorEntry, error) {
	raw, err := jsonMap(metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).MergeVisitorEntryMetadata(ctx, db.MergeVisitorEntryMetadataParams{
		ID: entryID, SocietyID: societyID, MetadataPatch: raw,
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) CountMemberApprovals(ctx context.Context, societyID int64, userID int64) (*models.MemberVisitorApprovalStatsResponse, error) {
	row, err := GetQueries(ctx, r.db).CountMemberVisitorApprovals(ctx, db.CountMemberVisitorApprovalsParams{
		SocietyID:  societyID,
		ApprovedBy: &userID,
	})
	if err != nil {
		return nil, err
	}
	return &models.MemberVisitorApprovalStatsResponse{
		ApprovedCount: row.ApprovedCount,
		RejectedCount: row.RejectedCount,
	}, nil
}

func (r *visitorEntryRepository) ListPending(ctx context.Context, societyID int64, flatID int64) ([]*models.VisitorEntry, error) {
	rows, err := GetQueries(ctx, r.db).ListPendingVisitorApprovals(ctx, db.ListPendingVisitorApprovalsParams{
		SocietyID: societyID,
		FlatID:    flatIDPtr(flatID),
		Limit:     100,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryFromPending(row))
	}
	return items, nil
}

func (r *visitorEntryRepository) Approve(ctx context.Context, societyID int64, entryID int64, actorUserID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).ApproveVisitorEntry(ctx, db.ApproveVisitorEntryParams{
		ID: entryID, SocietyID: societyID, ApprovedBy: &actorUserID, QrTokenHash: &qrHash, QrExpiresAt: visitorTimeToPgTimestamptz(qrExpiresAt),
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) Reject(ctx context.Context, societyID int64, entryID int64, actorUserID int64, reason string) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).RejectVisitorEntry(ctx, db.RejectVisitorEntryParams{
		ID: entryID, SocietyID: societyID, RejectedBy: &actorUserID, RejectionReason: &reason,
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) GenerateQR(ctx context.Context, societyID int64, entryID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).GenerateVisitorEntryQR(ctx, db.GenerateVisitorEntryQRParams{
		ID: entryID, SocietyID: societyID, QrTokenHash: &qrHash, QrExpiresAt: visitorTimeToPgTimestamptz(qrExpiresAt),
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) CheckIn(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).CheckInVisitorEntry(ctx, db.CheckInVisitorEntryParams{ID: entryID, SocietyID: societyID, HandledByGuardID: &guardUserID})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) CheckOut(ctx context.Context, societyID int64, entryID int64, guardUserID int64) (*models.VisitorEntry, error) {
	row, err := GetQueries(ctx, r.db).CheckOutVisitorEntry(ctx, db.CheckOutVisitorEntryParams{ID: entryID, SocietyID: societyID, HandledByGuardID: &guardUserID})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) UpdateDetails(ctx context.Context, societyID int64, entryID int64, req contracts.UpdateGuardVisitorEntryInput) (*models.VisitorEntry, error) {
	var companionDetails []byte
	if req.CompanionDetails != nil {
		raw, err := json.Marshal(req.CompanionDetails)
		if err != nil {
			return nil, err
		}
		if string(raw) == "null" {
			raw = []byte("[]")
		}
		companionDetails = raw
	}
	row, err := GetQueries(ctx, r.db).UpdateVisitorEntryDetails(ctx, db.UpdateVisitorEntryDetailsParams{
		ID:               entryID,
		SocietyID:        societyID,
		FlatID:           req.FlatID,
		VehicleNumber:    req.VehicleNumber,
		VehicleType:      dbVisitorVehicleTypePtr(req.VehicleType),
		CompanionsCount:  req.CompanionsCount,
		CompanionDetails: companionDetails,
		Notes:            req.Notes,
	})
	return visitorEntryFromDBNoRows(row, err)
}

func (r *visitorEntryRepository) AutoCloseExpired(ctx context.Context) error {
	return GetQueries(ctx, r.db).AutoCloseExpiredVisitorEntries(ctx)
}

func (r *visitorEntryRepository) ExpireStaleEntries(ctx context.Context) error {
	queries := GetQueries(ctx, r.db)
	if err := queries.ExpireStaleWaitingApprovalEntries(ctx); err != nil {
		return err
	}
	return queries.ExpireStaleApprovedEntries(ctx)
}

func (r *visitorEntryEventRepository) Create(ctx context.Context, entryID int64, societyID int64, actorUserID *int64, eventType models.VisitorEventType, message *string, metadata map[string]any) (*models.VisitorEntryEvent, error) {
	rawMetadata, err := jsonMap(metadata)
	if err != nil {
		return nil, err
	}
	row, err := GetQueries(ctx, r.db).CreateVisitorEntryEvent(ctx, db.CreateVisitorEntryEventParams{
		VisitorEntryID: entryID, SocietyID: societyID, ActorUserID: actorUserID,
		EventType: db.VisitorEventType(eventType), Message: message, Metadata: rawMetadata,
	})
	return visitorEntryEventFromDBNoRows(row, err)
}

func (r *visitorEntryEventRepository) List(ctx context.Context, societyID int64, entryID int64) ([]*models.VisitorEntryEvent, error) {
	rows, err := GetQueries(ctx, r.db).ListVisitorEntryEvents(ctx, db.ListVisitorEntryEventsParams{VisitorEntryID: entryID, SocietyID: societyID})
	if err != nil {
		return nil, err
	}
	items := make([]*models.VisitorEntryEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, visitorEntryEventFromDB(row))
	}
	return items, nil
}

func visitorFromDBNoRows(row db.Visitor, err error) (*models.Visitor, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorFromDB(row), nil
}

func visitorFromDB(row db.Visitor) *models.Visitor {
	return &models.Visitor{
		ID: row.ID, FullName: row.FullName, PhoneNumber: row.PhoneNumber, Email: row.Email, PhotoURL: row.PhotoUrl,
		Metadata: metadataMap(row.Metadata), CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func visitorInviteFromDBNoRows(row db.VisitorInvite, err error) (*models.VisitorInvite, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorInviteFromDB(row), nil
}

func visitorInviteFromDB(row db.VisitorInvite) *models.VisitorInvite {
	return &models.VisitorInvite{
		ID:        row.ID,
		SocietyID: row.SocietyID,
		FlatID:    row.FlatID,
		CreatedBy: row.CreatedBy,
		Purpose:   models.VisitorPurpose(row.Purpose),
		Status:    models.VisitorInviteStatus(row.Status),
		ExpiresAt: pgTimestamptzToTime(row.ExpiresAt),
		UsedAt:    pgTimestamptzToTimePtr(row.UsedAt),
		Metadata:  metadataMap(row.Metadata),
		CreatedAt: pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func visitorInviteHistoryFromDB(row db.ListVisitorInviteHistoryRow, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) *models.VisitorInviteHistoryItem {
	item := &models.VisitorInviteHistoryItem{
		ID:        row.ID,
		SocietyID: row.SocietyID,
		FlatID:    row.FlatID,
		CreatedBy: row.CreatedBy,
		Purpose:   models.VisitorPurpose(row.Purpose),
		Status:    models.VisitorInviteStatus(row.Status),
		ExpiresAt: pgTimestamptzToTime(row.ExpiresAt),
		UsedAt:    pgTimestamptzToTimePtr(row.UsedAt),
		Metadata:  metadataMap(row.Metadata),
		CreatedAt: pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
	if row.ShortCode != nil && linkBuilder != nil {
		item.ShortLink = linkBuilder(*row.ShortCode, pgTimestamptzToTimePtr(row.ShortLinkExpiresAt))
	}
	if row.VisitorFullName != nil {
		item.Visitor = &models.VisitorSummary{
			FullName:    *row.VisitorFullName,
			PhoneNumber: row.VisitorPhoneNumber,
			Email:       row.VisitorEmail,
			PhotoURL:    row.VisitorPhotoUrl,
		}
	}
	if row.EntryID != nil && row.EntryStatus != nil {
		item.Entry = &models.VisitorInviteEntry{
			ID:           *row.EntryID,
			Status:       models.VisitorStatus(*row.EntryStatus),
			CheckedInAt:  pgTimestamptzToTimePtr(row.CheckedInAt),
			CheckedOutAt: pgTimestamptzToTimePtr(row.CheckedOutAt),
			QRExpiresAt:  pgTimestamptzToTimePtr(row.QrExpiresAt),
			QRUsedAt:     pgTimestamptzToTimePtr(row.QrUsedAt),
		}
	}
	return item
}

func visitorInviteHistoryByIDFromDB(row db.GetVisitorInviteHistoryByIDRow, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) *models.VisitorInviteHistoryItem {
	return visitorInviteHistoryFromDB(db.ListVisitorInviteHistoryRow(row), linkBuilder)
}

func visitorEntryFromDBNoRows(row db.VisitorEntry, err error) (*models.VisitorEntry, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorEntryFromDB(row), nil
}

func visitorEntryFromDB(row db.VisitorEntry) *models.VisitorEntry {
	return visitorEntryFromParts(
		row.ID,
		row.SocietyID,
		flatIDValue(row.FlatID),
		row.VisitorID,
		row.InviteID,
		string(row.Source),
		string(row.Purpose),
		string(row.Status),
		row.VehicleNumber,
		row.VehicleType,
		row.CompanionsCount,
		row.CompanionDetails,
		row.ExpectedAt,
		row.ExpectedCheckoutAt,
		row.CheckedInAt,
		row.CheckedOutAt,
		row.AutoClosedAt,
		row.ApprovedBy,
		row.ApprovedAt,
		row.DeliveryPartner,
		row.ServiceProvider,
		row.RejectedBy,
		row.HandledByGuardID,
		row.CreatedBy,
		row.QrExpiresAt,
		row.QrUsedAt,
		row.Notes,
		row.RejectionReason,
		row.Metadata,
		row.CreatedAt,
		row.UpdatedAt,
		"",
		nil,
		nil,
		nil,
		"",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil)
}

func visitorEntryFromGetNoRows(row db.GetVisitorEntryRow, err error) (*models.VisitorEntry, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	entry := visitorEntryFromGet(row)
	return entry, nil
}

func visitorEntryFromQRNoRows(row db.GetVisitorEntryByQRHashRow, err error) (*models.VisitorEntry, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber), nil
}

func visitorEntryFromInviteNoRows(row db.GetVisitorEntryByInviteIDRow, err error) (*models.VisitorEntry, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber), nil
}

func visitorEntryFromGet(row db.GetVisitorEntryRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, // Approver
		row.ApprovedByFullName,
		row.ApprovedByPhoneNumber,

		// Guard
		row.HandledByGuardFullName,
		row.HandledByGuardPhoneNumber)
}

func visitorEntryFromList(row db.ListVisitorEntriesRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryFromRecent(row db.ListRecentVisitorEntriesByFlatRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryFromPending(row db.ListPendingVisitorApprovalsRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryFromSocietyPending(row db.ListSocietyPendingVisitorApprovalsRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryFromWaitingAtGate(row db.ListWaitingAtGateVisitorEntriesRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryFromExpectedGuest(row db.ListExpectedGuestEntriesRow) *models.VisitorEntry {
	return visitorEntryFromParts(row.ID, row.SocietyID, flatIDValue(row.FlatID), row.VisitorID, row.InviteID, string(row.Source), string(row.Purpose), string(row.Status), row.VehicleNumber, row.VehicleType, row.CompanionsCount, row.CompanionDetails, row.ExpectedAt, row.ExpectedCheckoutAt, row.CheckedInAt, row.CheckedOutAt, row.AutoClosedAt, row.ApprovedBy, row.ApprovedAt, row.DeliveryPartner, row.ServiceProvider, row.RejectedBy, row.HandledByGuardID, row.CreatedBy, row.QrExpiresAt, row.QrUsedAt, row.Notes, row.RejectionReason, row.Metadata, row.CreatedAt, row.UpdatedAt, row.VisitorFullName, row.VisitorPhoneNumber, row.VisitorEmail, row.VisitorPhotoUrl, flatNumberValue(row.FlatNumber), row.Block, row.Floor, row.ApprovedByFullName, row.ApprovedByPhoneNumber, row.HandledByGuardFullName, row.HandledByGuardPhoneNumber)
}

func visitorEntryListParams(filter models.VisitorEntryFilter) db.ListVisitorEntriesParams {
	return db.ListVisitorEntriesParams{
		SocietyID:   filter.SocietyID,
		FlatID:      filter.FlatID,
		Status:      dbVisitorStatusPtr(filter.Status),
		Source:      dbVisitorSourcePtr(filter.Source),
		Purpose:     dbVisitorPurposePtr(filter.Purpose),
		Block:       filter.Block,
		CreatedFrom: timePtrToPgTimestamptz(filter.CreatedFrom),
		CreatedTo:   timePtrToPgTimestamptz(filter.CreatedTo),
		Event:       visitorEntryEventPtr(filter.Event),
		EventFrom:   timePtrToPgTimestamptz(filter.EventFrom),
		EventTo:     timePtrToPgTimestamptz(filter.EventTo),
		Search:      filter.Search,
		Limit:       normalizeVisitorLimit(filter.Limit),
		Offset:      normalizeOffset(filter.Offset),
	}
}

func visitorEntryCountParams(filter models.VisitorEntryFilter) db.CountVisitorEntriesParams {
	return db.CountVisitorEntriesParams{
		SocietyID:   filter.SocietyID,
		FlatID:      filter.FlatID,
		Status:      dbVisitorStatusPtr(filter.Status),
		Source:      dbVisitorSourcePtr(filter.Source),
		Purpose:     dbVisitorPurposePtr(filter.Purpose),
		Block:       filter.Block,
		CreatedFrom: timePtrToPgTimestamptz(filter.CreatedFrom),
		CreatedTo:   timePtrToPgTimestamptz(filter.CreatedTo),
		Event:       visitorEntryEventPtr(filter.Event),
		EventFrom:   timePtrToPgTimestamptz(filter.EventFrom),
		EventTo:     timePtrToPgTimestamptz(filter.EventTo),
		Search:      filter.Search,
	}
}

func visitorEntryEventPtr(event *models.VisitorEntryEventFilter) *string {
	if event == nil {
		return nil
	}
	value := string(*event)
	return &value
}

func visitorEntryFromParts(
	id, societyID, flatID, visitorID int64,
	inviteID *int64,
	source, purpose, status string,
	vehicleNumber *string,
	vehicleType *db.VisitorVehicleType,
	companionsCount int32,
	companionDetails []byte,
	expectedAt, expectedCheckoutAt, checkedInAt, checkedOutAt, autoClosedAt pgtype.Timestamptz,
	approvedBy *int64,
	approvedAt pgtype.Timestamptz,
	deliveryPartner, serviceProvider *string,
	rejectedBy, handledByGuardID, createdBy *int64,
	qrExpiresAt, qrUsedAt pgtype.Timestamptz,
	notes, rejectionReason *string,
	metadata []byte,
	createdAt, updatedAt pgtype.Timestamptz,
	visitorFullName string,
	visitorPhone, visitorEmail, visitorPhoto *string,
	flatNumber string,
	block, floor *string,
	actorDetails ...*string,
) *models.VisitorEntry {
	var approvedByFullName, approvedByPhoneNumber, handledByGuardFullName, handledByGuardPhoneNumber *string
	if len(actorDetails) > 0 {
		approvedByFullName = actorDetails[0]
	}
	if len(actorDetails) > 1 {
		approvedByPhoneNumber = actorDetails[1]
	}
	if len(actorDetails) > 2 {
		handledByGuardFullName = actorDetails[2]
	}
	if len(actorDetails) > 3 {
		handledByGuardPhoneNumber = actorDetails[3]
	}

	return &models.VisitorEntry{
		ID:                        id,
		SocietyID:                 societyID,
		FlatID:                    flatID,
		VisitorID:                 visitorID,
		InviteID:                  inviteID,
		Source:                    models.VisitorEntrySource(source),
		Purpose:                   models.VisitorPurpose(purpose),
		Status:                    models.VisitorStatus(status),
		VehicleNumber:             vehicleNumber,
		VehicleType:               visitorVehicleTypeFromDB(vehicleType),
		CompanionsCount:           companionsCount,
		CompanionDetails:          companionDetailsFromJSON(companionDetails),
		ExpectedAt:                pgTimestamptzToTimePtr(expectedAt),
		ExpectedCheckoutAt:        pgTimestamptzToTimePtr(expectedCheckoutAt),
		CheckedInAt:               pgTimestamptzToTimePtr(checkedInAt),
		CheckedOutAt:              pgTimestamptzToTimePtr(checkedOutAt),
		AutoClosedAt:              pgTimestamptzToTimePtr(autoClosedAt),
		ApprovedBy:                approvedBy,
		ApproverName:              approvedByFullName,
		ApproverPhoneNumber:       approvedByPhoneNumber,
		ApprovedAt:                pgTimestamptzToTimePtr(approvedAt),
		DeliveryPartner:           deliveryPartner,
		ServiceProvider:           serviceProvider,
		RejectedBy:                rejectedBy,
		HandledByGuardID:          handledByGuardID,
		HandledByGuardName:        handledByGuardFullName,
		HandledByGuardPhoneNumber: handledByGuardPhoneNumber,
		CreatedBy:                 createdBy,
		QRExpiresAt:               pgTimestamptzToTimePtr(qrExpiresAt),
		QRUsedAt:                  pgTimestamptzToTimePtr(qrUsedAt),
		Notes:                     notes,
		RejectionReason:           rejectionReason,
		Metadata:                  metadataMap(metadata),
		CreatedAt:                 pgTimestamptzToTime(createdAt),
		UpdatedAt:                 pgTimestamptzToTime(updatedAt),
		Visitor:                   visitorSummary(visitorFullName, visitorPhone, visitorEmail, visitorPhoto),
		Flat:                      flatSummary(flatID, flatNumber, block, floor),
	}
}

func visitorEntryEventFromDBNoRows(row db.VisitorEntryEvent, err error) (*models.VisitorEntryEvent, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return visitorEntryEventFromDB(row), nil
}

func visitorEntryEventFromDB(row db.VisitorEntryEvent) *models.VisitorEntryEvent {
	return &models.VisitorEntryEvent{
		ID:             row.ID,
		VisitorEntryID: row.VisitorEntryID,
		SocietyID:      row.SocietyID,
		ActorUserID:    row.ActorUserID,
		EventType:      models.VisitorEventType(row.EventType),
		Message:        row.Message,
		Metadata:       metadataMap(row.Metadata),
		CreatedAt:      pgTimestamptzToTime(row.CreatedAt),
	}
}

func metadataMap(raw []byte) map[string]any {
	result := map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &result)
	}
	return result
}

func companionDetailsFromJSON(raw []byte) []map[string]any {
	result := []map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &result)
	}
	return result
}

func visitorSummary(name string, phone, email, photo *string) *models.VisitorSummary {
	if name == "" && phone == nil && email == nil && photo == nil {
		return nil
	}
	return &models.VisitorSummary{FullName: name, PhoneNumber: phone, Email: email, PhotoURL: photo}
}

func flatSummary(id int64, flatNumber string, block, floor *string) *models.VisitorFlatSummary {
	if id <= 0 || flatNumber == "" {
		return nil
	}
	return &models.VisitorFlatSummary{ID: id, FlatNumber: flatNumber, Block: block, Floor: floor}
}

func dbVisitorVehicleTypePtr(value *models.VisitorVehicleType) *db.VisitorVehicleType {
	if value == nil || *value == "" {
		return nil
	}
	converted := db.VisitorVehicleType(*value)
	return &converted
}

func visitorVehicleTypeFromDB(value *db.VisitorVehicleType) *models.VisitorVehicleType {
	if value == nil || *value == "" {
		return nil
	}
	converted := models.VisitorVehicleType(*value)
	return &converted
}

func dbVisitorStatusPtr(value *models.VisitorStatus) *db.VisitorStatus {
	if value == nil || *value == "" {
		return nil
	}
	converted := db.VisitorStatus(*value)
	return &converted
}

func dbVisitorSourcePtr(value *models.VisitorEntrySource) *db.VisitorSource {
	if value == nil || *value == "" {
		return nil
	}
	converted := db.VisitorSource(*value)
	return &converted
}

func dbVisitorPurposePtr(value *models.VisitorPurpose) *db.VisitorPurpose {
	if value == nil || *value == "" {
		return nil
	}
	converted := db.VisitorPurpose(*value)
	return &converted
}

func approvedByForCreate(status models.VisitorStatus, actorUserID *int64) *int64 {
	if status == models.VisitorStatusApproved {
		return actorUserID
	}
	return nil
}

func approvedAtForCreate(status models.VisitorStatus) pgtype.Timestamptz {
	if status == models.VisitorStatusApproved {
		return visitorTimeToPgTimestamptz(time.Now())
	}
	return pgtype.Timestamptz{}
}

func visitorTimeToPgTimestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func flatIDPtr(flatID int64) *int64 {
	if flatID <= 0 {
		return nil
	}
	value := flatID
	return &value
}

func flatIDValue(flatID *int64) int64 {
	if flatID == nil {
		return 0
	}
	return *flatID
}

func flatNumberValue(flatNumber *string) string {
	if flatNumber == nil {
		return ""
	}
	return *flatNumber
}

func normalizeVisitorLimit(limit int32) int32 {
	if limit <= 0 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}
