package handlers

import (
	"context"
	"go-server/internal/models"
	"go-server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type maintenancePayments interface {
	Settings(context.Context, int64, int64) (models.MaintenancePaymentSettings, error)
	SaveSettings(context.Context, int64, int64, string, models.MaintenancePaymentSettings) (models.MaintenancePaymentSettings, error)
	SettingsHistory(context.Context, models.UPIListFilter) (models.UPISettingsPage, error)
	RequestCommand(context.Context, int64, int64, int64, string) (models.UPIPaymentRequest, error)
	QR(context.Context, int64, int64, string) ([]byte, error)
	SubmitClaim(context.Context, int64, int64, int64, string, models.UPISubmitClaim) (models.UPIClaim, error)
	CloseClaim(context.Context, int64, int64, int64, string, models.UPIReason, bool) (models.UPIClaim, error)
	Claims(context.Context, models.UPIListFilter) (models.UPIClaimsPage, error)
	Claim(context.Context, int64, int64, int64) (models.UPIClaimDetail, error)
	Verify(context.Context, int64, int64, int64, string, models.UPIVerifyCredit) (models.UPIPayment, error)
	Record(context.Context, int64, int64, string, models.UPIDirectPayment) (models.UPIPayment, error)
	Reverse(context.Context, int64, int64, int64, string, models.UPIReason) (models.UPIPayment, error)
	Payments(context.Context, models.UPIListFilter) (models.UPIPaymentsPage, error)
	Payment(context.Context, int64, int64, int64, bool) (models.UPIPayment, error)
	CreateReport(context.Context, int64, int64, string, models.UPICreateReport, bool) (models.UPIReport, error)
	Reports(context.Context, models.UPIListFilter) (models.UPIReportsPage, error)
	Report(context.Context, int64, int64, int64) (models.UPIReportDetail, error)
	UpdateReport(context.Context, int64, int64, int64, string, models.UPIUpdateReport) (models.UPIReport, error)
	ReferenceHistory(context.Context, int64, int64, string) (models.UPIReferenceHistory, error)
	Summary(context.Context, int64, int64, string, ...int64) (models.UPICollectionSummary, error)
	Audit(context.Context, models.UPIListFilter, int64) (models.UPIAuditPage, error)
}
type MaintenancePaymentHandler struct{ service maintenancePayments }

func NewMaintenancePaymentHandler(s maintenancePayments) *MaintenancePaymentHandler {
	return &MaintenancePaymentHandler{s}
}
func paymentFilter(c *gin.Context, resident bool) (models.UPIListFilter, bool) {
	society, user, ok := maintenanceIDs(c)
	if !ok {
		return models.UPIListFilter{}, false
	}
	f := models.UPIListFilter{SocietyID: society, UserID: user, Resident: resident, Limit: 25, Status: c.Query("status"), Reference: c.Query("reference")}
	if v := c.Query("cursor"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			utils.BadRequestResponse(c, "cursor must be positive")
			return f, false
		}
		f.BeforeID = n
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n < 1 || n > 100 {
			utils.BadRequestResponse(c, "limit must be 1 to 100")
			return f, false
		}
		f.Limit = int32(n)
	}

	for name, target := range map[string]*int64{"bill_id": &f.BillID, "flat_id": &f.FlatID} {
		if raw := c.Query(name); raw != "" {
			n, err := strconv.ParseInt(raw, 10, 64)
			if err != nil || n <= 0 {
				utils.BadRequestResponse(c, "Filters require positive IDs")
				return f, false
			}
			*target = n
		}
	}
	f.BillingMonth = c.Query("billing_month")
	flatQuery, ok := maintenanceFlatQuery(c)
	if !ok {
		return f, false
	}
	f.Block = flatQuery.Block
	f.FlatNumber = flatQuery.FlatNumber
	f.Search = flatQuery.Search
	return f, true
}

// GetPaymentSettings godoc
// @Summary GetPaymentSettings
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Success 200 {object} models.APIResponse{data=models.MaintenancePaymentSettings}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-settings [get]
func (h *MaintenancePaymentHandler) GetPaymentSettings(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	result, err := h.service.Settings(c.Request.Context(), sid, user)
	maintenanceResult(c, result, err)
}

