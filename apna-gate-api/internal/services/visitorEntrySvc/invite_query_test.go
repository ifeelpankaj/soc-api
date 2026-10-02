package visitorentrysvc_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"go-server/internal/repositories/contracts"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	flatauthz "go-server/internal/services/flatAuthz"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
)

func hashInviteToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type inviteQuerySettingSvc struct {
	settings *models.SocietyVisitorSettingsResponse
}

func (s *inviteQuerySettingSvc) GetSocietySettings(context.Context, int64) (*models.SocietyVisitorSettingsResponse, error) {
	return s.settings, nil
}
func (s *inviteQuerySettingSvc) GetFlatSettings(context.Context, int64, int64) ([]models.FlatVisitorSettingsResponse, error) {
	return nil, nil
}
func (s *inviteQuerySettingSvc) EnsureDefaultFlatSettingsIfMissing(context.Context, int64, int64, int64) error {
	return nil
}
func (s *inviteQuerySettingSvc) ResolveApprovalRequirement(context.Context, int64, int64, models.VisitorPurpose, models.VisitorEntrySource) (bool, error) {
	return false, nil
}
func (s *inviteQuerySettingSvc) ResolveVisitDurationMinutes(context.Context, int64, int64, models.VisitorPurpose) (int32, error) {
	return 60, nil
}

type inviteQueryInviteRepo struct {
	invite            *models.VisitorInvite
	expectedTokenHash string
}

func (r *inviteQueryInviteRepo) Create(context.Context, int64, int64, models.VisitorPurpose, string, time.Time, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) GetByID(context.Context, int64, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) GetHistoryByID(context.Context, int64, int64, int64, func(string, *time.Time) *models.ShortLinkResponse) (*models.VisitorInviteHistoryItem, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) GetByTokenHash(_ context.Context, tokenHash string) (*models.VisitorInvite, error) {
	if r.invite != nil && r.expectedTokenHash == tokenHash {
		return r.invite, nil
	}
	return nil, nil
}
func (r *inviteQueryInviteRepo) ListHistory(context.Context, int64, int64, models.VisitorInviteHistoryFilter, func(string, *time.Time) *models.ShortLinkResponse) ([]*models.VisitorInviteHistoryItem, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) CountHistory(context.Context, int64, int64, models.VisitorInviteHistoryFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryInviteRepo) MarkUsed(context.Context, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) GetForUpdate(context.Context, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) Cancel(context.Context, int64, int64) (*models.VisitorInvite, error) {
	return nil, nil
}
func (r *inviteQueryInviteRepo) ExpireOld(context.Context) error { return nil }

type inviteQueryFlatRepo struct {
	flat *models.Flat
}

func (r *inviteQueryFlatRepo) Create(context.Context, *models.Flat) error { return nil }
func (r *inviteQueryFlatRepo) Get(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return r.flat, nil
}
func (r *inviteQueryFlatRepo) List(context.Context, *models.FlatFilter) ([]*models.Flat, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) Count(context.Context, *models.FlatFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryFlatRepo) Stats(context.Context, int64) (*models.FlatStatsResponse, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) Update(context.Context, *models.FlatFilter, *contracts.UpdateFlatInput) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) Deactivate(context.Context, *models.FlatFilter) error { return nil }
func (r *inviteQueryFlatRepo) Block(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) Unblock(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) MarkOccupied(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}
func (r *inviteQueryFlatRepo) MarkVacant(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}

type inviteQueryEntryRepo struct {
	entry *models.VisitorEntry
}

