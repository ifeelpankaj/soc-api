package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	maintenancesvc "go-server/internal/services/maintenanceSvc"
	"go-server/pkg/utils"
	"strconv"
)

type maintenanceService interface {
	Outstanding(context.Context, int64, int64, int64) (models.MaintenanceOutstanding, error)
	OutstandingFlats(context.Context, models.MaintenanceOutstandingFlatFilter) (models.MaintenanceOutstandingFlatList, error)
	BillingRun(context.Context, int64, int64, string) (models.MaintenanceRunStatus, error)
	Settings(context.Context, int64, int64) (models.MaintenanceSettings, error)
	Configure(context.Context, int64, int64, string, models.MaintenanceSettings) (models.MaintenanceSettings, error)
	PreviewCommand(context.Context, int64, int64, models.MaintenanceMonthRequest) (models.MaintenancePreview, error)
	GenerateCommand(context.Context, int64, int64, string, models.MaintenanceMonthRequest) (models.MaintenanceRunResult, error)
	List(context.Context, models.MaintenanceBillFilter) (models.MaintenanceBillList, error)
	Get(context.Context, models.MaintenanceBillFilter) (models.MaintenanceBill, error)
}
type MaintenanceHandler struct{ service maintenanceService }

func NewMaintenanceHandler(s maintenanceService) *MaintenanceHandler { return &MaintenanceHandler{s} }
func maintenanceIDs(c *gin.Context) (int64, int64, bool) {
	id, ok := parsePathInt64(c, "societyId")
	if !ok {
		return 0, 0, false
	}
	user, ok := currentUserID(c)
	return id, user, ok
}
func maintenanceResult(c *gin.Context, result any, err error) {
	var issues *maintenancesvc.ValidationIssues
	if errors.As(err, &issues) {
		utils.ErrorResponseWithDetails(c, 422, "MAINTENANCE_FLAT_DATA_REQUIRED", issues.Error(), gin.H{"issues": issues.Issues})
		return
	}
	if handleServiceError(c, err) {
		return
	}
	utils.SuccessResponse(c, 200, "Maintenance request completed", result)
}

// GetSettings godoc
// @Summary Get maintenance settings and eligibility policy
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.MaintenanceSettingsResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/settings [get]
func (h *MaintenanceHandler) GetSettings(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	result, err := h.service.Settings(c.Request.Context(), id, user)
	maintenanceResult(c, result, err)
}

