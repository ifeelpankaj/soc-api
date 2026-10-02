package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type billingRemindersRepository struct{ database *database.Database }

func (r *billingRemindersRepository) EnqueueMaintenanceReminder(ctx context.Context, arg contracts.EnqueueMaintenanceReminderInput) (int64, error) {
	value, err := GetQueries(ctx, r.database).EnqueueMaintenanceReminder(ctx, db.EnqueueMaintenanceReminderParams{EventKey: arg.EventKey, EventData: arg.EventData, BillID: arg.BillID, SocietyID: arg.SocietyID})
	return value, persistenceError(err)
}

func (r *billingRemindersRepository) ListMaintenanceReminderCandidates(ctx context.Context, arg contracts.ListMaintenanceReminderCandidatesInput) ([]contracts.ListMaintenanceReminderCandidatesRecord, error) {
	value, err := GetQueries(ctx, r.database).ListMaintenanceReminderCandidates(ctx, db.ListMaintenanceReminderCandidatesParams{SocietyID: arg.SocietyID, AfterID: arg.AfterID, AsOf: timestampParam(arg.AsOf)})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.ListMaintenanceReminderCandidatesRecord, len(value))
	for i := range value {
		out[i] = mapFinancialListMaintenanceReminderCandidatesRow(value[i])
	}
	return out, nil
}

var _ contracts.BillingRemindersRepository = (*billingRemindersRepository)(nil)