func (r *inviteQueryEntryRepo) Create(context.Context, contracts.VisitorFormInput, int64, *int64, int64, *int64, models.VisitorEntrySource, models.VisitorPurpose, models.VisitorStatus, *int64, *int64, *string, *time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) Get(context.Context, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetForUpdate(context.Context, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetByQRHash(context.Context, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetByInviteID(context.Context, int64) (*models.VisitorEntry, error) {
	return r.entry, nil
}
func (r *inviteQueryEntryRepo) List(context.Context, models.VisitorEntryFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) Count(context.Context, models.VisitorEntryFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) ListPending(context.Context, int64, int64) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) ListSocietyPending(context.Context, models.VisitorPendingFilter) ([]*models.VisitorPendingEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) CountSocietyPending(context.Context, models.VisitorPendingFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) ListRecentByFlat(context.Context, int64, int64, int32) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetStats(context.Context, int64) (*models.VisitorEntryStatsResponse, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetStatsInRange(context.Context, int64, time.Time, time.Time) (*models.VisitorEntryStatsResponse, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) GetDailyStatsCreated(context.Context, int64, int32) ([]models.VisitorDailyCountResponse, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) CountWaitingAtGate(context.Context, int64) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) ListWaitingAtGate(context.Context, models.WaitingAtGateFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) CountWaitingAtGateFiltered(context.Context, models.WaitingAtGateFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) CountExpectedGuests(context.Context, int64, time.Time, time.Time) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) ListExpectedGuests(context.Context, models.ExpectedGuestFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) CountExpectedGuestsFiltered(context.Context, models.ExpectedGuestFilter) (int64, error) {
	return 0, nil
}
func (r *inviteQueryEntryRepo) CountMemberApprovals(context.Context, int64, int64) (*models.MemberVisitorApprovalStatsResponse, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) Approve(context.Context, int64, int64, int64, string, time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) Reject(context.Context, int64, int64, int64, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) MergeMetadata(_ context.Context, _ int64, _ int64, metadata map[string]any) (*models.VisitorEntry, error) {
	if r.entry == nil {
		return nil, nil
	}
	entry := *r.entry
	if entry.Metadata == nil {
		entry.Metadata = map[string]any{}
	}
	for key, value := range metadata {
		entry.Metadata[key] = value
	}
	r.entry = &entry
	return &entry, nil
}

func (r *inviteQueryEntryRepo) GenerateQR(_ context.Context, societyID int64, entryID int64, qrHash string, qrExpiresAt time.Time) (*models.VisitorEntry, error) {
	entry := *r.entry
	entry.QRExpiresAt = &qrExpiresAt
	return &entry, nil
}
func (r *inviteQueryEntryRepo) CheckIn(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) CheckOut(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) UpdateDetails(context.Context, int64, int64, contracts.UpdateGuardVisitorEntryInput) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *inviteQueryEntryRepo) AutoCloseExpired(context.Context) error   { return nil }
func (r *inviteQueryEntryRepo) ExpireStaleEntries(context.Context) error { return nil }

type inviteQueryEventRepo struct{}

func (inviteQueryEventRepo) Create(context.Context, int64, int64, *int64, models.VisitorEventType, *string, map[string]any) (*models.VisitorEntryEvent, error) {
	return &models.VisitorEntryEvent{}, nil
}
func (inviteQueryEventRepo) List(context.Context, int64, int64) ([]*models.VisitorEntryEvent, error) {
	return nil, nil
}

func TestGetPublicInviteByTokenReturnsQRForUsedInviteWithValidQR(t *testing.T) {
	const rawToken = "test-invite-token"
	societyID := int64(10)
	flatID := int64(20)
	inviteID := int64(30)
	entryID := int64(40)
	qrExpiresAt := time.Now().Add(2 * time.Hour)

	invite := &models.VisitorInvite{
		ID:        inviteID,
		SocietyID: societyID,
		FlatID:    flatID,
		Purpose:   models.VisitorPurposeGuest,
		Status:    models.VisitorInviteStatusUsed,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	entry := &models.VisitorEntry{
		ID:          entryID,
		SocietyID:   societyID,
		FlatID:      flatID,
		InviteID:    &inviteID,
		Status:      models.VisitorStatusApproved,
		QRExpiresAt: &qrExpiresAt,
	}

	block := "left-wing"
	floor := "1"
	flat := &models.Flat{
		ID:         flatID,
		SocietyID:  societyID,
		FlatNumber: "G-01",
		Block:      &block,
		Floor:      &floor,
	}
	society := &models.Society{
		ID:          societyID,
		Name:        "Test Society",
		SocietyCode: "TS001",
		Status:      models.SocietyStatusActive,
	}

	inviteSvc, _ := visitorentrysvc.NewVisitorService(
		nil,
		&inviteQueryInviteRepo{invite: invite, expectedTokenHash: hashInviteToken(rawToken)},
		&inviteQueryEntryRepo{entry: entry},
		inviteQueryEventRepo{},
		&inviteQuerySettingSvc{settings: &models.SocietyVisitorSettingsResponse{
			SocietyID:       societyID,
			QRExpiryMinutes: 60,
			IsActive:        true,
		}},
		nil,
		nil,
		&inviteQueryFlatRepo{flat: flat},
		&guardDeskSocietyRepo{society: society},
		flatauthz.New(nil, nil, nil),
		nil,
		noopTxManager{},
	)

	got, err := inviteSvc.GetPublicInviteByToken(context.Background(), rawToken)
	if err != nil {
		t.Fatalf("GetPublicInviteByToken() error = %v", err)
	}
	if got.View != models.PublicVisitorInviteViewQR {
		t.Fatalf("view = %q, want %q", got.View, models.PublicVisitorInviteViewQR)
	}
	if got.QR == nil || got.QR.Token == "" {
		t.Fatalf("expected non-empty qr token, got %+v", got.QR)
	}
	if got.Entry == nil || got.Entry.ID != entryID {
		t.Fatalf("expected entry id %d, got %+v", entryID, got.Entry)
	}
}

var (
	_ repository.VisitorInviteRepository     = (*inviteQueryInviteRepo)(nil)
	_ repository.FlatRepository              = (*inviteQueryFlatRepo)(nil)
	_ repository.VisitorEntryRepository      = (*inviteQueryEntryRepo)(nil)
	_ repository.VisitorEntryEventRepository = inviteQueryEventRepo{}
)
