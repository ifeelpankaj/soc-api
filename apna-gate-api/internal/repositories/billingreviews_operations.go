package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type billingReviewsRepository struct{ database *database.Database }

func (r *billingReviewsRepository) GetMaintenanceReview(ctx context.Context, arg contracts.GetMaintenanceReviewInput) (contracts.GetMaintenanceReviewRecord, error) {
	value, err := GetQueries(ctx, r.database).GetMaintenanceReview(ctx, db.GetMaintenanceReviewParams{Token: pgtype.UUID{Bytes: arg.Token, Valid: arg.Token != uuid.Nil}, SocietyID: arg.SocietyID, ActorID: arg.ActorID, BillingMonth: pgtype.Date{Time: arg.BillingMonth, Valid: !arg.BillingMonth.IsZero()}})
	return mapFinancialGetMaintenanceReviewRow(value), persistenceError(err)
}

func (r *billingReviewsRepository) SaveMaintenanceReview(ctx context.Context, arg contracts.SaveMaintenanceReviewInput) error {
	return persistenceError(GetQueries(ctx, r.database).SaveMaintenanceReview(ctx, db.SaveMaintenanceReviewParams{Token: pgtype.UUID{Bytes: arg.Token, Valid: arg.Token != uuid.Nil}, SocietyID: arg.SocietyID, ActorID: arg.ActorID, BillingMonth: pgtype.Date{Time: arg.BillingMonth, Valid: !arg.BillingMonth.IsZero()}, SnapshotHash: arg.SnapshotHash, ExpiresAt: timestampParam(arg.ExpiresAt)}))
}

var _ contracts.BillingReviewsRepository = (*billingReviewsRepository)(nil)
