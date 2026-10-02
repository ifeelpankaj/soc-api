package repository

import (
	"context"
	"errors"
	"go-server/internal/repositories/contracts"

	"go-server/internal/db"
	"go-server/internal/models"
	"go-server/pkg/database"

	"github.com/jackc/pgx/v5"
)

type VisitorSettingRepository = contracts.VisitorSettingRepository

type visitorSettingRepository struct {
	db *database.Database
}

func NewVisitorSettingRepository(db *database.Database) VisitorSettingRepository {
	return &visitorSettingRepository{db: db}
}

func (r *visitorSettingRepository) CreateDefaultSociety(ctx context.Context, societyID int64, actorUserID int64) error {
	return GetQueries(ctx, r.db).CreateDefaultSocietyVisitorSettings(ctx, db.CreateDefaultSocietyVisitorSettingsParams{
		SocietyID: societyID,
		UpdatedBy: &actorUserID,
	})
}

func (r *visitorSettingRepository) GetSociety(ctx context.Context, societyID int64) (*models.SocietyVisitorSettings, error) {
	row, err := GetQueries(ctx, r.db).GetSocietyVisitorSettings(ctx, societyID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyVisitorSettingsFromDB(row), nil
}

func (r *visitorSettingRepository) UpdateSociety(ctx context.Context, societyID int64, req contracts.UpdateSocietyVisitorSettingsInput, actorUserID int64) (*models.SocietyVisitorSettings, error) {
	row, err := GetQueries(ctx, r.db).UpdateSocietyVisitorSettings(ctx, db.UpdateSocietyVisitorSettingsParams{
		ApprovalMode:                dbVisitorApprovalModePtr(req.ApprovalMode),
		DefaultVisitDurationMinutes: req.DefaultVisitDurationMinutes,
		GracePeriodMinutes:          req.GracePeriodMinutes,
		QrExpiryMinutes:             req.QRExpiryMinutes,
		AllowResidentPreApproval:    req.AllowResidentPreApproval,
		AllowPublicQrEntry:          req.AllowPublicQREntry,
		AllowGuardEntry:             req.AllowGuardEntry,
		AllowGuardOnBehalfApproval:  req.AllowGuardOnBehalfApproval,
		IsActive:                    req.IsActive,
		UpdatedBy:                   &actorUserID,
		SocietyID:                   societyID,
	})
	return societyVisitorSettingsFromDBNoRows(row, err)
}

func (r *visitorSettingRepository) CreateDefaultFlat(ctx context.Context, societyID int64, flatID int64, actorUserID int64) error {
	return GetQueries(ctx, r.db).CreateDefaultFlatVisitorSettings(ctx, db.CreateDefaultFlatVisitorSettingsParams{
		SocietyID: societyID,
		FlatID:    flatID,
		UpdatedBy: &actorUserID,
	})
}

func (r *visitorSettingRepository) ListFlat(ctx context.Context, societyID int64, flatID int64) ([]*models.FlatVisitorSettings, error) {
	rows, err := GetQueries(ctx, r.db).ListFlatVisitorSettings(ctx, db.ListFlatVisitorSettingsParams{SocietyID: societyID, FlatID: flatID})
	if err != nil {
		return nil, err
	}
	items := make([]*models.FlatVisitorSettings, 0, len(rows))
	for _, row := range rows {
		items = append(items, flatVisitorSettingsFromDB(row))
	}
	return items, nil
}

func (r *visitorSettingRepository) ListSocietyFlat(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) ([]*models.SocietyFlatVisitorSettingRow, error) {
	rows, err := GetQueries(ctx, r.db).ListSocietyFlatVisitorSettings(ctx, db.ListSocietyFlatVisitorSettingsParams{
		SocietyID: filter.SocietyID,
		FlatID:    filter.FlatID,
		Block:     filter.Block,
		Purpose:   dbVisitorPurposePtr(filter.Purpose),
		Limit:     normalizeVisitorLimit(filter.Limit),
		Offset:    normalizeOffset(filter.Offset),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*models.SocietyFlatVisitorSettingRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &models.SocietyFlatVisitorSettingRow{
			FlatID:                      row.FlatID,
			FlatNumber:                  row.FlatNumber,
			Block:                       row.Block,
			Purpose:                     models.VisitorPurpose(row.Purpose),
			ApprovalRequired:            row.ApprovalRequired,
			IsEnabled:                   row.IsEnabled,
			DefaultVisitDurationMinutes: row.DefaultVisitDurationMinutes,
		})
	}
	return items, nil
}

func (r *visitorSettingRepository) CountSocietyFlat(ctx context.Context, filter models.SocietyFlatVisitorSettingsFilter) (int64, error) {
	return GetQueries(ctx, r.db).CountSocietyFlatVisitorSettings(ctx, db.CountSocietyFlatVisitorSettingsParams{
		SocietyID: filter.SocietyID,
		FlatID:    filter.FlatID,
		Block:     filter.Block,
		Purpose:   dbVisitorPurposePtr(filter.Purpose),
	})
}

func (r *visitorSettingRepository) GetFlatPurpose(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose) (*models.FlatVisitorSettings, error) {
	row, err := GetQueries(ctx, r.db).GetFlatVisitorPurposeSetting(ctx, db.GetFlatVisitorPurposeSettingParams{
		SocietyID: societyID,
		FlatID:    flatID,
		Purpose:   db.VisitorPurpose(purpose),
	})
	return flatVisitorSettingsFromDBNoRows(row, err)
}

func (r *visitorSettingRepository) UpdateFlatPurpose(ctx context.Context, societyID int64, flatID int64, purpose models.VisitorPurpose, req contracts.UpdateFlatVisitorSettingInput, actorUserID int64) (*models.FlatVisitorSettings, error) {
	row, err := GetQueries(ctx, r.db).UpdateFlatVisitorPurposeSetting(ctx, db.UpdateFlatVisitorPurposeSettingParams{
		ApprovalRequired:            req.ApprovalRequired,
		DefaultVisitDurationMinutes: req.DefaultVisitDurationMinutes,
		IsEnabled:                   req.IsEnabled,
		UpdatedBy:                   &actorUserID,
		SocietyID:                   societyID,
		FlatID:                      flatID,
		Purpose:                     db.VisitorPurpose(purpose),
	})
	return flatVisitorSettingsFromDBNoRows(row, err)
}

func (r *visitorSettingRepository) DeleteFlat(ctx context.Context, societyID int64, flatID int64) error {
	return GetQueries(ctx, r.db).DeleteFlatVisitorSettings(ctx, db.DeleteFlatVisitorSettingsParams{SocietyID: societyID, FlatID: flatID})
}

func societyVisitorSettingsFromDBNoRows(row db.SocietyVisitorSetting, err error) (*models.SocietyVisitorSettings, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return societyVisitorSettingsFromDB(row), nil
}

func societyVisitorSettingsFromDB(row db.SocietyVisitorSetting) *models.SocietyVisitorSettings {
	return &models.SocietyVisitorSettings{
		ID: row.ID, SocietyID: row.SocietyID, ApprovalMode: models.VisitorApprovalMode(row.ApprovalMode),
		DefaultVisitDurationMinutes: row.DefaultVisitDurationMinutes, GracePeriodMinutes: row.GracePeriodMinutes,
		QRExpiryMinutes: row.QrExpiryMinutes, AllowResidentPreApproval: row.AllowResidentPreApproval,
		AllowPublicQREntry: row.AllowPublicQrEntry, AllowGuardEntry: row.AllowGuardEntry,
		AllowGuardOnBehalfApproval: row.AllowGuardOnBehalfApproval, IsActive: row.IsActive,
		UpdatedBy: row.UpdatedBy, CreatedAt: pgTimestamptzToTime(row.CreatedAt), UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func flatVisitorSettingsFromDBNoRows(row db.FlatVisitorSetting, err error) (*models.FlatVisitorSettings, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return flatVisitorSettingsFromDB(row), nil
}

func flatVisitorSettingsFromDB(row db.FlatVisitorSetting) *models.FlatVisitorSettings {
	return &models.FlatVisitorSettings{
		ID: row.ID, SocietyID: row.SocietyID, FlatID: row.FlatID, Purpose: models.VisitorPurpose(row.Purpose),
		ApprovalRequired: row.ApprovalRequired, DefaultVisitDurationMinutes: row.DefaultVisitDurationMinutes,
		IsEnabled: row.IsEnabled, UpdatedBy: row.UpdatedBy, CreatedAt: pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt: pgTimestamptzToTime(row.UpdatedAt),
	}
}

func dbVisitorApprovalModePtr(mode *models.VisitorApprovalMode) *db.VisitorApprovalMode {
	if mode == nil || *mode == "" {
		return nil
	}
	value := db.VisitorApprovalMode(*mode)
	return &value
}
