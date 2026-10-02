package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
)

func (r *HubRepository) HubActorLock(ctx context.Context, arg contracts.HubActorLockInput) error {
	return persistenceError(GetQueries(ctx, r.database).HubActorLock(ctx, db.HubActorLockParams{SocietyID: arg.SocietyID, UserID: arg.UserID}))
}

func (r *HubRepository) HubResidents(ctx context.Context, societyID int64) ([]int64, error) {
	value, err := GetQueries(ctx, r.database).HubResidents(ctx, societyID)
	return value, persistenceError(err)
}

func (r *HubRepository) HubRole(ctx context.Context, arg contracts.HubRoleInput) (string, error) {
	value, err := GetQueries(ctx, r.database).HubRole(ctx, db.HubRoleParams{SocietyID: arg.SocietyID, UserID: arg.UserID})
	return value, persistenceError(err)
}
