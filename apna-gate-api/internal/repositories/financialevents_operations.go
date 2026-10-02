package repository

import (
	"context"
	"go-server/internal/db"
	"go-server/internal/repositories/contracts"
	"go-server/pkg/database"
)

type financialEventsRepository struct{ database *database.Database }

func (r *financialEventsRepository) EnqueueUPIEvent(ctx context.Context, arg contracts.EnqueueUPIEventInput) error {
	return persistenceError(GetQueries(ctx, r.database).EnqueueUPIEvent(ctx, db.EnqueueUPIEventParams{BillID: arg.BillID, EventType: arg.EventType, EventKey: arg.EventKey, Audience: arg.Audience, EventData: arg.EventData, SocietyID: arg.SocietyID}))
}

func (r *financialEventsRepository) InsertUPIAudit(ctx context.Context, arg contracts.InsertUPIAuditInput) error {
	return persistenceError(GetQueries(ctx, r.database).InsertUPIAudit(ctx, db.InsertUPIAuditParams{SocietyID: arg.SocietyID, BillID: arg.BillID, ActorID: arg.ActorID, Action: arg.Action, EntityID: arg.EntityID, Details: arg.Details}))
}

func (r *financialEventsRepository) ListUPIAudit(ctx context.Context, arg contracts.ListUPIAuditInput) ([]contracts.MaintenancePaymentAuditEvent, error) {
	value, err := GetQueries(ctx, r.database).ListUPIAudit(ctx, db.ListUPIAuditParams{SocietyID: arg.SocietyID, BillID: arg.BillID, BeforeID: arg.BeforeID, Limit: arg.Limit})
	if err != nil {
		return nil, persistenceError(err)
	}
	out := make([]contracts.MaintenancePaymentAuditEvent, len(value))
	for i := range value {
		out[i] = mapFinancialMaintenancePaymentAuditEvent(value[i])
	}
	return out, nil
}

var _ contracts.FinancialEventsRepository = (*financialEventsRepository)(nil)
