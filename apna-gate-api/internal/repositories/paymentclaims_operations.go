package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type paymentClaimsRepository struct{ database *database.Database }

func (r *paymentClaimsRepository) GetPendingUPIClaim(ctx context.Context, arg contracts.GetPendingUPIClaimInput) (contracts.MaintenancePaymentClaim, error) {
	value, err := GetQueries(ctx, r.database).GetPendingUPIClaim(ctx, db.GetPendingUPIClaimParams{SocietyID: arg.SocietyID, BillID: arg.BillID})
	return mapFinancialMaintenancePaymentClaim(value), persistenceError(err)
}

func (r *paymentClaimsRepository) GetUPIClaim(ctx context.Context, arg contracts.GetUPIClaimInput) (contracts.MaintenancePaymentClaim, error) {
	value, err := GetQueries(ctx, r.database).GetUPIClaim(ctx, db.GetUPIClaimParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapFinancialMaintenancePaymentClaim(value), persistenceError(err)
}

func (r *paymentClaimsRepository) InsertUPIClaim(ctx context.Context, arg contracts.InsertUPIClaimInput) (contracts.MaintenancePaymentClaim, error) {
	value, err := GetQueries(ctx, r.database).InsertUPIClaim(ctx, db.InsertUPIClaimParams{SocietyID: arg.SocietyID, BillID: arg.BillID, RequestID: pgtype.UUID{Bytes: arg.RequestID, Valid: arg.RequestID != uuid.Nil}, SubmittedBy: arg.SubmittedBy, Reference: arg.Reference, PaymentDate: pgtype.Date{Time: arg.PaymentDate, Valid: !arg.PaymentDate.IsZero()}})
	return mapFinancialMaintenancePaymentClaim(value), persistenceError(err)
}

func (r *paymentClaimsRepository) ListUPIClaims(ctx context.Context, arg contracts.ListUPIClaimsInput) ([]contracts.MaintenancePaymentClaim, error) {
	value, err := GetQueries(ctx, r.database).ListUPIClaims(ctx, db.ListUPIClaimsParams{SocietyID: arg.SocietyID, UserID: arg.UserID, Status: arg.Status, Reference: arg.Reference, BeforeID: arg.BeforeID, BillID: arg.BillID, FlatID: arg.FlatID, BillingMonth: pgtype.Date{Time: arg.BillingMonth, Valid: !arg.BillingMonth.IsZero()}, Block: arg.Block, FlatNumber: arg.FlatNumber, Search: arg.Search, Limit: arg.Limit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.MaintenancePaymentClaim, len(value))
	for i := range value {
		out[i] = mapFinancialMaintenancePaymentClaim(value[i])
	}
	return out, nil
}

func (r *paymentClaimsRepository) MaintenancePendingClaimsCount(ctx context.Context) (int64, error) {
	value, err := GetQueries(ctx, r.database).MaintenancePendingClaimsCount(ctx)
	return value, persistenceError(err)
}

func (r *paymentClaimsRepository) ReleaseUPIClaimReference(ctx context.Context, arg contracts.ReleaseUPIClaimReferenceInput) error {
	return persistenceError(GetQueries(ctx, r.database).ReleaseUPIClaimReference(ctx, db.ReleaseUPIClaimReferenceParams{SocietyID: arg.SocietyID, ClaimID: arg.ClaimID}))
}

func (r *paymentClaimsRepository) ReviewUPIClaim(ctx context.Context, arg contracts.ReviewUPIClaimInput) (contracts.MaintenancePaymentClaim, error) {
	value, err := GetQueries(ctx, r.database).ReviewUPIClaim(ctx, db.ReviewUPIClaimParams{SocietyID: arg.SocietyID, ID: arg.ID, Status: arg.Status, ReviewedBy: arg.ReviewedBy, Reason: arg.Reason})
	return mapFinancialMaintenancePaymentClaim(value), persistenceError(err)
}

var _ contracts.PaymentClaimsRepository = (*paymentClaimsRepository)(nil)
