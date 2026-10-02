package routes

import (
	"go-server/internal/app"
	"go-server/internal/middlewares/guards"
	"go-server/internal/models"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutesV1(rg *gin.RouterGroup, h *app.V1Handlers, g *guards.Guards) {
	ownerAdminRoles := []string{
		string(models.SocietyMemberRoleOwner),
		string(models.SocietyMemberRoleAdmin),
	}
	ownerAdminStaffRoles := []string{
		string(models.SocietyMemberRoleOwner),
		string(models.SocietyMemberRoleAdmin),
		string(models.SocietyMemberRoleStaff),
	}

	ownerOnly := rg.Group("/societies/:societyId")
	ownerOnly.Use(g.SocietyRolesForParam("societyId", string(models.SocietyMemberRoleOwner))...)
	{
		ownerOnly.DELETE("", h.Society.DeleteSociety)
	}

	adminSetup := rg.Group("/societies/:societyId")
	adminSetup.Use(g.SocietyMembershipRolesForParam("societyId", ownerAdminRoles...)...)
	{
		adminSetup.GET("/onboarding/bootstrap", h.Society.GetOnboardingBootstrap)
	}

	admin := rg.Group("/societies/:societyId")
	admin.Use(g.SocietyRolesForParam("societyId", ownerAdminRoles...)...)
	admin.Use(g.OperationalSocietyForParam("societyId")...)
	{
		admin.GET("/maintenance/payment-settings", h.MaintenancePayment.GetPaymentSettings)
		admin.GET("/maintenance/bills/:id/invoice", h.MaintenanceDocument.InvoicePDF)
		admin.GET("/maintenance/payments/:id/receipt/pdf", h.MaintenanceDocument.ReceiptPDF)
		admin.PUT("/maintenance/payment-settings", h.MaintenancePayment.PutPaymentSettings)
		admin.GET("/maintenance/payment-settings/versions", h.MaintenancePayment.PaymentSettingsHistory)
		admin.GET("/maintenance/payment-claims", h.MaintenancePayment.PaymentClaims)
		admin.GET("/maintenance/payment-claims/:id", h.MaintenancePayment.PaymentClaim)
		admin.POST("/maintenance/payment-claims/:id/verify", h.MaintenancePayment.VerifyPaymentClaim)
		admin.POST("/maintenance/payment-claims/:id/reject", h.MaintenancePayment.RejectPaymentClaim)
		admin.POST("/maintenance/payments/manual", h.MaintenancePayment.RecordPayment)
		admin.GET("/maintenance/payments", h.MaintenancePayment.Payments)
		admin.GET("/maintenance/payments/:id", h.MaintenancePayment.Payment)
		admin.POST("/maintenance/payments/:id/reverse", h.MaintenancePayment.ReversePayment)
		admin.GET("/maintenance/payments/:id/receipt", h.MaintenancePayment.PaymentReceipt)
		admin.POST("/maintenance/payment-reports", h.MaintenancePayment.CreatePaymentReport)
		admin.GET("/maintenance/payment-reports", h.MaintenancePayment.PaymentReports)
		admin.GET("/maintenance/payment-reports/:id", h.MaintenancePayment.PaymentReport)
		admin.PATCH("/maintenance/payment-reports/:id", h.MaintenancePayment.UpdatePaymentReport)
		admin.GET("/maintenance/payment-reference-history", h.MaintenancePayment.PaymentReferenceHistory)
		admin.GET("/maintenance/reports/summary", h.MaintenancePayment.CollectionSummary)
		admin.GET("/maintenance/payment-audit", h.MaintenancePayment.PaymentAudit)
		admin.GET("/maintenance/billing-runs", h.Maintenance.BillingRun)
		admin.GET("/maintenance/settings", h.Maintenance.GetSettings)
		admin.PUT("/maintenance/settings", h.Maintenance.PutSettings)
		admin.POST("/maintenance/bills/preview", h.Maintenance.PreviewBills)
		admin.POST("/maintenance/bills/generate", h.Maintenance.GenerateBills)
		admin.GET("/maintenance/bills", h.Maintenance.ListBills)
		admin.GET("/maintenance/bills/:id", h.Maintenance.GetBill)
		admin.GET("/maintenance/outstanding-flats", h.Maintenance.ListOutstandingFlats)
		admin.PATCH("", h.Society.UpdateSociety)
		admin.GET("/visitor-settings", h.VisitorSetting.GetSocietySettings)
		admin.PATCH("/visitor-settings", h.VisitorSetting.UpdateSocietySettings)
		admin.GET("/visitor-settings/flats", h.VisitorSetting.ListSocietyFlatSettings)
		admin.GET("/members/:memberId/visitor-approval-stats", h.VisitorEntry.GetMemberVisitorApprovalStats)
		admin.POST("/members", h.Society.AddMember)
		admin.PATCH("/members/:userId/role", h.Society.ChangeMemberRole)
		admin.POST("/members/:userId/suspend", h.Society.SuspendMember)
		admin.POST("/members/:userId/reactivate", h.Society.ReactivateMember)
		admin.DELETE("/members/:userId", h.Society.RemoveMember)
		admin.POST("/transfer-ownership", h.Society.TransferOwnership)
		admin.GET("/members/summary", h.Society.GetSocietyMemberSummary)
		admin.GET("/members/search", h.Society.ListSocietyMembers)
		admin.GET("/members/:memberId", h.Society.GetSocietyMember)
		admin.GET("/members", h.Society.ListSocietyMembers)
		admin.GET("/dashboard/bootstrap", h.Society.GetDashboardBootstrap)

	}

	staffRead := rg.Group("/societies/:societyId")
	staffRead.Use(g.SocietyRolesForParam("societyId", ownerAdminStaffRoles...)...)
	staffRead.Use(g.OperationalSocietyForParam("societyId")...)
	{
		staffRead.GET("", h.Society.GetSociety)
		staffRead.GET("/flats", h.Flat.ListSocietyFlats)
		staffRead.GET("/flats/stats", h.Flat.GetFlatStats)
		staffRead.GET("/flats/:flatId", h.Flat.GetFlat)
		staffRead.GET("/flats/:flatId/residents", h.Flat.ListSocietyFlatResidents)
		staffRead.GET("/flats/:flatId/residents/:residentId", h.Flat.GetFlatResident)
		staffRead.GET("/flat-claims/search", h.Flat.ListSocietyFlatClaims)
		staffRead.GET("/flat-claims/stats", h.Flat.GetFlatClaimStats)
		staffRead.GET("/flat-claims", h.Flat.ListSocietyFlatClaims)
		staffRead.GET("/flat-claims/:claimId", h.Flat.GetSocietyFlatClaim)
		staffRead.POST("/visitor-entries/guard", h.VisitorEntry.CreateGuardEntry)
		staffRead.POST("/flats/:flatId/visitor-invites/staff", h.VisitorEntry.CreateStaffInvite)
		staffRead.GET("/visitor-invites/:token", h.VisitorEntry.GetGuardInviteByToken)
		staffRead.POST("/visitor-entries/check-in", h.VisitorEntry.CheckIn)
		staffRead.GET("/visitor-entries/stats", h.VisitorEntry.GetEntryStats)
		staffRead.GET("/visitor-entries/stats/daily", h.VisitorEntry.GetEntryDailyStats)
		staffRead.GET("/guard-desk/bootstrap", h.VisitorEntry.GetGuardDeskBootstrap)
		staffRead.GET("/visitor-entries/pending", h.VisitorEntry.ListSocietyPendingApprovals)
		staffRead.GET("/visitor-entries/waiting-at-gate", h.VisitorEntry.ListWaitingAtGate)
		staffRead.GET("/visitor-entries/expected-guests", h.VisitorEntry.ListExpectedGuests)
		staffRead.GET("/visitor-entries", h.VisitorEntry.ListEntries)
		staffRead.GET("/visitor-entries/:entryId", h.VisitorEntry.GetEntry)
		staffRead.PATCH("/visitor-entries/:entryId", h.VisitorEntry.UpdateGuardEntry)
		staffRead.GET("/visitor-entries/:entryId/events", h.VisitorEntry.ListEvents)
		staffRead.POST("/visitor-entries/:entryId/notify", h.VisitorEntry.NotifyPendingEntry)
		staffRead.POST("/visitor-entries/:entryId/guard-approve", h.VisitorEntry.GuardApproveEntry)
		staffRead.POST("/visitor-entries/:entryId/approve-and-check-in", h.VisitorEntry.GuardApproveAndCheckIn)
		staffRead.POST("/visitor-entries/:entryId/check-in", h.VisitorEntry.CheckInByEntryID)
		staffRead.POST("/visitor-entries/:entryId/check-out", h.VisitorEntry.CheckOut)
	}

	operationalAdmin := rg.Group("/societies/:societyId")
	operationalAdmin.Use(g.SocietyRolesForParam("societyId", ownerAdminRoles...)...)
	operationalAdmin.Use(g.SocietyOperationalForParam("societyId")...)
	{
		operationalAdmin.POST("/guards", h.Society.CreateGuard)

		operationalAdmin.POST("/flats", h.Flat.CreateFlat)
		operationalAdmin.POST("/flats/bulk", h.Flat.BulkCreateFlats)
		operationalAdmin.POST("/flats/generate", h.Flat.GenerateFlats)
		operationalAdmin.PATCH("/flats/:flatId", h.Flat.UpdateFlat)
		operationalAdmin.DELETE("/flats/:flatId", h.Flat.DeleteFlat)
		operationalAdmin.POST("/flats/:flatId/block", h.Flat.BlockFlat)
		operationalAdmin.POST("/flats/:flatId/unblock", h.Flat.UnblockFlat)

		operationalAdmin.POST("/flat-claims/:claimId/approve", h.Flat.ApproveFlatClaim)
		operationalAdmin.POST("/flat-claims/:claimId/reject", h.Flat.RejectFlatClaim)

		operationalAdmin.POST("/flats/:flatId/residents/users/:userId", h.Flat.AddResidentToFlat)
		operationalAdmin.DELETE("/flats/:flatId/residents/:residentId", h.Flat.RemoveResidentFromFlat)
		operationalAdmin.POST("/flats/:flatId/residents/:residentId/primary", h.Flat.ChangePrimaryResident)
		operationalAdmin.PATCH("/flats/:flatId/residents/:residentId/role", h.Flat.UpdateFlatResidentRole)
		operationalAdmin.POST("/flats/:flatId/residents/:residentId/move-out", h.Flat.MoveOutResident)
	}
}
