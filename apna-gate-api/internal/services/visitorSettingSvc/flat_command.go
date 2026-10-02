package visitorsettingsvc

import (
	"context"
	"go-server/internal/repositories/contracts"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) CreateDefaultFlatSettings(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.settingsRepo.CreateDefaultFlat(ctx, societyID, flatID, actorUserID)
}

func (s *VisitorSettingSvc) EnsureDefaultFlatSettingsIfMissing(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.settingsRepo.ListFlat(ctx, societyID, flatID)
	if err != nil {
		return err
	}
	if len(items) > 0 {
		return nil
	}

	return s.settingsRepo.CreateDefaultFlat(ctx, societyID, flatID, actorUserID)
}

func (s *VisitorSettingSvc) UpdateFlatPurposeSetting(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, req models.UpdateFlatVisitorSettingRequest, actorUserID int64) (*models.FlatVisitorSettingsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if !purpose.IsValid() {
		return nil, ErrInvalidVisitorSettings
	}
	if err := req.Validate(); err != nil {
		return nil, ErrInvalidVisitorSettings.WithCause(err)
	}
	if err := s.ensureFlatSettingsActor(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	settings, err := s.settingsRepo.UpdateFlatPurpose(ctx, societyID, flatID, purpose, contracts.UpdateFlatVisitorSettingInput(req), actorUserID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, ErrVisitorSettingsNotFound
	}
	return settings.ToResponse(), nil
}

func (s *VisitorSettingSvc) ResetFlatSettingsToDefault(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatSettingsActor(ctx, societyID, flatID, actorUserID); err != nil {
		return err
	}
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.settingsRepo.DeleteFlat(txCtx, societyID, flatID); err != nil {
			return err
		}
		return s.settingsRepo.CreateDefaultFlat(txCtx, societyID, flatID, actorUserID)
	})
}
