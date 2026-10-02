package visitorsettingsvc

import (
	"net/http"

	"go-server/internal/models"
)

var (
	ErrVisitorSettingsNotFound  = models.NewAppError("VISITOR_SETTINGS_NOT_FOUND", "visitor settings not found", http.StatusNotFound, nil)
	ErrInvalidVisitorSettings   = models.NewAppError("INVALID_VISITOR_SETTINGS", "invalid visitor settings request", http.StatusBadRequest, nil)
	ErrVisitorSettingsForbidden = models.NewAppError("VISITOR_SETTINGS_FORBIDDEN", "only an active flat resident or society admin can manage visitor settings", http.StatusForbidden, nil)
	ErrVisitorSourceDisabled    = models.NewAppError("VISITOR_SOURCE_DISABLED", "visitor entry source is disabled", http.StatusForbidden, nil)
	ErrVisitorPurposeDisabled   = models.NewAppError("VISITOR_PURPOSE_DISABLED", "visitor purpose is disabled", http.StatusConflict, nil)
	ErrVisitorSettingsInactive  = models.NewAppError("VISITOR_SETTINGS_INACTIVE", "visitor settings are inactive", http.StatusConflict, nil)
)
