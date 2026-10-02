package contracts

import (
	"context"
)

type TransactionManager interface {
	// WithTransaction reuses a transaction already present in ctx. Only the
	// outermost invocation begins, commits, or rolls back. Return nested errors
	// to abort the outer transaction; a nil outer result commits it.
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
