package visitorsettingsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) GetFlatSettings(ctx context.Context, societyID int64, flatID int64) ([]models.FlatVisitorSettingsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	items, err := s.settingsRepo.ListFlat(ctx, societyID, flatID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrVisitorSettingsNotFound
	}
	responses := make([]models.FlatVisitorSettingsResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, *item.ToResponse())
	}
	return responses, nil
}

func (s *VisitorSettingSvc) GetFlatSettingsForActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) ([]models.FlatVisitorSettingsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if err := s.ensureFlatSettingsActor(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	if err := s.EnsureDefaultFlatSettingsIfMissing(ctx, societyID, flatID, actorUserID); err != nil {
		return nil, err
	}
	return s.GetFlatSettings(ctx, societyID, flatID)
}
