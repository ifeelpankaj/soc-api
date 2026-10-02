package contracts

import (
	"context"
	"go-server/internal/models"
)

type SocietyMemberRepository interface {
	Add(ctx context.Context, member *models.SocietyMember) error
	Get(ctx context.Context, filter models.GetSocietyMemberFilter) (*models.SocietyMember, error)
	List(ctx context.Context, filter models.ListSocietyMembersFilter) ([]*models.SocietyMember, error)
	ListByUser(ctx context.Context, userID int64) ([]*models.SocietyMember, error)
	ListMySocietiesByUser(ctx context.Context, userID int64) ([]*models.MySocietyResponse, error)
	Count(ctx context.Context, filter models.ListSocietyMembersFilter) (int64, error)
	ChangeRole(ctx context.Context, societyID int64, userID int64, role models.SocietyMemberRole) (*models.SocietyMember, error)
	Suspend(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error)
	Reactivate(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error)
	Remove(ctx context.Context, societyID int64, userID int64, removedBy int64, reason string) error
	CountActiveOwners(ctx context.Context, societyID int64) (int64, error)
	DemoteActiveOwners(ctx context.Context, societyID int64, exceptUserID int64) error
	PromoteToOwner(ctx context.Context, societyID int64, userID int64) (*models.SocietyMember, error)
	UpsertResident(ctx context.Context, societyID int64, userID int64, invitedBy int64) (*models.SocietyMember, error)
}
