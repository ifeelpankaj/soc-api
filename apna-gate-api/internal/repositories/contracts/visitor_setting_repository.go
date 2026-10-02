package contracts

import (
	"context"
	"go-server/internal/models"
)

type VisitorSettingRepository interface {
	CreateDefaultSociety(ctx context.Context, societyID int64, actorUserID int64) error
	GetSociety(ctx context.Context, societyID int64) (*models.SocietyVisitorSettings, error)
	UpdateSociety(ctx context.Context, societyID int64, req UpdateSocietyVisitorSettingsInput, actorUserID int64) (*models.SocietyVisitorSettings, error)
	CreateDefaultFlat(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
	ListFlat(ctx context.Context, societyID int64, flatID int64) ([]*models.FlatVisitorSettings, error)
	ListSocietyFlat(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) ([]*models.SocietyFlatVisitorSettingRow, error)
	CountSocietyFlat(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) (int64, error)
	GetFlatPurpose(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose) (*models.FlatVisitorSettings, error)
	UpdateFlatPurpose(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, req UpdateFlatVisitorSettingInput, actorUserID int64) (*models.FlatVisitorSettings, error)
	DeleteFlat(ctx context.Context, societyID int64, flatID int64) error
}
