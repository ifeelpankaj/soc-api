package societysvc

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"
	"testing"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
)

type guardSocietyRepo struct {
	society    *models.Society
	getErr     error
	update     *models.Society
	updateErr  error
	softDelete error
	getCalls   int
	lastFilter models.GetSocietyFilter
}

func (r *guardSocietyRepo) Create(context.Context, *models.Society) error { panic("unused") }
func (r *guardSocietyRepo) Get(_ context.Context, filter models.GetSocietyFilter) (*models.Society, error) {
	r.getCalls++
	r.lastFilter = filter
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.society, nil
}
func (r *guardSocietyRepo) List(context.Context, models.ListSocietiesFilter) ([]*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Count(context.Context, models.ListSocietiesFilter) (int64, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Update(context.Context, int64, contracts.UpdateSocietyInput) (*models.Society, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	return r.update, nil
}
func (r *guardSocietyRepo) Approve(context.Context, int64, int64) (*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Reject(context.Context, int64, int64, string) (*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Suspend(context.Context, int64, int64, string) (*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Reactivate(context.Context, int64, int64) (*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) Restore(context.Context, int64) (*models.Society, error) {
	panic("unused")
}
func (r *guardSocietyRepo) SoftDelete(context.Context, int64) error { return r.softDelete }
func (r *guardSocietyRepo) CountPendingByCreator(context.Context, int64) (int64, error) {
	panic("unused")
}

type guardMemberRepo struct {
	members           map[[2]int64]*models.SocietyMember
	getErr            error
	addErr            error
	changeRoleErr     error
	changedMember     *models.SocietyMember
	activeOwnerCount  int64
	activeOwnerErr    error
	added             *models.SocietyMember
	getCalls          int
	countOwnerCalls   int
	lastGetFilter     models.GetSocietyMemberFilter
	lastAddedMember   *models.SocietyMember
	lastChangeRole    models.SocietyMemberRole
	lastChangeRoleUID int64
}

func (r *guardMemberRepo) Add(_ context.Context, member *models.SocietyMember) error {
	if r.addErr != nil {
		return r.addErr
	}
	if member.ID == 0 {
		member.ID = 900
	}
	r.added = member
	r.lastAddedMember = member
	if r.members == nil {
		r.members = map[[2]int64]*models.SocietyMember{}
	}
	r.members[[2]int64{member.SocietyID, member.UserID}] = member
	return nil
}
func (r *guardMemberRepo) Get(_ context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	r.getCalls++
	r.lastGetFilter = filter
	if r.getErr != nil {
		return nil, r.getErr
	}
	if filter.SocietyID == nil || filter.UserID == nil {
		return nil, nil
	}
	member := r.members[[2]int64{*filter.SocietyID, *filter.UserID}]
	if member == nil {
		return nil, nil
	}
	if filter.Status != nil && string(member.Status) != *filter.Status {
		return nil, nil
	}
	if filter.Role != nil && string(member.Role) != *filter.Role {
		return nil, nil
	}
	return member, nil
}
func (r *guardMemberRepo) List(context.Context, models.ListSocietyMembersFilter) ([]*models.SocietyMember, error) {
	panic("unused")
}
func (r *guardMemberRepo) ListByUser(context.Context, int64) ([]*models.SocietyMember, error) {
	panic("unused")
}
func (r *guardMemberRepo) ListMySocietiesByUser(context.Context, int64) ([]*models.MySocietyResponse, error) {
	panic("unused")
}
func (r *guardMemberRepo) Count(context.Context, models.ListSocietyMembersFilter) (int64, error) {
	panic("unused")
}
func (r *guardMemberRepo) ChangeRole(_ context.Context, societyID int64, userID int64, role models.SocietyMemberRole) (*models.SocietyMember, error) {
	r.lastChangeRoleUID = userID
	r.lastChangeRole = role
	if r.changeRoleErr != nil {
		return nil, r.changeRoleErr
	}
	if r.changedMember != nil {
		return r.changedMember, nil
	}
	member := r.members[[2]int64{societyID, userID}]
	if member == nil {
		return nil, nil
	}
	member.Role = role
	return member, nil
}
func (r *guardMemberRepo) Suspend(context.Context, int64, int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *guardMemberRepo) Reactivate(context.Context, int64, int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *guardMemberRepo) Remove(context.Context, int64, int64, int64, string) error {
	panic("unused")
}
func (r *guardMemberRepo) CountActiveOwners(context.Context, int64) (int64, error) {
	r.countOwnerCalls++
	if r.activeOwnerErr != nil {
		return 0, r.activeOwnerErr
	}
	return r.activeOwnerCount, nil
}
func (r *guardMemberRepo) DemoteActiveOwners(context.Context, int64, int64) error { panic("unused") }
func (r *guardMemberRepo) PromoteToOwner(context.Context, int64, int64) (*models.SocietyMember, error) {
	panic("unused")
}
func (r *guardMemberRepo) UpsertResident(context.Context, int64, int64, int64) (*models.SocietyMember, error) {
	panic("unused")
}

type guardUserRepo struct {
	user        *models.User
	getByIDErr  error
	emailExists bool
	phoneExists bool
	createErr   error
	nextID      int64
}

func (r *guardUserRepo) Create(_ context.Context, user *models.User) error {
	if r.createErr != nil {
		return r.createErr
	}
	if user.ID == 0 {
		r.nextID++
		user.ID = r.nextID
	}
	r.user = user
	return nil
}
func (r *guardUserRepo) GetByID(context.Context, int64) (*models.User, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.user, nil
}
func (r *guardUserRepo) GetByEmail(context.Context, string) (*models.User, error) { panic("unused") }
func (r *guardUserRepo) GetByPhoneNumber(context.Context, string) (*models.User, error) {
	panic("unused")
}
func (r *guardUserRepo) EmailExists(context.Context, string) (bool, error) { return r.emailExists, nil }
func (r *guardUserRepo) PhoneExists(context.Context, string) (bool, error) {
	return r.phoneExists, nil
}
func (r *guardUserRepo) MarkEmailVerified(context.Context, int64) error { panic("unused") }
func (r *guardUserRepo) UpdatePasswordHash(context.Context, int64, string, int64) error {
	panic("unused")
}
func (r *guardUserRepo) UpdateLastLogin(context.Context, int64) error { panic("unused") }
func (r *guardUserRepo) UpdateProfile(context.Context, int64, *contracts.UpdateUserInput) (*models.User, error) {
	panic("unused")
}

type guardSubscriptionSvc struct {
	adminErr error
	staffErr error
}

func (s *guardSubscriptionSvc) CanAddAdmin(context.Context, int64, int64) error {
	return s.adminErr
}
func (s *guardSubscriptionSvc) CanAddStaff(context.Context, int64, int64) error {
	return s.staffErr
}
func (s *guardSubscriptionSvc) ListSubscriptions(context.Context, *models.SubscriptionFilter) ([]*models.SocietySubscriptionResponse, error) {
	panic("unused")
}
func (s *guardSubscriptionSvc) GetSubscriptionStats(context.Context, *models.SubscriptionFilter) (*models.SubscriptionStatsResponse, error) {
	panic("unused")
}

type guardTxManager struct{}

func (guardTxManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestEnsureActiveSociety(t *testing.T) {
	ctx := context.Background()
	sentinelErr := errors.New("database unavailable")

	tests := []struct {
		name    string
		society *models.Society
		getErr  error
		wantErr *models.AppError
		wantRaw error
	}{
		{name: "active society succeeds", society: guardSociety(10, models.SocietyStatusActive)},
		{name: "missing society fails", wantErr: ErrSocietyNotFound},
		{name: "pending society fails", society: guardSociety(10, models.SocietyStatusPending), wantErr: ErrSocietyInactive},
		{name: "suspended society fails", society: guardSociety(10, models.SocietyStatusSuspended), wantErr: ErrSocietyInactive},
		{name: "rejected society fails", society: guardSociety(10, models.SocietyStatusRejected), wantErr: ErrSocietyInactive},
		{name: "repository error is propagated", getErr: sentinelErr, wantRaw: sentinelErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			societyRepo := &guardSocietyRepo{society: tt.society, getErr: tt.getErr}
			svc := &SocietySvc{societyRepo: societyRepo}

			err := svc.EnsureActiveSociety(ctx, 10)

			if tt.wantRaw != nil {
				if !errors.Is(err, tt.wantRaw) {
					t.Fatalf("error = %v, want %v", err, tt.wantRaw)
				}
			} else if tt.wantErr != nil {
				requireSocietyAppErrCode(t, err, tt.wantErr.Code)
			} else if err != nil {
				t.Fatalf("EnsureActiveSociety() error = %v", err)
			}
			if societyRepo.lastFilter.ID == nil || *societyRepo.lastFilter.ID != 10 {
				t.Fatalf("society filter ID = %v, want 10", societyRepo.lastFilter.ID)
			}
		})
	}
}

func TestEnsureActiveMember(t *testing.T) {
	ctx := context.Background()
	sentinelErr := errors.New("member repo unavailable")

	tests := []struct {
		name    string
		member  *models.SocietyMember
		getErr  error
		wantErr *models.AppError
		wantRaw error
	}{
		{name: "active owner succeeds", member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive)},
		{name: "active admin succeeds", member: guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive)},
		{name: "active staff succeeds", member: guardMember(10, 100, models.SocietyMemberRoleStaff, models.SocietyMemberStatusActive)},
		{name: "active resident succeeds", member: guardMember(10, 100, models.SocietyMemberRoleResident, models.SocietyMemberStatusActive)},
		{name: "missing member fails", wantErr: ErrMemberNotFound},
		{name: "pending member fails", member: guardMember(10, 100, models.SocietyMemberRoleResident, models.SocietyMemberStatusPending), wantErr: ErrMemberInactive},
		{name: "suspended member fails", member: guardMember(10, 100, models.SocietyMemberRoleResident, models.SocietyMemberStatusSuspended), wantErr: ErrMemberInactive},
		{name: "removed member fails", member: guardMember(10, 100, models.SocietyMemberRoleResident, models.SocietyMemberStatusRemoved), wantErr: ErrMemberInactive},
		{name: "repository error is propagated", getErr: sentinelErr, wantRaw: sentinelErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberRepo := guardMemberRepoWith(tt.member)
			memberRepo.getErr = tt.getErr
			svc := &SocietySvc{memberRepo: memberRepo}

			got, err := svc.EnsureActiveMember(ctx, 10, 100)

			if tt.wantRaw != nil {
				if !errors.Is(err, tt.wantRaw) {
					t.Fatalf("error = %v, want %v", err, tt.wantRaw)
				}
				return
			}
			if tt.wantErr != nil {
				requireSocietyAppErrCode(t, err, tt.wantErr.Code)
				return
			}
			if err != nil {
				t.Fatalf("EnsureActiveMember() error = %v", err)
			}
			if got == nil || got.SocietyID != 10 || got.UserID != 100 || got.Role != tt.member.Role || got.Status != models.SocietyMemberStatusActive {
				t.Fatalf("member response = %+v, want active %s member", got, tt.member.Role)
			}
			if memberRepo.lastGetFilter.SocietyID == nil || *memberRepo.lastGetFilter.SocietyID != 10 {
				t.Fatalf("member society filter = %v, want 10", memberRepo.lastGetFilter.SocietyID)
			}
			if memberRepo.lastGetFilter.UserID == nil || *memberRepo.lastGetFilter.UserID != 100 {
				t.Fatalf("member user filter = %v, want 100", memberRepo.lastGetFilter.UserID)
			}
		})
	}
}

func TestEnsureRole(t *testing.T) {
	ctx := context.Background()
	societyErr := errors.New("society lookup failed")
	memberErr := errors.New("member lookup failed")

	tests := []struct {
		name                  string
		society               *models.Society
		member                *models.SocietyMember
		roles                 []string
		societyErr            error
		memberErr             error
		wantErr               *models.AppError
		wantRaw               error
		wantMemberLookupCalls int
	}{
		{name: "owner allowed", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleOwner)}, wantMemberLookupCalls: 1},
		{name: "admin allowed by multi role list", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleOwner), string(models.SocietyMemberRoleAdmin)}, wantMemberLookupCalls: 1},
		{name: "staff allowed", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleStaff, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleStaff)}, wantMemberLookupCalls: 1},
		{name: "resident rejected for admin action", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleResident, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleAdmin)}, wantErr: ErrForbiddenSociety, wantMemberLookupCalls: 1},
		{name: "no roles rejects active member", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), wantErr: ErrForbiddenSociety, wantMemberLookupCalls: 1},
		{name: "missing society stops before member lookup", member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleOwner)}, wantErr: ErrSocietyNotFound},
		{name: "inactive society stops before member lookup", society: guardSociety(10, models.SocietyStatusSuspended), member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleOwner)}, wantErr: ErrSocietyInactive},
		{name: "missing member fails", society: guardSociety(10, models.SocietyStatusActive), roles: []string{string(models.SocietyMemberRoleOwner)}, wantErr: ErrMemberNotFound, wantMemberLookupCalls: 1},
		{name: "inactive member fails", society: guardSociety(10, models.SocietyStatusActive), member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusSuspended), roles: []string{string(models.SocietyMemberRoleOwner)}, wantErr: ErrMemberInactive, wantMemberLookupCalls: 1},
		{name: "society repo error is propagated", societyErr: societyErr, member: guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), roles: []string{string(models.SocietyMemberRoleOwner)}, wantRaw: societyErr},
		{name: "member repo error is propagated", society: guardSociety(10, models.SocietyStatusActive), memberErr: memberErr, roles: []string{string(models.SocietyMemberRoleOwner)}, wantRaw: memberErr, wantMemberLookupCalls: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberRepo := guardMemberRepoWith(tt.member)
			memberRepo.getErr = tt.memberErr
			svc := &SocietySvc{
				societyRepo: &guardSocietyRepo{society: tt.society, getErr: tt.societyErr},
				memberRepo:  memberRepo,
			}

			err := svc.EnsureRole(ctx, 10, 100, tt.roles...)

			if tt.wantRaw != nil {
				if !errors.Is(err, tt.wantRaw) {
					t.Fatalf("error = %v, want %v", err, tt.wantRaw)
				}
			} else if tt.wantErr != nil {
				requireSocietyAppErrCode(t, err, tt.wantErr.Code)
			} else if err != nil {
				t.Fatalf("EnsureRole() error = %v", err)
			}
			if memberRepo.getCalls != tt.wantMemberLookupCalls {
				t.Fatalf("member lookup calls = %d, want %d", memberRepo.getCalls, tt.wantMemberLookupCalls)
			}
		})
	}
}

