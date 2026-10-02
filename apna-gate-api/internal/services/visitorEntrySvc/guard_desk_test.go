package visitorentrysvc_test

import (
	"context"
	"go-server/internal/repositories/contracts"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
	flatauthz "go-server/internal/services/flatAuthz"
	visitorentrysvc "go-server/internal/services/visitorEntrySvc"
)

type guardDeskSocietyRepo struct {
	society *models.Society
}

func (r *guardDeskSocietyRepo) Create(context.Context, *models.Society) error { return nil }
func (r *guardDeskSocietyRepo) Get(_ context.Context, _ models.GetSocietyFilter) (*models.Society, error) {
	return r.society, nil
}
func (r *guardDeskSocietyRepo) List(context.Context, models.ListSocietiesFilter) ([]*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Count(context.Context, models.ListSocietiesFilter) (int64, error) {
	return 0, nil
}
func (r *guardDeskSocietyRepo) Update(context.Context, int64, contracts.UpdateSocietyInput) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Approve(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Reject(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Suspend(context.Context, int64, int64, string) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Reactivate(context.Context, int64, int64) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) Restore(context.Context, int64) (*models.Society, error) {
	return nil, nil
}
func (r *guardDeskSocietyRepo) SoftDelete(context.Context, int64) error { return nil }
func (r *guardDeskSocietyRepo) CountPendingByCreator(context.Context, int64) (int64, error) {
	return 0, nil
}

type guardDeskSettingSvc struct{}

func (guardDeskSettingSvc) GetSocietySettings(context.Context, int64) (*models.SocietyVisitorSettingsResponse, error) {
	return &models.SocietyVisitorSettingsResponse{AllowGuardEntry: true, AllowGuardOnBehalfApproval: true}, nil
}
func (guardDeskSettingSvc) GetFlatSettings(context.Context, int64, int64) ([]models.FlatVisitorSettingsResponse, error) {
	return nil, nil
}
func (guardDeskSettingSvc) EnsureDefaultFlatSettingsIfMissing(context.Context, int64, int64, int64) error {
	return nil
}
func (guardDeskSettingSvc) ResolveApprovalRequirement(context.Context, int64, int64, models.VisitorPurpose, models.VisitorEntrySource) (bool, error) {
	return false, nil
}
func (guardDeskSettingSvc) ResolveVisitDurationMinutes(context.Context, int64, int64, models.VisitorPurpose) (int32, error) {
	return 60, nil
}

type guardDeskEntryRepo struct {
	stats               *models.VisitorEntryStatsResponse
	expectedGuestsCount int64
	expectedGuests      []*models.VisitorEntry
	waitingAtGateCount  int64
	pending             []*models.VisitorPendingEntry
}

func (r *guardDeskEntryRepo) Create(context.Context, contracts.VisitorFormInput, int64, *int64, int64, *int64, models.VisitorEntrySource, models.VisitorPurpose, models.VisitorStatus, *int64, *int64, *string, *time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) Get(context.Context, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) GetForUpdate(context.Context, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) GetByQRHash(context.Context, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) GetByInviteID(context.Context, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) List(context.Context, models.VisitorEntryFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) Count(context.Context, models.VisitorEntryFilter) (int64, error) {
	return 0, nil
}
func (r *guardDeskEntryRepo) ListPending(context.Context, int64, int64) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) ListSocietyPending(context.Context, models.VisitorPendingFilter) ([]*models.VisitorPendingEntry, error) {
	return r.pending, nil
}
func (r *guardDeskEntryRepo) CountSocietyPending(context.Context, models.VisitorPendingFilter) (int64, error) {
	return int64(len(r.pending)), nil
}
func (r *guardDeskEntryRepo) ListRecentByFlat(context.Context, int64, int64, int32) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) GetStats(context.Context, int64) (*models.VisitorEntryStatsResponse, error) {
	return r.stats, nil
}
func (r *guardDeskEntryRepo) GetStatsInRange(context.Context, int64, time.Time, time.Time) (*models.VisitorEntryStatsResponse, error) {
	return r.stats, nil
}
func (r *guardDeskEntryRepo) GetDailyStatsCreated(context.Context, int64, int32) ([]models.VisitorDailyCountResponse, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) CountWaitingAtGate(context.Context, int64) (int64, error) {
	return r.waitingAtGateCount, nil
}
func (r *guardDeskEntryRepo) ListWaitingAtGate(context.Context, models.WaitingAtGateFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) CountWaitingAtGateFiltered(context.Context, models.WaitingAtGateFilter) (int64, error) {
	return r.waitingAtGateCount, nil
}
func (r *guardDeskEntryRepo) CountExpectedGuests(_ context.Context, _ int64, _, _ time.Time) (int64, error) {
	return r.expectedGuestsCount, nil
}
func (r *guardDeskEntryRepo) ListExpectedGuests(context.Context, models.ExpectedGuestFilter) ([]*models.VisitorEntry, error) {
	return r.expectedGuests, nil
}
func (r *guardDeskEntryRepo) CountExpectedGuestsFiltered(context.Context, models.ExpectedGuestFilter) (int64, error) {
	return r.expectedGuestsCount, nil
}
func (r *guardDeskEntryRepo) MergeMetadata(context.Context, int64, int64, map[string]any) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) CountMemberApprovals(context.Context, int64, int64) (*models.MemberVisitorApprovalStatsResponse, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) Approve(context.Context, int64, int64, int64, string, time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) Reject(context.Context, int64, int64, int64, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) GenerateQR(context.Context, int64, int64, string, time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) CheckIn(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) CheckOut(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) UpdateDetails(context.Context, int64, int64, contracts.UpdateGuardVisitorEntryInput) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *guardDeskEntryRepo) AutoCloseExpired(context.Context) error   { return nil }
func (r *guardDeskEntryRepo) ExpireStaleEntries(context.Context) error { return nil }

