package visitorsettingsvc

import (
	"context"

	"go-server/internal/models"
	service "go-server/internal/services"
)

func (s *VisitorSettingSvc) ResolveApprovalRequirement(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, source models.VisitorEntrySource) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	defer cancel()

	if !purpose.IsValid() || !source.IsValid() {
		return false, ErrInvalidVisitorSettings
	}
	societySettings, err := s.settingsRepo.GetSociety(ctx, societyID)
	if err != nil {
		return false, err
	}
	if societySettings == nil {
		return false, ErrVisitorSettingsNotFound
	}
	if !societySettings.IsActive {
		return false, ErrVisitorSettingsInactive
	}
	if err := ensureSourceEnabled(societySettings, source); err != nil {
		return false, err
	}
	if source == models.VisitorEntrySourceResidentLink {
		return false, nil
	}

	switch societySettings.ApprovalMode {
	case models.VisitorApprovalModeMandatory:
		return true, nil
	case models.VisitorApprovalModeOptional:
		return false, nil
	case models.VisitorApprovalModeHybrid:
		flatSettings, err := s.settingsRepo.GetFlatPurpose(ctx, societyID, flatID, purpose)
		if err != nil {
			return false, err
		}
		if flatSettings == nil {
			return false, ErrVisitorSettingsNotFound
		}
		if !flatSettings.IsEnabled {
			return false, ErrVisitorPurposeDisabled
		}
		return flatSettings.ApprovalRequired, nil
	default:
		return false, ErrInvalidVisitorSettings
	}
}

func ensureSourceEnabled(settings *models.SocietyVisitorSettings, source models.VisitorEntrySource) error {
	switch source {
	case models.VisitorEntrySourceResidentLink:
		if !settings.AllowResidentPreApproval {
			return ErrVisitorSourceDisabled
		}
	case models.VisitorEntrySourcePublicQR:
		if !settings.AllowPublicQREntry {
			return ErrVisitorSourceDisabled
		}
	case models.VisitorEntrySourceQuickLink:
		if !settings.AllowPublicQREntry {
			return ErrVisitorSourceDisabled
		}
	case models.VisitorEntrySourceGuardEntry:
		if !settings.AllowGuardEntry {
			return ErrVisitorSourceDisabled
		}
	default:
		return ErrInvalidVisitorSettings
	}
	return nil
}
