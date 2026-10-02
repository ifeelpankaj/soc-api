package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type idempotencyRepository struct{ database *database.Database }

func (r *idempotencyRepository) GetUPIIdempotency(ctx context.Context, arg contracts.GetUPIIdempotencyInput) (contracts.MaintenancePaymentIdempotency, error) {
	value, err := GetQueries(ctx, r.database).GetUPIIdempotency(ctx, db.GetUPIIdempotencyParams{SocietyID: arg.SocietyID, ActorID: arg.ActorID, Operation: arg.Operation, Key: arg.Key})
	return mapFinancialMaintenancePaymentIdempotency(value), persistenceError(err)
}

func (r *idempotencyRepository) InsertUPIIdempotency(ctx context.Context, arg contracts.InsertUPIIdempotencyInput) error {
	return persistenceError(GetQueries(ctx, r.database).InsertUPIIdempotency(ctx, db.InsertUPIIdempotencyParams{SocietyID: arg.SocietyID, ActorID: arg.ActorID, Operation: arg.Operation, Key: arg.Key, RequestHash: arg.RequestHash, Response: arg.Response}))
}

var _ contracts.IdempotencyRepository = (*idempotencyRepository)(nil)