// PutSettings godoc
// @Summary Replace maintenance configuration (owner/admin)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param Idempotency-Key header string true "Unique retry key"
// @Param request body models.MaintenanceSettings true "Configuration; eligible_flat_statuses is read-only"
// @Success 200 {object} models.MaintenanceSettingsResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/settings [put]
func (h *MaintenanceHandler) PutSettings(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.MaintenanceSettings
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.Configure(c.Request.Context(), id, user, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// PreviewBills godoc
// @Summary Preview missing current-month bills or an eligible unissued past month
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param request body models.MaintenanceMonthRequest true "Billing month YYYY-MM; catch_up requires acknowledgment and review token on generation"
// @Success 200 {object} models.MaintenancePreviewResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/bills/preview [post]
func (h *MaintenanceHandler) PreviewBills(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.MaintenanceMonthRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.PreviewCommand(c.Request.Context(), id, user, req)
	maintenanceResult(c, result, err)
}

// GenerateBills godoc
// @Summary Generate missing current-month bills or a reviewed unissued past month
// @Description Fresh current-month requests add only missing flats using original issued-month terms. Reusing an Idempotency-Key replays its response. Missing flat data returns 422; no partial additional batch is committed. Existing bills and the original run snapshot remain unchanged.
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param societyId path int true "Society ID"
// @Param Idempotency-Key header string true "Unique retry key"
// @Param request body models.MaintenanceMonthRequest true "Billing month YYYY-MM; catch_up requires acknowledgment and review token on generation"
// @Success 200 {object} models.MaintenanceRunResponse
// @Failure 422 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/bills/generate [post]
func (h *MaintenanceHandler) GenerateBills(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.MaintenanceMonthRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.GenerateCommand(c.Request.Context(), id, user, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

func (h *MaintenanceHandler) read(c *gin.Context, resident, detail bool) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	f := models.MaintenanceBillFilter{SocietyID: id, UserID: user, Resident: resident, Limit: 25}
	if detail {
		f.ID, ok = parsePathInt64(c, "id")
		if !ok {
			return
		}
		result, err := h.service.Get(c.Request.Context(), f)
		maintenanceResult(c, result, err)
		return
	}
	for key, dest := range map[string]*int64{"flat_id": &f.FlatID, "cursor": &f.BeforeID} {
		if value := c.Query(key); value != "" {
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil || n <= 0 {
				utils.BadRequestResponse(c, key+" must be a positive integer")
				return
			}
			*dest = n
		}
	}
	if value := c.Query("limit"); value != "" {
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil || n < 1 || n > 100 {
			utils.BadRequestResponse(c, "limit must be 1 to 100")
			return
		}
		f.Limit = int32(n)
	}
	f.Month = c.Query("billing_month")
	f.Status = c.Query("status")
	f.DisplayStatus = c.Query("payment_status")
	if value := c.Query("page"); value != "" {
		page, err := strconv.ParseInt(value, 10, 32)
		if err != nil || page < 1 || page > 100000 || f.BeforeID != 0 {
			utils.BadRequestResponse(c, "page must be 1 to 100000 and cannot be combined with cursor")
			return
		}
		f.Page = int32(page)
	}
	f.BillNumber = c.Query("bill_number")
	flatQuery, ok := maintenanceFlatQuery(c)
	if !ok {
		return
	}
	f.Block = flatQuery.Block
	f.FlatNumber = flatQuery.FlatNumber
	f.Search = flatQuery.Search
	result, err := h.service.List(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// ListBills godoc
// @Summary List society maintenance bills (owner/admin)
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param billing_month query string false "YYYY-MM"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block (partial match)"
// @Param flat_number query string false "Flat number (partial match)"
// @Param search query string false "Search block or flat number"
// @Param status query string false "unpaid or overdue"
// @Param payment_status query string false "unpaid, overdue, paid, pending_review, rejected or outstanding"
// @Param page query int false "1-based page; do not combine with cursor"
// @Param cursor query int false "Before bill ID"
// @Param limit query int false "1 to 100; default 25"
// @Success 200 {object} models.MaintenanceBillListResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/bills [get]
func (h *MaintenanceHandler) ListBills(c *gin.Context) { h.read(c, false, false) }

// GetBill godoc
// @Summary Get society bill details (owner/admin)
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Bill ID"
// @Success 200 {object} models.MaintenanceBillResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/bills/{id} [get]
func (h *MaintenanceHandler) GetBill(c *gin.Context) { h.read(c, false, true) }

// MyBills godoc
// @Summary List bills for the authenticated resident's active flat residencies
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param billing_month query string false "YYYY-MM"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block (partial match)"
// @Param flat_number query string false "Flat number (partial match)"
// @Param search query string false "Search block or flat number"
// @Param status query string false "unpaid or overdue"
// @Param payment_status query string false "unpaid, overdue, paid, pending_review, rejected or outstanding"
// @Param page query int false "1-based page; do not combine with cursor"
// @Param cursor query int false "Before bill ID"
// @Param limit query int false "1 to 100; default 25"
// @Success 200 {object} models.MaintenanceBillListResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/bills [get]
func (h *MaintenanceHandler) MyBills(c *gin.Context) { h.read(c, true, false) }

// MyBill godoc
// @Summary Get an accessible flat bill without historical billed-party personal details
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Bill ID"
// @Success 200 {object} models.MaintenanceBillResponse
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/bills/{id} [get]
func (h *MaintenanceHandler) MyBill(c *gin.Context) { h.read(c, true, true) }

// BillingRun godoc
// @Summary Get monthly billing run status, including completed empty runs
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param billing_month query string true "YYYY-MM"
// @Success 200 {object} models.APIResponse{data=models.MaintenanceRunStatus}
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/billing-runs [get]
func (h *MaintenanceHandler) BillingRun(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	result, err := h.service.BillingRun(c.Request.Context(), sid, user, c.Query("billing_month"))
	maintenanceResult(c, result, err)
}

// Outstanding godoc
// @Summary Outstanding maintenance across independent monthly bills
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param flatId path int true "Flat ID"
// @Success 200 {object} models.APIResponse{data=models.MaintenanceOutstanding}
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/flats/{flatId}/outstanding [get]
func (h *MaintenanceHandler) Outstanding(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	flat, ok := parsePathInt64(c, "flatId")
	if !ok {
		return
	}
	result, err := h.service.Outstanding(c.Request.Context(), id, user, flat)
	maintenanceResult(c, result, err)
}

// ListOutstandingFlats godoc
// @Summary List flats with outstanding maintenance (owner/admin)
// @Tags Maintenance
// @Produce json
// @Param societyId path int true "Society ID"
// @Param block query string false "Block"
// @Param flat_number query string false "Flat number"
// @Param search query string false "Search block or flat number"
// @Param cursor query int false "Before flat ID"
// @Param limit query int false "1 to 100; default 25"
// @Success 200 {object} models.APIResponse{data=models.MaintenanceOutstandingFlatList}
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/outstanding-flats [get]
func (h *MaintenanceHandler) ListOutstandingFlats(c *gin.Context) {
	id, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	f := models.MaintenanceOutstandingFlatFilter{SocietyID: id, UserID: user, Limit: 25}
	if v := c.Query("cursor"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			utils.BadRequestResponse(c, "cursor must be positive")
			return
		}
		f.BeforeID = n
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n < 1 || n > 100 {
			utils.BadRequestResponse(c, "limit must be 1 to 100")
			return
		}
		f.Limit = int32(n)
	}
	flatQuery, ok := maintenanceFlatQuery(c)
	if !ok {
		return
	}
	f.Block = flatQuery.Block
	f.FlatNumber = flatQuery.FlatNumber
	f.Search = flatQuery.Search
	result, err := h.service.OutstandingFlats(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}