type noopTxManager struct{}

func (noopTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestGetGuardDeskBootstrap(t *testing.T) {
	societyID := int64(10)
	society := &models.Society{ID: societyID, Name: "Test Society", SocietyCode: "TS001", Status: models.SocietyStatusActive}
	stats := &models.VisitorEntryStatsResponse{
		TodayVisitors:    5,
		VisitorsInside:   2,
		PendingApprovals: 1,
	}
	pending := []*models.VisitorPendingEntry{
		{VisitorEntry: &models.VisitorEntry{ID: 99, SocietyID: societyID, FlatID: 3, Status: models.VisitorStatusWaitingApproval}},
	}

	_, entrySvc := visitorentrysvc.NewVisitorService(
		nil,
		nil,
		&guardDeskEntryRepo{stats: stats, expectedGuestsCount: 2, waitingAtGateCount: 4, pending: pending},
		nil,
		guardDeskSettingSvc{},
		nil,
		nil,
		nil,
		&guardDeskSocietyRepo{society: society},
		flatauthz.New(nil, nil, nil),
		nil,
		noopTxManager{},
	)

	got, err := entrySvc.GetGuardDeskBootstrap(context.Background(), societyID)
	if err != nil {
		t.Fatalf("GetGuardDeskBootstrap() error = %v", err)
	}
	if got.Society == nil || got.Society.Name != "Test Society" {
		t.Fatalf("expected society name Test Society, got %+v", got.Society)
	}
	if got.Stats == nil || got.Stats.PendingApprovals != 1 {
		t.Fatalf("expected pending approvals 1, got %+v", got.Stats)
	}
	if got.ExpectedGuestsCount != 2 {
		t.Fatalf("expected guests count = 2, got %d", got.ExpectedGuestsCount)
	}
	if got.WaitingAtGateCount != 4 {
		t.Fatalf("waiting at gate count = 4, got %d", got.WaitingAtGateCount)
	}
	if len(got.PendingPreview) != 1 {
		t.Fatalf("expected 1 pending preview entry, got %d", len(got.PendingPreview))
	}
}

func TestListExpectedGuests(t *testing.T) {
	societyID := int64(10)
	entries := []*models.VisitorEntry{
		{ID: 1, SocietyID: societyID, FlatID: 3, Source: models.VisitorEntrySourceResidentLink, Status: models.VisitorStatusApproved},
	}

	_, entrySvc := visitorentrysvc.NewVisitorService(
		nil,
		nil,
		&guardDeskEntryRepo{expectedGuestsCount: 1, expectedGuests: entries},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		noopTxManager{},
	)

	got, err := entrySvc.ListExpectedGuests(context.Background(), models.ExpectedGuestFilter{
		SocietyID: societyID,
		Limit:     20,
	})
	if err != nil {
		t.Fatalf("ListExpectedGuests() error = %v", err)
	}
	if got.Total != 1 {
		t.Fatalf("expected total 1, got %d", got.Total)
	}
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
}

var _ repository.SocietyRepository = (*guardDeskSocietyRepo)(nil)
var _ repository.VisitorEntryRepository = (*guardDeskEntryRepo)(nil)
var _ repository.TransactionManager = noopTxManager{}
