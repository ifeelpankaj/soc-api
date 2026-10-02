package maintenancesvc

import (
	"context"
	"go-server/internal/models"
)

type Payments interface {
	PaymentsCommands
	PaymentsQueries
}

type PaymentsCommands interface {
	SubmitClaim(ctx context.Context, society, user, bill int64, key string, req models.UPISubmitClaim) (models.UPIClaim, error)
	CloseClaim(ctx context.Context, society, user, id int64, key string, req models.UPIReason, admin bool) (models.UPIClaim, error)
	Verify(ctx context.Context, society, user, claimID int64, key string, req models.UPIVerifyCredit) (models.UPIPayment, error)
	Record(ctx context.Context, society, user int64, key string, req models.UPIDirectPayment) (models.UPIPayment, error)
	Reverse(ctx context.Context, society, user, id int64, key string, req models.UPIReason) (models.UPIPayment, error)
	CreateReport(ctx context.Context, society, user int64, key string, req models.UPICreateReport, admin bool) (models.UPIReport, error)
	UpdateReport(ctx context.Context, society, user, id int64, key string, req models.UPIUpdateReport) (models.UPIReport, error)
	SaveSettings(ctx context.Context, society, user int64, key string, v models.MaintenancePaymentSettings) (models.MaintenancePaymentSettings, error)
	Request(ctx context.Context, society, user, bill int64) (models.UPIPaymentRequest, error)
	RequestCommand(ctx context.Context, society, user, bill int64, key string) (models.UPIPaymentRequest, error)
}

type PaymentsQueries interface {
	SettingsHistory(ctx context.Context, f models.UPIListFilter) (models.UPISettingsPage, error)
	Claims(ctx context.Context, f models.UPIListFilter) (models.UPIClaimsPage, error)
	Payments(ctx context.Context, f models.UPIListFilter) (models.UPIPaymentsPage, error)
	Payment(ctx context.Context, society, user, id int64, resident bool) (models.UPIPayment, error)
	ReferenceHistory(ctx context.Context, society, user int64, reference string) (models.UPIReferenceHistory, error)
	Claim(ctx context.Context, society, user, id int64) (models.UPIClaimDetail, error)
	Summary(ctx context.Context, society, user int64, month string, flatIDs ...int64) (models.UPICollectionSummary, error)
	Reports(ctx context.Context, f models.UPIListFilter) (models.UPIReportsPage, error)
	Report(ctx context.Context, society, user, id int64) (models.UPIReportDetail, error)
	Audit(ctx context.Context, f models.UPIListFilter, bill int64) (models.UPIAuditPage, error)
	Settings(ctx context.Context, society, user int64) (models.MaintenancePaymentSettings, error)
	QR(ctx context.Context, society, user int64, id string) ([]byte, error)
}

var _ Payments = (*PaymentService)(nil)
