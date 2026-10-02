package contracts

import (
	"context"
	"go-server/internal/models"
	"time"
)

type FlatMemberInviteRepository interface {
	Create(ctx context.Context, societyID int64, flatID int64, invitedBy int64, role models.FlatMemberInviteRole, phone, email *string, fullName, tokenHash string, expiresAt time.Time) (*models.FlatMemberInvite, error)
	GetByID(ctx context.Context, societyID int64, inviteID int64) (*models.FlatMemberInvite, error)
	GetHistoryByID(ctx context.Context, societyID int64, flatID int64, inviteID int64, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) (*models.FlatMemberInviteResponse, error)
	GetByIDAny(ctx context.Context, inviteID int64) (*models.FlatMemberInvite, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.FlatMemberInvite, error)
	ListPending(ctx context.Context, societyID int64, flatID int64) ([]*models.FlatMemberInvite, error)
	ListHistory(ctx context.Context, societyID int64, flatID int64, filter models.FlatMemberInviteHistoryFilter, linkBuilder func(string, *time.Time) *models.ShortLinkResponse) ([]*models.FlatMemberInviteResponse, error)
	CountHistory(ctx context.Context, societyID int64, flatID int64, filter models.FlatMemberInviteHistoryFilter) (int64, error)
	Cancel(ctx context.Context, societyID int64, flatID int64, inviteID int64) (*models.FlatMemberInvite, error)
	Accept(ctx context.Context, inviteID int64) (*models.FlatMemberInvite, error)
	ExpireOld(ctx context.Context) error
}
