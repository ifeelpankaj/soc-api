package flatauthz

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"testing"

	"go-server/internal/models"
)

type fakeMemberRepo struct {
	member *models.SocietyMember
	err    error
}

func (f *fakeMemberRepo) Add(context.Context, *models.SocietyMember) error { return nil }
func (f *fakeMemberRepo) Get(_ context.Context, _ models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	return f.member, f.err
}
func (f *fakeMemberRepo) List(context.Context, models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) ListByUser(context.Context, int64) ([]*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) ListMySocietiesByUser(context.Context, int64) ([]*models.MySocietyResponse, error) {
	return nil, nil
}
func (f *fakeMemberRepo) Count(context.Context, models.ListSocietyMembersFilter) (int64, error) {
	return 0, nil
}
func (f *fakeMemberRepo) ChangeRole(context.Context, int64, int64, models.SocietyMemberRole) (*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) Suspend(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) Reactivate(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) Remove(context.Context, int64, int64, int64, string) error { return nil }
func (f *fakeMemberRepo) CountActiveOwners(context.Context, int64) (int64, error)   { return 0, nil }
func (f *fakeMemberRepo) DemoteActiveOwners(context.Context, int64, int64) error    { return nil }
func (f *fakeMemberRepo) PromoteToOwner(context.Context, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}
func (f *fakeMemberRepo) UpsertResident(context.Context, int64, int64, int64) (*models.SocietyMember, error) {
	return nil, nil
}

type fakeResidentRepo struct {
	resident *models.FlatResident
	err      error
}

func (f *fakeResidentRepo) Add(context.Context, *models.FlatResident) error { return nil }
func (f *fakeResidentRepo) Get(_ context.Context, _ *models.FlatResidentFilter) (*models.FlatResident, error) {
	return f.resident, f.err
}
func (f *fakeResidentRepo) List(context.Context, *models.FlatResidentFilter) ([]*models.FlatResident, error) {
	return nil, nil
}
func (f *fakeResidentRepo) Remove(context.Context, *models.FlatResidentFilter) error { return nil }
func (f *fakeResidentRepo) MoveOut(context.Context, *models.FlatResidentFilter) (*models.FlatResident, error) {
	return nil, nil
}
func (f *fakeResidentRepo) ClearPrimary(context.Context, int64, int64) error { return nil }
func (f *fakeResidentRepo) SetPrimary(context.Context, int64, int64, int64) (*models.FlatResident, error) {
	return nil, nil
}
func (f *fakeResidentRepo) UpdateRole(context.Context, *models.FlatResidentFilter, models.FlatResidentRole) (*models.FlatResident, error) {
	return nil, nil
}
func (f *fakeResidentRepo) CountActive(context.Context, int64, int64) (int64, error) {
	return 0, nil
}
func (f *fakeResidentRepo) CountPrimary(context.Context, int64, int64) (int64, error) {
	return 0, nil
}

type fakeFlatRepo struct {
	flat *models.Flat
	err  error
}

func (f *fakeFlatRepo) Create(context.Context, *models.Flat) error { return nil }
func (f *fakeFlatRepo) Get(_ context.Context, _ *models.FlatFilter) (*models.Flat, error) {
	return f.flat, f.err
}
func (f *fakeFlatRepo) List(context.Context, *models.FlatFilter) ([]*models.Flat, error) {
	return nil, nil
}
func (f *fakeFlatRepo) Count(context.Context, *models.FlatFilter) (int64, error) { return 0, nil }
func (f *fakeFlatRepo) Stats(context.Context, int64) (*models.FlatStatsResponse, error) {
	return nil, nil
}
func (f *fakeFlatRepo) Update(context.Context, *models.FlatFilter, *contracts.UpdateFlatInput) (*models.Flat, error) {
	return nil, nil
}
func (f *fakeFlatRepo) Deactivate(context.Context, *models.FlatFilter) error { return nil }
func (f *fakeFlatRepo) MarkOccupied(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}
func (f *fakeFlatRepo) MarkVacant(context.Context, int64, int64) (*models.Flat, error) {
	return nil, nil
}
func (f *fakeFlatRepo) Block(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}
func (f *fakeFlatRepo) Unblock(context.Context, *models.FlatFilter) (*models.Flat, error) {
	return nil, nil
}

func testAuthz(member *models.SocietyMember, resident *models.FlatResident) *FlatVisitorAuthz {
	return New(
		&fakeMemberRepo{member: member},
		&fakeResidentRepo{resident: resident},
		&fakeFlatRepo{flat: &models.Flat{ID: 10, SocietyID: 1}},
	)
}

