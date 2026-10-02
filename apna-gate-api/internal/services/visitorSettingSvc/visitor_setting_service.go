package visitorsettingsvc

import (
	"context"

	"go-server/internal/models"
	repository "go-server/internal/repositories/contracts"
	flatauthz "go-server/internal/services/flatAuthz"
)

type VisitorSettingService interface {
	CreateDefaultSocietySettings(ctx context.Context, societyID int64, actorUserID int64) error
	GetSocietySettings(ctx context.Context, societyID int64) (*models.SocietyVisitorSettingsResponse, error)
	UpdateSocietySettings(ctx context.Context, societyID int64, req models.UpdateSocietyVisitorSettingsRequest, actorUserID int64) (*models.SocietyVisitorSettingsResponse, error)
	CreateDefaultFlatSettings(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
	EnsureDefaultFlatSettingsIfMissing(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
	GetFlatSettings(ctx context.Context, societyID int64, flatID int64) ([]models.FlatVisitorSettingsResponse, error)
	GetFlatSettingsForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) ([]models.FlatVisitorSettingsResponse, error)
	UpdateFlatPurposeSetting(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, req models.UpdateFlatVisitorSettingRequest, actorUserID int64) (*models.FlatVisitorSettingsResponse, error)
	ResetFlatSettingsToDefault(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error
	ResolveApprovalRequirement(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, source models.VisitorEntrySource) (bool, error)
	ResolveVisitDurationMinutes(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose) (int32, error)
	ListSocietyFlatSettings(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) (*models.SocietyFlatVisitorSettingsListResult, error)
}

type VisitorSettingSvc struct {
	settingsRepo repository.VisitorSettingRepository
	memberRepo   repository.SocietyMemberRepository
	residentRepo repository.FlatResidentRepository
	flatRepo     repository.FlatRepository
	flatAuthz    *flatauthz.FlatVisitorAuthz
	txManager    repository.TransactionManager
}

func NewVisitorSettingService(
	settingsRepo repository.VisitorSettingRepository,
	memberRepo repository.SocietyMemberRepository,
	residentRepo repository.FlatResidentRepository,
	flatRepo repository.FlatRepository,
	flatAuthz *flatauthz.FlatVisitorAuthz,
	txManager repository.TransactionManager,
) VisitorSettingService {
	return &VisitorSettingSvc{
		settingsRepo: settingsRepo,
		memberRepo:   memberRepo,
		residentRepo: residentRepo,
		flatRepo:     flatRepo,
		flatAuthz:    flatAuthz,
		txManager:    txManager,
	}
}

var _ VisitorSettingService = (*VisitorSettingSvc)(nil)
