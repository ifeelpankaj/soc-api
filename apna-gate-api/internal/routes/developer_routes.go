package routes

import (
	"go-server/internal/app"
	"go-server/internal/middlewares/guards"

	"github.com/gin-gonic/gin"
)

func SetupDeveloperRoutesV1(rg *gin.RouterGroup, h *app.V1Handlers, g *guards.Guards) {
	developer := rg.Group("")
	developer.Use(g.Developer()...)
	{
		developer.GET("/developer/dashboard/bootstrap", h.Society.GetDeveloperDashboardBootstrap)

		developer.GET("/societies/search", h.Society.ListSocieties)
		developer.GET("/societies", h.Society.ListSocieties)
		developer.GET("/societies/:societyId/allmember", h.Society.ListAllSocietyMember)
		developer.POST("/societies/:societyId/approve", h.Society.ApproveSociety)
		developer.POST("/societies/:societyId/reject", h.Society.RejectSociety)
		developer.POST("/societies/:societyId/suspend", h.Society.SuspendSociety)
		developer.POST("/societies/:societyId/reactivate", h.Society.ReactivateSociety)
		developer.POST("/societies/:societyId/restore", h.Society.RestoreSociety)

		developer.POST("/plans", h.Plan.CreatePlan)
		developer.PATCH("/plans/:planId", h.Plan.UpdatePlan)
		developer.POST("/plans/:planId/activate", h.Plan.ActivatePlan)
		developer.POST("/plans/:planId/deactivate", h.Plan.DeactivatePlan)

		developer.POST("/societies/:societyId/subscriptions/plans/:planId/pending", h.Subscription.CreatePendingSubscription)
		developer.POST("/societies/:societyId/subscriptions/plans/:planId/trial", h.Subscription.CreateTrialSubscription)
		developer.POST("/subscriptions/:subscriptionId/activate", h.Subscription.ActivateSubscription)
		developer.POST("/subscriptions/:subscriptionId/renew", h.Subscription.RenewSubscription)
		developer.POST("/subscriptions/:subscriptionId/cancel", h.Subscription.CancelSubscription)
		developer.POST("/subscriptions/:subscriptionId/expire", h.Subscription.ExpireSubscription)
		developer.POST("/subscriptions/:subscriptionId/plans/:planId", h.Subscription.ChangeSubscriptionPlan)
		developer.GET("/subscriptions/lookup", h.Subscription.GetSubscription)
		developer.GET("/subscriptions/search", h.Subscription.ListSubscriptions)
		developer.GET("/subscriptions", h.Subscription.ListSubscriptions)
		developer.GET("/subscriptions/stats", h.Subscription.GetSubscriptionStats)

		developer.GET("/flats", h.Flat.ListFlats)
		developer.GET("/flat-claims", h.Flat.ListFlatClaims)
		developer.GET("/flat-claims/:claimId", h.Flat.GetFlatClaim)
		developer.GET("/flat-residents/search", h.Flat.ListFlatResidents)
		developer.GET("/flat-residents", h.Flat.ListFlatResidents)
	}
}
