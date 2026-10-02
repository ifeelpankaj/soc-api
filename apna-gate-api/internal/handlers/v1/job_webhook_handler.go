package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go-server/internal/backup"
	"io"
	"net/http"

	"go-server/internal/jobs"
	"go-server/internal/models"
	"go-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type jobTriggerer interface {
	Trigger(name string) (jobs.Trigger, error)
}

type JobWebhookHandler struct {
	manager jobTriggerer
	backup  backupService
}

type backupService interface {
	Trigger(context.Context, string) (*backup.Run, error)
	Get(context.Context, string) (*backup.Run, error)
}

func NewJobWebhookHandler(manager jobTriggerer, backups ...backupService) *JobWebhookHandler {
	h := &JobWebhookHandler{manager: manager}
	if len(backups) > 0 {
		h.backup = backups[0]
	}
	return h
}

func (h *JobWebhookHandler) TriggerBackup(c *gin.Context) {
	// Decode an exact one-field object: also reject null, duplicate keys and trailing JSON.
	dec := json.NewDecoder(http.MaxBytesReader(c.Writer, c.Request.Body, 1024))
	token, err := dec.Token()
	if err != nil || token != json.Delim('{') {
		c.JSON(400, gin.H{"error": "Expected an object containing only type: daily or weekly"})
		return
	}
	var kind string
	seen := false
	for dec.More() {
		key, err := dec.Token()
		if err != nil || key != "type" || seen {
			c.JSON(400, gin.H{"error": "Only one type field is allowed"})
			return
		}
		if err := dec.Decode(&kind); err != nil {
			c.JSON(400, gin.H{"error": "Invalid backup type"})
			return
		}
		seen = true
	}
	end, err := dec.Token()
	var trailing any
	if err != nil || end != json.Delim('}') || dec.Decode(&trailing) != io.EOF || !seen || (kind != "daily" && kind != "weekly") {
		c.JSON(400, gin.H{"error": "Type must be daily or weekly"})
		return
	}
	if h.backup == nil {
		c.JSON(503, gin.H{"error": "Database backups are disabled"})
		return
	}
	r, err := h.backup.Trigger(c.Request.Context(), kind)
	if err != nil {
		if errors.Is(err, backup.ErrDisabled) || errors.Is(err, jobs.ErrJobManagerStopped) {
			c.JSON(503, gin.H{"error": "Database backups are unavailable"})
			return
		}
		c.JSON(500, gin.H{"error": "Could not persist backup request"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"run_id": r.ID, "type": r.Type, "status": r.Status})
}
func (h *JobWebhookHandler) GetBackup(c *gin.Context) {
	id := c.Param("run_id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid run ID"})
		return
	}
	if h.backup == nil {
		c.JSON(503, gin.H{"error": "Database backups are unavailable"})
		return
	}
	r, err := h.backup.Get(c.Request.Context(), id)
	if errors.Is(err, backup.ErrNotFound) {
		c.JSON(404, gin.H{"error": "Backup run not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not read backup status"})
		return
	}
	c.JSON(200, r)
}

func (h *JobWebhookHandler) TriggerCleanup(c *gin.Context) {
	h.trigger(c, jobs.JobCleanup)
}

func (h *JobWebhookHandler) TriggerMonthlyVisitorReport(c *gin.Context) {
	h.trigger(c, jobs.JobMonthlyVisitorReport)
}

func (h *JobWebhookHandler) trigger(c *gin.Context, name string) {
	trigger, err := h.manager.Trigger(name)
	if err != nil {
		if errors.Is(err, jobs.ErrJobManagerStopped) {
			utils.ErrorResponse(c, http.StatusServiceUnavailable, models.ErrCodeServiceUnavailable, "Job manager is shutting down", nil)
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusAccepted, trigger)
}

func (h *JobWebhookHandler) TriggerMaintenanceBilling(c *gin.Context) {
	h.trigger(c, jobs.JobMaintenanceBilling)
}

// TriggerMaintenanceReminders godoc
// @Summary Queue and deliver eligible unpaid maintenance reminders
// @Tags Internal Jobs
// @Param Authorization header string true "Bearer JOB_WEBHOOK_SECRET"
// @Success 202 {object} models.JobTriggerResponse
// @Failure 401 {object} models.APIResponse
// @Router /internal/jobs/maintenance-reminders [post]
func (h *JobWebhookHandler) TriggerMaintenanceReminders(c *gin.Context) {
	h.trigger(c, jobs.JobMaintenanceReminders)
}
