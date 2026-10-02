package maintenancesvc

import (
	"context"
	"go-server/internal/models"
)

type BillingService interface {
	BillingServiceCommands
	BillingServiceQueries
	BillingServiceJobs
}

type BillingServiceCommands interface {
	Configure(ctx context.Context, society, user int64, key string, v models.MaintenanceSettings) (models.MaintenanceSettings, error)
	GenerateCommand(ctx context.Context, id, user int64, key string, req models.MaintenanceMonthRequest) (models.MaintenanceRunResult, error)
	SaveSettings(ctx context.Context, id, user int64, v models.MaintenanceSettings) (models.MaintenanceSettings, error)
	Generate(ctx context.Context, id, user int64, month string) (models.MaintenanceRunResult, error)
}

type BillingServiceQueries interface {
	PreviewCommand(ctx context.Context, id, user int64, req models.MaintenanceMonthRequest) (models.MaintenancePreview, error)
	Outstanding(ctx context.Context, society, user, flat int64) (models.MaintenanceOutstanding, error)
	OutstandingFlats(ctx context.Context, f models.MaintenanceOutstandingFlatFilter) (models.MaintenanceOutstandingFlatList, error)
	Settings(ctx context.Context, id, user int64) (models.MaintenanceSettings, error)
	Preview(ctx context.Context, id, user int64, month string) (models.MaintenancePreview, error)
	List(ctx context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBillList, error)
	Get(ctx context.Context, f models.MaintenanceBillFilter) (models.MaintenanceBill, error)
	BillingRun(ctx context.Context, society, user int64, month string) (models.MaintenanceRunStatus, error)
}

type BillingServiceJobs interface {
	RunScheduled(ctx context.Context) error
	Deliver(ctx context.Context) error
	RunReminders(ctx context.Context) error
}

var _ BillingService = (*Service)(nil)
