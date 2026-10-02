package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type financialAccessRepository struct{ database *database.Database }

func (r *financialAccessRepository) MaintenanceAdmin(ctx context.Context, arg contracts.MaintenanceAdminInput) (bool, error) {
	value, err := GetQueries(ctx, r.database).MaintenanceAdmin(ctx, db.MaintenanceAdminParams{SocietyID: arg.SocietyID, UserID: arg.UserID})
	return value, persistenceError(err)
}

func (r *financialAccessRepository) MaintenanceLock(ctx context.Context, societyID int64) error {
	return persistenceError(GetQueries(ctx, r.database).MaintenanceLock(ctx, societyID))
}

func (r *financialAccessRepository) MaintenanceResidentAccess(ctx context.Context, arg contracts.MaintenanceResidentAccessInput) (bool, error) {
	value, err := GetQueries(ctx, r.database).MaintenanceResidentAccess(ctx, db.MaintenanceResidentAccessParams{SocietyID: arg.SocietyID, ID: arg.ID, UserID: arg.UserID})
	return value, persistenceError(err)
}

var _ contracts.FinancialAccessRepository = (*financialAccessRepository)(nil)