func TestEnsureMembershipRoleAllowsPreOperationalLifecycle(t *testing.T) {
	for _, status := range []models.SocietyStatus{
		models.SocietyStatusPending,
		models.SocietyStatusActive,
		models.SocietyStatusRejected,
		models.SocietyStatusSuspended,
	} {
		t.Run(string(status), func(t *testing.T) {
			svc := guardSocietySvc(
				guardSociety(10, status),
				guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive),
				nil,
			)

			if err := svc.EnsureMembershipRole(context.Background(), 10, 100, string(models.SocietyMemberRoleOwner)); err != nil {
				t.Fatalf("EnsureMembershipRole() error = %v", err)
			}
		})
	}
}

func TestUpdateAndDeleteSocietyUseRoleGuards(t *testing.T) {
	ctx := context.Background()
	updatedName := "Updated Society"

	for _, role := range []models.SocietyMemberRole{models.SocietyMemberRoleOwner, models.SocietyMemberRoleAdmin} {
		t.Run("update allowed for "+string(role), func(t *testing.T) {
			svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, role, models.SocietyMemberStatusActive), nil)
			svc.societyRepo.(*guardSocietyRepo).update = guardSociety(10, models.SocietyStatusActive)
			svc.societyRepo.(*guardSocietyRepo).update.Name = updatedName

			got, err := svc.UpdateSociety(ctx, 10, models.UpdateSocietyRequest{Name: &updatedName}, 100)
			if err != nil {
				t.Fatalf("UpdateSociety() error = %v", err)
			}
			if got.Name != updatedName {
				t.Fatalf("updated name = %q, want %q", got.Name, updatedName)
			}
		})
	}

	for _, role := range []models.SocietyMemberRole{models.SocietyMemberRoleStaff, models.SocietyMemberRoleResident} {
		t.Run("update rejected for "+string(role), func(t *testing.T) {
			svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, role, models.SocietyMemberStatusActive), nil)

			_, err := svc.UpdateSociety(ctx, 10, models.UpdateSocietyRequest{Name: &updatedName}, 100)
			requireSocietyAppErrCode(t, err, ErrForbiddenSociety.Code)
		})
	}

	t.Run("update rejected for inactive society", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusSuspended), guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), nil)

		_, err := svc.UpdateSociety(ctx, 10, models.UpdateSocietyRequest{Name: &updatedName}, 100)
		requireSocietyAppErrCode(t, err, ErrSocietyInactive.Code)
	})

	t.Run("delete allowed for owner only", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), nil)
		if err := svc.DeleteSociety(ctx, 10, 100); err != nil {
			t.Fatalf("DeleteSociety() error = %v", err)
		}
	})

	for _, role := range []models.SocietyMemberRole{models.SocietyMemberRoleAdmin, models.SocietyMemberRoleStaff, models.SocietyMemberRoleResident} {
		t.Run("delete rejected for "+string(role), func(t *testing.T) {
			svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, role, models.SocietyMemberStatusActive), nil)

			err := svc.DeleteSociety(ctx, 10, 100)
			requireSocietyAppErrCode(t, err, ErrForbiddenSociety.Code)
		})
	}
}

