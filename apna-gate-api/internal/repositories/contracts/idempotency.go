package contracts

import (
	"context"
)

type IdempotencyRepository interface {
	GetUPIIdempotency(ctx context.Context, arg GetUPIIdempotencyInput) (MaintenancePaymentIdempotency, error)
	InsertUPIIdempotency(ctx context.Context, arg InsertUPIIdempotencyInput) error
}
