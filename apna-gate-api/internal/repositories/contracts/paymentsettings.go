package contracts

import (
	"context"
)

type PaymentSettingsRepository interface {
	GetUPISettings(ctx context.Context, societyID int64) (MaintenancePaymentSettingsVersion, error)
	GetUPISettingsVersion(ctx context.Context, arg GetUPISettingsVersionInput) (MaintenancePaymentSettingsVersion, error)
	InsertUPISettingsVersion(ctx context.Context, arg InsertUPISettingsVersionInput) (MaintenancePaymentSettingsVersion, error)
	ListUPISettingsVersions(ctx context.Context, arg ListUPISettingsVersionsInput) ([]MaintenancePaymentSettingsVersion, error)
	PointUPISettings(ctx context.Context, arg PointUPISettingsInput) error
}
