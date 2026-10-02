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

type staffEntryRepo struct {
	entry *models.VisitorEntry
}

func (r *staffEntryRepo) Get(_ context.Context, societyID int64, entryID int64) (*models.VisitorEntry, error) {
	if r.entry != nil && r.entry.SocietyID == societyID && r.entry.ID == entryID {
		return r.entry, nil
	}
	return nil, nil
}

func (r *staffEntryRepo) List(_ context.Context, filter models.VisitorEntryFilter) ([]*models.VisitorEntry, error) {
	if r.entry == nil || r.entry.SocietyID != filter.SocietyID {
		return nil, nil
	}
	if filter.Purpose != nil && *filter.Purpose != r.entry.Purpose {
		return nil, nil
	}
	return []*models.VisitorEntry{r.entry}, nil
}

func (r *staffEntryRepo) Count(_ context.Context, filter models.VisitorEntryFilter) (int64, error) {
	entries, err := r.List(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return int64(len(entries)), nil
}

func (r *staffEntryRepo) Create(context.Context, contracts.VisitorFormInput, int64, *int64, int64, *int64, models.VisitorEntrySource, models.VisitorPurpose, models.VisitorStatus, *int64, *int64, *string, *time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetForUpdate(context.Context, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetByQRHash(context.Context, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetByInviteID(context.Context, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) ListPending(context.Context, int64, int64) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) ListSocietyPending(context.Context, models.VisitorPendingFilter) ([]*models.VisitorPendingEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) CountSocietyPending(context.Context, models.VisitorPendingFilter) (int64, error) {
	return 0, nil
}
func (r *staffEntryRepo) ListRecentByFlat(context.Context, int64, int64, int32) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetStats(context.Context, int64) (*models.VisitorEntryStatsResponse, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetStatsInRange(context.Context, int64, time.Time, time.Time) (*models.VisitorEntryStatsResponse, error) {
	return nil, nil
}
func (r *staffEntryRepo) GetDailyStatsCreated(context.Context, int64, int32) ([]models.VisitorDailyCountResponse, error) {
	return nil, nil
}
func (r *staffEntryRepo) CountWaitingAtGate(context.Context, int64) (int64, error) {
	return 0, nil
}
func (r *staffEntryRepo) ListWaitingAtGate(context.Context, models.WaitingAtGateFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) CountWaitingAtGateFiltered(context.Context, models.WaitingAtGateFilter) (int64, error) {
	return 0, nil
}
func (r *staffEntryRepo) CountExpectedGuests(context.Context, int64, time.Time, time.Time) (int64, error) {
	return 0, nil
}
func (r *staffEntryRepo) ListExpectedGuests(context.Context, models.ExpectedGuestFilter) ([]*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) CountExpectedGuestsFiltered(context.Context, models.ExpectedGuestFilter) (int64, error) {
	return 0, nil
}
func (r *staffEntryRepo) CountMemberApprovals(context.Context, int64, int64) (*models.MemberVisitorApprovalStatsResponse, error) {
	return nil, nil
}
func (r *staffEntryRepo) Approve(context.Context, int64, int64, int64, string, time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) MergeMetadata(context.Context, int64, int64, map[string]any) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) Reject(context.Context, int64, int64, int64, string) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) GenerateQR(context.Context, int64, int64, string, time.Time) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) CheckIn(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) CheckOut(context.Context, int64, int64, int64) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) UpdateDetails(context.Context, int64, int64, contracts.UpdateGuardVisitorEntryInput) (*models.VisitorEntry, error) {
	return nil, nil
}
func (r *staffEntryRepo) AutoCloseExpired(context.Context) error   { return nil }
func (r *staffEntryRepo) ExpireStaleEntries(context.Context) error { return nil }

func staffTestEntry() *models.VisitorEntry {
	return &models.VisitorEntry{
		ID:        15,
		SocietyID: 1,
		FlatID:    0,
		Flat:      nil,
		Source:    models.VisitorEntrySourceGuardEntry,
		Purpose:   models.VisitorPurposeStaff,
		Status:    models.VisitorStatusApproved,
		Visitor:   &models.VisitorSummary{FullName: "Staff Member", PhoneNumber: strPtr("9876543210")},
	}
}

func strPtr(value string) *string {
	return &value
}

func TestGetEntryReturnsStaffEntryWithoutFlat(t *testing.T) {
	entry := staffTestEntry()
	_, entrySvc := visitorentrysvc.NewVisitorService(
		nil,
		nil,
		&staffEntryRepo{entry: entry},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		flatauthz.New(nil, nil, nil),
		nil,
		noopTxManager{},
	)

	got, err := entrySvc.GetEntry(context.Background(), 1, 15)
	if err != nil {
		t.Fatalf("GetEntry() error = %v", err)
	}
	if got == nil || got.Purpose != models.VisitorPurposeStaff {
		t.Fatalf("expected staff entry, got %+v", got)
	}
	if got.Flat != nil {
		t.Fatalf("expected nil flat summary for staff entry, got %+v", got.Flat)
	}
}

func TestListEntriesPaginatedIncludesStaffEntry(t *testing.T) {
	entry := staffTestEntry()
	purpose := models.VisitorPurposeStaff
	_, entrySvc := visitorentrysvc.NewVisitorService(
		nil,
		nil,
		&staffEntryRepo{entry: entry},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		flatauthz.New(nil, nil, nil),
		nil,
		noopTxManager{},
	)

	got, err := entrySvc.ListEntriesPaginated(context.Background(), models.VisitorEntryFilter{
		SocietyID: 1,
		Purpose:   &purpose,
		Limit:     20,
	})
	if err != nil {
		t.Fatalf("ListEntriesPaginated() error = %v", err)
	}
	if got.Total != 1 || len(got.Entries) != 1 {
		t.Fatalf("expected one staff entry, got total=%d entries=%d", got.Total, len(got.Entries))
	}
	if got.Entries[0].Flat != nil {
		t.Fatalf("expected nil flat on listed staff entry, got %+v", got.Entries[0].Flat)
	}
}

var _ repository.VisitorEntryRepository = (*staffEntryRepo)(nil)