// PutPaymentSettings godoc
// @Summary PutPaymentSettings
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.MaintenancePaymentSettings true "Request"
// @Success 200 {object} models.APIResponse{data=models.MaintenancePaymentSettings}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-settings [put]
func (h *MaintenancePaymentHandler) PutPaymentSettings(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.MaintenancePaymentSettings
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.SaveSettings(c.Request.Context(), sid, user, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// PaymentSettingsHistory godoc
// @Summary PaymentSettingsHistory
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPISettingsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-settings/versions [get]
func (h *MaintenancePaymentHandler) PaymentSettingsHistory(c *gin.Context) {
	f, ok := paymentFilter(c, false)
	if !ok {
		return
	}
	result, err := h.service.SettingsHistory(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// CreatePaymentRequest godoc
// @Summary CreatePaymentRequest
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIPaymentRequest}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/bills/{id}/payment-request [post]
func (h *MaintenancePaymentHandler) CreatePaymentRequest(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.RequestCommand(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"))
	maintenanceResult(c, result, err)
}

// SubmitPaymentClaim godoc
// @Summary SubmitPaymentClaim
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPISubmitClaim true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIClaim}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/bills/{id}/payment-claims [post]
func (h *MaintenancePaymentHandler) SubmitPaymentClaim(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPISubmitClaim
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.SubmitClaim(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// MyPaymentClaims godoc
// @Summary MyPaymentClaims
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIClaimsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payment-claims [get]
func (h *MaintenancePaymentHandler) MyPaymentClaims(c *gin.Context) {
	f, ok := paymentFilter(c, true)
	if !ok {
		return
	}
	result, err := h.service.Claims(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// CancelPaymentClaim godoc
// @Summary CancelPaymentClaim
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIReason true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIClaim}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payment-claims/{id}/cancel [post]
func (h *MaintenancePaymentHandler) CancelPaymentClaim(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPIReason
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.CloseClaim(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req, false)
	maintenanceResult(c, result, err)
}

// PaymentClaims godoc
// @Summary PaymentClaims
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIClaimsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-claims [get]
// @Param bill_id query int false "Bill ID"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block (partial match)"
// @Param flat_number query string false "Flat number (partial match)"
// @Param search query string false "Search block or flat number"
// @Param billing_month query string false "YYYY-MM"
func (h *MaintenancePaymentHandler) PaymentClaims(c *gin.Context) {
	f, ok := paymentFilter(c, false)
	if !ok {
		return
	}
	result, err := h.service.Claims(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// PaymentClaim godoc
// @Summary PaymentClaim
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIClaimDetail}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-claims/{id} [get]
func (h *MaintenancePaymentHandler) PaymentClaim(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.Claim(c.Request.Context(), sid, user, id)
	maintenanceResult(c, result, err)
}

// VerifyPaymentClaim godoc
// @Summary VerifyPaymentClaim
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIVerifyCredit true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-claims/{id}/verify [post]
func (h *MaintenancePaymentHandler) VerifyPaymentClaim(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPIVerifyCredit
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.Verify(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// RejectPaymentClaim godoc
// @Summary RejectPaymentClaim
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIReason true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIClaim}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-claims/{id}/reject [post]
func (h *MaintenancePaymentHandler) RejectPaymentClaim(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPIReason
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.CloseClaim(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req, true)
	maintenanceResult(c, result, err)
}

// RecordPayment godoc
// @Summary RecordPayment
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIDirectPayment true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments/manual [post]
func (h *MaintenancePaymentHandler) RecordPayment(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.UPIDirectPayment
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.Record(c.Request.Context(), sid, user, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// Payments godoc
// @Summary Payments
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIPaymentsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments [get]
// @Param bill_id query int false "Bill ID"
// @Param flat_id query int false "Flat ID"
// @Param block query string false "Block (partial match)"
// @Param flat_number query string false "Flat number (partial match)"
// @Param search query string false "Search block or flat number"
// @Param billing_month query string false "YYYY-MM"
func (h *MaintenancePaymentHandler) Payments(c *gin.Context) {
	f, ok := paymentFilter(c, false)
	if !ok {
		return
	}
	result, err := h.service.Payments(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// Payment godoc
// @Summary Payment
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments/{id} [get]
func (h *MaintenancePaymentHandler) Payment(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.Payment(c.Request.Context(), sid, user, id, false)
	maintenanceResult(c, result, err)
}

// ReversePayment godoc
// @Summary ReversePayment
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIReason true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments/{id}/reverse [post]
func (h *MaintenancePaymentHandler) ReversePayment(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPIReason
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.Reverse(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// MyPayments godoc
// @Summary MyPayments
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIPaymentsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payments [get]
func (h *MaintenancePaymentHandler) MyPayments(c *gin.Context) {
	f, ok := paymentFilter(c, true)
	if !ok {
		return
	}
	result, err := h.service.Payments(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// MyPaymentReceipt godoc
// @Summary MyPaymentReceipt
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payments/{id}/receipt [get]
func (h *MaintenancePaymentHandler) MyPaymentReceipt(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.Payment(c.Request.Context(), sid, user, id, true)
	maintenanceResult(c, result, err)
}

// PaymentReceipt godoc
// @Summary PaymentReceipt
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIPayment}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments/{id}/receipt [get]
func (h *MaintenancePaymentHandler) PaymentReceipt(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.Payment(c.Request.Context(), sid, user, id, false)
	maintenanceResult(c, result, err)
}

// CreateMyPaymentReport godoc
// @Summary CreateMyPaymentReport
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPICreateReport true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIReport}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payment-reports [post]
func (h *MaintenancePaymentHandler) CreateMyPaymentReport(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.UPICreateReport
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.CreateReport(c.Request.Context(), sid, user, c.GetHeader("Idempotency-Key"), req, false)
	maintenanceResult(c, result, err)
}

// MyPaymentReports godoc
// @Summary MyPaymentReports
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIReportsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payment-reports [get]
func (h *MaintenancePaymentHandler) MyPaymentReports(c *gin.Context) {
	f, ok := paymentFilter(c, true)
	if !ok {
		return
	}
	result, err := h.service.Reports(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// CreatePaymentReport godoc
// @Summary CreatePaymentReport
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPICreateReport true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIReport}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-reports [post]
func (h *MaintenancePaymentHandler) CreatePaymentReport(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	var req models.UPICreateReport
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.CreateReport(c.Request.Context(), sid, user, c.GetHeader("Idempotency-Key"), req, true)
	maintenanceResult(c, result, err)
}

// PaymentReports godoc
// @Summary PaymentReports
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param cursor query int false "Before ID (version for settings history)"
// @Param limit query int false "1 to 100"
// @Param status query string false "Status filter"
// @Param reference query string false "Reference filter for admin claims/payments"
// @Success 200 {object} models.APIResponse{data=models.UPIReportsPage}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-reports [get]
func (h *MaintenancePaymentHandler) PaymentReports(c *gin.Context) {
	f, ok := paymentFilter(c, false)
	if !ok {
		return
	}
	result, err := h.service.Reports(c.Request.Context(), f)
	maintenanceResult(c, result, err)
}

// PaymentReport godoc
// @Summary PaymentReport
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Success 200 {object} models.APIResponse{data=models.UPIReportDetail}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-reports/{id} [get]
func (h *MaintenancePaymentHandler) PaymentReport(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	result, err := h.service.Report(c.Request.Context(), sid, user, id)
	maintenanceResult(c, result, err)
}

// UpdatePaymentReport godoc
// @Summary UpdatePaymentReport
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param id path int true "Resource ID"
// @Accept json
// @Param Idempotency-Key header string true "Unique retry key; reuse only with identical input"
// @Param request body models.UPIUpdateReport true "Request"
// @Success 200 {object} models.APIResponse{data=models.UPIReport}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-reports/{id} [patch]
func (h *MaintenancePaymentHandler) UpdatePaymentReport(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var req models.UPIUpdateReport
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.service.UpdateReport(c.Request.Context(), sid, user, id, c.GetHeader("Idempotency-Key"), req)
	maintenanceResult(c, result, err)
}

// PaymentReferenceHistory godoc
// @Summary PaymentReferenceHistory
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param reference query string true "UTR or bank reference"
// @Success 200 {object} models.APIResponse{data=models.UPIReferenceHistory}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-reference-history [get]
func (h *MaintenancePaymentHandler) PaymentReferenceHistory(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	result, err := h.service.ReferenceHistory(c.Request.Context(), sid, user, c.Query("reference"))
	maintenanceResult(c, result, err)
}

// CollectionSummary godoc
// @Summary CollectionSummary
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param billing_month query string false "YYYY-MM"
// @Param flat_id query int false "Flat ID"
// @Success 200 {object} models.APIResponse{data=models.UPICollectionSummary}
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/reports/summary [get]
func (h *MaintenancePaymentHandler) CollectionSummary(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	flatID := int64(0)
	if raw := c.Query("flat_id"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || n <= 0 {
			utils.BadRequestResponse(c, "flat_id must be positive")
			return
		}
		flatID = n
	}
	result, err := h.service.Summary(c.Request.Context(), sid, user, c.Query("billing_month"), flatID)
	maintenanceResult(c, result, err)
}

// PaymentQR godoc
// @Summary Get active UPI request QR as PNG
// @Tags Maintenance Payments
// @Produce png
// @Param societyId path int true "Society ID"
// @Param requestId path string true "Payment request UUID"
// @Success 200 {file} binary
// @Failure 400,401,403,404,409 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payment-requests/{requestId}/qr [get]
func (h *MaintenancePaymentHandler) PaymentQR(c *gin.Context) {
	sid, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	data, err := h.service.QR(c.Request.Context(), sid, user, c.Param("requestId"))
	if handleServiceError(c, err) {
		return
	}
	c.Header("Cache-Control", "private, no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(200, "image/png", data)
}

// PaymentAudit godoc
// @Summary Read payment audit trail (owner/admin)
// @Tags Maintenance Payments
// @Produce json
// @Param societyId path int true "Society ID"
// @Param bill_id query int false "Bill ID"
// @Param cursor query int false "Before event ID"
// @Param limit query int false "1 to 100"
// @Success 200 {object} models.APIResponse{data=models.UPIAuditPage}
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payment-audit [get]
func (h *MaintenancePaymentHandler) PaymentAudit(c *gin.Context) {
	f, ok := paymentFilter(c, false)
	if !ok {
		return
	}
	var bill int64
	if v := c.Query("bill_id"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			utils.BadRequestResponse(c, "bill_id must be positive")
			return
		}
		bill = n
	}
	result, err := h.service.Audit(c.Request.Context(), f, bill)
	maintenanceResult(c, result, err)
}
