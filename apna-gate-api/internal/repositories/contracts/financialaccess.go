package contracts

import (
	"context"
)

type FinancialAccessRepository interface {
	MaintenanceAdmin(ctx context.Context, arg MaintenanceAdminInput) (bool, error)
	MaintenanceLock(ctx context.Context, societyID int64) error
	MaintenanceResidentAccess(ctx context.Context, arg MaintenanceResidentAccessInput) (bool, error)
}
