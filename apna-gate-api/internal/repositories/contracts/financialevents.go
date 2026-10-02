package contracts

import (
	"context"
)

type FinancialEventsRepository interface {
	EnqueueUPIEvent(ctx context.Context, arg EnqueueUPIEventInput) error
	InsertUPIAudit(ctx context.Context, arg InsertUPIAuditInput) error
	ListUPIAudit(ctx context.Context, arg ListUPIAuditInput) ([]MaintenancePaymentAuditEvent, error)
}
