package visitorsettingsvc

import (
	"context"
	"go-server/internal/repositories/contracts"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) CreateDefaultSocietySettings(ctx context.Context, societyID int64, actorUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	return s.settingsRepo.CreateDefaultSociety(ctx, societyID, actorUserID)
}

func (s *VisitorSettingSvc) UpdateSocietySettings(ctx context.Context, societyID int64, req models.UpdateSocietyVisitorSettingsRequest, actorUserID int64) (*models.SocietyVisitorSettingsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := req.Validate(); err != nil {
		return nil, ErrInvalidVisitorSettings.WithCause(err)
	}
	settings, err := s.settingsRepo.UpdateSociety(ctx, societyID, contracts.UpdateSocietyVisitorSettingsInput(req), actorUserID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, ErrVisitorSettingsNotFound
	}
	return settings.ToResponse(), nil
}