func TestAddMemberUsesRoleGuardAndQuota(t *testing.T) {
	ctx := context.Background()
	targetUser := &models.User{ID: 200, FullName: "Target User"}

	for _, role := range []models.SocietyMemberRole{models.SocietyMemberRoleOwner, models.SocietyMemberRoleAdmin} {
		t.Run("admin add allowed for "+string(role), func(t *testing.T) {
			subscription := &guardSubscriptionSvc{}
			svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, role, models.SocietyMemberStatusActive), targetUser)
			svc.subscriptionSvc = subscription

			got, err := svc.AddMember(ctx, models.AddSocietyMemberRequest{SocietyID: 10, UserID: 200, Role: models.SocietyMemberRoleAdmin}, 100)
			if err != nil {
				t.Fatalf("AddMember() error = %v", err)
			}
			if got.UserID != 200 || got.Role != models.SocietyMemberRoleAdmin || got.Status != models.SocietyMemberStatusActive {
				t.Fatalf("member = %+v, want active admin user 200", got)
			}
		})
	}

	for _, role := range []models.SocietyMemberRole{models.SocietyMemberRoleStaff, models.SocietyMemberRoleResident} {
		t.Run("add rejected for "+string(role), func(t *testing.T) {
			svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, role, models.SocietyMemberStatusActive), targetUser)

			_, err := svc.AddMember(ctx, models.AddSocietyMemberRequest{SocietyID: 10, UserID: 200, Role: models.SocietyMemberRoleStaff}, 100)
			requireSocietyAppErrCode(t, err, ErrForbiddenSociety.Code)
		})
	}

	t.Run("quota error is propagated", func(t *testing.T) {
		quotaErr := models.NewAppError("QUOTA_TEST", "quota failed", 402, nil)
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive), targetUser)
		svc.subscriptionSvc = &guardSubscriptionSvc{staffErr: quotaErr}

		_, err := svc.AddMember(ctx, models.AddSocietyMemberRequest{SocietyID: 10, UserID: 200, Role: models.SocietyMemberRoleStaff}, 100)
		requireSocietyAppErrCode(t, err, quotaErr.Code)
	})
}

