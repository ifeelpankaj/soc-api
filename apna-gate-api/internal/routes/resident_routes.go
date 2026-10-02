package routes

import (
	"go-server/internal/app"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupResidentRoutesV1(rg *gin.RouterGroup, h *app.V1Handlers, g *guards.Guards) {
	plans := rg.Group("/plans")
	{
		plans.GET("", h.Plan.ListPlans)
		plans.GET("/lookup", h.Plan.GetPlan)
	}

	resident := rg.Group("")
	resident.Use(g.Authenticated()...)
	{
		resident.POST("/societies", h.Society.CreateSocietyRequest)
		resident.GET("/societies/my", h.Society.ListMySocieties)

		resident.POST("/flat-claims", h.Flat.SubmitFlatClaim)
		resident.POST("/flat-claims/:claimId/cancel", h.Flat.CancelMyFlatClaim)
		resident.GET("/me/flat-claims", h.Flat.ListMyFlatClaims)
		resident.GET("/me/residences", h.Flat.ListMyResidences)
	}

	residentSociety := rg.Group("")
	residentSociety.Use(g.Authenticated()...)
	residentSociety.Use(g.OperationalSocietyForParam("societyId")...)
	{
		residentSociety.GET("/societies/:societyId/maintenance/flats/:flatId/outstanding", h.Maintenance.Outstanding)
		residentSociety.POST("/societies/:societyId/maintenance/my/bills/:id/payment-request", h.MaintenancePayment.CreatePaymentRequest)
		residentSociety.GET("/societies/:societyId/maintenance/my/bills/:id/invoice", h.MaintenanceDocument.MyInvoicePDF)
		residentSociety.GET("/societies/:societyId/maintenance/my/payments/:id/receipt/pdf", h.MaintenanceDocument.MyReceiptPDF)
		residentSociety.POST("/societies/:societyId/maintenance/my/bills/:id/payment-claims", h.MaintenancePayment.SubmitPaymentClaim)
		residentSociety.GET("/societies/:societyId/maintenance/my/payment-claims", h.MaintenancePayment.MyPaymentClaims)
		residentSociety.POST("/societies/:societyId/maintenance/my/payment-claims/:id/cancel", h.MaintenancePayment.CancelPaymentClaim)
		residentSociety.GET("/societies/:societyId/maintenance/my/payments", h.MaintenancePayment.MyPayments)
		residentSociety.GET("/societies/:societyId/maintenance/my/payments/:id/receipt", h.MaintenancePayment.MyPaymentReceipt)
		residentSociety.POST("/societies/:societyId/maintenance/my/payment-reports", h.MaintenancePayment.CreateMyPaymentReport)
		residentSociety.GET("/societies/:societyId/maintenance/my/payment-reports", h.MaintenancePayment.MyPaymentReports)
		residentSociety.GET("/societies/:societyId/maintenance/my/payment-requests/:requestId/qr", h.MaintenancePayment.PaymentQR)
		residentSociety.GET("/societies/:societyId/maintenance/my/bills", h.Maintenance.MyBills)
		residentSociety.GET("/societies/:societyId/maintenance/my/bills/:id", h.Maintenance.MyBill)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-settings", h.VisitorSetting.GetFlatSettings)
		residentSociety.PATCH("/societies/:societyId/flats/:flatId/visitor-settings/:purpose", h.VisitorSetting.UpdateFlatPurposeSetting)
		residentSociety.POST("/societies/:societyId/flats/:flatId/visitor-settings/reset", h.VisitorSetting.ResetFlatSettingsToDefault)
		residentSociety.GET("/societies/:societyId/flats/:flatId/members", h.MemberInvite.ListFlatResidentsForResident)
		residentSociety.GET("/societies/:societyId/flats/:flatId/member-invites", h.MemberInvite.ListPendingMemberInvites)
		residentSociety.GET("/societies/:societyId/flats/:flatId/member-invites/:inviteId", h.MemberInvite.GetMemberInviteHistory)
		residentSociety.POST("/societies/:societyId/flats/:flatId/member-invites", h.MemberInvite.CreateMemberInvite)
		residentSociety.POST("/societies/:societyId/flats/:flatId/member-invites/:inviteId/cancel", h.MemberInvite.CancelMemberInvite)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-invites", h.VisitorEntry.ListVisitorInviteHistory)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-invites/:inviteId", h.VisitorEntry.GetVisitorInviteHistory)
		residentSociety.POST("/societies/:societyId/flats/:flatId/visitor-invites", h.VisitorEntry.CreateInvite)
		residentSociety.POST("/societies/:societyId/visitor-invites/:inviteId/cancel", h.VisitorEntry.CancelInvite)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-context", h.VisitorEntry.GetFlatVisitorContextForResident)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-entries/pending", h.VisitorEntry.ListPendingApprovals)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-entries", h.VisitorEntry.ListFlatVisitorEntries)
		residentSociety.GET("/societies/:societyId/flats/:flatId/visitor-entries/:entryId", h.VisitorEntry.GetFlatVisitorEntry)
		residentSociety.POST("/societies/:societyId/visitor-entries/:entryId/approve", h.VisitorEntry.ApproveEntry)
		residentSociety.POST("/societies/:societyId/visitor-entries/:entryId/reject", h.VisitorEntry.RejectEntry)
	}
}
