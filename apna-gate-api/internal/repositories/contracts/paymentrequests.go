package contracts

import (
	"context"
)

type PaymentRequestsRepository interface {
	CloseUPIRequests(ctx context.Context, arg CloseUPIRequestsInput) error
	GetActiveUPIRequest(ctx context.Context, arg GetActiveUPIRequestInput) (MaintenancePaymentRequest, error)
	GetUPIRequest(ctx context.Context, arg GetUPIRequestInput) (MaintenancePaymentRequest, error)
	InsertUPIRequest(ctx context.Context, arg InsertUPIRequestInput) (MaintenancePaymentRequest, error)
	SupersedeUPIRequests(ctx context.Context, societyID int64) error
}