func TestChangeMemberRoleBlocksLastOwnerDemotion(t *testing.T) {
	ctx := context.Background()
	memberRepo := guardMemberRepoWith(
		guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive),
		guardMember(10, 200, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive),
	)
	memberRepo.activeOwnerCount = 1
	svc := &SocietySvc{
		societyRepo: &guardSocietyRepo{society: guardSociety(10, models.SocietyStatusActive)},
		memberRepo:  memberRepo,
	}

	_, err := svc.ChangeMemberRole(ctx, models.ChangeSocietyMemberRoleRequest{SocietyID: 10, UserID: 200, Role: models.SocietyMemberRoleAdmin}, 100)
	requireSocietyAppErrCode(t, err, ErrOwnerProtection.Code)
}

func TestCreateGuardUsesRoleGuardAndDuplicateChecks(t *testing.T) {
	ctx := context.Background()
	req := models.CreateGuardRequest{FirstName: "Gate", LastName: "Staff", Email: "GUARD@EXAMPLE.COM", PhoneNumber: "9999999999", Password: "password123"}

	t.Run("owner can create guard", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleOwner, models.SocietyMemberStatusActive), nil)
		svc.userRepo = &guardUserRepo{nextID: 700}
		svc.txManager = guardTxManager{}

		got, err := svc.CreateGuard(ctx, 10, 100, req)
		if err != nil {
			t.Fatalf("CreateGuard() error = %v", err)
		}
		if got.User == nil || got.User.ID == 0 || got.Member == nil || got.Member.Role != models.SocietyMemberRoleStaff {
			t.Fatalf("guard response = %+v, want created staff user/member", got)
		}
	})

	t.Run("staff cannot create guard", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleStaff, models.SocietyMemberStatusActive), nil)

		_, err := svc.CreateGuard(ctx, 10, 100, req)
		requireSocietyAppErrCode(t, err, ErrForbiddenSociety.Code)
	})

	t.Run("duplicate email is mapped", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive), nil)
		svc.userRepo = &guardUserRepo{emailExists: true}

		_, err := svc.CreateGuard(ctx, 10, 100, req)
		requireSocietyAppErrCode(t, err, ErrDuplicateGuardEmail.Code)
	})

	t.Run("duplicate phone is mapped", func(t *testing.T) {
		svc := guardSocietySvc(guardSociety(10, models.SocietyStatusActive), guardMember(10, 100, models.SocietyMemberRoleAdmin, models.SocietyMemberStatusActive), nil)
		svc.userRepo = &guardUserRepo{phoneExists: true}

		_, err := svc.CreateGuard(ctx, 10, 100, req)
		requireSocietyAppErrCode(t, err, ErrDuplicateGuardPhone.Code)
	})
}

