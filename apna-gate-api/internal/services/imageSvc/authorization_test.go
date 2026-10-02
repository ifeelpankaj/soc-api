package imagesvc

import (
	"context"
	"errors"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"testing"
)

type imageUsers struct {
	repository.UserRepository
	user *models.User
}

func (r imageUsers) GetByID(context.Context, int64) (*models.User, error) { return r.user, nil }

type imageMembers struct {
	repository.SocietyMemberRepository
	member *models.SocietyMember
}

func (r imageMembers) Get(context.Context, models.GetSocietyMemberFilter) (*models.SocietyMember, error) {
	return r.member, nil
}

type imageEntries struct {
	repository.VisitorEntryRepository
	entry *models.VisitorEntry
}

func (r imageEntries) Get(_ context.Context, society, id int64) (*models.VisitorEntry, error) {
	if r.entry.SocietyID != society || r.entry.ID != id {
		return nil, nil
	}
	return r.entry, nil
}

type imageOperational struct{ err error }

func (r imageOperational) EnsureSocietyOperational(context.Context, int64) error { return r.err }

type imageFlatReader struct {
	called bool
	err    error
	entry  *models.VisitorEntry
}

func (r *imageFlatReader) GetFlatEntryForActor(context.Context, int64, int64, int64, int64) (*models.VisitorEntry, error) {
	r.called = true
	return r.entry, r.err
}

func TestImageAuthorization(t *testing.T) {
	entry := &models.VisitorEntry{ID: 3, SocietyID: 2, VisitorID: 4, Status: models.VisitorStatusApproved}
	for _, tc := range []struct {
		name            string
		active, blocked bool
		member          bool
		role            models.SocietyMemberRole
		target          models.ImageTarget
		upload          bool
		status          models.VisitorStatus
		want            error
	}{
		{name: "own avatar", active: true, target: models.ImageTarget{ActorID: 1}},
		{name: "inactive", target: models.ImageTarget{ActorID: 1}, want: ErrImageForbidden},
		{name: "blocked", active: true, blocked: true, target: models.ImageTarget{ActorID: 1}, want: ErrImageForbidden},
		{name: "staff", active: true, member: true, role: models.SocietyMemberRoleStaff, target: models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3}, upload: true, status: models.VisitorStatusApproved},
		{name: "background upload after check-in", active: true, member: true, role: models.SocietyMemberRoleStaff, target: models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3}, upload: true, status: models.VisitorStatusCheckedIn},
		{name: "cross society", active: true, member: true, role: models.SocietyMemberRoleStaff, target: models.ImageTarget{ActorID: 1, SocietyID: 99, EntryID: 3}, want: ErrImageTargetNotFound},
		{name: "not member", active: true, target: models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3}, want: ErrImageForbidden},
		{name: "terminal upload", active: true, member: true, role: models.SocietyMemberRoleOwner, target: models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3}, upload: true, status: models.VisitorStatusCheckedOut, want: ErrImageState},
		{name: "terminal removal", active: true, member: true, role: models.SocietyMemberRoleAdmin, target: models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3}, status: models.VisitorStatusCheckedOut},
	} {
		t.Run(tc.name, func(t *testing.T) {
			copyEntry := *entry
			copyEntry.Status = tc.status
			var member *models.SocietyMember
			if tc.member {
				member = &models.SocietyMember{Role: tc.role}
			}
			auth := NewAuthorization(imageUsers{user: &models.User{IsActive: tc.active, IsBlocked: tc.blocked}}, imageMembers{member: member}, imageEntries{entry: &copyEntry}, imageOperational{}, &imageFlatReader{})
			_, err := auth.Authorize(context.Background(), tc.target, tc.upload)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
		})
	}
	flat := &imageFlatReader{err: ErrImageForbidden}
	auth := NewAuthorization(imageUsers{user: &models.User{IsActive: true}}, nil, nil, imageOperational{}, flat)
	if _, err := auth.Authorize(context.Background(), models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: 3, FlatID: 5}, false); !errors.Is(err, ErrImageForbidden) || !flat.called {
		t.Fatal("resident authorization not reused")
	}
}
