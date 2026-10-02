package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type paymentRequestsRepository struct{ database *database.Database }

func (r *paymentRequestsRepository) CloseUPIRequests(ctx context.Context, arg contracts.CloseUPIRequestsInput) error {
	return persistenceError(GetQueries(ctx, r.database).CloseUPIRequests(ctx, db.CloseUPIRequestsParams{SocietyID: arg.SocietyID, BillID: arg.BillID}))
}

func (r *paymentRequestsRepository) GetActiveUPIRequest(ctx context.Context, arg contracts.GetActiveUPIRequestInput) (contracts.MaintenancePaymentRequest, error) {
	value, err := GetQueries(ctx, r.database).GetActiveUPIRequest(ctx, db.GetActiveUPIRequestParams{SocietyID: arg.SocietyID, BillID: arg.BillID})
	return mapFinancialMaintenancePaymentRequest(value), persistenceError(err)
}

func (r *paymentRequestsRepository) GetUPIRequest(ctx context.Context, arg contracts.GetUPIRequestInput) (contracts.MaintenancePaymentRequest, error) {
	value, err := GetQueries(ctx, r.database).GetUPIRequest(ctx, db.GetUPIRequestParams{SocietyID: arg.SocietyID, ID: pgtype.UUID{Bytes: arg.ID, Valid: arg.ID != uuid.Nil}})
	return mapFinancialMaintenancePaymentRequest(value), persistenceError(err)
}

func (r *paymentRequestsRepository) InsertUPIRequest(ctx context.Context, arg contracts.InsertUPIRequestInput) (contracts.MaintenancePaymentRequest, error) {
	value, err := GetQueries(ctx, r.database).InsertUPIRequest(ctx, db.InsertUPIRequestParams{ID: pgtype.UUID{Bytes: arg.ID, Valid: arg.ID != uuid.Nil}, SocietyID: arg.SocietyID, BillID: arg.BillID, SettingsVersion: arg.SettingsVersion, AmountPaise: arg.AmountPaise, Reference: arg.Reference, CreatedBy: arg.CreatedBy})
	return mapFinancialMaintenancePaymentRequest(value), persistenceError(err)
}

func (r *paymentRequestsRepository) SupersedeUPIRequests(ctx context.Context, societyID int64) error {
	return persistenceError(GetQueries(ctx, r.database).SupersedeUPIRequests(ctx, societyID))
}

var _ contracts.PaymentRequestsRepository = (*paymentRequestsRepository)(nil)