func assertAppErrorCode(t *testing.T, err error, wantErr error) {
	t.Helper()

	var appErr *models.AppError
	var wantAppErr *models.AppError
	if !errors.As(err, &appErr) || !errors.As(wantErr, &wantAppErr) || appErr.Code != wantAppErr.Code {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestCanManageFlatVisitors(t *testing.T) {
	ctx := context.Background()
	const societyID int64 = 1
	const flatID int64 = 10

	tests := []struct {
		name     string
		member   *models.SocietyMember
		resident *models.FlatResident
		wantErr  error
	}{
		{
			name: "society admin",
			member: &models.SocietyMember{
				SocietyID: societyID,
				UserID:    100,
				Role:      models.SocietyMemberRoleAdmin,
				Status:    models.SocietyMemberStatusActive,
			},
		},
		{
			name: "society owner",
			member: &models.SocietyMember{
				SocietyID: societyID,
				UserID:    100,
				Role:      models.SocietyMemberRoleOwner,
				Status:    models.SocietyMemberStatusActive,
			},
		},
		{
			name: "primary resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				IsPrimary: true,
				Role:      models.FlatResidentRoleFamily,
			},
		},
		{
			name: "flat owner role without primary",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				IsPrimary: false,
				Role:      models.FlatResidentRoleOwner,
			},
		},
		{
			name: "family resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleFamily,
			},
		},
		{
			name: "tenant resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleTenant,
			},
		},
		{
			name:    "unrelated user",
			wantErr: ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authz := testAuthz(tt.member, tt.resident)
			err := authz.CanManageFlatVisitors(ctx, societyID, flatID, 100)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				assertAppErrorCode(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCanViewFlatVisitors(t *testing.T) {
	ctx := context.Background()
	const societyID int64 = 1
	const flatID int64 = 10

	tests := []struct {
		name     string
		member   *models.SocietyMember
		resident *models.FlatResident
		wantErr  error
	}{
		{
			name: "society staff",
			member: &models.SocietyMember{
				SocietyID: societyID,
				UserID:    100,
				Role:      models.SocietyMemberRoleStaff,
				Status:    models.SocietyMemberStatusActive,
			},
		},
		{
			name: "family resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleFamily,
			},
		},
		{
			name:    "unrelated user",
			wantErr: ErrViewForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authz := testAuthz(tt.member, tt.resident)
			err := authz.CanViewFlatVisitors(ctx, societyID, flatID, 100)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				assertAppErrorCode(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCanManageFlatMembers(t *testing.T) {
	ctx := context.Background()
	const societyID int64 = 1
	const flatID int64 = 10

	tests := []struct {
		name     string
		resident *models.FlatResident
		wantErr  error
	}{
		{
			name: "primary resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleFamily,
				IsPrimary: true,
			},
		},
		{
			name: "owner role resident",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleOwner,
			},
		},
		{
			name: "family non-primary",
			resident: &models.FlatResident{
				SocietyID: societyID,
				FlatID:    flatID,
				UserID:    100,
				Status:    models.FlatResidentStatusActive,
				Role:      models.FlatResidentRoleFamily,
			},
			wantErr: ErrManageForbidden,
		},
		{
			name:    "unrelated user",
			wantErr: ErrManageForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authz := testAuthz(nil, tt.resident)
			err := authz.CanManageFlatMembers(ctx, societyID, flatID, 100)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				assertAppErrorCode(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCanViewFlatMemberInvites(t *testing.T) {
	ctx := context.Background()
	const societyID int64 = 1
	const flatID int64 = 10

	tests := []struct {
		name            string
		member          *models.SocietyMember
		resident        *models.FlatResident
		wantRevealLinks bool
		wantErr         error
	}{
		{
			name:   "society admin can audit without links",
			member: &models.SocietyMember{SocietyID: societyID, UserID: 100, Role: models.SocietyMemberRoleAdmin, Status: models.SocietyMemberStatusActive},
		},
		{
			name:            "flat owner can audit with links",
			resident:        &models.FlatResident{SocietyID: societyID, FlatID: flatID, UserID: 100, Status: models.FlatResidentStatusActive, Role: models.FlatResidentRoleOwner},
			wantRevealLinks: true,
		},
		{
			name:            "primary resident can audit with links",
			resident:        &models.FlatResident{SocietyID: societyID, FlatID: flatID, UserID: 100, Status: models.FlatResidentStatusActive, Role: models.FlatResidentRoleFamily, IsPrimary: true},
			wantRevealLinks: true,
		},
		{
			name:     "family resident is rejected",
			resident: &models.FlatResident{SocietyID: societyID, FlatID: flatID, UserID: 100, Status: models.FlatResidentStatusActive, Role: models.FlatResidentRoleFamily},
			wantErr:  ErrManageForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authz := testAuthz(tt.member, tt.resident)
			got, err := authz.CanViewFlatMemberInvites(ctx, societyID, flatID, 100)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected an authorization error")
				}
				assertAppErrorCode(t, err, tt.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantRevealLinks {
				t.Fatalf("canRevealLinks = %v, want %v", got, tt.wantRevealLinks)
			}
		})
	}
}
