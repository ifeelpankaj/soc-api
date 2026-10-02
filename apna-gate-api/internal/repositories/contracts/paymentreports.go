package contracts

import (
	"context"
)

type PaymentReportsRepository interface {
	GetUPIReport(ctx context.Context, arg GetUPIReportInput) (GetUPIReportRecord, error)
	GetUPIReportByFingerprint(ctx context.Context, arg GetUPIReportByFingerprintInput) (MaintenancePaymentReport, error)
	InsertUPIReport(ctx context.Context, arg InsertUPIReportInput) (MaintenancePaymentReport, error)
	ListUPIReports(ctx context.Context, arg ListUPIReportsInput) ([]ListUPIReportsRecord, error)
	UpdateUPIReport(ctx context.Context, arg UpdateUPIReportInput) (MaintenancePaymentReport, error)
}
