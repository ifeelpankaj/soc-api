package repository

import (
	"context"

	"go-server/pkg/database"
)

// MaintenancePaymentRepository implements financial persistence capabilities.
type MaintenancePaymentRepository struct {
	*financialRepository
	tx TransactionManager
}

func NewMaintenancePaymentRepository(d *database.Database, tx TransactionManager) *MaintenancePaymentRepository {
	return &MaintenancePaymentRepository{newFinancialRepository(d), tx}
}
func (r *MaintenancePaymentRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.tx.WithTransaction(ctx, fn)
}
