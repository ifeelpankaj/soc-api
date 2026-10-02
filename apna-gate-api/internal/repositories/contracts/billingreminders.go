package contracts

import (
	"context"
)

type BillingRemindersRepository interface {
	EnqueueMaintenanceReminder(ctx context.Context, arg EnqueueMaintenanceReminderInput) (int64, error)
	ListMaintenanceReminderCandidates(ctx context.Context, arg ListMaintenanceReminderCandidatesInput) ([]ListMaintenanceReminderCandidatesRecord, error)
}
