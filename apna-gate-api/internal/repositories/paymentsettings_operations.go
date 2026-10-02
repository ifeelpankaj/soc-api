package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type paymentSettingsRepository struct{ database *database.Database }

func (r *paymentSettingsRepository) GetUPISettings(ctx context.Context, societyID int64) (contracts.MaintenancePaymentSettingsVersion, error) {
	value, err := GetQueries(ctx, r.database).GetUPISettings(ctx, societyID)
	return mapFinancialMaintenancePaymentSettingsVersion(value), persistenceError(err)
}

func (r *paymentSettingsRepository) GetUPISettingsVersion(ctx context.Context, arg contracts.GetUPISettingsVersionInput) (contracts.MaintenancePaymentSettingsVersion, error) {
	value, err := GetQueries(ctx, r.database).GetUPISettingsVersion(ctx, db.GetUPISettingsVersionParams{SocietyID: arg.SocietyID, Version: arg.Version})
	return mapFinancialMaintenancePaymentSettingsVersion(value), persistenceError(err)
}

func (r *paymentSettingsRepository) InsertUPISettingsVersion(ctx context.Context, arg contracts.InsertUPISettingsVersionInput) (contracts.MaintenancePaymentSettingsVersion, error) {
	value, err := GetQueries(ctx, r.database).InsertUPISettingsVersion(ctx, db.InsertUPISettingsVersionParams{SocietyID: arg.SocietyID, Version: arg.Version, Enabled: arg.Enabled, UpiID: arg.UpiID, PayeeName: arg.PayeeName, CreatedBy: arg.CreatedBy})
	return mapFinancialMaintenancePaymentSettingsVersion(value), persistenceError(err)
}

func (r *paymentSettingsRepository) ListUPISettingsVersions(ctx context.Context, arg contracts.ListUPISettingsVersionsInput) ([]contracts.MaintenancePaymentSettingsVersion, error) {
	value, err := GetQueries(ctx, r.database).ListUPISettingsVersions(ctx, db.ListUPISettingsVersionsParams{SocietyID: arg.SocietyID, BeforeVersion: arg.BeforeVersion, Limit: arg.Limit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.MaintenancePaymentSettingsVersion, len(value))
	for i := range value {
		out[i] = mapFinancialMaintenancePaymentSettingsVersion(value[i])
	}
	return out, nil
}

func (r *paymentSettingsRepository) PointUPISettings(ctx context.Context, arg contracts.PointUPISettingsInput) error {
	return persistenceError(GetQueries(ctx, r.database).PointUPISettings(ctx, db.PointUPISettingsParams{SocietyID: arg.SocietyID, Version: arg.Version}))
}

var _ contracts.PaymentSettingsRepository = (*paymentSettingsRepository)(nil)
