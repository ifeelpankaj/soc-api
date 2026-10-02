package repository

import (
	"context"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

// HubRepository shares the transaction context used by society and notification repositories.
type HubRepository struct {
	database     *database.Database
	transactions TransactionManager
}

func NewHubRepository(database *database.Database, transactions TransactionManager) *HubRepository {
	return &HubRepository{database, transactions}
}
func (r *HubRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.transactions.WithTransaction(ctx, fn)
}

var _ contracts.HubStore = (*HubRepository)(nil)
