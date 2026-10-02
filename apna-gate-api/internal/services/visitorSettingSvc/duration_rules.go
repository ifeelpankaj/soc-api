package visitorsettingsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) ResolveVisitDurationMinutes(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose) (int32, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	settings, err := s.settingsRepo.GetSociety(ctx, societyID)
	if err != nil {
		return 0, err
	}
	if settings == nil {
		return 0, ErrVisitorSettingsNotFound
	}

	duration := settings.DefaultVisitDurationMinutes
	if flatID > 0 && purpose.IsValid() {
		flatSettings, err := s.settingsRepo.GetFlatPurpose(ctx, societyID, flatID, purpose)
		if err != nil {
			return 0, err
		}
		if flatSettings != nil && flatSettings.DefaultVisitDurationMinutes != nil {
			duration = *flatSettings.DefaultVisitDurationMinutes
		}
	}
	if duration <= 0 {
		return 0, ErrInvalidVisitorSettings
	}
	return duration, nil
}
