package contracts

import (
	"context"
)

type CollectionsRepository interface {
	GetActiveUPIPayment(ctx context.Context, arg GetActiveUPIPaymentInput) (MaintenancePayment, error)
	GetPaymentBill(ctx context.Context, arg GetPaymentBillInput) (MaintenanceBill, error)
	GetUPIPayment(ctx context.Context, arg GetUPIPaymentInput) (MaintenancePayment, error)
	HasReversedUPIReference(ctx context.Context, arg HasReversedUPIReferenceInput) (bool, error)
	InsertUPILedger(ctx context.Context, arg InsertUPILedgerInput) error
	InsertUPIPayment(ctx context.Context, arg InsertUPIPaymentInput) (MaintenancePayment, error)
	ListUPIPayments(ctx context.Context, arg ListUPIPaymentsInput) ([]MaintenancePayment, error)
	LockMaintenancePaymentBill(ctx context.Context, arg LockMaintenancePaymentBillInput) (MaintenanceBill, error)
	ReleaseUPIPaymentReference(ctx context.Context, arg ReleaseUPIPaymentReferenceInput) error
	ReserveUPIReference(ctx context.Context, arg ReserveUPIReferenceInput) error
	ReverseUPIPayment(ctx context.Context, arg ReverseUPIPaymentInput) (MaintenancePayment, error)
	UPICollectionSummary(ctx context.Context, arg UPICollectionSummaryInput) (UPICollectionSummaryRecord, error)
}
