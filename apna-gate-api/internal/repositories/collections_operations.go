package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type collectionsRepository struct{ database *database.Database }

func (r *collectionsRepository) UserName(ctx context.Context, id int64) (string, error) {
	var name string
	err := r.database.Pool.QueryRow(ctx, "SELECT full_name FROM users WHERE id = $1", id).Scan(&name)
	return name, err
}

func (r *collectionsRepository) GetActiveUPIPayment(ctx context.Context, arg contracts.GetActiveUPIPaymentInput) (contracts.MaintenancePayment, error) {
	value, err := GetQueries(ctx, r.database).GetActiveUPIPayment(ctx, db.GetActiveUPIPaymentParams{SocietyID: arg.SocietyID, BillID: arg.BillID})
	return mapFinancialMaintenancePayment(value), persistenceError(err)
}

func (r *collectionsRepository) GetPaymentBill(ctx context.Context, arg contracts.GetPaymentBillInput) (contracts.MaintenanceBill, error) {
	value, err := GetQueries(ctx, r.database).GetPaymentBill(ctx, db.GetPaymentBillParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapFinancialMaintenanceBill(value), persistenceError(err)
}

func (r *collectionsRepository) GetUPIPayment(ctx context.Context, arg contracts.GetUPIPaymentInput) (contracts.MaintenancePayment, error) {
	value, err := GetQueries(ctx, r.database).GetUPIPayment(ctx, db.GetUPIPaymentParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapFinancialMaintenancePayment(value), persistenceError(err)
}

func (r *collectionsRepository) HasReversedUPIReference(ctx context.Context, arg contracts.HasReversedUPIReferenceInput) (bool, error) {
	value, err := GetQueries(ctx, r.database).HasReversedUPIReference(ctx, db.HasReversedUPIReferenceParams{SocietyID: arg.SocietyID, Reference: arg.Reference})
	return value, persistenceError(err)
}

func (r *collectionsRepository) InsertUPILedger(ctx context.Context, arg contracts.InsertUPILedgerInput) error {
	return persistenceError(GetQueries(ctx, r.database).InsertUPILedger(ctx, db.InsertUPILedgerParams{SocietyID: arg.SocietyID, BillID: arg.BillID, PaymentID: arg.PaymentID, Kind: arg.Kind, AmountPaise: arg.AmountPaise, ActorID: arg.ActorID}))
}

func (r *collectionsRepository) InsertUPIPayment(ctx context.Context, arg contracts.InsertUPIPaymentInput) (contracts.MaintenancePayment, error) {
	value, err := GetQueries(ctx, r.database).InsertUPIPayment(ctx, db.InsertUPIPaymentParams{SocietyID: arg.SocietyID, BillID: arg.BillID, ClaimID: arg.ClaimID, PayerID: arg.PayerID, SettingsVersion: arg.SettingsVersion, AmountPaise: arg.AmountPaise, Reference: arg.Reference, CreditDate: pgtype.Date{Time: arg.CreditDate, Valid: !arg.CreditDate.IsZero()}, ReceiptNumber: pgtype.UUID{Bytes: arg.ReceiptNumber, Valid: arg.ReceiptNumber != uuid.Nil}, VerifiedBy: arg.VerifiedBy, EvidenceReference: arg.EvidenceReference})
	return mapFinancialMaintenancePayment(value), persistenceError(err)
}

func (r *collectionsRepository) ListUPIPayments(ctx context.Context, arg contracts.ListUPIPaymentsInput) ([]contracts.MaintenancePayment, error) {
	value, err := GetQueries(ctx, r.database).ListUPIPayments(ctx, db.ListUPIPaymentsParams{SocietyID: arg.SocietyID, Status: arg.Status, Reference: arg.Reference, BeforeID: arg.BeforeID, BillID: arg.BillID, FlatID: arg.FlatID, BillingMonth: pgtype.Date{Time: arg.BillingMonth, Valid: !arg.BillingMonth.IsZero()}, Block: arg.Block, FlatNumber: arg.FlatNumber, Search: arg.Search, UserID: arg.UserID, Limit: arg.Limit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.MaintenancePayment, len(value))
	for i := range value {
		out[i] = mapFinancialMaintenancePayment(value[i])
	}
	return out, nil
}

func (r *collectionsRepository) LockMaintenancePaymentBill(ctx context.Context, arg contracts.LockMaintenancePaymentBillInput) (contracts.MaintenanceBill, error) {
	value, err := GetQueries(ctx, r.database).LockMaintenancePaymentBill(ctx, db.LockMaintenancePaymentBillParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapFinancialMaintenanceBill(value), persistenceError(err)
}

func (r *collectionsRepository) ReleaseUPIPaymentReference(ctx context.Context, arg contracts.ReleaseUPIPaymentReferenceInput) error {
	return persistenceError(GetQueries(ctx, r.database).ReleaseUPIPaymentReference(ctx, db.ReleaseUPIPaymentReferenceParams{SocietyID: arg.SocietyID, PaymentID: arg.PaymentID}))
}

func (r *collectionsRepository) ReserveUPIReference(ctx context.Context, arg contracts.ReserveUPIReferenceInput) error {
	return persistenceError(GetQueries(ctx, r.database).ReserveUPIReference(ctx, db.ReserveUPIReferenceParams{SocietyID: arg.SocietyID, Reference: arg.Reference, BillID: arg.BillID, ClaimID: arg.ClaimID, PaymentID: arg.PaymentID}))
}

func (r *collectionsRepository) ReverseUPIPayment(ctx context.Context, arg contracts.ReverseUPIPaymentInput) (contracts.MaintenancePayment, error) {
	value, err := GetQueries(ctx, r.database).ReverseUPIPayment(ctx, db.ReverseUPIPaymentParams{SocietyID: arg.SocietyID, ID: arg.ID, ReversedBy: arg.ReversedBy, ReversalReason: arg.ReversalReason})
	return mapFinancialMaintenancePayment(value), persistenceError(err)
}

func (r *collectionsRepository) UPICollectionSummary(ctx context.Context, arg contracts.UPICollectionSummaryInput) (contracts.UPICollectionSummaryRecord, error) {
	value, err := GetQueries(ctx, r.database).UPICollectionSummary(ctx, db.UPICollectionSummaryParams{SocietyID: arg.SocietyID, FlatID: arg.FlatID, Month: pgtype.Date{Time: arg.Month, Valid: !arg.Month.IsZero()}})
	return mapFinancialUPICollectionSummaryRow(value), persistenceError(err)
}

var _ contracts.CollectionsRepository = (*collectionsRepository)(nil)
