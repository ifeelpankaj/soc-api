package visitorsettingsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) GetSocietySettings(ctx context.Context, societyID int64) (*models.SocietyVisitorSettingsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	settings, err := s.settingsRepo.GetSociety(ctx, societyID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, ErrVisitorSettingsNotFound
	}
	return settings.ToResponse(), nil
}

func (s *VisitorSettingSvc) ListSocietyFlatSettings(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) (*models.SocietyFlatVisitorSettingsListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	settings, err := s.settingsRepo.ListSocietyFlat(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.settingsRepo.CountSocietyFlat(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &models.SocietyFlatVisitorSettingsListResult{
		Settings: settings,
		Total:    total,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}
