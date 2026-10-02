package routes

import (
	handlers "go-server/internal/handlers/v1"
	middleware "go-server/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupInternalJobRoutes(rg *gin.RouterGroup, handler *handlers.JobWebhookHandler, secret string) {
	internalJobs := rg.Group("/jobs")
	internalJobs.Use(middleware.JobWebhookAuth(secret))
	{
		internalJobs.POST("/db-backup", handler.TriggerBackup)
		internalJobs.GET("/db-backup/:run_id", handler.GetBackup)
		internalJobs.POST("/maintenance-billing", handler.TriggerMaintenanceBilling)
		internalJobs.POST("/maintenance-reminders", handler.TriggerMaintenanceReminders)
		internalJobs.POST("/cleanup", handler.TriggerCleanup)
		internalJobs.POST("/monthly-visitor-report", handler.TriggerMonthlyVisitorReport)
	}
}