func guardSocietySvc(society *models.Society, actor *models.SocietyMember, user *models.User) *SocietySvc {
	return &SocietySvc{
		societyRepo: &guardSocietyRepo{society: society, update: society},
		memberRepo:  guardMemberRepoWith(actor),
		userRepo:    &guardUserRepo{user: user, nextID: 1000},
		txManager:   guardTxManager{},
	}
}

//nolint:unparam // Helper keeps call sites explicit about the society under test.
func guardSociety(id int64, status models.SocietyStatus) *models.Society {
	return &models.Society{ID: id, Name: "Test Society", SocietyCode: "TS001", Status: status, CreatedBy: 1}
}

//nolint:unparam // Helper keeps call sites explicit about the society under test.
func guardMember(societyID int64, userID int64, role models.SocietyMemberRole, status models.SocietyMemberStatus) *models.SocietyMember {
	return &models.SocietyMember{ID: userID + 1000, SocietyID: societyID, UserID: userID, Role: role, Status: status}
}

func guardMemberRepoWith(members ...*models.SocietyMember) *guardMemberRepo {
	repo := &guardMemberRepo{members: map[[2]int64]*models.SocietyMember{}, activeOwnerCount: 2}
	for _, member := range members {
		if member != nil {
			repo.members[[2]int64{member.SocietyID, member.UserID}] = member
		}
	}
	return repo
}

func requireSocietyAppErrCode(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error code %s, got nil", want)
	}
	var appErr *models.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError code %s, got %T: %v", want, err, err)
	}
	if appErr.Code != want {
		t.Fatalf("error code = %s, want %s", appErr.Code, want)
	}
}

var (
	_ repository.SocietyRepository       = (*guardSocietyRepo)(nil)
	_ repository.SocietyMemberRepository = (*guardMemberRepo)(nil)
	_ repository.UserRepository          = (*guardUserRepo)(nil)
	_ repository.TransactionManager      = guardTxManager{}
)
