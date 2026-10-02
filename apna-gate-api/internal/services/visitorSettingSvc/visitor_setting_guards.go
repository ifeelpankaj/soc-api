package visitorsettingsvc

import (
	"context"
	"errors"

	"go-server/internal/models"
	flatauthz "go-server/internal/services/flatAuthz"
)

func (s *VisitorSettingSvc) ensureFlatSettingsActor(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	return mapFlatAuthzError(s.flatAuthz.CanManageFlatVisitors(ctx, societyID, flatID, actorUserID))
}

func mapFlatAuthzError(err error) error {
	if err == nil {
		return nil
	}

	var appErr *models.AppError
	if !errors.As(err, &appErr) {
		return err
	}

	switch appErr.Code {
	case flatauthz.ErrForbidden.Code:
		return ErrVisitorSettingsForbidden
	case flatauthz.ErrFlatNotFound.Code:
		return ErrVisitorSettingsNotFound
	default:
		return err
	}
}
