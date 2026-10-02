package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-server/internal/models"
	"mime"
	"net/http"
)

type maintenanceDocuments interface {
	Invoice(context.Context, int64, int64, int64, bool) (models.MaintenancePDF, error)
	Receipt(context.Context, int64, int64, int64, bool) (models.MaintenancePDF, error)
}
type MaintenanceDocumentHandler struct{ service maintenanceDocuments }

func NewMaintenanceDocumentHandler(s maintenanceDocuments) *MaintenanceDocumentHandler {
	return &MaintenanceDocumentHandler{s}
}
func (h *MaintenanceDocumentHandler) download(c *gin.Context, resident, receipt bool) {
	society, user, ok := maintenanceIDs(c)
	if !ok {
		return
	}
	id, ok := parsePathInt64(c, "id")
	if !ok {
		return
	}
	var pdf models.MaintenancePDF
	var err error
	if receipt {
		pdf, err = h.service.Receipt(c.Request.Context(), society, user, id, resident)
	} else {
		pdf, err = h.service.Invoice(c.Request.Context(), society, user, id, resident)
	}
	if handleServiceError(c, err) {
		return
	}
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": pdf.Filename}))
	c.Header("Cache-Control", "private, no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "application/pdf", pdf.Bytes)
}

// InvoicePDF godoc
// @Summary Download an issued maintenance invoice (owner/admin)
// @Tags Maintenance Documents
// @Produce application/pdf
// @Param societyId path int true "Society ID"
// @Param id path int true "Bill ID"
// @Success 200 {file} binary
// @Failure 400,401,403,404,422,500 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/bills/{id}/invoice [get]
func (h *MaintenanceDocumentHandler) InvoicePDF(c *gin.Context) { h.download(c, false, false) }

// MyInvoicePDF godoc
// @Summary Download an invoice for a currently accessible flat
// @Tags Maintenance Documents
// @Produce application/pdf
// @Param societyId path int true "Society ID"
// @Param id path int true "Bill ID"
// @Success 200 {file} binary
// @Failure 400,401,403,404,422,500 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/bills/{id}/invoice [get]
func (h *MaintenanceDocumentHandler) MyInvoicePDF(c *gin.Context) { h.download(c, true, false) }

// ReceiptPDF godoc
// @Summary Download a verified or reversed payment receipt (owner/admin)
// @Tags Maintenance Documents
// @Produce application/pdf
// @Param societyId path int true "Society ID"
// @Param id path int true "Payment ID, not a claim or bill ID"
// @Success 200 {file} binary
// @Failure 400,401,403,404,422,500 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/payments/{id}/receipt/pdf [get]
func (h *MaintenanceDocumentHandler) ReceiptPDF(c *gin.Context) { h.download(c, false, true) }

// MyReceiptPDF godoc
// @Summary Download a receipt with the same privacy restrictions as the JSON receipt
// @Tags Maintenance Documents
// @Produce application/pdf
// @Param societyId path int true "Society ID"
// @Param id path int true "Payment ID, not a claim or bill ID"
// @Success 200 {file} binary
// @Failure 400,401,403,404,422,500 {object} models.ErrorResponseDoc
// @Security AccessToken
// @Router /v1/societies/{societyId}/maintenance/my/payments/{id}/receipt/pdf [get]
func (h *MaintenanceDocumentHandler) MyReceiptPDF(c *gin.Context) { h.download(c, true, true) }
