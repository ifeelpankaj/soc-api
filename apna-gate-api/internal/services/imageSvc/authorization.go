package imagesvc

import (
	"context"
	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
)

type OperationalGuard interface {
	EnsureSocietyOperational(context.Context, int64) error
}
type FlatEntryReader interface {
	GetFlatEntryForActor(context.Context, int64, int64, int64, int64) (*models.VisitorEntry, error)
}
type Authorizer interface {
	Authorize(context.Context, models.ImageTarget, bool) (models.ImageTarget, error)
}
type authorization struct {
	users       repository.UserRepository
	members     repository.SocietyMemberRepository
	entries     repository.VisitorEntryRepository
	operational OperationalGuard
	flats       FlatEntryReader
}

func NewAuthorization(users repository.UserRepository, members repository.SocietyMemberRepository, entries repository.VisitorEntryRepository, operational OperationalGuard, flats FlatEntryReader) Authorizer {
	return &authorization{users, members, entries, operational, flats}
}
func (a *authorization) Authorize(ctx context.Context, t models.ImageTarget, upload bool) (models.ImageTarget, error) {
	user, err := a.users.GetByID(ctx, t.ActorID)
	if err != nil {
		return t, err
	}
	if user == nil || !user.IsActive || user.IsBlocked || user.DeletedAt != nil {
		return t, ErrImageForbidden
	}
	if t.Avatar() {
		return t, nil
	}
	if err = a.operational.EnsureSocietyOperational(ctx, t.SocietyID); err != nil {
		return t, err
	}
	var entry *models.VisitorEntry
	if t.FlatID != 0 {
		if upload {
			return t, ErrImageForbidden
		}
		entry, err = a.flats.GetFlatEntryForActor(ctx, t.SocietyID, t.FlatID, t.EntryID, t.ActorID)
	} else {
		active := string(models.SocietyMemberStatusActive)
		member, memberErr := a.members.Get(ctx, models.GetSocietyMemberFilter{SocietyID: &t.SocietyID, UserID: &t.ActorID, Status: &active})
		if memberErr != nil {
			return t, memberErr
		}
		if member == nil {
			return t, ErrImageForbidden
		}
		switch member.Role {
		case models.SocietyMemberRoleOwner, models.SocietyMemberRoleAdmin, models.SocietyMemberRoleStaff:
		default:
			return t, ErrImageForbidden
		}
		entry, err = a.entries.Get(ctx, t.SocietyID, t.EntryID)
	}
	if err != nil {
		return t, err
	}
	if entry == nil {
		return t, ErrImageTargetNotFound
	}
	if upload {
		switch entry.Status {
		case models.VisitorStatusWaitingApproval, models.VisitorStatusApproved, models.VisitorStatusCheckedIn:
		default:
			return t, ErrImageState
		}
	}
	t.VisitorID = entry.VisitorID
	return t, nil
}
