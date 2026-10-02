package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type paymentReportsRepository struct{ database *database.Database }

func (r *paymentReportsRepository) GetUPIReport(ctx context.Context, arg contracts.GetUPIReportInput) (contracts.GetUPIReportRecord, error) {
	value, err := GetQueries(ctx, r.database).GetUPIReport(ctx, db.GetUPIReportParams{SocietyID: arg.SocietyID, ID: arg.ID})
	return mapFinancialGetUPIReportRow(value), persistenceError(err)
}

func (r *paymentReportsRepository) GetUPIReportByFingerprint(ctx context.Context, arg contracts.GetUPIReportByFingerprintInput) (contracts.MaintenancePaymentReport, error) {
	value, err := GetQueries(ctx, r.database).GetUPIReportByFingerprint(ctx, db.GetUPIReportByFingerprintParams{SocietyID: arg.SocietyID, ReportedBy: arg.ReportedBy, Fingerprint: arg.Fingerprint})
	return mapFinancialMaintenancePaymentReport(value), persistenceError(err)
}

func (r *paymentReportsRepository) InsertUPIReport(ctx context.Context, arg contracts.InsertUPIReportInput) (contracts.MaintenancePaymentReport, error) {
	value, err := GetQueries(ctx, r.database).InsertUPIReport(ctx, db.InsertUPIReportParams{SocietyID: arg.SocietyID, BillID: arg.BillID, RequestID: pgtype.UUID{Bytes: arg.RequestID, Valid: arg.RequestID != uuid.Nil}, ReportedBy: arg.ReportedBy, Reference: arg.Reference, AmountPaise: arg.AmountPaise, PaymentDate: pgtype.Date{Time: arg.PaymentDate, Valid: !arg.PaymentDate.IsZero()}, Explanation: arg.Explanation, Fingerprint: arg.Fingerprint})
	return mapFinancialMaintenancePaymentReport(value), persistenceError(err)
}

func (r *paymentReportsRepository) ListUPIReports(ctx context.Context, arg contracts.ListUPIReportsInput) ([]contracts.ListUPIReportsRecord, error) {
	value, err := GetQueries(ctx, r.database).ListUPIReports(ctx, db.ListUPIReportsParams{SocietyID: arg.SocietyID, UserID: arg.UserID, Status: arg.Status, BeforeID: arg.BeforeID, Limit: arg.Limit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.ListUPIReportsRecord, len(value))
	for i := range value {
		out[i] = mapFinancialListUPIReportsRow(value[i])
	}
	return out, nil
}

func (r *paymentReportsRepository) UpdateUPIReport(ctx context.Context, arg contracts.UpdateUPIReportInput) (contracts.MaintenancePaymentReport, error) {
	value, err := GetQueries(ctx, r.database).UpdateUPIReport(ctx, db.UpdateUPIReportParams{SocietyID: arg.SocietyID, ID: arg.ID, Status: arg.Status, ResolutionNote: arg.ResolutionNote, UpdatedBy: arg.UpdatedBy})
	return mapFinancialMaintenancePaymentReport(value), persistenceError(err)
}

var _ contracts.PaymentReportsRepository = (*paymentReportsRepository)(nil)
